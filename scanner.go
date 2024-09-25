package ibapi

import "fmt"

const NO_ROW_NUMBER_SPECIFIED int64 = -1

// ScanData .
type ScanData struct {
	Rank            int64
	ContractDetails *ContractDetails
	Distance        string
	Benchmark       string
	Projection      string
	LegsStr         string
}

func (s ScanData) String() string {
	return fmt.Sprintf("Rank: %d, Symbol: %s, SecType: %s, Currency: %s, Distance: %s, Benchmark: %s, Projection: %s, Legs String: %s",
		s.Rank, s.ContractDetails.Contract.Symbol, s.ContractDetails.Contract.SecType, s.ContractDetails.Contract.Currency,
		s.Distance, s.Benchmark, s.Projection, s.LegsStr)
}

// ScannerSubscription .
type ScannerSubscription struct {
	NumberOfRows             int64 `default:"NO_ROW_NUMBER_SPECIFIED"`
	Instrument               string
	LocationCode             string
	ScanCode                 string
	AbovePrice               float64 `default:"UNSET_FLOAT"`
	BelowPrice               float64 `default:"UNSET_FLOAT"`
	AboveVolume              int64   `default:"UNSET_INT"`
	MarketCapAbove           float64 `default:"UNSET_FLOAT"`
	MarketCapBelow           float64 `default:"UNSET_FLOAT"`
	MoodyRatingAbove         string
	MoodyRatingBelow         string
	SpRatingAbove            string
	SpRatingBelow            string
	MaturityDateAbove        string
	MaturityDateBelow        string
	CouponRateAbove          float64 `default:"UNSET_FLOAT"`
	CouponRateBelow          float64 `default:"UNSET_FLOAT"`
	ExcludeConvertible       bool
	AverageOptionVolumeAbove int64 `default:"UNSET_INT"`
	ScannerSettingPairs      string
	StockTypeFilter          string
}

func (s ScannerSubscription) String() string {
	return fmt.Sprintf("Instrument: %s, LocationCode: %s, ScanCode: %s", s.Instrument, s.LocationCode, s.ScanCode)
}

// NewScannerSubscription creates a default ScannerSubscription.
func NewScannerSubscription() *ScannerSubscription {
	scannerSubscription := &ScannerSubscription{}

	scannerSubscription.NumberOfRows = NO_ROW_NUMBER_SPECIFIED

	scannerSubscription.AbovePrice = UNSET_FLOAT
	scannerSubscription.BelowPrice = UNSET_FLOAT
	scannerSubscription.AboveVolume = UNSET_INT
	scannerSubscription.MarketCapAbove = UNSET_FLOAT
	scannerSubscription.MarketCapBelow = UNSET_FLOAT

	scannerSubscription.CouponRateAbove = UNSET_FLOAT
	scannerSubscription.CouponRateBelow = UNSET_FLOAT

	scannerSubscription.AverageOptionVolumeAbove = UNSET_INT

	return scannerSubscription
}
