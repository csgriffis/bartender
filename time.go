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

			// is the trade before the aligned start?
			if trade.Time.Before(alignedStart) {
				// then drop the trade
				continue
			}

			// is the trade beyond the current interval?
			if current != nil && trade.Time.Sub(current.Start.Add(c.interval)).Nanoseconds() >= 0 {
				// then finalize the current interval
				output <- current

				// is there a gap between the current interval and the trade?
				for nextStart := alignedStart.Add(c.interval); current.Start.Add(c.interval).Before(alignedStart); nextStart = nextStart.Add(c.interval) {
					emptyBar := &Bar{
						Open:  current.Close,
						High:  current.Close,
						Low:   current.Close,
						Close: current.Close,
						Start: current.Start.Add(c.interval),
					}
					output <- emptyBar

					current = emptyBar
				}

				// start a new bar
				newBar := &Bar{
					Open:  trade.Price,
					High:  trade.Price,
					Low:   trade.Price,
					Close: trade.Price,
					Start: alignedStart,
				}

				current = newBar
			}

			if current == nil {
				current = &Bar{
					Open:  trade.Price,
					High:  trade.Price,
					Low:   trade.Price,
					Start: alignedStart,
				}
			}

			// update the current bar
			current.Close = trade.Price
			current.High = decimal.Max(current.High, trade.Price)
			current.Low = decimal.Min(current.Low, trade.Price)
			current.Volume = current.Volume.Add(trade.Size)
		}

		// send the last bar
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
