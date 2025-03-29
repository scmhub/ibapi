package main

import (
	"time"

	"github.com/rs/zerolog"
	"github.com/scmhub/ibapi"
)

const (
	host = "localhost"
	port = 7497
)

var orderID int64

func nextID() int64 {
	orderID++
	return orderID
}

func main() {
	// We set logger for pretty logs to console
	log := ibapi.Logger()
	ibapi.SetLogLevel(int(zerolog.TraceLevel))
	ibapi.SetConsoleWriter()
	// ibapi.SetConnectionTimeout(1 * time.Second)

	// IB CLient
	ib := ibapi.NewEClient(nil)

	if err := ib.Connect(host, port, 5); err != nil { //rand.Int63n(999999)
		log.Error().Err(err).Msg("Connect")
		return
	}

	// Add a short delay to allow the connection to stabilize
	time.Sleep(100 * time.Millisecond)
	log.Info().Msg("Waited for connection to stabilize")

	// ib.SetConnectionOptions("+PACEAPI")

	// Logger test
	// log.Trace().Interface("Log level", log.GetLevel()).Msg("Logger Trace")
	// log.Debug().Interface("Log level", log.GetLevel()).Msg("Logger Debug")
	// log.Info().Interface("Log level", log.GetLevel()).Msg("Logger Info")
	// log.Warn().Interface("Log level", log.GetLevel()).Msg("Logger Warn")
	// log.Error().Interface("Log level", log.GetLevel()).Msg("Logger Error")

	// time.Sleep(1 * time.Second)
	// log.Print("Is connected: ", ib.IsConnected())
	// log.Print("Server Version: ", ib.ServerVersion())
	// log.Print("TWS Connection time: ", ib.TWSConnectionTime())

	// time.Sleep(1 * time.Second)
	// ib.ReqCurrentTime()

	// ########## account ##########
	ib.ReqManagedAccts()

	// ib.ReqAutoOpenOrders(false) // Only from clientID = 0
	// ib.ReqAutoOpenOrders(false)
	// ib.ReqAccountUpdates(true, "")
	// ib.ReqAllOpenOrders()
	// ib.ReqPositions()
	// ib.ReqCompletedOrders(false)

	// tags := []string{"AccountType", "NetLiquidation", "TotalCashValue", "SettledCash",
	// 	"sAccruedCash", "BuyingPower", "EquityWithLoanValue",
	// 	"PreviousEquityWithLoanValue", "GrossPositionValue", "ReqTEquity",
	// 	"ReqTMargin", "SMA", "InitMarginReq", "MaintMarginReq", "AvailableFunds",
	// 	"ExcessLiquidity", "Cushion", "FullInitMarginReq", "FullMaintMarginReq",
	// 	"FullAvailableFunds", "FullExcessLiquidity", "LookAheadNextChange",
	// 	"LookAheadInitMarginReq", "LookAheadMaintMarginReq",
	// 	"LookAheadAvailableFunds", "LookAheadExcessLiquidity",
	// 	"HighestSeverity", "DayTradesRemaining", "Leverage", "$LEDGER:ALL"}
	id := nextID()
	// ib.ReqAccountSummary(id, "All", strings.Join(tags, ","))
	// time.Sleep(10 * time.Second)
	// ib.CancelAccountSummary(id)

	// ib.ReqFamilyCodes()
	// ib.ReqScannerParameters()

	// ########## market data ##########
	eurusd := &ibapi.Contract{Symbol: "EUR", SecType: "CASH", Currency: "USD", Exchange: "IDEALPRO"}
	// id := nextID()
	// ib.ReqMktData(id, eurusd, "", false, false, nil)
	// time.Sleep(4 * time.Second)
	// ib.CancelMktData(id)

	// ########## real time bars ##########
	// aapl := &ibapi.Contract{ConID: 265598, Symbol: "AAPL", SecType: "STK", Exchange: "SMART"}
	// id := nextID()
	// ib.ReqRealTimeBars(id, aapl, 5, "TRADES", false, nil)
	// time.Sleep(10 * time.Second)
	// ib.CancelRealTimeBars(id)

	//  ########## contract ##########
	// ib.ReqContractDetails(nextID(), aapl)
	// ib.ReqMatchingSymbols(nextID(), "ibm")

	// ########## orders ##########
	// id := nextID()
	// eurusd := &ibapi.Contract{Symbol: "EUR", SecType: "CASH", Currency: "USD", Exchange: "IDEALPRO"}
	// limitOrder := ibapi.LimitOrder("BUY", ibapi.StringToDecimal("20000"), 1.08)
	// ib.PlaceOrder(id, eurusd, limitOrder)
	// time.Sleep(4 * time.Second)
	// ib.CancelOrder(id, ibapi.NewOrderCancel())
	// time.Sleep(4 * time.Second)
	// ib.ReqGlobalCancel()
	// Real time bars

	duration := "60 S"
	barSize := "5 secs"
	whatToShow := "MIDPOINT" // "TRADES", "MIDPOINT", "BID" or "ASK"
	ib.ReqHistoricalData(id, eurusd, "", duration, barSize, whatToShow, true, 1, true, nil)

	time.Sleep(30 * time.Second)
	ib.CancelHistoricalData(id)
	err := ib.Disconnect()
	if err != nil {
		log.Error().Err(err).Msg("Disconnect")
	}
	log.Info().Msg("Bye!!!!")
}
