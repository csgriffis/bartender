/*
Copyright © 2025 Chris Griffis <dev@chrisgriffis.com> and contributors.

All rights reserved.
Licensed under the MIT license. See LICENSE file in the project root for details.
*/

package bartender

type FilterFunc func(Trade) bool

// Filter returns a function that filters trades based on the provided filter function.
func Filter(filter func(Trade) bool) func(trades chan Trade) chan Trade {
	return func(trades chan Trade) chan Trade {
		output := make(chan Trade)

		go func() {
			defer close(output)

			for trade := range trades {
				if filter(trade) {
					output <- trade
				}
			}
		}()

		return output
	}
}
