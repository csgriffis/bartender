/*
Copyright © 2025 Chris Griffis <dev@chrisgriffis.com> and contributors.

All rights reserved.
Licensed under the MIT license. See LICENSE file in the project root for details.
*/

package bartender_test

import (
	"github.com/csgriffis/bartender"
	"reflect"
	"testing"
	"time"

	decimal "github.com/alpacahq/alpacadecimal"
)

func TestDollarImbalanceBarConfig_Process(t *testing.T) {
	tests := []struct {
		name               string
		imbalanceThreshold float64
		trades             []bartender.Trade
		want               []bartender.Bar
	}{
		{
			name:               "Single Bar with No Imbalance Trigger",
			imbalanceThreshold: 1000,
			trades: []bartender.Trade{
				{Price: decimal.NewFromInt(100), Size: decimal.NewFromInt(1), Side: bartender.SideBuy, Time: time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)},
				{Price: decimal.NewFromInt(101), Size: decimal.NewFromInt(1), Side: bartender.SideBuy, Time: time.Date(2025, 1, 1, 10, 0, 30, 0, time.UTC)},
				{Price: decimal.NewFromInt(102), Size: decimal.NewFromInt(1), Side: bartender.SideBuy, Time: time.Date(2025, 1, 1, 10, 0, 45, 0, time.UTC)},
			},
			want: []bartender.Bar{
				{
					Open:   decimal.NewFromInt(100),
					High:   decimal.NewFromInt(102),
					Low:    decimal.NewFromInt(100),
					Close:  decimal.NewFromInt(102),
					Volume: decimal.NewFromInt(3),
					Start:  time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			name:               "Multiple Bars with Imbalance Trigger",
			imbalanceThreshold: 300,
			trades: []bartender.Trade{
				{Price: decimal.NewFromInt(100), Size: decimal.NewFromInt(2), Side: bartender.SideBuy, Time: time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)},
				{Price: decimal.NewFromInt(101), Size: decimal.NewFromInt(3), Side: bartender.SideSell, Time: time.Date(2025, 1, 1, 10, 0, 30, 0, time.UTC)},
				{Price: decimal.NewFromInt(102), Size: decimal.NewFromInt(4), Side: bartender.SideBuy, Time: time.Date(2025, 1, 1, 10, 1, 0, 0, time.UTC)},
			},
			want: []bartender.Bar{
				{
					Open:   decimal.NewFromInt(100),
					High:   decimal.NewFromInt(101),
					Low:    decimal.NewFromInt(100),
					Close:  decimal.NewFromInt(101),
					Volume: decimal.NewFromInt(5),
					Start:  time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC),
				},
				{
					Open:   decimal.NewFromInt(102),
					High:   decimal.NewFromInt(102),
					Low:    decimal.NewFromInt(102),
					Close:  decimal.NewFromInt(102),
					Volume: decimal.NewFromInt(4),
					Start:  time.Date(2025, 1, 1, 10, 1, 0, 0, time.UTC),
				},
			},
		},
		{
			name:               "No Trades",
			imbalanceThreshold: 100,
			trades:             []bartender.Trade{},
			want:               []bartender.Bar{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := bartender.New(bartender.WithDollarImbalanceThreshold(tt.imbalanceThreshold))
			if err != nil {
				t.Errorf("New() error = %v", err)
				return
			}

			tradesChan := make(chan bartender.Trade)
			go func() {
				defer close(tradesChan)
				for _, trade := range tt.trades {
					tradesChan <- trade
				}
			}()

			barCh := config.Process(tradesChan)

			barsGot := make([]bartender.Bar, 0, len(tt.trades))
			for bar := range barCh {
				if bar != nil {
					barsGot = append(barsGot, *bar)
				}
			}

			if !reflect.DeepEqual(barsGot, tt.want) {
				t.Errorf("Process() = %v, want %v", barsGot, tt.want)
			}
		})
	}
}
