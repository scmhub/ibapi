package ibapi

import "fmt"

// HistoricalTickBidAsk is the historical tick's description.
// Used when requesting historical tick data with whatToShow = BID_ASK.
type HistoricalTickBidAsk struct {
	Time             int64
	TickAttribBidAsk TickAttribBidAsk
	PriceBid         float64
	PriceAsk         float64
	SizeBid          Decimal
	SizeAsk          Decimal
}

func NewHistoricalTickBidAsk() HistoricalTickBidAsk {
	htba := HistoricalTickBidAsk{}
	htba.SizeBid = UNSET_DECIMAL
	htba.SizeAsk = UNSET_DECIMAL
	return htba
}

func (h HistoricalTickBidAsk) String() string {
	return fmt.Sprintf("Time: %d, TickAttriBidAsk: %s, PriceBid: %f, PriceAsk: %f, SizeBid: %s, SizeAsk: %s",
		h.Time, h.TickAttribBidAsk, h.PriceBid, h.PriceAsk, DecimalMaxString(h.SizeBid), DecimalMaxString(h.SizeAsk))
}
