[![Go Report Card](https://goreportcard.com/badge/github.com/scmhub/ibapi)](https://goreportcard.com/report/github.com/scmhub/ibapi)
[![Go Reference](https://pkg.go.dev/badge/github.com/scmhub/ibapi.svg)](https://pkg.go.dev/github.com/scmhub/ibapi)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

# Unofficial Golang Interactive Brokers API

`ibapi` package provides an **unofficial** Golang implementation of the [Interactive Brokers](https://www.interactivebrokers.com/en/home.php) API. It is designed to mirror the official Python or C++ [tws-api](https://github.com/InteractiveBrokers) provided by Interactive Brokers.
**We will do our best to keep it in sync** with the official API releases to ensure compatibility and feature parity, but users should be aware that this is a community-driven project and may lag behind the official versions at times.

> [!CAUTION]
> This package is in the **beta phase**. While functional, it may still have bugs or incomplete features. Please test extensively in non-production environments.

## Getting Started

### Prerequisites
- **Go** version 1.23 or higher (recommended)
- An **Interactive Brokers** account with TWS or IB Gateway installed and running

### Installation
Install the package via `go get`:
```bash
go get -u github.com/scmhub/ibapi
```

## Usage
Here’s a basic example to connect and place an order using this package:
```go
package main

import (
	"math/rand"
	"time"

	"github.com/scmhub/ibapi"
)

const (
	IB_HOST = "127.0.0.1"
	IB_PORT = 7497
)

func main() {
	// We set logger for pretty logs to console
	log := ibapi.Logger()
	ibapi.SetConsoleWriter()

	// New IB CLient
	ib := ibapi.NewEClient(nil)
	
    // Connect client
	if err := ib.Connect(IB_HOST, IB_PORT, rand.Int63n(999999)); err != nil {
		log.Error().Err(err)
		return
	}

    // Create and place order
	id := 1
	eurusd := &ibapi.Contract{Symbol: "EUR", SecType: "CASH", Currency: "USD", Exchange: "IDEALPRO"}
	limitOrder := ibapi.LimitOrder("BUY", ibapi.StringToDecimal("20000"), 1.08)
	ib.PlaceOrder(id, eurusd, limitOrder)

	time.Sleep(1 * time.Second)

	err := ib.Disconnect()
	if err != nil {
		log.Error().Err(err).Msg("Disconnect")
	}
}
```

For more information on how to use this package, please refer to the [GoDoc](https://pkg.go.dev/github.com/scmhub/ibapi) documentation.

## Acknowledgments
- Some portions of the code were adapted from [hadrianl](https://github.com/hadrianl/ibapi). Thanks to them for their valuable work!
- Decimals are implemented with the [fixed](https://github.com/robaho/fixed) package

## Notice of Non-Affiliation and Disclaimer
> [!CAUTION]
> This project is in the **beta phase** and is still undergoing testing and development. Users are advised to thoroughly test the software in non-production environments before relying on it for live trading. Features may be incomplete, and bugs may exist. Use at your own risk.

> [!IMPORTANT]
>This project is **not affiliated** with Interactive Brokers Group, Inc. All references to Interactive Brokers, including trademarks, logos, and brand names, belong to their respective owners. The use of these names is purely for informational purposes and does not imply endorsement by Interactive Brokers.

> [!IMPORTANT]
>The authors of this package make **no guarantees** regarding the software's reliability, accuracy, or suitability for any particular purpose, including trading or financial decisions. **No liability** will be accepted for any financial losses, damages, or misinterpretations arising from the use of this software.

## License
Distributed under the MIT License. See [LICENSE](./LICENSE) for more information.

## Author
**Philippe Chavanne** - [contact](https://scm.cx/contact)
