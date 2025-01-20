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

func TestDollarThresholdBarConfig_Process(t *testing.T) {
	tt := []TestCase[float64]{
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
					Open:   decimal.NewFromInt(100),
					High:   decimal.NewFromInt(100),
					Low:    decimal.NewFromInt(100),
					Close:  decimal.NewFromInt(100),
					Volume: decimal.NewFromInt(1),
					Start:  time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			name:  "Multiple Bars with Threshold Trigger",
			input: 200,
			trades: []bartender.Trade{
				{Price: decimal.NewFromInt(100), Size: decimal.NewFromInt(1), Side: bartender.SideBuy, Time: time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)},
				{Price: decimal.NewFromInt(101), Size: decimal.NewFromInt(1), Side: bartender.SideBuy, Time: time.Date(2025, 1, 1, 10, 0, 30, 0, time.UTC)},
				{Price: decimal.NewFromInt(102), Size: decimal.NewFromInt(1), Side: bartender.SideBuy, Time: time.Date(2025, 1, 1, 10, 0, 45, 0, time.UTC)},
			},
			want: []bartender.Bar{
				{
					Open:   decimal.NewFromInt(100),
					High:   decimal.NewFromInt(101),
					Low:    decimal.NewFromInt(100),
					Close:  decimal.NewFromInt(101),
					Volume: decimal.NewFromInt(2),
					Start:  time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC),
				},
				{
					Open:   decimal.NewFromInt(102),
					High:   decimal.NewFromInt(102),
					Low:    decimal.NewFromInt(102),
					Close:  decimal.NewFromInt(102),
					Volume: decimal.NewFromInt(1),
					Start:  time.Date(2025, 1, 1, 10, 0, 45, 0, time.UTC),
				},
			},
		},
	}

	for _, tc := range tt {
		p, err := bartender.New(bartender.WithDollarThreshold(tc.input))
		if err != nil {
			t.Errorf("New() error = %v", err)
			return
		}

		tc.Run(t, p)
	}
}

func TestDollarImbalanceBarConfig_Process(t *testing.T) {
	tt := []TestCase[float64]{
		{
			name:  "Single Bar with No Imbalance Trigger",
			input: 1000,
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
			name:  "Multiple Bars with Imbalance Trigger",
			input: 300,
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
			name:   "No Trades",
			input:  100,
			trades: []bartender.Trade{},
			want:   []bartender.Bar{},
		},
	}

	for _, tc := range tt {
		p, err := bartender.New(bartender.WithDollarImbalanceThreshold(tc.input))
		if err != nil {
			t.Errorf("New() error = %v", err)
			return
		}

		tc.Run(t, p)
	}
}

func TestDollarRunBarConfig_Process(t *testing.T) {
	tt := []TestCase[float64]{
		{
			name:  "Single Bar with No Runs Trigger",
			input: 100,
			trades: []bartender.Trade{
				{Price: decimal.NewFromInt(100), Size: decimal.NewFromInt(1), Side: bartender.SideBuy, Time: time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)},
			},
			want: []bartender.Bar{
				{
					Open:   decimal.NewFromInt(100),
					High:   decimal.NewFromInt(100),
					Low:    decimal.NewFromInt(100),
					Close:  decimal.NewFromInt(100),
					Volume: decimal.NewFromInt(1),
					Start:  time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			name:  "Multiple Bars with Runs Trigger",
			input: 200,
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
			name:   "No Trades",
			input:  100,
			trades: []bartender.Trade{},
			want:   []bartender.Bar{},
		},
	}

	for _, tc := range tt {
		p, err := bartender.New(bartender.WithDollarRunThreshold(tc.input))
		if err != nil {
			t.Errorf("New() error = %v", err)
			return
		}

		tc.Run(t, p)
	}
}
