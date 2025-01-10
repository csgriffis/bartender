/*
Copyright © 2025 Chris Griffis <dev@chrisgriffis.com> and contributors.

All rights reserved.
Licensed under the MIT license. See LICENSE file in the project root for details.
*/

package bartender

import (
	"time"

	decimal "github.com/alpacahq/alpacadecimal"
)

func WithInterval(interval time.Duration) Option[TimeBarConfig] {
	return func(v *TimeBarConfig) {
		v.interval = interval
	}
}

type TimeBarConfig struct {
	interval time.Duration `validate:"required"`
}

func (c TimeBarConfig) Process(trades <-chan Trade) chan *Bar {
	output := make(chan *Bar)

	go func() {
		defer close(output)

		var current *Bar
		for trade := range trades {
			alignedStart := calculateAlignedStart(trade.Time, c.interval)

			// check if a new bar should be started
			if current == nil || trade.Time.Sub(current.Start.Add(c.interval)).Nanoseconds() >= 0 || trade.Time.Before(current.Start) {
				// if there is an existing bar, finalize it
				finalizedBar := current

				// start a new bar only if the trade falls on or after the aligned start
				if trade.Time.Before(alignedStart) {
					// skip trades before the first aligned start
					output <- finalizedBar
					continue
				}

				// start a new bar
				newBar := &Bar{
					Open:   trade.Price,
					High:   trade.Price,
					Low:    trade.Price,
					Close:  trade.Price,
					Volume: trade.Size,
					Start:  trade.Time,
				}

				output <- finalizedBar

				// Handle missing intervals between the current bar and the new trade
				if current != nil && trade.Time.Sub(current.Start.Add(c.interval)).Nanoseconds() >= 0 {
					for nextStart := current.Start.Add(c.interval); nextStart.Before(trade.Time); nextStart = nextStart.Add(c.interval) {
						emptyBar := &Bar{
							Open:  current.Close,
							High:  current.Close,
							Low:   current.Close,
							Close: current.Close,
							Start: nextStart,
						}
						output <- emptyBar
					}
				}

				current = newBar
				continue
			}

			// update the current bar
			current.Close = trade.Price
			current.High = decimal.Max(current.High, trade.Price)
			current.Low = decimal.Min(current.Low, trade.Price)
			current.Volume = current.Volume.Add(trade.Size)
		}

		if current != nil {
			output <- current
		}
	}()

	return output
}

// calculateAlignedStart determines the start time of a trade interval
func calculateAlignedStart(t time.Time, interval time.Duration) time.Time {
	intervalSeconds := int64(interval.Seconds())
	timestampSeconds := t.Unix()
	alignedSeconds := (timestampSeconds / intervalSeconds) * intervalSeconds
	return time.Unix(alignedSeconds, 0).UTC()
}

// Interface guards
var _ Processor = (*TimeBarConfig)(nil)
