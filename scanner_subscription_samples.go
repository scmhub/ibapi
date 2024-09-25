package ibapi

// HotUSStkByVolume .
// Hot US stocks by volume
func HotUSStkByVolume() *ScannerSubscription {

	scanSub := NewScannerSubscription()
	scanSub.Instrument = "STK"
	scanSub.LocationCode = "STK.US.MAJOR"
	scanSub.ScanCode = "HOT_BY_VOLUME"

	return scanSub
}

// TopPercentGainersIbis .
// Top % gainers at IBIS
func TopPercentGainersIbis() *ScannerSubscription {

	scanSub := NewScannerSubscription()
	scanSub.Instrument = "STOCK.EU"
	scanSub.LocationCode = "STK.EU.IBIS"
	scanSub.ScanCode = "TOP_PERC_GAIN"

	return scanSub
}

// MostActiveFutEurex .
// Most active futures at EUREX
func MostActiveFutEurex() *ScannerSubscription {

	scanSub := NewScannerSubscription()
	scanSub.Instrument = "FUT.EU"
	scanSub.LocationCode = "FUT.EU.EUREX"
	scanSub.ScanCode = "MOST_ACTIVE"

	return scanSub
}

// HighOptVolumePCRatioUSIndexes .
// High option volume P/C ratio US indexes
func HighOptVolumePCRatioUSIndexes() *ScannerSubscription {

	scanSub := NewScannerSubscription()
	scanSub.Instrument = "IND.US"
	scanSub.LocationCode = "IND.US"
	scanSub.ScanCode = "HIGH_OPT_VOLUME_PUT_CALL_RATIO"

	return scanSub
}

// ComplexOrdersAndTrades .
// Complex orders and trades scan, latest trades
func ComplexOrdersAndTrades() *ScannerSubscription {

	scanSub := NewScannerSubscription()
	scanSub.Instrument = "NATCOMB"
	scanSub.LocationCode = "NATCOMB.OPT.US"
	scanSub.ScanCode = "COMBO_LATEST_TRADE"

	return scanSub
}
