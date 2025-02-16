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

func TestTimeBarConfig_Process(t *testing.T) {
	tt := []TestCase[time.Duration]{
		{
			name:  "Single Bar",
			input: 1 * time.Minute,
			trades: []bartender.Trade{
				{Price: decimal.NewFromInt(100), Side: bartender.SideBuy, Size: decimal.NewFromInt(1), Time: time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)},
				{Price: decimal.NewFromInt(101), Side: bartender.SideBuy, Size: decimal.NewFromInt(2), Time: time.Date(2025, 1, 1, 10, 0, 30, 0, time.UTC)},
				{Price: decimal.NewFromInt(102), Side: bartender.SideBuy, Size: decimal.NewFromInt(3), Time: time.Date(2025, 1, 1, 10, 0, 45, 0, time.UTC)},
			},
			want: []bartender.Bar{
				{
					Open:      decimal.NewFromInt(100),
					High:      decimal.NewFromInt(102),
					Low:       decimal.NewFromInt(100),
					Close:     decimal.NewFromInt(102),
					Volume:    decimal.NewFromInt(6),
					Start:     time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC),
					BuyVolume: decimal.NewFromInt(6),
					Ticks:     3,
					Upticks:   2,
				},
			},
			wantErr: false,
		},
		{
			name:  "Multiple Bars",
			input: 1 * time.Minute,
			trades: []bartender.Trade{
				{Price: decimal.NewFromInt(100), Side: bartender.SideBuy, Size: decimal.NewFromInt(1), Time: time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)},
				{Price: decimal.NewFromInt(101), Side: bartender.SideBuy, Size: decimal.NewFromInt(2), Time: time.Date(2025, 1, 1, 10, 1, 0, 0, time.UTC)},
				{Price: decimal.NewFromInt(102), Side: bartender.SideBuy, Size: decimal.NewFromInt(3), Time: time.Date(2025, 1, 1, 10, 2, 0, 0, time.UTC)},
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
				{
					Open:      decimal.NewFromInt(101),
					High:      decimal.NewFromInt(101),
					Low:       decimal.NewFromInt(101),
					Close:     decimal.NewFromInt(101),
					Volume:    decimal.NewFromInt(2),
					Start:     time.Date(2025, 1, 1, 10, 1, 0, 0, time.UTC),
					BuyVolume: decimal.NewFromInt(2),
					Ticks:     1,
				},
				{
					Open:      decimal.NewFromInt(102),
					High:      decimal.NewFromInt(102),
					Low:       decimal.NewFromInt(102),
					Close:     decimal.NewFromInt(102),
					Volume:    decimal.NewFromInt(3),
					Start:     time.Date(2025, 1, 1, 10, 2, 0, 0, time.UTC),
					BuyVolume: decimal.NewFromInt(3),
					Ticks:     1,
				},
			},
			wantErr: false,
		},
		{
			name:  "Multiple Bars across days",
			input: 60 * time.Minute,
			trades: []bartender.Trade{
				{Price: decimal.NewFromInt(100), Side: bartender.SideBuy, Size: decimal.NewFromInt(1), Time: time.Date(2025, 1, 1, 23, 45, 0, 0, time.UTC)},
				{Price: decimal.NewFromInt(101), Side: bartender.SideBuy, Size: decimal.NewFromInt(2), Time: time.Date(2025, 1, 1, 23, 46, 0, 0, time.UTC)},
				{Price: decimal.NewFromInt(102), Side: bartender.SideBuy, Size: decimal.NewFromInt(3), Time: time.Date(2025, 1, 1, 23, 47, 0, 0, time.UTC)},
				{Price: decimal.NewFromInt(102), Side: bartender.SideBuy, Size: decimal.NewFromInt(3), Time: time.Date(2025, 1, 2, 0, 2, 0, 0, time.UTC)},
			},
			want: []bartender.Bar{
				{
					Open:      decimal.NewFromInt(100),
					High:      decimal.NewFromInt(102),
					Low:       decimal.NewFromInt(100),
					Close:     decimal.NewFromInt(102),
					Volume:    decimal.NewFromInt(6),
					Start:     time.Date(2025, 1, 1, 23, 00, 0, 0, time.UTC),
					BuyVolume: decimal.NewFromInt(6),
					Ticks:     3,
					Upticks:   2,
				},
				{
					Open:      decimal.NewFromInt(102),
					High:      decimal.NewFromInt(102),
					Low:       decimal.NewFromInt(102),
					Close:     decimal.NewFromInt(102),
					Volume:    decimal.NewFromInt(3),
					Start:     time.Date(2025, 1, 2, 0, 00, 0, 0, time.UTC),
					BuyVolume: decimal.NewFromInt(3),
					Ticks:     1,
				},
			},
			wantErr: false,
		},
		{
			name:  "Missing Trades for interval",
			input: 1 * time.Minute,
			trades: []bartender.Trade{
				{Price: decimal.NewFromInt(100), Side: bartender.SideBuy, Size: decimal.NewFromInt(1), Time: time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)},
				{Price: decimal.NewFromInt(101), Side: bartender.SideBuy, Size: decimal.NewFromInt(2), Time: time.Date(2025, 1, 1, 10, 1, 0, 0, time.UTC)},
				{Price: decimal.NewFromInt(102), Side: bartender.SideBuy, Size: decimal.NewFromInt(3), Time: time.Date(2025, 1, 1, 10, 3, 0, 0, time.UTC)},
				{Price: decimal.NewFromInt(103), Side: bartender.SideBuy, Size: decimal.NewFromInt(6), Time: time.Date(2025, 1, 1, 10, 6, 0, 0, time.UTC)},
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
				{
					Open:      decimal.NewFromInt(101),
					High:      decimal.NewFromInt(101),
					Low:       decimal.NewFromInt(101),
					Close:     decimal.NewFromInt(101),
					Volume:    decimal.NewFromInt(2),
					Start:     time.Date(2025, 1, 1, 10, 1, 0, 0, time.UTC),
					BuyVolume: decimal.NewFromInt(2),
					Ticks:     1,
				},
				{
					Open:   decimal.NewFromInt(101),
					High:   decimal.NewFromInt(101),
					Low:    decimal.NewFromInt(101),
					Close:  decimal.NewFromInt(101),
					Volume: decimal.NewFromInt(0),
					Start:  time.Date(2025, 1, 1, 10, 2, 0, 0, time.UTC),
				},
				{
					Open:      decimal.NewFromInt(102),
					High:      decimal.NewFromInt(102),
					Low:       decimal.NewFromInt(102),
					Close:     decimal.NewFromInt(102),
					Volume:    decimal.NewFromInt(3),
					Start:     time.Date(2025, 1, 1, 10, 3, 0, 0, time.UTC),
					BuyVolume: decimal.NewFromInt(3),
					Ticks:     1,
				},
				{
					Open:   decimal.NewFromInt(102),
					High:   decimal.NewFromInt(102),
					Low:    decimal.NewFromInt(102),
					Close:  decimal.NewFromInt(102),
					Volume: decimal.NewFromInt(0),
					Start:  time.Date(2025, 1, 1, 10, 4, 0, 0, time.UTC),
				},
				{
					Open:   decimal.NewFromInt(102),
					High:   decimal.NewFromInt(102),
					Low:    decimal.NewFromInt(102),
					Close:  decimal.NewFromInt(102),
					Volume: decimal.NewFromInt(0),
					Start:  time.Date(2025, 1, 1, 10, 5, 0, 0, time.UTC),
				},
				{
					Open:      decimal.NewFromInt(103),
					High:      decimal.NewFromInt(103),
					Low:       decimal.NewFromInt(103),
					Close:     decimal.NewFromInt(103),
					Volume:    decimal.NewFromInt(6),
					Start:     time.Date(2025, 1, 1, 10, 6, 0, 0, time.UTC),
					BuyVolume: decimal.NewFromInt(6),
					Ticks:     1,
				},
			},
			wantErr: false,
		},
		{
			name:    "No Trades",
			input:   1 * time.Minute,
			trades:  []bartender.Trade{},
			want:    []bartender.Bar{},
			wantErr: false,
		},
	}

	for _, tc := range tt {
		p, err := bartender.New(bartender.WithInterval(tc.input))
		if err != nil {
			t.Errorf("New() error = %v", err)
			return
		}

		tc.Run(t, p)
	}
}
