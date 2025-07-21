package ibapi

import (
	"strconv"
	"strings"

	"github.com/scmhub/ibapi/protobuf"
)

func decodeContract(contractProto *protobuf.Contract) *Contract {
	contract := NewContract()
	if contractProto.ConId != nil {
		contract.ConID = int64(contractProto.GetConId())
	}
	if contractProto.Symbol != nil {
		contract.Symbol = contractProto.GetSymbol()
	}
	if contractProto.SecType != nil {
		contract.SecType = contractProto.GetSecType()
	}
	if contractProto.LastTradeDateOrContractMonth != nil {
		contract.LastTradeDateOrContractMonth = contractProto.GetLastTradeDateOrContractMonth()
	}
	if contractProto.Strike != nil {
		contract.Strike = contractProto.GetStrike()
	}
	if contractProto.Right != nil {
		contract.Right = contractProto.GetRight()
	}
	if contractProto.Multiplier != nil {
		contract.Multiplier = FloatMaxString(contractProto.GetMultiplier())
	}
	if contractProto.Exchange != nil {
		contract.Exchange = contractProto.GetExchange()
	}
	if contractProto.PrimaryExch != nil {
		contract.PrimaryExchange = contractProto.GetPrimaryExch()
	}
	if contractProto.Currency != nil {
		contract.Currency = contractProto.GetCurrency()
	}
	if contractProto.LocalSymbol != nil {
		contract.LocalSymbol = contractProto.GetLocalSymbol()
	}
	if contractProto.TradingClass != nil {
		contract.TradingClass = contractProto.GetTradingClass()
	}
	if contractProto.SecIdType != nil {
		contract.SecIDType = contractProto.GetSecIdType()
	}
	if contractProto.SecId != nil {
		contract.SecID = contractProto.GetSecId()
	}
	if contractProto.Description != nil {
		contract.Description = contractProto.GetDescription()
	}
	if contractProto.IssuerId != nil {
		contract.IssuerID = contractProto.GetIssuerId()
	}
	if contractProto.DeltaNeutralContract != nil {
		contract.DeltaNeutralContract = decodeDeltaNeutralContract(contractProto)
	}
	if contractProto.IncludeExpired != nil {
		contract.IncludeExpired = contractProto.GetIncludeExpired()
	}
	if contractProto.ComboLegsDescrip != nil {
		contract.ComboLegsDescrip = contractProto.GetComboLegsDescrip()
	}
	if contractProto.ComboLegs != nil {
		contract.ComboLegs = decodeComboLegs(contractProto)
	}
	if contractProto.LastTradeDate != nil {
		contract.LastTradeDate = contractProto.GetLastTradeDate()
	}

	return contract
}

func decodeComboLegs(contractProto *protobuf.Contract) []ComboLeg {
	var comboLegs []ComboLeg
	if len(contractProto.GetComboLegs()) > 0 {
		for _, comboLegProto := range contractProto.GetComboLegs() {
			comboLeg := ComboLeg{}
			if comboLegProto.ConId != nil {
				comboLeg.ConID = int64(comboLegProto.GetConId())
			}
			if comboLegProto.Ratio != nil {
				comboLeg.Ratio = int64(comboLegProto.GetRatio())
			}
			if comboLegProto.Action != nil {
				comboLeg.Action = comboLegProto.GetAction()
			}
			if comboLegProto.Exchange != nil {
				comboLeg.Exchange = comboLegProto.GetExchange()
			}
			if comboLegProto.OpenClose != nil {
				comboLeg.OpenClose = int64(comboLegProto.GetOpenClose())
			}
			if comboLegProto.ShortSalesSlot != nil {
				comboLeg.ShortSaleSlot = int64(comboLegProto.GetShortSalesSlot())
			}
			if comboLegProto.DesignatedLocation != nil {
				comboLeg.DesignatedLocation = comboLegProto.GetDesignatedLocation()
			}
			if comboLegProto.ExemptCode != nil {
				comboLeg.ExemptCode = int64(comboLegProto.GetExemptCode())
			}
			comboLegs = append(comboLegs, comboLeg)
		}
	}
	return comboLegs
}

func decodeOrderComboLegs(contractProto *protobuf.Contract) []OrderComboLeg {
	var orderComboLegs []OrderComboLeg
	if len(contractProto.GetComboLegs()) > 0 {
		for _, comboLegProto := range contractProto.GetComboLegs() {
			orderComboLeg := OrderComboLeg{}
			if comboLegProto.PerLegPrice != nil {
				orderComboLeg.Price = comboLegProto.GetPerLegPrice()
			}
			orderComboLegs = append(orderComboLegs, orderComboLeg)
		}
	}
	return orderComboLegs
}

func decodeDeltaNeutralContract(contractProto *protobuf.Contract) *DeltaNeutralContract {
	if contractProto.DeltaNeutralContract != nil {
		dncProto := contractProto.GetDeltaNeutralContract()
		dnc := &DeltaNeutralContract{}
		if dncProto.ConId != nil {
			dnc.ConID = int64(dncProto.GetConId())
		}
		if dncProto.Delta != nil {
			dnc.Delta = dncProto.GetDelta()
		}
		if dncProto.Price != nil {
			dnc.Price = dncProto.GetPrice()
		}
		return dnc
	}
	return nil
}

func decodeExecution(executionProto *protobuf.Execution) *Execution {
	execution := NewExecution()
	if executionProto.OrderId != nil {
		execution.OrderID = int64(executionProto.GetOrderId())
	}
	if executionProto.ClientId != nil {
		execution.ClientID = int64(executionProto.GetClientId())
	}
	if executionProto.ExecId != nil {
		execution.ExecID = executionProto.GetExecId()
	}
	if executionProto.Time != nil {
		execution.Time = executionProto.GetTime()
	}
	if executionProto.AcctNumber != nil {
		execution.AcctNumber = executionProto.GetAcctNumber()
	}
	if executionProto.Exchange != nil {
		execution.Exchange = executionProto.GetExchange()
	}
	if executionProto.Side != nil {
		execution.Side = executionProto.GetSide()
	}
	if executionProto.Shares != nil {
		execution.Shares = StringToDecimal(executionProto.GetShares())
	}
	if executionProto.Price != nil {
		execution.Price = executionProto.GetPrice()
	}
	if executionProto.PermId != nil {
		execution.PermID = int64(executionProto.GetPermId())
	}
	if executionProto.IsLiquidation != nil {
		execution.Liquidation = BoolToInt64(executionProto.GetIsLiquidation())
	}
	if executionProto.CumQty != nil {
		execution.CumQty = StringToDecimal(executionProto.GetCumQty())
	}
	if executionProto.AvgPrice != nil {
		execution.AvgPrice = executionProto.GetAvgPrice()
	}
	if executionProto.OrderRef != nil {
		execution.OrderRef = executionProto.GetOrderRef()
	}
	if executionProto.EvRule != nil {
		execution.EVRule = executionProto.GetEvRule()
	}
	if executionProto.EvMultiplier != nil {
		execution.EVMultiplier = executionProto.GetEvMultiplier()
	}
	if executionProto.ModelCode != nil {
		execution.ModelCode = executionProto.GetModelCode()
	}
	if executionProto.LastLiquidity != nil {
		execution.LastLiquidity = int64(executionProto.GetLastLiquidity())
	}
	if executionProto.IsPriceRevisionPending != nil {
		execution.PendingPriceRevision = executionProto.GetIsPriceRevisionPending()
	}
	if executionProto.Submitter != nil {
		execution.Submitter = executionProto.GetSubmitter()
	}
	if executionProto.OptExerciseOrLapseType != nil {
		execution.OptExerciseOrLapseType = OptionExerciseType(executionProto.GetOptExerciseOrLapseType())
	}
	return execution
}

func decodeOrder(orderID int64, contractProto *protobuf.Contract, orderProto *protobuf.Order) *Order {
	order := NewOrder()
	// order ids
	if orderProto.ClientId != nil {
		order.ClientID = int64(orderProto.GetClientId())
	}
	if isValidInt64Value(orderID) {
		order.OrderID = orderID
	}
	if orderProto.OrderId != nil {
		order.OrderID = int64(orderProto.GetOrderId())
	}
	if orderProto.PermId != nil {
		order.PermID = int64(orderProto.GetPermId())
	}
	if orderProto.ParentId != nil {
		order.ParentID = int64(orderProto.GetParentId())
	}
	// primary attributes
	if orderProto.Action != nil {
		order.Action = orderProto.GetAction()
	}
	if orderProto.TotalQuantity != nil {
		order.TotalQuantity = StringToDecimal(orderProto.GetTotalQuantity())
	}
	if orderProto.DisplaySize != nil {
		order.DisplaySize = int64(orderProto.GetDisplaySize())
	}
	if orderProto.OrderType != nil {
		order.OrderType = orderProto.GetOrderType()
	}
	if orderProto.LmtPrice != nil {
		order.LmtPrice = orderProto.GetLmtPrice()
	}
	if orderProto.AuxPrice != nil {
		order.AuxPrice = orderProto.GetAuxPrice()
	}
	if orderProto.Tif != nil {
		order.TIF = orderProto.GetTif()
	}
	// clearing info
	if orderProto.Account != nil {
		order.Account = orderProto.GetAccount()
	}
	if orderProto.SettlingFirm != nil {
		order.SettlingFirm = orderProto.GetSettlingFirm()
	}
	if orderProto.ClearingAccount != nil {
		order.ClearingAccount = orderProto.GetClearingAccount()
	}
	if orderProto.ClearingIntent != nil {
		order.ClearingIntent = orderProto.GetClearingIntent()
	}
	// secondary attributes
	if orderProto.AllOrNone != nil {
		order.AllOrNone = orderProto.GetAllOrNone()
	}
	if orderProto.BlockOrder != nil {
		order.BlockOrder = orderProto.GetBlockOrder()
	}
	if orderProto.Hidden != nil {
		order.Hidden = orderProto.GetHidden()
	}
	if orderProto.OutsideRth != nil {
		order.OutsideRTH = orderProto.GetOutsideRth()
	}
	if orderProto.SweepToFill != nil {
		order.SweepToFill = orderProto.GetSweepToFill()
	}
	if orderProto.PercentOffset != nil {
		order.PercentOffset = orderProto.GetPercentOffset()
	}
	if orderProto.TrailingPercent != nil {
		order.TrailingPercent = orderProto.GetTrailingPercent()
	}
	if orderProto.TrailStopPrice != nil {
		order.TrailStopPrice = orderProto.GetTrailStopPrice()
	}
	if orderProto.MinQty != nil {
		order.MinQty = int64(orderProto.GetMinQty())
	}
	if orderProto.GoodAfterTime != nil {
		order.GoodAfterTime = orderProto.GetGoodAfterTime()
	}
	if orderProto.GoodTillDate != nil {
		order.GoodTillDate = orderProto.GetGoodTillDate()
	}
	if orderProto.OcaGroup != nil {
		order.OCAGroup = orderProto.GetOcaGroup()
	}
	if orderProto.OrderRef != nil {
		order.OrderRef = orderProto.GetOrderRef()
	}
	if orderProto.Rule80A != nil {
		order.Rule80A = orderProto.GetRule80A()
	}
	if orderProto.OcaType != nil {
		order.OCAType = int64(orderProto.GetOcaType())
	}
	if orderProto.TriggerMethod != nil {
		order.TriggerMethod = int64(orderProto.GetTriggerMethod())
	}
	// extended order fields
	if orderProto.ActiveStartTime != nil {
		order.ActiveStartTime = orderProto.GetActiveStartTime()
	}
	if orderProto.ActiveStopTime != nil {
		order.ActiveStopTime = orderProto.GetActiveStopTime()
	}
	// advisor allocation orders
	if orderProto.FaGroup != nil {
		order.FAGroup = orderProto.GetFaGroup()
	}
	if orderProto.FaMethod != nil {
		order.FAMethod = orderProto.GetFaMethod()
	}
	if orderProto.FaPercentage != nil {
		order.FAPercentage = orderProto.GetFaPercentage()
	}
	// volatility orders
	if orderProto.Volatility != nil {
		order.Volatility = orderProto.GetVolatility()
	}
	if orderProto.VolatilityType != nil {
		order.VolatilityType = int64(orderProto.GetVolatilityType())
	}
	if orderProto.ContinuousUpdate != nil {
		order.ContinuousUpdate = orderProto.GetContinuousUpdate()
	}
	if orderProto.ReferencePriceType != nil {
		order.ReferencePriceType = int64(orderProto.GetReferencePriceType())
	}
	if orderProto.DeltaNeutralOrderType != nil {
		order.DeltaNeutralOrderType = orderProto.GetDeltaNeutralOrderType()
	}
	if orderProto.DeltaNeutralAuxPrice != nil {
		order.DeltaNeutralAuxPrice = orderProto.GetDeltaNeutralAuxPrice()
	}
	if orderProto.DeltaNeutralConId != nil {
		order.DeltaNeutralConID = int64(orderProto.GetDeltaNeutralConId())
	}
	if orderProto.DeltaNeutralOpenClose != nil {
		order.DeltaNeutralOpenClose = orderProto.GetDeltaNeutralOpenClose()
	}
	if orderProto.DeltaNeutralShortSale != nil {
		order.DeltaNeutralShortSale = orderProto.GetDeltaNeutralShortSale()
	}
	if orderProto.DeltaNeutralShortSaleSlot != nil {
		order.DeltaNeutralShortSaleSlot = int64(orderProto.GetDeltaNeutralShortSaleSlot())
	}
	if orderProto.DeltaNeutralDesignatedLocation != nil {
		order.DeltaNeutralDesignatedLocation = orderProto.GetDeltaNeutralDesignatedLocation()
	}
	// scale orders
	if orderProto.ScaleInitLevelSize != nil {
		order.ScaleInitLevelSize = int64(orderProto.GetScaleInitLevelSize())
	}
	if orderProto.ScaleSubsLevelSize != nil {
		order.ScaleSubsLevelSize = int64(orderProto.GetScaleSubsLevelSize())
	}
	if orderProto.ScalePriceIncrement != nil {
		order.ScalePriceIncrement = orderProto.GetScalePriceIncrement()
	}
	if orderProto.ScalePriceAdjustValue != nil {
		order.ScalePriceAdjustValue = orderProto.GetScalePriceAdjustValue()
	}
	if orderProto.ScalePriceAdjustInterval != nil {
		order.ScalePriceAdjustInterval = int64(orderProto.GetScalePriceAdjustInterval())
	}
	if orderProto.ScaleProfitOffset != nil {
		order.ScaleProfitOffset = orderProto.GetScaleProfitOffset()
	}
	if orderProto.ScaleAutoReset != nil {
		order.ScaleAutoReset = orderProto.GetScaleAutoReset()
	}
	if orderProto.ScaleInitPosition != nil {
		order.ScaleInitPosition = int64(orderProto.GetScaleInitPosition())
	}
	if orderProto.ScaleInitFillQty != nil {
		order.ScaleInitFillQty = int64(orderProto.GetScaleInitFillQty())
	}
	if orderProto.ScaleRandomPercent != nil {
		order.ScaleRandomPercent = orderProto.GetScaleRandomPercent()
	}
	if orderProto.ScaleTable != nil {
		order.ScaleTable = orderProto.GetScaleTable()
	}
	// hedge orders
	if orderProto.HedgeType != nil {
		order.HedgeType = orderProto.GetHedgeType()
	}
	if orderProto.HedgeParam != nil {
		order.HedgeParam = orderProto.GetHedgeParam()
	}
	// algo orders
	if orderProto.AlgoStrategy != nil {
		order.AlgoStrategy = orderProto.GetAlgoStrategy()
		order.AlgoParams = decodeTagValueList(orderProto.AlgoParams)
	}
	if orderProto.AlgoId != nil {
		order.AlgoID = orderProto.GetAlgoId()
	}
	// combo orders
	order.OrderComboLegs = decodeOrderComboLegs(contractProto)
	order.SmartComboRoutingParams = decodeTagValueList(orderProto.SmartComboRoutingParams)
	// processing control
	if orderProto.WhatIf != nil {
		order.WhatIf = orderProto.GetWhatIf()
	}
	if orderProto.Transmit != nil {
		order.Transmit = orderProto.GetTransmit()
	}
	if orderProto.OverridePercentageConstraints != nil {
		order.OverridePercentageConstraints = orderProto.GetOverridePercentageConstraints()
	}
	// Institutional orders only
	if orderProto.OpenClose != nil {
		order.OpenClose = orderProto.GetOpenClose()
	}
	if orderProto.Origin != nil {
		order.Origin = int64(orderProto.GetOrigin())
	}
	if orderProto.ShortSaleSlot != nil {
		order.ShortSaleSlot = int64(orderProto.GetShortSaleSlot())
	}
	if orderProto.DesignatedLocation != nil {
		order.DesignatedLocation = orderProto.GetDesignatedLocation()
	}
	if orderProto.ExemptCode != nil {
		order.ExemptCode = int64(orderProto.GetExemptCode())
	}
	if orderProto.DeltaNeutralSettlingFirm != nil {
		order.DeltaNeutralSettlingFirm = orderProto.GetDeltaNeutralSettlingFirm()
	}
	if orderProto.DeltaNeutralClearingAccount != nil {
		order.DeltaNeutralClearingAccount = orderProto.GetDeltaNeutralClearingAccount()
	}
	if orderProto.DeltaNeutralClearingIntent != nil {
		order.DeltaNeutralClearingIntent = orderProto.GetDeltaNeutralClearingIntent()
	}
	// SMART routing only
	if orderProto.DiscretionaryAmt != nil {
		order.DiscretionaryAmt = orderProto.GetDiscretionaryAmt()
	}
	if orderProto.OptOutSmartRouting != nil {
		order.OptOutSmartRouting = orderProto.GetOptOutSmartRouting()
	}
	// BOX ORDERS ONLY
	if orderProto.StartingPrice != nil {
		order.StartingPrice = orderProto.GetStartingPrice()
	}
	if orderProto.StockRefPrice != nil {
		order.StockRefPrice = orderProto.GetStockRefPrice()
	}
	if orderProto.Delta != nil {
		order.Delta = orderProto.GetDelta()
	}
	// pegged to stock or VOL orders
	if orderProto.StockRangeLower != nil {
		order.StockRangeLower = orderProto.GetStockRangeLower()
	}
	if orderProto.StockRangeUpper != nil {
		order.StockRangeUpper = orderProto.GetStockRangeUpper()
	}
	// Not Held
	if orderProto.NotHeld != nil {
		order.NotHeld = orderProto.GetNotHeld()
	}
	//order algo id
	if orderProto.Solicited != nil {
		order.Solicited = orderProto.GetSolicited()
	}
	if orderProto.RandomizeSize != nil {
		order.RandomizeSize = orderProto.GetRandomizeSize()
	}
	if orderProto.RandomizePrice != nil {
		order.RandomizePrice = orderProto.GetRandomizePrice()
	}
	// PEG2BENCH fields
	if orderProto.ReferenceContractId != nil {
		order.ReferenceContractID = int64(orderProto.GetReferenceContractId())
	}
	if orderProto.PeggedChangeAmount != nil {
		order.PeggedChangeAmount = orderProto.GetPeggedChangeAmount()
	}
	if orderProto.IsPeggedChangeAmountDecrease != nil {
		order.IsPeggedChangeAmountDecrease = orderProto.GetIsPeggedChangeAmountDecrease()
	}
	if orderProto.ReferenceChangeAmount != nil {
		order.ReferenceChangeAmount = orderProto.GetReferenceChangeAmount()
	}
	if orderProto.ReferenceExchangeId != nil {
		order.ReferenceExchangeID = orderProto.GetReferenceExchangeId()
	}
	if orderProto.AdjustedOrderType != nil {
		order.AdjustedOrderType = orderProto.GetAdjustedOrderType()
	}
	if orderProto.TriggerPrice != nil {
		order.TriggerPrice = orderProto.GetTriggerPrice()
	}
	if orderProto.AdjustedStopPrice != nil {
		order.AdjustedStopPrice = orderProto.GetAdjustedStopPrice()
	}
	if orderProto.AdjustedStopLimitPrice != nil {
		order.AdjustedStopLimitPrice = orderProto.GetAdjustedStopLimitPrice()
	}
	if orderProto.AdjustedTrailingAmount != nil {
		order.AdjustedTrailingAmount = orderProto.GetAdjustedTrailingAmount()
	}
	if orderProto.AdjustableTrailingUnit != nil {
		order.AdjustableTrailingUnit = int64(orderProto.GetAdjustableTrailingUnit())
	}
	if orderProto.LmtPriceOffset != nil {
		order.LmtPriceOffset = orderProto.GetLmtPriceOffset()
	}
	order.Conditions = decodeConditions(orderProto)
	if orderProto.ConditionsCancelOrder != nil {
		order.ConditionsCancelOrder = orderProto.GetConditionsCancelOrder()
	}
	if orderProto.ConditionsIgnoreRth != nil {
		order.ConditionsIgnoreRth = orderProto.GetConditionsIgnoreRth()
	}
	// models
	if orderProto.ModelCode != nil {
		order.ModelCode = orderProto.GetModelCode()
	}
	if orderProto.ExtOperator != nil {
		order.ExtOperator = orderProto.GetExtOperator()
	}
	order.SoftDollarTier = decodeSoftDollarTier(orderProto)
	// native cash quantity
	if orderProto.CashQty != nil {
		order.CashQty = orderProto.GetCashQty()
	}
	if orderProto.Mifid2DecisionMaker != nil {
		order.Mifid2DecisionMaker = orderProto.GetMifid2DecisionMaker()
	}
	if orderProto.Mifid2DecisionAlgo != nil {
		order.Mifid2DecisionAlgo = orderProto.GetMifid2DecisionAlgo()
	}
	if orderProto.Mifid2ExecutionTrader != nil {
		order.Mifid2ExecutionTrader = orderProto.GetMifid2ExecutionTrader()
	}
	if orderProto.Mifid2ExecutionAlgo != nil {
		order.Mifid2ExecutionAlgo = orderProto.GetMifid2ExecutionAlgo()
	}
	// don't use auto price for hedge
	if orderProto.DontUseAutoPriceForHedge != nil {
		order.DontUseAutoPriceForHedge = orderProto.GetDontUseAutoPriceForHedge()
	}
	if orderProto.IsOmsContainer != nil {
		order.IsOmsContainer = orderProto.GetIsOmsContainer()
	}
	if orderProto.DiscretionaryUpToLimitPrice != nil {
		order.DiscretionaryUpToLimitPrice = orderProto.GetDiscretionaryUpToLimitPrice()
	}
	if orderProto.AutoCancelDate != nil {
		order.AutoCancelDate = orderProto.GetAutoCancelDate()
	}
	if orderProto.FilledQuantity != nil {
		order.FilledQuantity = StringToDecimal(orderProto.GetFilledQuantity())
	}
	if orderProto.RefFuturesConId != nil {
		order.RefFuturesConID = int64(orderProto.GetRefFuturesConId())
	}
	if orderProto.AutoCancelParent != nil {
		order.AutoCancelParent = orderProto.GetAutoCancelParent()
	}
	if orderProto.Shareholder != nil {
		order.Shareholder = orderProto.GetShareholder()
	}
	if orderProto.ImbalanceOnly != nil {
		order.ImbalanceOnly = orderProto.GetImbalanceOnly()
	}
	if orderProto.RouteMarketableToBbo != nil {
		order.RouteMarketableToBbo = orderProto.GetRouteMarketableToBbo()
	}
	if orderProto.ParentPermId != nil {
		order.ParentPermID = orderProto.GetParentPermId()
	}
	if orderProto.UsePriceMgmtAlgo != nil {
		order.UsePriceMgmtAlgo = int64(orderProto.GetUsePriceMgmtAlgo())
	}
	if orderProto.Duration != nil {
		order.Duration = int64(orderProto.GetDuration())
	}
	if orderProto.PostToAts != nil {
		order.PostToAts = int64(orderProto.GetPostToAts())
	}
	if orderProto.AdvancedErrorOverride != nil {
		order.AdvancedErrorOverride = orderProto.GetAdvancedErrorOverride()
	}
	if orderProto.ManualOrderTime != nil {
		order.ManualOrderTime = orderProto.GetManualOrderTime()
	}
	if orderProto.MinTradeQty != nil {
		order.MinTradeQty = int64(orderProto.GetMinTradeQty())
	}
	if orderProto.MinCompeteSize != nil {
		order.MinCompeteSize = int64(orderProto.GetMinCompeteSize())
	}
	if orderProto.CompeteAgainstBestOffset != nil {
		order.CompeteAgainstBestOffset = orderProto.GetCompeteAgainstBestOffset()
	}
	if orderProto.MidOffsetAtWhole != nil {
		order.MidOffsetAtWhole = orderProto.GetMidOffsetAtWhole()
	}
	if orderProto.MidOffsetAtHalf != nil {
		order.MidOffsetAtHalf = orderProto.GetMidOffsetAtHalf()
	}
	if orderProto.CustomerAccount != nil {
		order.CustomerAccount = orderProto.GetCustomerAccount()
	}
	if orderProto.ProfessionalCustomer != nil {
		order.ProfessionalCustomer = orderProto.GetProfessionalCustomer()
	}
	if orderProto.BondAccruedInterest != nil {
		order.BondAccruedInterest = orderProto.GetBondAccruedInterest()
	}
	if orderProto.IncludeOvernight != nil {
		order.IncludeOvernight = orderProto.GetIncludeOvernight()
	}
	if orderProto.ManualOrderIndicator != nil {
		order.ManualOrderIndicator = int64(orderProto.GetManualOrderIndicator())
	}
	if orderProto.Submitter != nil {
		order.Submitter = orderProto.GetSubmitter()
	}

	return order
}

func decodeConditions(orderProto *protobuf.Order) []OrderCondition {
	var orderConditions []OrderCondition
	for _, condProto := range orderProto.GetConditions() {
		var cond OrderCondition
		conditionType := OrderConditionType(condProto.GetType())
		switch conditionType {
		case PriceOrderCondition:
			cond = decodePriceCondition(condProto)
		case TimeOrderCondition:
			cond = decodeTimeCondition(condProto)
		case MarginOrderCondition:
			cond = decodeMarginCondition(condProto)
		case ExecutionOrderCondition:
			cond = decodeExecutionCondition(condProto)
		case VolumeOrderCondition:
			cond = decodeVolumeCondition(condProto)
		case PercentChangeOrderCondition:
			cond = decodePercentChangeCondition(condProto)
		}
		if cond != nil {
			orderConditions = append(orderConditions, cond)
		}
	}
	return orderConditions
}

// --- Condition field helpers ---

func setConditionFields(condProto *protobuf.OrderCondition, cond OrderCondition) {
	if condProto.IsConjunctionConnection != nil {
		cond.SetIsConjunctionConnection(condProto.GetIsConjunctionConnection())
	}
}

func setOperatorConditionFields(condProto *protobuf.OrderCondition, cond *operatorCondition) {
	setConditionFields(condProto, cond)
	if condProto.IsMore != nil {
		cond.IsMore = condProto.GetIsMore()
	}
}

func setContractConditionFields(condProto *protobuf.OrderCondition, cond *contractCondition) {
	setOperatorConditionFields(condProto, cond.operatorCondition)
	if condProto.ConId != nil {
		cond.ConID = int64(condProto.GetConId())
	}
	if condProto.Exchange != nil {
		cond.Exchange = condProto.GetExchange()
	}
}

// --- Concrete condition decoders ---

func decodePriceCondition(condProto *protobuf.OrderCondition) *PriceCondition {
	cond := &PriceCondition{}
	setContractConditionFields(condProto, cond.contractCondition)
	if condProto.Price != nil {
		cond.Price = condProto.GetPrice()
	}
	if condProto.TriggerMethod != nil {
		cond.TriggerMethod = TriggerMethod(condProto.GetTriggerMethod())
	}
	return cond
}

func decodeTimeCondition(condProto *protobuf.OrderCondition) *TimeCondition {
	cond := &TimeCondition{}
	setOperatorConditionFields(condProto, cond.operatorCondition)
	if condProto.Time != nil {
		cond.Time = condProto.GetTime()
	}
	return cond
}

func decodeMarginCondition(condProto *protobuf.OrderCondition) *MarginCondition {
	cond := &MarginCondition{}
	setOperatorConditionFields(condProto, cond.operatorCondition)
	if condProto.Percent != nil {
		cond.Percent = int64(condProto.GetPercent())
	}
	return cond
}

func decodeExecutionCondition(condProto *protobuf.OrderCondition) *ExecutionCondition {
	cond := &ExecutionCondition{}
	setConditionFields(condProto, cond)
	if condProto.SecType != nil {
		cond.SecType = condProto.GetSecType()
	}
	if condProto.Exchange != nil {
		cond.Exchange = condProto.GetExchange()
	}
	if condProto.Symbol != nil {
		cond.Symbol = condProto.GetSymbol()
	}
	return cond
}

func decodeVolumeCondition(condProto *protobuf.OrderCondition) *VolumeCondition {
	cond := &VolumeCondition{}
	setContractConditionFields(condProto, cond.contractCondition)
	if condProto.Volume != nil {
		cond.Volume = int64(condProto.GetVolume())
	}
	return cond
}

func decodePercentChangeCondition(condProto *protobuf.OrderCondition) *PercentChangeCondition {
	cond := &PercentChangeCondition{}
	setContractConditionFields(condProto, cond.contractCondition)
	if condProto.ChangePercent != nil {
		cond.ChangePercent = condProto.GetChangePercent()
	}
	return cond
}

func decodeSoftDollarTier(orderProto *protobuf.Order) SoftDollarTier {
	var softDollarTier SoftDollarTier
	if orderProto.SoftDollarTier != nil {
		tierProto := orderProto.GetSoftDollarTier()
		var name string
		if tierProto.Name != nil {
			name = tierProto.GetName()
		}
		var value string
		if tierProto.Value != nil {
			value = tierProto.GetValue()
		}
		var displayName string
		if tierProto.DisplayName != nil {
			displayName = tierProto.GetDisplayName()
		}
		softDollarTier = SoftDollarTier{
			Name:        name,
			Value:       value,
			DisplayName: displayName,
		}
	}
	return softDollarTier
}

func decodeTagValueList(stringStringMap map[string]string) []TagValue {
	var params []TagValue
	for k, v := range stringStringMap {
		params = append(params, TagValue{Tag: k, Value: v})
	}
	return params
}

func decodeOrderState(orderStateProto *protobuf.OrderState) *OrderState {
	orderState := &OrderState{}
	if orderStateProto.Status != nil {
		orderState.Status = orderStateProto.GetStatus()
	}
	if orderStateProto.InitMarginBefore != nil {
		orderState.InitMarginBefore = FloatToString(orderStateProto.GetInitMarginBefore())
	}
	if orderStateProto.MaintMarginBefore != nil {
		orderState.MaintMarginBefore = FloatToString(orderStateProto.GetMaintMarginBefore())
	}
	if orderStateProto.EquityWithLoanBefore != nil {
		orderState.EquityWithLoanBefore = FloatToString(orderStateProto.GetEquityWithLoanBefore())
	}
	if orderStateProto.InitMarginChange != nil {
		orderState.InitMarginChange = FloatToString(orderStateProto.GetInitMarginChange())
	}
	if orderStateProto.MaintMarginChange != nil {
		orderState.MaintMarginChange = FloatToString(orderStateProto.GetMaintMarginChange())
	}
	if orderStateProto.EquityWithLoanChange != nil {
		orderState.EquityWithLoanChange = FloatToString(orderStateProto.GetEquityWithLoanChange())
	}
	if orderStateProto.InitMarginAfter != nil {
		orderState.InitMarginAfter = FloatToString(orderStateProto.GetInitMarginAfter())
	}
	if orderStateProto.MaintMarginAfter != nil {
		orderState.MaintMarginAfter = FloatToString(orderStateProto.GetMaintMarginAfter())
	}
	if orderStateProto.EquityWithLoanAfter != nil {
		orderState.EquityWithLoanAfter = FloatToString(orderStateProto.GetEquityWithLoanAfter())
	}
	if orderStateProto.CommissionAndFees != nil {
		orderState.CommissionAndFees = orderStateProto.GetCommissionAndFees()
	}
	if orderStateProto.MinCommissionAndFees != nil {
		orderState.MinCommissionAndFees = orderStateProto.GetMinCommissionAndFees()
	}
	if orderStateProto.MaxCommissionAndFees != nil {
		orderState.MaxCommissionAndFees = orderStateProto.GetMaxCommissionAndFees()
	}
	if orderStateProto.CommissionAndFeesCurrency != nil {
		orderState.CommissionAndFeesCurrency = orderStateProto.GetCommissionAndFeesCurrency()
	}
	if orderStateProto.MarginCurrency != nil {
		orderState.MarginCurrency = orderStateProto.GetMarginCurrency()
	}
	if orderStateProto.InitMarginBeforeOutsideRTH != nil {
		orderState.InitMarginBeforeOutsideRTH = orderStateProto.GetInitMarginBeforeOutsideRTH()
	}
	if orderStateProto.MaintMarginBeforeOutsideRTH != nil {
		orderState.MaintMarginBeforeOutsideRTH = orderStateProto.GetMaintMarginBeforeOutsideRTH()
	}
	if orderStateProto.EquityWithLoanBeforeOutsideRTH != nil {
		orderState.EquityWithLoanBeforeOutsideRTH = orderStateProto.GetEquityWithLoanBeforeOutsideRTH()
	}
	if orderStateProto.InitMarginChangeOutsideRTH != nil {
		orderState.InitMarginChangeOutsideRTH = orderStateProto.GetInitMarginChangeOutsideRTH()
	}
	if orderStateProto.MaintMarginChangeOutsideRTH != nil {
		orderState.MaintMarginChangeOutsideRTH = orderStateProto.GetMaintMarginChangeOutsideRTH()
	}
	if orderStateProto.EquityWithLoanChangeOutsideRTH != nil {
		orderState.EquityWithLoanChangeOutsideRTH = orderStateProto.GetEquityWithLoanChangeOutsideRTH()
	}
	if orderStateProto.InitMarginAfterOutsideRTH != nil {
		orderState.InitMarginAfterOutsideRTH = orderStateProto.GetInitMarginAfterOutsideRTH()
	}
	if orderStateProto.MaintMarginAfterOutsideRTH != nil {
		orderState.MaintMarginAfterOutsideRTH = orderStateProto.GetMaintMarginAfterOutsideRTH()
	}
	if orderStateProto.EquityWithLoanAfterOutsideRTH != nil {
		orderState.EquityWithLoanAfterOutsideRTH = orderStateProto.GetEquityWithLoanAfterOutsideRTH()
	}
	if orderStateProto.SuggestedSize != nil {
		orderState.SuggestedSize = StringToDecimal(orderStateProto.GetSuggestedSize())
	}
	if orderStateProto.RejectReason != nil {
		orderState.RejectReason = orderStateProto.GetRejectReason()
	}
	orderState.OrderAllocations = decodeOrderAllocations(orderStateProto)
	if orderStateProto.WarningText != nil {
		orderState.WarningText = orderStateProto.GetWarningText()
	}
	if orderStateProto.CompletedTime != nil {
		orderState.CompletedTime = orderStateProto.GetCompletedTime()
	}
	if orderStateProto.CompletedStatus != nil {
		orderState.CompletedStatus = orderStateProto.GetCompletedStatus()
	}
	return orderState
}

func decodeOrderAllocations(orderStateProto *protobuf.OrderState) []*OrderAllocation {
	var orderAllocations []*OrderAllocation
	for _, allocProto := range orderStateProto.GetOrderAllocations() {
		orderAllocation := NewOrderAllocation()
		if allocProto.Account != nil {
			orderAllocation.Account = allocProto.GetAccount()
		}
		if allocProto.Position != nil {
			orderAllocation.Position = StringToDecimal(allocProto.GetPosition())
		}
		if allocProto.PositionDesired != nil {
			orderAllocation.PositionDesired = StringToDecimal(allocProto.GetPositionDesired())
		}
		if allocProto.PositionAfter != nil {
			orderAllocation.PositionAfter = StringToDecimal(allocProto.GetPositionAfter())
		}
		if allocProto.DesiredAllocQty != nil {
			orderAllocation.DesiredAllocQty = StringToDecimal(allocProto.GetDesiredAllocQty())
		}
		if allocProto.AllowedAllocQty != nil {
			orderAllocation.AllowedAllocQty = StringToDecimal(allocProto.GetAllowedAllocQty())
		}
		if allocProto.IsMonetary != nil {
			orderAllocation.IsMonetary = allocProto.GetIsMonetary()
		}
		orderAllocations = append(orderAllocations, orderAllocation)
	}
	return orderAllocations
}

func decodeContractDetails(contractProto *protobuf.Contract, contractDetailsProto *protobuf.ContractDetails, isBond bool) *ContractDetails {
	contractDetails := &ContractDetails{}
	contractDetails.Contract = *decodeContract(contractProto)

	if contractDetailsProto.MarketName != nil {
		contractDetails.MarketName = contractDetailsProto.GetMarketName()
	}
	if contractDetailsProto.MinTick != nil {
		contractDetails.MinTick, _ = strconv.ParseFloat(contractDetailsProto.GetMinTick(), 64)
	}
	if contractDetailsProto.OrderTypes != nil {
		contractDetails.OrderTypes = contractDetailsProto.GetOrderTypes()
	}
	if contractDetailsProto.ValidExchanges != nil {
		contractDetails.ValidExchanges = contractDetailsProto.GetValidExchanges()
	}
	if contractDetailsProto.PriceMagnifier != nil {
		contractDetails.PriceMagnifier = int64(contractDetailsProto.GetPriceMagnifier())
	}
	if contractDetailsProto.UnderConId != nil {
		contractDetails.UnderConID = int64(contractDetailsProto.GetUnderConId())
	}
	if contractDetailsProto.LongName != nil {
		contractDetails.LongName = contractDetailsProto.GetLongName()
	}
	if contractDetailsProto.ContractMonth != nil {
		contractDetails.ContractMonth = contractDetailsProto.GetContractMonth()
	}
	if contractDetailsProto.Industry != nil {
		contractDetails.Industry = contractDetailsProto.GetIndustry()
	}
	if contractDetailsProto.Category != nil {
		contractDetails.Category = contractDetailsProto.GetCategory()
	}
	if contractDetailsProto.Subcategory != nil {
		contractDetails.Subcategory = contractDetailsProto.GetSubcategory()
	}
	if contractDetailsProto.TimeZoneId != nil {
		contractDetails.TimeZoneID = contractDetailsProto.GetTimeZoneId()
	}
	if contractDetailsProto.TradingHours != nil {
		contractDetails.TradingHours = contractDetailsProto.GetTradingHours()
	}
	if contractDetailsProto.LiquidHours != nil {
		contractDetails.LiquidHours = contractDetailsProto.GetLiquidHours()
	}
	if contractDetailsProto.EvRule != nil {
		contractDetails.EVRule = contractDetailsProto.GetEvRule()
	}
	if contractDetailsProto.EvMultiplier != nil {
		contractDetails.EVMultiplier = int64(contractDetailsProto.GetEvMultiplier())
	}

	contractDetails.SecIDList = decodeTagValueList(contractDetailsProto.GetSecIdList())

	if contractDetailsProto.AggGroup != nil {
		contractDetails.AggGroup = int64(contractDetailsProto.GetAggGroup())
	}
	if contractDetailsProto.UnderSymbol != nil {
		contractDetails.UnderSymbol = contractDetailsProto.GetUnderSymbol()
	}
	if contractDetailsProto.UnderSecType != nil {
		contractDetails.UnderSecType = contractDetailsProto.GetUnderSecType()
	}
	if contractDetailsProto.MarketRuleIds != nil {
		contractDetails.MarketRuleIDs = contractDetailsProto.GetMarketRuleIds()
	}
	if contractDetailsProto.RealExpirationDate != nil {
		contractDetails.RealExpirationDate = contractDetailsProto.GetRealExpirationDate()
	}
	if contractDetailsProto.StockType != nil {
		contractDetails.StockType = contractDetailsProto.GetStockType()
	}
	if contractDetailsProto.MinSize != nil {
		contractDetails.MinSize = StringToDecimal(contractDetailsProto.GetMinSize())
	}
	if contractDetailsProto.SizeIncrement != nil {
		contractDetails.SizeIncrement = StringToDecimal(contractDetailsProto.GetSizeIncrement())
	}
	if contractDetailsProto.SuggestedSizeIncrement != nil {
		contractDetails.SuggestedSizeIncrement = StringToDecimal(contractDetailsProto.GetSuggestedSizeIncrement())
	}

	setLastTradeDate(contractDetails.Contract.LastTradeDateOrContractMonth, contractDetails, isBond)

	// fund	fields
	if contractDetailsProto.FundName != nil {
		contractDetails.FundName = contractDetailsProto.GetFundName()
	}
	if contractDetailsProto.FundFamily != nil {
		contractDetails.FundFamily = contractDetailsProto.GetFundFamily()
	}
	if contractDetailsProto.FundType != nil {
		contractDetails.FundType = contractDetailsProto.GetFundType()
	}
	if contractDetailsProto.FundFrontLoad != nil {
		contractDetails.FundFrontLoad = contractDetailsProto.GetFundFrontLoad()
	}
	if contractDetailsProto.FundBackLoad != nil {
		contractDetails.FundBackLoad = contractDetailsProto.GetFundBackLoad()
	}
	if contractDetailsProto.FundBackLoadTimeInterval != nil {
		contractDetails.FundBackLoadTimeInterval = contractDetailsProto.GetFundBackLoadTimeInterval()
	}
	if contractDetailsProto.FundManagementFee != nil {
		contractDetails.FundManagementFee = contractDetailsProto.GetFundManagementFee()
	}
	if contractDetailsProto.FundClosed != nil {
		contractDetails.FundClosed = contractDetailsProto.GetFundClosed()
	}
	if contractDetailsProto.FundClosedForNewInvestors != nil {
		contractDetails.FundClosedForNewInvestors = contractDetailsProto.GetFundClosedForNewInvestors()
	}
	if contractDetailsProto.FundClosedForNewMoney != nil {
		contractDetails.FundClosedForNewMoney = contractDetailsProto.GetFundClosedForNewMoney()
	}
	if contractDetailsProto.FundNotifyAmount != nil {
		contractDetails.FundNotifyAmount = contractDetailsProto.GetFundNotifyAmount()
	}
	if contractDetailsProto.FundMinimumInitialPurchase != nil {
		contractDetails.FundMinimumInitialPurchase = contractDetailsProto.GetFundMinimumInitialPurchase()
	}
	if contractDetailsProto.FundMinimumSubsequentPurchase != nil {
		contractDetails.FundSubsequentMinimumPurchase = contractDetailsProto.GetFundMinimumSubsequentPurchase()
	}
	if contractDetailsProto.FundBlueSkyStates != nil {
		contractDetails.FundBlueSkyStates = contractDetailsProto.GetFundBlueSkyStates()
	}
	if contractDetailsProto.FundBlueSkyTerritories != nil {
		contractDetails.FundBlueSkyTerritories = contractDetailsProto.GetFundBlueSkyTerritories()
	}
	if contractDetailsProto.FundDistributionPolicyIndicator != nil {
		contractDetails.FundDistributionPolicyIndicator = getFundDistributionPolicyIndicator(contractDetailsProto.GetFundDistributionPolicyIndicator())
	}
	if contractDetailsProto.FundAssetType != nil {
		contractDetails.FundAssetType = getFundAssetType(contractDetailsProto.GetFundAssetType())
	}

	// bond fields
	if contractDetailsProto.Cusip != nil {
		contractDetails.Cusip = contractDetailsProto.GetCusip()
	}
	if contractDetailsProto.IssueDate != nil {
		contractDetails.IssueDate = contractDetailsProto.GetIssueDate()
	}
	if contractDetailsProto.Ratings != nil {
		contractDetails.Ratings = contractDetailsProto.GetRatings()
	}
	if contractDetailsProto.BondType != nil {
		contractDetails.BondType = contractDetailsProto.GetBondType()
	}
	if contractDetailsProto.Coupon != nil {
		contractDetails.Coupon = contractDetailsProto.GetCoupon()
	}
	if contractDetailsProto.CouponType != nil {
		contractDetails.CouponType = contractDetailsProto.GetCouponType()
	}
	if contractDetailsProto.Convertible != nil {
		contractDetails.Convertible = contractDetailsProto.GetConvertible()
	}
	if contractDetailsProto.Callable != nil {
		contractDetails.Callable = contractDetailsProto.GetCallable()
	}
	if contractDetailsProto.Puttable != nil {
		contractDetails.Putable = contractDetailsProto.GetPuttable()
	}
	if contractDetailsProto.DescAppend != nil {
		contractDetails.DescAppend = contractDetailsProto.GetDescAppend()
	}
	if contractDetailsProto.NextOptionDate != nil {
		contractDetails.NextOptionDate = contractDetailsProto.GetNextOptionDate()
	}
	if contractDetailsProto.NextOptionType != nil {
		contractDetails.NextOptionType = contractDetailsProto.GetNextOptionType()
	}
	if contractDetailsProto.NextOptionPartial != nil {
		contractDetails.NextOptionPartial = contractDetailsProto.GetNextOptionPartial()
	}
	if contractDetailsProto.BondNotes != nil {
		contractDetails.Notes = contractDetailsProto.GetBondNotes()
	}

	contractDetails.IneligibilityReasonList = decodeIneligibilityReasonList(contractDetailsProto)

	return contractDetails
}

func decodeIneligibilityReasonList(proto *protobuf.ContractDetails) []IneligibilityReason {
	var reasons []IneligibilityReason
	for _, reasonProto := range proto.GetIneligibilityReasonList() {
		reason := IneligibilityReason{}
		if reasonProto.Id != nil {
			reason.ID = reasonProto.GetId()
		}
		if reasonProto.Description != nil {
			reason.Description = reasonProto.GetDescription()
		}
		reasons = append(reasons, reason)
	}
	return reasons
}

func setLastTradeDate(lastTradeDateOrContractMonth string, contract *ContractDetails, isBond bool) {
	if lastTradeDateOrContractMonth == "" {
		return
	}
	var split []string
	splitWith := " "
	if strings.Contains(lastTradeDateOrContractMonth, "-") {
		splitWith = "-"
	}
	split = strings.Split(lastTradeDateOrContractMonth, splitWith)

	if len(split) > 0 {
		if isBond {
			contract.Maturity = split[0]
		} else {
			contract.Contract.LastTradeDateOrContractMonth = split[0]
		}
	}
	if len(split) > 1 {
		contract.LastTradeTime = split[1]
	}
	if isBond && len(split) > 2 {
		contract.TimeZoneID = split[2]
	}
}

// Helper for float to string conversion, if needed
func FloatToString(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

// decodeHistoricalTick translates protobuf HistoricalTick to Go HistoricalTick
func decodeHistoricalTick(proto *protobuf.HistoricalTick) *HistoricalTick {
	tick := &HistoricalTick{}
	if proto.Time != nil {
		tick.Time = proto.GetTime()
	}
	if proto.Price != nil {
		tick.Price = proto.GetPrice()
	}
	if proto.Size != nil {
		tick.Size = StringToDecimal(proto.GetSize())
	}
	return tick
}

// decodeHistoricalTickBidAsk translates protobuf HistoricalTickBidAsk to Go HistoricalTickBidAsk
func decodeHistoricalTickBidAsk(proto *protobuf.HistoricalTickBidAsk) *HistoricalTickBidAsk {
	tick := &HistoricalTickBidAsk{}
	if proto.Time != nil {
		tick.Time = proto.GetTime()
	}
	if proto.TickAttribBidAsk != nil {
		attrib := proto.GetTickAttribBidAsk()
		tick.TickAttribBidAsk.AskPastHigh = attrib.GetAskPastHigh()
		tick.TickAttribBidAsk.BidPastLow = attrib.GetBidPastLow()
	}
	if proto.PriceBid != nil {
		tick.PriceBid = proto.GetPriceBid()
	}
	if proto.PriceAsk != nil {
		tick.PriceAsk = proto.GetPriceAsk()
	}
	if proto.SizeBid != nil {
		tick.SizeBid = StringToDecimal(proto.GetSizeBid())
	}
	if proto.SizeAsk != nil {
		tick.SizeAsk = StringToDecimal(proto.GetSizeAsk())
	}
	return tick
}

// decodeHistoricalTickLast translates protobuf HistoricalTickLast to Go HistoricalTickLast
func decodeHistoricalTickLast(proto *protobuf.HistoricalTickLast) *HistoricalTickLast {
	tick := &HistoricalTickLast{}
	if proto.Time != nil {
		tick.Time = proto.GetTime()
	}
	if proto.TickAttribLast != nil {
		attrib := proto.GetTickAttribLast()
		tick.TickAttribLast.PastLimit = attrib.GetPastLimit()
		tick.TickAttribLast.Unreported = attrib.GetUnreported()
	}
	if proto.Price != nil {
		tick.Price = proto.GetPrice()
	}
	if proto.Size != nil {
		tick.Size = StringToDecimal(proto.GetSize())
	}
	if proto.Exchange != nil {
		tick.Exchange = proto.GetExchange()
	}
	if proto.SpecialConditions != nil {
		tick.SpecialConditions = proto.GetSpecialConditions()
	}
	return tick
}

// decodeHistogramDataEntry translates protobuf HistogramDataEntry to Go HistogramEntry
func decodeHistogramDataEntry(proto *protobuf.HistogramDataEntry) *HistogramData {
	entry := &HistogramData{}
	if proto.Price != nil {
		entry.Price = proto.GetPrice()
	}
	if proto.Size != nil {
		entry.Size = StringToDecimal(proto.GetSize())
	}
	return entry
}

// decodeHistoricalDataBar translates protobuf HistoricalDataBar to Go Bar
func decodeHistoricalDataBar(proto *protobuf.HistoricalDataBar) *Bar {
	bar := &Bar{}
	if proto.Date != nil {
		bar.Date = proto.GetDate()
	}
	if proto.Open != nil {
		bar.Open = proto.GetOpen()
	}
	if proto.High != nil {
		bar.High = proto.GetHigh()
	}
	if proto.Low != nil {
		bar.Low = proto.GetLow()
	}
	if proto.Close != nil {
		bar.Close = proto.GetClose()
	}
	if proto.Volume != nil {
		bar.Volume = StringToDecimal(proto.GetVolume())
	}
	if proto.BarCount != nil {
		bar.BarCount = int64(proto.GetBarCount())
	}
	if proto.WAP != nil {
		bar.Wap = StringToDecimal(proto.GetWAP())
	}
	return bar
}
