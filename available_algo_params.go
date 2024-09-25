package ibapi

import (
	"strconv"
)

func int64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}
func float64ToString(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}
func boolToString(b bool) string {
	if b {
		return "1"
	}
	return "0"
}

// FillArrivalPriceParams .
func FillArrivalPriceParams(baseOrder *Order, maxPctVol float64, riskAversion, startTime, endTime string, forceCompletion, allowPastTime bool) {
	baseOrder.AlgoStrategy = "ArrivalPx"
	baseOrder.AlgoParams = []TagValue{}
	tag1 := TagValue{Tag: "maxPctVol", Value: float64ToString(maxPctVol)}
	tag2 := TagValue{Tag: "riskAversion", Value: riskAversion}
	tag3 := TagValue{Tag: "startTime", Value: startTime}
	tag4 := TagValue{Tag: "endTime", Value: endTime}
	tag5 := TagValue{Tag: "forceCompletion", Value: boolToString(forceCompletion)}
	tag6 := TagValue{Tag: "allowPastEndTime", Value: boolToString(allowPastTime)}
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag1)
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag2)
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag3)
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag4)
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag5)
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag6)
}

// FillDarkIceParams .
func FillDarkIceParams(baseOrder *Order, displaySize int64, startTime, endTime string, allowPastEndTime bool) {
	baseOrder.AlgoStrategy = "DarkIce"
	baseOrder.AlgoParams = []TagValue{}
	tag1 := TagValue{Tag: "displaySize", Value: int64ToString(displaySize)}
	tag2 := TagValue{Tag: "startTime", Value: startTime}
	tag3 := TagValue{Tag: "endTime", Value: endTime}
	tag4 := TagValue{Tag: "allowPastEndTime", Value: boolToString(allowPastEndTime)}
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag1)
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag2)
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag3)
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag4)
}

// FillPctVolParams .
func FillPctVolParams(baseOrder *Order, pctVol float64, startTime, endTime string, noTakeLiq bool) {
	baseOrder.AlgoStrategy = "PctVol"
	baseOrder.AlgoParams = []TagValue{}
	tag1 := TagValue{Tag: "pctVol", Value: float64ToString(pctVol)}
	tag2 := TagValue{Tag: "startTime", Value: startTime}
	tag3 := TagValue{Tag: "endTime", Value: endTime}
	tag4 := TagValue{Tag: "noTakeLiq", Value: boolToString(noTakeLiq)}
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag1)
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag2)
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag3)
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag4)
}

// FillTwapParams .
func FillTwapParams(baseOrder *Order, strategyType string, startTime, endTime string, allowPastEndTime bool) {
	baseOrder.AlgoStrategy = "Twap"
	baseOrder.AlgoParams = []TagValue{}
	tag1 := TagValue{Tag: "strategyType", Value: strategyType}
	tag2 := TagValue{Tag: "startTime", Value: startTime}
	tag3 := TagValue{Tag: "endTime", Value: endTime}
	tag4 := TagValue{Tag: "allowPastEndTime", Value: boolToString(allowPastEndTime)}
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag1)
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag2)
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag3)
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag4)
}

// FillVwapParams .
func FillVwapParams(baseOrder *Order, maxPctVol float64, startTime, endTime string, allowPastEndTime, noTakeLiq, speedUp bool) {
	baseOrder.AlgoStrategy = "Vwap"
	baseOrder.AlgoParams = []TagValue{}
	tag1 := TagValue{Tag: "maxPctVol", Value: float64ToString(maxPctVol)}
	tag2 := TagValue{Tag: "startTime", Value: startTime}
	tag3 := TagValue{Tag: "endTime", Value: endTime}
	tag4 := TagValue{Tag: "allowPastEndTime", Value: boolToString(allowPastEndTime)}
	tag5 := TagValue{Tag: "noTakeLiq", Value: boolToString(noTakeLiq)}
	tag6 := TagValue{Tag: "speedUp", Value: boolToString(speedUp)}
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag1)
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag2)
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag3)
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag4)
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag5)
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag6)
}

// FillAccumulateDistributeParams .
func FillAccumulateDistributeParams(baseOrder *Order, componentSize, timeBetweenOrders int64, randomizeTime20, randomizeSize55 bool,
	giveUp int64, catchUp, waitForFill bool, startTime, endTime string) {
	baseOrder.AlgoStrategy = "AD"
	baseOrder.AlgoParams = []TagValue{}
	tag1 := TagValue{Tag: "componentSize", Value: int64ToString(componentSize)}
	tag2 := TagValue{Tag: "timeBetweenOrders", Value: int64ToString(timeBetweenOrders)}
	tag3 := TagValue{Tag: "randomizeTime20", Value: boolToString(randomizeTime20)}
	tag4 := TagValue{Tag: "randomizeSize55", Value: boolToString(randomizeSize55)}
	tag5 := TagValue{Tag: "giveUp", Value: int64ToString(giveUp)}
	tag6 := TagValue{Tag: "catchUp", Value: boolToString(catchUp)}
	tag7 := TagValue{Tag: "waitForFill", Value: boolToString(waitForFill)}
	tag8 := TagValue{Tag: "activeTimeStart", Value: startTime}
	tag9 := TagValue{Tag: "activeTimeEnd", Value: endTime}
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag1)
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag2)
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag3)
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag4)
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag5)
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag6)
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag7)
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag8)
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag9)
}

// FillBalanceImpactRiskParams .
func FillBalanceImpactRiskParams(baseOrder *Order, maxPctVol float64, riskAversion string, forceCompletion bool) {
	baseOrder.AlgoStrategy = "BalanceImpactRisk"
	baseOrder.AlgoParams = []TagValue{}
	tag1 := TagValue{Tag: "maxPctVol", Value: float64ToString(maxPctVol)}
	tag2 := TagValue{Tag: "riskAversion", Value: riskAversion}
	tag3 := TagValue{Tag: "forceCompletion", Value: boolToString(forceCompletion)}
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag1)
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag2)
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag3)
}

// FillMinImpactParams .
func FillMinImpactParams(baseOrder *Order, maxPctVol float64) {
	baseOrder.AlgoStrategy = "MinImpact"
	baseOrder.AlgoParams = []TagValue{}
	tag1 := TagValue{Tag: "maxPctVol", Value: float64ToString(maxPctVol)}
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag1)
}

// FillAdaptiveParams .
func FillAdaptiveParams(baseOrder *Order, priority string) {
	baseOrder.AlgoStrategy = "Adaptive"
	baseOrder.AlgoParams = []TagValue{}
	tag1 := TagValue{Tag: "adaptivePriority", Value: priority}
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag1)
}

// FillClosePriceParams .
func FillClosePriceParams(baseOrder *Order, maxPctVol float64, riskAversion, startTime string, forceCompletion bool) {
	baseOrder.AlgoStrategy = "ClosePx"
	baseOrder.AlgoParams = []TagValue{}
	tag1 := TagValue{Tag: "maxPctVol", Value: float64ToString(maxPctVol)}
	tag2 := TagValue{Tag: "riskAversion", Value: riskAversion}
	tag3 := TagValue{Tag: "startTime", Value: startTime}
	tag4 := TagValue{Tag: "forceCompletion", Value: boolToString(forceCompletion)}
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag1)
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag2)
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag3)
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag4)
}

// FillPriceVariantPctVolParams .
func FillPriceVariantPctVolParams(baseOrder *Order, pctVol, deltaPctVol, minPctVol4Px, maxPctVol4Px float64, startTime string, endTime string, noTakeLiq bool) {
	baseOrder.AlgoStrategy = "PctVolPx"
	baseOrder.AlgoParams = []TagValue{}
	tag1 := TagValue{Tag: "pctVol", Value: float64ToString(pctVol)}
	tag2 := TagValue{Tag: "deltaPctVol", Value: float64ToString(deltaPctVol)}
	tag3 := TagValue{Tag: "minPctVol4Px", Value: float64ToString(minPctVol4Px)}
	tag4 := TagValue{Tag: "maxPctVol4Px", Value: float64ToString(maxPctVol4Px)}
	tag5 := TagValue{Tag: "startTime", Value: startTime}
	tag6 := TagValue{Tag: "endTime", Value: endTime}
	tag7 := TagValue{Tag: "noTakeLiq", Value: boolToString(noTakeLiq)}
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag1)
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag2)
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag3)
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag4)
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag5)
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag6)
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag7)
}

// FillSizeVariantPctVolParams .
func FillSizeVariantPctVolParams(baseOrder *Order, startPctVol, endPctVol float64, startTime, endTime string, noTakeLiq bool) {
	baseOrder.AlgoStrategy = "PctVolSz"
	baseOrder.AlgoParams = []TagValue{}
	tag1 := TagValue{Tag: "startPctVol", Value: float64ToString(startPctVol)}
	tag2 := TagValue{Tag: "endPctVol", Value: float64ToString(endPctVol)}
	tag3 := TagValue{Tag: "startTime", Value: startTime}
	tag4 := TagValue{Tag: "endTime", Value: endTime}
	tag5 := TagValue{Tag: "noTakeLiq", Value: boolToString(noTakeLiq)}
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag1)
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag2)
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag3)
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag4)
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag5)
}

// FillTimeVariantPctVolParams .
func FillTimeVariantPctVolParams(baseOrder *Order, startPctVol, endPctVol float64, startTime, endTime string, noTakeLiq bool) {
	baseOrder.AlgoStrategy = "PctVolTm"
	baseOrder.AlgoParams = []TagValue{}
	tag1 := TagValue{Tag: "startPctVol", Value: float64ToString(startPctVol)}
	tag2 := TagValue{Tag: "endPctVol", Value: float64ToString(endPctVol)}
	tag3 := TagValue{Tag: "startTime", Value: startTime}
	tag4 := TagValue{Tag: "endTime", Value: endTime}
	tag5 := TagValue{Tag: "noTakeLiq", Value: boolToString(noTakeLiq)}
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag1)
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag2)
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag3)
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag4)
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag5)
}

// FillJefferiesVWAPParams .
func FillJefferiesVWAPParams(baseOrder *Order, startTime, endTime string, relativeLimit, maxVolumeRate float64, excludeAuctions string,
	triggerPrice, wowPrice float64, minFillSize int64, wowOrderPct float64, wowMode string, isBuyBack bool, wowReference string) {
	baseOrder.AlgoStrategy = "VWAP"
	baseOrder.AlgoParams = []TagValue{}
	tag1 := TagValue{Tag: "StartTime", Value: startTime}
	tag2 := TagValue{Tag: "EndTime", Value: endTime}
	tag3 := TagValue{Tag: "RelativeLimit", Value: float64ToString(relativeLimit)}
	tag4 := TagValue{Tag: "MaxVolumeRate", Value: float64ToString(maxVolumeRate)}
	tag5 := TagValue{Tag: "ExcludeAuctions", Value: excludeAuctions}
	tag6 := TagValue{Tag: "TriggerPrice", Value: float64ToString(triggerPrice)}
	tag7 := TagValue{Tag: "WowPrice", Value: float64ToString(wowPrice)}
	tag8 := TagValue{Tag: "MinFillSize", Value: int64ToString(minFillSize)}
	tag9 := TagValue{Tag: "WowOrderPct", Value: float64ToString(wowOrderPct)}
	tag10 := TagValue{Tag: "WowMode", Value: wowMode}
	tag11 := TagValue{Tag: "IsBuyBack", Value: boolToString(isBuyBack)}
	tag12 := TagValue{Tag: "WowReference", Value: wowReference}
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag1)
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag2)
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag3)
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag4)
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag5)
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag6)
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag7)
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag8)
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag9)
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag10)
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag11)
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag12)
}

// FillCSFBInlineParams .
func FillCSFBInlineParams(baseOrder *Order, startTime, endTime, execStyle string, minPercent, maxPercent, displaySize int64, auction string, blockFinder bool,
	blockPrice float64, minBlockSize, maxBlockSize int64, iWouldPrice float64) {

	// must be direct-routed to "CSFBALGO"

	baseOrder.AlgoStrategy = "INLINE"
	baseOrder.AlgoParams = []TagValue{}
	tag1 := TagValue{Tag: "StartTime", Value: startTime}
	tag2 := TagValue{Tag: "EndTime", Value: endTime}
	tag3 := TagValue{Tag: "ExecStyle", Value: execStyle}
	tag4 := TagValue{Tag: "MinPercent", Value: int64ToString(minPercent)}
	tag5 := TagValue{Tag: "MaxPercent", Value: int64ToString(maxPercent)}
	tag6 := TagValue{Tag: "DisplaySize", Value: int64ToString(displaySize)}
	tag7 := TagValue{Tag: "Auction", Value: auction}
	tag8 := TagValue{Tag: "BlockFinder", Value: boolToString(blockFinder)}
	tag9 := TagValue{Tag: "BlockPrice", Value: float64ToString(blockPrice)}
	tag10 := TagValue{Tag: "MinBlockSize", Value: int64ToString(minBlockSize)}
	tag11 := TagValue{Tag: "MaxBlockSize", Value: int64ToString(maxBlockSize)}
	tag12 := TagValue{Tag: "IWouldPrice", Value: float64ToString(iWouldPrice)}
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag1)
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag2)
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag3)
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag4)
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag5)
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag6)
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag7)
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag8)
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag9)
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag10)
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag11)
	baseOrder.AlgoParams = append(baseOrder.AlgoParams, tag12)
}
