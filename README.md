# Bartender

[![Go Report Card](https://goreportcard.com/badge/github.com/csgriffis/bartender)](https://goreportcard.com/report/github.com/csgriffis/bartender)
[![Build Status](https://github.com/csgriffis/bartender/actions/workflows/go.yml/badge.svg)](https://github.com/csgriffis/bartender/actions/workflows/go.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](./LICENSE)
[![GoDoc](https://pkg.go.dev/badge/github.com/csgriffis/bartender.svg)](https://pkg.go.dev/github.com/csgriffis/bartender)
[![Code Coverage](https://img.shields.io/codecov/c/github/csgriffis/bartender)](https://codecov.io/gh/csgriffis/bartender)

---

Bartender is a library for generating candlesticks for different aggregation types, including Time, Tick, Volume, and Information bars.

---

## Install

```bash
go get github.com/csgriffis/bartender
```

## Example


### Async Using Channels

```go
import (
    "fmt"
    "time"

    "github.com/csgriffis/bartender"
)

generator, err := bartender.New(bartender.WithInterval(tt.interval))
check(err)

tradesStream := make(chan bartender.Trade)

go func() {
	defer close(tradesStream)

	// Simulate a stream of trades
	trades := []bartender.Trade{}
	for _, trade := range trades {
		tradesStream <- trade
	}
}()

barStream, err := bartender.GenerateStream(tradesStream, generator)
check(err)

for bar := range barStream {
    fmt.Printf("Bar: %v\n", bar)
}

```

### Synchronous

```go
import (
    "fmt"
    "time"

    "github.com/csgriffis/bartender"
)

generator, err := bartender.New(bartender.WithInterval(tt.interval))
check(err)

trades := []bartender.Trade{}
bars, err := bartender.Generate(trades, generator)
check(err)

for _, bar := range bars {
    fmt.Printf("Bar: %v\n", bar)
}
```

### Candlestick Generation

This library uses generics when creating the bar generator. The type of aggregation will be inferred from the
configuration options provided to the generator. Each option returns a distinct generator.

The following are the available options:

#### Dollar Bars
- `WithDollarThreshold`: Aggregates bars based on the dollar volume of the trades.
- `WithDollarImbalanceThreshold`: Aggregates bars based on the dollar imbalance of the trades.
- `WithDollarRunThreshold`: Aggregates bars based on the running dollar volume of the trades.

#### Tick Bars
- `WithTickThreshold`: Aggregates bars based on the number of ticks.
- `WithTickImbalanceThreshold`: Aggregates bars based on the tick imbalance.
- `WithTickRunThreshold`: Aggregates bars based on the running tick volume.

#### Volume Bars
- `WithVolumeThreshold`: Aggregates bars based on the volume of the trades.
- `WithVolumeImbalanceThreshold`: Aggregates bars based on the volume imbalance.
- `WithVolumeRunThreshold`: Aggregates bars based on the running volume.

#### Time Bars
- `WithInterval`: Aggregates bars based on the time interval.

---
## Contributing

### Setup
Install the following tools:

```bash
brew install golangci-lint pre-commit
```

Configure `pre-commit` from the root of the repository:
```bash
pre-commit install
```

---

### License

This project is licensed under the MIT License.
