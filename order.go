package ibapi

import "fmt"

type AuctionStrategy = int64

const (
	AUCTION_UNSET       AuctionStrategy = 0
	AUCTION_MATCH       AuctionStrategy = 1
	AUCTION_IMPROVEMENT AuctionStrategy = 2
	AUCTION_TRANSPARENT AuctionStrategy = 3
)

type UsePriceMmgtAlgo = int64

const (
	USE_PRICE_MGMT_ALGO_DONT_USE UsePriceMmgtAlgo = 0
	USE_PRICE_MGMT_ALGO_USE      UsePriceMmgtAlgo = 1
	USE_PRICE_MGMT_ALGO_DEFAULT  UsePriceMmgtAlgo = UNSET_INT
)

var COMPETE_AGAINST_BEST_OFFSET_UP_TO_MID = INFINITY_FLOAT

type ThreeStateBoolean int64

const (
	STATE_NO      ThreeStateBoolean = 0
	STATE_YES     ThreeStateBoolean = 1
	STATE_DEFAULT ThreeStateBoolean = ThreeStateBoolean(UNSET_INT) // or some other sentinel value
)

// OrderComboLeg .
type OrderComboLeg struct {
	Price float64 `default:"UNSET_FLOAT"`
}

// NewOrder creates a default OrderComboLeg.
func NewOrderComboLeg() OrderComboLeg {
	ocl := OrderComboLeg{}
	ocl.Price = UNSET_FLOAT
	return ocl
}

func (o OrderComboLeg) String() string {
	return fmt.Sprintf("%s ", FloatMaxString(o.Price))
}

// Order .
type Order struct {
	// order identifier
	OrderID  int64
	ClientID int64
	PermID   int64

	// main order fields
	Action        string
	TotalQuantity Decimal `default:"UNSET_DECIMAL"`
	OrderType     string
	LmtPrice      float64 `default:"UNSET_FLOAT"`
	AuxPrice      float64 `default:"UNSET_FLOAT"`

	// extended order fields
	TIF                           string // "Time in Force" - DAY, GTC, etc.
	ActiveStartTime               string // for GTC orders
	ActiveStopTime                string // for GTC orders
	OCAGroup                      string // one cancels all group name
	OCAType                       int64  // 1 = CANCEL_WITH_BLOCK, 2 = REDUCE_WITH_BLOCK, 3 = REDUCE_NON_BLOCK
	OrderRef                      string // order reference
	Transmit                      bool   `default:"true"` // if false, order will be created but not transmitted
	ParentID                      int64  // Parent order Id, to associate Auto STP or TRAIL orders with the original order.
	BlockOrder                    bool
	SweepToFill                   bool
	DisplaySize                   int64
	TriggerMethod                 int64 // 0=Default, 1=Double_Bid_Ask, 2=Last, 3=Double_Last, 4=Bid_Ask, 7=Last_or_Bid_Ask, 8=Mid-point
	OutsideRTH                    bool
	Hidden                        bool
	GoodAfterTime                 string // Format: 20060505 08:00:00 {time zone}
	GoodTillDate                  string // Format: 20060505 08:00:00 {time zone}
	Rule80A                       string // Individual = 'I', Agency = 'A', AgentOtherMember = 'W', IndividualPTIA = 'J', AgencyPTIA = 'U', AgentOtherMemberPTIA = 'M', IndividualPT = 'K', AgencyPT = 'Y', AgentOtherMemberPT = 'N'
	AllOrNone                     bool
	MinQty                        int64   `default:"UNSET_INT"`
	PercentOffset                 float64 `default:"UNSET_FLOAT"` // REL orders only
	OverridePercentageConstraints bool
	TrailStopPrice                float64 `default:"UNSET_FLOAT"` // TRAILLIMIT orders only
	TrailingPercent               float64 `default:"UNSET_FLOAT"`

	// financial advisors only
	FAGroup      string
	FAMethod     string
	FAPercentage string

	// institutional (ie non-cleared) only
	OpenClose          string // O=Open, C=Close
	Origin             int64  // 0=Customer, 1=Firm
	ShortSaleSlot      int64  // 1 if you hold the shares, 2 if they will be delivered from elsewhere.  Only for Action=SSHORT
	DesignatedLocation string // set when slot=2 only.
	ExemptCode         int64  `default:"-1"`

	// SMART routing only
	DiscretionaryAmt   float64
	OptOutSmartRouting bool

	// BOX exchange orders only
	AuctionStrategy AuctionStrategy // AUCTION_UNSET, AUCTION_MATCH, AUCTION_IMPROVEMENT, AUCTION_TRANSPARENT
	StartingPrice   float64         `default:"UNSET_FLOAT"`
	StockRefPrice   float64         `default:"UNSET_FLOAT"`
	Delta           float64         `default:"UNSET_FLOAT"`

	// pegged to stock and VOL orders only
	StockRangeLower float64 `default:"UNSET_FLOAT"`
	StockRangeUpper float64 `default:"UNSET_FLOAT"`

	RandomizeSize  bool
	RandomizePrice bool

	// VOLATILITY ORDERS ONLY
	Volatility                     float64 `default:"UNSET_FLOAT"`
	VolatilityType                 int64   `default:"UNSET_INT"`
	DeltaNeutralOrderType          string
	DeltaNeutralAuxPrice           float64 `default:"UNSET_FLOAT"`
	DeltaNeutralConID              int64
	DeltaNeutralSettlingFirm       string
	DeltaNeutralClearingAccount    string
	DeltaNeutralClearingIntent     string
	DeltaNeutralOpenClose          string
	DeltaNeutralShortSale          bool
	DeltaNeutralShortSaleSlot      int64
	DeltaNeutralDesignatedLocation string
	ContinuousUpdate               bool
	ReferencePriceType             int64 `default:"UNSET_INT"` // 1=Average, 2 = BidOrAsk

	// COMBO ORDERS ONLY
	BasisPoints     float64 `default:"UNSET_FLOAT"` // EFP orders only
	BasisPointsType int64   `default:"UNSET_INT"`   // EFP orders only

	// SCALE ORDERS ONLY
	ScaleInitLevelSize       int64   `default:"UNSET_INT"`
	ScaleSubsLevelSize       int64   `default:"UNSET_INT"`
	ScalePriceIncrement      float64 `default:"UNSET_FLOAT"`
	ScalePriceAdjustValue    float64 `default:"UNSET_FLOAT"`
	ScalePriceAdjustInterval int64   `default:"UNSET_INT"`
	ScaleProfitOffset        float64 `default:"UNSET_FLOAT"`
	ScaleAutoReset           bool
	ScaleInitPosition        int64 `default:"UNSET_INT"`
	ScaleInitFillQty         int64 `default:"UNSET_INT"`
	ScaleRandomPercent       bool
	ScaleTable               string

	// HEDGE ORDERS
	HedgeType  string // 'D' - delta, 'B' - beta, 'F' - FX, 'P' - pair
	HedgeParam string // 'beta=X' value for beta hedge, 'ratio=Y' for pair hedge

	// Clearing info
	Account         string // IB account
	SettlingFirm    string
	ClearingAccount string // True beneficiary of the order
	ClearingIntent  string // "" (Default), "IB", "Away", "PTA" (PostTrade)

	// ALGO ORDERS ONLY
	AlgoStrategy string

	AlgoParams              []TagValue
	SmartComboRoutingParams []TagValue

	AlgoID string

	// What-if
	WhatIf bool

	// Not Held
	NotHeld   bool
	Solicited bool

	// models
	ModelCode string

	// order combo legs
	OrderComboLegs   []OrderComboLeg
	OrderMiscOptions []TagValue

	//VER PEG2BENCH fields:
	ReferenceContractID          int64
	PeggedChangeAmount           float64
	IsPeggedChangeAmountDecrease bool
	ReferenceChangeAmount        float64
	ReferenceExchangeID          string
	AdjustedOrderType            string
	TriggerPrice                 float64 `default:"UNSET_FLOAT"`
	AdjustedStopPrice            float64 `default:"UNSET_FLOAT"`
	AdjustedStopLimitPrice       float64 `default:"UNSET_FLOAT"`
	AdjustedTrailingAmount       float64 `default:"UNSET_FLOAT"`
	AdjustableTrailingUnit       int64
	LmtPriceOffset               float64 `default:"UNSET_FLOAT"`

	Conditions            []OrderCondition
	ConditionsCancelOrder bool
	ConditionsIgnoreRth   bool

	// ext operator
	ExtOperator string

	SoftDollarTier SoftDollarTier

	// native cash quantity
	CashQty float64 `default:"UNSET_FLOAT"`

	Mifid2DecisionMaker   string
	Mifid2DecisionAlgo    string
	Mifid2ExecutionTrader string
	Mifid2ExecutionAlgo   string

	// don't use auto price for hedge
	DontUseAutoPriceForHedge bool

	IsOmsContainer bool

	DiscretionaryUpToLimitPrice bool

	AutoCancelDate       string
	FilledQuantity       Decimal `default:"UNSET_DECIMAL"`
	RefFuturesConID      int64
	AutoCancelParent     bool
	Shareholder          string
	ImbalanceOnly        bool
	RouteMarketableToBbo bool
	ParentPermID         int64

	UsePriceMgmtAlgo         UsePriceMmgtAlgo
	Duration                 int64 `default:"UNSET_INT"`
	PostToAts                int64 `default:"UNSET_INT"`
	AdvancedErrorOverride    string
	ManualOrderTime          string
	MinTradeQty              int64   `default:"UNSET_INT"`
	MinCompeteSize           int64   `default:"UNSET_INT"`
	CompeteAgainstBestOffset float64 `default:"UNSET_FLOAT"`
	MidOffsetAtWhole         float64 `default:"UNSET_FLOAT"`
	MidOffsetAtHalf          float64 `default:"UNSET_FLOAT"`
	CustomerAccount          string
	ProfessionalCustomer     bool
	BondAccruedInterest      string
	IncludeOvernight         bool
	ManualOrderIndicator     int64 `default:"UNSET_INT"`
	Submitter                string
	PostOnly                 bool
	AllowPreOpen             bool
	IgnoreOpenAuction        bool
	Deactivate               bool
	SeekPriceImprovement     ThreeStateBoolean
	WhatIfType               int64
}

// NewOrder creates a default Order.
func NewOrder() *Order {
	order := &Order{}
	order.TotalQuantity = UNSET_DECIMAL
	order.LmtPrice = UNSET_FLOAT
	order.AuxPrice = UNSET_FLOAT

	order.Transmit = true
	order.MinQty = UNSET_INT
	order.PercentOffset = UNSET_FLOAT
	order.TrailStopPrice = UNSET_FLOAT
	order.TrailingPercent = UNSET_FLOAT

	order.ExemptCode = -1

	order.AuctionStrategy = AUCTION_UNSET
	order.StartingPrice = UNSET_FLOAT
	order.StockRefPrice = UNSET_FLOAT
	order.Delta = UNSET_FLOAT

	order.StockRangeLower = UNSET_FLOAT
	order.StockRangeUpper = UNSET_FLOAT

	order.Volatility = UNSET_FLOAT
	order.VolatilityType = UNSET_INT
	order.DeltaNeutralAuxPrice = UNSET_FLOAT
	order.ReferencePriceType = UNSET_INT

	order.BasisPoints = UNSET_FLOAT
	order.BasisPointsType = UNSET_INT

	order.ScaleInitLevelSize = UNSET_INT
	order.ScaleSubsLevelSize = UNSET_INT
	order.ScalePriceIncrement = UNSET_FLOAT
	order.ScalePriceAdjustValue = UNSET_FLOAT
	order.ScalePriceAdjustInterval = UNSET_INT
	order.ScaleProfitOffset = UNSET_FLOAT
	order.ScaleInitPosition = UNSET_INT
	order.ScaleInitFillQty = UNSET_INT

	order.TriggerPrice = UNSET_FLOAT
	order.AdjustedStopPrice = UNSET_FLOAT
	order.AdjustedStopLimitPrice = UNSET_FLOAT
	order.AdjustedTrailingAmount = UNSET_FLOAT
	order.LmtPriceOffset = UNSET_FLOAT

	order.CashQty = UNSET_FLOAT

	order.FilledQuantity = UNSET_DECIMAL

	order.UsePriceMgmtAlgo = USE_PRICE_MGMT_ALGO_DEFAULT
	order.Duration = UNSET_INT
	order.PostToAts = UNSET_INT
	order.MinTradeQty = UNSET_INT
	order.MinCompeteSize = UNSET_INT
	order.CompeteAgainstBestOffset = UNSET_FLOAT
	order.MidOffsetAtWhole = UNSET_FLOAT
	order.MidOffsetAtHalf = UNSET_FLOAT
	order.ManualOrderIndicator = UNSET_INT
	order.SeekPriceImprovement = STATE_DEFAULT
	order.WhatIfType = UNSET_INT

	return order
}

func (o *Order) HasSameID(other *Order) bool {
	if o.PermID != 0 && other.PermID != 0 {
		return o.PermID == other.PermID
	}
	return o.OrderID == other.OrderID && o.ClientID == other.ClientID
}

func (o Order) String() string {
	s := fmt.Sprintf("%s, %s, %s: %s %s %s@%s %s",
		IntMaxString(o.OrderID),
		IntMaxString(o.ClientID),
		LongMaxString(o.PermID),
		o.OrderType,
		o.Action,
		DecimalMaxString(o.TotalQuantity),
		FloatMaxString(o.LmtPrice),
		o.TIF,
	)

	if len(o.OrderComboLegs) > 0 {
		s += " CMB("
		for _, leg := range o.OrderComboLegs {
			s += fmt.Sprintf("%s,", leg)
		}
		s += ")"
	}
	if len(o.Conditions) > 0 {
		s += " COND("
		for _, cond := range o.Conditions {
			s += fmt.Sprintf("%s,", cond)
		}
		s += ")"
	}

	return s
}
