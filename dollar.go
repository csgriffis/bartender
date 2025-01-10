/*
Copyright © 2025 Chris Griffis <dev@chrisgriffis.com> and contributors.

All rights reserved.
Licensed under the MIT license. See LICENSE file in the project root for details.
*/

package bartender

import (
	decimal "github.com/alpacahq/alpacadecimal"
)

func WithDollarThreshold(threshold float64) Option[DollarBarConfig] {
	return func(d *DollarBarConfig) {
		d.dollarThreshold = decimal.NewFromFloat(threshold)
	}
}

type DollarBarConfig struct {
	dollarThreshold decimal.Decimal `validate:"required"`
}

func (c DollarBarConfig) Process(trades <-chan Trade) chan *Bar {
	output := make(chan *Bar)

	go func() {
		defer close(output)

		var current *Bar
		var dollar decimal.Decimal

		for trade := range trades {
			if current == nil {
				current = &Bar{
					Open:  trade.Price,
					High:  trade.Price,
					Low:   trade.Price,
					Close: trade.Price,
					Start: trade.Time,
				}
			}

			// update the current bar
			current.Close = trade.Price
			current.High = decimal.Max(current.High, trade.Price)
			current.Low = decimal.Min(current.Low, trade.Price)
			current.Volume = current.Volume.Add(trade.Size)

			// increment tracker
			dollar.Add(trade.Price.Mul(trade.Size))

			if dollar.GreaterThanOrEqual(c.dollarThreshold) {
				finalizedBar := current

				output <- finalizedBar

				// reset the current bar
				current = nil
				// reset the dollar tracker
				dollar = decimal.Zero
			}
		}

		if current != nil {
			output <- current
		}
	}()

	return output
}

func WithDollarImbalanceThreshold(threshold float64) Option[DollarImbalanceBarConfig] {
	return func(d *DollarImbalanceBarConfig) {
		d.imbalanceThreshold = decimal.NewFromFloat(threshold)
	}
}

type DollarImbalanceBarConfig struct {
	imbalanceThreshold decimal.Decimal `validate:"required"`
}

func (c DollarImbalanceBarConfig) Process(trades <-chan Trade) chan *Bar {
	output := make(chan *Bar)

	go func() {
		defer close(output)

		var current *Bar
		var netImbalance decimal.Decimal
		var prevPrice decimal.Decimal

		for trade := range trades {
			if prevPrice.IsZero() {
				prevPrice = trade.Price
			}

			if current == nil {
				current = &Bar{
					Open:  trade.Price,
					High:  trade.Price,
					Low:   trade.Price,
					Close: trade.Price,
					Start: trade.Time,
				}
			}

			// update the current bar
			current.Close = trade.Price
			current.High = decimal.Max(current.High, trade.Price)
			current.Low = decimal.Min(current.Low, trade.Price)
			current.Volume = current.Volume.Add(trade.Size)

			// update net imbalance
			if trade.Price.GreaterThan(prevPrice) {
				netImbalance = netImbalance.Add(trade.Price.Mul(trade.Size))
			} else if trade.Price.LessThan(prevPrice) {
				netImbalance = netImbalance.Sub(trade.Price.Mul(trade.Size))
			}

			prevPrice = trade.Price

			if netImbalance.Abs().GreaterThanOrEqual(c.imbalanceThreshold) {
				finalizedBar := current

				output <- finalizedBar

				// reset the current bar
				current = nil
				// reset the net imbalance
				netImbalance = decimal.Zero
			}
		}

		if current != nil {
			output <- current
		}
	}()

	return output
}

func WithDollarRunThreshold(dollarThreshold float64) Option[DollarRunBarConfig] {
	return func(d *DollarRunBarConfig) {
		d.runDollarThreshold = decimal.NewFromFloat(dollarThreshold)
	}
}

type DollarRunBarConfig struct {
	runDollarThreshold decimal.Decimal `validate:"required,nonzero"`
}

func (c DollarRunBarConfig) Process(trades <-chan Trade) chan *Bar {
	output := make(chan *Bar)

	go func() {
		defer close(output)

		var current *Bar
		var upwardDollarRun, downwardDollarRun decimal.Decimal
		var prevPrice decimal.Decimal

		for trade := range trades {
			// initialize the last price if not already set
			if prevPrice.IsZero() {
				prevPrice = trade.Price
			}

			if current == nil {
				current = &Bar{
					Open:  trade.Price,
					High:  trade.Price,
					Low:   trade.Price,
					Close: trade.Price,
					Start: trade.Time,
				}
			}

			// update the current bar
			current.Close = trade.Price
			current.High = decimal.Max(current.High, trade.Price)
			current.Low = decimal.Min(current.Low, trade.Price)
			current.Volume = current.Volume.Add(trade.Size)

			// calculate the dollar value of the trade (Price * Size)
			tradeDollarValue := trade.Price.Mul(trade.Size)

			// determine the direction of the tick and update dollar runs
			if trade.Price.GreaterThan(prevPrice) {
				upwardDollarRun = upwardDollarRun.Add(tradeDollarValue)
				downwardDollarRun = decimal.Zero
			} else if trade.Price.LessThan(prevPrice) {
				downwardDollarRun = downwardDollarRun.Add(tradeDollarValue)
				upwardDollarRun = decimal.Zero
			}

			prevPrice = trade.Price // update last price

			// check if a new bar should be created based on the dollar run threshold
			if upwardDollarRun.Abs().GreaterThanOrEqual(c.runDollarThreshold) || downwardDollarRun.Abs().GreaterThanOrEqual(c.runDollarThreshold) {
				finalizedBar := current

				output <- finalizedBar

				// reset the current bar
				current = nil
				// reset the dollar runs
				upwardDollarRun = decimal.Zero
				// reset the dollar runs
				downwardDollarRun = decimal.Zero
			}
		}

		if current != nil {
			output <- current
		}
	}()

	return output
}

// Interface guards
var _ Processor = (*DollarBarConfig)(nil)
var _ Processor = (*DollarImbalanceBarConfig)(nil)
var _ Processor = (*DollarRunBarConfig)(nil)
