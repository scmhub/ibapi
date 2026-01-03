/*
The EDecoder knows how to transform a message's payload into higher level IB message (eg: order info, mkt data, etc).
It will call the corresponding method from the EWrapper so that customer's code (eg: class derived from EWrapper) can make further use of the data.
*/
package ibapi

import (
	"strconv"
	"strings"

	"github.com/scmhub/ibapi/protobuf"
	"google.golang.org/protobuf/proto"
)

// EDecoder transforms a message's payload into higher level IB message.
type EDecoder struct {
	wrapper       EWrapper
	serverVersion Version
}

func (d *EDecoder) parseAndProcessMsg(msgBytes []byte) {

	msgBuf := NewMsgBuffer(msgBytes)

	if msgBuf.Len() == 0 {
		log.Warn().Msg("message has no fields")
		return
	}

	var msgID int64
	if d.serverVersion >= MIN_SERVER_VER_PROTOBUF {
		msgID = msgBuf.decodeRawInt64()
	} else {
		msgID = msgBuf.decodeInt64()
	}

	var useProtoBuf bool
	if msgID >= PROTOBUF_MSG_ID {
		useProtoBuf = true
		msgID -= PROTOBUF_MSG_ID
	}

	if useProtoBuf {
		switch msgID {
		case ORDER_STATUS:
			d.processOrderStatusMsgProtoBuf(msgBuf)
		case ERR_MSG:
			d.processErrorMsgProtoBuf(msgBuf)
		case OPEN_ORDER:
			d.processOpenOrderMsgProtoBuf(msgBuf)
		case EXECUTION_DATA:
			d.processExecutionDetailsMsgProtoBuf(msgBuf)
		case OPEN_ORDER_END:
			d.processOpenOrderEndMsgProtoBuf(msgBuf)
		case EXECUTION_DATA_END:
			d.processExecutionDetailsEndMsgProtoBuf(msgBuf)
		case COMPLETED_ORDER:
			d.processCompletedOrderMsgProtoBuf(msgBuf)
		case COMPLETED_ORDERS_END:
			d.processCompletedOrdersEndMsgProtoBuf(msgBuf)
		case ORDER_BOUND:
			d.processOrderBoundMsgProtoBuf(msgBuf)
		case CONTRACT_DATA:
			d.processContractDataMsgProtoBuf(msgBuf)
		case BOND_CONTRACT_DATA:
			d.processBondContractDataMsgProtoBuf(msgBuf)
		case CONTRACT_DATA_END:
			d.processContractDataEndMsgProtoBuf(msgBuf)
		case TICK_PRICE:
			d.processTickPriceMsgProtoBuf(msgBuf)
		case TICK_SIZE:
			d.processTickSizeMsgProtoBuf(msgBuf)
		case MARKET_DEPTH:
			d.processMarketDepthMsgProtoBuf(msgBuf)
		case MARKET_DEPTH_L2:
			d.processMarketDepthL2MsgProtoBuf(msgBuf)
		case TICK_OPTION_COMPUTATION:
			d.processTickOptionComputationMsgProtoBuf(msgBuf)
		case TICK_GENERIC:
			d.processTickGenericMsgProtoBuf(msgBuf)
		case TICK_STRING:
			d.processTickStringMsgProtoBuf(msgBuf)
		case TICK_SNAPSHOT_END:
			d.processTickSnapshotEndMsgProtoBuf(msgBuf)
		case MARKET_DATA_TYPE:
			d.processMarketDataTypeMsgProtoBuf(msgBuf)
		case TICK_REQ_PARAMS:
			d.processTickReqParamsMsgProtoBuf(msgBuf)
		case ACCT_VALUE:
			d.processAccountValueMsgProtoBuf(msgBuf)
		case PORTFOLIO_VALUE:
			d.processPortfolioValueMsgProtoBuf(msgBuf)
		case ACCT_UPDATE_TIME:
			d.processAcctUpdateTimeMsgProtoBuf(msgBuf)
		case ACCT_DOWNLOAD_END:
			d.processAccountDataEndMsgProtoBuf(msgBuf)
		case MANAGED_ACCTS:
			d.processManagedAccountsMsgProtoBuf(msgBuf)
		case POSITION_DATA:
			d.processPositionMsgProtoBuf(msgBuf)
		case POSITION_END:
			d.processPositionEndMsgProtoBuf(msgBuf)
		case ACCOUNT_SUMMARY:
			d.processAccountSummaryMsgProtoBuf(msgBuf)
		case ACCOUNT_SUMMARY_END:
			d.processAccountSummaryEndMsgProtoBuf(msgBuf)
		case POSITION_MULTI:
			d.processPositionMultiMsgProtoBuf(msgBuf)
		case POSITION_MULTI_END:
			d.processPositionMultiEndMsgProtoBuf(msgBuf)
		case ACCOUNT_UPDATE_MULTI:
			d.processAccountUpdateMultiMsgProtoBuf(msgBuf)
		case ACCOUNT_UPDATE_MULTI_END:
			d.processAccountUpdateMultiEndMsgProtoBuf(msgBuf)
		case HISTORICAL_DATA:
			d.processHistoricalDataMsgProtoBuf(msgBuf)
		case HISTORICAL_DATA_UPDATE:
			d.processHistoricalDataUpdateMsgProtoBuf(msgBuf)
		case HISTORICAL_DATA_END:
			d.processHistoricalDataEndMsgProtoBuf(msgBuf)
		case REAL_TIME_BARS:
			d.processRealTimeBarsMsgProtoBuf(msgBuf)
		case HEAD_TIMESTAMP:
			d.processHeadTimestampMsgProtoBuf(msgBuf)
		case HISTOGRAM_DATA:
			d.processHistogramDataMsgProtoBuf(msgBuf)
		case HISTORICAL_TICKS:
			d.processHistoricalTicksMsgProtoBuf(msgBuf)
		case HISTORICAL_TICKS_BID_ASK:
			d.processHistoricalTicksBidAskMsgProtoBuf(msgBuf)
		case HISTORICAL_TICKS_LAST:
			d.processHistoricalTicksLastMsgProtoBuf(msgBuf)
		case TICK_BY_TICK:
			d.processTickByTickMsgProtoBuf(msgBuf)
		case NEWS_BULLETINS:
			d.processNewsBulletinMsgProtoBuf(msgBuf)
		case NEWS_ARTICLE:
			d.processNewsArticleMsgProtoBuf(msgBuf)
		case NEWS_PROVIDERS:
			d.processNewsProvidersMsgProtoBuf(msgBuf)
		case HISTORICAL_NEWS:
			d.processHistoricalNewsMsgProtoBuf(msgBuf)
		case HISTORICAL_NEWS_END:
			d.processHistoricalNewsEndMsgProtoBuf(msgBuf)
		case WSH_META_DATA:
			d.processWshMetaDataMsgProtoBuf(msgBuf)
		case WSH_EVENT_DATA:
			d.processWshEventDataMsgProtoBuf(msgBuf)
		case TICK_NEWS:
			d.processTickNewsMsgProtoBuf(msgBuf)
		case SCANNER_PARAMETERS:
			d.processScannerParametersMsgProtoBuf(msgBuf)
		case SCANNER_DATA:
			d.processScannerDataMsgProtoBuf(msgBuf)
		case FUNDAMENTAL_DATA:
			d.processFundamentalsDataMsgProtoBuf(msgBuf)
		case PNL:
			d.processPnLMsgProtoBuf(msgBuf)
		case PNL_SINGLE:
			d.processPnLSingleMsgProtoBuf(msgBuf)
		case RECEIVE_FA:
			d.processReceiveFAMsgProtoBuf(msgBuf)
		case REPLACE_FA_END:
			d.processReplaceFAEndMsgProtoBuf(msgBuf)
		case COMMISSION_AND_FEES_REPORT:
			d.processCommissionAndFeesReportMsgProtoBuf(msgBuf)
		case HISTORICAL_SCHEDULE:
			d.processHistoricalScheduleMsgProtoBuf(msgBuf)
		case REROUTE_MKT_DATA_REQ:
			d.processRerouteMktDataReqMsgProtoBuf(msgBuf)
		case REROUTE_MKT_DEPTH_REQ:
			d.processRerouteMktDepthReqMsgProtoBuf(msgBuf)
		case SECURITY_DEFINITION_OPTION_PARAMETER:
			d.processSecurityDefinitionOptionParameterMsgProtoBuf(msgBuf)
		case SECURITY_DEFINITION_OPTION_PARAMETER_END:
			d.processSecurityDefinitionOptionParameterEndMsgProtoBuf(msgBuf)
		case SOFT_DOLLAR_TIERS:
			d.processSoftDollarTiersMsgProtoBuf(msgBuf)
		case FAMILY_CODES:
			d.processFamilyCodesMsgProtoBuf(msgBuf)
		case SYMBOL_SAMPLES:
			d.processSymbolSamplesMsgProtoBuf(msgBuf)
		case SMART_COMPONENTS:
			d.processSmartComponentsMsgProtoBuf(msgBuf)
		case MARKET_RULE:
			d.processMarketRuleMsgProtoBuf(msgBuf)
		case USER_INFO:
			d.processUserInfoMsgProtoBuf(msgBuf)
		case NEXT_VALID_ID:
			d.processNextValidIdMsgProtoBuf(msgBuf)
		case CURRENT_TIME:
			d.processCurrentTimeMsgProtoBuf(msgBuf)
		case CURRENT_TIME_IN_MILLIS:
			d.processCurrentTimeInMillisMsgProtoBuf(msgBuf)
		case VERIFY_MESSAGE_API:
			d.processVerifyMessageApiMsgProtoBuf(msgBuf)
		case VERIFY_COMPLETED:
			d.processVerifyCompletedMsgProtoBuf(msgBuf)
		case DISPLAY_GROUP_LIST:
			d.processDisplayGroupListMsgProtoBuf(msgBuf)
		case DISPLAY_GROUP_UPDATED:
			d.processDisplayGroupUpdatedMsgProtoBuf(msgBuf)
		case MKT_DEPTH_EXCHANGES:
			d.processMktDepthExchangesMsgProtoBuf(msgBuf)
		case CONFIG_RESPONSE:
			d.processConfigResponseMsgProtoBuf(msgBuf)
		default:
			d.wrapper.Error(msgID, currentTimeMillis(), UNKNOWN_ID.Code, UNKNOWN_ID.Msg, "")
		}
	} else {
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
			d.processErrorMsg(msgBuf)
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
		case COMMISSION_AND_FEES_REPORT:
			d.processCommissionAndFeesReportMsg(msgBuf)
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
		case HISTORICAL_DATA_END:
			d.processHistoricalDataEndMsg(msgBuf)
		case CURRENT_TIME_IN_MILLIS:
			d.processCurrentTimeInMillisMsg(msgBuf)
		default:
			d.wrapper.Error(msgID, currentTimeMillis(), BAD_MESSAGE.Code, BAD_MESSAGE.Msg, "")
		}
	}
}

func (d *EDecoder) processTickPriceMsg(msgBuf *MsgBuffer) {

	msgBuf.decode() // version

	reqID := msgBuf.decodeInt64()
	tickType := msgBuf.decodeInt64()
	price := msgBuf.decodeFloat64()
	size := msgBuf.decodeDecimal()   // ver 2 field
	attrMask := msgBuf.decodeInt64() // ver 3 field

	attrib := NewTickAttrib()
	attrib.CanAutoExecute = attrMask == 1

	if d.serverVersion >= MIN_SERVER_VER_PAST_LIMIT {
		attrib.CanAutoExecute = (attrMask & (1 << 0)) != 0
		attrib.PastLimit = (attrMask & (1 << 1)) != 0
		if d.serverVersion >= MIN_SERVER_VER_PRE_OPEN_BID_ASK {
			attrib.PreOpen = (attrMask & (1 << 2)) != 0
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

func (d *EDecoder) processTickPriceMsgProtoBuf(msgBuf *MsgBuffer) {

	var tickPriceProto protobuf.TickPrice
	err := proto.Unmarshal(msgBuf.bs, &tickPriceProto)
	if err != nil {
		log.Error().Err(err).Msg("failed to unmarshal TickPrice message")
		return
	}

	d.wrapper.TickPriceProtoBuf(&tickPriceProto)

	reqID := NO_VALID_ID
	if tickPriceProto.ReqId != nil {
		reqID = int64(tickPriceProto.GetReqId())
	}
	var tickType TickType
	if tickPriceProto.TickType != nil {
		tickType = TickType(tickPriceProto.GetTickType())
	}
	var price float64
	if tickPriceProto.Price != nil {
		price = tickPriceProto.GetPrice()
	}
	var size Decimal
	if tickPriceProto.Size != nil { // Assuming Size is named Size_ in proto to avoid conflict
		size = StringToDecimal(tickPriceProto.GetSize())
	}
	var attrMask int64
	if tickPriceProto.AttrMask != nil {
		attrMask = int64(tickPriceProto.GetAttrMask())
	}

	attrib := NewTickAttrib()
	attrib.CanAutoExecute = (attrMask & (1 << 0)) != 0 // Check the 0th bit
	attrib.PastLimit = (attrMask & (1 << 1)) != 0      // Check the 1st bit
	attrib.PreOpen = (attrMask & (1 << 2)) != 0        // Check the 2nd bit

	d.wrapper.TickPrice(reqID, tickType, price, attrib)

	// process size tick
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

	msgBuf.decode() // version

	reqID := msgBuf.decodeInt64()
	sizeTickType := msgBuf.decodeInt64()
	size := msgBuf.decodeDecimal()

	if sizeTickType != NOT_SET {
		d.wrapper.TickSize(reqID, sizeTickType, size)
	}
}

func (d *EDecoder) processTickSizeMsgProtoBuf(msgBuf *MsgBuffer) {
	var tickSizeProto protobuf.TickSize
	err := proto.Unmarshal(msgBuf.bs, &tickSizeProto)
	if err != nil {
		log.Error().Err(err).Msg("failed to unmarshal TickSize message")
		return
	}

	d.wrapper.TickSizeProtoBuf(&tickSizeProto)

	reqID := NO_VALID_ID
	if tickSizeProto.ReqId != nil {
		reqID = int64(tickSizeProto.GetReqId())
	}
	var sizeTickType int64
	if tickSizeProto.TickType != nil {
		sizeTickType = int64(tickSizeProto.GetTickType())
	}
	var size Decimal
	if tickSizeProto.Size != nil {
		size = StringToDecimal(tickSizeProto.GetSize())
	}

	if sizeTickType != NOT_SET { // NOT_SET would be a constant defined in your Go code
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

func (d *EDecoder) processTickOptionComputationMsgProtoBuf(msgBuf *MsgBuffer) {
	var tickOptionComputationProto protobuf.TickOptionComputation
	err := proto.Unmarshal(msgBuf.bs, &tickOptionComputationProto)
	if err != nil {
		log.Error().Err(err).Msg("failed to unmarshal TickOptionComputation message")
		return
	}

	d.wrapper.TickOptionComputationProtoBuf(&tickOptionComputationProto)

	reqID := NO_VALID_ID
	if tickOptionComputationProto.ReqId != nil {
		reqID = int64(tickOptionComputationProto.GetReqId())
	}

	var tickType TickType
	if tickOptionComputationProto.TickType != nil {
		tickType = TickType(tickOptionComputationProto.GetTickType())
	}

	var tickAttrib int64
	if tickOptionComputationProto.TickAttrib != nil {
		tickAttrib = int64(tickOptionComputationProto.GetTickAttrib())
	}

	var impliedVol float64
	if tickOptionComputationProto.ImpliedVol != nil {
		impliedVol = tickOptionComputationProto.GetImpliedVol()
		if impliedVol == -1 { // -1 is the "not computed" indicator
			impliedVol = UNSET_FLOAT
		}
	}

	var delta float64
	if tickOptionComputationProto.Delta != nil {
		delta = tickOptionComputationProto.GetDelta()
		if delta == -2 { // -2 is the "not computed" indicator
			delta = UNSET_FLOAT
		}
	}

	var optPrice float64
	if tickOptionComputationProto.OptPrice != nil {
		optPrice = tickOptionComputationProto.GetOptPrice()
		if optPrice == -1 { // -1 is the "not computed" indicator
			optPrice = UNSET_FLOAT
		}
	}

	var pvDividend float64
	if tickOptionComputationProto.PvDividend != nil {
		pvDividend = tickOptionComputationProto.GetPvDividend()
		if pvDividend == -1 { // -1 is the "not computed" indicator
			pvDividend = UNSET_FLOAT
		}
	}

	var gamma float64
	if tickOptionComputationProto.Gamma != nil {
		gamma = tickOptionComputationProto.GetGamma()
		if gamma == -2 { // -2 is the "not yet computed" indicator
			gamma = UNSET_FLOAT
		}
	}

	var vega float64
	if tickOptionComputationProto.Vega != nil {
		vega = tickOptionComputationProto.GetVega()
		if vega == -2 { // -2 is the "not yet computed" indicator
			vega = UNSET_FLOAT
		}
	}

	var theta float64
	if tickOptionComputationProto.Theta != nil {
		theta = tickOptionComputationProto.GetTheta()
		if theta == -2 { // -2 is the "not yet computed" indicator
			theta = UNSET_FLOAT
		}
	}

	var undPrice float64
	if tickOptionComputationProto.UndPrice != nil {
		undPrice = tickOptionComputationProto.GetUndPrice()
		if undPrice == -1 { // -1 is the "not computed" indicator
			undPrice = UNSET_FLOAT
		}
	}

	d.wrapper.TickOptionComputation(reqID, tickType, tickAttrib, impliedVol, delta, optPrice, pvDividend, gamma, vega, theta, undPrice)
}

func (d *EDecoder) processTickGenericMsg(msgBuf *MsgBuffer) {

	msgBuf.decode() // version

	reqID := msgBuf.decodeInt64()
	tickType := msgBuf.decodeInt64()
	value := msgBuf.decodeFloat64()

	d.wrapper.TickGeneric(reqID, tickType, value)
}

func (d *EDecoder) processTickGenericMsgProtoBuf(msgBuf *MsgBuffer) {
	var tickGenericProto protobuf.TickGeneric
	err := proto.Unmarshal(msgBuf.bs, &tickGenericProto)
	if err != nil {
		log.Error().Err(err).Msg("failed to unmarshal TickGeneric message")
		return
	}

	d.wrapper.TickGenericProtoBuf(&tickGenericProto)

	reqID := NO_VALID_ID
	if tickGenericProto.ReqId != nil {
		reqID = int64(tickGenericProto.GetReqId())
	}

	var tickType TickType
	if tickGenericProto.TickType != nil {
		tickType = TickType(tickGenericProto.GetTickType())
	}

	var value float64
	if tickGenericProto.Value != nil {
		value = tickGenericProto.GetValue()
	}

	d.wrapper.TickGeneric(reqID, tickType, value)
}

func (d *EDecoder) processTickStringMsg(msgBuf *MsgBuffer) {

	msgBuf.decode() // version

	reqID := msgBuf.decodeInt64()
	tickType := msgBuf.decodeInt64()
	value := msgBuf.decodeString()

	d.wrapper.TickString(reqID, tickType, value)
}

func (d *EDecoder) processTickStringMsgProtoBuf(msgBuf *MsgBuffer) {
	var tickStringProto protobuf.TickString
	err := proto.Unmarshal(msgBuf.bs, &tickStringProto)
	if err != nil {
		log.Error().Err(err).Msg("failed to unmarshal TickString message")
		return
	}

	d.wrapper.TickStringProtoBuf(&tickStringProto)

	reqID := NO_VALID_ID
	if tickStringProto.ReqId != nil {
		reqID = int64(tickStringProto.GetReqId())
	}
	var tickType int64
	if tickStringProto.TickType != nil {
		tickType = int64(tickStringProto.GetTickType())
	}
	var value string
	if tickStringProto.Value != nil {
		value = tickStringProto.GetValue()
	}

	d.wrapper.TickString(reqID, tickType, value)
}

func (d *EDecoder) processTickEfpMsg(msgBuf *MsgBuffer) {

	msgBuf.decode() // version

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
		msgBuf.decode() // version
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

func (d *EDecoder) processOrderStatusMsgProtoBuf(msgBuf *MsgBuffer) {

	var orderStatusProto protobuf.OrderStatus
	err := proto.Unmarshal(msgBuf.bs, &orderStatusProto)
	if err != nil {
		log.Error().Err(err).Msg("failed to unmarshal OrderStatus message")
		return
	}

	d.wrapper.OrderStatusProtoBuf(&orderStatusProto)

	var orderID int64
	if orderStatusProto.OrderId != nil {
		orderID = int64(orderStatusProto.GetOrderId())
	}
	var status string
	if orderStatusProto.Status != nil {
		status = orderStatusProto.GetStatus()
	}
	var filled Decimal
	if orderStatusProto.Filled != nil {
		filled = StringToDecimal(orderStatusProto.GetFilled())
	}
	var remaining Decimal
	if orderStatusProto.Remaining != nil {
		remaining = StringToDecimal(orderStatusProto.GetRemaining())
	}
	var avgFillPrice float64
	if orderStatusProto.AvgFillPrice != nil {
		avgFillPrice = orderStatusProto.GetAvgFillPrice()
	}
	var permID int64
	if orderStatusProto.PermId != nil {
		permID = int64(orderStatusProto.GetPermId())
	}
	var parentID int64
	if orderStatusProto.ParentId != nil {
		parentID = int64(orderStatusProto.GetParentId())
	}
	var lastFillPrice float64
	if orderStatusProto.LastFillPrice != nil {
		lastFillPrice = orderStatusProto.GetLastFillPrice()
	}
	var clientID int64
	if orderStatusProto.ClientId != nil {
		clientID = int64(orderStatusProto.GetClientId())
	}
	var whyHeld string
	if orderStatusProto.WhyHeld != nil {
		whyHeld = orderStatusProto.GetWhyHeld()
	}
	var mktCapPrice float64
	if orderStatusProto.MktCapPrice != nil {
		mktCapPrice = orderStatusProto.GetMktCapPrice()
	}

	d.wrapper.OrderStatus(orderID, status, filled, remaining, avgFillPrice, permID, parentID, lastFillPrice, clientID, whyHeld, mktCapPrice)
}

func (d *EDecoder) processErrorMsg(msgBuf *MsgBuffer) {

	if d.serverVersion < MIN_SERVER_VER_ERROR_TIME {
		msgBuf.decode() // version
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

func (d *EDecoder) processErrorMsgProtoBuf(msgBuf *MsgBuffer) {

	var errorMessageProto protobuf.ErrorMessage
	err := proto.Unmarshal(msgBuf.bs, &errorMessageProto)
	if err != nil {
		log.Error().Err(err).Msg("failed to unmarshal ErrorMessage")
		return
	}

	d.wrapper.ErrorProtoBuf(&errorMessageProto)

	var reqID int64
	if errorMessageProto.Id != nil {
		reqID = int64(errorMessageProto.GetId())
	}
	var errorTime int64
	if errorMessageProto.ErrorTime != nil {
		errorTime = int64(errorMessageProto.GetErrorTime())
	}
	var errorCode int64
	if errorMessageProto.ErrorCode != nil {
		errorCode = int64(errorMessageProto.GetErrorCode())
	}
	var errorString string
	if errorMessageProto.ErrorMsg != nil {
		errorString = errorMessageProto.GetErrorMsg()
	}
	var advancedOrderRejectJson string
	if errorMessageProto.AdvancedOrderRejectJson != nil {
		advancedOrderRejectJson = errorMessageProto.GetAdvancedOrderRejectJson()
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
	orderDecoder.decodeWhatIfInfoAndCommissionAndFees(msgBuf)
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
	orderDecoder.decodeSubmitter(msgBuf)
	orderDecoder.decodeImbalanceOnly(msgBuf, MIN_SERVER_VER_IMBALANCE_ONLY)

	d.wrapper.OpenOrder(order.OrderID, contract, order, orderState)
}

func (d *EDecoder) processOpenOrderMsgProtoBuf(msgBuf *MsgBuffer) {

	var openOrderProto protobuf.OpenOrder
	err := proto.Unmarshal(msgBuf.bs, &openOrderProto)
	if err != nil {
		log.Error().Err(err).Msg("failed to unmarshal OpenOrder message")
		return
	}

	d.wrapper.OpenOrderProtoBuf(&openOrderProto)

	var orderID int64
	if openOrderProto.OrderId != nil {
		orderID = int64(openOrderProto.GetOrderId())
	}

	var contract *Contract
	if openOrderProto.Contract != nil {
		contract = decodeContract(openOrderProto.GetContract())
	}
	var order *Order
	if openOrderProto.Order != nil {
		order = decodeOrder(orderID, openOrderProto.GetContract(), openOrderProto.GetOrder())
	}
	var orderState *OrderState
	if openOrderProto.OrderState != nil {
		orderState = decodeOrderState(openOrderProto.GetOrderState())
	}

	d.wrapper.OpenOrder(orderID, contract, order, orderState)
}

func (d *EDecoder) processAcctValueMsg(msgBuf *MsgBuffer) {

	msgBuf.decode() // version

	tag := msgBuf.decodeString()
	val := msgBuf.decodeString()
	currency := msgBuf.decodeString()
	accountName := msgBuf.decodeString()

	d.wrapper.UpdateAccountValue(tag, val, currency, accountName)
}

func (d *EDecoder) processAccountValueMsgProtoBuf(msgBuf *MsgBuffer) {

	var accountValueProto protobuf.AccountValue
	err := proto.Unmarshal(msgBuf.bs, &accountValueProto)
	if err != nil {
		log.Error().Err(err).Msg("failed to unmarshal AccountValue message")
		return
	}

	d.wrapper.UpdateAccountValueProtoBuf(&accountValueProto)

	var key string
	if accountValueProto.Key != nil {
		key = accountValueProto.GetKey()
	}

	var value string
	if accountValueProto.Value != nil {
		value = accountValueProto.GetValue()
	}

	var currency string
	if accountValueProto.Currency != nil {
		currency = accountValueProto.GetCurrency()
	}

	var accountName string
	if accountValueProto.AccountName != nil {
		accountName = accountValueProto.GetAccountName()
	}

	d.wrapper.UpdateAccountValue(key, value, currency, accountName)
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

func (d *EDecoder) processPortfolioValueMsgProtoBuf(msgBuf *MsgBuffer) {

	var portfolioValueProto protobuf.PortfolioValue
	err := proto.Unmarshal(msgBuf.bs, &portfolioValueProto)
	if err != nil {
		log.Error().Err(err).Msg("failed to unmarshal PortfolioValue message")
		return
	}

	d.wrapper.UpdatePortfolioProtoBuf(&portfolioValueProto)

	if portfolioValueProto.Contract == nil {
		return
	}

	contract := decodeContract(portfolioValueProto.GetContract())

	var position Decimal
	if portfolioValueProto.Position != nil {
		position = StringToDecimal(portfolioValueProto.GetPosition())
	} else {
		position = UNSET_DECIMAL
	}

	var marketPrice float64
	if portfolioValueProto.MarketPrice != nil {
		marketPrice = portfolioValueProto.GetMarketPrice()
	}

	var marketValue float64
	if portfolioValueProto.MarketValue != nil {
		marketValue = portfolioValueProto.GetMarketValue()
	}

	var averageCost float64
	if portfolioValueProto.AverageCost != nil {
		averageCost = portfolioValueProto.GetAverageCost()
	}

	var unrealizedPNL float64
	if portfolioValueProto.UnrealizedPNL != nil {
		unrealizedPNL = portfolioValueProto.GetUnrealizedPNL()
	}

	var realizedPNL float64
	if portfolioValueProto.RealizedPNL != nil {
		realizedPNL = portfolioValueProto.GetRealizedPNL()
	}

	var accountName string
	if portfolioValueProto.AccountName != nil {
		accountName = portfolioValueProto.GetAccountName()
	}

	d.wrapper.UpdatePortfolio(contract, position, marketPrice, marketValue, averageCost, unrealizedPNL, realizedPNL, accountName)
}

func (d *EDecoder) processAcctUpdateTimeMsg(msgBuf *MsgBuffer) {

	msgBuf.decode() // version

	timeStamp := msgBuf.decodeString()

	d.wrapper.UpdateAccountTime(timeStamp)
}

func (d *EDecoder) processAcctUpdateTimeMsgProtoBuf(msgBuf *MsgBuffer) {

	var accountUpdateTimeProto protobuf.AccountUpdateTime
	err := proto.Unmarshal(msgBuf.bs, &accountUpdateTimeProto)
	if err != nil {
		log.Error().Err(err).Msg("failed to unmarshal AccountUpdateTime message")
		return
	}

	d.wrapper.UpdateAccountTimeProtoBuf(&accountUpdateTimeProto)

	var timeStamp string
	if accountUpdateTimeProto.TimeStamp != nil {
		timeStamp = accountUpdateTimeProto.GetTimeStamp()
	}

	d.wrapper.UpdateAccountTime(timeStamp)
}

func (d *EDecoder) processNextValidIdMsg(msgBuf *MsgBuffer) {

	msgBuf.decode() // version

	reqID := msgBuf.decodeInt64()

	d.wrapper.NextValidID(reqID)
}

func (d *EDecoder) processNextValidIdMsgProtoBuf(msgBuf *MsgBuffer) {

	var protoMsg protobuf.NextValidId
	if err := proto.Unmarshal(msgBuf.bs, &protoMsg); err != nil {
		log.Error().Err(err).Msg("failed to unmarshal NextValidId")
		return
	}

	d.wrapper.NextValidIdProtoBuf(&protoMsg)

	id := NO_VALID_ID
	if protoMsg.OrderId != nil {
		id = int64(protoMsg.GetOrderId())
	}

	d.wrapper.NextValidID(id)
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
	d.decodeLastTradeDate(msgBuf, cd, false)
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

func (d *EDecoder) processContractDataMsgProtoBuf(msgBuf *MsgBuffer) {
	var contractDataProto protobuf.ContractData
	err := proto.Unmarshal(msgBuf.bs, &contractDataProto)
	if err != nil {
		log.Error().Err(err).Msg("failed to unmarshal ContractData message")
		return
	}

	d.wrapper.ContractDataProtoBuf(&contractDataProto)

	reqID := NO_VALID_ID
	if contractDataProto.ReqId != nil {
		reqID = int64(contractDataProto.GetReqId())
	}
	var contractDetails *ContractDetails
	if contractDataProto.Contract != nil && contractDataProto.ContractDetails != nil {
		contractDetails = decodeContractDetails(contractDataProto.GetContract(), contractDataProto.GetContractDetails(), false)
	}

	d.wrapper.ContractDetails(reqID, contractDetails)
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
	d.decodeLastTradeDate(msgBuf, contract, true)
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

	if d.serverVersion >= MIN_SERVER_VER_SUBMITTER {
		execution.Submitter = msgBuf.decodeString()
	}

	d.wrapper.ExecDetails(reqID, contract, execution)
}

func (d *EDecoder) processBondContractDataMsgProtoBuf(msgBuf *MsgBuffer) {
	var contractDataProto protobuf.ContractData
	err := proto.Unmarshal(msgBuf.bs, &contractDataProto)
	if err != nil {
		log.Error().Err(err).Msg("failed to unmarshal BondContractData message")
		return
	}

	d.wrapper.BondContractDataProtoBuf(&contractDataProto)

	reqID := NO_VALID_ID
	if contractDataProto.ReqId != nil {
		reqID = int64(contractDataProto.GetReqId())
	}
	var contractDetails *ContractDetails
	if contractDataProto.Contract != nil && contractDataProto.ContractDetails != nil {
		contractDetails = decodeContractDetails(contractDataProto.GetContract(), contractDataProto.GetContractDetails(), true)
	}

	d.wrapper.BondContractDetails(reqID, contractDetails)
}

func (d *EDecoder) processExecutionDetailsMsgProtoBuf(msgBuf *MsgBuffer) {

	var executionDetailsProto protobuf.ExecutionDetails
	err := proto.Unmarshal(msgBuf.Bytes(), &executionDetailsProto)
	if err != nil {
		log.Panic().Err(err).Msg("processExecutionDetailsMsgProtoBuf unmarshal error")
	}

	d.wrapper.ExecDetailsProtoBuf(&executionDetailsProto)

	var reqID int64 = int64(executionDetailsProto.GetReqId())

	var contract *Contract
	if executionDetailsProto.Contract != nil {
		contract = decodeContract(executionDetailsProto.GetContract())
	}

	var execution *Execution
	if executionDetailsProto.Execution != nil {
		execution = decodeExecution(executionDetailsProto.GetExecution())
	}

	d.wrapper.ExecDetails(reqID, contract, execution)
}

func (d *EDecoder) processMarketDepthMsg(msgBuf *MsgBuffer) {

	msgBuf.decode() // version

	tickerID := msgBuf.decodeInt64()

	position := msgBuf.decodeInt64()
	operation := msgBuf.decodeInt64()
	side := msgBuf.decodeInt64()
	price := msgBuf.decodeFloat64()
	size := msgBuf.decodeDecimal()

	d.wrapper.UpdateMktDepth(tickerID, position, operation, side, price, size)
}

func (d *EDecoder) processMarketDepthMsgProtoBuf(msgBuf *MsgBuffer) {
	var marketDepthProto protobuf.MarketDepth
	err := proto.Unmarshal(msgBuf.bs, &marketDepthProto)
	if err != nil {
		log.Error().Err(err).Msg("failed to unmarshal MarketDepth message")
		return
	}

	d.wrapper.UpdateMarketDepthProtoBuf(&marketDepthProto) // Assuming you have a wrapper method for the raw protobuf message

	reqID := NO_VALID_ID
	if marketDepthProto.ReqId != nil {
		reqID = int64(marketDepthProto.GetReqId())
	}
	if marketDepthProto.GetMarketDepthData() == nil {
		return
	}
	marketDepthData := marketDepthProto.GetMarketDepthData()
	var position int64
	if marketDepthData.Position != nil {
		position = int64(marketDepthData.GetPosition())
	}
	var operation int64
	if marketDepthData.Operation != nil {
		operation = int64(marketDepthData.GetOperation())
	}
	var side int64
	if marketDepthData.Side != nil {
		side = int64(marketDepthData.GetSide())
	}
	var price float64
	if marketDepthData.Price != nil {
		price = marketDepthData.GetPrice()
	}
	var size Decimal
	if marketDepthData.Size != nil {
		size = StringToDecimal(marketDepthData.GetSize())
	}

	d.wrapper.UpdateMktDepth(reqID, position, operation, side, price, size)
}

func (d *EDecoder) processMarketDepthL2Msg(msgBuf *MsgBuffer) {

	msgBuf.decode() // version

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

func (d *EDecoder) processMarketDepthL2MsgProtoBuf(msgBuf *MsgBuffer) {
	var marketDepthL2Proto protobuf.MarketDepthL2
	err := proto.Unmarshal(msgBuf.bs, &marketDepthL2Proto)
	if err != nil {
		log.Error().Err(err).Msg("failed to unmarshal MarketDepthL2 message")
		return
	}

	d.wrapper.UpdateMarketDepthL2ProtoBuf(&marketDepthL2Proto) // Assuming you have a wrapper method for the raw protobuf message

	reqID := NO_VALID_ID
	if marketDepthL2Proto.ReqId != nil {
		reqID = int64(marketDepthL2Proto.GetReqId())
	}
	if marketDepthL2Proto.GetMarketDepthData() == nil {
		return
	}
	marketDepthData := marketDepthL2Proto.GetMarketDepthData()
	var position int64
	if marketDepthData.Position != nil {
		position = int64(marketDepthData.GetPosition())
	}
	var marketMaker string
	if marketDepthData.MarketMaker != nil {
		marketMaker = marketDepthData.GetMarketMaker()
	}
	var operation int64
	if marketDepthData.Operation != nil {
		operation = int64(marketDepthData.GetOperation())
	}
	var side int64
	if marketDepthData.Side != nil {
		side = int64(marketDepthData.GetSide())
	}
	var price float64
	if marketDepthData.Price != nil {
		price = marketDepthData.GetPrice()
	}
	var size Decimal
	if marketDepthData.Size != nil {
		size = StringToDecimal(marketDepthData.GetSize())
	}
	var isSmartDepth bool
	if marketDepthData.IsSmartDepth != nil {
		isSmartDepth = marketDepthData.GetIsSmartDepth()
	}

	d.wrapper.UpdateMktDepthL2(reqID, position, marketMaker, operation, side, price, size, isSmartDepth)
}

func (d *EDecoder) processNewsBulletinsMsg(msgBuf *MsgBuffer) {

	msgBuf.decode() // version

	msgID := msgBuf.decodeInt64()
	msgType := msgBuf.decodeInt64()
	newsMessage := msgBuf.decodeString()
	originExch := msgBuf.decodeString()

	d.wrapper.UpdateNewsBulletin(msgID, msgType, newsMessage, originExch)
}

func (d *EDecoder) processNewsBulletinMsgProtoBuf(msgBuf *MsgBuffer) {
	var newsBulletinProto protobuf.NewsBulletin
	if err := proto.Unmarshal(msgBuf.bs, &newsBulletinProto); err != nil {
		log.Error().Err(err).Msg("failed to unmarshal NewsBulletin message")
		return
	}
	d.wrapper.UpdateNewsBulletinProtoBuf(&newsBulletinProto)

	msgID := NO_VALID_ID
	if newsBulletinProto.NewsMsgId != nil {
		msgID = int64(newsBulletinProto.GetNewsMsgId())
	}
	var msgType int64
	if newsBulletinProto.NewsMsgType != nil {
		msgType = int64(newsBulletinProto.GetNewsMsgType())
	}
	newsMessage := ""
	if newsBulletinProto.NewsMessage != nil {
		newsMessage = newsBulletinProto.GetNewsMessage()
	}
	originExch := ""
	if newsBulletinProto.OriginatingExch != nil {
		originExch = newsBulletinProto.GetOriginatingExch()
	}

	d.wrapper.UpdateNewsBulletin(msgID, msgType, newsMessage, originExch)
}

func (d *EDecoder) processManagedAcctsMsg(msgBuf *MsgBuffer) {

	msgBuf.decode() // version

	accountsNames := msgBuf.decodeString()
	accountsList := strings.Split(accountsNames, ",")

	d.wrapper.ManagedAccounts(accountsList)
}

func (d *EDecoder) processManagedAccountsMsgProtoBuf(msgBuf *MsgBuffer) {

	var managedAccountsProto protobuf.ManagedAccounts
	err := proto.Unmarshal(msgBuf.bs, &managedAccountsProto)
	if err != nil {
		log.Error().Err(err).Msg("failed to unmarshal ManagedAccounts message")
		return
	}

	d.wrapper.ManagedAccountsProtoBuf(&managedAccountsProto)

	var accounts string
	if managedAccountsProto.AccountsList != nil {
		accounts = managedAccountsProto.GetAccountsList()
	}
	accountsList := strings.Split(accounts, ",")

	d.wrapper.ManagedAccounts(accountsList)
}

func (d *EDecoder) processReceiveFaMsg(msgBuf *MsgBuffer) {

	msgBuf.decode() // version

	faDataType := FaDataType(msgBuf.decodeInt64())
	cxml := msgBuf.decodeString()

	d.wrapper.ReceiveFA(faDataType, cxml)
}

func (d *EDecoder) processReceiveFAMsgProtoBuf(msgBuf *MsgBuffer) {
	var receiveFAProto protobuf.ReceiveFA
	err := proto.Unmarshal(msgBuf.bs, &receiveFAProto)
	if err != nil {
		log.Error().Err(err).Msg("failed to unmarshal ReceiveFA message")
		return
	}

	d.wrapper.ReceiveFAProtoBuf(&receiveFAProto)

	faDataTypeInt := 0
	if receiveFAProto.FaDataType != nil {
		faDataTypeInt = int(receiveFAProto.GetFaDataType())
	}
	xml := ""
	if receiveFAProto.Xml != nil {
		xml = receiveFAProto.GetXml()
	}

	d.wrapper.ReceiveFA(FaDataType(faDataTypeInt), xml)
}

func (d *EDecoder) processHistoricalDataMsg(msgBuf *MsgBuffer) {

	if d.serverVersion < MIN_SERVER_VER_SYNT_REALTIME_BARS {
		msgBuf.decode() // version
	}

	reqID := msgBuf.decodeInt64()

	var startDateStr, endDateStr string
	if d.serverVersion < MIN_SERVER_VER_HISTORICAL_DATA_END {
		startDateStr = msgBuf.decodeString()
		endDateStr = msgBuf.decodeString()
	}

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
			msgBuf.decode()
		}

		bar.BarCount = msgBuf.decodeInt64()

		d.wrapper.HistoricalData(reqID, &bar)
	}

	if d.serverVersion < MIN_SERVER_VER_HISTORICAL_DATA_END {
		d.wrapper.HistoricalDataEnd(reqID, startDateStr, endDateStr)
	}
}

func (d *EDecoder) processHistoricalDataMsgProtoBuf(msgBuf *MsgBuffer) {
	var historicalDataProto protobuf.HistoricalData
	if err := proto.Unmarshal(msgBuf.bs, &historicalDataProto); err != nil {
		log.Error().Err(err).Msg("failed to unmarshal HistoricalData message")
		return
	}

	d.wrapper.HistoricalDataProtoBuf(&historicalDataProto)

	reqID := NO_VALID_ID
	if historicalDataProto.ReqId != nil {
		reqID = int64(historicalDataProto.GetReqId())
	}

	for _, barProto := range historicalDataProto.GetHistoricalDataBars() {
		bar := decodeHistoricalDataBar(barProto)
		d.wrapper.HistoricalData(reqID, bar)
	}
}

func (d *EDecoder) processHistoricalDataEndMsg(msgBuf *MsgBuffer) {

	reqID := msgBuf.decodeInt64()
	startDateStr := msgBuf.decodeString()
	endDateStr := msgBuf.decodeString()

	d.wrapper.HistoricalDataEnd(reqID, startDateStr, endDateStr)
}

func (d *EDecoder) processHistoricalDataEndMsgProtoBuf(msgBuf *MsgBuffer) {
	var historicalDataEndProto protobuf.HistoricalDataEnd
	if err := proto.Unmarshal(msgBuf.bs, &historicalDataEndProto); err != nil {
		log.Error().Err(err).Msg("failed to unmarshal HistoricalDataEnd message")
		return
	}

	d.wrapper.HistoricalDataEndProtoBuf(&historicalDataEndProto)

	reqID := NO_VALID_ID
	if historicalDataEndProto.ReqId != nil {
		reqID = int64(historicalDataEndProto.GetReqId())
	}

	var startDate, endDate string
	if historicalDataEndProto.StartDateStr != nil {
		startDate = historicalDataEndProto.GetStartDateStr()
	}
	if historicalDataEndProto.EndDateStr != nil {
		endDate = historicalDataEndProto.GetEndDateStr()
	}

	d.wrapper.HistoricalDataEnd(reqID, startDate, endDate)
}

func (d *EDecoder) processScannerDataMsg(msgBuf *MsgBuffer) {

	msgBuf.decode() // version

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

func (d *EDecoder) processScannerDataMsgProtoBuf(msgBuf *MsgBuffer) {
	var scannerDataProto protobuf.ScannerData
	if err := proto.Unmarshal(msgBuf.bs, &scannerDataProto); err != nil {
		log.Error().Err(err).Msg("failed to unmarshal ScannerData message")
		return
	}

	d.wrapper.ScannerDataProtoBuf(&scannerDataProto)

	reqID := NO_VALID_ID
	if scannerDataProto.ReqId != nil {
		reqID = int64(scannerDataProto.GetReqId())
	}

	for _, element := range scannerDataProto.GetScannerDataElement() {
		rank := int64(element.GetRank())

		// Decode contract details
		var cd ContractDetails
		if element.Contract != nil {
			contract := decodeContract(element.GetContract())
			cd.Contract = *contract
		}
		cd.MarketName = element.GetMarketName()

		distance := element.GetDistance()
		benchmark := element.GetBenchmark()
		projection := element.GetProjection()
		comboKey := element.GetComboKey()

		d.wrapper.ScannerData(reqID, rank, &cd, distance, benchmark, projection, comboKey)
	}

	d.wrapper.ScannerDataEnd(reqID)
}

func (d *EDecoder) processScannerParametersMsg(msgBuf *MsgBuffer) {

	msgBuf.decode() // version

	xml := msgBuf.decodeString()

	d.wrapper.ScannerParameters(xml)
}

func (d *EDecoder) processScannerParametersMsgProtoBuf(msgBuf *MsgBuffer) {
	var scannerParametersProto protobuf.ScannerParameters
	if err := proto.Unmarshal(msgBuf.bs, &scannerParametersProto); err != nil {
		log.Error().Err(err).Msg("failed to unmarshal ScannerParameters message")
		return
	}

	d.wrapper.ScannerParametersProtoBuf(&scannerParametersProto)

	xml := ""
	if scannerParametersProto.Xml != nil {
		xml = scannerParametersProto.GetXml()
	}

	d.wrapper.ScannerParameters(xml)
}

func (d *EDecoder) processCurrentTimeMsg(msgBuf *MsgBuffer) {

	msgBuf.decode() // version

	t := msgBuf.decodeInt64()

	d.wrapper.CurrentTime(t)
}

func (d *EDecoder) processCurrentTimeMsgProtoBuf(msgBuf *MsgBuffer) {

	var protoMsg protobuf.CurrentTime
	if err := proto.Unmarshal(msgBuf.bs, &protoMsg); err != nil {
		log.Error().Err(err).Msg("failed to unmarshal CurrentTime")
		return
	}

	d.wrapper.CurrentTimeProtoBuf(&protoMsg)

	ts := int64(0)
	if protoMsg.CurrentTime != nil {
		ts = protoMsg.GetCurrentTime()
	}

	d.wrapper.CurrentTime(ts)
}

func (d *EDecoder) processRealTimeBarsMsg(msgBuf *MsgBuffer) {

	msgBuf.decode() // version

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

func (d *EDecoder) processRealTimeBarsMsgProtoBuf(msgBuf *MsgBuffer) {
	var realTimeBarTickProto protobuf.RealTimeBarTick
	if err := proto.Unmarshal(msgBuf.bs, &realTimeBarTickProto); err != nil {
		log.Error().Err(err).Msg("failed to unmarshal RealTimeBarTick message")
		return
	}

	d.wrapper.RealTimeBarTickProtoBuf(&realTimeBarTickProto)

	reqID := NO_VALID_ID
	if realTimeBarTickProto.ReqId != nil {
		reqID = int64(realTimeBarTickProto.GetReqId())
	}

	var t int64
	if realTimeBarTickProto.Time != nil {
		t = realTimeBarTickProto.GetTime()
	}

	var open, high, low, close float64
	if realTimeBarTickProto.Open != nil {
		open = realTimeBarTickProto.GetOpen()
	}
	if realTimeBarTickProto.High != nil {
		high = realTimeBarTickProto.GetHigh()
	}
	if realTimeBarTickProto.Low != nil {
		low = realTimeBarTickProto.GetLow()
	}
	if realTimeBarTickProto.Close != nil {
		close = realTimeBarTickProto.GetClose()
	}

	var volume, wap Decimal
	if realTimeBarTickProto.Volume != nil {
		volume = StringToDecimal(realTimeBarTickProto.GetVolume())
	}
	if realTimeBarTickProto.WAP != nil {
		wap = StringToDecimal(realTimeBarTickProto.GetWAP())
	}

	var count int64
	if realTimeBarTickProto.Count != nil {
		count = int64(realTimeBarTickProto.GetCount())
	}

	d.wrapper.RealtimeBar(reqID, t, open, high, low, close, volume, wap, count)
}

func (d *EDecoder) processFundamentalDataMsg(msgBuf *MsgBuffer) {

	msgBuf.decode() // version

	reqID := msgBuf.decodeInt64()
	data := msgBuf.decodeString()

	d.wrapper.FundamentalData(reqID, data)
}

func (d *EDecoder) processContractDataEndMsg(msgBuf *MsgBuffer) {

	msgBuf.decode() // version

	reqID := msgBuf.decodeInt64()

	d.wrapper.ContractDetailsEnd(reqID)
}

func (d *EDecoder) processFundamentalsDataMsgProtoBuf(msgBuf *MsgBuffer) {
	var fundamentalsDataProto protobuf.FundamentalsData
	if err := proto.Unmarshal(msgBuf.bs, &fundamentalsDataProto); err != nil {
		log.Error().Err(err).Msg("failed to unmarshal FundamentalsData message")
		return
	}

	d.wrapper.FundamentalsDataProtoBuf(&fundamentalsDataProto)

	reqID := NO_VALID_ID
	if fundamentalsDataProto.ReqId != nil {
		reqID = int64(fundamentalsDataProto.GetReqId())
	}
	data := ""
	if fundamentalsDataProto.Data != nil {
		data = fundamentalsDataProto.GetData()
	}

	d.wrapper.FundamentalData(reqID, data)
}

func (d *EDecoder) processContractDataEndMsgProtoBuf(msgBuf *MsgBuffer) {
	var contractDataEndProto protobuf.ContractDataEnd
	err := proto.Unmarshal(msgBuf.bs, &contractDataEndProto)
	if err != nil {
		log.Error().Err(err).Msg("failed to unmarshal ContractDataEnd message")
		return
	}

	d.wrapper.ContractDataEndProtoBuf(&contractDataEndProto)

	reqID := NO_VALID_ID
	if contractDataEndProto.ReqId != nil {
		reqID = int64(contractDataEndProto.GetReqId())
	}

	d.wrapper.ContractDetailsEnd(reqID)
}

func (d *EDecoder) processOpenOrderEndMsg(msgBuf *MsgBuffer) {

	msgBuf.decode() // version

	d.wrapper.OpenOrderEnd()
}

func (d *EDecoder) processOpenOrderEndMsgProtoBuf(msgBuf *MsgBuffer) {

	var openOrdersEndProto protobuf.OpenOrdersEnd
	err := proto.Unmarshal(msgBuf.Bytes(), &openOrdersEndProto)
	if err != nil {
		log.Panic().Err(err).Msg("processOpenOrderEndMsgProtoBuf unmarshal error")
	}

	d.wrapper.OpenOrdersEndProtoBuf(&openOrdersEndProto)

	d.wrapper.OpenOrderEnd()
}

func (d *EDecoder) processAcctDownloadEndMsg(msgBuf *MsgBuffer) {

	msgBuf.decode() // version

	accountName := msgBuf.decodeString()

	d.wrapper.AccountDownloadEnd(accountName)
}

func (d *EDecoder) processExecutionDetailsEndMsg(msgBuf *MsgBuffer) {

	msgBuf.decode() // version

	reqID := msgBuf.decodeInt64()

	d.wrapper.ExecDetailsEnd(reqID)
}

func (d *EDecoder) processAccountDataEndMsgProtoBuf(msgBuf *MsgBuffer) {

	var accountDataEndProto protobuf.AccountDataEnd
	err := proto.Unmarshal(msgBuf.bs, &accountDataEndProto)
	if err != nil {
		log.Error().Err(err).Msg("failed to unmarshal AccountDataEnd message")
		return
	}

	d.wrapper.AccountDataEndProtoBuf(&accountDataEndProto)

	var accountName string
	if accountDataEndProto.AccountName != nil {
		accountName = accountDataEndProto.GetAccountName()
	}

	d.wrapper.AccountDownloadEnd(accountName)
}

func (d *EDecoder) processExecutionDetailsEndMsgProtoBuf(msgBuf *MsgBuffer) {

	var executionDetailsEndProto protobuf.ExecutionDetailsEnd
	err := proto.Unmarshal(msgBuf.Bytes(), &executionDetailsEndProto)
	if err != nil {
		log.Panic().Err(err).Msg("processExecutionDetailsEndMsgProtoBuf unmarshal error")
	}

	d.wrapper.ExecDetailsEndProtoBuf(&executionDetailsEndProto)

	reqID := NO_VALID_ID
	if executionDetailsEndProto.ReqId != nil {
		reqID = int64(executionDetailsEndProto.GetReqId())
	}

	d.wrapper.ExecDetailsEnd(reqID)
}

func (d *EDecoder) processDeltaNeutralValidationMsg(msgBuf *MsgBuffer) {

	msgBuf.decode() // version

	reqID := msgBuf.decodeInt64()

	deltaNeutralContract := NewDeltaNeutralContract()

	deltaNeutralContract.ConID = msgBuf.decodeInt64()
	deltaNeutralContract.Delta = msgBuf.decodeFloat64()
	deltaNeutralContract.Price = msgBuf.decodeFloat64()

	d.wrapper.DeltaNeutralValidation(reqID, deltaNeutralContract)
}

func (d *EDecoder) processTickSnapshotEndMsg(msgBuf *MsgBuffer) {

	msgBuf.decode() // version

	reqID := msgBuf.decodeInt64()

	d.wrapper.TickSnapshotEnd(reqID)
}

func (d *EDecoder) processTickSnapshotEndMsgProtoBuf(msgBuf *MsgBuffer) {
	var tickSnapshotEndProto protobuf.TickSnapshotEnd
	err := proto.Unmarshal(msgBuf.bs, &tickSnapshotEndProto)
	if err != nil {
		log.Error().Err(err).Msg("failed to unmarshal TickSnapshotEnd message")
		return
	}

	d.wrapper.TickSnapshotEndProtoBuf(&tickSnapshotEndProto)

	reqID := NO_VALID_ID
	if tickSnapshotEndProto.ReqId != nil {
		reqID = int64(tickSnapshotEndProto.GetReqId())
	}

	d.wrapper.TickSnapshotEnd(reqID)
}

func (d *EDecoder) processMarketDataTypeMsg(msgBuf *MsgBuffer) {

	msgBuf.decode() // version

	reqID := msgBuf.decodeInt64()
	marketDataType := msgBuf.decodeInt64()

	d.wrapper.MarketDataType(reqID, marketDataType)
}

func (d *EDecoder) processMarketDataTypeMsgProtoBuf(msgBuf *MsgBuffer) {
	var marketDataTypeProto protobuf.MarketDataType
	err := proto.Unmarshal(msgBuf.bs, &marketDataTypeProto)
	if err != nil {
		log.Error().Err(err).Msg("failed to unmarshal MarketDataType message")
		return
	}

	d.wrapper.MarketDataTypeProtoBuf(&marketDataTypeProto)

	reqID := NO_VALID_ID
	if marketDataTypeProto.ReqId != nil {
		reqID = int64(marketDataTypeProto.GetReqId())
	}

	var marketDataType int64
	if marketDataTypeProto.MarketDataType != nil {
		marketDataType = int64(marketDataTypeProto.GetMarketDataType())
	}

	d.wrapper.MarketDataType(reqID, marketDataType)
}

func (d *EDecoder) processCommissionAndFeesReportMsg(msgBuf *MsgBuffer) {

	msgBuf.decode() // version

	commissionAndFeesReport := NewCommissionAndFeesReport()
	commissionAndFeesReport.ExecID = msgBuf.decodeString()
	commissionAndFeesReport.CommissionAndFees = msgBuf.decodeFloat64()
	commissionAndFeesReport.Currency = msgBuf.decodeString()
	commissionAndFeesReport.RealizedPNL = msgBuf.decodeFloat64()
	commissionAndFeesReport.Yield = msgBuf.decodeFloat64()
	commissionAndFeesReport.YieldRedemptionDate = msgBuf.decodeInt64()

	d.wrapper.CommissionAndFeesReport(commissionAndFeesReport)
}

func (d *EDecoder) processCommissionAndFeesReportMsgProtoBuf(msgBuf *MsgBuffer) {
	var commissionAndFeesReportProto protobuf.CommissionAndFeesReport
	err := proto.Unmarshal(msgBuf.bs, &commissionAndFeesReportProto)
	if err != nil {
		log.Error().Err(err).Msg("failed to unmarshal CommissionAndFeesReport message")
		return
	}

	d.wrapper.CommissionAndFeesReportProtoBuf(&commissionAndFeesReportProto)

	// mirror C++ field-by-field defaults
	report := CommissionAndFeesReport{
		ExecID:              "",
		CommissionAndFees:   0,
		Currency:            "",
		RealizedPNL:         0,
		Yield:               0,
		YieldRedemptionDate: 0,
	}
	if commissionAndFeesReportProto.ExecId != nil {
		report.ExecID = commissionAndFeesReportProto.GetExecId()
	}
	if commissionAndFeesReportProto.CommissionAndFees != nil {
		report.CommissionAndFees = commissionAndFeesReportProto.GetCommissionAndFees()
	}
	if commissionAndFeesReportProto.Currency != nil {
		report.Currency = commissionAndFeesReportProto.GetCurrency()
	}
	if commissionAndFeesReportProto.RealizedPNL != nil {
		report.RealizedPNL = commissionAndFeesReportProto.GetRealizedPNL()
	}
	if commissionAndFeesReportProto.BondYield != nil {
		report.Yield = commissionAndFeesReportProto.GetBondYield()
	}
	if commissionAndFeesReportProto.YieldRedemptionDate != nil {
		yrd, err := strconv.ParseInt(commissionAndFeesReportProto.GetYieldRedemptionDate(), 10, 64)
		if err != nil {
			log.Error().Err(err).Msg("failed to parse CommissionAndFeesReport YieldRedemptionDate")
			return
		}
		report.YieldRedemptionDate = yrd
	}

	d.wrapper.CommissionAndFeesReport(report)
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

func (d *EDecoder) processPositionMsgProtoBuf(msgBuf *MsgBuffer) {

	var positionProto protobuf.Position
	err := proto.Unmarshal(msgBuf.bs, &positionProto)
	if err != nil {
		log.Error().Err(err).Msg("failed to unmarshal Position message")
		return
	}

	d.wrapper.PositionProtoBuf(&positionProto)

	if positionProto.Contract == nil {
		return
	}

	contract := decodeContract(positionProto.GetContract())

	var position Decimal
	if positionProto.Position != nil {
		position = StringToDecimal(positionProto.GetPosition())
	} else {
		position = UNSET_DECIMAL
	}

	var avgCost float64
	if positionProto.AvgCost != nil {
		avgCost = positionProto.GetAvgCost()
	}

	var account string
	if positionProto.Account != nil {
		account = positionProto.GetAccount()
	}

	d.wrapper.Position(account, contract, position, avgCost)
}

func (d *EDecoder) processPositionEndMsg(*MsgBuffer) {

	d.wrapper.PositionEnd()
}

func (d *EDecoder) processPositionEndMsgProtoBuf(msgBuf *MsgBuffer) {
	var positionEndProto protobuf.PositionEnd
	err := proto.Unmarshal(msgBuf.bs, &positionEndProto)
	if err != nil {
		log.Error().Err(err).Msg("failed to unmarshal PositionEnd message")
		return
	}

	d.wrapper.PositionEndProtoBuf(&positionEndProto)

	d.wrapper.PositionEnd()
}

func (d *EDecoder) processAccountSummaryMsg(msgBuf *MsgBuffer) {

	msgBuf.decode() // version

	reqID := msgBuf.decodeInt64()
	account := msgBuf.decodeString()
	tag := msgBuf.decodeString()
	value := msgBuf.decodeString()
	currency := msgBuf.decodeString()

	d.wrapper.AccountSummary(reqID, account, tag, value, currency)
}

func (d *EDecoder) processAccountSummaryMsgProtoBuf(msgBuf *MsgBuffer) {
	var accountSummaryProto protobuf.AccountSummary
	err := proto.Unmarshal(msgBuf.bs, &accountSummaryProto)
	if err != nil {
		log.Error().Err(err).Msg("failed to unmarshal AccountSummary message")
		return
	}

	d.wrapper.AccountSummaryProtoBuf(&accountSummaryProto)

	var reqId int64
	if accountSummaryProto.ReqId != nil {
		reqId = int64(accountSummaryProto.GetReqId())
	} else {
		reqId = NO_VALID_ID
	}

	var account string
	if accountSummaryProto.Account != nil {
		account = accountSummaryProto.GetAccount()
	}

	var tag string
	if accountSummaryProto.Tag != nil {
		tag = accountSummaryProto.GetTag()
	}

	var value string
	if accountSummaryProto.Value != nil {
		value = accountSummaryProto.GetValue()
	}

	var currency string
	if accountSummaryProto.Currency != nil {
		currency = accountSummaryProto.GetCurrency()
	}

	d.wrapper.AccountSummary(reqId, account, tag, value, currency)
}

func (d *EDecoder) processAccountSummaryEndMsg(msgBuf *MsgBuffer) {

	msgBuf.decode() // version

	reqID := msgBuf.decodeInt64()

	d.wrapper.AccountSummaryEnd(reqID)
}

func (d *EDecoder) processAccountSummaryEndMsgProtoBuf(msgBuf *MsgBuffer) {
	var accountSummaryEndProto protobuf.AccountSummaryEnd
	err := proto.Unmarshal(msgBuf.bs, &accountSummaryEndProto)
	if err != nil {
		log.Error().Err(err).Msg("failed to unmarshal AccountSummaryEnd message")
		return
	}

	d.wrapper.AccountSummaryEndProtoBuf(&accountSummaryEndProto)

	var reqId int64
	if accountSummaryEndProto.ReqId != nil {
		reqId = int64(accountSummaryEndProto.GetReqId())
	} else {
		reqId = NO_VALID_ID
	}

	d.wrapper.AccountSummaryEnd(reqId)
}

func (d *EDecoder) processVerifyMessageApiMsg(msgBuf *MsgBuffer) {

	msgBuf.decode() // version

	apiData := msgBuf.decodeString()

	d.wrapper.VerifyMessageAPI(apiData)
}

func (d *EDecoder) processVerifyMessageApiMsgProtoBuf(msgBuf *MsgBuffer) {

	var protoMsg protobuf.VerifyMessageApi
	if err := proto.Unmarshal(msgBuf.bs, &protoMsg); err != nil {
		log.Error().Err(err).Msg("failed to unmarshal VerifyMessageApi")
		return
	}

	d.wrapper.VerifyMessageApiProtoBuf(&protoMsg)

	data := ""
	if protoMsg.ApiData != nil {
		data = protoMsg.GetApiData()
	}

	d.wrapper.VerifyMessageAPI(data)
}

func (d *EDecoder) processVerifyCompletedMsg(msgBuf *MsgBuffer) {

	msgBuf.decode() // version

	isSuccessful := msgBuf.decodeBool()
	errorText := msgBuf.decodeString()

	d.wrapper.VerifyCompleted(isSuccessful, errorText)
}

func (d *EDecoder) processVerifyCompletedMsgProtoBuf(msgBuf *MsgBuffer) {

	var protoMsg protobuf.VerifyCompleted
	if err := proto.Unmarshal(msgBuf.bs, &protoMsg); err != nil {
		log.Error().Err(err).Msg("failed to unmarshal VerifyCompleted")
		return
	}

	d.wrapper.VerifyCompletedProtoBuf(&protoMsg)

	ok := false
	if protoMsg.IsSuccessful != nil {
		ok = protoMsg.GetIsSuccessful()
	}
	errText := ""
	if protoMsg.ErrorText != nil {
		errText = protoMsg.GetErrorText()
	}
	d.wrapper.VerifyCompleted(ok, errText)
}

func (d *EDecoder) processDisplayGroupListMsg(msgBuf *MsgBuffer) {

	msgBuf.decode() // version

	reqID := msgBuf.decodeInt64()
	groups := msgBuf.decodeString()

	d.wrapper.DisplayGroupList(reqID, groups)
}

func (d *EDecoder) processDisplayGroupListMsgProtoBuf(msgBuf *MsgBuffer) {

	var protoMsg protobuf.DisplayGroupList
	if err := proto.Unmarshal(msgBuf.bs, &protoMsg); err != nil {
		log.Error().Err(err).Msg("failed to unmarshal DisplayGroupList")
		return
	}

	d.wrapper.DisplayGroupListProtoBuf(&protoMsg)

	reqID := NO_VALID_ID
	if protoMsg.ReqId != nil {
		reqID = int64(protoMsg.GetReqId())
	}
	groups := ""
	if protoMsg.Groups != nil {
		groups = protoMsg.GetGroups()
	}

	d.wrapper.DisplayGroupList(reqID, groups)
}

func (d *EDecoder) processDisplayGroupUpdatedMsg(msgBuf *MsgBuffer) {

	msgBuf.decode() // version

	reqID := msgBuf.decodeInt64()
	contractInfo := msgBuf.decodeString()

	d.wrapper.DisplayGroupUpdated(reqID, contractInfo)
}

func (d *EDecoder) processDisplayGroupUpdatedMsgProtoBuf(msgBuf *MsgBuffer) {

	var protoMsg protobuf.DisplayGroupUpdated
	if err := proto.Unmarshal(msgBuf.bs, &protoMsg); err != nil {
		log.Error().Err(err).Msg("failed to unmarshal DisplayGroupUpdated")
		return
	}

	d.wrapper.DisplayGroupUpdatedProtoBuf(&protoMsg)

	reqID := NO_VALID_ID
	if protoMsg.ReqId != nil {
		reqID = int64(protoMsg.GetReqId())
	}
	contractInfo := ""
	if protoMsg.ContractInfo != nil {
		contractInfo = protoMsg.GetContractInfo()
	}

	d.wrapper.DisplayGroupUpdated(reqID, contractInfo)
}

func (d *EDecoder) processVerifyAndAuthMessageApiMsg(msgBuf *MsgBuffer) {

	msgBuf.decode() // version

	apiData := msgBuf.decodeString()
	xyzChallange := msgBuf.decodeString()

	d.wrapper.VerifyAndAuthMessageAPI(apiData, xyzChallange)
}

func (d *EDecoder) processVerifyAndAuthCompletedMsg(msgBuf *MsgBuffer) {

	msgBuf.decode() // version

	isSuccessful := msgBuf.decodeBool()
	errorText := msgBuf.decodeString()

	d.wrapper.VerifyAndAuthCompleted(isSuccessful, errorText)
}

func (d *EDecoder) processPositionMultiMsg(msgBuf *MsgBuffer) {

	msgBuf.decode() // version

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

func (d *EDecoder) processPositionMultiMsgProtoBuf(msgBuf *MsgBuffer) {
	var positionMultiProto protobuf.PositionMulti
	err := proto.Unmarshal(msgBuf.bs, &positionMultiProto)
	if err != nil {
		log.Error().Err(err).Msg("failed to unmarshal PositionMulti message")
		return
	}

	d.wrapper.PositionMultiProtoBuf(&positionMultiProto)

	var reqId int64
	if positionMultiProto.ReqId != nil {
		reqId = int64(positionMultiProto.GetReqId())
	} else {
		reqId = NO_VALID_ID
	}

	var account string
	if positionMultiProto.Account != nil {
		account = positionMultiProto.GetAccount()
	}

	var modelCode string
	if positionMultiProto.ModelCode != nil {
		modelCode = positionMultiProto.GetModelCode()
	}

	if positionMultiProto.Contract == nil {
		return
	}

	contract := decodeContract(positionMultiProto.GetContract())

	var position Decimal
	if positionMultiProto.Position != nil {
		position = StringToDecimal(positionMultiProto.GetPosition())
	} else {
		position = UNSET_DECIMAL
	}

	var avgCost float64
	if positionMultiProto.AvgCost != nil {
		avgCost = positionMultiProto.GetAvgCost()
	}

	d.wrapper.PositionMulti(reqId, account, modelCode, contract, position, avgCost)
}

func (d *EDecoder) processPositionMultiEndMsg(msgBuf *MsgBuffer) {

	msgBuf.decode() // version

	reqID := msgBuf.decodeInt64()

	d.wrapper.PositionMultiEnd(reqID)
}

func (d *EDecoder) processPositionMultiEndMsgProtoBuf(msgBuf *MsgBuffer) {
	var positionMultiEndProto protobuf.PositionMultiEnd
	err := proto.Unmarshal(msgBuf.bs, &positionMultiEndProto)
	if err != nil {
		log.Error().Err(err).Msg("failed to unmarshal PositionMultiEnd message")
		return
	}

	d.wrapper.PositionMultiEndProtoBuf(&positionMultiEndProto)

	var reqId int64
	if positionMultiEndProto.ReqId != nil {
		reqId = int64(positionMultiEndProto.GetReqId())
	} else {
		reqId = NO_VALID_ID
	}

	d.wrapper.PositionMultiEnd(reqId)
}

func (d *EDecoder) processAccountUpdateMultiMsg(msgBuf *MsgBuffer) {

	msgBuf.decode() // version

	reqID := msgBuf.decodeInt64()
	account := msgBuf.decodeString()
	modelCode := msgBuf.decodeString()
	key := msgBuf.decodeString()
	value := msgBuf.decodeString()
	currency := msgBuf.decodeString()

	d.wrapper.AccountUpdateMulti(reqID, account, modelCode, key, value, currency)
}

func (d *EDecoder) processAccountUpdateMultiMsgProtoBuf(msgBuf *MsgBuffer) {
	var accountUpdateMultiProto protobuf.AccountUpdateMulti
	err := proto.Unmarshal(msgBuf.bs, &accountUpdateMultiProto)
	if err != nil {
		log.Error().Err(err).Msg("failed to unmarshal AccountUpdateMulti message")
		return
	}

	d.wrapper.AccountUpdateMultiProtoBuf(&accountUpdateMultiProto)

	var reqId int64
	if accountUpdateMultiProto.ReqId != nil {
		reqId = int64(accountUpdateMultiProto.GetReqId())
	} else {
		reqId = NO_VALID_ID
	}

	var account string
	if accountUpdateMultiProto.Account != nil {
		account = accountUpdateMultiProto.GetAccount()
	}

	var modelCode string
	if accountUpdateMultiProto.ModelCode != nil {
		modelCode = accountUpdateMultiProto.GetModelCode()
	}

	var key string
	if accountUpdateMultiProto.Key != nil {
		key = accountUpdateMultiProto.GetKey()
	}

	var value string
	if accountUpdateMultiProto.Value != nil {
		value = accountUpdateMultiProto.GetValue()
	}

	var currency string
	if accountUpdateMultiProto.Currency != nil {
		currency = accountUpdateMultiProto.GetCurrency()
	}

	d.wrapper.AccountUpdateMulti(reqId, account, modelCode, key, value, currency)
}

func (d *EDecoder) processAccountUpdateMultiEndMsg(msgBuf *MsgBuffer) {

	msgBuf.decode() // version

	reqID := msgBuf.decodeInt64()

	d.wrapper.AccountUpdateMultiEnd(reqID)
}

func (d *EDecoder) processAccountUpdateMultiEndMsgProtoBuf(msgBuf *MsgBuffer) {
	var accountUpdateMultiEndProto protobuf.AccountUpdateMultiEnd
	err := proto.Unmarshal(msgBuf.bs, &accountUpdateMultiEndProto)
	if err != nil {
		log.Error().Err(err).Msg("failed to unmarshal AccountUpdateMultiEnd message")
		return
	}

	d.wrapper.AccountUpdateMultiEndProtoBuf(&accountUpdateMultiEndProto)

	var reqId int64
	if accountUpdateMultiEndProto.ReqId != nil {
		reqId = int64(accountUpdateMultiEndProto.GetReqId())
	} else {
		reqId = NO_VALID_ID
	}

	d.wrapper.AccountUpdateMultiEnd(reqId)
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

func (d *EDecoder) processSecurityDefinitionOptionParameterMsgProtoBuf(msgBuf *MsgBuffer) {
	var protoMsg protobuf.SecDefOptParameter
	if err := proto.Unmarshal(msgBuf.bs, &protoMsg); err != nil {
		log.Error().Err(err).Msg("failed to unmarshal SecurityDefinitionOptionParameter")
		return
	}
	d.wrapper.SecDefOptParameterProtoBuf(&protoMsg)

	reqID := NO_VALID_ID
	if protoMsg.ReqId != nil {
		reqID = int64(protoMsg.GetReqId())
	}
	exchange := ""
	if protoMsg.Exchange != nil {
		exchange = protoMsg.GetExchange()
	}
	underlyingConID := int64(0)
	if protoMsg.UnderlyingConId != nil {
		underlyingConID = int64(protoMsg.GetUnderlyingConId())
	}
	tradingClass := ""
	if protoMsg.TradingClass != nil {
		tradingClass = protoMsg.GetTradingClass()
	}
	multiplier := ""
	if protoMsg.Multiplier != nil {
		multiplier = protoMsg.GetMultiplier()
	}
	expirations := make([]string, 0, len(protoMsg.GetExpirations()))
	expirations = append(expirations, protoMsg.GetExpirations()...)

	strikes := make([]float64, 0, len(protoMsg.GetStrikes()))
	strikes = append(strikes, protoMsg.GetStrikes()...)

	d.wrapper.SecurityDefinitionOptionParameter(
		reqID, exchange, underlyingConID, tradingClass, multiplier,
		expirations, strikes,
	)
}

func (d *EDecoder) processSecurityDefinitionOptionalParameterEndMsg(msgBuf *MsgBuffer) {

	reqID := msgBuf.decodeInt64()

	d.wrapper.SecurityDefinitionOptionParameterEnd(reqID)
}

func (d *EDecoder) processSecurityDefinitionOptionParameterEndMsgProtoBuf(msgBuf *MsgBuffer) {
	var protoMsg protobuf.SecDefOptParameterEnd
	if err := proto.Unmarshal(msgBuf.bs, &protoMsg); err != nil {
		log.Error().Err(err).Msg("failed to unmarshal SecurityDefinitionOptionParameterEnd")
		return
	}
	d.wrapper.SecDefOptParameterEndProtoBuf(&protoMsg)

	reqID := NO_VALID_ID
	if protoMsg.ReqId != nil {
		reqID = int64(protoMsg.GetReqId())
	}
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

func (d *EDecoder) processSoftDollarTiersMsgProtoBuf(msgBuf *MsgBuffer) {
	var protoMsg protobuf.SoftDollarTiers
	if err := proto.Unmarshal(msgBuf.bs, &protoMsg); err != nil {
		log.Error().Err(err).Msg("failed to unmarshal SoftDollarTierList")
		return
	}

	d.wrapper.SoftDollarTiersProtoBuf(&protoMsg)

	reqID := NO_VALID_ID
	if protoMsg.ReqId != nil {
		reqID = int64(protoMsg.GetReqId())
	}
	tiers := make([]SoftDollarTier, 0, len(protoMsg.GetSoftDollarTiers()))
	for _, t := range protoMsg.GetSoftDollarTiers() {
		tiers = append(tiers, decodeSoftDollarTier(t))
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

func (d *EDecoder) processFamilyCodesMsgProtoBuf(msgBuf *MsgBuffer) {
	var protoMsg protobuf.FamilyCodes
	if err := proto.Unmarshal(msgBuf.bs, &protoMsg); err != nil {
		log.Error().Err(err).Msg("failed to unmarshal FamilyCodeList")
		return
	}
	d.wrapper.FamilyCodesProtoBuf(&protoMsg)

	codes := make([]FamilyCode, 0, len(protoMsg.GetFamilyCodes()))
	for _, fc := range protoMsg.GetFamilyCodes() {
		codes = append(codes, *decodeFamilyCode(fc))
	}
	d.wrapper.FamilyCodes(codes)
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

func (d *EDecoder) processSymbolSamplesMsgProtoBuf(msgBuf *MsgBuffer) {
	var protoMsg protobuf.SymbolSamples
	if err := proto.Unmarshal(msgBuf.bs, &protoMsg); err != nil {
		log.Error().Err(err).Msg("failed to unmarshal SymbolSamples")
		return
	}
	d.wrapper.SymbolSamplesProtoBuf(&protoMsg)

	reqID := NO_VALID_ID
	if protoMsg.ReqId != nil {
		reqID = int64(protoMsg.GetReqId())
	}
	descs := make([]ContractDescription, 0, len(protoMsg.GetContractDescriptions()))
	for _, cd := range protoMsg.GetContractDescriptions() {
		contract := Contract{}
		if cd.Contract != nil {
			contract = *decodeContract(cd.GetContract())
		}
		derivativeSecTypes := cd.GetDerivativeSecTypes()
		descs = append(descs, ContractDescription{
			Contract:           contract,
			DerivativeSecTypes: append([]string{}, derivativeSecTypes...),
		})
	}
	d.wrapper.SymbolSamples(reqID, descs)
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

func (d *EDecoder) processMktDepthExchangesMsgProtoBuf(msgBuf *MsgBuffer) {

	var protoMsg protobuf.MarketDepthExchanges
	if err := proto.Unmarshal(msgBuf.bs, &protoMsg); err != nil {
		log.Error().Err(err).Msg("failed to unmarshal MarketDepthExchanges")
		return
	}

	d.wrapper.MarketDepthExchangesProtoBuf(&protoMsg)

	descs := make([]DepthMktDataDescription, 0, len(protoMsg.GetDepthMarketDataDescriptions()))
	for _, dd := range protoMsg.GetDepthMarketDataDescriptions() {
		descs = append(descs, *decodeDepthMarketDataDescription(dd))
	}
	d.wrapper.MktDepthExchanges(descs)
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

func (d *EDecoder) processTickNewsMsgProtoBuf(msgBuf *MsgBuffer) {
	var tickNewsProto protobuf.TickNews
	if err := proto.Unmarshal(msgBuf.bs, &tickNewsProto); err != nil {
		log.Error().Err(err).Msg("failed to unmarshal TickNews message")
		return
	}
	d.wrapper.TickNewsProtoBuf(&tickNewsProto)

	reqID := NO_VALID_ID
	if tickNewsProto.ReqId != nil {
		reqID = int64(tickNewsProto.GetReqId())
	}
	timestamp := int64(0)
	if tickNewsProto.Timestamp != nil {
		timestamp = tickNewsProto.GetTimestamp()
	}
	providerCode := ""
	if tickNewsProto.ProviderCode != nil {
		providerCode = tickNewsProto.GetProviderCode()
	}
	articleID := ""
	if tickNewsProto.ArticleId != nil {
		articleID = tickNewsProto.GetArticleId()
	}
	headline := ""
	if tickNewsProto.Headline != nil {
		headline = tickNewsProto.GetHeadline()
	}
	extraData := ""
	if tickNewsProto.ExtraData != nil {
		extraData = tickNewsProto.GetExtraData()
	}

	d.wrapper.TickNews(reqID, timestamp, providerCode, articleID, headline, extraData)
}

func (d *EDecoder) processTickReqParamsMsg(msgBuf *MsgBuffer) {

	tickerID := msgBuf.decodeInt64()

	minTick := msgBuf.decodeFloat64()
	bboExchange := msgBuf.decodeString()
	snapshotPermissions := msgBuf.decodeInt64()

	d.wrapper.TickReqParams(tickerID, minTick, bboExchange, snapshotPermissions)
}

func (d *EDecoder) processTickReqParamsMsgProtoBuf(msgBuf *MsgBuffer) {
	var tickReqParamsProto protobuf.TickReqParams
	err := proto.Unmarshal(msgBuf.bs, &tickReqParamsProto)
	if err != nil {
		log.Error().Err(err).Msg("failed to unmarshal TickReqParams message")
		return
	}

	d.wrapper.TickReqParamsProtoBuf(&tickReqParamsProto) // Assuming you have a wrapper method for the raw protobuf message

	reqID := NO_VALID_ID
	if tickReqParamsProto.ReqId != nil {
		reqID = int64(tickReqParamsProto.GetReqId())
	}

	minTick := UNSET_FLOAT
	if tickReqParamsProto.MinTick != nil {
		minTick, err = strconv.ParseFloat(tickReqParamsProto.GetMinTick(), 64)
		if err != nil {
			log.Error().Err(err).Msg("failed to convert TickReqParams minTick")
			minTick = UNSET_FLOAT
		}
	}

	var bboExchange string
	if tickReqParamsProto.BboExchange != nil {
		bboExchange = tickReqParamsProto.GetBboExchange()
	}

	snapshotPermissions := UNSET_INT
	if tickReqParamsProto.SnapshotPermissions != nil {
		snapshotPermissions = int64(tickReqParamsProto.GetSnapshotPermissions())
	}

	d.wrapper.TickReqParams(reqID, minTick, bboExchange, snapshotPermissions)
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

func (d *EDecoder) processSmartComponentsMsgProtoBuf(msgBuf *MsgBuffer) {
	var protoMsg protobuf.SmartComponents
	if err := proto.Unmarshal(msgBuf.bs, &protoMsg); err != nil {
		log.Error().Err(err).Msg("failed to unmarshal SmartComponents")
		return
	}
	d.wrapper.SmartComponentsProtoBuf(&protoMsg)

	reqID := NO_VALID_ID
	if protoMsg.ReqId != nil {
		reqID = int64(protoMsg.GetReqId())
	}
	comps := decodeSmartComponents(&protoMsg)

	d.wrapper.SmartComponents(reqID, comps)
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

func (d *EDecoder) processNewsProvidersMsgProtoBuf(msgBuf *MsgBuffer) {
	var newsProvidersProto protobuf.NewsProviders
	if err := proto.Unmarshal(msgBuf.bs, &newsProvidersProto); err != nil {
		log.Error().Err(err).Msg("failed to unmarshal NewsProviders message")
		return
	}
	d.wrapper.NewsProvidersProtoBuf(&newsProvidersProto)

	list := make([]NewsProvider, 0, len(newsProvidersProto.GetNewsProviders()))
	for _, np := range newsProvidersProto.GetNewsProviders() {
		provider := NewNewsProvider()
		provider.Code = np.GetProviderCode()
		provider.Name = np.GetProviderName()
		list = append(list, provider)
	}
	d.wrapper.NewsProviders(list)
}

func (d *EDecoder) processNewsArticleMsg(msgBuf *MsgBuffer) {

	reqID := msgBuf.decodeInt64()

	articleType := msgBuf.decodeInt64()
	articleText := msgBuf.decodeString()

	d.wrapper.NewsArticle(reqID, articleType, articleText)
}

func (d *EDecoder) processNewsArticleMsgProtoBuf(msgBuf *MsgBuffer) {
	var newsArticleProto protobuf.NewsArticle
	if err := proto.Unmarshal(msgBuf.bs, &newsArticleProto); err != nil {
		log.Error().Err(err).Msg("failed to unmarshal NewsArticle message")
		return
	}
	d.wrapper.NewsArticleProtoBuf(&newsArticleProto)

	reqID := NO_VALID_ID
	if newsArticleProto.ReqId != nil {
		reqID = int64(newsArticleProto.GetReqId())
	}
	articleType := int64(0)
	if newsArticleProto.ArticleType != nil {
		articleType = int64(newsArticleProto.GetArticleType())
	}
	articleText := ""
	if newsArticleProto.ArticleText != nil {
		articleText = newsArticleProto.GetArticleText()
	}

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

func (d *EDecoder) processHistoricalNewsMsgProtoBuf(msgBuf *MsgBuffer) {
	var historicalNewsProto protobuf.HistoricalNews
	if err := proto.Unmarshal(msgBuf.bs, &historicalNewsProto); err != nil {
		log.Error().Err(err).Msg("failed to unmarshal HistoricalNews message")
		return
	}
	d.wrapper.HistoricalNewsProtoBuf(&historicalNewsProto)

	reqID := NO_VALID_ID
	if historicalNewsProto.ReqId != nil {
		reqID = int64(historicalNewsProto.GetReqId())
	}
	timeStr := ""
	if historicalNewsProto.Time != nil {
		timeStr = historicalNewsProto.GetTime()
	}
	providerCode := ""
	if historicalNewsProto.ProviderCode != nil {
		providerCode = historicalNewsProto.GetProviderCode()
	}
	articleID := ""
	if historicalNewsProto.ArticleId != nil {
		articleID = historicalNewsProto.GetArticleId()
	}
	headline := ""
	if historicalNewsProto.Headline != nil {
		headline = historicalNewsProto.GetHeadline()
	}

	d.wrapper.HistoricalNews(reqID, timeStr, providerCode, articleID, headline)
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

func (d *EDecoder) processHistoricalNewsEndMsgProtoBuf(msgBuf *MsgBuffer) {
	var historicalNewsEndProto protobuf.HistoricalNewsEnd
	if err := proto.Unmarshal(msgBuf.bs, &historicalNewsEndProto); err != nil {
		log.Error().Err(err).Msg("failed to unmarshal HistoricalNewsEnd message")
		return
	}
	d.wrapper.HistoricalNewsEndProtoBuf(&historicalNewsEndProto)

	reqID := NO_VALID_ID
	if historicalNewsEndProto.ReqId != nil {
		reqID = int64(historicalNewsEndProto.GetReqId())
	}
	hasMore := false
	if historicalNewsEndProto.HasMore != nil {
		hasMore = historicalNewsEndProto.GetHasMore()
	}

	d.wrapper.HistoricalNewsEnd(reqID, hasMore)
}

func (d *EDecoder) processHeadTimestampMsgProtoBuf(msgBuf *MsgBuffer) {
	var headTimestampProto protobuf.HeadTimestamp
	if err := proto.Unmarshal(msgBuf.bs, &headTimestampProto); err != nil {
		log.Error().Err(err).Msg("failed to unmarshal HeadTimestamp message")
		return
	}

	d.wrapper.HeadTimestampProtoBuf(&headTimestampProto)

	reqID := NO_VALID_ID
	if headTimestampProto.ReqId != nil {
		reqID = int64(headTimestampProto.GetReqId())
	}

	var timestamp string
	if headTimestampProto.HeadTimestamp != nil {
		timestamp = headTimestampProto.GetHeadTimestamp()
	}

	d.wrapper.HeadTimestamp(reqID, timestamp)
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

func (d *EDecoder) processHistogramDataMsgProtoBuf(msgBuf *MsgBuffer) {
	var histogramDataProto protobuf.HistogramData
	if err := proto.Unmarshal(msgBuf.bs, &histogramDataProto); err != nil {
		log.Error().Err(err).Msg("failed to unmarshal HistogramData message")
		return
	}

	d.wrapper.HistogramDataProtoBuf(&histogramDataProto)

	reqID := NO_VALID_ID
	if histogramDataProto.ReqId != nil {
		reqID = int64(histogramDataProto.GetReqId())
	}

	var histogramData []HistogramData
	for _, entryProto := range histogramDataProto.GetHistogramDataEntries() {
		histogramEntry := decodeHistogramDataEntry(entryProto)
		histogramData = append(histogramData, *histogramEntry)
	}

	d.wrapper.HistogramData(reqID, histogramData)
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

func (d *EDecoder) processHistoricalDataUpdateMsgProtoBuf(msgBuf *MsgBuffer) {
	var historicalDataUpdateProto protobuf.HistoricalDataUpdate
	if err := proto.Unmarshal(msgBuf.bs, &historicalDataUpdateProto); err != nil {
		log.Error().Err(err).Msg("failed to unmarshal HistoricalDataUpdate message")
		return
	}

	d.wrapper.HistoricalDataUpdateProtoBuf(&historicalDataUpdateProto)

	reqID := NO_VALID_ID
	if historicalDataUpdateProto.ReqId != nil {
		reqID = int64(historicalDataUpdateProto.GetReqId())
	}

	if historicalDataUpdateProto.GetHistoricalDataBar() != nil {
		bar := decodeHistoricalDataBar(historicalDataUpdateProto.GetHistoricalDataBar())
		d.wrapper.HistoricalDataUpdate(reqID, bar)
	}
}

func (d *EDecoder) processRerouteMktDataReqMsg(msgBuf *MsgBuffer) {

	reqID := msgBuf.decodeInt64()

	conID := msgBuf.decodeInt64()
	exchange := msgBuf.decodeString()

	d.wrapper.RerouteMktDataReq(reqID, conID, exchange)
}

func (d *EDecoder) processRerouteMktDataReqMsgProtoBuf(msgBuf *MsgBuffer) {
	var rerouteProto protobuf.RerouteMarketDataRequest
	err := proto.Unmarshal(msgBuf.bs, &rerouteProto)
	if err != nil {
		log.Error().Err(err).Msg("failed to unmarshal RerouteMarketDataRequest message")
		return
	}

	d.wrapper.RerouteMarketDataRequestProtoBuf(&rerouteProto)

	reqID := NO_VALID_ID
	if rerouteProto.ReqId != nil {
		reqID = int64(rerouteProto.GetReqId())
	}
	var conID int64
	if rerouteProto.ConId != nil {
		conID = int64(rerouteProto.GetConId())
	}
	var exchange string
	if rerouteProto.Exchange != nil {
		exchange = rerouteProto.GetExchange()
	}

	d.wrapper.RerouteMktDataReq(reqID, conID, exchange)
}

func (d *EDecoder) processRerouteMktDepthReqMsg(msgBuf *MsgBuffer) {

	reqID := msgBuf.decodeInt64()

	conID := msgBuf.decodeInt64()
	exchange := msgBuf.decodeString()

	d.wrapper.RerouteMktDepthReq(reqID, conID, exchange)
}

func (d *EDecoder) processRerouteMktDepthReqMsgProtoBuf(msgBuf *MsgBuffer) {
	var rerouteDepthProto protobuf.RerouteMarketDepthRequest
	err := proto.Unmarshal(msgBuf.bs, &rerouteDepthProto)
	if err != nil {
		log.Error().Err(err).Msg("failed to unmarshal RerouteMarketDepthRequest message")
		return
	}

	d.wrapper.RerouteMarketDepthRequestProtoBuf(&rerouteDepthProto)

	reqID := NO_VALID_ID
	if rerouteDepthProto.ReqId != nil {
		reqID = int64(rerouteDepthProto.GetReqId())
	}
	var conID int64
	if rerouteDepthProto.ConId != nil {
		conID = int64(rerouteDepthProto.GetConId())
	}
	var exchange string
	if rerouteDepthProto.Exchange != nil {
		exchange = rerouteDepthProto.GetExchange()
	}

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

func (d *EDecoder) processMarketRuleMsgProtoBuf(msgBuf *MsgBuffer) {
	var protoMsg protobuf.MarketRule
	if err := proto.Unmarshal(msgBuf.bs, &protoMsg); err != nil {
		log.Error().Err(err).Msg("failed to unmarshal MarketRule")
		return
	}
	d.wrapper.MarketRuleProtoBuf(&protoMsg)

	ruleID := int64(0)
	if protoMsg.MarketRuleId != nil {
		ruleID = int64(protoMsg.GetMarketRuleId())
	}
	imps := make([]PriceIncrement, 0, len(protoMsg.GetPriceIncrements()))
	for _, p := range protoMsg.GetPriceIncrements() {
		imps = append(imps, *decodePriceIncrement(p))
	}
	d.wrapper.MarketRule(ruleID, imps)
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

func (d *EDecoder) processPnLMsgProtoBuf(msgBuf *MsgBuffer) {
	var pnlProto protobuf.PnL
	if err := proto.Unmarshal(msgBuf.bs, &pnlProto); err != nil {
		log.Error().Err(err).Msg("failed to unmarshal PnL message")
		return
	}

	d.wrapper.PnLProtoBuf(&pnlProto)

	reqID := NO_VALID_ID
	if pnlProto.ReqId != nil {
		reqID = int64(pnlProto.GetReqId())
	}
	dailyPnL := 0.0
	if pnlProto.DailyPnL != nil {
		dailyPnL = pnlProto.GetDailyPnL()
	}
	unrealizedPnL := 0.0
	if pnlProto.UnrealizedPnL != nil {
		unrealizedPnL = pnlProto.GetUnrealizedPnL()
	}
	realizedPnL := 0.0
	if pnlProto.RealizedPnL != nil {
		realizedPnL = pnlProto.GetRealizedPnL()
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

func (d *EDecoder) processPnLSingleMsgProtoBuf(msgBuf *MsgBuffer) {
	var pnlSingleProto protobuf.PnLSingle
	if err := proto.Unmarshal(msgBuf.bs, &pnlSingleProto); err != nil {
		log.Error().Err(err).Msg("failed to unmarshal PnLSingle message")
		return
	}
	d.wrapper.PnLSingleProtoBuf(&pnlSingleProto)

	reqID := NO_VALID_ID
	if pnlSingleProto.ReqId != nil {
		reqID = int64(pnlSingleProto.GetReqId())
	}
	var pos Decimal
	if pnlSingleProto.Position != nil {
		pos = StringToDecimal(pnlSingleProto.GetPosition())
	}
	dailyPnL := 0.0
	if pnlSingleProto.DailyPnL != nil {
		dailyPnL = pnlSingleProto.GetDailyPnL()
	}
	unrealizedPnL := 0.0
	if pnlSingleProto.UnrealizedPnL != nil {
		unrealizedPnL = pnlSingleProto.GetUnrealizedPnL()
	}
	realizedPnL := 0.0
	if pnlSingleProto.RealizedPnL != nil {
		realizedPnL = pnlSingleProto.GetRealizedPnL()
	}
	value := 0.0
	if pnlSingleProto.Value != nil {
		value = pnlSingleProto.GetValue()
	}

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
		msgBuf.decode()
		historicalTick.Price = msgBuf.decodeFloat64()
		historicalTick.Size = msgBuf.decodeDecimal()
		ticks = append(ticks, historicalTick)
	}

	done := msgBuf.decodeBool()

	d.wrapper.HistoricalTicks(reqID, ticks, done)
}

func (d *EDecoder) processHistoricalTicksMsgProtoBuf(msgBuf *MsgBuffer) {
	var historicalTicksProto protobuf.HistoricalTicks
	if err := proto.Unmarshal(msgBuf.bs, &historicalTicksProto); err != nil {
		log.Error().Err(err).Msg("failed to unmarshal HistoricalTicks message")
		return
	}

	d.wrapper.HistoricalTicksProtoBuf(&historicalTicksProto)

	reqID := NO_VALID_ID
	if historicalTicksProto.ReqId != nil {
		reqID = int64(historicalTicksProto.GetReqId())
	}

	isDone := false
	if historicalTicksProto.IsDone != nil {
		isDone = historicalTicksProto.GetIsDone()
	}

	var historicalTicks []HistoricalTick
	for _, tickProto := range historicalTicksProto.GetHistoricalTicks() {
		tick := decodeHistoricalTick(tickProto)
		historicalTicks = append(historicalTicks, *tick)
	}

	d.wrapper.HistoricalTicks(reqID, historicalTicks, isDone)
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
		historicalTickBidAsk.TickAttribBidAsk = tickAttribBidAsk
		historicalTickBidAsk.PriceBid = msgBuf.decodeFloat64()
		historicalTickBidAsk.PriceAsk = msgBuf.decodeFloat64()
		historicalTickBidAsk.SizeBid = msgBuf.decodeDecimal()
		historicalTickBidAsk.SizeAsk = msgBuf.decodeDecimal()
		ticks = append(ticks, historicalTickBidAsk)
	}

	done := msgBuf.decodeBool()

	d.wrapper.HistoricalTicksBidAsk(reqID, ticks, done)
}

func (d *EDecoder) processHistoricalTicksBidAskMsgProtoBuf(msgBuf *MsgBuffer) {
	var protoBA protobuf.HistoricalTicksBidAsk
	if err := proto.Unmarshal(msgBuf.bs, &protoBA); err != nil {
		log.Error().Err(err).Msg("failed to unmarshal HistoricalTicksBidAsk message")
		return
	}

	d.wrapper.HistoricalTicksBidAskProtoBuf(&protoBA)

	reqID := NO_VALID_ID
	if protoBA.ReqId != nil {
		reqID = int64(protoBA.GetReqId())
	}

	done := false
	if protoBA.IsDone != nil {
		done = protoBA.GetIsDone()
	}

	var ticksBA []HistoricalTickBidAsk
	for _, baProto := range protoBA.GetHistoricalTicksBidAsk() {
		ba := decodeHistoricalTickBidAsk(baProto)
		ticksBA = append(ticksBA, *ba)
	}

	d.wrapper.HistoricalTicksBidAsk(reqID, ticksBA, done)
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

func (d *EDecoder) processHistoricalTicksLastMsgProtoBuf(msgBuf *MsgBuffer) {
	var protoLast protobuf.HistoricalTicksLast
	if err := proto.Unmarshal(msgBuf.bs, &protoLast); err != nil {
		log.Error().Err(err).Msg("failed to unmarshal HistoricalTicksLast message")
		return
	}

	d.wrapper.HistoricalTicksLastProtoBuf(&protoLast)

	reqID := NO_VALID_ID
	if protoLast.ReqId != nil {
		reqID = int64(protoLast.GetReqId())
	}

	done := false
	if protoLast.IsDone != nil {
		done = protoLast.GetIsDone()
	}

	var ticksLast []HistoricalTickLast
	for _, lastProto := range protoLast.GetHistoricalTicksLast() {
		last := decodeHistoricalTickLast(lastProto)
		ticksLast = append(ticksLast, *last)
	}

	d.wrapper.HistoricalTicksLast(reqID, ticksLast, done)
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

func (d *EDecoder) processTickByTickMsgProtoBuf(msgBuf *MsgBuffer) {
	var tickByTickDataProto protobuf.TickByTickData
	if err := proto.Unmarshal(msgBuf.bs, &tickByTickDataProto); err != nil {
		log.Error().Err(err).Msg("failed to unmarshal TickByTickData message")
		return
	}

	d.wrapper.TickByTickDataProtoBuf(&tickByTickDataProto)

	reqID := NO_VALID_ID
	if tickByTickDataProto.ReqId != nil {
		reqID = int64(tickByTickDataProto.GetReqId())
	}
	var tickType int64
	if tickByTickDataProto.TickType != nil {
		tickType = int64(tickByTickDataProto.GetTickType())
	}

	switch tickType {
	case 1, 2: // Last or AllLast
		if tickByTickDataProto.GetHistoricalTickLast() != nil {
			lastProto := tickByTickDataProto.GetHistoricalTickLast()
			ht := decodeHistoricalTickLast(lastProto)
			d.wrapper.TickByTickAllLast(
				reqID,
				tickType,
				ht.Time,
				ht.Price,
				ht.Size,
				ht.TickAttribLast,
				ht.Exchange,
				ht.SpecialConditions,
			)
		}
	case 3: // BidAsk
		if tickByTickDataProto.GetHistoricalTickBidAsk() != nil {
			baProto := tickByTickDataProto.GetHistoricalTickBidAsk()
			hba := decodeHistoricalTickBidAsk(baProto)
			d.wrapper.TickByTickBidAsk(
				reqID,
				hba.Time,
				hba.PriceBid,
				hba.PriceAsk,
				hba.SizeBid,
				hba.SizeAsk,
				hba.TickAttribBidAsk,
			)
		}
	case 4: // MidPoint
		if tickByTickDataProto.GetHistoricalTickMidPoint() != nil {
			midProto := tickByTickDataProto.GetHistoricalTickMidPoint()
			hm := decodeHistoricalTick(midProto)
			d.wrapper.TickByTickMidPoint(
				reqID,
				hm.Time,
				hm.Price,
			)
		}
	}
}

func (d *EDecoder) processOrderBoundMsg(msgBuf *MsgBuffer) {

	permID := msgBuf.decodeInt64()
	clientId := msgBuf.decodeInt64()
	orderId := msgBuf.decodeInt64()

	d.wrapper.OrderBound(permID, clientId, orderId)
}

func (d *EDecoder) processOrderBoundMsgProtoBuf(msgBuf *MsgBuffer) {

	var orderBoundProto protobuf.OrderBound
	err := proto.Unmarshal(msgBuf.bs, &orderBoundProto)
	if err != nil {
		log.Error().Err(err).Msg("failed to unmarshal OrderBound message")
		return
	}

	d.wrapper.OrderBoundProtoBuf(&orderBoundProto)

	var permID int64
	if orderBoundProto.PermId != nil {
		permID = int64(orderBoundProto.GetPermId())
	}
	var clientID int64
	if orderBoundProto.PermId != nil {
		clientID = int64(orderBoundProto.GetClientId())
	}
	var orderID int64
	if orderBoundProto.PermId != nil {
		orderID = int64(orderBoundProto.GetOrderId())
	}

	d.wrapper.OrderBound(permID, clientID, orderID)
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
	orderDecoder.decodeImbalanceOnly(msgBuf, MIN_CLIENT_VER)
	orderDecoder.decodeRouteMarketableToBbo(msgBuf)
	orderDecoder.decodeParentPermId(msgBuf)
	orderDecoder.decodeCompletedTime(msgBuf)
	orderDecoder.decodeCompletedStatus(msgBuf)
	orderDecoder.decodePegBestPegMidOrderAttributes(msgBuf)
	orderDecoder.decodeCustomerAccount(msgBuf)
	orderDecoder.decodeProfessionalCustomer(msgBuf)
	orderDecoder.decodeSubmitter(msgBuf)

	d.wrapper.CompletedOrder(contract, order, orderState)
}

func (d *EDecoder) processCompletedOrderMsgProtoBuf(msgBuf *MsgBuffer) {
	var completedOrderProto protobuf.CompletedOrder
	err := proto.Unmarshal(msgBuf.Bytes(), &completedOrderProto)
	if err != nil {
		log.Panic().Err(err).Msg("processOpenOrderEndMsgProtoBuf unmarshal error")
	}

	d.wrapper.CompletedOrderProtoBuf(&completedOrderProto)

	var contract *Contract
	if completedOrderProto.Contract != nil {
		contract = decodeContract(completedOrderProto.GetContract())
	}
	var order *Order
	if completedOrderProto.Order != nil {
		order = decodeOrder(UNSET_INT, completedOrderProto.GetContract(), completedOrderProto.GetOrder())
	}
	var orderState *OrderState
	if completedOrderProto.OrderState != nil {
		orderState = decodeOrderState(completedOrderProto.GetOrderState())
	}

	d.wrapper.CompletedOrder(contract, order, orderState)
}

func (d *EDecoder) processCompletedOrdersEndMsg(*MsgBuffer) {
	d.wrapper.CompletedOrdersEnd()
}

func (d *EDecoder) processCompletedOrdersEndMsgProtoBuf(msgBuf *MsgBuffer) {
	var completedOrdersEndProto protobuf.CompletedOrdersEnd
	err := proto.Unmarshal(msgBuf.Bytes(), &completedOrdersEndProto)
	if err != nil {
		log.Panic().Err(err).Msg("processOpenOrderEndMsgProtoBuf unmarshal error")
	}

	d.wrapper.CompletedOrdersEndProtoBuf(&completedOrdersEndProto)

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

func (d *EDecoder) processReplaceFAEndMsgProtoBuf(msgBuf *MsgBuffer) {
	var replaceFAEndProto protobuf.ReplaceFAEnd
	err := proto.Unmarshal(msgBuf.bs, &replaceFAEndProto)
	if err != nil {
		log.Error().Err(err).Msg("failed to unmarshal ReplaceFAEnd message")
		return
	}

	d.wrapper.ReplaceFAEndProtoBuf(&replaceFAEndProto)

	reqID := NO_VALID_ID
	if replaceFAEndProto.ReqId != nil {
		reqID = int64(replaceFAEndProto.GetReqId())
	}
	text := ""
	if replaceFAEndProto.Text != nil {
		text = replaceFAEndProto.GetText()
	}

	d.wrapper.ReplaceFAEnd(reqID, text)
}

func (d *EDecoder) processWshMetaDataMsgProtoBuf(msgBuf *MsgBuffer) {
	var wshMetaDataProto protobuf.WshMetaData
	if err := proto.Unmarshal(msgBuf.bs, &wshMetaDataProto); err != nil {
		log.Error().Err(err).Msg("failed to unmarshal WshMetaData message")
		return
	}
	d.wrapper.WshMetaDataProtoBuf(&wshMetaDataProto)

	reqID := NO_VALID_ID
	if wshMetaDataProto.ReqId != nil {
		reqID = int64(wshMetaDataProto.GetReqId())
	}
	dataJSON := ""
	if wshMetaDataProto.DataJson != nil {
		dataJSON = wshMetaDataProto.GetDataJson()
	}

	d.wrapper.WshMetaData(reqID, dataJSON)
}

func (d *EDecoder) processWshEventData(msgBuf *MsgBuffer) {

	reqID := msgBuf.decodeInt64()
	dataJSON := msgBuf.decodeString()

	d.wrapper.WshEventData(reqID, dataJSON)
}

func (d *EDecoder) processWshEventDataMsgProtoBuf(msgBuf *MsgBuffer) {
	var wshEventDataProto protobuf.WshEventData
	if err := proto.Unmarshal(msgBuf.bs, &wshEventDataProto); err != nil {
		log.Error().Err(err).Msg("failed to unmarshal WshEventData message")
		return
	}
	d.wrapper.WshEventDataProtoBuf(&wshEventDataProto)

	reqID := NO_VALID_ID
	if wshEventDataProto.ReqId != nil {
		reqID = int64(wshEventDataProto.GetReqId())
	}
	dataJSON := ""
	if wshEventDataProto.DataJson != nil {
		dataJSON = wshEventDataProto.GetDataJson()
	}

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

func (d *EDecoder) processHistoricalScheduleMsgProtoBuf(msgBuf *MsgBuffer) {
	var historicalScheduleProto protobuf.HistoricalSchedule
	err := proto.Unmarshal(msgBuf.bs, &historicalScheduleProto)
	if err != nil {
		log.Error().Err(err).Msg("failed to unmarshal HistoricalSchedule message")
		return
	}

	d.wrapper.HistoricalScheduleProtoBuf(&historicalScheduleProto)

	reqID := NO_VALID_ID
	if historicalScheduleProto.ReqId != nil {
		reqID = int64(historicalScheduleProto.GetReqId())
	}
	startDT := ""
	if historicalScheduleProto.StartDateTime != nil {
		startDT = historicalScheduleProto.GetEndDateTime()
	}
	endDT := ""
	if historicalScheduleProto.EndDateTime != nil {
		endDT = historicalScheduleProto.GetEndDateTime()
	}
	tz := ""
	if historicalScheduleProto.TimeZone != nil {
		tz = historicalScheduleProto.GetTimeZone()
	}

	var sessions []HistoricalSession
	for _, s := range historicalScheduleProto.GetHistoricalSessions() {
		var hs HistoricalSession
		if s.StartDateTime != nil {
			hs.StartDateTime = s.GetStartDateTime()
		}
		if s.EndDateTime != nil {
			hs.EndDateTime = s.GetEndDateTime()
		}
		if s.RefDate != nil {
			hs.RefDate = s.GetRefDate()
		}
		sessions = append(sessions, hs)
	}

	d.wrapper.HistoricalSchedule(reqID, startDT, endDT, tz, sessions)
}

func (d *EDecoder) processUserInfo(msgBuf *MsgBuffer) {

	reqID := msgBuf.decodeInt64()
	whiteBrandingId := msgBuf.decodeString()

	d.wrapper.UserInfo(reqID, whiteBrandingId)
}

func (d *EDecoder) processUserInfoMsgProtoBuf(msgBuf *MsgBuffer) {
	var protoMsg protobuf.UserInfo
	if err := proto.Unmarshal(msgBuf.bs, &protoMsg); err != nil {
		log.Error().Err(err).Msg("failed to unmarshal UserInfo")
		return
	}
	d.wrapper.UserInfoProtoBuf(&protoMsg)

	reqID := NO_VALID_ID
	if protoMsg.ReqId != nil {
		reqID = int64(protoMsg.GetReqId())
	}
	whiteBrandingId := ""
	if protoMsg.WhiteBrandingId != nil {
		whiteBrandingId = protoMsg.GetWhiteBrandingId()
	}

	d.wrapper.UserInfo(reqID, whiteBrandingId)
}

func (d *EDecoder) processCurrentTimeInMillisMsg(msgBuf *MsgBuffer) {

	timeInMillis := msgBuf.decodeInt64()

	d.wrapper.CurrentTimeInMillis(timeInMillis)
}

func (d *EDecoder) processCurrentTimeInMillisMsgProtoBuf(msgBuf *MsgBuffer) {

	var protoMsg protobuf.CurrentTimeInMillis
	if err := proto.Unmarshal(msgBuf.bs, &protoMsg); err != nil {
		log.Error().Err(err).Msg("failed to unmarshal CurrentTimeMillis")
		return
	}

	d.wrapper.CurrentTimeInMillisProtoBuf(&protoMsg)

	ms := int64(0)
	if protoMsg.CurrentTimeInMillis != nil {
		ms = protoMsg.GetCurrentTimeInMillis()
	}

	d.wrapper.CurrentTimeInMillis(ms)
}

func (d *EDecoder) processConfigResponseMsgProtoBuf(msgBuf *MsgBuffer) {

	var protoMsg protobuf.ConfigResponse
	if err := proto.Unmarshal(msgBuf.bs, &protoMsg); err != nil {
		log.Error().Err(err).Msg("failed to unmarshal ConfigResponse")
		return
	}

	d.wrapper.ConfigResponseProtoBuf(&protoMsg)
}

//
//		Helpers
//

func (d *EDecoder) decodeLastTradeDate(msgBuf *MsgBuffer, contract *ContractDetails, isBond bool) {
	lastTradeDateOrContractMonth := msgBuf.decodeString()
	setLastTradeDate(lastTradeDateOrContractMonth, contract, isBond)
}
