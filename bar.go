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
