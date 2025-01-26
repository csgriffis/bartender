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
	Symbol string
	Price  decimal.Decimal
	Size   decimal.Decimal
	Side   Side
	Time   time.Time
}
