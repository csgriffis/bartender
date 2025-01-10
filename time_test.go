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

type TestCase struct {
	name     string
	interval time.Duration
	trades   []bartender.Trade
	want     []bartender.Bar
	wantErr  bool
}

func TestTimeBarGenerator_Process(t *testing.T) {
	tests := []TestCase{
		{
			name:     "Single Bar",
			interval: 1 * time.Minute,
			trades: []bartender.Trade{
				{Price: decimal.NewFromInt(100), Size: decimal.NewFromInt(1), Time: time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)},
				{Price: decimal.NewFromInt(101), Size: decimal.NewFromInt(2), Time: time.Date(2025, 1, 1, 10, 0, 30, 0, time.UTC)},
				{Price: decimal.NewFromInt(102), Size: decimal.NewFromInt(3), Time: time.Date(2025, 1, 1, 10, 0, 45, 0, time.UTC)},
			},
			want: []bartender.Bar{
				{
					Open:   decimal.NewFromInt(100),
					High:   decimal.NewFromInt(102),
					Low:    decimal.NewFromInt(100),
					Close:  decimal.NewFromInt(102),
					Volume: decimal.NewFromInt(6),
					Start:  time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC),
				},
			},
			wantErr: false,
		},
		{
			name:     "Multiple Bars",
			interval: 1 * time.Minute,
			trades: []bartender.Trade{
				{Price: decimal.NewFromInt(100), Size: decimal.NewFromInt(1), Time: time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)},
				{Price: decimal.NewFromInt(101), Size: decimal.NewFromInt(2), Time: time.Date(2025, 1, 1, 10, 1, 0, 0, time.UTC)},
				{Price: decimal.NewFromInt(102), Size: decimal.NewFromInt(3), Time: time.Date(2025, 1, 1, 10, 2, 0, 0, time.UTC)},
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
				{
					Open:   decimal.NewFromInt(101),
					High:   decimal.NewFromInt(101),
					Low:    decimal.NewFromInt(101),
					Close:  decimal.NewFromInt(101),
					Volume: decimal.NewFromInt(2),
					Start:  time.Date(2025, 1, 1, 10, 1, 0, 0, time.UTC),
				},
				{
					Open:   decimal.NewFromInt(102),
					High:   decimal.NewFromInt(102),
					Low:    decimal.NewFromInt(102),
					Close:  decimal.NewFromInt(102),
					Volume: decimal.NewFromInt(3),
					Start:  time.Date(2025, 1, 1, 10, 2, 0, 0, time.UTC),
				},
			},
			wantErr: false,
		},
		{
			name:     "Missing Trades for interval",
			interval: 1 * time.Minute,
			trades: []bartender.Trade{
				{Price: decimal.NewFromInt(100), Size: decimal.NewFromInt(1), Time: time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)},
				{Price: decimal.NewFromInt(101), Size: decimal.NewFromInt(2), Time: time.Date(2025, 1, 1, 10, 1, 0, 0, time.UTC)},
				{Price: decimal.NewFromInt(102), Size: decimal.NewFromInt(3), Time: time.Date(2025, 1, 1, 10, 3, 0, 0, time.UTC)},
				{Price: decimal.NewFromInt(103), Size: decimal.NewFromInt(6), Time: time.Date(2025, 1, 1, 10, 6, 0, 0, time.UTC)},
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
				{
					Open:   decimal.NewFromInt(101),
					High:   decimal.NewFromInt(101),
					Low:    decimal.NewFromInt(101),
					Close:  decimal.NewFromInt(101),
					Volume: decimal.NewFromInt(2),
					Start:  time.Date(2025, 1, 1, 10, 1, 0, 0, time.UTC),
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
					Open:   decimal.NewFromInt(102),
					High:   decimal.NewFromInt(102),
					Low:    decimal.NewFromInt(102),
					Close:  decimal.NewFromInt(102),
					Volume: decimal.NewFromInt(3),
					Start:  time.Date(2025, 1, 1, 10, 3, 0, 0, time.UTC),
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
					Open:   decimal.NewFromInt(103),
					High:   decimal.NewFromInt(103),
					Low:    decimal.NewFromInt(103),
					Close:  decimal.NewFromInt(103),
					Volume: decimal.NewFromInt(6),
					Start:  time.Date(2025, 1, 1, 10, 6, 0, 0, time.UTC),
				},
			},
			wantErr: false,
		},
		{
			name:     "No Trades",
			interval: 1 * time.Minute,
			trades:   []bartender.Trade{},
			want:     []bartender.Bar{},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g, err := bartender.New(bartender.WithInterval(tt.interval))
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

			barCh, err := bartender.GenerateStream(tradesChan, g)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateStream() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			barsGot := make([]bartender.Bar, 0, len(tt.trades))
			for bar := range barCh {
				barsGot = append(barsGot, bar)
			}

			if !reflect.DeepEqual(barsGot, tt.want) {
				t.Errorf("GenerateStream() = %v, want %v", barsGot, tt.want)
			}
		})
	}
}
