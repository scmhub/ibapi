package ibapi

import (
	"time"
)

// EWrapper contains the function to handle incoming messages from TWS or Gateway
type EWrapper interface {
	// TickPrice handles all price related ticks. Every tickPrice callback is followed by a tickSize.
	// A tickPrice value of -1 or 0 followed by a tickSize of 0 indicates there is no data for this field currently available, whereas a tickPrice with a positive tickSize indicates an active quote of 0 (typically for a combo contract).
	TickPrice(reqID TickerID, tickType TickType, price float64, attrib TickAttrib)
	// TickSize handles all size related ticks.
	TickSize(reqID TickerID, tickType TickType, size Decimal)
	// TickOptionComputation is called when the market in an option or its underlier moves.
	// TWS's option model volatilities, prices, and deltas, along with the present value of dividends expected on that options underlier are received.
	TickOptionComputation(reqID TickerID, tickType TickType, tickAttrib int64, impliedVol float64, delta float64, optPrice float64, pvDividend float64, gamma float64, vega float64, theta float64, undPrice float64)
	// TickGeneric .
	TickGeneric(reqID TickerID, tickType TickType, value float64)
	// TickString .
	TickString(reqID TickerID, tickType TickType, value string)
	// TickEFP handles market for Exchange for Physical.
	// tickerId is the request's identifier.
	// tickType is the type of tick being received.
	// basisPoints is the annualized basis points, which is representative of the financing rate that can be directly compared to broker rates.
	// formattedBasisPoints is the annualized basis points as a formatted string that depicts them in percentage form.
	// impliedFuture is the implied Futures price.
	// holdDays is the number of hold days until the lastTradeDate of the EFP.
	// futureLastTradeDate is the expiration date of the single stock future.
	// dividendImpact is the dividend impact upon the annualized basis points interest rate.
	// dividendsToLastTradeDate is the dividends expected until the expiration of the single stock future.
	TickEFP(reqID TickerID, tickType TickType, basisPoints float64, formattedBasisPoints string, totalDividends float64, holdDays int64, futureLastTradeDate string, dividendImpact float64, dividendsToLastTradeDate float64)
	// OrderStatus is called whenever the status of an order changes.
	// It is also fired after reconnecting to TWS if the client has any open orders.
	// OrderID is the order ID that was specified previously in the	call to placeOrder().
	// status is the order status. Possible values include:
	//		PendingSubmit - indicates that you have transmitted the order, but have not  yet received confirmation that it has been accepted by the order destination. NOTE: This order status is not sent by TWS and should be explicitly set by the API developer when an order is submitted.
	//		PendingCancel - indicates that you have sent a request to cancel the order but have not yet received cancel confirmation from the order destination. At this point, your order is not confirmed canceled. You may still receive an execution while your cancellation request is pending. NOTE: This order status is not sent by TWS and should be explicitly set by the API developer when an order is canceled.
	//		PreSubmitted - indicates that a simulated order type has been accepted by the IB system and that this order has yet to be elected. The order is held in the IB system until the election criteria are met. At that time the order is transmitted to the order destination as specified.
	//		Submitted - indicates that your order has been accepted at the order destination and is working.
	//		Cancelled - indicates that the balance of your order has been confirmed canceled by the IB system. This could occur unexpectedly when IB or the destination has rejected your order.
	//		Filled - indicates that the order has been completely filled.
	//		Inactive - indicates that the order has been accepted by the system (simulated orders) or an exchange (native orders) but that currently the order is inactive due to system, exchange or other issues.
	// filled specifies the number of shares that have been executed. For more information about partial fills, see Order Status for Partial Fills.
	// remaining specifies the number of shares still outstanding.
	// avgFillPrice is the average price of the shares that have been executed. This parameter is valid only if the filled parameter value is greater than zero. Otherwise, the price parameter will be zero.
	// permId is the TWS id used to identify orders. Remains the same over TWS sessions.
	// parentId is the order ID of the parent order, used for bracket and auto trailing stop orders.
	// lastFilledPrice is the last price of the shares that have been executed. This parameter is valid only if the filled parameter value is greater than zero. Otherwise, the price parameter will be zero.
	// clientId is the ID of the client (or TWS) that placed the order. Note that TWS orders have a fixed clientId and OrderID of 0 that distinguishes them from API orders.
	// whyHeld is the field used to identify an order held when TWS is trying to locate shares for a short sell. The value used to indicate this is 'locate'.
	OrderStatus(orderID OrderID, status string, filled Decimal, remaining Decimal, avgFillPrice float64, permID int64, parentID int64, lastFillPrice float64, clientID int64, whyHeld string, mktCapPrice float64)
	// OpenOrder is called to feed in open orders.
	// orderID: OrderId - The order ID assigned by TWS. Use to cancel or update TWS order.
	// contract: Contract - The Contract class attributes describe the contract.
	// order: Order - The Order class gives the details of the open order.
	// orderState: OrderState - The orderState class includes attributes Used for both pre and post trade margin and commission data.
	OpenOrder(orderID OrderID, contract *Contract, order *Order, orderState *OrderState)
	// OpenOrderEnd is called at the end of a given request for open orders.
	OpenOrderEnd()
	// WinError .
	WinError(text string, lastError int64)
	// ConnectionClosed is called when TWS closes the sockets connection with the ActiveX control, or when TWS is shut down.
	ConnectionClosed()
	// UpdateAccountValue is called only when reqAccountUpdates() has been called.
	UpdateAccountValue(tag string, val string, currency string, accountName string)
	// UpdatePortfolio is called only when reqAccountUpdates() has been called.
	UpdatePortfolio(contract *Contract, position Decimal, marketPrice float64, marketValue float64, averageCost float64, unrealizedPNL float64, realizedPNL float64, accountName string)
	// UpdateAccountTime .
	UpdateAccountTime(timeStamp string)
	// AccountDownloadEnd is called after a batch updateAccountValue() and updatePortfolio() is sent.
	AccountDownloadEnd(accountName string)
	// NextValidID Receives next valid order id. NOT THREAD-SAFE.
	NextValidID(reqID int64)
	// ContractDetails Receives the full contract's definitions. This method will return all contracts matching the requested via reqContractDetails().
	// For example, one can obtain the whole option chain with it.
	ContractDetails(reqID int64, contractDetails *ContractDetails)
	// BondContractDetails is called when reqContractDetails function has been called for bonds.
	BondContractDetails(reqID int64, contractDetails *ContractDetails)
	// ContractDetailsEnd is called once all contract details for a given request are received.
	// This helps to define the end of an option chain.
	ContractDetailsEnd(reqID int64)
	// ExecDetails is called when the reqExecutions() functions is invoked, or when an order is filled.
	ExecDetails(reqID int64, contract *Contract, execution *Execution)
	// ExecDetailsEnd is called once all executions have been sent to a client in response to reqExecutions().
	ExecDetailsEnd(reqID int64)
	// Error is called when there is an error with the communication or when TWS wants to send a message to the client.
	Error(reqID TickerID, errTime int64, errCode int64, errString string, advancedOrderRejectJson string)
	// UpdateMktDepth returns the order book.
	// 	TickerID -  the request's identifier.
	// 	position -  the order book's row being updated.
	// 	operation - how to refresh the row:
	// 		0 = insert (insert this new order into the row identified by 'position').
	// 		1 = update (update the existing order in the row identified by 'position').
	// 		2 = delete (delete the existing order at the row identified by 'position').
	// 	side -  0 for ask, 1 for bid.
	// 	price - the order's price.
	// 	size -  the order's size.
	UpdateMktDepth(TickerID TickerID, position int64, operation int64, side int64, price float64, size Decimal)
	// UpdateMktDepthL2 returns the order book.
	// 	TickerID -  the request's identifier.
	//  position -  the order book's row being updated.
	//  marketMaker - the exchange holding the order.
	//  operation - how to refresh the row:
	//  	0 = insert (insert this new order into the row identified by 'position').
	//      1 = update (update the existing order in the row identified by 'position').
	//      2 = delete (delete the existing order at the row identified by 'position').
	//  side -  0 for ask, 1 for bid.
	//  price - the order's price.
	//  size -  the order's size.
	//  isSmartDepth - is SMART Depth request.
	UpdateMktDepthL2(TickerID TickerID, position int64, marketMaker string, operation int64, side int64, price float64, size Decimal, isSmartDepth bool)
	// UpdateNewsBulletin provides IB's bulletins.
	// 	msgID - the bulletin's identifier.
	// 	msgType - one of: 1 - Regular news bulletin 2 - Exchange no longer available for trading 3 - Exchange is available for trading.
	// 	newsMessage - the message.
	// 	originExch -    the exchange where the message comes from.
	UpdateNewsBulletin(msgID int64, msgType int64, newsMessage string, originExch string)
	// ManagedAccounts Receives a comma-separated string with the managed account ids.
	ManagedAccounts(accountsList []string)
	// ReceiveFA receives the Financial Advisor's configuration available in the TWS
	//  faDataType - one of:
	// 		GROUPS: offer traders a way to create a group of accounts and apply a single allocation method to all accounts in the group.
	// 		ALIASES: let you easily identify the accounts by meaningful names rather than account numbers.
	// faXmlData -  the xml-formatted configuration
	ReceiveFA(faDataType FaDataType, cxml string)
	// HistoricalData returns the requested historical data bars
	// reqID - the request's identifier
	// bar - The Bar
	HistoricalData(reqID int64, bar *Bar)
	// HistoricalDataEnd is called when historical bars reception is ending.
	HistoricalDataEnd(reqID int64, startDateStr string, endDateStr string)
	// ScannerParameters Provides the xml-formatted parameters available to create a market scanner.
	ScannerParameters(xml string)
	// ScannerData Provides the data resulting from the market scanner request.
	// reqID - the request's identifier.
	// rank -  the ranking within the response of this bar.
	// contractDetails - the data's ContractDetails
	// distance -      according to query.
	// benchmark -     according to query.
	// projection -    according to query.
	// legStr - describes the combo legs when the scanner is returning EFP
	ScannerData(reqID int64, rank int64, contractDetails *ContractDetails, distance string, benchmark string, projection string, legsStr string)
	// ScannerDataEnd indicates that the scanner data reception has terminated.
	ScannerDataEnd(reqID int64)
	// RealtimeBar updates the real time 5 seconds bars
	// reqID - the request's identifier
	// time  - start of bar in unix (or 'epoch') time
	// open_  - the bar's open value
	// high  - the bar's high value
	// low   - the bar's low value
	// close - the bar's closing value
	// volume - the bar's traded volume if available
	// wap   - the bar's Weighted Average Price
	// count - the number of trades during the bar's timespan (only available for TRADES).
	RealtimeBar(reqID TickerID, time int64, open float64, high float64, low float64, close float64, volume Decimal, wap Decimal, count int64)
	// CurrentTime will receive IB server's system current time after the invokation of reqCurrentTime.
	CurrentTime(t int64)
	// FundamentalData
	FundamentalData(reqID TickerID, data string)
	// DeltaNeutralValidation
	DeltaNeutralValidation(reqID int64, deltaNeutralContract DeltaNeutralContract)
	// TickSnapshotEnd indicates the snapshot reception is finished.
	TickSnapshotEnd(reqID int64)
	// MarketDataType is called when market data switches between real-time and frozen.
	// The marketDataType( ) callback accepts a reqId parameter and is sent per every subscription because different contracts can generally trade on a different schedule
	MarketDataType(reqID TickerID, marketDataType int64)
	// CommissionReport is called immediately after a trade execution or by calling reqExecutions().
	CommissionReport(commissionReport CommissionReport)
	// Position returns real-time positions for all accounts in response to the reqPositions() method.
	Position(account string, contract *Contract, position Decimal, avgCost float64)
	// PositionEnd is called once all position data for a given request are received and functions as an end marker for the position() data.
	PositionEnd()
	// AccountSummary returns the data from the TWS Account Window Summary tab in response to reqAccountSummary().
	AccountSummary(reqID int64, account string, tag string, value string, currency string)
	// AccountSummaryEnd is called once all account summary data for a given request are received.
	AccountSummaryEnd(reqID int64)
	// VerifyMessageAPI .
	VerifyMessageAPI(apiData string)
	// VerifyCompleted .
	VerifyCompleted(isSuccessful bool, errorText string)
	// DisplayGroupList is a one-time response to queryDisplayGroups().
	// reqID - The reqID specified in queryDisplayGroups().
	// groups - A list of integers representing visible group ID separated by the | character, and sorted by most used group first. This list will
	//      not change during TWS session (in other words, user cannot add a new group; sorting can change though).
	DisplayGroupList(reqID int64, groups string)
	// DisplayGroupUpdated .
	DisplayGroupUpdated(reqID int64, contractInfo string)
	// VerifyAndAuthMessageAPI .
	VerifyAndAuthMessageAPI(apiData string, xyzChallange string)
	// VerifyAndAuthCompleted .
	VerifyAndAuthCompleted(isSuccessful bool, errorText string)
	// ConnectAck is called on completion of successful connection.
	ConnectAck()
	// PositionMulti .
	PositionMulti(reqID int64, account string, modelCode string, contract *Contract, pos Decimal, avgCost float64)
	// PositionMultiEnd .
	PositionMultiEnd(reqID int64)
	// AccountUpdateMulti .
	AccountUpdateMulti(reqID int64, account string, modleCode string, key string, value string, currency string)
	// AccountUpdateMultiEnd .
	AccountUpdateMultiEnd(reqID int64)
	// SecurityDefinitionOptionParameter returns the option chain for an underlying on an exchange specified in reqSecDefOptParams.
	// There will be multiple callbacks to securityDefinitionOptionParameter if multiple exchanges are specified in reqSecDefOptParams.
	// reqId - ID of the request initiating the callback.
	// underlyingConId - The conID of the underlying security.
	// tradingClass -  the option trading class.
	// multiplier -    the option multiplier.
	// expirations - a list of the expiries for the options of this underlying on this exchange.
	// strikes - a list of the possible strikes for options of this underlying on this exchange.
	SecurityDefinitionOptionParameter(reqID int64, exchange string, underlyingConID int64, tradingClass string, multiplier string, expirations []string, strikes []float64)
	// SecurityDefinitionOptionParameterEnd is called when all callbacks to securityDefinitionOptionParameter are completed.
	SecurityDefinitionOptionParameterEnd(reqID int64)
	// SoftDollarTiers is called when receives Soft Dollar Tier configuration information.
	// reqID - The request ID used in the call to reqSoftDollarTiers()
	// tiers - Stores a list of SoftDollarTier that contains all Soft Dollar Tiers information
	SoftDollarTiers(reqID int64, tiers []SoftDollarTier)
	// FamilyCodes .
	FamilyCodes(familyCodes []FamilyCode)
	// SymbolSamples .
	SymbolSamples(reqID int64, contractDescriptions []ContractDescription)
	// MktDepthExchanges .
	MktDepthExchanges(depthMktDataDescriptions []DepthMktDataDescription)
	// TickNews .
	TickNews(TickerID TickerID, timeStamp int64, providerCode string, articleID string, headline string, extraData string)
	// SmartComponents .
	SmartComponents(reqID int64, smartComponents []SmartComponent)
	// TickReqParams .
	TickReqParams(TickerID TickerID, minTick float64, bboExchange string, snapshotPermissions int64)
	// NewsProviders .
	NewsProviders(newsProviders []NewsProvider)
	// NewsArticle .
	NewsArticle(requestID int64, articleType int64, articleText string)
	// HistoricalNews returns historical news headlines.
	HistoricalNews(requestID int64, time string, providerCode string, articleID string, headline string)
	// HistoricalNewsEnd signals end of historical news.
	HistoricalNewsEnd(requestID int64, hasMore bool)
	// HeadTimestamp returns earliest available data of a type of data for a particular contract.
	HeadTimestamp(reqID int64, headTimestamp string)
	// HistogramData returns histogram data for a contract.
	HistogramData(reqID int64, data []HistogramData)
	// HistoricalDataUpdate .
	HistoricalDataUpdate(reqID int64, bar *Bar)
	// RerouteMktDataReq .
	RerouteMktDataReq(reqID int64, conID int64, exchange string)
	// RerouteMktDepthReq .
	RerouteMktDepthReq(reqID int64, conID int64, exchange string)
	// MarketRule .
	MarketRule(marketRuleID int64, priceIncrements []PriceIncrement)
	// Pnl returns the daily PnL for the account.
	Pnl(reqID int64, dailyPnL float64, unrealizedPnL float64, realizedPnL float64)
	// PnlSingle returns the daily PnL for a single position in the account.
	PnlSingle(reqID int64, pos Decimal, dailyPnL float64, unrealizedPnL float64, realizedPnL float64, value float64)
	// HistoricalTicks returns historical tick data when whatToShow=MIDPOINT.
	HistoricalTicks(reqID int64, ticks []HistoricalTick, done bool)
	// HistoricalTicksBidAsk returns historical tick data when whatToShow=BID_ASK.
	HistoricalTicksBidAsk(reqID int64, ticks []HistoricalTickBidAsk, done bool)
	// HistoricalTicksLast returns historical tick data when whatToShow=TRADES
	HistoricalTicksLast(reqID int64, ticks []HistoricalTickLast, done bool)
	// TickByTickAllLast returns tick-by-tick data for tickType = "Last" or "AllLast"
	TickByTickAllLast(reqID int64, tickType int64, time int64, price float64, size Decimal, tickAttribLast TickAttribLast, exchange string, specialConditions string)
	// TickByTickBidAsk .
	TickByTickBidAsk(reqID int64, time int64, bidPrice float64, askPrice float64, bidSize Decimal, askSize Decimal, tickAttribBidAsk TickAttribBidAsk)
	// TickByTickMidPoint .
	TickByTickMidPoint(reqID int64, time int64, midPoint float64)
	// OrderBound returns orderBound notification .
	OrderBound(permID int64, clientID int64, orderID int64)
	// CompletedOrder is called to feed in completed orders.
	CompletedOrder(contract *Contract, order *Order, orderState *OrderState)
	// CompletedOrdersEnd is called at the end of a given request for completed orders.
	CompletedOrdersEnd()
	// ReplaceFAEnd is called at the end of a replace FA.
	ReplaceFAEnd(reqID int64, text string)
	// WshMetaData .
	WshMetaData(reqID int64, dataJson string)
	// WshEventData .
	WshEventData(reqID int64, dataJson string)
	// HistoricalSchedule returns historical schedule for historical data request with whatToShow=SCHEDULE
	HistoricalSchedule(reqID int64, startDarteTime, endDateTime, timeZone string, sessions []HistoricalSession)
	// UserInfo returns user info.
	UserInfo(reqID int64, whiteBrandingId string)
}

var _ EWrapper = (*Wrapper)(nil)

// Wrapper is the default implementation of the EWrapper interface.
type Wrapper struct {
}

func (w Wrapper) TickPrice(reqID TickerID, tickType TickType, price float64, attrib TickAttrib) {
	log.Info().Int64("reqID", reqID).Int64("tickType", tickType).Str("price", FloatMaxString(price)).Bool("CanAutoExecute", attrib.CanAutoExecute).Bool("PastLimit", attrib.PastLimit).Bool("PreOpen", attrib.PreOpen).Msg("<TickPrice>")
}

func (w Wrapper) TickSize(reqID TickerID, tickType TickType, size Decimal) {
	log.Info().Int64("reqID", reqID).Int64("tickType", tickType).Str("size", DecimalMaxString(size)).Msg("<TickSize>")
}

func (w Wrapper) TickOptionComputation(reqID TickerID, tickType TickType, tickAttrib int64, impliedVol float64, delta float64, optPrice float64, pvDividend float64, gamma float64, vega float64, theta float64, undPrice float64) {
	log.Info().Int64("reqID", reqID).Int64("tickType", tickType).Str("tickAttrib", IntMaxString(tickAttrib)).Str("impliedVol", FloatMaxString(impliedVol)).Str("delta", FloatMaxString(delta)).Str("optPrice", FloatMaxString(optPrice)).Str("pvDividend", FloatMaxString(pvDividend)).Str("gamma", FloatMaxString(gamma)).Str("vega", FloatMaxString(vega)).Str("theta", FloatMaxString(theta)).Str("undPrice", FloatMaxString(undPrice)).Msg("<TickOptionComputation>")
}

func (w Wrapper) TickGeneric(reqID TickerID, tickType TickType, value float64) {
	log.Info().Int64("reqID", reqID).Int64("tickType", tickType).Str("value", FloatMaxString(value)).Msg("<TickGeneric>")
}

func (w Wrapper) TickString(reqID TickerID, tickType TickType, value string) {
	log.Info().Int64("reqID", reqID).Int64("tickType", tickType).Str("value", value).Msg("<TickString>")
}

func (w Wrapper) TickEFP(reqID TickerID, tickType TickType, basisPoints float64, formattedBasisPoints string, totalDividends float64, holdDays int64, futureLastTradeDate string, dividendImpact float64, dividendsToLastTradeDate float64) {
	log.Info().Int64("reqID", reqID).Int64("tickType", tickType).Float64("basisPoints", basisPoints).Str("formattedBasisPoints", formattedBasisPoints).Float64("totalDividends", totalDividends).Int64("holdDays", holdDays).Str("futureLastTradeDate", futureLastTradeDate).Float64("dividendImpact", dividendImpact).Float64("dividendsToLastTradeDate", dividendsToLastTradeDate).Msg("<TickEFP>")
}

func (w Wrapper) OrderStatus(orderID OrderID, status string, filled Decimal, remaining Decimal, avgFillPrice float64, permID int64, parentID int64, lastFillPrice float64, clientID int64, whyHeld string, mktCapPrice float64) {
	log.Info().Int64("orderID", orderID).Str("status", status).Stringer("filled", filled).Stringer("remaining", remaining).Float64("avgFillPrice", avgFillPrice).Int64("permID", permID).Int64("parentID", parentID).Float64("lastFillPrice", lastFillPrice).Int64("clientID", clientID).Str("whyHeld", whyHeld).Float64("mktCapPrice", mktCapPrice).Msg("<OrderStatus>")
}

func (w Wrapper) OpenOrder(orderID OrderID, contract *Contract, order *Order, orderState *OrderState) {
	log.Info().Int64("orderID", orderID).Stringer("contract", contract).Stringer("order", order).Stringer("orderState", orderState).Msg("<OpenOrder>")
}

func (w Wrapper) OpenOrderEnd() {
	log.Info().Msg("<OpenOrderEnd>")
}

func (w Wrapper) WinError(text string, lastError int64) {
	log.Info().Str("text", text).Int64("lastError", lastError).Msg("<WinError>")
}

func (w Wrapper) ConnectionClosed() {
	log.Info().Msg("<ConnectionClosed>...")
}

func (w Wrapper) UpdateAccountValue(tag string, value string, currency string, accountName string) {
	log.Info().Str("tag", tag).Str("value", value).Str("currency", currency).Str("accountName", accountName).Msg("<UpdateAccountValue>")
}

func (w Wrapper) UpdatePortfolio(contract *Contract, position Decimal, marketPrice float64, marketValue float64, averageCost float64, unrealizedPNL float64, realizedPNL float64, accountName string) {
	log.Info().Str("Symbol", contract.Symbol).Str("secType", contract.SecType).Str("exchange", contract.Exchange).Discard().Str("position", DecimalMaxString(position)).Str("marketPrice", FloatMaxString(marketPrice)).Str("marketValue", FloatMaxString(marketValue)).Str("averageCost", FloatMaxString(averageCost)).Str("unrealizedPNL", FloatMaxString(unrealizedPNL)).Str("realizedPNL", FloatMaxString(realizedPNL)).Str("accountName", accountName).Msg("<UpdatePortfolio>")
}

func (w Wrapper) UpdateAccountTime(timeStamp string) {
	log.Info().Str("timeStamp", timeStamp).Msg("<UpdateAccountTime>")
}

func (w Wrapper) AccountDownloadEnd(accountName string) {
	log.Info().Str("accountName", accountName).Msg("<AccountDownloadEnd>")
}

func (w Wrapper) NextValidID(reqID int64) {
	log.Info().Int64("reqID", reqID).Msg("<NextValidID>")
}

func (w Wrapper) ContractDetails(reqID int64, contractDetails *ContractDetails) {
	log.Info().Int64("reqID", reqID).Stringer("contractDetails", contractDetails).Msg("<ContractDetails>")
}

func (w Wrapper) BondContractDetails(reqID int64, contractDetails *ContractDetails) {
	log.Info().Int64("reqID", reqID).Stringer("contractDetails", contractDetails).Msg("<BondContractDetails>")
}

func (w Wrapper) ContractDetailsEnd(reqID int64) {
	log.Info().Int64("reqID", reqID).Msg("<ContractDetailsEnd>")
}

func (w Wrapper) ExecDetails(reqID int64, contract *Contract, execution *Execution) {
	log.Info().Int64("reqID", reqID).Stringer("contract", contract).Stringer("execution", execution).Msg("<ExecDetails>")
}

func (w Wrapper) ExecDetailsEnd(reqID int64) {
	log.Info().Int64("reqID", reqID).Msg("<ExecDetailsEnd>")
}

func (w Wrapper) Error(reqID TickerID, errorTime int64, errCode int64, errString string, advancedOrderRejectJson string) {
	logger := log.Error().Int64("reqID", reqID).Int64("errorTime", errorTime).Int64("errCode", errCode).Str("errString", errString)
	if advancedOrderRejectJson != "" {
		logger = logger.Str("advancedOrderRejectJson", advancedOrderRejectJson)
	}
	logger.Msg("<Error>")
}

func (w Wrapper) UpdateMktDepth(TickerID TickerID, position int64, operation int64, side int64, price float64, size Decimal) {
	log.Info().Int64("TickerID", TickerID).Int64("position", position).Int64("operation", operation).Int64("side", side).Str("price", FloatMaxString(price)).Str("size", DecimalMaxString(size)).Msg("<UpdateMktDepth>")
}

func (w Wrapper) UpdateMktDepthL2(TickerID TickerID, position int64, marketMaker string, operation int64, side int64, price float64, size Decimal, isSmartDepth bool) {
	log.Info().Int64("TickerID", TickerID).Int64("position", position).Str("marketMaker", marketMaker).Int64("operation", operation).Int64("side", side).Str("price", FloatMaxString(price)).Str("size", DecimalMaxString(size)).Bool("isSmartDepth", isSmartDepth).Msg("<UpdateMktDepthL2>")
}

func (w Wrapper) UpdateNewsBulletin(msgID int64, msgType int64, newsMessage string, originExch string) {
	log.Info().Int64("msgID", msgID).Int64("msgType", msgType).Str("newsMessage", newsMessage).Str("originExch", originExch).Msg("<UpdateNewsBulletin>")
}

func (w Wrapper) ManagedAccounts(accountsList []string) {
	log.Info().Strs("accountsList", accountsList).Msg("<ManagedAccounts>")
}

func (w Wrapper) ReceiveFA(faDataType FaDataType, cxml string) {
	log.Info().Stringer("faDataType", faDataType).Str("cxml", cxml).Msg("<ReceiveFA>")
}

func (w Wrapper) HistoricalData(reqID int64, bar *Bar) {
	log.Info().Int64("reqID", reqID).Stringer("bar", bar).Msg("<HistoricalData>")
}

func (w Wrapper) HistoricalDataEnd(reqID int64, startDateStr string, endDateStr string) {
	log.Info().Int64("reqID", reqID).Str("startDateStr", startDateStr).Str("endDateStr", endDateStr).Msg("<HistoricalDataEnd>")
}

func (w Wrapper) ScannerParameters(xml string) {
	log.Info().Str("xml", xml[:50]).Msg("<ScannerParameters>")
}

func (w Wrapper) ScannerData(reqID int64, rank int64, contractDetails *ContractDetails, distance string, benchmark string, projection string, legsStr string) {
	log.Info().Int64("reqID", reqID).Int64("rank", rank).Stringer("contractDetails", contractDetails).Str("distance", distance).Str("benchmark", benchmark).Str("projection", projection).Str("legsStr", legsStr).Msg("<ScannerData>")
}

func (w Wrapper) ScannerDataEnd(reqID int64) {
	log.Info().Int64("reqID", reqID).Msg("<ScannerDataEnd>")
}

func (w Wrapper) RealtimeBar(reqID int64, time int64, open float64, high float64, low float64, close float64, volume Decimal, wap Decimal, count int64) {
	log.Info().Int64("reqID", reqID).Int64("bar time", time).Float64("open", open).Float64("high", high).Float64("low", low).Float64("close", close).Stringer("volume", volume).Stringer("wap", wap).Int64("count", count).Msg("<RealtimeBar>")
}

func (w Wrapper) CurrentTime(t int64) {
	log.Info().Time("Server Time", time.Unix(t, 0)).Msg("<CurrentTime>")
}

func (w Wrapper) FundamentalData(reqID int64, data string) {
	log.Info().Int64("reqID", reqID).Str("data", data).Msg("<FundamentalData>")
}

func (w Wrapper) DeltaNeutralValidation(reqID int64, deltaNeutralContract DeltaNeutralContract) {
	log.Info().Int64("reqID", reqID).Stringer("deltaNeutralContract", deltaNeutralContract).Msg("<DeltaNeutralValidation>")
}

func (w Wrapper) TickSnapshotEnd(reqID int64) {
	log.Info().Int64("reqID", reqID).Msg("<TickSnapshotEnd>")
}

func (w Wrapper) MarketDataType(reqID int64, marketDataType int64) {
	log.Info().Int64("reqID", reqID).Int64("marketDataType", marketDataType).Msg("<MarketDataType>")
}

func (w Wrapper) CommissionReport(commissionReport CommissionReport) {
	log.Info().Stringer("commissionReport", commissionReport).Msg("<CommissionReport>")
}

func (w Wrapper) Position(account string, contract *Contract, position Decimal, avgCost float64) {
	log.Info().Str("account", account).Stringer("contract", contract).Str("position", DecimalMaxString(position)).Str("avgCost", FloatMaxString(avgCost)).Msg("<Position>")
}

func (w Wrapper) PositionEnd() {
	log.Info().Msg("<PositionEnd>")
}

func (w Wrapper) AccountSummary(reqID int64, account string, tag string, value string, currency string) {
	log.Info().Int64("reqID", reqID).Str("account", account).Str("tag", tag).Str("value", value).Str("currency", currency).Msg("<AccountSummary>")
}

func (w Wrapper) AccountSummaryEnd(reqID int64) {
	log.Info().Int64("reqID", reqID).Msg("<AccountSummaryEnd>")
}

func (w Wrapper) VerifyMessageAPI(apiData string) {
	log.Info().Str("apiData", apiData).Msg("<VerifyMessageAPI>")
}

func (w Wrapper) VerifyCompleted(isSuccessful bool, errorText string) {
	log.Info().Bool("isSuccessful", isSuccessful).Str("errorText", errorText).Msg("<VerifyCompleted>")
}

func (w Wrapper) DisplayGroupList(reqID int64, groups string) {
	log.Info().Int64("reqID", reqID).Str("groups", groups).Msg("<DisplayGroupList>")
}

func (w Wrapper) DisplayGroupUpdated(reqID int64, contractInfo string) {
	log.Info().Int64("reqID", reqID).Str("contractInfo", contractInfo).Msg("<DisplayGroupUpdated>")
}

func (w Wrapper) VerifyAndAuthMessageAPI(apiData string, xyzChallange string) {
	log.Info().Str("apiData", apiData).Str("xyzChallange", xyzChallange).Msg("<VerifyAndAuthMessageAPI>")
}

func (w Wrapper) VerifyAndAuthCompleted(isSuccessful bool, errorText string) {
	log.Info().Bool("isSuccessful", isSuccessful).Str("errorText", errorText).Msg("<VerifyAndAuthCompleted>")
}

func (w Wrapper) ConnectAck() {
	log.Info().Msg("<ConnectAck>...")
}

func (w Wrapper) PositionMulti(reqID int64, account string, modelCode string, contract *Contract, pos Decimal, avgCost float64) {
	log.Info().Int64("reqID", reqID).Str("account", account).Str("modelCode", modelCode).Stringer("contract", contract).Str("position", DecimalMaxString(pos)).Str("avgCost", FloatMaxString(avgCost)).Msg("<PositionMulti>")
}

func (w Wrapper) PositionMultiEnd(reqID int64) {
	log.Info().Int64("reqID", reqID).Msg("<PositionMultiEnd>")
}

func (w Wrapper) AccountUpdateMulti(reqID int64, account string, modelCode string, key string, value string, currency string) {
	log.Info().Int64("reqID", reqID).Str("account", account).Str("modelCode", modelCode).Str("key", key).Str("value", value).Str("currency", currency).Msg("<AccountUpdateMulti>")
}

func (w Wrapper) AccountUpdateMultiEnd(reqID int64) {
	log.Info().Int64("reqID", reqID).Msg("<AccountUpdateMultiEnd>")
}

func (w Wrapper) SecurityDefinitionOptionParameter(reqID int64, exchange string, underlyingConID int64, tradingClass string, multiplier string, expirations []string, strikes []float64) {
	log.Info().Int64("reqID", reqID).Str("exchange", exchange).Str("underlyingConID", IntMaxString(underlyingConID)).Str("tradingClass", tradingClass).Str("multiplier", multiplier).Strs("expirations", expirations).Floats64("strikes", strikes).Msg("<SecurityDefinitionOptionParameter>")
}

func (w Wrapper) SecurityDefinitionOptionParameterEnd(reqID int64) {
	log.Info().Int64("reqID", reqID).Msg("<SecurityDefinitionOptionParameterEnd>")
}

func (w Wrapper) SoftDollarTiers(reqID int64, tiers []SoftDollarTier) {
	for _, sdt := range tiers {
		log.Info().Int64("reqID", reqID).Stringer("softDollarTier", sdt).Msg("<SoftDollarTiers>")
	}
}

func (w Wrapper) FamilyCodes(familyCodes []FamilyCode) {
	for _, fc := range familyCodes {
		log.Info().Stringer("familyCode", fc).Msg("<FamilyCodes>")
	}
}

func (w Wrapper) SymbolSamples(reqID int64, contractDescriptions []ContractDescription) {
	log.Info().Int("nb_samples", len(contractDescriptions)).Int64("reqID", reqID).Msg("<SymbolSamples>")
	for i, cd := range contractDescriptions {
		log.Info().Stringer("contract", cd.Contract).Msgf("<Sample %v>", i)
	}
}

func (w Wrapper) MktDepthExchanges(depthMktDataDescriptions []DepthMktDataDescription) {
	log.Info().Any("depthMktDataDescriptions", depthMktDataDescriptions).Msg("<MktDepthExchanges>")
}

func (w Wrapper) TickNews(TickerID TickerID, timeStamp int64, providerCode string, articleID string, headline string, extraData string) {
	log.Info().Int64("TickerID", TickerID).Str("timeStamp", IntMaxString(timeStamp)).Str("providerCode", providerCode).Str("articleID", articleID).Str("headline", headline).Str("extraData", extraData).Msg("<TickNews>")
}

func (w Wrapper) SmartComponents(reqID int64, smartComponents []SmartComponent) {
	log.Info().Int64("reqID", reqID).Msg("<SmartComponents>")
	for i, sc := range smartComponents {
		log.Info().Stringer("smartComponent", sc).Msgf("<Sample %v>", i)
	}
}

func (w Wrapper) TickReqParams(TickerID TickerID, minTick float64, bboExchange string, snapshotPermissions int64) {
	log.Info().Int64("TickerID", TickerID).Str("minTick", FloatMaxString(minTick)).Str("bboExchange", bboExchange).Str("snapshotPermissions", IntMaxString(snapshotPermissions)).Msg("<TickReqParams>")
}

func (w Wrapper) NewsProviders(newsProviders []NewsProvider) {
	for _, np := range newsProviders {
		log.Info().Stringer("newsProvider", np).Msg("<NewsProviders>")
	}
}

func (w Wrapper) NewsArticle(requestID int64, articleType int64, articleText string) {
	log.Info().Int64("requestID", requestID).Int64("articleType", articleType).Str("articleText", articleText).Msg("<NewsArticle>")
}

func (w Wrapper) HistoricalNews(requestID int64, time string, providerCode string, articleID string, headline string) {
	log.Info().Int64("requestID", requestID).Str("news time", time).Str("providerCode", providerCode).Str("providerCode", providerCode).Str("headline", headline).Msg("<HistoricalNews>")
}

func (w Wrapper) HistoricalNewsEnd(requestID int64, hasMore bool) {
	log.Info().Int64("requestID", requestID).Bool("hasMore", hasMore).Msg("<HistoricalNewsEnd>")
}

func (w Wrapper) HeadTimestamp(reqID int64, headTimestamp string) {
	log.Info().Int64("reqID", reqID).Str("headTimestamp", headTimestamp).Msg("<HeadTimestamp>")
}

func (w Wrapper) HistogramData(reqID int64, data []HistogramData) {
	log.Info().Int64("reqID", reqID).Any("data", data).Msg("<HistogramData>")
}

func (w Wrapper) HistoricalDataUpdate(reqID int64, bar *Bar) {
	log.Info().Int64("reqID", reqID).Stringer("bar", bar).Msg("<HistoricalDataUpdate>")
}

func (w Wrapper) RerouteMktDataReq(reqID int64, conID int64, exchange string) {
	log.Info().Int64("reqID", reqID).Int64("conID", conID).Str("exchange", exchange).Msg("<RerouteMktDataReq>")
}

func (w Wrapper) RerouteMktDepthReq(reqID int64, conID int64, exchange string) {
	log.Info().Int64("reqID", reqID).Int64("conID", conID).Str("exchange", exchange).Msg("<RerouteMktDepthReq>")
}

func (w Wrapper) MarketRule(marketRuleID int64, priceIncrements []PriceIncrement) {
	log.Info().Int64("marketRuleID", marketRuleID).Any("priceIncrements", priceIncrements).Msg("<MarketRule>")
}

func (w Wrapper) Pnl(reqID int64, dailyPnL float64, unrealizedPnL float64, realizedPnL float64) {
	log.Info().Int64("reqID", reqID).Str("dailyPnL", FloatMaxString(dailyPnL)).Str("unrealizedPnL", FloatMaxString(unrealizedPnL)).Str("realizedPnL", FloatMaxString(realizedPnL)).Msg("<Pnl>")
}

func (w Wrapper) PnlSingle(reqID int64, pos Decimal, dailyPnL float64, unrealizedPnL float64, realizedPnL float64, value float64) {
	log.Info().Int64("reqID", reqID).Str("position", DecimalMaxString(pos)).Str("dailyPnL", FloatMaxString(dailyPnL)).Str("unrealizedPnL", FloatMaxString(unrealizedPnL)).Str("realizedPnL", FloatMaxString(realizedPnL)).Str("value", FloatMaxString(value)).Msg("<PnlSingle>")
}

func (w Wrapper) HistoricalTicks(reqID int64, ticks []HistoricalTick, done bool) {
	log.Info().Int64("reqID", reqID).Bool("done", done).Any("ticks", ticks).Msg("<HistoricalTicks>")
}

func (w Wrapper) HistoricalTicksBidAsk(reqID int64, ticks []HistoricalTickBidAsk, done bool) {
	log.Info().Int64("reqID", reqID).Bool("done", done).Any("ticks", ticks).Msg("<HistoricalTicksBidAsk>")
}

func (w Wrapper) HistoricalTicksLast(reqID int64, ticks []HistoricalTickLast, done bool) {
	log.Info().Int64("reqID", reqID).Bool("done", done).Any("ticks", ticks).Msg("<HistoricalTicksLast>")
}

func (w Wrapper) TickByTickAllLast(reqID int64, tickType int64, time int64, price float64, size Decimal, tickAttribLast TickAttribLast, exchange string, specialConditions string) {
	log.Info().Int64("reqID", reqID).Int64("tickType", tickType).Int64("tick time", time).Str("price", FloatMaxString(price)).Str("size", DecimalMaxString(size)).Bool("PastLimit", tickAttribLast.PastLimit).Bool("Unreported", tickAttribLast.Unreported).Str("exchange", exchange).Str("specialConditions", specialConditions).Msg("<TickByTickAllLast>")
}

func (w Wrapper) TickByTickBidAsk(reqID int64, time int64, bidPrice float64, askPrice float64, bidSize Decimal, askSize Decimal, tickAttribBidAsk TickAttribBidAsk) {
	log.Info().Int64("reqID", reqID).Int64("tick time", time).Str("bidPrice", FloatMaxString(bidPrice)).Str("askPrice", FloatMaxString(askPrice)).Str("bidSize", DecimalMaxString(bidSize)).Str("askSize", DecimalMaxString(askSize)).Bool("AskPastHigh", tickAttribBidAsk.AskPastHigh).Bool("BidPastLow", tickAttribBidAsk.BidPastLow).Msg("<TickByTickBidAsk>")
}

func (w Wrapper) TickByTickMidPoint(reqID int64, time int64, midPoint float64) {
	log.Info().Int64("reqID", reqID).Int64("tick time", time).Str("midPoint", FloatMaxString(midPoint)).Msg("<TickByTickMidPoint>")
}

func (w Wrapper) OrderBound(permID int64, clientID int64, orderID int64) {
	log.Info().Str("permID", LongMaxString(permID)).Str("clientID", IntMaxString(clientID)).Str("OrderID", IntMaxString(orderID)).Msg("<OrderBound>")
}

func (w Wrapper) CompletedOrder(contract *Contract, order *Order, orderState *OrderState) {
	logger := log.Info().Str("account", order.Account).Str("PermID", LongMaxString(order.PermID)).Str("parentPermID", LongMaxString(order.ParentPermID)).Str("symbol", contract.Symbol).Str("secType", contract.SecType).Str("exchange", contract.Exchange).Str("action", order.Action).Str("orderType", order.OrderType).Str("totalQuantity", DecimalMaxString(order.TotalQuantity))
	logger = logger.Str("cashQty", FloatMaxString(order.CashQty)).Str("filledQuantity", DecimalMaxString(order.FilledQuantity)).Str("lmtPrice", FloatMaxString(order.LmtPrice)).Str("auxPrice", FloatMaxString(order.AuxPrice)).Str("Status", orderState.Status)
	logger = logger.Str("completedTime", orderState.CompletedTime).Str("CompletedStatus", orderState.CompletedStatus).Str("MinTradeQty", IntMaxString(order.MinTradeQty)).Str("MinCompeteSize", IntMaxString(order.MinCompeteSize))
	logger.Msg("<CompletedOrder>")
}

func (w Wrapper) CompletedOrdersEnd() {
	log.Info().Msg("<CompletedOrdersEnd>")
}

func (w Wrapper) ReplaceFAEnd(reqID int64, text string) {
	log.Info().Int64("reqID", reqID).Str("text", text).Msg("<ReplaceFAEnd>")
}

func (w Wrapper) WshMetaData(reqID int64, dataJson string) {
	log.Info().Int64("reqID", reqID).Str("dataJson", dataJson).Msg("<WshMetaData>")
}

func (w Wrapper) WshEventData(reqID int64, dataJson string) {
	log.Info().Int64("reqID", reqID).Str("dataJson", dataJson).Msg("<WshEventData>")
}

func (w Wrapper) HistoricalSchedule(reqID int64, startDarteTime, endDateTime, timeZone string, sessions []HistoricalSession) {
	log.Info().Int64("reqID", reqID).Str("startDarteTime", startDarteTime).Str("endDateTime", endDateTime).Str("timeZone", timeZone).Msg("<HistoricalSchedule>")
}

func (w Wrapper) UserInfo(reqID int64, whiteBrandingId string) {
	log.Info().Int64("reqID", reqID).Str("whiteBrandingId", whiteBrandingId).Msg("<UserInfo>")
}
