/*
Copyright © 2025 Chris Griffis <dev@chrisgriffis.com> and contributors.

All rights reserved.
Licensed under the MIT license. See LICENSE file in the project root for details.
*/

package bartender

import (
	decimal "github.com/alpacahq/alpacadecimal"
)

func WithTickThreshold(threshold int64) Option[TickBarConfig] {
	return func(t *TickBarConfig) {
		t.tickThreshold = decimal.NewFromInt(threshold)
	}
}

type TickBarConfig struct {
	tickThreshold decimal.Decimal `validate:"required"`
}

func (c TickBarConfig) Process(trades <-chan Trade) chan *Bar {
	output := make(chan *Bar)

	go func() {
		defer close(output)

		var current *Bar
		var tradeCount decimal.Decimal

		for trade := range trades {
			// initialize the current bar if it doesn't exist
			if current == nil {
				current = &Bar{}
			}

			// check if the trade is on a new day
			if !current.Start.IsZero() && current.Start.Weekday() != trade.Time.Weekday() {
				finalizedBar := current

				output <- finalizedBar

				// reset the current bar
				current = &Bar{}
				tradeCount = decimal.Zero
			}

			current.applyTrade(trade)

			// increment counter
			tradeCount = tradeCount.Add(decimal.NewFromInt(1))

			if tradeCount.GreaterThanOrEqual(c.tickThreshold) {
				finalizedBar := current
				output <- finalizedBar

				current = nil
				tradeCount = decimal.Zero
			}
		}

		if current != nil {
			output <- current
		}
	}()

	return output
}

func WithTickImbalanceThreshold(threshold int64) Option[TickImbalanceBarConfig] {
	return func(t *TickImbalanceBarConfig) {
		t.imbalanceThreshold = decimal.NewFromInt(threshold)
	}
}

type TickImbalanceBarConfig struct {
	imbalanceThreshold decimal.Decimal `validate:"required"`
}

func (c TickImbalanceBarConfig) Process(trades <-chan Trade) chan *Bar {
	output := make(chan *Bar)

	go func() {
		defer close(output)

		var current *Bar
		var netImbalance decimal.Decimal
		var prevPrice decimal.Decimal

		for trade := range trades {
			// initialize the current bar if it doesn't exist
			if current == nil {
				current = &Bar{}
			}

			// initialize the previous price if it doesn't exist
			if prevPrice.IsZero() {
				prevPrice = trade.Price

				// set first imbalance value based on side of first trade
				if trade.Side == SideBuy {
					netImbalance = decimal.NewFromInt(1)
				} else {
					netImbalance = decimal.NewFromInt(-1)
				}
			}

			// check if the trade is on a new day
			if !current.Start.IsZero() && current.Start.Weekday() != trade.Time.Weekday() {
				finalizedBar := current

				output <- finalizedBar

				// reset the current bar
				current = &Bar{}
				netImbalance = decimal.Zero
			}

			current.applyTrade(trade)

			// update net imbalance
			if trade.Price.GreaterThan(prevPrice) {
				netImbalance = netImbalance.Add(decimal.NewFromInt(1))
			} else if trade.Price.LessThan(prevPrice) {
				netImbalance = netImbalance.Sub(decimal.NewFromInt(1))
			}

			prevPrice = trade.Price // update the previous price

			if netImbalance.Abs().GreaterThanOrEqual(c.imbalanceThreshold) {
				finalizedBar := current
				output <- finalizedBar

				current = nil
				netImbalance = decimal.Zero
			}
		}

		if current != nil {
			output <- current
		}
	}()

	return output
}

func WithTickRunThreshold(threshold int64) Option[TickRunsBarConfig] {
	return func(t *TickRunsBarConfig) {
		t.runsLengthThreshold = decimal.NewFromInt(threshold)
	}
}

type TickRunsBarConfig struct {
	runsLengthThreshold decimal.Decimal `validate:"required"`
}

func (c TickRunsBarConfig) Process(trades <-chan Trade) chan *Bar {
	output := make(chan *Bar)

	go func() {
		defer close(output)

		var current *Bar
		var upwardRun, downwardRun decimal.Decimal
		var prevPrice decimal.Decimal

		for trade := range trades {
			// initialize the current bar if it doesn't exist
			if current == nil {
				current = &Bar{}
			}

			// initialize the last price if not already set
			if prevPrice.IsZero() {
				prevPrice = trade.Price
			}

			// check if the trade is on a new day
			if !current.Start.IsZero() && current.Start.Weekday() != trade.Time.Weekday() {
				finalizedBar := current

				output <- finalizedBar

				// reset the current bar
				current = &Bar{}
				upwardRun = decimal.Zero
				downwardRun = decimal.Zero
			}

			current.applyTrade(trade)

			// determine the direction of the tick and update runs
			if trade.Price.GreaterThan(prevPrice) {
				upwardRun = upwardRun.Add(decimal.NewFromInt(1))
				downwardRun = decimal.Zero
			} else if trade.Price.LessThan(prevPrice) {
				downwardRun = downwardRun.Add(decimal.NewFromInt(1))
				upwardRun = decimal.Zero
			}

			prevPrice = trade.Price // update last price

			// check if a new bar should be created based on the run threshold
			if upwardRun.GreaterThanOrEqual(c.runsLengthThreshold) || downwardRun.GreaterThanOrEqual(c.runsLengthThreshold) {
				finalizedBar := current
				output <- finalizedBar

				current = nil
				upwardRun = decimal.Zero
				downwardRun = decimal.Zero
			}
		}

		if current != nil {
			output <- current
		}
	}()

	return output
}

// Interface guards
var _ Processor = (*TickBarConfig)(nil)
var _ Processor = (*TickImbalanceBarConfig)(nil)
var _ Processor = (*TickRunsBarConfig)(nil)
