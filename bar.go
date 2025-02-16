/*
Copyright © 2025 Chris Griffis <dev@chrisgriffis.com> and contributors.

All rights reserved.
Licensed under the MIT license. See LICENSE file in the project root for details.
*/

package bartender

import (
	"fmt"
	"strconv"
	"time"

	decimal "github.com/alpacahq/alpacadecimal"
)

type Bar struct {
	Symbol string          `json:"symbol"`
	Open   decimal.Decimal `json:"open"`
	High   decimal.Decimal `json:"high"`
	Low    decimal.Decimal `json:"low"`
	Close  decimal.Decimal `json:"close"`
	Volume decimal.Decimal `json:"volume"`
	Start  time.Time       `json:"start"`

	// Intra-bar statistics
	BuyVolume  decimal.Decimal `json:"buy_volume"`
	SellVolume decimal.Decimal `json:"sell_volume"`
	Ticks      int             `json:"ticks"`
	Upticks    int             `json:"upticks"`

	prevPrice decimal.Decimal
}

func (b *Bar) applyTrade(t Trade) {
	if b.Symbol == "" {
		b.Symbol = t.Symbol
	}

	// is this the first trade?
	if b.Open.IsZero() && b.Upticks == 0 {
		b.Open = t.Price
		b.High = t.Price
		b.Low = t.Price
	}

	// all trades increment the tick count
	b.Ticks++

	if b.prevPrice.IsZero() {
		b.prevPrice = t.Price
	}

	if b.Start.IsZero() {
		b.Start = t.Time
	}

	// only increment Upticks if the price has increased
	if t.Price.GreaterThan(b.prevPrice) {
		b.Upticks++
	}

	if t.Side == SideBuy {
		b.BuyVolume = b.BuyVolume.Add(t.Size)
	} else {
		b.SellVolume = b.SellVolume.Add(t.Size)
	}

	b.Close = t.Price
	b.High = decimal.Max(b.High, t.Price)
	b.Low = decimal.Min(b.Low, t.Price)
	b.Volume = b.Volume.Add(t.Size)
}

func (b *Bar) UnmarshalCSV(record []string) error {
	var err error

	b.Symbol = record[0]

	b.Start, err = time.Parse(time.RFC3339Nano, record[1])
	if err != nil {
		return fmt.Errorf("failed to parse start time: %w", err)
	}

	b.Open, err = decimal.NewFromString(record[2])
	if err != nil {
		return fmt.Errorf("failed to parse open price: %w", err)
	}

	b.High, err = decimal.NewFromString(record[3])
	if err != nil {
		return fmt.Errorf("failed to parse high price: %w", err)
	}

	b.Low, err = decimal.NewFromString(record[4])
	if err != nil {
		return fmt.Errorf("failed to parse low price: %w", err)
	}

	b.Close, err = decimal.NewFromString(record[5])
	if err != nil {
		return fmt.Errorf("failed to parse close price: %w", err)
	}

	b.Volume, err = decimal.NewFromString(record[6])
	if err != nil {
		return fmt.Errorf("failed to parse volume: %w", err)
	}

	b.Ticks, err = strconv.Atoi(record[7])
	if err != nil {
		return fmt.Errorf("failed to parse Ticks: %w", err)
	}

	b.Upticks, err = strconv.Atoi(record[8])
	if err != nil {
		return fmt.Errorf("failed to parse Upticks: %w", err)
	}

	b.BuyVolume, err = decimal.NewFromString(record[9])
	if err != nil {
		return fmt.Errorf("failed to parse buy volume: %w", err)
	}

	b.SellVolume, err = decimal.NewFromString(record[10])
	if err != nil {
		return fmt.Errorf("failed to parse sell volume: %w", err)
	}

	return nil
}

func (b *Bar) MarshalCSV() ([]string, error) {
	return []string{
		b.Symbol,
		b.Start.Format(time.RFC3339Nano),
		b.Open.StringFixed(2),
		b.High.StringFixed(2),
		b.Low.StringFixed(2),
		b.Close.StringFixed(2),
		b.Volume.String(),
		strconv.Itoa(b.Ticks),
		strconv.Itoa(b.Upticks),
		b.BuyVolume.String(),
		b.SellVolume.String(),
	}, nil
}
