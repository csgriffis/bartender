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
	Open   decimal.Decimal
	High   decimal.Decimal
	Low    decimal.Decimal
	Close  decimal.Decimal
	Volume decimal.Decimal
	Start  time.Time

	prevPrice  decimal.Decimal
	buyVolume  decimal.Decimal
	sellVolume decimal.Decimal

	ticks   int
	upticks int
}

func (b *Bar) applyTrade(t Trade) {
	// is this the first trade?
	if b.Open.IsZero() && b.upticks == 0 {
		b.Open = t.Price
		b.High = t.Price
		b.Low = t.Price
	}

	// all trades increment the tick count
	b.ticks++

	if b.prevPrice.IsZero() {
		b.prevPrice = t.Price
	}

	if b.Start.IsZero() {
		b.Start = t.Time
	}

	// only increment upticks if the price has increased
	if t.Price.GreaterThan(b.prevPrice) {
		b.upticks++
	}

	if t.Side == SideBuy {
		b.buyVolume = b.buyVolume.Add(t.Size)
	} else {
		b.sellVolume = b.sellVolume.Add(t.Size)
	}

	b.Close = t.Price
	b.High = decimal.Max(b.High, t.Price)
	b.Low = decimal.Min(b.Low, t.Price)
	b.Volume = b.Volume.Add(t.Size)
}

func (b *Bar) BuyVolume() decimal.Decimal {
	return b.buyVolume
}

func (b *Bar) SellVolume() decimal.Decimal {
	return b.sellVolume
}

func (b *Bar) Ticks() int {
	return b.ticks
}

func (b *Bar) Upticks() int {
	return b.upticks
}

func (b *Bar) UnmarshalCSV(record []string) error {
	var err error

	b.Start, err = time.Parse(time.RFC3339Nano, record[0])
	if err != nil {
		return fmt.Errorf("failed to parse start time: %w", err)
	}

	b.Open, err = decimal.NewFromString(record[1])
	if err != nil {
		return fmt.Errorf("failed to parse open price: %w", err)
	}

	b.High, err = decimal.NewFromString(record[2])
	if err != nil {
		return fmt.Errorf("failed to parse high price: %w", err)
	}

	b.Low, err = decimal.NewFromString(record[3])
	if err != nil {
		return fmt.Errorf("failed to parse low price: %w", err)
	}

	b.Close, err = decimal.NewFromString(record[4])
	if err != nil {
		return fmt.Errorf("failed to parse close price: %w", err)
	}

	b.Volume, err = decimal.NewFromString(record[5])
	if err != nil {
		return fmt.Errorf("failed to parse volume: %w", err)
	}

	b.ticks, err = strconv.Atoi(record[6])
	if err != nil {
		return fmt.Errorf("failed to parse ticks: %w", err)
	}

	b.upticks, err = strconv.Atoi(record[7])
	if err != nil {
		return fmt.Errorf("failed to parse upticks: %w", err)
	}

	b.buyVolume, err = decimal.NewFromString(record[8])
	if err != nil {
		return fmt.Errorf("failed to parse buy volume: %w", err)
	}

	b.sellVolume, err = decimal.NewFromString(record[9])
	if err != nil {
		return fmt.Errorf("failed to parse sell volume: %w", err)
	}

	return nil
}

func (b *Bar) MarshalCSV() ([]string, error) {
	return []string{
		b.Start.Format(time.RFC3339Nano),
		b.Open.StringFixed(2),
		b.High.StringFixed(2),
		b.Low.StringFixed(2),
		b.Close.StringFixed(2),
		b.Volume.String(),
		strconv.Itoa(b.ticks),
		strconv.Itoa(b.upticks),
		b.buyVolume.String(),
		b.sellVolume.String(),
	}, nil
}
