package ibapi

import (
	"fmt"
	"strconv"

	"github.com/scmhub/ibapi/protobuf"
)

func createExecutionFilterProto(execFilter *ExecutionFilter) *protobuf.ExecutionFilter {
	executionFilterProto := &protobuf.ExecutionFilter{}
	if execFilter == nil {
		return executionFilterProto
	}
	if isValidInt64Value(execFilter.ClientID) {
		clientID := int32(execFilter.ClientID)
		executionFilterProto.ClientId = &clientID
	}
	if !stringIsEmpty(execFilter.AcctCode) {
		executionFilterProto.AcctCode = &execFilter.AcctCode
	}
	if !stringIsEmpty(execFilter.Time) {
		executionFilterProto.Time = &execFilter.Time
	}
	if !stringIsEmpty(execFilter.Symbol) {
		executionFilterProto.Symbol = &execFilter.Symbol
	}
	if !stringIsEmpty(execFilter.SecType) {
		executionFilterProto.SecType = &execFilter.SecType
	}
	if !stringIsEmpty(execFilter.Exchange) {
		executionFilterProto.Exchange = &execFilter.Exchange
	}
	if !stringIsEmpty(execFilter.Side) {
		executionFilterProto.Side = &execFilter.Side
	}
	if isValidInt64Value(execFilter.LastNDays) {
		lastNDays := int32(execFilter.LastNDays)
		executionFilterProto.LastNDays = &lastNDays
	}
	if len(execFilter.SpecificDates) > 0 {
		for _, date := range execFilter.SpecificDates {
			executionFilterProto.SpecificDates = append(executionFilterProto.SpecificDates, int32(date))
		}
	}
	return executionFilterProto
}

func createExecutionRequestProto(reqID int64, execFilter *ExecutionFilter) *protobuf.ExecutionRequest {
	executionFilterProto := createExecutionFilterProto(execFilter)
	executionRequestProto := &protobuf.ExecutionRequest{}
	id := int32(reqID)
	executionRequestProto.ReqId = &id
	executionRequestProto.ExecutionFilter = executionFilterProto
	return executionRequestProto
}

func createPlaceOrderRequestProto(orderID OrderID, contract *Contract, order *Order) (*protobuf.PlaceOrderRequest, error) {
	var err error
	placeOrderRequestProto := &protobuf.PlaceOrderRequest{}
	if isValidInt64Value(orderID) {
		orderIDProto := int32(orderID)
		placeOrderRequestProto.OrderId = &orderIDProto
	}
	placeOrderRequestProto.Contract = createContractProto(contract, order)
	placeOrderRequestProto.Order, err = createOrderProto(order)
	if err != nil {
		return nil, err
	}
	placeOrderRequestProto.AttachedOrders = createAttachedOrdersProto(order)
	return placeOrderRequestProto, nil
}

func createAttachedOrdersProto(order *Order) *protobuf.AttachedOrders {
	attachedOrdersProto := &protobuf.AttachedOrders{}

	// Stop Loss Order
	if isValidInt64Value(order.SLOrderID) {
		slOrderID := int32(order.SLOrderID)
		attachedOrdersProto.SlOrderId = &slOrderID
	}
	if !stringIsEmpty(order.SLOrderType) {
		attachedOrdersProto.SlOrderType = &order.SLOrderType
	}

	// Profit Target Order
	if isValidInt64Value(order.PTOrderID) {
		ptOrderID := int32(order.PTOrderID)
		attachedOrdersProto.PtOrderId = &ptOrderID
	}
	if !stringIsEmpty(order.PTOrderType) {
		attachedOrdersProto.PtOrderType = &order.PTOrderType
	}

	return attachedOrdersProto
}

func createOrderProto(order *Order) (*protobuf.Order, error) {
	var err error
	orderProto := &protobuf.Order{}
	// order ids
	if isValidInt64Value(order.ClientID) {
		clientID := int32(order.ClientID)
		orderProto.ClientId = &clientID
	}
	if isValidInt64Value(order.OrderID) {
		orderID := int32(order.OrderID)
		orderProto.OrderId = &orderID
	}
	if isValidInt64Value(order.PermID) {
		orderProto.PermId = &order.PermID
	}
	if isValidInt64Value(order.ParentID) {
		parentID := int32(order.ParentID)
		orderProto.ParentId = &parentID
	}
	// primary attributes
	if !stringIsEmpty(order.Action) {
		orderProto.Action = &order.Action
	}
	if isValidDecimalValue(order.TotalQuantity) {
		totalQuantity := DecimalToString(order.TotalQuantity)
		orderProto.TotalQuantity = &totalQuantity
	}
	if isValidInt64Value(order.DisplaySize) {
		displaySize := int32(order.DisplaySize)
		orderProto.DisplaySize = &displaySize
	}
	if !stringIsEmpty(order.OrderType) {
		orderProto.OrderType = &order.OrderType
	}
	if isValidFloat64Value(order.LmtPrice) {
		orderProto.LmtPrice = &order.LmtPrice
	}
	if isValidFloat64Value(order.AuxPrice) {
		orderProto.AuxPrice = &order.AuxPrice
	}
	if !stringIsEmpty(order.TIF) {
		orderProto.Tif = &order.TIF
	}
	// clearing info
	if !stringIsEmpty(order.Account) {
		orderProto.Account = &order.Account
	}
	if !stringIsEmpty(order.SettlingFirm) {
		orderProto.SettlingFirm = &order.SettlingFirm
	}
	if !stringIsEmpty(order.ClearingAccount) {
		orderProto.ClearingAccount = &order.ClearingAccount
	}
	if !stringIsEmpty(order.ClearingIntent) {
		orderProto.ClearingIntent = &order.ClearingIntent
	}
	// secondary attributes
	if order.AllOrNone {
		orderProto.AllOrNone = &order.AllOrNone
	}
	if order.BlockOrder {
		orderProto.BlockOrder = &order.BlockOrder
	}
	if order.Hidden {
		orderProto.Hidden = &order.Hidden
	}
	if order.OutsideRTH {
		orderProto.OutsideRth = &order.OutsideRTH
	}
	if order.SweepToFill {
		orderProto.SweepToFill = &order.SweepToFill
	}
	if isValidFloat64Value(order.PercentOffset) {
		orderProto.PercentOffset = &order.PercentOffset
	}
	if isValidFloat64Value(order.TrailingPercent) {
		orderProto.TrailingPercent = &order.TrailingPercent
	}
	if isValidFloat64Value(order.TrailStopPrice) {
		orderProto.TrailStopPrice = &order.TrailStopPrice
	}
	if isValidInt64Value(order.MinQty) {
		minQty := int32(order.MinQty)
		orderProto.MinQty = &minQty
	}
	if !stringIsEmpty(order.GoodAfterTime) {
		orderProto.GoodAfterTime = &order.GoodAfterTime
	}
	if !stringIsEmpty(order.GoodTillDate) {
		orderProto.GoodTillDate = &order.GoodTillDate
	}
	if !stringIsEmpty(order.OCAGroup) {
		orderProto.OcaGroup = &order.OCAGroup
	}
	if !stringIsEmpty(order.OrderRef) {
		orderProto.OrderRef = &order.OrderRef
	}
	if !stringIsEmpty(order.Rule80A) {
		orderProto.Rule80A = &order.Rule80A
	}
	if isValidInt64Value(order.OCAType) {
		OCAType := int32(order.OCAType)
		orderProto.OcaType = &OCAType
	}
	if isValidInt64Value(order.TriggerMethod) {
		triggerMethod := int32(order.TriggerMethod)
		orderProto.TriggerMethod = &triggerMethod
	}
	// extended order fieldss
	if !stringIsEmpty(order.ActiveStartTime) {
		orderProto.ActiveStartTime = &order.ActiveStartTime
	}
	if !stringIsEmpty(order.ActiveStopTime) {
		orderProto.ActiveStopTime = &order.ActiveStopTime
	}
	// advisor allocation orders
	if !stringIsEmpty(order.FAGroup) {
		orderProto.FaGroup = &order.FAGroup
	}
	if !stringIsEmpty(order.FAMethod) {
		orderProto.FaMethod = &order.FAMethod
	}
	if !stringIsEmpty(order.FAPercentage) {
		orderProto.FaPercentage = &order.FAPercentage
	}
	// volatility orders
	if isValidFloat64Value(order.Volatility) {
		orderProto.Volatility = &order.Volatility
	}
	if isValidInt64Value(order.VolatilityType) {
		volatilityType := int32(order.VolatilityType)
		orderProto.VolatilityType = &volatilityType
	}
	if order.ContinuousUpdate {
		orderProto.ContinuousUpdate = &order.ContinuousUpdate
	}
	if isValidInt64Value(order.ReferencePriceType) {
		referencePriceType := int32(order.ReferencePriceType)
		orderProto.ReferencePriceType = &referencePriceType
	}
	if !stringIsEmpty(order.DeltaNeutralOrderType) {
		orderProto.DeltaNeutralOrderType = &order.DeltaNeutralOrderType
	}
	if isValidFloat64Value(order.DeltaNeutralAuxPrice) {
		orderProto.DeltaNeutralAuxPrice = &order.DeltaNeutralAuxPrice
	}
	if isValidInt64Value(order.DeltaNeutralConID) {
		deltaNeutralConID := int32(order.DeltaNeutralConID)
		orderProto.DeltaNeutralConId = &deltaNeutralConID
	}
	if !stringIsEmpty(order.DeltaNeutralOpenClose) {
		orderProto.DeltaNeutralOpenClose = &order.DeltaNeutralOpenClose
	}
	if order.DeltaNeutralShortSale {
		orderProto.DeltaNeutralShortSale = &order.DeltaNeutralShortSale
	}
	if isValidInt64Value(order.DeltaNeutralShortSaleSlot) {
		deltaNeutralShortSaleSlot := int32(order.DeltaNeutralShortSaleSlot)
		orderProto.DeltaNeutralShortSaleSlot = &deltaNeutralShortSaleSlot
	}
	if !stringIsEmpty(order.DeltaNeutralDesignatedLocation) {
		orderProto.DeltaNeutralDesignatedLocation = &order.DeltaNeutralDesignatedLocation
	}
	// scale orders
	if isValidInt64Value(order.ScaleInitLevelSize) {
		scaleInitLevelSize := int32(order.ScaleInitLevelSize)
		orderProto.ScaleInitLevelSize = &scaleInitLevelSize
	}
	if isValidInt64Value(order.ScaleSubsLevelSize) {
		scaleSubsLevelSize := int32(order.ScaleSubsLevelSize)
		orderProto.ScaleSubsLevelSize = &scaleSubsLevelSize
	}
	if isValidFloat64Value(order.ScalePriceIncrement) {
		orderProto.ScalePriceIncrement = &order.ScalePriceIncrement
	}
	if isValidFloat64Value(order.ScalePriceAdjustValue) {
		orderProto.ScalePriceAdjustValue = &order.ScalePriceAdjustValue
	}
	if isValidInt64Value(order.ScalePriceAdjustInterval) {
		scalePriceAdjustInterval := int32(order.ScalePriceAdjustInterval)
		orderProto.ScalePriceAdjustInterval = &scalePriceAdjustInterval
	}
	if isValidFloat64Value(order.ScaleProfitOffset) {
		orderProto.ScaleProfitOffset = &order.ScaleProfitOffset
	}
	if order.ScaleAutoReset {
		orderProto.ScaleAutoReset = &order.ScaleAutoReset
	}
	if isValidInt64Value(order.ScaleInitPosition) {
		scaleInitPosition := int32(order.ScaleInitPosition)
		orderProto.ScaleInitPosition = &scaleInitPosition
	}
	if isValidInt64Value(order.ScaleInitFillQty) {
		scaleInitFillQty := int32(order.ScaleInitFillQty)
		orderProto.ScaleInitFillQty = &scaleInitFillQty
	}
	if order.ScaleRandomPercent {
		orderProto.ScaleRandomPercent = &order.ScaleRandomPercent
	}
	if !stringIsEmpty(order.ScaleTable) {
		orderProto.ScaleTable = &order.ScaleTable
	}
	// hedge orders
	if !stringIsEmpty(order.HedgeType) {
		orderProto.HedgeType = &order.HedgeType
	}
	if !stringIsEmpty(order.HedgeParam) {
		orderProto.HedgeParam = &order.HedgeParam
	}

	// algo orders
	if !stringIsEmpty(order.AlgoStrategy) {
		orderProto.AlgoStrategy = &order.AlgoStrategy
		orderProto.AlgoParams = createStringStringMap(order.AlgoParams)
	}
	if !stringIsEmpty(order.AlgoID) {
		orderProto.AlgoId = &order.AlgoID
	}
	// combo orders
	if order.SmartComboRoutingParams != nil {
		orderProto.SmartComboRoutingParams = createStringStringMap(order.SmartComboRoutingParams)
	}

	// processing control
	if order.WhatIf {
		orderProto.WhatIf = &order.WhatIf
	}
	if order.Transmit {
		orderProto.Transmit = &order.Transmit
	}
	if order.OverridePercentageConstraints {
		orderProto.OverridePercentageConstraints = &order.OverridePercentageConstraints
	}

	// Institutional orders only
	if !stringIsEmpty(order.OpenClose) {
		orderProto.OpenClose = &order.OpenClose
	}
	if isValidInt64Value(order.Origin) {
		origin := int32(order.Origin)
		orderProto.Origin = &origin
	}
	if isValidInt64Value(order.ShortSaleSlot) {
		shortSaleSlot := int32(order.ShortSaleSlot)
		orderProto.ShortSaleSlot = &shortSaleSlot
	}
	if !stringIsEmpty(order.DesignatedLocation) {
		orderProto.DesignatedLocation = &order.DesignatedLocation
	}
	if isValidInt64Value(order.ExemptCode) {
		exemptCode := int32(order.ExemptCode)
		orderProto.ExemptCode = &exemptCode
	}
	if !stringIsEmpty(order.DeltaNeutralSettlingFirm) {
		orderProto.DeltaNeutralSettlingFirm = &order.DeltaNeutralSettlingFirm
	}
	if !stringIsEmpty(order.DeltaNeutralClearingAccount) {
		orderProto.DeltaNeutralClearingAccount = &order.DeltaNeutralClearingAccount
	}
	if !stringIsEmpty(order.DeltaNeutralClearingIntent) {
		orderProto.DeltaNeutralClearingIntent = &order.DeltaNeutralClearingIntent
	}
	// SMART routing only
	if isValidFloat64Value(order.DiscretionaryAmt) {
		orderProto.DiscretionaryAmt = &order.DiscretionaryAmt
	}
	if order.OptOutSmartRouting {
		orderProto.OptOutSmartRouting = &order.OptOutSmartRouting
	}

	// BOX ORDERS ONLY
	if isValidFloat64Value(order.StartingPrice) {
		orderProto.StartingPrice = &order.StartingPrice
	}
	if isValidFloat64Value(order.StockRefPrice) {
		orderProto.StockRefPrice = &order.StockRefPrice
	}
	if isValidFloat64Value(order.Delta) {
		orderProto.Delta = &order.Delta
	}

	// pegged to stock or VOL orders
	if isValidFloat64Value(order.StockRangeLower) {
		orderProto.StockRangeLower = &order.StockRangeLower
	}
	if isValidFloat64Value(order.StockRangeUpper) {
		orderProto.StockRangeUpper = &order.StockRangeUpper
	}

	// Not Held
	if order.NotHeld {
		orderProto.NotHeld = &order.NotHeld
	}

	// order misc options
	if len(order.OrderMiscOptions) > 0 {
		orderProto.OrderMiscOptions = createStringStringMap(order.OrderMiscOptions)
	}

	// order algo id
	if order.Solicited {
		orderProto.Solicited = &order.Solicited
	}

	if order.RandomizeSize {
		orderProto.RandomizeSize = &order.RandomizeSize
	}
	if order.RandomizePrice {
		orderProto.RandomizePrice = &order.RandomizePrice
	}

	// PEG2BENCH fields
	if isValidInt64Value(order.ReferenceContractID) {
		refId := int32(order.ReferenceContractID)
		orderProto.ReferenceContractId = &refId
	}
	if isValidFloat64Value(order.PeggedChangeAmount) {
		orderProto.PeggedChangeAmount = &order.PeggedChangeAmount
	}
	if order.IsPeggedChangeAmountDecrease {
		orderProto.IsPeggedChangeAmountDecrease = &order.IsPeggedChangeAmountDecrease
	}
	if isValidFloat64Value(order.ReferenceChangeAmount) {
		orderProto.ReferenceChangeAmount = &order.ReferenceChangeAmount
	}
	if !stringIsEmpty(order.ReferenceExchangeID) {
		orderProto.ReferenceExchangeId = &order.ReferenceExchangeID
	}
	if !stringIsEmpty(order.AdjustedOrderType) {
		orderProto.AdjustedOrderType = &order.AdjustedOrderType
	}
	if isValidFloat64Value(order.TriggerPrice) {
		orderProto.TriggerPrice = &order.TriggerPrice
	}
	if isValidFloat64Value(order.AdjustedStopPrice) {
		orderProto.AdjustedStopPrice = &order.AdjustedStopPrice
	}
	if isValidFloat64Value(order.AdjustedStopLimitPrice) {
		orderProto.AdjustedStopLimitPrice = &order.AdjustedStopLimitPrice
	}
	if isValidFloat64Value(order.AdjustedTrailingAmount) {
		orderProto.AdjustedTrailingAmount = &order.AdjustedTrailingAmount
	}
	if isValidInt64Value(order.AdjustableTrailingUnit) {
		unit := int32(order.AdjustableTrailingUnit)
		orderProto.AdjustableTrailingUnit = &unit
	}
	if isValidFloat64Value(order.LmtPriceOffset) {
		orderProto.LmtPriceOffset = &order.LmtPriceOffset
	}

	if order.Conditions != nil {
		orderProto.Conditions, err = createConditionsProto(order)
		if err != nil {
			return nil, err
		}
	}
	if order.ConditionsCancelOrder {
		orderProto.ConditionsCancelOrder = &order.ConditionsCancelOrder
	}
	if order.ConditionsIgnoreRth {
		orderProto.ConditionsIgnoreRth = &order.ConditionsIgnoreRth
	}

	// models
	if !stringIsEmpty(order.ModelCode) {
		orderProto.ModelCode = &order.ModelCode
	}
	if !stringIsEmpty(order.ExtOperator) {
		orderProto.ExtOperator = &order.ExtOperator
	}
	orderProto.SoftDollarTier = createSoftDollarTierProto(order)

	// native cash quantity
	if isValidFloat64Value(order.CashQty) {
		orderProto.CashQty = &order.CashQty
	}

	if !stringIsEmpty(order.Mifid2DecisionMaker) {
		orderProto.Mifid2DecisionMaker = &order.Mifid2DecisionMaker
	}
	if !stringIsEmpty(order.Mifid2DecisionAlgo) {
		orderProto.Mifid2DecisionAlgo = &order.Mifid2DecisionAlgo
	}
	if !stringIsEmpty(order.Mifid2ExecutionTrader) {
		orderProto.Mifid2ExecutionTrader = &order.Mifid2ExecutionTrader
	}
	if !stringIsEmpty(order.Mifid2ExecutionAlgo) {
		orderProto.Mifid2ExecutionAlgo = &order.Mifid2ExecutionAlgo
	}

	// don't use auto price for hedge
	if order.DontUseAutoPriceForHedge {
		orderProto.DontUseAutoPriceForHedge = &order.DontUseAutoPriceForHedge
	}

	if order.IsOmsContainer {
		orderProto.IsOmsContainer = &order.IsOmsContainer
	}
	if order.DiscretionaryUpToLimitPrice {
		orderProto.DiscretionaryUpToLimitPrice = &order.DiscretionaryUpToLimitPrice
	}
	// //////
	// if !stringIsEmpty(order.AutoCancelDate) {
	// 	orderProto.AutoCancelDate = &order.AutoCancelDate
	// }
	// if isValidDecimalValue(order.FilledQuantity) {
	// 	filledQty := DecimalToString(order.FilledQuantity)
	// 	orderProto.FilledQuantity = &filledQty
	// }
	// if isValidInt64Value(order.RefFuturesConID) {
	// 	refFut := int32(order.RefFuturesConID)
	// 	orderProto.RefFuturesConId = &refFut
	// }
	// if !stringIsEmpty(order.Shareholder) {
	// 	orderProto.Shareholder = &order.Shareholder
	// }
	if order.RouteMarketableToBbo != STATE_DEFAULT {
		val := int32(order.RouteMarketableToBbo)
		orderProto.RouteMarketableToBbo = &val
	}
	// if isValidInt64Value(order.ParentPermID) {
	// 	orderProto.ParentPermId = &order.ParentPermID
	// }
	// /////
	if isValidInt64Value(order.UsePriceMgmtAlgo) {
		val := int32(order.UsePriceMgmtAlgo)
		orderProto.UsePriceMgmtAlgo = &val
	}
	if isValidInt64Value(order.Duration) {
		val := int32(order.Duration)
		orderProto.Duration = &val
	}
	if isValidInt64Value(order.PostToAts) {
		val := int32(order.PostToAts)
		orderProto.PostToAts = &val
	}
	if !stringIsEmpty(order.AdvancedErrorOverride) {
		orderProto.AdvancedErrorOverride = &order.AdvancedErrorOverride
	}
	if !stringIsEmpty(order.ManualOrderTime) {
		orderProto.ManualOrderTime = &order.ManualOrderTime
	}
	if isValidInt64Value(order.MinTradeQty) {
		val := int32(order.MinTradeQty)
		orderProto.MinTradeQty = &val
	}
	if isValidInt64Value(order.MinCompeteSize) {
		val := int32(order.MinCompeteSize)
		orderProto.MinCompeteSize = &val
	}
	if isValidFloat64Value(order.CompeteAgainstBestOffset) {
		orderProto.CompeteAgainstBestOffset = &order.CompeteAgainstBestOffset
	}
	if isValidFloat64Value(order.MidOffsetAtWhole) {
		orderProto.MidOffsetAtWhole = &order.MidOffsetAtWhole
	}
	if isValidFloat64Value(order.MidOffsetAtHalf) {
		orderProto.MidOffsetAtHalf = &order.MidOffsetAtHalf
	}
	if !stringIsEmpty(order.CustomerAccount) {
		orderProto.CustomerAccount = &order.CustomerAccount
	}
	if order.ProfessionalCustomer {
		orderProto.ProfessionalCustomer = &order.ProfessionalCustomer
	}
	if !stringIsEmpty(order.BondAccruedInterest) {
		orderProto.BondAccruedInterest = &order.BondAccruedInterest
	}
	if order.IncludeOvernight {
		orderProto.IncludeOvernight = &order.IncludeOvernight
	}
	if isValidInt64Value(order.ManualOrderIndicator) {
		val := int32(order.ManualOrderIndicator)
		orderProto.ManualOrderIndicator = &val
	}
	if !stringIsEmpty(order.Submitter) {
		orderProto.Submitter = &order.Submitter
	}
	if order.AutoCancelParent {
		orderProto.AutoCancelParent = &order.AutoCancelParent
	}
	if order.ImbalanceOnly {
		orderProto.ImbalanceOnly = &order.ImbalanceOnly
	}
	if order.PostOnly {
		orderProto.PostOnly = &order.PostOnly
	}
	if order.AllowPreOpen {
		orderProto.AllowPreOpen = &order.AllowPreOpen
	}
	if order.IgnoreOpenAuction {
		orderProto.IgnoreOpenAuction = &order.IgnoreOpenAuction
	}
	if order.Deactivate {
		orderProto.Deactivate = &order.Deactivate
	}
	if order.SeekPriceImprovement != STATE_DEFAULT {
		val := int32(order.SeekPriceImprovement)
		orderProto.SeekPriceImprovement = &val
	}
	if isValidInt64Value(order.WhatIfType) {
		val := int32(order.WhatIfType)
		orderProto.WhatIfType = &val
	}
	return orderProto, nil
}

func createConditionsProto(order *Order) ([]*protobuf.OrderCondition, error) {
	var orderConditionList []*protobuf.OrderCondition
	for _, condition := range order.Conditions {
		var protoCond *protobuf.OrderCondition
		switch condition.Type() {
		case PriceOrderCondition:
			protoCond = createPriceConditionProto(condition)
		case TimeOrderCondition:
			protoCond = createTimeConditionProto(condition)
		case MarginOrderCondition:
			protoCond = createMarginConditionProto(condition)
		case ExecutionOrderCondition:
			protoCond = createExecutionConditionProto(condition)
		case VolumeOrderCondition:
			protoCond = createVolumeConditionProto(condition)
		case PercentChangeOrderCondition:
			protoCond = createPercentChangeConditionProto(condition)
		default:
			return nil, fmt.Errorf("unknown condition type: %v", condition)
		}
		orderConditionList = append(orderConditionList, protoCond)
	}
	return orderConditionList, nil
}

// Base
func createOrderConditionProto(cond OrderCondition) *protobuf.OrderCondition {
	protoCond := &protobuf.OrderCondition{}
	if isValidInt64Value(int64(cond.Type())) {
		t := int32(cond.Type())
		protoCond.Type = &t
	}
	isConjunctionConnection := cond.IsConjunctionConnection()
	protoCond.IsConjunctionConnection = &isConjunctionConnection
	return protoCond
}

// Operator
func createOperatorConditionProto(cond OrderCondition) *protobuf.OrderCondition {
	protoCond := createOrderConditionProto(cond)
	if op, ok := cond.(*operatorCondition); ok {
		isMore := op.IsMore
		protoCond.IsMore = &isMore
	}
	return protoCond
}

// Contract
func createContractConditionProto(cond OrderCondition) *protobuf.OrderCondition {
	protoCond := createOperatorConditionProto(cond)
	if cc, ok := cond.(*contractCondition); ok {
		if isValidInt64Value(cc.ConID) {
			conID := int32(cc.ConID)
			protoCond.ConId = &conID
		}
		if !stringIsEmpty(cc.Exchange) {
			protoCond.Exchange = &cc.Exchange
		}
	}
	return protoCond
}

// Price
func createPriceConditionProto(cond OrderCondition) *protobuf.OrderCondition {
	protoCond := createContractConditionProto(cond)
	if pc, ok := cond.(*PriceCondition); ok {
		if isValidFloat64Value(pc.Price) {
			protoCond.Price = &pc.Price
		}
		if isValidInt64Value(int64(pc.TriggerMethod)) {
			tm := int32(pc.TriggerMethod)
			protoCond.TriggerMethod = &tm
		}
	}
	return protoCond
}

// Time
func createTimeConditionProto(cond OrderCondition) *protobuf.OrderCondition {
	protoCond := createOperatorConditionProto(cond)
	if tc, ok := cond.(*TimeCondition); ok {
		if !stringIsEmpty(tc.Time) {
			protoCond.Time = &tc.Time
		}
	}
	return protoCond
}

// Margin
func createMarginConditionProto(cond OrderCondition) *protobuf.OrderCondition {
	protoCond := createOperatorConditionProto(cond)
	if mc, ok := cond.(*MarginCondition); ok {
		if isValidInt64Value(mc.Percent) {
			percent := int32(mc.Percent)
			protoCond.Percent = &percent
		}
	}
	return protoCond
}

// Execution
func createExecutionConditionProto(cond OrderCondition) *protobuf.OrderCondition {
	protoCond := createOrderConditionProto(cond)
	if ec, ok := cond.(*ExecutionCondition); ok {
		if !stringIsEmpty(ec.SecType) {
			protoCond.SecType = &ec.SecType
		}
		if !stringIsEmpty(ec.Exchange) {
			protoCond.Exchange = &ec.Exchange
		}
		if !stringIsEmpty(ec.Symbol) {
			protoCond.Symbol = &ec.Symbol
		}
	}
	return protoCond
}

// Volume
func createVolumeConditionProto(cond OrderCondition) *protobuf.OrderCondition {
	protoCond := createContractConditionProto(cond)
	if vc, ok := cond.(*VolumeCondition); ok {
		if isValidInt64Value(vc.Volume) {
			volume := int32(vc.Volume)
			protoCond.Volume = &volume
		}
	}
	return protoCond
}

// PercentChange
func createPercentChangeConditionProto(cond OrderCondition) *protobuf.OrderCondition {
	protoCond := createContractConditionProto(cond)
	if pc, ok := cond.(*PercentChangeCondition); ok {
		if isValidFloat64Value(pc.ChangePercent) {
			protoCond.ChangePercent = &pc.ChangePercent
		}
	}
	return protoCond
}

// SoftDollarTier
func createSoftDollarTierProto(order *Order) *protobuf.SoftDollarTier {
	tier := order.SoftDollarTier
	softDollarTierProto := &protobuf.SoftDollarTier{}
	if !stringIsEmpty(tier.Name) {
		softDollarTierProto.Name = &tier.Name
	}
	if !stringIsEmpty(tier.Value) {
		softDollarTierProto.Value = &tier.Value
	}
	if !stringIsEmpty(tier.DisplayName) {
		softDollarTierProto.DisplayName = &tier.DisplayName
	}
	return softDollarTierProto
}

func createStringStringMap(tagValueList []TagValue) map[string]string {
	stringStringMap := make(map[string]string)
	for _, tagValue := range tagValueList {
		stringStringMap[tagValue.Tag] = tagValue.Value
	}
	return stringStringMap
}

// Contract
func createContractProto(contract *Contract, order *Order) *protobuf.Contract {
	contractProto := &protobuf.Contract{}

	if isValidInt64Value(contract.ConID) {
		conId := int32(contract.ConID)
		contractProto.ConId = &conId
	}
	if !stringIsEmpty(contract.Symbol) {
		contractProto.Symbol = &contract.Symbol
	}
	if !stringIsEmpty(contract.SecType) {
		contractProto.SecType = &contract.SecType
	}
	if !stringIsEmpty(contract.LastTradeDateOrContractMonth) {
		contractProto.LastTradeDateOrContractMonth = &contract.LastTradeDateOrContractMonth
	}
	if isValidFloat64Value(contract.Strike) {
		contractProto.Strike = &contract.Strike
	}
	if !stringIsEmpty(contract.Right) {
		contractProto.Right = &contract.Right
	}
	if !stringIsEmpty(contract.Multiplier) {
		if multiplier, err := strconv.ParseFloat(contract.Multiplier, 64); err == nil {
			contractProto.Multiplier = &multiplier
		}
	}
	if !stringIsEmpty(contract.Exchange) {
		contractProto.Exchange = &contract.Exchange
	}
	if !stringIsEmpty(contract.PrimaryExchange) {
		contractProto.PrimaryExch = &contract.PrimaryExchange
	}
	if !stringIsEmpty(contract.Currency) {
		contractProto.Currency = &contract.Currency
	}
	if !stringIsEmpty(contract.LocalSymbol) {
		contractProto.LocalSymbol = &contract.LocalSymbol
	}
	if !stringIsEmpty(contract.TradingClass) {
		contractProto.TradingClass = &contract.TradingClass
	}
	if !stringIsEmpty(contract.SecIDType) {
		contractProto.SecIdType = &contract.SecIDType
	}
	if !stringIsEmpty(contract.SecID) {
		contractProto.SecId = &contract.SecID
	}
	if !stringIsEmpty(contract.Description) {
		contractProto.Description = &contract.Description
	}
	if !stringIsEmpty(contract.IssuerID) {
		contractProto.IssuerId = &contract.IssuerID
	}
	if contract.DeltaNeutralContract != nil {
		contractProto.DeltaNeutralContract = createDeltaNeutralContractProto(contract.DeltaNeutralContract)
	}
	if contract.IncludeExpired {
		contractProto.IncludeExpired = &contract.IncludeExpired
	}
	if !stringIsEmpty(contract.ComboLegsDescrip) {
		contractProto.ComboLegsDescrip = &contract.ComboLegsDescrip
	}
	if len(contract.ComboLegs) > 0 {
		var legs []OrderComboLeg
		if order != nil {
			legs = order.OrderComboLegs
		}
		contractProto.ComboLegs = createComboLegsProto(contract.ComboLegs, legs)
	}

	return contractProto
}

// Delta Neutral Contract
func createDeltaNeutralContractProto(dnc *DeltaNeutralContract) *protobuf.DeltaNeutralContract {
	if dnc == nil {
		return nil
	}
	dncProto := &protobuf.DeltaNeutralContract{}
	if isValidInt64Value(dnc.ConID) {
		conId := int32(dnc.ConID)
		dncProto.ConId = &conId
	}
	if isValidFloat64Value(dnc.Delta) {
		dncProto.Delta = &dnc.Delta
	}
	if isValidFloat64Value(dnc.Price) {
		dncProto.Price = &dnc.Price
	}
	return dncProto
}

func createComboLegsProto(comboLegs []ComboLeg, orderComboLegs []OrderComboLeg) []*protobuf.ComboLeg {
	var comboLegProtoList []*protobuf.ComboLeg
	for i, comboLeg := range comboLegs {
		var perLegPrice float64
		if i < len(orderComboLegs) {
			perLegPrice = orderComboLegs[i].Price
		} else {
			perLegPrice = UNSET_FLOAT // define this as needed, e.g. math.NaN()
		}
		comboLegProto := createComboLegProto(&comboLeg, perLegPrice)
		comboLegProtoList = append(comboLegProtoList, comboLegProto)
	}
	return comboLegProtoList
}

func createComboLegProto(comboLeg *ComboLeg, perLegPrice float64) *protobuf.ComboLeg {
	comboLegProto := &protobuf.ComboLeg{}
	if isValidInt64Value(comboLeg.ConID) {
		conId := int32(comboLeg.ConID)
		comboLegProto.ConId = &conId
	}
	if isValidInt64Value(comboLeg.Ratio) {
		ratio := int32(comboLeg.Ratio)
		comboLegProto.Ratio = &ratio
	}
	if !stringIsEmpty(comboLeg.Action) {
		comboLegProto.Action = &comboLeg.Action
	}
	if !stringIsEmpty(comboLeg.Exchange) {
		comboLegProto.Exchange = &comboLeg.Exchange
	}
	if isValidInt64Value(comboLeg.OpenClose) {
		openClose := int32(comboLeg.OpenClose)
		comboLegProto.OpenClose = &openClose
	}
	if isValidInt64Value(comboLeg.ShortSaleSlot) {
		shortSaleSlot := int32(comboLeg.ShortSaleSlot)
		comboLegProto.ShortSalesSlot = &shortSaleSlot
	}
	if !stringIsEmpty(comboLeg.DesignatedLocation) {
		comboLegProto.DesignatedLocation = &comboLeg.DesignatedLocation
	}
	if isValidInt64Value(comboLeg.ExemptCode) {
		exemptCode := int32(comboLeg.ExemptCode)
		comboLegProto.ExemptCode = &exemptCode
	}
	if isValidFloat64Value(perLegPrice) {
		comboLegProto.PerLegPrice = &perLegPrice
	}
	return comboLegProto
}

func createCancelOrderRequestProto(orderID OrderID, orderCancel *OrderCancel) *protobuf.CancelOrderRequest {
	cancelOrderRequestProto := &protobuf.CancelOrderRequest{}
	if isValidInt64Value(orderID) {
		id := int32(orderID)
		cancelOrderRequestProto.OrderId = &id
	}
	cancelOrderRequestProto.OrderCancel = createOrderCancelProto(orderCancel)
	return cancelOrderRequestProto
}

func createGlobalCancelRequestProto(orderCancel *OrderCancel) *protobuf.GlobalCancelRequest {
	globalCancelRequestProto := &protobuf.GlobalCancelRequest{}
	globalCancelRequestProto.OrderCancel = createOrderCancelProto(orderCancel)
	return globalCancelRequestProto
}

func createOrderCancelProto(orderCancel *OrderCancel) *protobuf.OrderCancel {
	orderCancelProto := &protobuf.OrderCancel{}
	if !stringIsEmpty(orderCancel.ManualOrderCancelTime) {
		orderCancelProto.ManualOrderCancelTime = &orderCancel.ManualOrderCancelTime
	}
	if !stringIsEmpty(orderCancel.ExtOperator) {
		orderCancelProto.ExtOperator = &orderCancel.ExtOperator
	}
	if isValidInt64Value(orderCancel.ManualOrderIndicator) {
		indicator := int32(orderCancel.ManualOrderIndicator)
		orderCancelProto.ManualOrderIndicator = &indicator
	}
	return orderCancelProto
}

func createAllOpenOrdersRequestProto() *protobuf.AllOpenOrdersRequest {
	return &protobuf.AllOpenOrdersRequest{}
}

func createAutoOpenOrdersRequestProto(autoBind bool) *protobuf.AutoOpenOrdersRequest {
	autoOpenOrdersRequestProto := &protobuf.AutoOpenOrdersRequest{}
	if autoBind {
		autoOpenOrdersRequestProto.AutoBind = &autoBind
	}
	return autoOpenOrdersRequestProto
}

func createOpenOrdersRequestProto() *protobuf.OpenOrdersRequest {
	return &protobuf.OpenOrdersRequest{}
}

func createCompletedOrdersRequestProto(apiOnly bool) *protobuf.CompletedOrdersRequest {
	completedOrdersRequestProto := &protobuf.CompletedOrdersRequest{}
	if apiOnly {
		completedOrdersRequestProto.ApiOnly = &apiOnly
	}
	return completedOrdersRequestProto
}

func createContractDataRequestProto(reqID int64, contract *Contract) *protobuf.ContractDataRequest {
	contractDataRequestProto := &protobuf.ContractDataRequest{}
	if isValidInt64Value(reqID) {
		id := int32(reqID)
		contractDataRequestProto.ReqId = &id
	}
	contractDataRequestProto.Contract = createContractProto(contract, nil)
	return contractDataRequestProto
}

func createMarketDataRequestProto(reqID int64, contract *Contract, genericTickList string, snapshot bool, regulatorySnapshot bool, marketDataOptionsList []TagValue) *protobuf.MarketDataRequest {
	marketDataRequestProto := &protobuf.MarketDataRequest{}
	if isValidInt64Value(reqID) {
		id := int32(reqID)
		marketDataRequestProto.ReqId = &id
	}
	marketDataRequestProto.Contract = createContractProto(contract, nil)
	if !stringIsEmpty(genericTickList) {
		marketDataRequestProto.GenericTickList = &genericTickList
	}

	if snapshot {
		marketDataRequestProto.Snapshot = &snapshot
	}

	if regulatorySnapshot {
		marketDataRequestProto.RegulatorySnapshot = &regulatorySnapshot
	}

	marketDataOptionsMap := createStringStringMap(marketDataOptionsList)
	if len(marketDataOptionsMap) > 0 {
		if marketDataRequestProto.MarketDataOptions == nil {
			marketDataRequestProto.MarketDataOptions = make(map[string]string)
		}
		for key, value := range marketDataOptionsMap {
			marketDataRequestProto.MarketDataOptions[key] = value
		}
	}

	return marketDataRequestProto
}

func createMarketDepthRequestProto(reqID int64, contract *Contract, numRows int64, isSmartDepth bool, marketDepthOptionsList []TagValue) *protobuf.MarketDepthRequest {
	marketDepthRequestProto := &protobuf.MarketDepthRequest{}

	if isValidInt64Value(reqID) {
		id := int32(reqID)
		marketDepthRequestProto.ReqId = &id
	}

	marketDepthRequestProto.Contract = createContractProto(contract, nil)

	if isValidInt64Value(numRows) {
		nr := int32(numRows)
		marketDepthRequestProto.NumRows = &nr
	}

	if isSmartDepth {
		marketDepthRequestProto.IsSmartDepth = &isSmartDepth
	}

	marketDepthOptionsMap := createStringStringMap(marketDepthOptionsList)
	if len(marketDepthOptionsMap) > 0 {
		if marketDepthRequestProto.MarketDepthOptions == nil {
			marketDepthRequestProto.MarketDepthOptions = make(map[string]string)
		}
		for key, value := range marketDepthOptionsMap {
			marketDepthRequestProto.MarketDepthOptions[key] = value
		}
	}

	return marketDepthRequestProto
}

func createMarketDataTypeRequestProto(marketDataType int64) *protobuf.MarketDataTypeRequest {
	marketDataTypeRequestProto := &protobuf.MarketDataTypeRequest{}

	if isValidInt64Value(marketDataType) {
		mdt := int32(marketDataType)
		marketDataTypeRequestProto.MarketDataType = &mdt
	}

	return marketDataTypeRequestProto
}

func createCancelMarketDataProto(reqID int64) *protobuf.CancelMarketData {
	cancelMarketDataProto := &protobuf.CancelMarketData{}

	if isValidInt64Value(reqID) {
		id := int32(reqID)
		cancelMarketDataProto.ReqId = &id
	}

	return cancelMarketDataProto
}

func createCancelMarketDepthProto(reqID int64, isSmartDepth bool) *protobuf.CancelMarketDepth {
	cancelMarketDepthProto := &protobuf.CancelMarketDepth{}

	if isValidInt64Value(reqID) {
		id := int32(reqID)
		cancelMarketDepthProto.ReqId = &id
	}

	if isSmartDepth {
		cancelMarketDepthProto.IsSmartDepth = &isSmartDepth
	}

	return cancelMarketDepthProto
}

func createAccountDataRequestProto(subscribe bool, acctCode string) *protobuf.AccountDataRequest {
	accountDataRequestProto := &protobuf.AccountDataRequest{}

	if subscribe {
		accountDataRequestProto.Subscribe = &subscribe
	}

	if !stringIsEmpty(acctCode) {
		accountDataRequestProto.AcctCode = &acctCode
	}

	return accountDataRequestProto
}

func createManagedAccountsRequestProto() *protobuf.ManagedAccountsRequest {
	return &protobuf.ManagedAccountsRequest{}
}

func createPositionsRequestProto() *protobuf.PositionsRequest {
	return &protobuf.PositionsRequest{}
}

func createCancelPositionsRequestProto() *protobuf.CancelPositions {
	return &protobuf.CancelPositions{}
}

func createAccountSummaryRequestProto(reqID int64, group string, tags string) *protobuf.AccountSummaryRequest {
	accountSummaryRequestProto := &protobuf.AccountSummaryRequest{}

	if isValidInt64Value(reqID) {
		id := int32(reqID)
		accountSummaryRequestProto.ReqId = &id
	}

	if !stringIsEmpty(group) {
		accountSummaryRequestProto.Group = &group
	}

	if !stringIsEmpty(tags) {
		accountSummaryRequestProto.Tags = &tags
	}

	return accountSummaryRequestProto
}

func createCancelAccountSummaryRequestProto(reqID int64) *protobuf.CancelAccountSummary {
	cancelAccountSummaryProto := &protobuf.CancelAccountSummary{}

	if isValidInt64Value(reqID) {
		id := int32(reqID)
		cancelAccountSummaryProto.ReqId = &id
	}

	return cancelAccountSummaryProto
}

func createPositionsMultiRequestProto(reqID int64, account string, modelCode string) *protobuf.PositionsMultiRequest {
	positionsMultiRequestProto := &protobuf.PositionsMultiRequest{}

	if isValidInt64Value(reqID) {
		id := int32(reqID)
		positionsMultiRequestProto.ReqId = &id
	}

	if !stringIsEmpty(account) {
		positionsMultiRequestProto.Account = &account
	}

	if !stringIsEmpty(modelCode) {
		positionsMultiRequestProto.ModelCode = &modelCode
	}

	return positionsMultiRequestProto
}

func createCancelPositionsMultiRequestProto(reqID int64) *protobuf.CancelPositionsMulti {
	cancelPositionsMultiProto := &protobuf.CancelPositionsMulti{}

	if isValidInt64Value(reqID) {
		id := int32(reqID)
		cancelPositionsMultiProto.ReqId = &id
	}

	return cancelPositionsMultiProto
}

func createAccountUpdatesMultiRequestProto(reqID int64, account string, modelCode string, ledgerAndNLV bool) *protobuf.AccountUpdatesMultiRequest {
	accountUpdatesMultiRequestProto := &protobuf.AccountUpdatesMultiRequest{}

	if isValidInt64Value(reqID) {
		id := int32(reqID)
		accountUpdatesMultiRequestProto.ReqId = &id
	}

	if !stringIsEmpty(account) {
		accountUpdatesMultiRequestProto.Account = &account
	}

	if !stringIsEmpty(modelCode) {
		accountUpdatesMultiRequestProto.ModelCode = &modelCode
	}

	if ledgerAndNLV {
		accountUpdatesMultiRequestProto.LedgerAndNLV = &ledgerAndNLV
	}

	return accountUpdatesMultiRequestProto
}

func createCancelAccountUpdatesMultiRequestProto(reqID int64) *protobuf.CancelAccountUpdatesMulti {
	cancelAccountUpdatesMultiProto := &protobuf.CancelAccountUpdatesMulti{}

	if isValidInt64Value(reqID) {
		id := int32(reqID)
		cancelAccountUpdatesMultiProto.ReqId = &id
	}

	return cancelAccountUpdatesMultiProto
}

func createHistoricalDataRequestProto(reqID int64, contract *Contract, endDateTime string, duration string,
	barSizeSetting string, whatToShow string, useRTH bool, formatDate int, keepUpToDate bool, chartOptionsList []TagValue) *protobuf.HistoricalDataRequest {

	historicalDataRequestProto := &protobuf.HistoricalDataRequest{}

	if isValidInt64Value(reqID) {
		id := int32(reqID)
		historicalDataRequestProto.ReqId = &id
	}

	order := &Order{}
	historicalDataRequestProto.Contract = createContractProto(contract, order)

	if !stringIsEmpty(endDateTime) {
		historicalDataRequestProto.EndDateTime = &endDateTime
	}

	if !stringIsEmpty(duration) {
		historicalDataRequestProto.Duration = &duration
	}

	if !stringIsEmpty(barSizeSetting) {
		historicalDataRequestProto.BarSizeSetting = &barSizeSetting
	}

	if !stringIsEmpty(whatToShow) {
		historicalDataRequestProto.WhatToShow = &whatToShow
	}

	if useRTH {
		historicalDataRequestProto.UseRTH = &useRTH
	}

	if isValidInt64Value(int64(formatDate)) {
		formatDate := int32(formatDate)
		historicalDataRequestProto.FormatDate = &formatDate
	}

	if keepUpToDate {
		historicalDataRequestProto.KeepUpToDate = &keepUpToDate
	}

	chartOptionsMap := createStringStringMap(chartOptionsList)
	if len(chartOptionsMap) > 0 {
		historicalDataRequestProto.ChartOptions = chartOptionsMap
	}

	return historicalDataRequestProto
}

func createRealTimeBarsRequestProto(reqID int64, contract *Contract, barSize int, whatToShow string, useRTH bool,
	realTimeBarsOptionsList []TagValue) *protobuf.RealTimeBarsRequest {

	realTimeBarsRequestProto := &protobuf.RealTimeBarsRequest{}

	if isValidInt64Value(reqID) {
		id := int32(reqID)
		realTimeBarsRequestProto.ReqId = &id
	}

	order := &Order{}
	realTimeBarsRequestProto.Contract = createContractProto(contract, order)

	if isValidInt64Value(int64(barSize)) {
		barSize := int32(barSize)
		realTimeBarsRequestProto.BarSize = &barSize
	}

	if !stringIsEmpty(whatToShow) {
		realTimeBarsRequestProto.WhatToShow = &whatToShow
	}

	if useRTH {
		realTimeBarsRequestProto.UseRTH = &useRTH
	}

	realTimeBarsOptionsMap := createStringStringMap(realTimeBarsOptionsList)
	if len(realTimeBarsOptionsMap) > 0 {
		realTimeBarsRequestProto.RealTimeBarsOptions = realTimeBarsOptionsMap
	}

	return realTimeBarsRequestProto
}

func createHeadTimestampRequestProto(reqID int64, contract *Contract, whatToShow string, useRTH bool, formatDate int) *protobuf.HeadTimestampRequest {
	headTimestampRequestProto := &protobuf.HeadTimestampRequest{}

	if isValidInt64Value(reqID) {
		id := int32(reqID)
		headTimestampRequestProto.ReqId = &id
	}

	order := &Order{}
	headTimestampRequestProto.Contract = createContractProto(contract, order)

	if !stringIsEmpty(whatToShow) {
		headTimestampRequestProto.WhatToShow = &whatToShow
	}

	if useRTH {
		headTimestampRequestProto.UseRTH = &useRTH
	}

	if isValidInt64Value(int64(formatDate)) {
		formatDate := int32(formatDate)
		headTimestampRequestProto.FormatDate = &formatDate
	}

	return headTimestampRequestProto
}

func createHistogramDataRequestProto(reqID int64, contract *Contract, useRTH bool, timePeriod string) *protobuf.HistogramDataRequest {
	histogramDataRequestProto := &protobuf.HistogramDataRequest{}

	if isValidInt64Value(reqID) {
		id := int32(reqID)
		histogramDataRequestProto.ReqId = &id
	}

	order := &Order{}
	histogramDataRequestProto.Contract = createContractProto(contract, order)

	if useRTH {
		histogramDataRequestProto.UseRTH = &useRTH
	}

	if !stringIsEmpty(timePeriod) {
		histogramDataRequestProto.TimePeriod = &timePeriod
	}

	return histogramDataRequestProto
}

func createHistoricalTicksRequestProto(reqID int64, contract *Contract, startDateTime string,
	endDateTime string, numberOfTicks int, whatToShow string, useRTH bool, ignoreSize bool, miscOptionsList []TagValue) *protobuf.HistoricalTicksRequest {

	historicalTicksRequestProto := &protobuf.HistoricalTicksRequest{}

	if isValidInt64Value(reqID) {
		id := int32(reqID)
		historicalTicksRequestProto.ReqId = &id
	}

	order := &Order{}
	historicalTicksRequestProto.Contract = createContractProto(contract, order)

	if !stringIsEmpty(startDateTime) {
		historicalTicksRequestProto.StartDateTime = &startDateTime
	}

	if !stringIsEmpty(endDateTime) {
		historicalTicksRequestProto.EndDateTime = &endDateTime
	}

	if isValidInt64Value(int64(numberOfTicks)) {
		numberOfTicks := int32(numberOfTicks)
		historicalTicksRequestProto.NumberOfTicks = &numberOfTicks
	}

	if !stringIsEmpty(whatToShow) {
		historicalTicksRequestProto.WhatToShow = &whatToShow
	}

	if useRTH {
		historicalTicksRequestProto.UseRTH = &useRTH
	}

	if ignoreSize {
		historicalTicksRequestProto.IgnoreSize = &ignoreSize
	}

	miscOptionsMap := createStringStringMap(miscOptionsList)
	if len(miscOptionsMap) > 0 {
		historicalTicksRequestProto.MiscOptions = miscOptionsMap
	}

	return historicalTicksRequestProto
}

func createTickByTickRequestProto(reqID int64, contract *Contract, tickType string, numberOfTicks int64, ignoreSize bool) *protobuf.TickByTickRequest {
	tickByTickRequestProto := &protobuf.TickByTickRequest{}

	if isValidInt64Value(reqID) {
		id := int32(reqID)
		tickByTickRequestProto.ReqId = &id
	}

	order := &Order{}
	tickByTickRequestProto.Contract = createContractProto(contract, order)

	if !stringIsEmpty(tickType) {
		tickByTickRequestProto.TickType = &tickType
	}

	if isValidInt64Value(numberOfTicks) {
		numberOfTicks := int32(numberOfTicks)
		tickByTickRequestProto.NumberOfTicks = &numberOfTicks
	}

	if ignoreSize {
		tickByTickRequestProto.IgnoreSize = &ignoreSize
	}

	return tickByTickRequestProto
}

func createCancelHistoricalDataProto(reqID int64) *protobuf.CancelHistoricalData {
	cancelHistoricalDataProto := &protobuf.CancelHistoricalData{}

	if isValidInt64Value(reqID) {
		id := int32(reqID)
		cancelHistoricalDataProto.ReqId = &id
	}

	return cancelHistoricalDataProto
}

func createCancelRealTimeBarsProto(reqID int64) *protobuf.CancelRealTimeBars {
	cancelRealTimeBarsProto := &protobuf.CancelRealTimeBars{}

	if isValidInt64Value(reqID) {
		id := int32(reqID)
		cancelRealTimeBarsProto.ReqId = &id
	}

	return cancelRealTimeBarsProto
}

func createCancelHeadTimestampProto(reqID int64) *protobuf.CancelHeadTimestamp {
	cancelHeadTimestampProto := &protobuf.CancelHeadTimestamp{}

	if isValidInt64Value(reqID) {
		id := int32(reqID)
		cancelHeadTimestampProto.ReqId = &id
	}

	return cancelHeadTimestampProto
}

func createCancelHistogramDataProto(reqID int64) *protobuf.CancelHistogramData {
	cancelHistogramDataProto := &protobuf.CancelHistogramData{}

	if isValidInt64Value(reqID) {
		id := int32(reqID)
		cancelHistogramDataProto.ReqId = &id
	}

	return cancelHistogramDataProto
}

func createCancelTickByTickProto(reqID int64) *protobuf.CancelTickByTick {
	cancelTickByTickProto := &protobuf.CancelTickByTick{}

	if isValidInt64Value(reqID) {
		id := int32(reqID)
		cancelTickByTickProto.ReqId = &id
	}

	return cancelTickByTickProto
}

func createNewsBulletinsRequestProto(allMessages bool) *protobuf.NewsBulletinsRequest {
	newsBulletinsRequestProto := &protobuf.NewsBulletinsRequest{}

	if allMessages {
		newsBulletinsRequestProto.AllMessages = &allMessages
	}

	return newsBulletinsRequestProto
}

func createCancelNewsBulletinsProto() *protobuf.CancelNewsBulletins {
	return &protobuf.CancelNewsBulletins{}
}

func createNewsArticleRequestProto(reqID int64, providerCode, articleID string, newsArticleOptionsList []TagValue) *protobuf.NewsArticleRequest {
	newsArticleRequestProto := &protobuf.NewsArticleRequest{}

	if isValidInt64Value(reqID) {
		id := int32(reqID)
		newsArticleRequestProto.ReqId = &id
	}

	if !stringIsEmpty(providerCode) {
		newsArticleRequestProto.ProviderCode = &providerCode
	}

	if !stringIsEmpty(articleID) {
		newsArticleRequestProto.ArticleId = &articleID
	}

	newsArticleOptionsMap := createStringStringMap(newsArticleOptionsList)
	if len(newsArticleOptionsMap) > 0 {
		newsArticleRequestProto.NewsArticleOptions = newsArticleOptionsMap
	}

	return newsArticleRequestProto
}

func createNewsProvidersRequestProto() *protobuf.NewsProvidersRequest {
	return &protobuf.NewsProvidersRequest{}
}

func createHistoricalNewsRequestProto(reqID int64, conID int64, providerCodes, startDateTime, endDateTime string, totalResults int64, historicalNewsOptionsList []TagValue) *protobuf.HistoricalNewsRequest {
	historicalNewsRequestProto := &protobuf.HistoricalNewsRequest{}

	if isValidInt64Value(reqID) {
		id := int32(reqID)
		historicalNewsRequestProto.ReqId = &id
	}

	if isValidInt64Value(conID) {
		cid := int32(conID)
		historicalNewsRequestProto.ConId = &cid
	}

	if !stringIsEmpty(providerCodes) {
		historicalNewsRequestProto.ProviderCodes = &providerCodes
	}

	if !stringIsEmpty(startDateTime) {
		historicalNewsRequestProto.StartDateTime = &startDateTime
	}

	if !stringIsEmpty(endDateTime) {
		historicalNewsRequestProto.EndDateTime = &endDateTime
	}

	if isValidInt64Value(totalResults) {
		tr := int32(totalResults)
		historicalNewsRequestProto.TotalResults = &tr
	}

	historicalNewsOptionsMap := createStringStringMap(historicalNewsOptionsList)
	if len(historicalNewsOptionsMap) > 0 {
		historicalNewsRequestProto.HistoricalNewsOptions = historicalNewsOptionsMap
	}

	return historicalNewsRequestProto
}

func createWshMetaDataRequestProto(reqID int64) *protobuf.WshMetaDataRequest {
	wshMetaDataRequestProto := &protobuf.WshMetaDataRequest{}

	if isValidInt64Value(reqID) {
		id := int32(reqID)
		wshMetaDataRequestProto.ReqId = &id
	}

	return wshMetaDataRequestProto
}

func createCancelWshMetaDataProto(reqID int64) *protobuf.CancelWshMetaData {
	cancelWshMetaDataProto := &protobuf.CancelWshMetaData{}

	if isValidInt64Value(reqID) {
		id := int32(reqID)
		cancelWshMetaDataProto.ReqId = &id
	}

	return cancelWshMetaDataProto
}

func createWshEventDataRequestProto(reqID int64, wshEventData *WshEventData) *protobuf.WshEventDataRequest {
	wshEventDataRequestProto := &protobuf.WshEventDataRequest{}

	if isValidInt64Value(reqID) {
		id := int32(reqID)
		wshEventDataRequestProto.ReqId = &id
	}

	if isValidInt64Value(int64(wshEventData.ConID)) {
		cid := int32(wshEventData.ConID)
		wshEventDataRequestProto.ConId = &cid
	}

	if !stringIsEmpty(wshEventData.Filter) {
		wshEventDataRequestProto.Filter = &wshEventData.Filter
	}

	if wshEventData.FillWatchList {
		wshEventDataRequestProto.FillWatchlist = &wshEventData.FillWatchList
	}

	if wshEventData.FillPortfolio {
		wshEventDataRequestProto.FillPortfolio = &wshEventData.FillPortfolio
	}

	if wshEventData.FillCompetitors {
		wshEventDataRequestProto.FillCompetitors = &wshEventData.FillCompetitors
	}

	if !stringIsEmpty(wshEventData.StartDate) {
		wshEventDataRequestProto.StartDate = &wshEventData.StartDate
	}

	if !stringIsEmpty(wshEventData.EndDate) {
		wshEventDataRequestProto.EndDate = &wshEventData.EndDate
	}

	if isValidInt64Value(int64(wshEventData.TotalLimit)) {
		total := int32(wshEventData.TotalLimit)
		wshEventDataRequestProto.TotalLimit = &total
	}

	return wshEventDataRequestProto
}

func createCancelWshEventDataProto(reqID int64) *protobuf.CancelWshEventData {
	cancelWshEventDataProto := &protobuf.CancelWshEventData{}

	if isValidInt64Value(reqID) {
		id := int32(reqID)
		cancelWshEventDataProto.ReqId = &id
	}

	return cancelWshEventDataProto
}

func createScannerParametersRequestProto() *protobuf.ScannerParametersRequest {
	return &protobuf.ScannerParametersRequest{}
}

func createScannerSubscriptionRequestProto(reqID int64, subscription *ScannerSubscription, scannerSubscriptionOptionsList []TagValue, scannerSubscriptionFilterOptionsList []TagValue) *protobuf.ScannerSubscriptionRequest {
	scannerSubscriptionRequestProto := &protobuf.ScannerSubscriptionRequest{}
	if isValidInt64Value(reqID) {
		id := int32(reqID)
		scannerSubscriptionRequestProto.ReqId = &id
	}
	scannerSubscriptionRequestProto.ScannerSubscription = createScannerSubscriptionProto(
		subscription, scannerSubscriptionOptionsList, scannerSubscriptionFilterOptionsList,
	)
	return scannerSubscriptionRequestProto
}

func createScannerSubscriptionProto(subscription *ScannerSubscription, scannerSubscriptionOptionsList []TagValue, scannerSubscriptionFilterOptionsList []TagValue) *protobuf.ScannerSubscription {
	proto := &protobuf.ScannerSubscription{}
	if isValidInt64Value(subscription.NumberOfRows) {
		n := int32(subscription.NumberOfRows)
		proto.NumberOfRows = &n
	}
	if !stringIsEmpty(subscription.Instrument) {
		proto.Instrument = &subscription.Instrument
	}
	if !stringIsEmpty(subscription.LocationCode) {
		proto.LocationCode = &subscription.LocationCode
	}
	if !stringIsEmpty(subscription.ScanCode) {
		proto.ScanCode = &subscription.ScanCode
	}
	if isValidFloat64Value(subscription.AbovePrice) {
		proto.AbovePrice = &subscription.AbovePrice
	}
	if isValidFloat64Value(subscription.BelowPrice) {
		proto.BelowPrice = &subscription.BelowPrice
	}
	if isValidInt64Value(subscription.AboveVolume) {
		proto.AboveVolume = &subscription.AboveVolume
	}
	if isValidInt64Value(subscription.AverageOptionVolumeAbove) {
		proto.AverageOptionVolumeAbove = &subscription.AverageOptionVolumeAbove
	}
	if isValidFloat64Value(subscription.MarketCapAbove) {
		proto.MarketCapAbove = &subscription.MarketCapAbove
	}
	if isValidFloat64Value(subscription.MarketCapBelow) {
		proto.MarketCapBelow = &subscription.MarketCapBelow
	}
	if !stringIsEmpty(subscription.MoodyRatingAbove) {
		proto.MoodyRatingAbove = &subscription.MoodyRatingAbove
	}
	if !stringIsEmpty(subscription.MoodyRatingBelow) {
		proto.MoodyRatingBelow = &subscription.MoodyRatingBelow
	}
	if !stringIsEmpty(subscription.SpRatingAbove) {
		proto.SpRatingAbove = &subscription.SpRatingAbove
	}
	if !stringIsEmpty(subscription.SpRatingBelow) {
		proto.SpRatingBelow = &subscription.SpRatingBelow
	}
	if !stringIsEmpty(subscription.MaturityDateAbove) {
		proto.MaturityDateAbove = &subscription.MaturityDateAbove
	}
	if !stringIsEmpty(subscription.MaturityDateBelow) {
		proto.MaturityDateBelow = &subscription.MaturityDateBelow
	}
	if isValidFloat64Value(subscription.CouponRateAbove) {
		proto.CouponRateAbove = &subscription.CouponRateAbove
	}
	if isValidFloat64Value(subscription.CouponRateBelow) {
		proto.CouponRateBelow = &subscription.CouponRateBelow
	}
	if subscription.ExcludeConvertible {
		proto.ExcludeConvertible = &subscription.ExcludeConvertible
	}
	if !stringIsEmpty(subscription.ScannerSettingPairs) {
		proto.ScannerSettingPairs = &subscription.ScannerSettingPairs
	}
	if !stringIsEmpty(subscription.StockTypeFilter) {
		proto.StockTypeFilter = &subscription.StockTypeFilter
	}
	opts := createStringStringMap(scannerSubscriptionOptionsList)
	if len(opts) > 0 {
		proto.ScannerSubscriptionOptions = opts
	}
	filt := createStringStringMap(scannerSubscriptionFilterOptionsList)
	if len(filt) > 0 {
		proto.ScannerSubscriptionFilterOptions = filt
	}
	return proto
}

func createFundamentalsDataRequestProto(reqID int64, contract *Contract, reportType string, fundamentalsDataOptionsList []TagValue) *protobuf.FundamentalsDataRequest {
	req := &protobuf.FundamentalsDataRequest{}
	if isValidInt64Value(reqID) {
		id := int32(reqID)
		req.ReqId = &id
	}
	req.Contract = createContractProto(contract, &Order{})
	if !stringIsEmpty(reportType) {
		req.ReportType = &reportType
	}
	opts := createStringStringMap(fundamentalsDataOptionsList)
	if len(opts) > 0 {
		req.FundamentalsDataOptions = opts
	}
	return req
}

func createCancelFundamentalsDataProto(reqID int64) *protobuf.CancelFundamentalsData {
	cancel := &protobuf.CancelFundamentalsData{}
	if isValidInt64Value(reqID) {
		id := int32(reqID)
		cancel.ReqId = &id
	}
	return cancel
}

func createPnLRequestProto(reqID int64, account, modelCode string) *protobuf.PnLRequest {
	req := &protobuf.PnLRequest{}
	if isValidInt64Value(reqID) {
		id := int32(reqID)
		req.ReqId = &id
	}
	if !stringIsEmpty(account) {
		req.Account = &account
	}
	if !stringIsEmpty(modelCode) {
		req.ModelCode = &modelCode
	}
	return req
}

func createCancelPnLProto(reqID int64) *protobuf.CancelPnL {
	cancel := &protobuf.CancelPnL{}
	if isValidInt64Value(reqID) {
		id := int32(reqID)
		cancel.ReqId = &id
	}
	return cancel
}

func createPnLSingleRequestProto(reqID int64, account, modelCode string, conID int64) *protobuf.PnLSingleRequest {
	req := &protobuf.PnLSingleRequest{}
	if isValidInt64Value(reqID) {
		id := int32(reqID)
		req.ReqId = &id
	}
	if !stringIsEmpty(account) {
		req.Account = &account
	}
	if !stringIsEmpty(modelCode) {
		req.ModelCode = &modelCode
	}
	if isValidInt64Value(conID) {
		cid := int32(conID)
		req.ConId = &cid
	}
	return req
}

func createCancelPnLSingleProto(reqID int64) *protobuf.CancelPnLSingle {
	cancel := &protobuf.CancelPnLSingle{}
	if isValidInt64Value(reqID) {
		id := int32(reqID)
		cancel.ReqId = &id
	}
	return cancel
}

func createCancelScannerSubscriptionProto(reqID int64) *protobuf.CancelScannerSubscription {
	cancel := &protobuf.CancelScannerSubscription{}
	if isValidInt64Value(reqID) {
		id := int32(reqID)
		cancel.ReqId = &id
	}
	return cancel
}

func createFARequestProto(faDataType int64) *protobuf.FARequest {
	req := &protobuf.FARequest{}
	if isValidInt64Value(faDataType) {
		t := int32(faDataType)
		req.FaDataType = &t
	}
	return req
}

func createFAReplaceProto(reqID int64, faDataType int64, xml string) *protobuf.FAReplace {
	req := &protobuf.FAReplace{}
	if isValidInt64Value(reqID) {
		id := int32(reqID)
		req.ReqId = &id
	}
	if isValidInt64Value(faDataType) {
		t := int32(faDataType)
		req.FaDataType = &t
	}
	if !stringIsEmpty(xml) {
		req.Xml = &xml
	}
	return req
}

func createExerciseOptionsRequestProto(orderID int64, contract *Contract, exerciseAction, exerciseQuantity int64, account string, override bool, manualOrderTime, customerAccount string, professionalCustomer bool) *protobuf.ExerciseOptionsRequest {
	req := &protobuf.ExerciseOptionsRequest{}
	if isValidInt64Value(orderID) {
		id := int32(orderID)
		req.OrderId = &id
	}
	req.Contract = createContractProto(contract, &Order{})

	if isValidInt64Value(exerciseAction) {
		action := int32(exerciseAction)
		req.ExerciseAction = &action
	}
	if isValidInt64Value(exerciseQuantity) {
		qty := int32(exerciseQuantity)
		req.ExerciseQuantity = &qty
	}
	if !stringIsEmpty(account) {
		req.Account = &account
	}
	if override {
		req.Override = &override
	}
	if !stringIsEmpty(manualOrderTime) {
		req.ManualOrderTime = &manualOrderTime
	}
	if !stringIsEmpty(customerAccount) {
		req.CustomerAccount = &customerAccount
	}
	if professionalCustomer {
		req.ProfessionalCustomer = &professionalCustomer
	}
	return req
}

func createCalculateImpliedVolatilityRequestProto(reqID int64, contract *Contract, optionPrice, underPrice float64, impliedVolatilityOptionsList []TagValue) *protobuf.CalculateImpliedVolatilityRequest {
	req := &protobuf.CalculateImpliedVolatilityRequest{}
	if isValidInt64Value(reqID) {
		id := int32(reqID)
		req.ReqId = &id
	}
	req.Contract = createContractProto(contract, &Order{})

	if isValidFloat64Value(optionPrice) {
		req.OptionPrice = &optionPrice
	}
	if isValidFloat64Value(underPrice) {
		req.UnderPrice = &underPrice
	}

	opts := createStringStringMap(impliedVolatilityOptionsList)
	if len(opts) > 0 {
		req.ImpliedVolatilityOptions = opts
	}
	return req
}

func createCancelCalculateImpliedVolatilityProto(reqID int64) *protobuf.CancelCalculateImpliedVolatility {
	cancel := &protobuf.CancelCalculateImpliedVolatility{}
	if isValidInt64Value(reqID) {
		id := int32(reqID)
		cancel.ReqId = &id
	}
	return cancel
}

func createCalculateOptionPriceRequestProto(reqID int64, contract *Contract, volatility, underPrice float64, optionPriceOptionsList []TagValue) *protobuf.CalculateOptionPriceRequest {
	req := &protobuf.CalculateOptionPriceRequest{}
	if isValidInt64Value(reqID) {
		id := int32(reqID)
		req.ReqId = &id
	}
	req.Contract = createContractProto(contract, &Order{})

	if isValidFloat64Value(volatility) {
		req.Volatility = &volatility
	}
	if isValidFloat64Value(underPrice) {
		req.UnderPrice = &underPrice
	}

	opts := createStringStringMap(optionPriceOptionsList)
	if len(opts) > 0 {
		req.OptionPriceOptions = opts
	}
	return req
}

func createCancelCalculateOptionPriceProto(reqID int64) *protobuf.CancelCalculateOptionPrice {
	cancel := &protobuf.CancelCalculateOptionPrice{}
	if isValidInt64Value(reqID) {
		id := int32(reqID)
		cancel.ReqId = &id
	}
	return cancel
}

func createSecDefOptParamsRequestProto(reqID int64, underlyingSymbol, futFopExchange, underlyingSecType string, underlyingConID int64) *protobuf.SecDefOptParamsRequest {
	req := &protobuf.SecDefOptParamsRequest{}
	if isValidInt64Value(reqID) {
		id := int32(reqID)
		req.ReqId = &id
	}
	if !stringIsEmpty(underlyingSymbol) {
		req.UnderlyingSymbol = &underlyingSymbol
	}
	if !stringIsEmpty(futFopExchange) {
		req.FutFopExchange = &futFopExchange
	}
	if !stringIsEmpty(underlyingSecType) {
		req.UnderlyingSecType = &underlyingSecType
	}
	if isValidInt64Value(underlyingConID) {
		cid := int32(underlyingConID)
		req.UnderlyingConId = &cid
	}
	return req
}

func createSoftDollarTiersRequestProto(reqID int64) *protobuf.SoftDollarTiersRequest {
	req := &protobuf.SoftDollarTiersRequest{}
	if isValidInt64Value(reqID) {
		id := int32(reqID)
		req.ReqId = &id
	}
	return req
}

func createFamilyCodesRequestProto() *protobuf.FamilyCodesRequest {
	return &protobuf.FamilyCodesRequest{}
}

func createMatchingSymbolsRequestProto(reqID int64, pattern string) *protobuf.MatchingSymbolsRequest {
	req := &protobuf.MatchingSymbolsRequest{}
	if isValidInt64Value(reqID) {
		id := int32(reqID)
		req.ReqId = &id
	}
	if !stringIsEmpty(pattern) {
		req.Pattern = &pattern
	}
	return req
}

func createSmartComponentsRequestProto(reqID int64, bboExchange string) *protobuf.SmartComponentsRequest {
	req := &protobuf.SmartComponentsRequest{}
	if isValidInt64Value(reqID) {
		id := int32(reqID)
		req.ReqId = &id
	}
	if !stringIsEmpty(bboExchange) {
		req.BboExchange = &bboExchange
	}
	return req
}

func createMarketRuleRequestProto(marketRuleID int64) *protobuf.MarketRuleRequest {
	req := &protobuf.MarketRuleRequest{}
	if isValidInt64Value(marketRuleID) {
		id := int32(marketRuleID)
		req.MarketRuleId = &id
	}
	return req
}

func createUserInfoRequestProto(reqID int64) *protobuf.UserInfoRequest {
	req := &protobuf.UserInfoRequest{}
	if isValidInt64Value(reqID) {
		id := int32(reqID)
		req.ReqId = &id
	}
	return req
}

func createIdsRequestProto(numIds int64) *protobuf.IdsRequest {
	req := &protobuf.IdsRequest{}
	if isValidInt64Value(numIds) {
		n := int32(numIds)
		req.NumIds = &n
	}
	return req
}

func createCurrentTimeRequestProto() *protobuf.CurrentTimeRequest {
	return &protobuf.CurrentTimeRequest{}
}

func createCurrentTimeInMillisRequestProto() *protobuf.CurrentTimeInMillisRequest {
	return &protobuf.CurrentTimeInMillisRequest{}
}

func createStartApiRequestProto(clientID int64, optionalCapabilities string) *protobuf.StartApiRequest {
	req := &protobuf.StartApiRequest{}
	if isValidInt64Value(clientID) {
		id := int32(clientID)
		req.ClientId = &id
	}
	if !stringIsEmpty(optionalCapabilities) {
		req.OptionalCapabilities = &optionalCapabilities
	}
	return req
}

func createSetServerLogLevelRequestProto(logLevel int64) *protobuf.SetServerLogLevelRequest {
	req := &protobuf.SetServerLogLevelRequest{}
	if isValidInt64Value(logLevel) {
		level := int32(logLevel)
		req.LogLevel = &level
	}
	return req
}

func createVerifyRequestProto(apiName, apiVersion string) *protobuf.VerifyRequest {
	req := &protobuf.VerifyRequest{}
	if !stringIsEmpty(apiName) {
		req.ApiName = &apiName
	}
	if !stringIsEmpty(apiVersion) {
		req.ApiVersion = &apiVersion
	}
	return req
}

func createVerifyMessageRequestProto(apiData string) *protobuf.VerifyMessageRequest {
	req := &protobuf.VerifyMessageRequest{}
	if !stringIsEmpty(apiData) {
		req.ApiData = &apiData
	}
	return req
}

func createQueryDisplayGroupsRequestProto(reqID int64) *protobuf.QueryDisplayGroupsRequest {
	req := &protobuf.QueryDisplayGroupsRequest{}
	if isValidInt64Value(reqID) {
		id := int32(reqID)
		req.ReqId = &id
	}
	return req
}

func createSubscribeToGroupEventsRequestProto(reqID int64, groupID int64) *protobuf.SubscribeToGroupEventsRequest {
	req := &protobuf.SubscribeToGroupEventsRequest{}
	if isValidInt64Value(reqID) {
		id := int32(reqID)
		req.ReqId = &id
	}
	if isValidInt64Value(groupID) {
		gid := int32(groupID)
		req.GroupId = &gid
	}
	return req
}

func createUpdateDisplayGroupRequestProto(reqID int64, contractInfo string) *protobuf.UpdateDisplayGroupRequest {
	req := &protobuf.UpdateDisplayGroupRequest{}
	if isValidInt64Value(reqID) {
		id := int32(reqID)
		req.ReqId = &id
	}
	if !stringIsEmpty(contractInfo) {
		req.ContractInfo = &contractInfo
	}
	return req
}

func createUnsubscribeFromGroupEventsRequestProto(reqID int64) *protobuf.UnsubscribeFromGroupEventsRequest {
	req := &protobuf.UnsubscribeFromGroupEventsRequest{}
	if isValidInt64Value(reqID) {
		id := int32(reqID)
		req.ReqId = &id
	}
	return req
}

func createMarketDepthExchangesRequestProto() *protobuf.MarketDepthExchangesRequest {
	return &protobuf.MarketDepthExchangesRequest{}
}

func createCancelContractDataProto(reqID int64) *protobuf.CancelContractData {
	req := &protobuf.CancelContractData{}
	if isValidInt64Value(reqID) {
		id := int32(reqID)
		req.ReqId = &id
	}
	return req
}

func createCancelHistoricalTicksProto(reqID int64) *protobuf.CancelHistoricalTicks {
	req := &protobuf.CancelHistoricalTicks{}
	if isValidInt64Value(reqID) {
		id := int32(reqID)
		req.ReqId = &id
	}
	return req
}
