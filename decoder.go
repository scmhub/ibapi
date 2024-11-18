/*
The EDecoder knows how to transform a message's payload into higher level IB message (eg: order info, mkt data, etc).
It will call the corresponding method from the EWrapper so that customer's code (eg: class derived from EWrapper) can make further use of the data.
*/
package ibapi

import (
	"strings"
)

// EDecoder transforms a message's payload into higher level IB message.
type EDecoder struct {
	wrapper       EWrapper
	serverVersion Version
}

func (d *EDecoder) interpret(msgBytes []byte) {

	msgBuf := NewMsgBuffer(msgBytes)
	log.Warn()
	if msgBuf.Len() == 0 {
		log.Debug().Msg("no fields")
		return
	}

	// read the msg type
	msgID := msgBuf.decodeInt64()

	switch msgID {
	case TICK_PRICE:
		d.processTickPriceMsg(msgBuf)
	case TICK_SIZE:
		d.processTickSizeMsg(msgBuf)
	case TICK_OPTION_COMPUTATION:
		d.processTickOptionComputationMsg(msgBuf)
	case TICK_GENERIC:
		d.processTickGenericMsg(msgBuf)
	case TICK_STRING:
		d.processTickStringMsg(msgBuf)
	case TICK_EFP:
		d.processTickEfpMsg(msgBuf)
	case ORDER_STATUS:
		d.processOrderStatusMsg(msgBuf)
	case ERR_MSG:
		d.processErrMsg(msgBuf)
	case OPEN_ORDER:
		d.processOpenOrderMsg(msgBuf)
	case ACCT_VALUE:
		d.processAcctValueMsg(msgBuf)
	case PORTFOLIO_VALUE:
		d.processPortfolioValueMsg(msgBuf)
	case ACCT_UPDATE_TIME:
		d.processAcctUpdateTimeMsg(msgBuf)
	case NEXT_VALID_ID:
		d.processNextValidIdMsg(msgBuf)
	case CONTRACT_DATA:
		d.processContractDataMsg(msgBuf)
	case BOND_CONTRACT_DATA:
		d.processBondContractDataMsg(msgBuf)
	case EXECUTION_DATA:
		d.processExecutionDetailsMsg(msgBuf)
	case MARKET_DEPTH:
		d.processMarketDepthMsg(msgBuf)
	case MARKET_DEPTH_L2:
		d.processMarketDepthL2Msg(msgBuf)
	case NEWS_BULLETINS:
		d.processNewsBulletinsMsg(msgBuf)
	case MANAGED_ACCTS:
		d.processManagedAcctsMsg(msgBuf)
	case RECEIVE_FA:
		d.processReceiveFaMsg(msgBuf)
	case HISTORICAL_DATA:
		d.processHistoricalDataMsg(msgBuf)
	case SCANNER_DATA:
		d.processScannerDataMsg(msgBuf)
	case SCANNER_PARAMETERS:
		d.processScannerParametersMsg(msgBuf)
	case CURRENT_TIME:
		d.processCurrentTimeMsg(msgBuf)
	case REAL_TIME_BARS:
		d.processRealTimeBarsMsg(msgBuf)
	case FUNDAMENTAL_DATA:
		d.processFundamentalDataMsg(msgBuf)
	case CONTRACT_DATA_END:
		d.processContractDataEndMsg(msgBuf)
	case OPEN_ORDER_END:
		d.processOpenOrderEndMsg(msgBuf)
	case ACCT_DOWNLOAD_END:
		d.processAcctDownloadEndMsg(msgBuf)
	case EXECUTION_DATA_END:
		d.processExecutionDetailsEndMsg(msgBuf)
	case DELTA_NEUTRAL_VALIDATION:
		d.processDeltaNeutralValidationMsg(msgBuf)
	case TICK_SNAPSHOT_END:
		d.processTickSnapshotEndMsg(msgBuf)
	case MARKET_DATA_TYPE:
		d.processMarketDataTypeMsg(msgBuf)
	case COMMISSION_REPORT:
		d.processCommissionReportMsg(msgBuf)
	case POSITION_DATA:
		d.processPositionDataMsg(msgBuf)
	case POSITION_END:
		d.processPositionEndMsg(msgBuf)
	case ACCOUNT_SUMMARY:
		d.processAccountSummaryMsg(msgBuf)
	case ACCOUNT_SUMMARY_END:
		d.processAccountSummaryEndMsg(msgBuf)
	case VERIFY_MESSAGE_API:
		d.processVerifyMessageApiMsg(msgBuf)
	case VERIFY_COMPLETED:
		d.processVerifyCompletedMsg(msgBuf)
	case DISPLAY_GROUP_LIST:
		d.processDisplayGroupListMsg(msgBuf)
	case DISPLAY_GROUP_UPDATED:
		d.processDisplayGroupUpdatedMsg(msgBuf)
	case VERIFY_AND_AUTH_MESSAGE_API:
		d.processVerifyAndAuthMessageApiMsg(msgBuf)
	case VERIFY_AND_AUTH_COMPLETED:
		d.processVerifyAndAuthCompletedMsg(msgBuf)
	case POSITION_MULTI:
		d.processPositionMultiMsg(msgBuf)
	case POSITION_MULTI_END:
		d.processPositionMultiEndMsg(msgBuf)
	case ACCOUNT_UPDATE_MULTI:
		d.processAccountUpdateMultiMsg(msgBuf)
	case ACCOUNT_UPDATE_MULTI_END:
		d.processAccountUpdateMultiEndMsg(msgBuf)
	case SECURITY_DEFINITION_OPTION_PARAMETER:
		d.processSecurityDefinitionOptionalParameterMsg(msgBuf)
	case SECURITY_DEFINITION_OPTION_PARAMETER_END:
		d.processSecurityDefinitionOptionalParameterEndMsg(msgBuf)
	case SOFT_DOLLAR_TIERS:
		d.processSoftDollarTiersMsg(msgBuf)
	case FAMILY_CODES:
		d.processFamilyCodesMsg(msgBuf)
	case SMART_COMPONENTS:
		d.processSmartComponentsMsg(msgBuf)
	case TICK_REQ_PARAMS:
		d.processTickReqParamsMsg(msgBuf)
	case SYMBOL_SAMPLES:
		d.processSymbolSamplesMsg(msgBuf)
	case MKT_DEPTH_EXCHANGES:
		d.processMktDepthExchangesMsg(msgBuf)
	case TICK_NEWS:
		d.processTickNewsMsg(msgBuf)
	case NEWS_PROVIDERS:
		d.processNewsProvidersMsg(msgBuf)
	case NEWS_ARTICLE:
		d.processNewsArticleMsg(msgBuf)
	case HISTORICAL_NEWS:
		d.processHistoricalNewsMsg(msgBuf)
	case HISTORICAL_NEWS_END:
		d.processHistoricalNewsEndMsg(msgBuf)
	case HEAD_TIMESTAMP:
		d.processHeadTimestampMsg(msgBuf)
	case HISTOGRAM_DATA:
		d.processHistogramDataMsg(msgBuf)
	case HISTORICAL_DATA_UPDATE:
		d.processHistoricalDataUpdateMsg(msgBuf)
	case REROUTE_MKT_DATA_REQ:
		d.processRerouteMktDataReqMsg(msgBuf)
	case REROUTE_MKT_DEPTH_REQ:
		d.processRerouteMktDepthReqMsg(msgBuf)
	case MARKET_RULE:
		d.processMarketRuleMsg(msgBuf)
	case PNL:
		d.processPnLMsg(msgBuf)
	case PNL_SINGLE:
		d.processPnLSingleMsg(msgBuf)
	case HISTORICAL_TICKS:
		d.processHistoricalTicks(msgBuf)
	case HISTORICAL_TICKS_BID_ASK:
		d.processHistoricalTicksBidAsk(msgBuf)
	case HISTORICAL_TICKS_LAST:
		d.processHistoricalTicksLast(msgBuf)
	case TICK_BY_TICK:
		d.processTickByTickDataMsg(msgBuf)
	case ORDER_BOUND:
		d.processOrderBoundMsg(msgBuf)
	case COMPLETED_ORDER:
		d.processCompletedOrderMsg(msgBuf)
	case COMPLETED_ORDERS_END:
		d.processCompletedOrdersEndMsg(msgBuf)
	case REPLACE_FA_END:
		d.processReplaceFAEndMsg(msgBuf)
	case WSH_META_DATA:
		d.processWshMetaData(msgBuf)
	case WSH_EVENT_DATA:
		d.processWshEventData(msgBuf)
	case HISTORICAL_SCHEDULE:
		d.processHistoricalSchedule(msgBuf)
	case USER_INFO:
		d.processUserInfo(msgBuf)
	default:
		d.wrapper.Error(NO_VALID_ID, currentTimeMillis(), BAD_MESSAGE.Code, BAD_MESSAGE.Msg, "")
	}
}

func (d *EDecoder) processTickPriceMsg(msgBuf *MsgBuffer) {

	_ = msgBuf.decodeString()

	reqID := msgBuf.decodeInt64()
	tickType := msgBuf.decodeInt64()
	price := msgBuf.decodeFloat64()
	size := msgBuf.decodeDecimal()   // ver 2 field
	attrMask := msgBuf.decodeInt64() // ver 3 field

	attrib := NewTickAttrib()
	attrib.CanAutoExecute = attrMask == 1

	if d.serverVersion >= MIN_SERVER_VER_PAST_LIMIT {
		attrib.CanAutoExecute = attrMask&0x1 != 0
		attrib.PastLimit = attrMask&0x2 != 0
		if d.serverVersion >= MIN_SERVER_VER_PRE_OPEN_BID_ASK {
			attrib.PreOpen = attrMask&0x4 != 0
		}
	}

	d.wrapper.TickPrice(reqID, tickType, price, attrib)

	var sizeTickType int64
	switch tickType {
	case BID:
		sizeTickType = BID_SIZE
	case ASK:
		sizeTickType = ASK_SIZE
	case LAST:
		sizeTickType = LAST_SIZE
	case DELAYED_BID:
		sizeTickType = DELAYED_BID_SIZE
	case DELAYED_ASK:
		sizeTickType = DELAYED_ASK_SIZE
	case DELAYED_LAST:
		sizeTickType = DELAYED_LAST_SIZE
	default:
		sizeTickType = NOT_SET
	}

	if sizeTickType != NOT_SET {
		d.wrapper.TickSize(reqID, sizeTickType, size)
	}
}

func (d *EDecoder) processTickSizeMsg(msgBuf *MsgBuffer) {

	_ = msgBuf.decodeString()

	reqID := msgBuf.decodeInt64()
	sizeTickType := msgBuf.decodeInt64()
	size := msgBuf.decodeDecimal()

	if sizeTickType != NOT_SET {
		d.wrapper.TickSize(reqID, sizeTickType, size)
	}
}

func (d *EDecoder) processTickOptionComputationMsg(msgBuf *MsgBuffer) {
	optPrice := UNSET_FLOAT
	pvDividend := UNSET_FLOAT
	gamma := UNSET_FLOAT
	vega := UNSET_FLOAT
	theta := UNSET_FLOAT
	undPrice := UNSET_FLOAT

	version := d.serverVersion
	if d.serverVersion < MIN_SERVER_VER_PRICE_BASED_VOLATILITY {
		version = Version(msgBuf.decodeInt64())
	}

	reqID := msgBuf.decodeInt64()
	tickType := msgBuf.decodeInt64()

	var tickAttrib int64
	if d.serverVersion >= MIN_SERVER_VER_PRICE_BASED_VOLATILITY {
		tickAttrib = msgBuf.decodeInt64()
	}

	impliedVol := msgBuf.decodeFloat64()
	if impliedVol < 0 { // -1 is the "not computed" indicator
		impliedVol = UNSET_FLOAT
	}

	delta := msgBuf.decodeFloat64()
	if delta == -2 { // -2 is the "not computed" indicator
		delta = UNSET_FLOAT
	}

	if version >= 6 || tickType == MODEL_OPTION || tickType == DELAYED_MODEL_OPTION {
		optPrice = msgBuf.decodeFloat64()
		if optPrice == -1 { // -1 is the "not computed" indicator
			optPrice = UNSET_FLOAT
		}
		pvDividend = msgBuf.decodeFloat64()
		if pvDividend == -1 { // -1 is the "not computed" indicator
			pvDividend = UNSET_FLOAT
		}
	}

	if version >= 6 {
		gamma = msgBuf.decodeFloat64()
		if gamma == -2 { // -2 is the "not yet computed" indicator
			gamma = UNSET_FLOAT
		}
		vega = msgBuf.decodeFloat64()
		if vega == -2 { // -2 is the "not yet computed" indicator
			vega = UNSET_FLOAT
		}
		theta = msgBuf.decodeFloat64()
		if theta == -2 { // -2 is the "not yet computed" indicator
			theta = UNSET_FLOAT
		}
		undPrice = msgBuf.decodeFloat64()
		if undPrice == -1 { // -1 is the "not computed" indicator
			undPrice = UNSET_FLOAT
		}
	}

	d.wrapper.TickOptionComputation(reqID, tickType, tickAttrib, impliedVol, delta, optPrice, pvDividend, gamma, vega, theta, undPrice)

}

func (d *EDecoder) processTickGenericMsg(msgBuf *MsgBuffer) {

	_ = msgBuf.decodeString()

	reqID := msgBuf.decodeInt64()
	tickType := msgBuf.decodeInt64()
	value := msgBuf.decodeFloat64()

	d.wrapper.TickGeneric(reqID, tickType, value)
}

func (d *EDecoder) processTickStringMsg(msgBuf *MsgBuffer) {

	_ = msgBuf.decodeString()

	reqID := msgBuf.decodeInt64()
	tickType := msgBuf.decodeInt64()
	value := msgBuf.decodeString()

	d.wrapper.TickString(reqID, tickType, value)
}

func (d *EDecoder) processTickEfpMsg(msgBuf *MsgBuffer) {

	_ = msgBuf.decodeString()

	reqID := msgBuf.decodeInt64()
	tickType := msgBuf.decodeInt64()
	basisPoints := msgBuf.decodeFloat64()
	formattedBasisPoints := msgBuf.decodeString()
	totalDividends := msgBuf.decodeFloat64()
	holdDays := msgBuf.decodeInt64()
	futureLastTradeDate := msgBuf.decodeString()
	dividendImpact := msgBuf.decodeFloat64()
	dividendsToLastTradeDate := msgBuf.decodeFloat64()

	d.wrapper.TickEFP(reqID, tickType, basisPoints, formattedBasisPoints, totalDividends, holdDays, futureLastTradeDate, dividendImpact, dividendsToLastTradeDate)
}

func (d *EDecoder) processOrderStatusMsg(msgBuf *MsgBuffer) {

	if d.serverVersion < MIN_SERVER_VER_MARKET_CAP_PRICE {
		_ = msgBuf.decodeString()
	}

	orderID := msgBuf.decodeInt64()
	status := msgBuf.decodeString()
	filled := msgBuf.decodeDecimal()
	remaining := msgBuf.decodeDecimal()
	avgFilledPrice := msgBuf.decodeFloat64()

	permID := msgBuf.decodeInt64()          // ver 2 field
	parentID := msgBuf.decodeInt64()        // ver 3 field
	lastFillPrice := msgBuf.decodeFloat64() // ver 4 field
	clientID := msgBuf.decodeInt64()        // ver 5 field
	whyHeld := msgBuf.decodeString()        // ver 6 field

	mktCapPrice := 0.0
	if d.serverVersion >= MIN_SERVER_VER_MARKET_CAP_PRICE {
		mktCapPrice = msgBuf.decodeFloat64()
	}

	d.wrapper.OrderStatus(orderID, status, filled, remaining, avgFilledPrice, permID, parentID, lastFillPrice, clientID, whyHeld, mktCapPrice)
}

func (d *EDecoder) processErrMsg(msgBuf *MsgBuffer) {

	if d.serverVersion < MIN_SERVER_VER_ERROR_TIME {
		_ = msgBuf.decodeString()
	}

	reqID := msgBuf.decodeInt64()

	errorCode := msgBuf.decodeInt64()
	errorString := msgBuf.decodeString()

	advancedOrderRejectJson := ""
	if d.serverVersion >= MIN_SERVER_VER_ADVANCED_ORDER_REJECT {
		advancedOrderRejectJson = msgBuf.decodeString()
	}
	var errorTime int64
	if d.serverVersion >= MIN_SERVER_VER_ERROR_TIME {
		errorTime = msgBuf.decodeInt64()
	}

	d.wrapper.Error(reqID, errorTime, errorCode, errorString, advancedOrderRejectJson)
}

func (d *EDecoder) processOpenOrderMsg(msgBuf *MsgBuffer) {

	order := NewOrder()
	contract := NewContract()
	orderState := NewOrderState()

	version := d.serverVersion
	if d.serverVersion < MIN_SERVER_VER_ORDER_CONTAINER {
		version = Version(msgBuf.decodeInt64())
	}

	orderDecoder := &OrderDecoder{order, contract, orderState, version, d.serverVersion}

	// read orderID
	orderDecoder.decodeOrderId(msgBuf)

	// read contract fields
	orderDecoder.decodeContractFields(msgBuf)

	// read order fields
	orderDecoder.decodeAction(msgBuf)
	orderDecoder.decodeTotalQuantity(msgBuf)
	orderDecoder.decodeOrderType(msgBuf)
	orderDecoder.decodeLmtPrice(msgBuf)
	orderDecoder.decodeAuxPrice(msgBuf)
	orderDecoder.decodeTIF(msgBuf)
	orderDecoder.decodeOcaGroup(msgBuf)
	orderDecoder.decodeAccount(msgBuf)
	orderDecoder.decodeOpenClose(msgBuf)
	orderDecoder.decodeOrigin(msgBuf)
	orderDecoder.decodeOrderRef(msgBuf)
	orderDecoder.decodeClientId(msgBuf)
	orderDecoder.decodePermId(msgBuf)
	orderDecoder.decodeOutsideRth(msgBuf)
	orderDecoder.decodeHidden(msgBuf)
	orderDecoder.decodeDiscretionaryAmount(msgBuf)
	orderDecoder.decodeGoodAfterTime(msgBuf)
	orderDecoder.skipSharesAllocation(msgBuf)
	orderDecoder.decodeFAParams(msgBuf)
	orderDecoder.decodeModelCode(msgBuf)
	orderDecoder.decodeGoodTillDate(msgBuf)
	orderDecoder.decodeRule80A(msgBuf)
	orderDecoder.decodePercentOffset(msgBuf)
	orderDecoder.decodeSettlingFirm(msgBuf)
	orderDecoder.decodeShortSaleParams(msgBuf)
	orderDecoder.decodeAuctionStrategy(msgBuf)
	orderDecoder.decodeBoxOrderParams(msgBuf)
	orderDecoder.decodePegToStkOrVolOrderParams(msgBuf)
	orderDecoder.decodeDisplaySize(msgBuf)
	orderDecoder.decodeBlockOrder(msgBuf)
	orderDecoder.decodeSweepToFill(msgBuf)
	orderDecoder.decodeAllOrNone(msgBuf)
	orderDecoder.decodeMinQty(msgBuf)
	orderDecoder.decodeOcaType(msgBuf)
	orderDecoder.skipETradeOnly(msgBuf)
	orderDecoder.skipFirmQuoteOnly(msgBuf)
	orderDecoder.skipNbboPriceCap(msgBuf)
	orderDecoder.decodeParentId(msgBuf)
	orderDecoder.decodeTriggerMethod(msgBuf)
	orderDecoder.decodeVolOrderParams(msgBuf, true)
	orderDecoder.decodeTrailParams(msgBuf)
	orderDecoder.decodeBasisPoints(msgBuf)
	orderDecoder.decodeComboLegs(msgBuf)
	orderDecoder.decodeSmartComboRoutingParams(msgBuf)
	orderDecoder.decodeScaleOrderParams(msgBuf)
	orderDecoder.decodeHedgeParams(msgBuf)
	orderDecoder.decodeOptOutSmartRouting(msgBuf)
	orderDecoder.decodeClearingParams(msgBuf)
	orderDecoder.decodeNotHeld(msgBuf)
	orderDecoder.decodeDeltaNeutral(msgBuf)
	orderDecoder.decodeAlgoParams(msgBuf)
	orderDecoder.decodeSolicited(msgBuf)
	orderDecoder.decodeWhatIfInfoAndCommission(msgBuf)
	orderDecoder.decodeVolRandomizeFlags(msgBuf)
	orderDecoder.decodePegBenchParams(msgBuf)
	orderDecoder.decodeConditions(msgBuf)
	orderDecoder.decodeAdjustedOrderParams(msgBuf)
	orderDecoder.decodeSoftDollarTier(msgBuf)
	orderDecoder.decodeCashQty(msgBuf)
	orderDecoder.decodeDontUseAutoPriceForHedge(msgBuf)
	orderDecoder.decodeIsOmsContainer(msgBuf)
	orderDecoder.decodeDiscretionaryUpToLimitPrice(msgBuf)
	orderDecoder.decodeUsePriceMgmtAlgo(msgBuf)
	orderDecoder.decodeDuration(msgBuf)
	orderDecoder.decodePostToAts(msgBuf)
	orderDecoder.decodeAutoCancelParent(msgBuf, MIN_SERVER_VER_AUTO_CANCEL_PARENT)
	orderDecoder.decodePegBestPegMidOrderAttributes(msgBuf)
	orderDecoder.decodeCustomerAccount(msgBuf)
	orderDecoder.decodeProfessionalCustomer(msgBuf)
	orderDecoder.decodeBondAccruedInterest(msgBuf)
	orderDecoder.decodeIncludeOvernight(msgBuf)
	orderDecoder.decodeCMETaggingFields(msgBuf)

	d.wrapper.OpenOrder(order.OrderID, contract, order, orderState)
}

func (d *EDecoder) processAcctValueMsg(msgBuf *MsgBuffer) {

	_ = msgBuf.decodeString()

	tag := msgBuf.decodeString()
	val := msgBuf.decodeString()
	currency := msgBuf.decodeString()
	accountName := msgBuf.decodeString()

	d.wrapper.UpdateAccountValue(tag, val, currency, accountName)
}

func (d *EDecoder) processPortfolioValueMsg(msgBuf *MsgBuffer) {

	version := msgBuf.decodeInt64()

	// read contract fields
	contract := NewContract()
	contract.ConID = msgBuf.decodeInt64() // ver 6 field
	contract.Symbol = msgBuf.decodeString()
	contract.SecType = msgBuf.decodeString()
	contract.LastTradeDateOrContractMonth = msgBuf.decodeString()
	contract.Strike = msgBuf.decodeFloat64()
	contract.Right = msgBuf.decodeString()

	if version >= 7 {
		contract.Multiplier = msgBuf.decodeString()
		contract.PrimaryExchange = msgBuf.decodeString()
	}

	contract.Currency = msgBuf.decodeString()
	contract.LocalSymbol = msgBuf.decodeString() // ver 2 field
	if version >= 8 {
		contract.TradingClass = msgBuf.decodeString()
	}
	position := msgBuf.decodeDecimal()

	marketPrice := msgBuf.decodeFloat64()
	marketValue := msgBuf.decodeFloat64()
	averageCost := msgBuf.decodeFloat64()   // ver 3 field
	unrealizedPNL := msgBuf.decodeFloat64() // ver 3 field
	realizedPNL := msgBuf.decodeFloat64()   // ver 3 field

	accountName := msgBuf.decodeString() // ver 4 field

	if version == 6 && d.serverVersion == 39 {
		contract.PrimaryExchange = msgBuf.decodeString()
	}

	d.wrapper.UpdatePortfolio(contract, position, marketPrice, marketValue, averageCost, unrealizedPNL, realizedPNL, accountName)

}

func (d *EDecoder) processAcctUpdateTimeMsg(msgBuf *MsgBuffer) {

	_ = msgBuf.decodeString()

	timeStamp := msgBuf.decodeString()

	d.wrapper.UpdateAccountTime(timeStamp)
}

func (d *EDecoder) processNextValidIdMsg(msgBuf *MsgBuffer) {

	_ = msgBuf.decodeString()

	reqID := msgBuf.decodeInt64()

	d.wrapper.NextValidID(reqID)
}

func (d *EDecoder) processContractDataMsg(msgBuf *MsgBuffer) {

	var version int64 = 8
	if d.serverVersion < MIN_SERVER_VER_SIZE_RULES {
		version = msgBuf.decodeInt64()
	}

	var reqID int64 = -1
	if version >= 3 {
		reqID = msgBuf.decodeInt64()
	}

	cd := NewContractDetails()
	cd.Contract = *NewContract()
	cd.Contract.Symbol = msgBuf.decodeString()
	cd.Contract.SecType = msgBuf.decodeString()
	d.readLastTradeDate(msgBuf, cd, false)
	if d.serverVersion >= MIN_SERVER_VER_LAST_TRADE_DATE {
		cd.Contract.LastTradeDate = msgBuf.decodeString()
	}
	cd.Contract.Strike = msgBuf.decodeFloat64()
	cd.Contract.Right = msgBuf.decodeString()
	cd.Contract.Exchange = msgBuf.decodeString()
	cd.Contract.Currency = msgBuf.decodeString()
	cd.Contract.LocalSymbol = msgBuf.decodeString()
	cd.MarketName = msgBuf.decodeString()
	cd.Contract.TradingClass = msgBuf.decodeString()
	cd.Contract.ConID = msgBuf.decodeInt64()
	cd.MinTick = msgBuf.decodeFloat64()
	if d.serverVersion >= MIN_SERVER_VER_MD_SIZE_MULTIPLIER && d.serverVersion < MIN_SERVER_VER_SIZE_RULES {
		_ = msgBuf.decodeInt64() // MdSizeMultiplier - not used anymore
	}
	cd.Contract.Multiplier = msgBuf.decodeString()
	cd.OrderTypes = msgBuf.decodeString()
	cd.ValidExchanges = msgBuf.decodeString()
	cd.PriceMagnifier = msgBuf.decodeInt64()
	if version >= 4 {
		cd.UnderConID = msgBuf.decodeInt64()
	}
	if version >= 5 {
		if d.serverVersion >= MIN_SERVER_VER_ENCODE_MSG_ASCII7 {
			cd.LongName = msgBuf.decodeStringUnescaped()
		} else {
			cd.LongName = msgBuf.decodeString()
		}
		cd.Contract.PrimaryExchange = msgBuf.decodeString()
	}
	if version >= 6 {
		cd.ContractMonth = msgBuf.decodeString()
		cd.Industry = msgBuf.decodeString()
		cd.Category = msgBuf.decodeString()
		cd.Subcategory = msgBuf.decodeString()
		cd.TimeZoneID = msgBuf.decodeString()
		cd.TradingHours = msgBuf.decodeString()
		cd.LiquidHours = msgBuf.decodeString()
	}
	if version >= 8 {
		cd.EVRule = msgBuf.decodeString()
		cd.EVMultiplier = msgBuf.decodeInt64()
	}
	if version >= 7 {
		secIDListCount := msgBuf.decodeInt64()
		cd.SecIDList = make([]TagValue, 0, secIDListCount)
		var i int64
		for i = 0; i < secIDListCount; i++ {
			tagValue := NewTagValue()
			tagValue.Tag = msgBuf.decodeString()
			tagValue.Value = msgBuf.decodeString()
			cd.SecIDList = append(cd.SecIDList, tagValue)
		}
	}

	if d.serverVersion >= MIN_SERVER_VER_AGG_GROUP {
		cd.AggGroup = msgBuf.decodeInt64()
	}

	if d.serverVersion >= MIN_SERVER_VER_UNDERLYING_INFO {
		cd.UnderSymbol = msgBuf.decodeString()
		cd.UnderSecType = msgBuf.decodeString()
	}

	if d.serverVersion >= MIN_SERVER_VER_MARKET_RULES {
		cd.MarketRuleIDs = msgBuf.decodeString()
	}

	if d.serverVersion >= MIN_SERVER_VER_REAL_EXPIRATION_DATE {
		cd.RealExpirationDate = msgBuf.decodeString()
	}

	if d.serverVersion >= MIN_SERVER_VER_STOCK_TYPE {
		cd.StockType = msgBuf.decodeString()
	}

	if d.serverVersion >= MIN_SERVER_VER_FRACTIONAL_SIZE_SUPPORT && d.serverVersion < MIN_SERVER_VER_SIZE_RULES {
		_ = msgBuf.decodeDecimal() // sizeMinTick - not used anymore
	}

	if d.serverVersion >= MIN_SERVER_VER_SIZE_RULES {
		cd.MinSize = msgBuf.decodeDecimal()
		cd.SizeIncrement = msgBuf.decodeDecimal()
		cd.SuggestedSizeIncrement = msgBuf.decodeDecimal()
	}

	if d.serverVersion >= MIN_SERVER_VER_FUND_DATA_FIELDS && cd.Contract.SecType == "FUND" {
		cd.FundName = msgBuf.decodeString()
		cd.FundFamily = msgBuf.decodeString()
		cd.FundType = msgBuf.decodeString()
		cd.FundFrontLoad = msgBuf.decodeString()
		cd.FundBackLoad = msgBuf.decodeString()
		cd.FundBackLoadTimeInterval = msgBuf.decodeString()
		cd.FundManagementFee = msgBuf.decodeString()
		cd.FundClosed = msgBuf.decodeBool()
		cd.FundClosedForNewInvestors = msgBuf.decodeBool()
		cd.FundClosedForNewMoney = msgBuf.decodeBool()
		cd.FundNotifyAmount = msgBuf.decodeString()
		cd.FundMinimumInitialPurchase = msgBuf.decodeString()
		cd.FundSubsequentMinimumPurchase = msgBuf.decodeString()
		cd.FundBlueSkyStates = msgBuf.decodeString()
		cd.FundBlueSkyTerritories = msgBuf.decodeString()
		cd.FundDistributionPolicyIndicator = getFundDistributionPolicyIndicator(msgBuf.decodeString())
		cd.FundAssetType = getFundAssetType(msgBuf.decodeString())
	}

	if d.serverVersion >= MIN_SERVER_VER_INELIGIBILITY_REASONS {
		ineligibilityReasonListCount := msgBuf.decodeInt64()
		if ineligibilityReasonListCount > 0 {
			cd.IneligibilityReasonList = make([]IneligibilityReason, ineligibilityReasonListCount)
			var i int64
			for i = 0; i < ineligibilityReasonListCount; i++ {
				ineligibilityReason := IneligibilityReason{}
				ineligibilityReason.ID = msgBuf.decodeString()
				ineligibilityReason.Description = msgBuf.decodeString()
				cd.IneligibilityReasonList = append(cd.IneligibilityReasonList, ineligibilityReason)
			}
		}
	}

	d.wrapper.ContractDetails(reqID, cd)
}

func (d *EDecoder) processBondContractDataMsg(msgBuf *MsgBuffer) {

	var version Version = 6
	if d.serverVersion < MIN_SERVER_VER_SIZE_RULES {
		version = Version(msgBuf.decodeInt64())
	}

	var reqID int64 = -1
	if version >= 3 {
		reqID = msgBuf.decodeInt64()
	}

	contract := NewContractDetails()
	contract.Contract.Symbol = msgBuf.decodeString()
	contract.Contract.SecType = msgBuf.decodeString()
	contract.Cusip = msgBuf.decodeString()
	contract.Coupon = msgBuf.decodeFloat64()
	d.readLastTradeDate(msgBuf, contract, true)
	contract.IssueDate = msgBuf.decodeString()
	contract.Ratings = msgBuf.decodeString()
	contract.BondType = msgBuf.decodeString()
	contract.CouponType = msgBuf.decodeString()
	contract.Convertible = msgBuf.decodeBool()
	contract.Callable = msgBuf.decodeBool()
	contract.Putable = msgBuf.decodeBool()
	contract.DescAppend = msgBuf.decodeString()
	contract.Contract.Exchange = msgBuf.decodeString()
	contract.Contract.Currency = msgBuf.decodeString()
	contract.MarketName = msgBuf.decodeString()
	contract.Contract.TradingClass = msgBuf.decodeString()
	contract.Contract.ConID = msgBuf.decodeInt64()
	contract.MinTick = msgBuf.decodeFloat64()

	if d.serverVersion >= MIN_SERVER_VER_MD_SIZE_MULTIPLIER && d.serverVersion < MIN_SERVER_VER_SIZE_RULES {
		_ = msgBuf.decodeInt64() // mdSizeMultiplier - not used anymore
	}

	contract.OrderTypes = msgBuf.decodeString()
	contract.ValidExchanges = msgBuf.decodeString()
	contract.NextOptionDate = msgBuf.decodeString()
	contract.NextOptionType = msgBuf.decodeString()
	contract.NextOptionPartial = msgBuf.decodeBool()
	contract.Notes = msgBuf.decodeString()

	if version >= 4 {
		contract.LongName = msgBuf.decodeString()
	}

	if d.serverVersion >= MIN_SERVER_VER_BOND_TRADING_HOURS {
		contract.TimeZoneID = msgBuf.decodeString()
		contract.TradingHours = msgBuf.decodeString()
		contract.LiquidHours = msgBuf.decodeString()
	}

	if version >= 6 {
		contract.EVRule = msgBuf.decodeString()
		contract.EVMultiplier = msgBuf.decodeInt64()
	}

	if version >= 5 {
		secIDListCount := msgBuf.decodeInt64()
		contract.SecIDList = make([]TagValue, 0, secIDListCount)
		var i int64
		for i = 0; i < secIDListCount; i++ {
			tagValue := NewTagValue()
			tagValue.Tag = msgBuf.decodeString()
			tagValue.Value = msgBuf.decodeString()
			contract.SecIDList = append(contract.SecIDList, tagValue)
		}
	}

	if d.serverVersion >= MIN_SERVER_VER_AGG_GROUP {
		contract.AggGroup = msgBuf.decodeInt64()
	}

	if d.serverVersion >= MIN_SERVER_VER_MARKET_RULES {
		contract.MarketRuleIDs = msgBuf.decodeString()
	}

	if d.serverVersion >= MIN_SERVER_VER_SIZE_RULES {
		contract.MinSize = msgBuf.decodeDecimal()
		contract.SizeIncrement = msgBuf.decodeDecimal()
		contract.SuggestedSizeIncrement = msgBuf.decodeDecimal()
	}

	d.wrapper.BondContractDetails(reqID, contract)
}

func (d *EDecoder) processExecutionDetailsMsg(msgBuf *MsgBuffer) {

	version := d.serverVersion
	if d.serverVersion < MIN_SERVER_VER_LAST_LIQUIDITY {
		version = Version(msgBuf.decodeInt64())
	}

	var reqID int64 = -1
	if version >= 7 {
		reqID = msgBuf.decodeInt64()
	}

	orderID := msgBuf.decodeInt64()

	// decode contact fields
	contract := NewContract()
	contract.ConID = msgBuf.decodeInt64()
	contract.Symbol = msgBuf.decodeString()
	contract.SecType = msgBuf.decodeString()
	contract.LastTradeDateOrContractMonth = msgBuf.decodeString()
	contract.Strike = msgBuf.decodeFloat64()
	contract.Right = msgBuf.decodeString()

	if version >= 9 {
		contract.Multiplier = msgBuf.decodeString()
	}

	contract.Exchange = msgBuf.decodeString()
	contract.Currency = msgBuf.decodeString()
	contract.LocalSymbol = msgBuf.decodeString()

	if version >= 10 {
		contract.TradingClass = msgBuf.decodeString()
	}

	// read execution fields
	execution := NewExecution()
	execution.OrderID = orderID
	execution.ExecID = msgBuf.decodeString()
	execution.Time = msgBuf.decodeString()
	execution.AcctNumber = msgBuf.decodeString()
	execution.Exchange = msgBuf.decodeString()
	execution.Side = msgBuf.decodeString()
	execution.Shares = msgBuf.decodeDecimal()
	execution.Price = msgBuf.decodeFloat64()
	execution.PermID = msgBuf.decodeInt64()
	execution.ClientID = msgBuf.decodeInt64()
	execution.Liquidation = msgBuf.decodeInt64()

	if version >= 6 {
		execution.CumQty = msgBuf.decodeDecimal()
		execution.AvgPrice = msgBuf.decodeFloat64()
	}

	if version >= 8 {
		execution.OrderRef = msgBuf.decodeString()
	}

	if version >= 9 {
		execution.EVRule = msgBuf.decodeString()
		execution.EVMultiplier = msgBuf.decodeFloat64()
	}

	if d.serverVersion >= MIN_SERVER_VER_MODELS_SUPPORT {
		execution.ModelCode = msgBuf.decodeString()
	}

	if d.serverVersion >= MIN_SERVER_VER_LAST_LIQUIDITY {
		execution.LastLiquidity = msgBuf.decodeInt64()
	}
	if d.serverVersion >= MIN_SERVER_VER_PENDING_PRICE_REVISION {
		execution.PendingPriceRevision = msgBuf.decodeBool()
	}

	d.wrapper.ExecDetails(reqID, contract, execution)
}

func (d *EDecoder) processMarketDepthMsg(msgBuf *MsgBuffer) {

	_ = msgBuf.decodeString()

	tickerID := msgBuf.decodeInt64()

	position := msgBuf.decodeInt64()
	operation := msgBuf.decodeInt64()
	side := msgBuf.decodeInt64()
	price := msgBuf.decodeFloat64()
	size := msgBuf.decodeDecimal()

	d.wrapper.UpdateMktDepth(tickerID, position, operation, side, price, size)
}

func (d *EDecoder) processMarketDepthL2Msg(msgBuf *MsgBuffer) {

	_ = msgBuf.decodeString()

	tickerID := msgBuf.decodeInt64()

	position := msgBuf.decodeInt64()
	marketMaker := msgBuf.decodeString()
	operation := msgBuf.decodeInt64()
	side := msgBuf.decodeInt64()
	price := msgBuf.decodeFloat64()
	size := msgBuf.decodeDecimal()

	isSmartDepth := false
	if d.serverVersion >= MIN_SERVER_VER_SMART_DEPTH {
		isSmartDepth = msgBuf.decodeBool()
	}

	d.wrapper.UpdateMktDepthL2(tickerID, position, marketMaker, operation, side, price, size, isSmartDepth)
}

func (d *EDecoder) processNewsBulletinsMsg(msgBuf *MsgBuffer) {

	_ = msgBuf.decodeString()

	msgID := msgBuf.decodeInt64()
	msgType := msgBuf.decodeInt64()
	newsMessage := msgBuf.decodeString()
	originExch := msgBuf.decodeString()

	d.wrapper.UpdateNewsBulletin(msgID, msgType, newsMessage, originExch)
}

func (d *EDecoder) processManagedAcctsMsg(msgBuf *MsgBuffer) {

	_ = msgBuf.decodeString()

	accountsNames := msgBuf.decodeString()
	accountsList := strings.Split(accountsNames, ",")

	d.wrapper.ManagedAccounts(accountsList)
}

func (d *EDecoder) processReceiveFaMsg(msgBuf *MsgBuffer) {

	_ = msgBuf.decodeString()

	faDataType := FaDataType(msgBuf.decodeInt64())
	cxml := msgBuf.decodeString()

	d.wrapper.ReceiveFA(faDataType, cxml)
}

func (d *EDecoder) processHistoricalDataMsg(msgBuf *MsgBuffer) {

	if d.serverVersion < MIN_SERVER_VER_SYNT_REALTIME_BARS {
		_ = msgBuf.decodeString()
	}

	reqID := msgBuf.decodeInt64()
	startDateStr := msgBuf.decodeString()
	endDateStr := msgBuf.decodeString()

	itemCount := msgBuf.decodeInt64()

	var i int64
	for i = 0; i < itemCount; i++ {
		bar := NewBar()
		bar.Date = msgBuf.decodeString()
		bar.Open = msgBuf.decodeFloat64()
		bar.High = msgBuf.decodeFloat64()
		bar.Low = msgBuf.decodeFloat64()
		bar.Close = msgBuf.decodeFloat64()
		bar.Volume = msgBuf.decodeDecimal()
		bar.Wap = msgBuf.decodeDecimal()

		if d.serverVersion < MIN_SERVER_VER_SYNT_REALTIME_BARS {
			_ = msgBuf.decodeString()
		}

		bar.BarCount = msgBuf.decodeInt64()

		d.wrapper.HistoricalData(reqID, &bar)
	}

	d.wrapper.HistoricalDataEnd(reqID, startDateStr, endDateStr)
}

func (d *EDecoder) processScannerDataMsg(msgBuf *MsgBuffer) {

	_ = msgBuf.decodeString()

	reqID := msgBuf.decodeInt64()

	numberOfElements := msgBuf.decodeInt64()

	var i int64
	for i = 0; i < numberOfElements; i++ {

		contractDetails := NewContractDetails()

		rank := msgBuf.decodeInt64()
		contractDetails.Contract.ConID = msgBuf.decodeInt64()
		contractDetails.Contract.Symbol = msgBuf.decodeString()
		contractDetails.Contract.SecType = msgBuf.decodeString()
		contractDetails.Contract.LastTradeDateOrContractMonth = msgBuf.decodeString()
		contractDetails.Contract.Strike = msgBuf.decodeFloat64()
		contractDetails.Contract.Right = msgBuf.decodeString()
		contractDetails.Contract.Exchange = msgBuf.decodeString()
		contractDetails.Contract.Currency = msgBuf.decodeString()
		contractDetails.Contract.LocalSymbol = msgBuf.decodeString()
		contractDetails.MarketName = msgBuf.decodeString()
		contractDetails.Contract.TradingClass = msgBuf.decodeString()
		distance := msgBuf.decodeString()
		benchmark := msgBuf.decodeString()
		projection := msgBuf.decodeString()
		legsStr := msgBuf.decodeString()

		d.wrapper.ScannerData(reqID, rank, contractDetails, distance, benchmark, projection, legsStr)

	}

	d.wrapper.ScannerDataEnd(reqID)
}

func (d *EDecoder) processScannerParametersMsg(msgBuf *MsgBuffer) {

	_ = msgBuf.decodeString()

	xml := msgBuf.decodeString()

	d.wrapper.ScannerParameters(xml)
}

func (d *EDecoder) processCurrentTimeMsg(msgBuf *MsgBuffer) {

	_ = msgBuf.decodeString()

	t := msgBuf.decodeInt64()

	d.wrapper.CurrentTime(t)
}

func (d *EDecoder) processRealTimeBarsMsg(msgBuf *MsgBuffer) {

	_ = msgBuf.decodeString()

	reqID := msgBuf.decodeInt64()

	time := msgBuf.decodeInt64()
	open := msgBuf.decodeFloat64()
	high := msgBuf.decodeFloat64()
	low := msgBuf.decodeFloat64()
	close := msgBuf.decodeFloat64()
	volume := msgBuf.decodeDecimal()
	wap := msgBuf.decodeDecimal()
	count := msgBuf.decodeInt64()

	d.wrapper.RealtimeBar(reqID, time, open, high, low, close, volume, wap, count)

}

func (d *EDecoder) processFundamentalDataMsg(msgBuf *MsgBuffer) {

	_ = msgBuf.decodeString()

	reqID := msgBuf.decodeInt64()
	data := msgBuf.decodeString()

	d.wrapper.FundamentalData(reqID, data)
}

func (d *EDecoder) processContractDataEndMsg(msgBuf *MsgBuffer) {

	_ = msgBuf.decodeString()

	reqID := msgBuf.decodeInt64()

	d.wrapper.ContractDetailsEnd(reqID)
}

func (d *EDecoder) processOpenOrderEndMsg(*MsgBuffer) {

	d.wrapper.OpenOrderEnd()
}

func (d *EDecoder) processAcctDownloadEndMsg(msgBuf *MsgBuffer) {

	_ = msgBuf.decodeString()

	accountName := msgBuf.decodeString()

	d.wrapper.AccountDownloadEnd(accountName)
}

func (d *EDecoder) processExecutionDetailsEndMsg(msgBuf *MsgBuffer) {

	_ = msgBuf.decodeString()

	reqID := msgBuf.decodeInt64()

	d.wrapper.ExecDetailsEnd(reqID)
}

func (d *EDecoder) processDeltaNeutralValidationMsg(msgBuf *MsgBuffer) {

	_ = msgBuf.decodeString()

	reqID := msgBuf.decodeInt64()

	deltaNeutralContract := NewDeltaNeutralContract()

	deltaNeutralContract.ConID = msgBuf.decodeInt64()
	deltaNeutralContract.Delta = msgBuf.decodeFloat64()
	deltaNeutralContract.Price = msgBuf.decodeFloat64()

	d.wrapper.DeltaNeutralValidation(reqID, deltaNeutralContract)
}

func (d *EDecoder) processTickSnapshotEndMsg(msgBuf *MsgBuffer) {

	_ = msgBuf.decodeString()

	reqID := msgBuf.decodeInt64()

	d.wrapper.TickSnapshotEnd(reqID)
}

func (d *EDecoder) processMarketDataTypeMsg(msgBuf *MsgBuffer) {

	_ = msgBuf.decodeString()

	reqID := msgBuf.decodeInt64()
	marketDataType := msgBuf.decodeInt64()

	d.wrapper.MarketDataType(reqID, marketDataType)
}

func (d *EDecoder) processCommissionReportMsg(msgBuf *MsgBuffer) {

	_ = msgBuf.decodeString()

	commissionReport := NewCommissionReport()
	commissionReport.ExecID = msgBuf.decodeString()
	commissionReport.Commission = msgBuf.decodeFloat64()
	commissionReport.Currency = msgBuf.decodeString()
	commissionReport.RealizedPNL = msgBuf.decodeFloat64()
	commissionReport.Yield = msgBuf.decodeFloat64()
	commissionReport.YieldRedemptionDate = msgBuf.decodeInt64()

	d.wrapper.CommissionReport(commissionReport)
}

func (d *EDecoder) processPositionDataMsg(msgBuf *MsgBuffer) {

	version := msgBuf.decodeInt64()

	account := msgBuf.decodeString()

	// decode contract fields
	contract := NewContract()
	contract.ConID = msgBuf.decodeInt64()
	contract.Symbol = msgBuf.decodeString()
	contract.SecType = msgBuf.decodeString()
	contract.LastTradeDateOrContractMonth = msgBuf.decodeString()
	contract.Strike = msgBuf.decodeFloat64()
	contract.Right = msgBuf.decodeString()
	contract.Multiplier = msgBuf.decodeString()
	contract.Exchange = msgBuf.decodeString()
	contract.Currency = msgBuf.decodeString()
	contract.LocalSymbol = msgBuf.decodeString()
	if version >= 2 {
		contract.TradingClass = msgBuf.decodeString()
	}

	position := msgBuf.decodeDecimal()

	var avgCost float64
	if version >= 3 {
		avgCost = msgBuf.decodeFloat64()
	}

	d.wrapper.Position(account, contract, position, avgCost)
}

func (d *EDecoder) processPositionEndMsg(*MsgBuffer) {

	d.wrapper.PositionEnd()
}

func (d *EDecoder) processAccountSummaryMsg(msgBuf *MsgBuffer) {

	_ = msgBuf.decodeString()

	reqID := msgBuf.decodeInt64()
	account := msgBuf.decodeString()
	tag := msgBuf.decodeString()
	value := msgBuf.decodeString()
	currency := msgBuf.decodeString()

	d.wrapper.AccountSummary(reqID, account, tag, value, currency)
}

func (d *EDecoder) processAccountSummaryEndMsg(msgBuf *MsgBuffer) {

	_ = msgBuf.decodeString()

	reqID := msgBuf.decodeInt64()

	d.wrapper.AccountSummaryEnd(reqID)
}

func (d *EDecoder) processVerifyMessageApiMsg(msgBuf *MsgBuffer) {

	_ = msgBuf.decodeString()

	apiData := msgBuf.decodeString()

	d.wrapper.VerifyMessageAPI(apiData)
}

func (d *EDecoder) processVerifyCompletedMsg(msgBuf *MsgBuffer) {

	_ = msgBuf.decodeString()

	isSuccessful := msgBuf.decodeBool()
	errorText := msgBuf.decodeString()

	d.wrapper.VerifyCompleted(isSuccessful, errorText)
}

func (d *EDecoder) processDisplayGroupListMsg(msgBuf *MsgBuffer) {

	_ = msgBuf.decodeString()

	reqID := msgBuf.decodeInt64()
	groups := msgBuf.decodeString()

	d.wrapper.DisplayGroupList(reqID, groups)
}

func (d *EDecoder) processDisplayGroupUpdatedMsg(msgBuf *MsgBuffer) {

	_ = msgBuf.decodeString()

	reqID := msgBuf.decodeInt64()
	contractInfo := msgBuf.decodeString()

	d.wrapper.DisplayGroupUpdated(reqID, contractInfo)
}

func (d *EDecoder) processVerifyAndAuthMessageApiMsg(msgBuf *MsgBuffer) {

	_ = msgBuf.decodeString()

	apiData := msgBuf.decodeString()
	xyzChallange := msgBuf.decodeString()

	d.wrapper.VerifyAndAuthMessageAPI(apiData, xyzChallange)
}

func (d *EDecoder) processVerifyAndAuthCompletedMsg(msgBuf *MsgBuffer) {

	_ = msgBuf.decodeString()

	isSuccessful := msgBuf.decodeBool()
	errorText := msgBuf.decodeString()

	d.wrapper.VerifyAndAuthCompleted(isSuccessful, errorText)
}

func (d *EDecoder) processPositionMultiMsg(msgBuf *MsgBuffer) {

	_ = msgBuf.decodeString()

	reqID := msgBuf.decodeInt64()
	account := msgBuf.decodeString()

	// decode contract fields
	contract := &Contract{}
	contract.ConID = msgBuf.decodeInt64()
	contract.Symbol = msgBuf.decodeString()
	contract.SecType = msgBuf.decodeString()
	contract.LastTradeDateOrContractMonth = msgBuf.decodeString()
	contract.Strike = msgBuf.decodeFloat64()
	contract.Right = msgBuf.decodeString()
	contract.Multiplier = msgBuf.decodeString()
	contract.Exchange = msgBuf.decodeString()
	contract.Currency = msgBuf.decodeString()
	contract.LocalSymbol = msgBuf.decodeString()
	contract.TradingClass = msgBuf.decodeString()

	pos := msgBuf.decodeDecimal()

	avgCost := msgBuf.decodeFloat64()
	modelCode := msgBuf.decodeString()

	d.wrapper.PositionMulti(reqID, account, modelCode, contract, pos, avgCost)
}

func (d *EDecoder) processPositionMultiEndMsg(msgBuf *MsgBuffer) {

	_ = msgBuf.decodeString()

	reqID := msgBuf.decodeInt64()

	d.wrapper.PositionMultiEnd(reqID)
}

func (d *EDecoder) processAccountUpdateMultiMsg(msgBuf *MsgBuffer) {

	_ = msgBuf.decodeString()

	reqID := msgBuf.decodeInt64()
	account := msgBuf.decodeString()
	modelCode := msgBuf.decodeString()
	key := msgBuf.decodeString()
	value := msgBuf.decodeString()
	currency := msgBuf.decodeString()

	d.wrapper.AccountUpdateMulti(reqID, account, modelCode, key, value, currency)
}

func (d *EDecoder) processAccountUpdateMultiEndMsg(msgBuf *MsgBuffer) {

	_ = msgBuf.decodeString()

	reqID := msgBuf.decodeInt64()

	d.wrapper.AccountUpdateMultiEnd(reqID)
}

func (d *EDecoder) processSecurityDefinitionOptionalParameterMsg(msgBuf *MsgBuffer) {

	reqID := msgBuf.decodeInt64()

	exchange := msgBuf.decodeString()
	underlyingConID := msgBuf.decodeInt64()
	tradingClass := msgBuf.decodeString()
	multiplier := msgBuf.decodeString()

	expCount := msgBuf.decodeInt64()
	expirations := make([]string, 0, expCount)
	var i int64
	for i = 0; i < expCount; i++ {
		expiration := msgBuf.decodeString()
		expirations = append(expirations, expiration)
	}

	strikeCount := msgBuf.decodeInt64()
	strikes := make([]float64, 0, strikeCount)
	for i = 0; i < strikeCount; i++ {
		strike := msgBuf.decodeFloat64()
		strikes = append(strikes, strike)
	}

	d.wrapper.SecurityDefinitionOptionParameter(reqID, exchange, underlyingConID, tradingClass, multiplier, expirations, strikes)

}

func (d *EDecoder) processSecurityDefinitionOptionalParameterEndMsg(msgBuf *MsgBuffer) {

	reqID := msgBuf.decodeInt64()

	d.wrapper.SecurityDefinitionOptionParameterEnd(reqID)
}

func (d *EDecoder) processSoftDollarTiersMsg(msgBuf *MsgBuffer) {

	reqID := msgBuf.decodeInt64()

	tiersCount := msgBuf.decodeInt64()
	tiers := make([]SoftDollarTier, 0, tiersCount)
	var i int64
	for i = 0; i < tiersCount; i++ {
		tier := NewSoftDollarTier()
		tier.Name = msgBuf.decodeString()
		tier.Value = msgBuf.decodeString()
		tier.DisplayName = msgBuf.decodeString()
		tiers = append(tiers, tier)
	}

	d.wrapper.SoftDollarTiers(reqID, tiers)
}

func (d *EDecoder) processFamilyCodesMsg(msgBuf *MsgBuffer) {

	familyCodesCount := msgBuf.decodeInt64()
	familyCodes := make([]FamilyCode, 0, familyCodesCount)
	var i int64
	for i = 0; i < familyCodesCount; i++ {
		familyCode := NewFamilyCode()
		familyCode.AccountID = msgBuf.decodeString()
		familyCode.FamilyCodeStr = msgBuf.decodeString()
		familyCodes = append(familyCodes, familyCode)
	}

	d.wrapper.FamilyCodes(familyCodes)
}

func (d *EDecoder) processSymbolSamplesMsg(msgBuf *MsgBuffer) {

	reqID := msgBuf.decodeInt64()

	contractDescriptionsCount := msgBuf.decodeInt64()
	contractDescriptions := make([]ContractDescription, 0, contractDescriptionsCount)
	var i int64
	for i = 0; i < contractDescriptionsCount; i++ {
		conDesc := NewContractDescription()
		conDesc.Contract.ConID = msgBuf.decodeInt64()
		conDesc.Contract.Symbol = msgBuf.decodeString()
		conDesc.Contract.SecType = msgBuf.decodeString()
		conDesc.Contract.PrimaryExchange = msgBuf.decodeString()
		conDesc.Contract.Currency = msgBuf.decodeString()

		DerivativeSecTypesCount := msgBuf.decodeInt64()
		conDesc.DerivativeSecTypes = make([]string, 0, DerivativeSecTypesCount)
		var j int64
		for j = 0; j < DerivativeSecTypesCount; j++ {
			derivativeSecType := msgBuf.decodeString()
			conDesc.DerivativeSecTypes = append(conDesc.DerivativeSecTypes, derivativeSecType)
		}
		if d.serverVersion >= MIN_SERVER_VER_BOND_ISSUERID {
			conDesc.Contract.Description = msgBuf.decodeString()
			conDesc.Contract.IssuerID = msgBuf.decodeString()
		}
		contractDescriptions = append(contractDescriptions, conDesc)
	}
	d.wrapper.SymbolSamples(reqID, contractDescriptions)
}

func (d *EDecoder) processMktDepthExchangesMsg(msgBuf *MsgBuffer) {

	depthMktDataDescriptionsCount := msgBuf.decodeInt64()
	depthMktDataDescriptions := make([]DepthMktDataDescription, 0, depthMktDataDescriptionsCount)

	var i int64
	for i = 0; i < depthMktDataDescriptionsCount; i++ {
		desc := NewDepthMktDataDescription()
		desc.Exchange = msgBuf.decodeString()
		desc.SecType = msgBuf.decodeString()
		if d.serverVersion >= MIN_SERVER_VER_SERVICE_DATA_TYPE {
			desc.ListingExch = msgBuf.decodeString()
			desc.SecType = msgBuf.decodeString()
			desc.AggGroup = msgBuf.decodeInt64()
		} else {
			_ = msgBuf.decodeInt64() // boolean notSuppIsL2
		}

		depthMktDataDescriptions = append(depthMktDataDescriptions, desc)
	}

	d.wrapper.MktDepthExchanges(depthMktDataDescriptions)
}

func (d *EDecoder) processTickNewsMsg(msgBuf *MsgBuffer) {

	tickerID := msgBuf.decodeInt64()

	timeStamp := msgBuf.decodeInt64()
	providerCode := msgBuf.decodeString()
	articleID := msgBuf.decodeString()
	headline := msgBuf.decodeString()
	extraData := msgBuf.decodeString()

	d.wrapper.TickNews(tickerID, timeStamp, providerCode, articleID, headline, extraData)
}

func (d *EDecoder) processTickReqParamsMsg(msgBuf *MsgBuffer) {

	tickerID := msgBuf.decodeInt64()

	minTick := msgBuf.decodeFloat64()
	bboExchange := msgBuf.decodeString()
	snapshotPermissions := msgBuf.decodeInt64()

	d.wrapper.TickReqParams(tickerID, minTick, bboExchange, snapshotPermissions)
}

func (d *EDecoder) processSmartComponentsMsg(msgBuf *MsgBuffer) {

	reqID := msgBuf.decodeInt64()

	smartComponentsCount := msgBuf.decodeInt64()
	smartComponents := make([]SmartComponent, 0, smartComponentsCount)
	var i int64
	for i = 0; i < smartComponentsCount; i++ {
		smartComponent := NewSmartComponent()
		smartComponent.BitNumber = msgBuf.decodeInt64()
		smartComponent.Exchange = msgBuf.decodeString()
		smartComponent.ExchangeLetter = msgBuf.decodeString()
		smartComponents = append(smartComponents, smartComponent)
	}

	d.wrapper.SmartComponents(reqID, smartComponents)
}

func (d *EDecoder) processNewsProvidersMsg(msgBuf *MsgBuffer) {

	newsProvidersCount := msgBuf.decodeInt64()
	newsProviders := make([]NewsProvider, 0, newsProvidersCount)
	var i int64
	for i = 0; i < newsProvidersCount; i++ {
		provider := NewNewsProvider()
		provider.Name = msgBuf.decodeString()
		provider.Code = msgBuf.decodeString()
		newsProviders = append(newsProviders, provider)
	}

	d.wrapper.NewsProviders(newsProviders)
}

func (d *EDecoder) processNewsArticleMsg(msgBuf *MsgBuffer) {

	reqID := msgBuf.decodeInt64()

	articleType := msgBuf.decodeInt64()
	articleText := msgBuf.decodeString()

	d.wrapper.NewsArticle(reqID, articleType, articleText)
}

func (d *EDecoder) processHistoricalNewsMsg(msgBuf *MsgBuffer) {

	requestID := msgBuf.decodeInt64()

	time := msgBuf.decodeString()
	providerCode := msgBuf.decodeString()
	articleID := msgBuf.decodeString()
	headline := msgBuf.decodeString()

	d.wrapper.HistoricalNews(requestID, time, providerCode, articleID, headline)
}

func (d *EDecoder) processHistoricalNewsEndMsg(msgBuf *MsgBuffer) {

	requestID := msgBuf.decodeInt64()

	hasMore := msgBuf.decodeBool()

	d.wrapper.HistoricalNewsEnd(requestID, hasMore)
}

func (d *EDecoder) processHeadTimestampMsg(msgBuf *MsgBuffer) {

	reqID := msgBuf.decodeInt64()

	headTimestamp := msgBuf.decodeString()

	d.wrapper.HeadTimestamp(reqID, headTimestamp)
}

func (d *EDecoder) processHistogramDataMsg(msgBuf *MsgBuffer) {

	reqID := msgBuf.decodeInt64()

	numPoints := msgBuf.decodeInt64()
	data := make([]HistogramData, 0, numPoints)
	var i int64
	for i = 0; i < numPoints; i++ {
		p := HistogramData{}
		p.Price = msgBuf.decodeFloat64()
		p.Size = msgBuf.decodeDecimal()
		data = append(data, p)
	}

	d.wrapper.HistogramData(reqID, data)
}

func (d *EDecoder) processHistoricalDataUpdateMsg(msgBuf *MsgBuffer) {

	reqID := msgBuf.decodeInt64()

	bar := NewBar()
	bar.BarCount = msgBuf.decodeInt64()
	bar.Date = msgBuf.decodeString()
	bar.Open = msgBuf.decodeFloat64()
	bar.Close = msgBuf.decodeFloat64()
	bar.High = msgBuf.decodeFloat64()
	bar.Low = msgBuf.decodeFloat64()
	bar.Wap = msgBuf.decodeDecimal()
	bar.Volume = msgBuf.decodeDecimal()

	d.wrapper.HistoricalDataUpdate(reqID, &bar)
}

func (d *EDecoder) processRerouteMktDataReqMsg(msgBuf *MsgBuffer) {

	reqID := msgBuf.decodeInt64()

	conID := msgBuf.decodeInt64()
	exchange := msgBuf.decodeString()

	d.wrapper.RerouteMktDataReq(reqID, conID, exchange)
}

func (d *EDecoder) processRerouteMktDepthReqMsg(msgBuf *MsgBuffer) {

	reqID := msgBuf.decodeInt64()

	conID := msgBuf.decodeInt64()
	exchange := msgBuf.decodeString()

	d.wrapper.RerouteMktDepthReq(reqID, conID, exchange)
}

func (d *EDecoder) processMarketRuleMsg(msgBuf *MsgBuffer) {

	marketRuleID := msgBuf.decodeInt64()

	priceIncrementsCount := msgBuf.decodeInt64()
	priceIncrements := make([]PriceIncrement, 0, priceIncrementsCount)

	var i int64
	for i = 0; i < priceIncrementsCount; i++ {
		priceInc := NewPriceIncrement()
		priceInc.LowEdge = msgBuf.decodeFloat64()
		priceInc.Increment = msgBuf.decodeFloat64()
		priceIncrements = append(priceIncrements, priceInc)
	}

	d.wrapper.MarketRule(marketRuleID, priceIncrements)
}

func (d *EDecoder) processPnLMsg(msgBuf *MsgBuffer) {

	reqID := msgBuf.decodeInt64()

	dailyPnL := msgBuf.decodeFloat64()
	var unrealizedPnL float64
	var realizedPnL float64

	if d.serverVersion >= MIN_SERVER_VER_UNREALIZED_PNL {
		unrealizedPnL = msgBuf.decodeFloat64()
	}

	if d.serverVersion >= MIN_SERVER_VER_REALIZED_PNL {
		realizedPnL = msgBuf.decodeFloat64()
	}

	d.wrapper.Pnl(reqID, dailyPnL, unrealizedPnL, realizedPnL)
}

func (d *EDecoder) processPnLSingleMsg(msgBuf *MsgBuffer) {

	reqID := msgBuf.decodeInt64()

	pos := msgBuf.decodeDecimal()
	dailyPnL := msgBuf.decodeFloat64()
	var unrealizedPnL float64
	var realizedPnL float64

	if d.serverVersion >= MIN_SERVER_VER_UNREALIZED_PNL {
		unrealizedPnL = msgBuf.decodeFloat64()
	}

	if d.serverVersion >= MIN_SERVER_VER_REALIZED_PNL {
		realizedPnL = msgBuf.decodeFloat64()
	}

	value := msgBuf.decodeFloat64()

	d.wrapper.PnlSingle(reqID, pos, dailyPnL, unrealizedPnL, realizedPnL, value)
}

func (d *EDecoder) processHistoricalTicks(msgBuf *MsgBuffer) {

	reqID := msgBuf.decodeInt64()

	tickCount := msgBuf.decodeInt64()
	ticks := make([]HistoricalTick, 0, tickCount)

	var i int64
	for i = 0; i < tickCount; i++ {
		historicalTick := NewHistoricalTick()
		historicalTick.Time = msgBuf.decodeInt64()
		_ = msgBuf.decodeString()
		historicalTick.Price = msgBuf.decodeFloat64()
		historicalTick.Size = msgBuf.decodeDecimal()
		ticks = append(ticks, historicalTick)
	}

	done := msgBuf.decodeBool()

	d.wrapper.HistoricalTicks(reqID, ticks, done)
}

func (d *EDecoder) processHistoricalTicksBidAsk(msgBuf *MsgBuffer) {

	reqID := msgBuf.decodeInt64()

	tickCount := msgBuf.decodeInt64()
	ticks := make([]HistoricalTickBidAsk, 0, tickCount)

	var i int64
	for i = 0; i < tickCount; i++ {
		historicalTickBidAsk := NewHistoricalTickBidAsk()
		historicalTickBidAsk.Time = msgBuf.decodeInt64()
		mask := msgBuf.decodeInt64()
		tickAttribBidAsk := NewTickAttribBidAsk()
		tickAttribBidAsk.AskPastHigh = mask&1 != 0
		tickAttribBidAsk.BidPastLow = mask&2 != 0
		historicalTickBidAsk.TickAttirbBidAsk = tickAttribBidAsk
		historicalTickBidAsk.PriceBid = msgBuf.decodeFloat64()
		historicalTickBidAsk.PriceAsk = msgBuf.decodeFloat64()
		historicalTickBidAsk.SizeBid = msgBuf.decodeDecimal()
		historicalTickBidAsk.SizeAsk = msgBuf.decodeDecimal()
		ticks = append(ticks, historicalTickBidAsk)
	}

	done := msgBuf.decodeBool()

	d.wrapper.HistoricalTicksBidAsk(reqID, ticks, done)
}

func (d *EDecoder) processHistoricalTicksLast(msgBuf *MsgBuffer) {

	reqID := msgBuf.decodeInt64()

	tickCount := msgBuf.decodeInt64()
	ticks := make([]HistoricalTickLast, 0, tickCount)

	var i int64
	for i = 0; i < tickCount; i++ {
		historicalTickLast := NewHistoricalTickLast()
		historicalTickLast.Time = msgBuf.decodeInt64()

		mask := msgBuf.decodeInt64()
		tickAttribLast := NewTickAttribLast()
		tickAttribLast.PastLimit = mask&1 != 0
		tickAttribLast.Unreported = mask&2 != 0

		historicalTickLast.TickAttribLast = tickAttribLast
		historicalTickLast.Price = msgBuf.decodeFloat64()
		historicalTickLast.Size = msgBuf.decodeDecimal()
		historicalTickLast.Exchange = msgBuf.decodeString()
		historicalTickLast.SpecialConditions = msgBuf.decodeString()
		ticks = append(ticks, historicalTickLast)
	}

	done := msgBuf.decodeBool()

	d.wrapper.HistoricalTicksLast(reqID, ticks, done)
}

func (d *EDecoder) processTickByTickDataMsg(msgBuf *MsgBuffer) {

	reqID := msgBuf.decodeInt64()

	tickType := msgBuf.decodeInt64()
	time := msgBuf.decodeInt64()

	switch tickType {
	case 0: // None
	case 1, 2: // Last or AllLast
		price := msgBuf.decodeFloat64()
		size := msgBuf.decodeDecimal()
		mask := msgBuf.decodeInt64()

		tickAttribLast := NewTickAttribLast()
		tickAttribLast.PastLimit = mask&1 != 0
		tickAttribLast.Unreported = mask&2 != 0

		exchange := msgBuf.decodeString()
		specialConditions := msgBuf.decodeString()

		d.wrapper.TickByTickAllLast(reqID, tickType, time, price, size, tickAttribLast, exchange, specialConditions)

	case 3: // BidAsk
		bidPrice := msgBuf.decodeFloat64()
		askPrice := msgBuf.decodeFloat64()
		bidSize := msgBuf.decodeDecimal()
		askSize := msgBuf.decodeDecimal()
		mask := msgBuf.decodeInt64()

		tickAttribBidAsk := NewTickAttribBidAsk()
		tickAttribBidAsk.BidPastLow = mask&1 != 0
		tickAttribBidAsk.AskPastHigh = mask&2 != 0

		d.wrapper.TickByTickBidAsk(reqID, time, bidPrice, askPrice, bidSize, askSize, tickAttribBidAsk)

	case 4: // MidPoint
		midPoint := msgBuf.decodeFloat64()

		d.wrapper.TickByTickMidPoint(reqID, time, midPoint)
	}
}

func (d *EDecoder) processOrderBoundMsg(msgBuf *MsgBuffer) {

	permID := msgBuf.decodeInt64()
	clientId := msgBuf.decodeInt64()
	orderId := msgBuf.decodeInt64()

	d.wrapper.OrderBound(permID, clientId, orderId)
}

func (d *EDecoder) processCompletedOrderMsg(msgBuf *MsgBuffer) {

	order := NewOrder()
	contract := NewContract()
	orderState := NewOrderState()

	orderDecoder := &OrderDecoder{order, contract, orderState, Version(UNSET_INT), d.serverVersion}

	// read contract fields
	orderDecoder.decodeContractFields(msgBuf)

	// read order fields
	orderDecoder.decodeAction(msgBuf)
	orderDecoder.decodeTotalQuantity(msgBuf)
	orderDecoder.decodeOrderType(msgBuf)
	orderDecoder.decodeLmtPrice(msgBuf)
	orderDecoder.decodeAuxPrice(msgBuf)
	orderDecoder.decodeTIF(msgBuf)
	orderDecoder.decodeOcaGroup(msgBuf)
	orderDecoder.decodeAccount(msgBuf)
	orderDecoder.decodeOpenClose(msgBuf)
	orderDecoder.decodeOrigin(msgBuf)
	orderDecoder.decodeOrderRef(msgBuf)
	orderDecoder.decodePermId(msgBuf)
	orderDecoder.decodeOutsideRth(msgBuf)
	orderDecoder.decodeHidden(msgBuf)
	orderDecoder.decodeDiscretionaryAmount(msgBuf)
	orderDecoder.decodeGoodAfterTime(msgBuf)
	orderDecoder.decodeFAParams(msgBuf)
	orderDecoder.decodeModelCode(msgBuf)
	orderDecoder.decodeGoodTillDate(msgBuf)
	orderDecoder.decodeRule80A(msgBuf)
	orderDecoder.decodePercentOffset(msgBuf)
	orderDecoder.decodeSettlingFirm(msgBuf)
	orderDecoder.decodeShortSaleParams(msgBuf)
	orderDecoder.decodeBoxOrderParams(msgBuf)
	orderDecoder.decodePegToStkOrVolOrderParams(msgBuf)
	orderDecoder.decodeDisplaySize(msgBuf)
	orderDecoder.decodeSweepToFill(msgBuf)
	orderDecoder.decodeAllOrNone(msgBuf)
	orderDecoder.decodeMinQty(msgBuf)
	orderDecoder.decodeOcaType(msgBuf)
	orderDecoder.decodeTriggerMethod(msgBuf)
	orderDecoder.decodeVolOrderParams(msgBuf, false)
	orderDecoder.decodeTrailParams(msgBuf)
	orderDecoder.decodeComboLegs(msgBuf)
	orderDecoder.decodeSmartComboRoutingParams(msgBuf)
	orderDecoder.decodeScaleOrderParams(msgBuf)
	orderDecoder.decodeHedgeParams(msgBuf)
	orderDecoder.decodeClearingParams(msgBuf)
	orderDecoder.decodeNotHeld(msgBuf)
	orderDecoder.decodeDeltaNeutral(msgBuf)
	orderDecoder.decodeAlgoParams(msgBuf)
	orderDecoder.decodeSolicited(msgBuf)
	orderDecoder.decodeOrderStatus(msgBuf)
	orderDecoder.decodeVolRandomizeFlags(msgBuf)
	orderDecoder.decodePegBenchParams(msgBuf)
	orderDecoder.decodeConditions(msgBuf)
	orderDecoder.decodeStopPriceAndLmtPriceOffset(msgBuf)
	orderDecoder.decodeCashQty(msgBuf)
	orderDecoder.decodeDontUseAutoPriceForHedge(msgBuf)
	orderDecoder.decodeIsOmsContainer(msgBuf)
	orderDecoder.decodeAutoCancelDate(msgBuf)
	orderDecoder.decodeFilledQuantity(msgBuf)
	orderDecoder.decodeRefFuturesConId(msgBuf)
	orderDecoder.decodeAutoCancelParent(msgBuf, MIN_CLIENT_VER)
	orderDecoder.decodeShareholder(msgBuf)
	orderDecoder.decodeImbalanceOnly(msgBuf)
	orderDecoder.decodeRouteMarketableToBbo(msgBuf)
	orderDecoder.decodeParentPermId(msgBuf)
	orderDecoder.decodeCompletedTime(msgBuf)
	orderDecoder.decodeCompletedStatus(msgBuf)
	orderDecoder.decodePegBestPegMidOrderAttributes(msgBuf)
	orderDecoder.decodeCustomerAccount(msgBuf)
	orderDecoder.decodeProfessionalCustomer(msgBuf)

	d.wrapper.CompletedOrder(contract, order, orderState)
}

func (d *EDecoder) processCompletedOrdersEndMsg(*MsgBuffer) {
	d.wrapper.CompletedOrdersEnd()
}

func (d *EDecoder) processReplaceFAEndMsg(msgBuf *MsgBuffer) {

	reqID := msgBuf.decodeInt64()
	text := msgBuf.decodeString()

	d.wrapper.ReplaceFAEnd(reqID, text)
}

func (d *EDecoder) processWshMetaData(msgBuf *MsgBuffer) {

	reqID := msgBuf.decodeInt64()
	dataJSON := msgBuf.decodeString()

	d.wrapper.WshMetaData(reqID, dataJSON)
}

func (d *EDecoder) processWshEventData(msgBuf *MsgBuffer) {

	reqID := msgBuf.decodeInt64()
	dataJSON := msgBuf.decodeString()

	d.wrapper.WshEventData(reqID, dataJSON)
}

func (d *EDecoder) processHistoricalSchedule(msgBuf *MsgBuffer) {

	reqID := msgBuf.decodeInt64()
	startDateTime := msgBuf.decodeString()
	endDateTime := msgBuf.decodeString()
	timeZone := msgBuf.decodeString()
	sessionsCount := msgBuf.decodeInt64()
	sessions := make([]HistoricalSession, 0, sessionsCount)
	var i int64
	for i = 0; i < sessionsCount; i++ {
		historicalSession := NewHistoricalSession()
		historicalSession.StartDateTime = msgBuf.decodeString()
		historicalSession.EndDateTime = msgBuf.decodeString()
		historicalSession.RefDate = msgBuf.decodeString()
		sessions = append(sessions, historicalSession)
	}

	d.wrapper.HistoricalSchedule(reqID, startDateTime, endDateTime, timeZone, sessions)
}

func (d *EDecoder) processUserInfo(msgBuf *MsgBuffer) {

	reqID := msgBuf.decodeInt64()
	whiteBrandingId := msgBuf.decodeString()

	d.wrapper.UserInfo(reqID, whiteBrandingId)
}

//
//		Helpers
//

func (d *EDecoder) readLastTradeDate(msgBuf *MsgBuffer, contract *ContractDetails, isBond bool) {
	lastTradeDateOrContractMonth := msgBuf.decodeString()
	if lastTradeDateOrContractMonth != "" {
		var splitted []string
		if strings.Contains(lastTradeDateOrContractMonth, "-") {
			splitted = strings.Split(lastTradeDateOrContractMonth, "-")
		} else {
			splitted = strings.Split(lastTradeDateOrContractMonth, " ")
		}

		if len(splitted) > 0 {
			if isBond {
				contract.Maturity = splitted[0]
			} else {
				contract.Contract.LastTradeDateOrContractMonth = splitted[0]
			}
		}
		if len(splitted) > 1 {
			contract.LastTradeTime = splitted[1]
		}
		if isBond && len(splitted) > 2 {
			contract.TimeZoneID = splitted[2]
		}
	}
}
