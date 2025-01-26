/*
Copyright © 2025 Chris Griffis <dev@chrisgriffis.com> and contributors.

All rights reserved.
Licensed under the MIT license. See LICENSE file in the project root for details.
*/

package bartender

import (
	"fmt"
	"time"

	decimal "github.com/alpacahq/alpacadecimal"
	"github.com/go-playground/validator/v10"
)

type Processor interface {
	// Process is the handler passed to the Generate functions.
	//
	// This function is responsible for the lifecycle management of the Bar channel it returns. It should watch for
	// the Trade channel to be closed and close the resulting Bar channel when it has completed processing all
	// remaining trades.
	Process(<-chan Trade) chan *Bar
}

type Option[T any] func(*T)

type Bar struct {
	Open   decimal.Decimal
	High   decimal.Decimal
	Low    decimal.Decimal
	Close  decimal.Decimal
	Volume decimal.Decimal
	Start  time.Time
}

func New[T Processor](options ...Option[T]) (Processor, error) {
	var cfg T

	validate := validator.New()

	for _, option := range options {
		option(&cfg)
	}

	if err := validate.Struct(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Generate processes trades synchronously. It accepts all trades to process and returns all bars generated from
// the provided trades.
func Generate(trades []Trade, processor Processor) ([]Bar, error) {
	if len(trades) == 0 {
		return nil, fmt.Errorf("no trades provided")
	}

	bars := make([]Bar, 0, len(trades))
	tradesCh := make(chan Trade)

	go func() {
		defer close(tradesCh)
		for _, trade := range trades {
			tradesCh <- trade
		}
	}()

	for bar := range processor.Process(tradesCh) {
		if bar != nil {
			bars = append(bars, *bar)
		}
	}

	return bars, nil
}

// GenerateStream processes a channel of trades and returns completed bars on the response channel.
func GenerateStream(trades <-chan Trade, processor Processor) (<-chan Bar, error) {
	if trades == nil {
		return nil, fmt.Errorf("trades channel is nil")
	}

	bars := make(chan Bar)

	go func(trades <-chan Trade) {
		defer close(bars)
		for bar := range processor.Process(trades) {
			if bar != nil {
				bars <- *bar
			}
		}
	}(trades)

	return bars, nil
}
