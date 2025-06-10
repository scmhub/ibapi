package ibapi

import (
	"strconv"

	"github.com/scmhub/ibapi/protobuf"
)

func createExecutionFilterProto(execFilter *ExecutionFilter) *protobuf.ExecutionFilter {
	executionFilterProto := &protobuf.ExecutionFilter{}
	if isValidInt64Value(execFilter.ClientID) {
		clientID := int32(execFilter.ClientID)
		executionFilterProto.LastNDays = &clientID
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

func createPlaceOrderRequestProto(orderID OrderID, contract *Contract, order *Order) *protobuf.PlaceOrderRequest {
	placeOrderRequestProto := &protobuf.PlaceOrderRequest{}
	placeOrderRequestProto.Order = createOrderProto(order)
	if isValidInt64Value(orderID) {
		orderIDProto := int32(orderID)
		placeOrderRequestProto.OrderId = &orderIDProto
	}
	placeOrderRequestProto.Contract = createContractProto(contract, order)
	placeOrderRequestProto.Order = createOrderProto(order)

	return placeOrderRequestProto
}

func createOrderProto(order *Order) *protobuf.Order {
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
		orderProto.Conditions = createConditionsProto(order)
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
	// if order.AutoCancelParent {
	// 	orderProto.AutoCancelParent = &order.AutoCancelParent
	// }
	// if !stringIsEmpty(order.Shareholder) {
	// 	orderProto.Shareholder = &order.Shareholder
	// }
	// if order.ImbalanceOnly {
	// 	orderProto.ImbalanceOnly = &order.ImbalanceOnly
	// }
	// if order.RouteMarketableToBbo {
	// 	orderProto.RouteMarketableToBbo = &order.RouteMarketableToBbo
	// }
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

	return orderProto
}

func createConditionsProto(order *Order) []*protobuf.OrderCondition {
	var orderConditionList []*protobuf.OrderCondition
	for _, cond := range order.Conditions {
		var protoCond *protobuf.OrderCondition
		switch cond.Type() {
		case PriceOrderCondition:
			protoCond = createPriceConditionProto(cond)
		case TimeOrderCondition:
			protoCond = createTimeConditionProto(cond)
		case MarginOrderCondition:
			protoCond = createMarginConditionProto(cond)
		case ExecutionOrderCondition:
			protoCond = createExecutionConditionProto(cond)
		case VolumeOrderCondition:
			protoCond = createVolumeConditionProto(cond)
		case PercentChangeOrderCondition:
			protoCond = createPercentChangeConditionProto(cond)
		default:
			continue
		}
		orderConditionList = append(orderConditionList, protoCond)
	}
	return orderConditionList
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
		multiplier, _ := strconv.ParseFloat(contract.Multiplier, 64)
		contractProto.Multiplier = &multiplier
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
		contractProto.ComboLegs = createComboLegsProto(contract.ComboLegs, order.OrderComboLegs)
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
