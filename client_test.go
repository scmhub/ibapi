package ibapi

import (
	"math/rand"
	"testing"
	"time"
)

const (
	host    = "localhost"
	port    = 7496
	account = "DUD00029"
)

var clientID = rand.Int63n(999999)

func TestClient(t *testing.T) {
	// We set logger for pretty logs to console
	// SetLogLevel(int(zerolog.TraceLevel))
	// SetConsoleWriter()
	// IB CLient
	ib := NewEClient(nil)

	ib.SetConnectionOptions("+PACEAPI")

	if err := ib.Connect(host, port, clientID); err != nil {
		log.Error().Err(err).Msg("connect")
		return
	}
	// IsConnected
	if !ib.IsConnected() {
		t.Error("not connected")
		return
	}

	reqCurrentTime(ib)
	// pnlSingleOperation(ib)
	// pnlOperation(ib)
	// tickDataOperation(ib)
	// tickOptionComputationOperation(ib)
	// delayedTickDataOperation(ib)
	// marketDepthOperations(ib)
	// realTimeBars(ib)
	// marketDataType(ib)
	// historicalDataRequests(ib)
	// optionsOperations(ib)
	// contractOperations(ib)
	// marketScanners(ib)
	// fundamentals(ib)
	// bulletins(ib)
	// accountOperations(ib)
	// orderOperations(ib)
	// ocaSamples(ib)
	// conditionSamples(ib)
	// bracketSample(ib)
	// hedgeSample(ib)
	// testAlgoSamples(ib)
	// financialAdvisorOrderSamples(ib)
	// financialAdvisorOperations(ib)
	// testDisplayGroups(ib)
	// miscellaneous(ib)
	// reqFamilyCodes(ib)
	// reqMatchingSymbols(ib)
	// reqMktDepthExchanges(ib)
	// reqNewsTicks(ib)
	// reqSmartComponents(ib)
	// reqNewsProviders(ib)
	// reqNewsArticle(ib)
	// reqHistoricalNews(ib)
	// reqHeadTimestamp(ib)
	// reqHistogramData(ib)
	// rerouteCFDOperations(ib)
	// marketRuleOperations(ib)
	// continuousFuturesOperations(ib)
	// reqHistoricalTicks(ib)
	// reqTickByTickData(ib)
	// whatIfSamples(ib)
	// ibkratsSample(ib)
	// wshCalendarOperations(ib)

	time.Sleep(1 * time.Second)

	t.Error("host:", host, "port:", port)
}

func reqCurrentTime(ib *EClient) {
	ib.ReqCurrentTime()
}

func pnlSingleOperation(ib *EClient) {
	ib.ReqPnLSingle(7002, account, "", 268084)
	time.Sleep(2 * time.Second)
	ib.CancelPnLSingle(7002)
}

func pnlOperation(ib *EClient) {
	ib.ReqPnL(7001, account, "")
	time.Sleep(2 * time.Second)
	ib.CancelPnL(7001)
}

func tickDataOperation(ib *EClient) {
	ib.ReqMarketDataType(4)
	time.Sleep(1 * time.Second)
	// ReqMktData
	ib.ReqMktData(1001, StockComboContract(), "", false, false, nil)
	ib.ReqMktData(1002, OptionWithLocalSymbol(), "", false, false, nil)
	// ReqMktData - snapshot
	ib.ReqMktData(1003, FutureComboContract(), "", true, false, nil)
	// ReqMktData - regulatory snapshot - Each regulatory snapshot request incurs a 0.01 USD fee
	// ib.ReqMktData(1013, USStock(), "", false, true, nil)
	// ReqMktData - genticks - Requesting RTVolume (Time & Sales) and shortable generic ticks
	ib.ReqMktData(1004, USStockAtSmart(), "233,236", false, false, nil)
	// ReqMktData - contract news - Without the API news subscription this will generate an "invalid tick type" error
	ib.ReqMktData(1005, USStock(), "mdoff,292:BZ", false, false, nil)
	ib.ReqMktData(1006, USStock(), "mdoff,292:BT", false, false, nil)
	ib.ReqMktData(1007, USStock(), "mdoff,292:FLY", false, false, nil)
	ib.ReqMktData(1008, USStock(), "mdoff,292:DJ-RT", false, false, nil)
	// ReqMktData - broad tape news
	ib.ReqMktData(1009, BTbroadtapeNewsFeed(), "mdoff,292", false, false, nil)
	ib.ReqMktData(1010, BZbroadtapeNewsFeed(), "mdoff,292", false, false, nil)
	ib.ReqMktData(1011, FLYbroadtapeNewsFeed(), "mdoff,292", false, false, nil)
	// ReqMktData - option data genticks - Requesting data for an option contract will return the greek values
	ib.ReqMktData(1013, USOptionContract(), "", false, false, nil)
	// ReqMktData - futures open interest - Requesting data for a futures contract will return the futures open interest
	ib.ReqMktData(1014, SimpleFuture(), "mdoff,588", false, false, nil)
	// ReqMktData - preopen bid/ask - Requesting data for a futures contract will return the pre-open bid/ask flag
	ib.ReqMktData(1015, SimpleFuture(), "", false, false, nil)
	// ReqMktData - avg opt volume - Requesting data for a stock will return the average option volume
	ib.ReqMktData(1016, USStockAtSmart(), "mdoff,105", false, false, nil)
	// ReqMktData - etf ticks
	ib.ReqMktData(1017, ETF(), "mdoff,577,614,623", false, false, nil)
	// ReqMktData - crypto
	ib.ReqMktData(1018, CryptoContract(), "", false, false, nil)
	// ReqMktData - IPO price
	ib.ReqMktData(1019, StockWithIPOPrice(), "mdoff,586", false, false, nil)
	// ReqMktData - yield bid/ask
	ib.ReqMktData(1020, Bond(), "", false, false, nil)

	time.Sleep(1 * time.Second)

	ib.CancelMktData(1001)
	ib.CancelMktData(1002)
	ib.CancelMktData(1003)
	ib.CancelMktData(1014)
	ib.CancelMktData(1015)
	ib.CancelMktData(1016)
	ib.CancelMktData(1017)
	ib.CancelMktData(1018)
	ib.CancelMktData(1019)
	ib.CancelMktData(1020)

}

func tickOptionComputationOperation(ib *EClient) {
	time.Sleep(1 * time.Second)
	ib.ReqMarketDataType(4)
	ib.ReqMktData(2001, OptionWithLocalSymbol(), "", false, false, nil)
	time.Sleep(10 * time.Second)
	ib.CancelMktData(2001)
}

func delayedTickDataOperation(ib *EClient) {
	ib.ReqMarketDataType(4)
	ib.ReqMktData(1013, HKStk(), "", false, false, nil)
	ib.ReqMktData(1014, USOptionContract(), "", false, false, nil)
	time.Sleep(10 * time.Second)
	ib.CancelMktData(1013)
	ib.CancelMktData(1014)
}

func marketDepthOperations(ib *EClient) {
	ib.ReqMktDepth(2001, EurGbpFx(), 5, false, nil)
	time.Sleep(2 * time.Second)
	ib.CancelMktDepth(2001, false)
	ib.ReqMktDepth(2002, EurGbpFx(), 5, true, nil)
	time.Sleep(5 * time.Second)
	ib.CancelMktDepth(2001, true)
}

func realTimeBars(ib *EClient) {
	ib.ReqRealTimeBars(3001, EurGbpFx(), 5, "MIDPOINT", true, nil)
	time.Sleep(2 * time.Second)
	ib.CancelRealTimeBars(3001)
}

func marketDataType(ib *EClient) {
	// By default only real-time (1) market data is enabled
	// Sending frozen (2) enables frozen market data
	// Sending delayed (3) enables delayed market data and disables delayed-frozen market data
	// Sending delayed-frozen (4) enables delayed and delayed-frozen market data
	// Sending real-time (1) disables frozen, delayed and delayed-frozen market data
	ib.ReqMarketDataType(2)
}

func historicalDataRequests(ib *EClient) {
	queryTime := time.Now().AddDate(0, 0, -180).Format("20060102-15:04:05")
	ib.ReqHistoricalData(4001, EurGbpFx(), queryTime, "1 M", "1 day", "MIDPOINT", true, 1, false, nil)
	ib.ReqHistoricalData(4002, EuropeanStock(), queryTime, "10 D", "1 min", "TRADES", true, 1, false, nil)
	ib.ReqHistoricalData(4003, USStockAtSmart(), queryTime, "1 M", "1 day", "SCHEDULE", true, 1, false, nil)
	time.Sleep(2 * time.Second)
	ib.CancelHistoricalData(4001)
	ib.CancelHistoricalData(4002)
	ib.CancelHistoricalData(4003)
}

func optionsOperations(ib *EClient) {
	ib.ReqSecDefOptParams(0, "IBM", "", "STK", 8314)
	ib.CalculateImpliedVolatility(5001, OptionWithLocalSymbol(), 0.5, 55, nil)
	ib.CancelCalculateImpliedVolatility(5001)
	ib.CalculateOptionPrice(5002, OptionWithLocalSymbol(), 0.6, 55, nil)
	ib.CancelCalculateOptionPrice(5002)
	ib.ExerciseOptions(5003, OptionWithTradingClass(), 1, 1, "", 1, "20231018-12:00:00", "CustAcct", true)
}

func contractOperations(ib *EClient) {
	ib.ReqContractDetails(209, EurGbpFx())
	time.Sleep(2 * time.Second)
	// ReqContractDetails
	ib.ReqContractDetails(210, OptionForQuery())
	ib.ReqContractDetails(212, IBMBond())
	ib.ReqContractDetails(213, IBKRStk())
	ib.ReqContractDetails(214, Bond())
	ib.ReqContractDetails(215, FuturesOnOptions())
	ib.ReqContractDetails(216, SimpleFuture())
	ib.ReqContractDetails(219, Fund())
	ib.ReqContractDetails(220, USStock())
	ib.ReqContractDetails(221, USStockAtSmart())
	// ReqContractDetails - news
	ib.ReqContractDetails(211, NewsFeedForQuery())
	// ReqContractDetails - crypto
	ib.ReqContractDetails(217, CryptoContract())
	// ReqContractDetails - by isssuer id
	ib.ReqContractDetails(211, ByIssuerId())
}

func marketScanners(ib *EClient) {
	// ReqScannerParameters - Requesting all available parameters which can be used to build a scanner request
	ib.ReqScannerParameters()
	time.Sleep(2 * time.Second)
	// ReqScannerSubscription - Triggering a scanner subscription
	ib.ReqScannerSubscription(7001, HotUSStkByVolume(), nil, nil)
	// ReqScannerSubscription
	TagValues := []TagValue{
		{Tag: "usdMarketCapAbove", Value: "10000"},
		{Tag: "optVolumeAbove", Value: "1000"},
		{Tag: "usdMarketCapAbove", Value: "100000000"},
	}
	ib.ReqScannerSubscription(7002, HotUSStkByVolume(), nil, TagValues) // requires TWS v973+
	// ReqScannerSubscription - complex scanner
	AAPLConIDTag := []TagValue{
		{Tag: "underConID", Value: "265598"},
	}
	ib.ReqScannerSubscription(7003, ComplexOrdersAndTrades(), nil, AAPLConIDTag)

	time.Sleep(2 * time.Second)
	ib.CancelScannerSubscription(7001)
	ib.CancelScannerSubscription(7002)
}

func fundamentals(ib *EClient) {
	ib.ReqFundamentalData(8001, USStock(), "ReportsFinSummary", nil)
	time.Sleep(2 * time.Second)
	ib.CancelFundamentalData(8001)
}

func bulletins(ib *EClient) {
	ib.ReqNewsBulletins(true)
	time.Sleep(2 * time.Second)
	ib.CancelNewsBulletins()
}

func accountOperations(ib *EClient) {
	ib.ReqManagedAccts()
	time.Sleep(2 * time.Second)
	ib.ReqAccountSummary(9001, "All", GetAllTags())
	time.Sleep(2 * time.Second)
	ib.ReqAccountSummary(9002, "All", "$LEDGER")
	time.Sleep(2 * time.Second)
	ib.ReqAccountSummary(9003, "All", "$LEDGER:EUR")
	time.Sleep(2 * time.Second)
	ib.ReqAccountSummary(9004, "All", "$LEDGER:ALL")
	time.Sleep(2 * time.Second)
	ib.CancelAccountSummary(9001)
	ib.CancelAccountSummary(9002)
	ib.CancelAccountSummary(9003)
	ib.CancelAccountSummary(9004)
	time.Sleep(2 * time.Second)
	ib.ReqAccountUpdates(true, account)
	time.Sleep(2 * time.Second)
	ib.ReqAccountUpdates(false, account)
	time.Sleep(2 * time.Second)
	ib.ReqAccountUpdatesMulti(9005, account, "EUstocks", true)
	time.Sleep(2 * time.Second)
	ib.ReqPositions()
	time.Sleep(2 * time.Second)
	ib.CancelPositions()
	ib.ReqPositionsMulti(9006, account, "EUstocks")
	time.Sleep(2 * time.Second)
	ib.CancelPositionsMulti(9006)
	ib.ReqUserInfo(9007)
}

func orderOperations(ib *EClient) {
	// Requesting the next valid id
	ib.ReqIDs(-1)
	// Requesting Orders
	ib.ReqAllOpenOrders()
	ib.ReqAutoOpenOrders(true)
	ib.ReqOpenOrders()
	// Placing/modifying an order - remember to ALWAYS increment the nextValidId after placing an order so it can be used for the next one!
	ib.PlaceOrder(ib.NextID(), USStock(), LimitOrder("SELL", ONE, 50))

	// ib.PlaceOrder(ib.NextID(), OptionAtBox(), Block("BUY", StringToDecimal("50"), 20))
	// ib.PlaceOrder(ib.NextID(), OptionAtBox(), BoxTop("SELL", StringToDecimal("10")))

	// ib.PlaceOrder(ib.NextID(), FutureComboContract(), ComboLimitOrder("SELL", ONE, 1, false))
	// ib.PlaceOrder(ib.NextID(), StockComboContract(), ComboMarketOrder("BUY", ONE, false))
	// ib.PlaceOrder(ib.NextID(), OptionComboContract(), ComboMarketOrder("BUY", ONE, true))
	// ib.PlaceOrder(ib.NextID(), StockComboContract(), LimitOrderForComboWithLegPrices("BUY", ONE, []float64{10, 5}, true))
	// ib.PlaceOrder(ib.NextID(), USStock(), Discretionary("SELL", ONE, 45, 0.5))
	// ib.PlaceOrder(ib.NextID(), OptionAtBox(), LimitIfTouched("SELL", ONE, 30, 34))
	// ib.PlaceOrder(ib.NextID(), USStock(), LimitOnClose("SELL", ONE, 34))
	// ib.PlaceOrder(ib.NextID(), USStock(), LimitOnOpen("BUY", ONE, 35))
	// ib.PlaceOrder(ib.NextID(), USStock(), MarketIfTouched("BUY", ONE, 35))
	// ib.PlaceOrder(ib.NextID(), USStock(), MarketOnClose("SELL", ONE))
	// ib.PlaceOrder(ib.NextID(), USStock(), MarketOnOpen("BUY", ONE))
	// ib.PlaceOrder(ib.NextID(), USStock(), MarketOrder("SELL", ONE))
	// ib.PlaceOrder(ib.NextID(), USStock(), MarketToLimit("BUY", ONE))
	// ib.PlaceOrder(ib.NextID(), OptionAtIse(), MidpointMatch("BUY", ONE))
	// ib.PlaceOrder(ib.NextID(), USStock(), Stop("SELL", ONE, 34.4))
	// ib.PlaceOrder(ib.NextID(), USStock(), StopLimit("BUY", ONE, 35, 33))
	// ib.PlaceOrder(ib.NextID(), USStock(), StopWithProtection("SELL", ONE, 45))
	// ib.PlaceOrder(ib.NextID(), USStock(), SweepToFill("BUY", ONE, 35))
	// ib.PlaceOrder(ib.NextID(), USStock(), TrailingStop("SELL", ONE, 0.5, 30))
	// ib.PlaceOrder(ib.NextID(), USStock(), TrailingStopLimit("BUY", ONE, 2, 5, 50))

	// mid price
	ib.PlaceOrder(ib.NextID(), USStockAtSmart(), Midprice("BUY", ONE, 150))
	// with cash Qty
	ib.PlaceOrder(ib.NextID(), USStockAtSmart(), LimitOrderWithCashQty("BUY", 111.11, 5000))

	time.Sleep(1 * time.Second)

	// Cancel one order
	ib.CancelOrder(ib.nextID, CancelOrderEmpty())

	// cancel all orders for all accounts
	ib.ReqGlobalCancel(CancelOrderEmpty())

	// request the day's execution
	ib.ReqExecutions(100001, ExecutionFilter{})

	// request completed orders
	ib.ReqCompletedOrders(false)

	// order submission
	ib.PlaceOrder(ib.NextID(), CryptoContract(), LimitOrder("BUY", StringToDecimal("0.12345678"), 3700))

	// order time
	ib.PlaceOrder(ib.NextID(), USStockAtSmart(), LimitOrderWithManualOrderTime("BUY", StringToDecimal("100"), 111.11, "20240714-13:00:00"))
	// Cancel one order
	ib.CancelOrder(ib.nextID, CancelOrderWithManualTime("20240914-00:00:05"))

	// peg best to mid order submission
	ib.PlaceOrder(ib.NextID(), IBKRATSContract(), PegBestUpToMidOrder("BUY", StringToDecimal("100"), 111.11, 100, 200, 0.02, 0.025))

	// peg best order submission
	ib.PlaceOrder(ib.NextID(), IBKRATSContract(), PegBestOrder("BUY", StringToDecimal("100"), 111.11, 100, 200, 0.03))

	// peg mid order submission
	ib.PlaceOrder(ib.NextID(), IBKRATSContract(), PegMidOrder("BUY", StringToDecimal("100"), 111.11, 100, 200, 0.025))

	// limit with customer account order submission
	ib.PlaceOrder(ib.NextID(), USStockAtSmart(), LimitOrderWithCustomerAccount("BUY", StringToDecimal("100"), 111.11, "CustAcct"))

	// limit with include overnight
	ib.PlaceOrder(ib.NextID(), USStockAtSmart(), LimitOrderWithIncludeOvernight("BUY", StringToDecimal("100"), 111.11))

	// limit with CME Tag
	ib.PlaceOrder(ib.NextID(), SimpleFuture(), LimitOrderWithCmeTaggingFields("BUY", StringToDecimal("1"), 5333, "ABCD", 1))
	time.Sleep(5 * time.Second)
	ib.CancelOrder(ib.nextID, OrderCancelWithCmeTaggingFields("BCDE", 0))
	time.Sleep(2 * time.Second)
	ib.PlaceOrder(ib.NextID(), SimpleFuture(), LimitOrderWithCmeTaggingFields("BUY", StringToDecimal("1"), 5333, "CDEF", 0))
	time.Sleep(5 * time.Second)
	ib.CancelOrder(ib.nextID, OrderCancelWithCmeTaggingFields("DEFG", 1))
}

func ocaSamples(ib *EClient) {
	orders := []*Order{}
	orders = append(orders, LimitOrder("BUY", ONE, 10))
	orders = append(orders, LimitOrder("BUY", ONE, 11))
	orders = append(orders, LimitOrder("BUY", ONE, 12))
	for _, order := range orders {
		OneCancelsAll("TestOca", order, 2)
		ib.PlaceOrder(ib.NextID(), USStock(), order)
	}
}

func conditionSamples(ib *EClient) {
	// Order conditioning activate - Order will become active if conditioning criteria is met
	lmt := LimitOrder("BUY", StringToDecimal("100"), 20)
	lmt.Conditions = append(lmt.Conditions, NewPriceCondition(DefaultTriggerMethod, 208813720, "SMART", 600, false, false))
	lmt.Conditions = append(lmt.Conditions, NewExecutionCondition("EUR.USD", "CASH", "IDEALPRO", true))
	lmt.Conditions = append(lmt.Conditions, NewMarginCondition(30, true, false))
	lmt.Conditions = append(lmt.Conditions, NewPercentageChangeCondition(15.0, 208813720, "SMART", true, true))
	lmt.Conditions = append(lmt.Conditions, NewTimeCondition("20220808 10:00:00 US/Eastern", true, false))
	lmt.Conditions = append(lmt.Conditions, NewVolumeCondition(208813720, "SMART", false, 100, true))
	ib.PlaceOrder(ib.NextID(), USStock(), lmt)

	// Conditions can make the order active or cancel it. Only LMT orders can be conditionally canceled.
	lmt2 := LimitOrder("BUY", StringToDecimal("100"), 20)
	lmt2.ConditionsCancelOrder = true
	lmt2.Conditions = append(lmt2.Conditions, NewPriceCondition(DefaultTriggerMethod, 208813720, "SMART", 600, false, false))
	ib.PlaceOrder(ib.NextID(), EuropeanStock(), lmt2)
}

func bracketSample(ib *EClient) {
	parent, takeProfit, stopLoss := BracketOrder(ib.NextID(), "BUY", StringToDecimal("100"), 30, 40, 20)
	ib.PlaceOrder(parent.OrderID, EuropeanStock(), parent)
	ib.PlaceOrder(takeProfit.OrderID, EuropeanStock(), takeProfit)
	ib.PlaceOrder(stopLoss.OrderID, EuropeanStock(), stopLoss)
}

func hedgeSample(ib *EClient) {
	//F Hedge order
	//Parent order on a contract which currency differs from your base currency
	parent := LimitOrder("BUY", StringToDecimal("100"), 10)
	parent.OrderID = ib.NextID()
	parent.Transmit = false
	// Hedge on the currency conversion
	hedge := MarketFHedge(parent.OrderID, "BUY")
	// Place the parent first...
	ib.PlaceOrder(parent.OrderID, EuropeanStock(), parent)
	// Then the hedge order
	ib.PlaceOrder(ib.NextID(), EurGbpFx(), hedge)
}

func testAlgoSamples(ib *EClient) {
	// base order
	baseOrder := LimitOrder("BUY", StringToDecimal("1000"), 1)
	// arrival px
	FillArrivalPriceParams(baseOrder, 0.1, "Aggressive", "09:00:00 US/Eastern", "16:00:00 US/Eastern", true, true)
	ib.PlaceOrder(ib.NextID(), USStockAtSmart(), baseOrder)
	// dark ice
	FillDarkIceParams(baseOrder, 10, "09:00:00 US/Eastern", "16:00:00 US/Eastern", true)
	ib.PlaceOrder(ib.NextID(), USStockAtSmart(), baseOrder)
	// accumulate/distribute - The Time Zone in "startTime" and "endTime" attributes is ignored and always defaulted to GMT
	FillAccumulateDistributeParams(baseOrder, 10, 60, true, true, 1, true, true, "12:00:00", "16:00:00")
	ib.PlaceOrder(ib.NextID(), USStockAtSmart(), baseOrder)
	// twap
	FillTwapParams(baseOrder, "Marketable", "09:00:00 US/Eastern", "16:00:00 US/Eastern", true)
	ib.PlaceOrder(ib.NextID(), USStockAtSmart(), baseOrder)
	// vwap
	FillVwapParams(baseOrder, 0.2, "09:00:00 US/Eastern", "16:00:00 US/Eastern", true, true, true)
	ib.PlaceOrder(ib.NextID(), USStockAtSmart(), baseOrder)
	// balance impact risk
	FillBalanceImpactRiskParams(baseOrder, 0.1, "Aggressive", true)
	ib.PlaceOrder(ib.NextID(), USStockAtSmart(), baseOrder)
	// mini impact
	FillMinImpactParams(baseOrder, 0.3)
	ib.PlaceOrder(ib.NextID(), USStockAtSmart(), baseOrder)
	// adaptive
	FillAdaptiveParams(baseOrder, "Normal")
	ib.PlaceOrder(ib.NextID(), USStockAtSmart(), baseOrder)
	// close price
	FillClosePriceParams(baseOrder, 0.5, "Neutral", "12:00:00 US/Eastern", true)
	ib.PlaceOrder(ib.NextID(), USStockAtSmart(), baseOrder)
	// percentage of volume
	FillPctVolParams(baseOrder, 0.5, "12:00:00 US/Eastern", "14:00:00 US/Eastern", true)
	ib.PlaceOrder(ib.NextID(), USStockAtSmart(), baseOrder)
	// price variant percentage of volume
	FillPriceVariantPctVolParams(baseOrder, 0.1, 0.05, 0.01, 0.2, "12:00:00 US/Eastern", "14:00:00 US/Eastern", true)
	ib.PlaceOrder(ib.NextID(), USStockAtSmart(), baseOrder)
	// size variant percentage of volume
	FillSizeVariantPctVolParams(baseOrder, 0.2, 0.4, "12:00:00 US/Eastern", "14:00:00 US/Eastern", true)
	ib.PlaceOrder(ib.NextID(), USStockAtSmart(), baseOrder)
	// time variant percentage of volume
	FillTimeVariantPctVolParams(baseOrder, 0.2, 0.4, "12:00:00 US/Eastern", "14:00:00 US/Eastern", true)
	ib.PlaceOrder(ib.NextID(), USStockAtSmart(), baseOrder)
	// Jefferies vwap
	FillJefferiesVWAPParams(baseOrder, "10:00:00 US/Eastern", "16:00:00 US/Eastern", 10, 10, "Exclude_Both", 130, 135, 1, 10, "Patience", false, "Midpoint")
	ib.PlaceOrder(ib.NextID(), JefferiesContract(), baseOrder)
	// CSFB Inline
	FillCSFBInlineParams(baseOrder, "10:00:00 US/Eastern", "16:00:00 US/Eastern", "Patient", 10, 20, 100, "Default", false, 40, 100, 100, 35)
	ib.PlaceOrder(ib.NextID(), CSFBContract(), baseOrder)
}

func financialAdvisorOrderSamples(ib *EClient) {
	// FA order on one account
	faOrderOneAccount := MarketOrder("BUY", StringToDecimal("100"))
	// Specify the Account Number directly
	faOrderOneAccount.Account = "DU119915"
	ib.PlaceOrder(ib.NextID(), USStock(), faOrderOneAccount)
	time.Sleep(1 * time.Second)
	// FA order group
	faOrderGroup := LimitOrder("BUY", StringToDecimal("200"), 10)
	faOrderGroup.FAGroup = "MyTestGroup1"
	faOrderGroup.FAMethod = "AvailableEquity"
	ib.PlaceOrder(ib.NextID(), USStock(), faOrderGroup)
	time.Sleep(1 * time.Second)
	// FA order user defined group
	faOrderUserDefinedGroup := LimitOrder("BUY", StringToDecimal("200"), 10)
	faOrderUserDefinedGroup.FAGroup = "MyTestProfile1"
	ib.PlaceOrder(ib.NextID(), USStock(), faOrderUserDefinedGroup)
	time.Sleep(1 * time.Second)
	// model order
	modelOrder := LimitOrder("BUY", StringToDecimal("200"), 100)
	modelOrder.Account = "DF12345"
	modelOrder.ModelCode = "Technology"
	ib.PlaceOrder(ib.NextID(), USStock(), modelOrder)
	time.Sleep(1 * time.Second)
}

func financialAdvisorOperations(ib *EClient) {
	// Requesting FA information
	ib.RequestFA(ALIASES)
	ib.RequestFA(GROUPS)
	// Replacing FA information - Fill in with the appropriate XML string.
	ib.ReplaceFA(1000, GROUPS, FAUpdatedGroup())
	// soft dollar tier
	ib.ReqSoftDollarTiers(4001)
}

func testDisplayGroups(ib *EClient) {
	ib.QueryDisplayGroups(9001)
	time.Sleep(1 * time.Second)
	ib.SubscribeToGroupEvents(9002, 1)
	time.Sleep(1 * time.Second)
	ib.UpdateDisplayGroup(9002, "8314@SMART")
	time.Sleep(1 * time.Second)
	ib.UnsubscribeFromGroupEvents(9002)
}

func miscellaneous(ib *EClient) {
	// Request TWS' current time
	ib.ReqCurrentTime()
	// Setting TWS logging level
	ib.SetServerLogLevel(5)
}

func reqFamilyCodes(ib *EClient) {
	// Request TWS' family codes
	ib.ReqFamilyCodes()
}

func reqMatchingSymbols(ib *EClient) {
	// Request TWS' mathing symbols
	ib.ReqMatchingSymbols(11001, "IBM")
}

func reqMktDepthExchanges(ib *EClient) {
	// Request TWS' market depth exchanges
	ib.ReqMktDepthExchanges()
}

func reqNewsTicks(ib *EClient) {
	ib.ReqMktData(12001, USStockAtSmart(), "mdoff,292", false, false, nil)
	time.Sleep(5 * time.Second)
	ib.CancelMktData(12001)
}

func reqSmartComponents(ib *EClient) {
	// TODO
}

func reqNewsProviders(ib *EClient) {
	// Request TWS' news providers
	ib.ReqNewsProviders()
}

func reqNewsArticle(ib *EClient) {
	// Request TWS' news article
	list := []TagValue{}
	ib.ReqNewsArticle(12001, "MST", "MST$06f53098", list)
}

func reqHistoricalNews(ib *EClient) {
	// Request TWS' historical news
	list := []TagValue{}
	list = append(list, TagValue{Tag: "manual", Value: "1"})
	ib.ReqHistoricalNews(12001, 8314, "BZ+FLY", "", "", 5, list)
	time.Sleep(1 * time.Second)
}

func reqHeadTimestamp(ib *EClient) {
	ib.ReqHeadTimeStamp(14001, EurGbpFx(), "MIDPOINT", true, 1)
	time.Sleep(1 * time.Second)
	ib.CancelHeadTimeStamp(14001)
}

func reqHistogramData(ib *EClient) {
	ib.ReqHistogramData(15001, IBMUSStockAtSmart(), false, "1 weeks")
	time.Sleep(2 * time.Second)
	ib.CancelHistogramData(15001)
}

func rerouteCFDOperations(ib *EClient) {
	ib.ReqMktData(16001, USStockCFD(), "", false, false, nil)
	time.Sleep(1 * time.Second)
	ib.ReqMktData(16002, EuropeanStockCFD(), "", false, false, nil)
	time.Sleep(1 * time.Second)
	ib.ReqMktData(16003, CashCFD(), "", false, false, nil)
	time.Sleep(1 * time.Second)

	ib.ReqMktDepth(16004, USStockCFD(), 10, false, nil)
	time.Sleep(1 * time.Second)
	ib.ReqMktDepth(16004, EuropeanStockCFD(), 10, false, nil)
	time.Sleep(1 * time.Second)
	ib.ReqMktDepth(16004, CashCFD(), 10, false, nil)
	time.Sleep(1 * time.Second)
}

func marketRuleOperations(ib *EClient) {
	ib.ReqContractDetails(17001, IBMBond())
	ib.ReqContractDetails(17002, IBKRStk())
	time.Sleep(2 * time.Second)
	ib.ReqMarketRule(26)
	ib.ReqMarketRule(635)
	ib.ReqMarketRule(1388)
}

func continuousFuturesOperations(ib *EClient) {
	ib.ReqContractDetails(18001, ContFut())
	queryTime := time.Now().Format("20060102-15:04:05")
	ib.ReqHistoricalData(18002, ContFut(), queryTime, "1 Y", "1 month", "TRADES", false, 1, false, nil)
	time.Sleep(10 * time.Second)
	ib.CancelHistoricalData(18002)
}

func reqHistoricalTicks(ib *EClient) {
	ib.ReqHistoricalTicks(19001, IBMUSStockAtSmart(), "20170621 09:38:33 US/Eastern", "", 10, "BID_ASK", true, true, nil)
	ib.ReqHistoricalTicks(19002, IBMUSStockAtSmart(), "20170621 09:38:33 US/Eastern", "", 10, "MIDPOINT", true, true, nil)
	ib.ReqHistoricalTicks(19003, IBMUSStockAtSmart(), "20170621 09:38:33 US/Eastern", "", 10, "TRADES", true, true, nil)
}

func reqTickByTickData(ib *EClient) {
	//  Requesting tick-by-tick data (only refresh)
	ib.ReqTickByTickData(20001, EuropeanStock(), "Last", 0, false)
	ib.ReqTickByTickData(20002, EuropeanStock(), "AllLast", 0, false)
	ib.ReqTickByTickData(20003, EuropeanStock(), "BidAsk", 0, true)
	ib.ReqTickByTickData(20004, EurGbpFx(), "MidPoint", 0, false)
	time.Sleep(10 * time.Second)
	ib.CancelTickByTickData(20001)
	ib.CancelTickByTickData(20002)
	ib.CancelTickByTickData(20003)
	ib.CancelTickByTickData(20004)
	// Requesting tick-by-tick data (historical + refresh)
	ib.ReqTickByTickData(20005, EuropeanStock(), "Last", 10, false)
	ib.ReqTickByTickData(20006, EuropeanStock(), "AllLast", 10, false)
	ib.ReqTickByTickData(20007, EuropeanStock(), "BidAsk", 10, false)
	ib.ReqTickByTickData(200048, EurGbpFx(), "MidPoint", 10, true)
	time.Sleep(10 * time.Second)
	ib.CancelTickByTickData(20005)
	ib.CancelTickByTickData(20006)
	ib.CancelTickByTickData(20007)
	ib.CancelTickByTickData(20008)
}

func whatIfSamples(ib *EClient) {
	// Placing what-if order
	ib.PlaceOrder(ib.NextID(), BondWithCusip(), WhatIfLimitOrder("BUY", StringToDecimal("100"), 20))
}

func ibkratsSample(ib *EClient) {
	ib.PlaceOrder(ib.NextID(), IBKRATSContract(), LimitIBKRATS("BUY", StringToDecimal("100"), 330))
}

func wshCalendarOperations(ib *EClient) {
	// request WSH meta data
	ib.ReqWshMetaData(30001)
	time.Sleep(10 * time.Second)
	ib.CancelWshMetaData(30001)
	// request event data
	wed1 := NewWshEventData()
	wed1.ConID = 8314
	wed1.StartDate = "20220511"
	wed1.TotalLimit = 5
	ib.ReqWshEventData(30002, wed1)
	time.Sleep(3 * time.Second)
	wed2 := NewWshEventData()
	wed2.Filter = "{\"watchlist\":[\"8314\"]}"
	wed2.EndDate = "20220512"
	ib.ReqWshEventData(30002, wed2)
	time.Sleep(10 * time.Second)
	ib.CancelWshEventData(30002)
	ib.CancelWshEventData(30003)
}
