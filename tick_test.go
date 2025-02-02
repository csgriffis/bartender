/*
Copyright © 2025 Chris Griffis <dev@chrisgriffis.com> and contributors.

All rights reserved.
Licensed under the MIT license. See LICENSE file in the project root for details.
*/

package bartender_test

import (
	"testing"
	"time"

	decimal "github.com/alpacahq/alpacadecimal"
	"github.com/csgriffis/bartender"
)

func TestTickThresholdBarConfig_Process(t *testing.T) {
	tt := []TestCase[int64]{
		{
			name:   "Single Bar with No Threshold Trigger",
			input:  100,
			trades: []bartender.Trade{},
			want:   []bartender.Bar{},
		},
		{
			name:  "Single Bar with Threshold Trigger",
			input: 100,
			trades: []bartender.Trade{
				{Price: decimal.NewFromInt(100), Size: decimal.NewFromInt(1), Side: bartender.SideBuy, Time: time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)},
			},
			want: []bartender.Bar{
				{
					Open:      decimal.NewFromInt(100),
					High:      decimal.NewFromInt(100),
					Low:       decimal.NewFromInt(100),
					Close:     decimal.NewFromInt(100),
					Volume:    decimal.NewFromInt(1),
					Start:     time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC),
					BuyVolume: decimal.NewFromInt(1),
					Ticks:     1,
				},
			},
		},
		{
			name:  "Multiple Bars with Threshold Trigger",
			input: 2,
			trades: []bartender.Trade{
				{Price: decimal.NewFromInt(100), Size: decimal.NewFromInt(1), Side: bartender.SideBuy, Time: time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)},
				{Price: decimal.NewFromInt(101), Size: decimal.NewFromInt(1), Side: bartender.SideBuy, Time: time.Date(2025, 1, 1, 10, 0, 30, 0, time.UTC)},
				{Price: decimal.NewFromInt(102), Size: decimal.NewFromInt(1), Side: bartender.SideBuy, Time: time.Date(2025, 1, 1, 10, 0, 45, 0, time.UTC)},
			},
			want: []bartender.Bar{
				{
					Open:      decimal.NewFromInt(100),
					High:      decimal.NewFromInt(101),
					Low:       decimal.NewFromInt(100),
					Close:     decimal.NewFromInt(101),
					Volume:    decimal.NewFromInt(2),
					Start:     time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC),
					BuyVolume: decimal.NewFromInt(2),
					Ticks:     2,
					Upticks:   1,
				},
				{
					Open:      decimal.NewFromInt(102),
					High:      decimal.NewFromInt(102),
					Low:       decimal.NewFromInt(102),
					Close:     decimal.NewFromInt(102),
					Volume:    decimal.NewFromInt(1),
					Start:     time.Date(2025, 1, 1, 10, 0, 45, 0, time.UTC),
					BuyVolume: decimal.NewFromInt(1),
					Ticks:     1,
				},
			},
		},
		{
			name:   "No Trades",
			input:  100,
			trades: []bartender.Trade{},
			want:   []bartender.Bar{},
		},
	}

	for _, tc := range tt {
		p, err := bartender.New(bartender.WithTickThreshold(tc.input))
		if err != nil {
			t.Errorf("New() error = %v", err)
			return
		}

		tc.Run(t, p)
	}
}

func TestTickImbalanceBarConfig_Process(t *testing.T) {
	tt := []TestCase[int64]{
		{
			name:  "Single Bar with No Imbalance Trigger",
			input: 100,
			trades: []bartender.Trade{
				{Price: decimal.NewFromInt(100), Size: decimal.NewFromInt(1), Side: bartender.SideBuy, Time: time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)},
			},
			want: []bartender.Bar{
				{
					Open:      decimal.NewFromInt(100),
					High:      decimal.NewFromInt(100),
					Low:       decimal.NewFromInt(100),
					Close:     decimal.NewFromInt(100),
					Volume:    decimal.NewFromInt(1),
					Start:     time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC),
					BuyVolume: decimal.NewFromInt(1),
					Ticks:     1,
				},
			},
		},
		{
			name:  "Single Bar with Imbalance Trigger",
			input: 2,
			trades: []bartender.Trade{
				{Price: decimal.NewFromInt(100), Size: decimal.NewFromInt(1), Side: bartender.SideBuy, Time: time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)},
				{Price: decimal.NewFromInt(101), Size: decimal.NewFromInt(1), Side: bartender.SideBuy, Time: time.Date(2025, 1, 1, 10, 0, 30, 0, time.UTC)},
				{Price: decimal.NewFromInt(102), Size: decimal.NewFromInt(1), Side: bartender.SideBuy, Time: time.Date(2025, 1, 1, 10, 0, 45, 0, time.UTC)},
			},
			want: []bartender.Bar{
				{
					Open:      decimal.NewFromInt(100),
					High:      decimal.NewFromInt(101),
					Low:       decimal.NewFromInt(100),
					Close:     decimal.NewFromInt(101),
					Volume:    decimal.NewFromInt(2),
					Start:     time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC),
					BuyVolume: decimal.NewFromInt(2),
					Ticks:     2,
					Upticks:   1,
				},
				{
					Open:      decimal.NewFromInt(102),
					High:      decimal.NewFromInt(102),
					Low:       decimal.NewFromInt(102),
					Close:     decimal.NewFromInt(102),
					Volume:    decimal.NewFromInt(1),
					Start:     time.Date(2025, 1, 1, 10, 0, 45, 0, time.UTC),
					BuyVolume: decimal.NewFromInt(1),
					Ticks:     1,
				},
			},
		},
		{
			name:   "No Trades",
			input:  100,
			trades: []bartender.Trade{},
			want:   []bartender.Bar{},
		},
	}

	for _, tc := range tt {
		p, err := bartender.New(bartender.WithTickImbalanceThreshold(tc.input))
		if err != nil {
			t.Errorf("New() error = %v", err)
			return
		}

		tc.Run(t, p)
	}
}

func TestTickRunBarConfig_Process(t *testing.T) {
	tt := []TestCase[int64]{
		{
			name:   "Single Bar with No Runs Trigger",
			input:  100,
			trades: []bartender.Trade{},
			want:   []bartender.Bar{},
		},
		{
			name:  "Multiple Bars with Runs Trigger",
			input: 2,
			trades: []bartender.Trade{
				{Price: decimal.NewFromInt(100), Size: decimal.NewFromInt(1), Side: bartender.SideBuy, Time: time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)},
				{Price: decimal.NewFromInt(101), Size: decimal.NewFromInt(1), Side: bartender.SideBuy, Time: time.Date(2025, 1, 1, 10, 0, 30, 0, time.UTC)},
				{Price: decimal.NewFromInt(102), Size: decimal.NewFromInt(1), Side: bartender.SideBuy, Time: time.Date(2025, 1, 1, 10, 0, 45, 0, time.UTC)},
			},
			want: []bartender.Bar{
				{
					Open:      decimal.NewFromInt(100),
					High:      decimal.NewFromInt(102),
					Low:       decimal.NewFromInt(100),
					Close:     decimal.NewFromInt(102),
					Volume:    decimal.NewFromInt(3),
					Start:     time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC),
					BuyVolume: decimal.NewFromInt(3),
					Ticks:     3,
					Upticks:   2,
				},
			},
		},
		{
			name:   "No Trades",
			input:  100,
			trades: []bartender.Trade{},
			want:   []bartender.Bar{},
		},
	}

	for _, tc := range tt {
		p, err := bartender.New(bartender.WithTickRunThreshold(tc.input))
		if err != nil {
			t.Errorf("New() error = %v", err)
			return
		}

		tc.Run(t, p)
	}
}
