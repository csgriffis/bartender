/*
Copyright © 2025 Chris Griffis <dev@chrisgriffis.com> and contributors.

All rights reserved.
Licensed under the MIT license. See LICENSE file in the project root for details.
*/

package bartender_test

import (
	"reflect"
	"testing"

	"github.com/csgriffis/bartender"
)

type TestCase[T any] struct {
	name    string
	input   T
	trades  []bartender.Trade
	want    []bartender.Bar
	wantErr bool
}

// Run is a helper method to run a test case
func (tc TestCase[T]) Run(t *testing.T, p bartender.Processor) {
	t.Run(tc.name, func(t *testing.T) {
		tradesChan := make(chan bartender.Trade)
		go func() {
			defer close(tradesChan)
			for _, trade := range tc.trades {
				tradesChan <- trade
			}
		}()

		barCh, err := bartender.GenerateStream(tradesChan, p)
		if (err != nil) != tc.wantErr {
			t.Errorf("GenerateStream() error = %v, wantErr %v", err, tc.wantErr)
			return
		}

		barsGot := make([]bartender.Bar, 0, len(tc.trades))
		for bar := range barCh {
			barsGot = append(barsGot, bar)
		}

		if !reflect.DeepEqual(barsGot, tc.want) {
			t.Errorf("GenerateStream() = %v, want %v", barsGot, tc.want)
		}
	})
}
