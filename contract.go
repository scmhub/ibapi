package ibapi

import (
	"fmt"
)

// LegOpenClose .
type LegOpenClose int64

const (
	SAME_POS    LegOpenClose = 0
	OPEN_POS    LegOpenClose = 1
	CLOSE_POS   LegOpenClose = 2
	UNKNOWN_POS LegOpenClose = 3
)

// ComboLeg .
type ComboLeg struct {
	ConID     int64
	Ratio     int64
	Action    string // BUY/SELL/SSHORT
	Exchange  string
	OpenClose int64
	// for stock legs when doing short sale
	ShortSaleSlot      int64 // 1 = clearing broker, 2 = third party
	DesignatedLocation string
	ExemptCode         int64
}

// NewComboLeg creates a default ComboLeg.
func NewComboLeg() ComboLeg {
	cl := ComboLeg{}
	cl.ExemptCode = -1
	return cl
}

func (c ComboLeg) String() string {
	return fmt.Sprintf("%d, %d, %s, %s, %d, %d, %s, %d",
		c.ConID, c.Ratio, c.Action, c.Exchange, c.OpenClose, c.ShortSaleSlot, c.DesignatedLocation, c.ExemptCode)
}

// DeltaNeutralContract .
type DeltaNeutralContract struct {
	ConID int64
	Delta float64
	Price float64
}

func NewDeltaNeutralContract() DeltaNeutralContract {
	return DeltaNeutralContract{}
}

func (c DeltaNeutralContract) String() string {
	return fmt.Sprintf("%d, %f, %f", c.ConID, c.Delta, c.Price)
}

// Contract describes an instrument's definition.
type Contract struct {
	ConID                        int64
	Symbol                       string
	SecType                      string
	LastTradeDateOrContractMonth string
	LastTradeDate                string
	Strike                       float64 // UNSET_FLOAT
	Right                        string
	Multiplier                   string
	Exchange                     string
	PrimaryExchange              string // pick an actual (ie non-aggregate) exchange that the contract trades on.  DO NOT SET TO SMART.
	Currency                     string
	LocalSymbol                  string
	TradingClass                 string
	IncludeExpired               bool
	SecIDType                    string // CUSIP;SEDOL;ISIN;RIC
	SecID                        string
	Description                  string
	IssuerID                     string

	// combo legs
	ComboLegsDescrip string // received in open order 14 and up for all combos
	ComboLegs        []ComboLeg

	// delta neutral contract
	DeltaNeutralContract *DeltaNeutralContract
}

func NewContract() *Contract {
	return &Contract{
		Strike: UNSET_FLOAT,
	}
}

func (c *Contract) Equal(other *Contract) bool {
	if c.ConID != 0 && other.ConID != 0 {
		return c.ConID == other.ConID
	}
	if c.SecIDType != "" && other.SecIDType != "" && c.SecIDType == other.SecIDType {
		return c.SecID == other.SecID
	}
	return c.Symbol == other.Symbol &&
		c.SecType == other.SecType &&
		c.Exchange == other.Exchange &&
		c.Currency == other.Currency &&
		c.LastTradeDate == other.LastTradeDate &&
		c.Strike == other.Strike &&
		c.Right == other.Right
}

func (c Contract) String() string {
	s := fmt.Sprintf("%d, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %t, %s, %s, %s, %s",
		c.ConID,
		c.Symbol,
		c.SecType,
		c.LastTradeDateOrContractMonth,
		c.LastTradeDate,
		FloatMaxString(c.Strike),
		c.Right,
		c.Multiplier,
		c.Exchange,
		c.PrimaryExchange,
		c.Currency,
		c.LocalSymbol,
		c.TradingClass,
		c.IncludeExpired,
		c.SecIDType,
		c.SecID,
		c.Description,
		c.IssuerID,
	)
	if len(c.ComboLegs) > 1 {
		s += ", combo:" + c.ComboLegsDescrip
		for _, leg := range c.ComboLegs {
			s += fmt.Sprintf("; %s", leg)
		}
	}

	if c.DeltaNeutralContract != nil {
		s += fmt.Sprintf("; %s", c.DeltaNeutralContract)
	}

	return s
}

// ContractDetails .
type ContractDetails struct {
	Contract               Contract
	MarketName             string
	MinTick                float64
	OrderTypes             string
	ValidExchanges         string
	PriceMagnifier         int64
	UnderConID             int64
	LongName               string
	ContractMonth          string
	Industry               string
	Category               string
	Subcategory            string
	TimeZoneID             string
	TradingHours           string
	LiquidHours            string
	EVRule                 string
	EVMultiplier           int64
	AggGroup               int64
	UnderSymbol            string
	UnderSecType           string
	MarketRuleIDs          string
	RealExpirationDate     string
	LastTradeTime          string
	StockType              string
	MinSize                Decimal
	SizeIncrement          Decimal
	SuggestedSizeIncrement Decimal

	SecIDList []TagValue

	// BOND values
	Cusip             string
	Ratings           string
	DescAppend        string
	BondType          string
	CouponType        string
	Callable          bool
	Putable           bool
	Coupon            float64
	Convertible       bool
	Maturity          string
	IssueDate         string
	NextOptionDate    string
	NextOptionType    string
	NextOptionPartial bool
	Notes             string

	// FUND values
	FundName                        string
	FundFamily                      string
	FundType                        string
	FundFrontLoad                   string
	FundBackLoad                    string
	FundBackLoadTimeInterval        string
	FundManagementFee               string
	FundClosed                      bool
	FundClosedForNewInvestors       bool
	FundClosedForNewMoney           bool
	FundNotifyAmount                string
	FundMinimumInitialPurchase      string
	FundSubsequentMinimumPurchase   string
	FundBlueSkyStates               string
	FundBlueSkyTerritories          string
	FundDistributionPolicyIndicator FundDistributionPolicyIndicator
	FundAssetType                   FundAssetType
	IneligibilityReasonList         []IneligibilityReason
}

func NewContractDetails() *ContractDetails {
	cd := &ContractDetails{}
	cd.MinSize = UNSET_DECIMAL
	cd.SizeIncrement = UNSET_DECIMAL
	cd.SuggestedSizeIncrement = UNSET_DECIMAL
	return cd
}

func (c ContractDetails) String() string {
	return fmt.Sprintf("%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %t, %t, %f, %t, %s, %s, %s, %s, %t, %s, %s, %s, %s",
		c.Contract,
		c.MarketName,
		FloatMaxString(c.MinTick),
		c.OrderTypes,
		c.ValidExchanges,
		IntMaxString(c.PriceMagnifier),
		IntMaxString(c.UnderConID),
		c.LongName,
		c.ContractMonth,
		c.Industry,
		c.Category,
		c.Subcategory,
		c.TimeZoneID,
		c.TradingHours,
		c.LiquidHours,
		c.EVRule,
		IntMaxString(c.EVMultiplier),
		c.UnderSymbol,
		c.UnderSecType,
		c.MarketRuleIDs,
		IntMaxString(c.AggGroup),
		c.SecIDList,
		c.RealExpirationDate,
		c.StockType,
		// Bond
		c.Cusip,
		c.Ratings,
		c.DescAppend,
		c.BondType,
		c.CouponType,
		c.Callable,
		c.Putable,
		c.Coupon,
		c.Convertible,
		c.Maturity,
		c.IssueDate,
		c.NextOptionDate,
		c.NextOptionType,
		c.NextOptionPartial,
		c.Notes,
		DecimalMaxString(c.MinSize),
		DecimalMaxString(c.SizeIncrement),
		DecimalMaxString(c.SuggestedSizeIncrement),
	)
}

// ContractDescription includes contract and DerivativeSecTypes.
type ContractDescription struct {
	Contract           Contract
	DerivativeSecTypes []string
}

func NewContractDescription() ContractDescription {
	return ContractDescription{}
}

type FundAssetType string

const (
	FundAssetTypeNone        FundAssetType = ""
	FundAssetTypeOthers      FundAssetType = "000"
	FundAssetTypeMoneyMarket FundAssetType = "001"
	FundAssetTypeFixedIncome FundAssetType = "002"
	FundAssetTypeMultiAsset  FundAssetType = "003"
	FundAssetTypeEquity      FundAssetType = "004"
	FundAssetTypeSector      FundAssetType = "005"
	FundAssetTypeGuaranteed  FundAssetType = "006"
	FundAssetTypeAlternative FundAssetType = "007"
)

func getFundAssetType(fat string) FundAssetType {
	switch fat {
	case string(FundAssetTypeOthers):
		return FundAssetTypeOthers
	case string(FundAssetTypeMoneyMarket):
		return FundAssetTypeMoneyMarket
	case string(FundAssetTypeFixedIncome):
		return FundAssetTypeFixedIncome
	case string(FundAssetTypeMultiAsset):
		return FundAssetTypeMultiAsset
	case string(FundAssetTypeEquity):
		return FundAssetTypeEquity
	case string(FundAssetTypeSector):
		return FundAssetTypeSector
	case string(FundAssetTypeGuaranteed):
		return FundAssetTypeGuaranteed
	case string(FundAssetTypeAlternative):
		return FundAssetTypeAlternative
	default:
		return FundAssetTypeNone
	}
}

type FundDistributionPolicyIndicator string

const (
	FundDistributionPolicyIndicatorNone             FundDistributionPolicyIndicator = ""
	FundDistributionPolicyIndicatorAccumulationFund FundDistributionPolicyIndicator = "N"
	FundDistributionPolicyIndicatorIncomeFund       FundDistributionPolicyIndicator = "Y"
)

func getFundDistributionPolicyIndicator(fat string) FundDistributionPolicyIndicator {
	switch fat {
	case string(FundDistributionPolicyIndicatorAccumulationFund):
		return FundDistributionPolicyIndicatorAccumulationFund
	case string(FundDistributionPolicyIndicatorIncomeFund):
		return FundDistributionPolicyIndicatorIncomeFund
	default:
		return FundDistributionPolicyIndicatorNone
	}
}
