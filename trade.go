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

type Side string

const (
	SideBuy  Side = "buy"
	SideSell Side = "sell"
)

type Trade struct {
	Symbol string          `json:"symbol"`
	Price  decimal.Decimal `json:"price"`
	Size   decimal.Decimal `json:"size"`
	Side   Side            `json:"side"`
	Time   time.Time       `json:"time"`
}
