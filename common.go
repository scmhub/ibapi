package ibapi

import (
	"fmt"
	"math"
	"time"
)

const (
	UNSET_INT       int64   = math.MaxInt64
	UNSET_FLOAT     float64 = math.MaxFloat64
	INFINITY_STRING string  = "Infinity"
)

var INFINITY_FLOAT float64 = math.Inf(1)

type TickerID = int64

type OrderID = int64

type FaDataType int64

const (
	GROUPS  FaDataType = 1
	ALIASES FaDataType = 3
)

func (fa FaDataType) String() string {
	switch fa {
	case GROUPS:
		return "GROUPS"
	case ALIASES:
		return "ALIASES"
	default:
		return ""
	}
}

type MarketDataType int64

const (
	REALTIME       MarketDataType = 1
	FROZEN         MarketDataType = 2
	DELAYED        MarketDataType = 3
	DELAYED_FROZEN MarketDataType = 4
)

func (mdt MarketDataType) String() string {
	switch mdt {
	case REALTIME:
		return "REALTIME"
	case FROZEN:
		return "FROZEN"
	case DELAYED:
		return "DELAYED"
	case DELAYED_FROZEN:
		return "DELAYED_FROZEN"
	}
	return ""
}

// Bar .
type Bar struct {
	Date     string
	Open     float64
	High     float64
	Low      float64
	Close    float64
	Volume   Decimal
	Wap      Decimal
	BarCount int64
}

func NewBar() Bar {
	b := Bar{}
	b.Volume = UNSET_DECIMAL
	b.Wap = UNSET_DECIMAL
	return b
}

func (b Bar) String() string {
	return fmt.Sprintf("Date: %s, Open: %f, High: %f, Low: %f, Close: %f, Volume: %s, WAP: %s, BarCount: %d",
		b.Date, b.Open, b.High, b.Low, b.Close, decimalMaxString(b.Volume), decimalMaxString(b.Wap), b.BarCount)
}

// RealTimeBar .
type RealTimeBar struct {
	Time    int64
	EndTime int64
	Open    float64
	High    float64
	Low     float64
	Close   float64
	Volume  Decimal
	Wap     Decimal
	Count   int64
}

func NewRealTimeBar() RealTimeBar {
	rtb := RealTimeBar{}
	rtb.Volume = UNSET_DECIMAL
	rtb.Wap = UNSET_DECIMAL
	return rtb
}

func (rb RealTimeBar) String() string {
	return fmt.Sprintf("Time: %d, Open: %f, High: %f, Low: %f, Close: %f, Volume: %s, Wap: %s, Count: %d",
		rb.Time, rb.Open, rb.High, rb.Low, rb.Close, decimalMaxString(rb.Volume), decimalMaxString(rb.Wap), rb.Count)
}

// HistogramData .
type HistogramData struct {
	Price float64
	Size  Decimal
}

func NewHistogramData() HistogramData {
	hd := HistogramData{}
	hd.Size = UNSET_DECIMAL
	return hd
}

func (hd HistogramData) String() string {
	return fmt.Sprintf("Price: %v, Size: %v", hd.Price, hd.Size)
}

// NewsProvider .
type NewsProvider struct {
	Code string
	Name string
}

func NewNewsProvider() NewsProvider {
	return NewsProvider{}
}

func (np NewsProvider) String() string {
	return fmt.Sprintf("Code: %s, Name: %s", np.Code, np.Name)
}

// DepthMktDataDescription .
type DepthMktDataDescription struct {
	Exchange        string
	SecType         string
	ListingExch     string
	ServiceDataType string
	AggGroup        int64
}

func NewDepthMktDataDescription() DepthMktDataDescription {
	dmdd := DepthMktDataDescription{}
	dmdd.AggGroup = UNSET_INT
	return dmdd
}

// DepthMktDataDescription .
func (d DepthMktDataDescription) String() string {
	return fmt.Sprintf("Exchange: %s, SecType: %s, ListingExchange: %s, ServiceDataType: %s, AggGroup: %s",
		d.Exchange, d.SecType, d.ListingExch, d.ServiceDataType, intMaxString(d.AggGroup))
}

// SmartComponent .
type SmartComponent struct {
	BitNumber      int64
	Exchange       string
	ExchangeLetter string
}

func NewSmartComponent() SmartComponent {
	return SmartComponent{}
}

func (s SmartComponent) String() string {
	return fmt.Sprintf("BitNumber: %d, Exchange: %s, ExchangeLetter: %s", s.BitNumber, s.Exchange, s.ExchangeLetter)
}

// TickAttrib .
type TickAttrib struct {
	CanAutoExecute bool
	PastLimit      bool
	PreOpen        bool
}

func NewTickAttrib() TickAttrib {
	return TickAttrib{}
}

func (t TickAttrib) String() string {
	return fmt.Sprintf("CanAutoExecute: %t, PastLimit: %t, PreOpen: %t", t.CanAutoExecute, t.PastLimit, t.PreOpen)
}

// TickAttribBidAsk .
type TickAttribBidAsk struct {
	BidPastLow  bool
	AskPastHigh bool
}

func NewTickAttribBidAsk() TickAttribBidAsk {
	return TickAttribBidAsk{}
}

func (t TickAttribBidAsk) String() string {
	return fmt.Sprintf("BidPastLow: %t, AskPastHigh: %t", t.BidPastLow, t.AskPastHigh)
}

// TickAttribLast .
type TickAttribLast struct {
	PastLimit  bool
	Unreported bool
}

func NewTickAttribLast() TickAttribLast {
	return TickAttribLast{}
}

func (t TickAttribLast) String() string {
	return fmt.Sprintf("PastLimit: %t, Unreported: %t", t.PastLimit, t.Unreported)
}

// FamilyCode .
type FamilyCode struct {
	AccountID     string
	FamilyCodeStr string
}

func NewFamilyCode() FamilyCode {
	return FamilyCode{}
}

func (f FamilyCode) String() string {
	return fmt.Sprintf("AccountId: %s, FamilyCodeStr: %s", f.AccountID, f.FamilyCodeStr)
}

// PriceIncrement .
type PriceIncrement struct {
	LowEdge   float64
	Increment float64
}

func NewPriceIncrement() PriceIncrement {
	return PriceIncrement{}
}

func (p PriceIncrement) String() string {
	return fmt.Sprintf("LowEdge: %f, Increment: %f", p.LowEdge, p.Increment)
}

// HistoricalTick is the historical tick's description.
// Used when requesting historical tick data with whatToShow = MIDPOINT.
type HistoricalTick struct {
	Time  int64
	Price float64
	Size  Decimal
}

func NewHistoricalTick() HistoricalTick {
	ht := HistoricalTick{}
	ht.Size = UNSET_DECIMAL
	return ht
}

func (h HistoricalTick) String() string {
	return fmt.Sprintf("Time: %d, Price: %f, Size: %s", h.Time, h.Price, decimalMaxString(h.Size))
}

// HistoricalTickBidAsk is the historical tick's description.
// Used when requesting historical tick data with whatToShow = BID_ASK.
type HistoricalTickBidAsk struct {
	Time             int64
	TickAttirbBidAsk TickAttribBidAsk
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
		h.Time, h.TickAttirbBidAsk, h.PriceBid, h.PriceAsk, decimalMaxString(h.SizeBid), decimalMaxString(h.SizeAsk))
}

// HistoricalTickLast is the historical last tick's description.
// Used when requesting historical tick data with whatToShow = TRADES.
type HistoricalTickLast struct {
	Time              int64
	TickAttribLast    TickAttribLast
	Price             float64
	Size              Decimal
	Exchange          string
	SpecialConditions string
}

func NewHistoricalTickLast() HistoricalTickLast {
	htl := HistoricalTickLast{}
	htl.Size = UNSET_DECIMAL
	return htl
}

func (h HistoricalTickLast) String() string {
	return fmt.Sprintf("Time: %d, TickAttribLast: %s, Price: %f, Size: %s, Exchange: %s, SpecialConditions: %s",
		h.Time, h.TickAttribLast, h.Price, decimalMaxString(h.Size), h.Exchange, h.SpecialConditions)
}

func (h HistoricalTickLast) Timestamp() time.Time {
	return time.Unix(h.Time, 0)
}

// HistoricalSession .
type HistoricalSession struct {
	StartDateTime string
	EndDateTime   string
	RefDate       string
}

func NewHistoricalSession() HistoricalSession {
	return HistoricalSession{}
}

func (h HistoricalSession) String() string {
	return fmt.Sprintf("Start: %s, End: %s, Ref Date: %s", h.StartDateTime, h.EndDateTime, h.RefDate)
}

// WshEventData .
type WshEventData struct {
	ConID           int64
	Filter          string
	FillWatchList   bool
	FillPortfolio   bool
	FillCompetitors bool
	StartDate       string
	EndDate         string
	TotalLimit      int64
}

func NewWshEventData() WshEventData {
	wed := WshEventData{}
	wed.ConID = UNSET_INT
	wed.TotalLimit = UNSET_INT
	return wed
}

func (w WshEventData) String() string {
	return fmt.Sprintf("WshEventData. ConId: %s, Filter: %s, Fill Watchlist: %t, Fill Portfolio: %t, Fill Competitors: %t",
		intMaxString(w.ConID), w.Filter, w.FillWatchList, w.FillPortfolio, w.FillCompetitors)
}
