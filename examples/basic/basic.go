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
	//ibapi.SetLogLevel(int(zerolog.TraceLevel))
	ibapi.SetConsoleWriter()

	// IB CLient
	ib := ibapi.NewEClient(nil)

	if err := ib.Connect(IB_HOST, IB_PORT, rand.Int63n(999999)); err != nil {
		log.Error().Err(err)
		return
	}

	ib.SetConnectionOptions("+PACEAPI")

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
	log.Info().Int64("ID", ib.NextID()).Msg("Next ID")
	log.Info().Int64("ID", ib.NextID()).Msg("Next ID")

	// time.Sleep(1 * time.Second)
	ib.ReqCurrentTime()

	// ########## account ##########
	ib.ReqManagedAccts()

	// ib.ReqAutoOpenOrders(true)
	// ib.ReqAutoOpenOrders(false)
	// ib.ReqAccountUpdates(true, "")
	ib.ReqAllOpenOrders()
	ib.ReqPositions()

	// tags := []string{"AccountType", "NetLiquidation", "TotalCashValue", "SettledCash",
	// 	"sAccruedCash", "BuyingPower", "EquityWithLoanValue",
	// 	"PreviousEquityWithLoanValue", "GrossPositionValue", "ReqTEquity",
	// 	"ReqTMargin", "SMA", "InitMarginReq", "MaintMarginReq", "AvailableFunds",
	// 	"ExcessLiquidity", "Cushion", "FullInitMarginReq", "FullMaintMarginReq",
	// 	"FullAvailableFunds", "FullExcessLiquidity", "LookAheadNextChange",
	// 	"LookAheadInitMarginReq", "LookAheadMaintMarginReq",
	// 	"LookAheadAvailableFunds", "LookAheadExcessLiquidity",
	// 	"HighestSeverity", "DayTradesRemaining", "Leverage", "$LEDGER:ALL"}
	// id := ib.NextID()
	// ib.ReqAccountSummary(id, "All", strings.Join(tags, ","))
	// time.Sleep(10 * time.Second)
	// ib.CancelAccountSummary(id)

	// ib.ReqFamilyCodes()
	// ib.ReqScannerParameters()

	// ########## market data ##########
	// eurusd := &ibapi.Contract{Symbol: "EUR", SecType: "CASH", Currency: "USD", Exchange: "IDEALPRO"}
	// id := ib.NextID()
	// ib.ReqMktData(id, eurusd, "", false, false, nil)
	// time.Sleep(4 * time.Second)
	// ib.CancelMktData(id)

	// ########## real time bars ##########
	// aapl := &ibapi.Contract{ConID: 265598, Symbol: "AAPL", SecType: "STK", Exchange: "SMART"}
	// id := ib.NextID()
	// ib.ReqRealTimeBars(id, aapl, 5, "TRADES", false, nil)
	// time.Sleep(10 * time.Second)
	// ib.CancelRealTimeBars(id)

	//  ########## contract ##########
	// ib.ReqContractDetails(ib.NextID(), aapl)
	// ib.ReqMatchingSymbols(ib.NextID(), "ibm")

	// ########## orders ##########
	// id := ib.NextID()
	// eurusd := &ibapi.Contract{Symbol: "EUR", SecType: "CASH", Currency: "USD", Exchange: "IDEALPRO"}
	// limitOrder := ibapi.LimitOrder("BUY", ibapi.StringToDecimal("20000"), 1.08)
	// ib.PlaceOrder(id, eurusd, limitOrder)
	// time.Sleep(4 * time.Second)
	// ib.CancelOrder(id, ibapi.NewOrderCancel())
	// time.Sleep(4 * time.Second)
	// ib.ReqGlobalCancel()

	time.Sleep(4 * time.Second)
	err := ib.Disconnect()
	if err != nil {
		log.Error().Err(err).Msg("Disconnect")
	}
	log.Info().Msg("Bye!!!!")
}
