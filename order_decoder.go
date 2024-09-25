package ibapi

// OrderDecoder .
type OrderDecoder struct {
	order         *Order
	contract      *Contract
	orderState    *OrderState
	version       Version
	serverVersion Version
}

func (d *OrderDecoder) decodeOrderId(msgBuf *MsgBuffer) {
	d.order.OrderID = msgBuf.decodeInt64()
}

func (d *OrderDecoder) decodeContractFields(msgBuf *MsgBuffer) {
	d.contract.ConID = msgBuf.decodeInt64()
	d.contract.Symbol = msgBuf.decodeString()
	d.contract.SecType = msgBuf.decodeString()
	d.contract.LastTradeDateOrContractMonth = msgBuf.decodeString()
	d.contract.Strike = msgBuf.decodeFloat64()
	d.contract.Right = msgBuf.decodeString()
	if d.version >= 32 {
		d.contract.Multiplier = msgBuf.decodeString()
	}
	d.contract.Exchange = msgBuf.decodeString()
	d.contract.Currency = msgBuf.decodeString()
	d.contract.LocalSymbol = msgBuf.decodeString()
	if d.version >= 32 {
		d.contract.TradingClass = msgBuf.decodeString()
	}
}

func (d *OrderDecoder) decodeAction(msgBuf *MsgBuffer) {
	d.order.Action = msgBuf.decodeString()
}

func (d *OrderDecoder) decodeTotalQuantity(msgBuf *MsgBuffer) {
	d.order.TotalQuantity = msgBuf.decodeDecimal()
}

func (d *OrderDecoder) decodeOrderType(msgBuf *MsgBuffer) {
	d.order.OrderType = msgBuf.decodeString()
}

func (d *OrderDecoder) decodeLmtPrice(msgBuf *MsgBuffer) {
	if d.version < 29 {
		d.order.LmtPrice = msgBuf.decodeFloat64()
	} else {
		d.order.LmtPrice = msgBuf.decodeFloat64ShowUnset()
	}
}

func (d *OrderDecoder) decodeAuxPrice(msgBuf *MsgBuffer) {
	if d.version < 30 {
		d.order.AuxPrice = msgBuf.decodeFloat64()
	} else {
		d.order.AuxPrice = msgBuf.decodeFloat64ShowUnset()
	}
}

func (d *OrderDecoder) decodeTIF(msgBuf *MsgBuffer) {
	d.order.TIF = msgBuf.decodeString()
}

func (d *OrderDecoder) decodeOcaGroup(msgBuf *MsgBuffer) {
	d.order.OCAGroup = msgBuf.decodeString()
}

func (d *OrderDecoder) decodeAccount(msgBuf *MsgBuffer) {
	d.order.Account = msgBuf.decodeString()
}

func (d *OrderDecoder) decodeOpenClose(msgBuf *MsgBuffer) {
	d.order.OpenClose = msgBuf.decodeString()
}

func (d *OrderDecoder) decodeOrigin(msgBuf *MsgBuffer) {
	d.order.Origin = msgBuf.decodeInt64()
}

func (d *OrderDecoder) decodeOrderRef(msgBuf *MsgBuffer) {
	d.order.OrderRef = msgBuf.decodeString()
}

func (d *OrderDecoder) decodeClientId(msgBuf *MsgBuffer) {
	d.order.ClientID = msgBuf.decodeInt64()
}

func (d *OrderDecoder) decodePermId(msgBuf *MsgBuffer) {
	d.order.PermID = msgBuf.decodeInt64()
}

func (d *OrderDecoder) decodeOutsideRth(msgBuf *MsgBuffer) {
	d.order.OutsideRTH = msgBuf.decodeBool()
}

func (d *OrderDecoder) decodeHidden(msgBuf *MsgBuffer) {
	d.order.Hidden = msgBuf.decodeBool()
}

func (d *OrderDecoder) decodeDiscretionaryAmount(msgBuf *MsgBuffer) {
	d.order.DiscretionaryAmt = msgBuf.decodeFloat64()
}

func (d *OrderDecoder) decodeGoodAfterTime(msgBuf *MsgBuffer) {
	d.order.GoodAfterTime = msgBuf.decodeString()
}

func (d *OrderDecoder) skipSharesAllocation(msgBuf *MsgBuffer) {
	_ = msgBuf.decodeString()
}

func (d *OrderDecoder) decodeFAParams(msgBuf *MsgBuffer) {
	d.order.FAGroup = msgBuf.decodeString()
	d.order.FAMethod = msgBuf.decodeString()
	d.order.FAPercentage = msgBuf.decodeString()
	if d.serverVersion < MIN_SERVER_VER_FA_PROFILE_DESUPPORT {
		_ = msgBuf.decodeString() // skip deprecated FAProfile
	}
}

func (d *OrderDecoder) decodeModelCode(msgBuf *MsgBuffer) {
	if d.serverVersion >= MIN_SERVER_VER_MODELS_SUPPORT {
		d.order.ModelCode = msgBuf.decodeString()
	}
}

func (d *OrderDecoder) decodeGoodTillDate(msgBuf *MsgBuffer) {
	d.order.GoodTillDate = msgBuf.decodeString()
}

func (d *OrderDecoder) decodeRule80A(msgBuf *MsgBuffer) {
	d.order.Rule80A = msgBuf.decodeString()
}

func (d *OrderDecoder) decodePercentOffset(msgBuf *MsgBuffer) {
	d.order.PercentOffset = msgBuf.decodeFloat64ShowUnset()
}

func (d *OrderDecoder) decodeSettlingFirm(msgBuf *MsgBuffer) {
	d.order.SettlingFirm = msgBuf.decodeString()
}

func (d *OrderDecoder) decodeShortSaleParams(msgBuf *MsgBuffer) {
	d.order.ShortSaleSlot = msgBuf.decodeInt64()
	d.order.DesignatedLocation = msgBuf.decodeString()
	if d.serverVersion == MIN_SERVER_VER_SSHORTX_OLD {
		_ = msgBuf.decodeString()
	} else if d.version >= 23 {
		d.order.ExemptCode = msgBuf.decodeInt64()
	}
}

func (d *OrderDecoder) decodeAuctionStrategy(msgBuf *MsgBuffer) {
	d.order.AuctionStrategy = msgBuf.decodeInt64()
}

func (d *OrderDecoder) decodeBoxOrderParams(msgBuf *MsgBuffer) {
	d.order.StartingPrice = msgBuf.decodeFloat64ShowUnset()
	d.order.StockRefPrice = msgBuf.decodeFloat64ShowUnset()
	d.order.Delta = msgBuf.decodeFloat64ShowUnset()
}

func (d *OrderDecoder) decodePegToStkOrVolOrderParams(msgBuf *MsgBuffer) {
	d.order.StockRangeLower = msgBuf.decodeFloat64ShowUnset()
	d.order.StockRangeUpper = msgBuf.decodeFloat64ShowUnset()
}

func (d *OrderDecoder) decodeDisplaySize(msgBuf *MsgBuffer) {
	d.order.DisplaySize = msgBuf.decodeInt64ShowUnset() //show_unset
}

func (d *OrderDecoder) decodeBlockOrder(msgBuf *MsgBuffer) {
	d.order.BlockOrder = msgBuf.decodeBool()
}

func (d *OrderDecoder) decodeSweepToFill(msgBuf *MsgBuffer) {
	d.order.SweepToFill = msgBuf.decodeBool()
}

func (d *OrderDecoder) decodeAllOrNone(msgBuf *MsgBuffer) {
	d.order.AllOrNone = msgBuf.decodeBool()
}

func (d *OrderDecoder) decodeMinQty(msgBuf *MsgBuffer) {
	d.order.MinQty = msgBuf.decodeInt64ShowUnset()
}

func (d *OrderDecoder) decodeOcaType(msgBuf *MsgBuffer) {
	d.order.OCAType = msgBuf.decodeInt64()
}

func (d *OrderDecoder) skipETradeOnly(msgBuf *MsgBuffer) {
	_ = msgBuf.decodeBool() // deprecated order.ETradeOnly
}

func (d *OrderDecoder) skipFirmQuoteOnly(msgBuf *MsgBuffer) {
	_ = msgBuf.decodeBool() // deprecated order.FirmQuoteOnly
}

func (d *OrderDecoder) skipNbboPriceCap(msgBuf *MsgBuffer) {
	_ = msgBuf.decodeFloat64ShowUnset() // depracated order.NBBOPriceCap
}

func (d *OrderDecoder) decodeParentId(msgBuf *MsgBuffer) {
	d.order.ParentID = msgBuf.decodeInt64()
}

func (d *OrderDecoder) decodeTriggerMethod(msgBuf *MsgBuffer) {
	d.order.TriggerMethod = msgBuf.decodeInt64()
}

func (d *OrderDecoder) decodeVolOrderParams(msgBuf *MsgBuffer) {
	d.order.Volatility = msgBuf.decodeFloat64ShowUnset()
	d.order.VolatilityType = msgBuf.decodeInt64()
	d.order.DeltaNeutralOrderType = msgBuf.decodeString()
	d.order.DeltaNeutralAuxPrice = msgBuf.decodeFloat64ShowUnset()

	if d.version >= 27 && d.order.DeltaNeutralOrderType != "" {
		d.order.DeltaNeutralConID = msgBuf.decodeInt64()
		d.order.DeltaNeutralSettlingFirm = msgBuf.decodeString()
		d.order.DeltaNeutralClearingAccount = msgBuf.decodeString()
		d.order.DeltaNeutralClearingIntent = msgBuf.decodeString()
	}
	if d.version >= 31 && d.order.DeltaNeutralOrderType != "" {
		d.order.DeltaNeutralOpenClose = msgBuf.decodeString()
		d.order.DeltaNeutralShortSale = msgBuf.decodeBool()
		d.order.DeltaNeutralShortSaleSlot = msgBuf.decodeInt64()
		d.order.DeltaNeutralDesignatedLocation = msgBuf.decodeString()
	}

	d.order.ContinuousUpdate = msgBuf.decodeBool()
	d.order.ReferencePriceType = msgBuf.decodeInt64()
}

func (d *OrderDecoder) decodeTrailParams(msgBuf *MsgBuffer) {
	d.order.TrailStopPrice = msgBuf.decodeFloat64ShowUnset()
	if d.version >= 30 {
		d.order.TrailingPercent = msgBuf.decodeFloat64ShowUnset()
	}
}

func (d *OrderDecoder) decodeBasisPoints(msgBuf *MsgBuffer) {
	d.order.BasisPoints = msgBuf.decodeFloat64ShowUnset()
	d.order.BasisPointsType = msgBuf.decodeInt64ShowUnset()
}

func (d *OrderDecoder) decodeComboLegs(msgBuf *MsgBuffer) {
	d.contract.ComboLegsDescrip = msgBuf.decodeString()
	if d.version >= 29 {
		comboLegsCount := msgBuf.decodeInt64()
		d.contract.ComboLegs = make([]ComboLeg, 0, comboLegsCount)
		var i int64
		for i = 0; i < comboLegsCount; i++ {
			comboleg := ComboLeg{}
			comboleg.ConID = msgBuf.decodeInt64()
			comboleg.Ratio = msgBuf.decodeInt64()
			comboleg.Action = msgBuf.decodeString()
			comboleg.Exchange = msgBuf.decodeString()
			comboleg.OpenClose = msgBuf.decodeInt64()
			comboleg.ShortSaleSlot = msgBuf.decodeInt64()
			comboleg.DesignatedLocation = msgBuf.decodeString()
			comboleg.ExemptCode = msgBuf.decodeInt64()
			d.contract.ComboLegs = append(d.contract.ComboLegs, comboleg)
		}
		orderComboLegsCount := msgBuf.decodeInt64()
		d.order.OrderComboLegs = make([]OrderComboLeg, 0, orderComboLegsCount)
		for i = 0; i < orderComboLegsCount; i++ {
			orderComboLeg := OrderComboLeg{}
			orderComboLeg.Price = msgBuf.decodeFloat64ShowUnset()
			d.order.OrderComboLegs = append(d.order.OrderComboLegs, orderComboLeg)
		}
	}
}

func (d *OrderDecoder) decodeSmartComboRoutingParams(msgBuf *MsgBuffer) {
	if d.version >= 26 {
		smartComboRoutingParamsCount := msgBuf.decodeInt64()
		d.order.SmartComboRoutingParams = make([]TagValue, 0, smartComboRoutingParamsCount)
		var i int64
		for i = 0; i < smartComboRoutingParamsCount; i++ {
			tagValue := TagValue{}
			tagValue.Tag = msgBuf.decodeString()
			tagValue.Value = msgBuf.decodeString()
			d.order.SmartComboRoutingParams = append(d.order.SmartComboRoutingParams, tagValue)
		}
	}
}

func (d *OrderDecoder) decodeScaleOrderParams(msgBuf *MsgBuffer) {
	if d.version >= 20 {
		d.order.ScaleInitLevelSize = msgBuf.decodeInt64ShowUnset()
		d.order.ScaleSubsLevelSize = msgBuf.decodeInt64ShowUnset()
	} else {
		_ = msgBuf.decodeInt64ShowUnset() // deprecated notSuppScaleNumComponents
		d.order.ScaleInitLevelSize = msgBuf.decodeInt64ShowUnset()
	}
	d.order.ScalePriceIncrement = msgBuf.decodeFloat64ShowUnset()
	if d.version >= 28 && d.order.ScalePriceIncrement != UNSET_FLOAT && d.order.ScalePriceIncrement > 0.0 {
		d.order.ScalePriceAdjustValue = msgBuf.decodeFloat64ShowUnset()
		d.order.ScalePriceAdjustInterval = msgBuf.decodeInt64ShowUnset()
		d.order.ScaleProfitOffset = msgBuf.decodeFloat64ShowUnset()
		d.order.ScaleAutoReset = msgBuf.decodeBool()
		d.order.ScaleInitPosition = msgBuf.decodeInt64ShowUnset()
		d.order.ScaleInitFillQty = msgBuf.decodeInt64ShowUnset()
		d.order.ScaleRandomPercent = msgBuf.decodeBool()
	}
}

func (d *OrderDecoder) decodeHedgeParams(msgBuf *MsgBuffer) {
	if d.version >= 24 {
		d.order.HedgeType = msgBuf.decodeString()
		if d.order.HedgeType != "" {
			d.order.HedgeParam = msgBuf.decodeString()
		}
	}

}

func (d *OrderDecoder) decodeOptOutSmartRouting(msgBuf *MsgBuffer) {
	if d.version >= 25 {
		d.order.OptOutSmartRouting = msgBuf.decodeBool()
	}
}

func (d *OrderDecoder) decodeClearingParams(msgBuf *MsgBuffer) {
	d.order.ClearingAccount = msgBuf.decodeString()
	d.order.ClearingIntent = msgBuf.decodeString()
}

func (d *OrderDecoder) decodeNotHeld(msgBuf *MsgBuffer) {
	if d.version >= 22 {
		d.order.NotHeld = msgBuf.decodeBool()
	}
}

func (d *OrderDecoder) decodeDeltaNeutral(msgBuf *MsgBuffer) {
	if d.version >= 20 {
		deltaNeutralContractPresent := msgBuf.decodeBool()
		if deltaNeutralContractPresent {
			d.contract.DeltaNeutralContract = &DeltaNeutralContract{}
			d.contract.DeltaNeutralContract.ConID = msgBuf.decodeInt64()
			d.contract.DeltaNeutralContract.Delta = msgBuf.decodeFloat64()
			d.contract.DeltaNeutralContract.Price = msgBuf.decodeFloat64()
		}
	}
}

func (d *OrderDecoder) decodeAlgoParams(msgBuf *MsgBuffer) {
	if d.version >= 21 {
		d.order.AlgoStrategy = msgBuf.decodeString()
		if d.order.AlgoStrategy != "" {
			AlgoParamsCount := msgBuf.decodeInt64()
			d.order.AlgoParams = make([]TagValue, 0, AlgoParamsCount)
			var i int64
			for i = 0; i < AlgoParamsCount; i++ {
				tagValue := TagValue{}
				tagValue.Tag = msgBuf.decodeString()
				tagValue.Value = msgBuf.decodeString()
				d.order.AlgoParams = append(d.order.AlgoParams, tagValue)
			}
		}
	}
}

func (d *OrderDecoder) decodeSolicited(msgBuf *MsgBuffer) {
	if d.version >= 33 {
		d.order.Solictied = msgBuf.decodeBool()
	}

}

func (d *OrderDecoder) decodeWhatIfInfoAndCommission(msgBuf *MsgBuffer) {
	d.order.WhatIf = msgBuf.decodeBool()
	d.decodeOrderStatus(msgBuf)
	if d.serverVersion >= MIN_SERVER_VER_WHAT_IF_EXT_FIELDS {
		d.orderState.InitMarginBefore = msgBuf.decodeString()
		d.orderState.MaintMarginBefore = msgBuf.decodeString()
		d.orderState.EquityWithLoanBefore = msgBuf.decodeString()
		d.orderState.InitMarginChange = msgBuf.decodeString()
		d.orderState.MaintMarginChange = msgBuf.decodeString()
		d.orderState.EquityWithLoanChange = msgBuf.decodeString()
	}

	d.orderState.InitMarginAfter = msgBuf.decodeString()
	d.orderState.MaintMarginAfter = msgBuf.decodeString()
	d.orderState.EquityWithLoanAfter = msgBuf.decodeString()

	d.orderState.Commission = msgBuf.decodeFloat64ShowUnset()
	d.orderState.MinCommission = msgBuf.decodeFloat64ShowUnset()
	d.orderState.MaxCommission = msgBuf.decodeFloat64ShowUnset()
	d.orderState.CommissionCurrency = msgBuf.decodeString()
	d.orderState.WarningText = msgBuf.decodeString()
}

func (d *OrderDecoder) decodeOrderStatus(msgBuf *MsgBuffer) {
	d.orderState.Status = msgBuf.decodeString()
}

func (d *OrderDecoder) decodeVolRandomizeFlags(msgBuf *MsgBuffer) {
	if d.version >= 34 {
		d.order.RandomizeSize = msgBuf.decodeBool()
		d.order.RandomizePrice = msgBuf.decodeBool()
	}
}

func (d *OrderDecoder) decodePegBenchParams(msgBuf *MsgBuffer) {
	if d.serverVersion >= MIN_SERVER_VER_PEGGED_TO_BENCHMARK {
		if d.order.OrderType == "PEG BENCH" {
			d.order.ReferenceContractID = msgBuf.decodeInt64()
			d.order.IsPeggedChangeAmountDecrease = msgBuf.decodeBool()
			d.order.PeggedChangeAmount = msgBuf.decodeFloat64()
			d.order.ReferenceChangeAmount = msgBuf.decodeFloat64()
			d.order.ReferenceExchangeID = msgBuf.decodeString()
		}
	}
}

func (d *OrderDecoder) decodeConditions(msgBuf *MsgBuffer) {
	if d.serverVersion >= MIN_SERVER_VER_PEGGED_TO_BENCHMARK {
		conditionsSize := msgBuf.decodeInt64()
		d.order.Conditions = make([]OrderCondition, 0, conditionsSize)
		if conditionsSize > 0 {
			var i int64
			for i = 0; i < conditionsSize; i++ {
				conditionType := msgBuf.decodeInt64()
				cond := CreateOrderCondition(conditionType)
				cond.decode(msgBuf)

				d.order.Conditions = append(d.order.Conditions, cond)
			}
			d.order.ConditionsIgnoreRth = msgBuf.decodeBool()
			d.order.ConditionsCancelOrder = msgBuf.decodeBool()
		}
	}
}

func (d *OrderDecoder) decodeAdjustedOrderParams(msgBuf *MsgBuffer) {
	if d.serverVersion >= MIN_SERVER_VER_PEGGED_TO_BENCHMARK {
		d.order.AdjustedOrderType = msgBuf.decodeString()
		d.order.TriggerPrice = msgBuf.decodeFloat64()
		d.decodeStopPriceAndLmtPriceOffset(msgBuf)
		d.order.AdjustedStopPrice = msgBuf.decodeFloat64()
		d.order.AdjustedStopLimitPrice = msgBuf.decodeFloat64()
		d.order.AdjustedTrailingAmount = msgBuf.decodeFloat64()
		d.order.AdjustableTrailingUnit = msgBuf.decodeInt64()
	}
}

func (d *OrderDecoder) decodeStopPriceAndLmtPriceOffset(msgBuf *MsgBuffer) {
	d.order.TrailStopPrice = msgBuf.decodeFloat64()
	d.order.LmtPriceOffset = msgBuf.decodeFloat64()
}

func (d *OrderDecoder) decodeSoftDollarTier(msgBuf *MsgBuffer) {
	if d.serverVersion >= MIN_SERVER_VER_SOFT_DOLLAR_TIER {
		name := msgBuf.decodeString()
		value := msgBuf.decodeString()
		displayName := msgBuf.decodeString()
		d.order.SoftDollarTier = SoftDollarTier{name, value, displayName}
	}
}

func (d *OrderDecoder) decodeCashQty(msgBuf *MsgBuffer) {
	if d.serverVersion >= MIN_SERVER_VER_CASH_QTY {
		d.order.CashQty = msgBuf.decodeFloat64()
	}
}

func (d *OrderDecoder) decodeDontUseAutoPriceForHedge(msgBuf *MsgBuffer) {
	if d.serverVersion >= MIN_SERVER_VER_AUTO_PRICE_FOR_HEDGE {
		d.order.DontUseAutoPriceForHedge = msgBuf.decodeBool()
	}
}

func (d *OrderDecoder) decodeIsOmsContainer(msgBuf *MsgBuffer) {
	if d.serverVersion >= MIN_SERVER_VER_ORDER_CONTAINER {
		d.order.IsOmsContainer = msgBuf.decodeBool()
	}
}

func (d *OrderDecoder) decodeDiscretionaryUpToLimitPrice(msgBuf *MsgBuffer) {
	if d.serverVersion >= MIN_SERVER_VER_D_PEG_ORDERS {
		d.order.DiscretionaryUpToLimitPrice = msgBuf.decodeBool()
	}
}

func (d *OrderDecoder) decodeAutoCancelDate(msgBuf *MsgBuffer) {
	d.order.AutoCancelDate = msgBuf.decodeString()
}

func (d *OrderDecoder) decodeFilledQuantity(msgBuf *MsgBuffer) {
	d.order.FilledQuantity = msgBuf.decodeDecimal()
}

func (d *OrderDecoder) decodeRefFuturesConId(msgBuf *MsgBuffer) {
	d.order.RefFuturesConID = msgBuf.decodeInt64()
}

func (d *OrderDecoder) decodeAutoCancelParent(msgBuf *MsgBuffer, minVersionAutoCancelParent Version) {
	if d.serverVersion >= minVersionAutoCancelParent {
		d.order.AutoCancelParent = msgBuf.decodeBool()
	}
}

func (d *OrderDecoder) decodeShareholder(msgBuf *MsgBuffer) {
	d.order.Shareholder = msgBuf.decodeString()
}

func (d *OrderDecoder) decodeImbalanceOnly(msgBuf *MsgBuffer) {
	d.order.ImbalanceOnly = msgBuf.decodeBool()
}

func (d *OrderDecoder) decodeRouteMarketableToBbo(msgBuf *MsgBuffer) {
	d.order.RouteMarketableToBbo = msgBuf.decodeBool()
}

func (d *OrderDecoder) decodeParentPermId(msgBuf *MsgBuffer) {
	d.order.ParentPermID = msgBuf.decodeInt64()
}

func (d *OrderDecoder) decodeCompletedTime(msgBuf *MsgBuffer) {
	d.orderState.CompletedTime = msgBuf.decodeString()
}

func (d *OrderDecoder) decodeCompletedStatus(msgBuf *MsgBuffer) {
	d.orderState.CompletedStatus = msgBuf.decodeString()
}

func (d *OrderDecoder) decodeUsePriceMgmtAlgo(msgBuf *MsgBuffer) {
	if d.serverVersion >= MIN_SERVER_VER_PRICE_MGMT_ALGO {
		d.order.UsePriceMgmtAlgo = msgBuf.decodeBool()
	}
}

func (d *OrderDecoder) decodeDuration(msgBuf *MsgBuffer) {
	if d.serverVersion >= MIN_SERVER_VER_DURATION {
		d.order.Duration = msgBuf.decodeInt64ShowUnset()
	}
}

func (d *OrderDecoder) decodePostToAts(msgBuf *MsgBuffer) {
	if d.serverVersion >= MIN_SERVER_VER_POST_TO_ATS {
		d.order.PostToAts = msgBuf.decodeInt64ShowUnset()
	}
}

func (d *OrderDecoder) decodePegBestPegMidOrderAttributes(msgBuf *MsgBuffer) {
	if d.serverVersion >= MIN_SERVER_VER_PEGBEST_PEGMID_OFFSETS {
		d.order.MinTradeQty = msgBuf.decodeInt64ShowUnset()
		d.order.MinCompeteSize = msgBuf.decodeInt64ShowUnset()
		d.order.CompeteAgainstBestOffset = msgBuf.decodeFloat64ShowUnset()
		d.order.MidOffsetAtWhole = msgBuf.decodeFloat64ShowUnset()
		d.order.MidOffsetAtHalf = msgBuf.decodeFloat64ShowUnset()
	}
}

func (d *OrderDecoder) decodeCustomerAccount(msgBuf *MsgBuffer) {
	if d.serverVersion >= MIN_SERVER_VER_CUSTOMER_ACCOUNT {
		d.order.CustomerAccount = msgBuf.decodeString()
	}
}

func (d *OrderDecoder) decodeProfessionalCustomer(msgBuf *MsgBuffer) {
	if d.serverVersion >= MIN_SERVER_VER_PROFESSIONAL_CUSTOMER {
		d.order.ProfessionalCustomer = msgBuf.decodeBool()
	}
}

func (d *OrderDecoder) decodeBondAccruedInterest(msgBuf *MsgBuffer) {
	if d.serverVersion >= MIN_SERVER_VER_BOND_ACCRUED_INTEREST {
		d.order.BondAccruedInterest = msgBuf.decodeString()
	}
}

func (d *OrderDecoder) decodeIncludeOvernight(msgBuf *MsgBuffer) {
	if d.serverVersion >= MIN_SERVER_VER_INCLUDE_OVERNIGHT {
		d.order.IncludeOvernight = msgBuf.decodeBool()
	}
}
