/*
Copyright © 2025 Chris Griffis <dev@chrisgriffis.com> and contributors.

All rights reserved.
Licensed under the MIT license. See LICENSE file in the project root for details.
*/

package bartender

import (
	decimal "github.com/alpacahq/alpacadecimal"
)

func WithVolumeThreshold(threshold float64) Option[VolumeBarConfig] {
	return func(v *VolumeBarConfig) {
		v.volumeThreshold = decimal.NewFromFloat(threshold)
	}
}

type VolumeBarConfig struct {
	volumeThreshold decimal.Decimal
}

func (c VolumeBarConfig) Process(trades <-chan Trade) chan *Bar {
	output := make(chan *Bar)

	go func() {
		defer close(output)

		var current *Bar

		for trade := range trades {
			// initialize the current bar if it doesn't exist
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

			if current.Volume.GreaterThanOrEqual(c.volumeThreshold) {
				finalizedBar := current
				output <- finalizedBar

				// reset the current bar
				current = nil
			}
		}

		if current != nil {
			output <- current
		}
	}()

	return output
}

func WithVolumeImbalanceThreshold(threshold float64) Option[VolumeImbalanceBarConfig] {
	return func(v *VolumeImbalanceBarConfig) {
		v.imbalanceThreshold = decimal.NewFromFloat(threshold)
	}
}

type VolumeImbalanceBarConfig struct {
	imbalanceThreshold decimal.Decimal
}

func (c VolumeImbalanceBarConfig) Process(trades <-chan Trade) chan *Bar {
	output := make(chan *Bar)

	go func() {
		defer close(output)

		var current *Bar
		var netImbalance decimal.Decimal
		var prevPrice decimal.Decimal

		for trade := range trades {
			// initialize the last price if not already set
			if prevPrice.IsZero() {
				prevPrice = trade.Price
			}

			// initialize the current bar if it doesn't exist
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
				netImbalance = netImbalance.Add(trade.Size)
			} else if trade.Price.LessThan(prevPrice) {
				netImbalance = netImbalance.Sub(trade.Size)
			}

			prevPrice = trade.Price // update the previous price

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

func WithVolumeRunThreshold(threshold float64) Option[VolumeRunsBarConfig] {
	return func(v *VolumeRunsBarConfig) {
		v.runVolumeThreshold = decimal.NewFromFloat(threshold)
	}
}

type VolumeRunsBarConfig struct {
	runVolumeThreshold decimal.Decimal
}

func (c VolumeRunsBarConfig) Process(trades <-chan Trade) chan *Bar {
	output := make(chan *Bar)

	go func() {
		defer close(output)

		var current *Bar
		var upwardVolumeRun, downwardVolumeRun decimal.Decimal
		var prevPrice decimal.Decimal

		for trade := range trades {
			// initialize the last price if not already set
			if prevPrice.IsZero() {
				prevPrice = trade.Price
			}

			// initialize the current bar if it doesn't exist
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

			// determine the direction of the tick and update volume runs
			if trade.Price.GreaterThan(prevPrice) {
				upwardVolumeRun = upwardVolumeRun.Add(trade.Size)
				downwardVolumeRun = decimal.Zero
			} else if trade.Price.LessThan(prevPrice) {
				downwardVolumeRun = downwardVolumeRun.Add(trade.Size)
				upwardVolumeRun = decimal.Zero
			}

			prevPrice = trade.Price // update last price

			// check if a new bar should be created based on the volume run threshold
			if upwardVolumeRun.GreaterThan(c.runVolumeThreshold) || downwardVolumeRun.GreaterThan(c.runVolumeThreshold) {
				finalizedBar := current
				output <- finalizedBar

				// reset the current bar
				current = nil
				// reset the volume runs
				upwardVolumeRun = decimal.Zero
				downwardVolumeRun = decimal.Zero
			}
		}

		if current != nil {
			output <- current
		}
	}()

	return output
}

// Interface guards
var _ Processor = (*VolumeBarConfig)(nil)
var _ Processor = (*VolumeImbalanceBarConfig)(nil)
var _ Processor = (*VolumeRunsBarConfig)(nil)
