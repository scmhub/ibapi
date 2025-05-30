package ibapi

import (
	"encoding/json"
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
	ConID     int64  `json:"conId,omitempty"`
	Ratio     int64  `json:"ratio,omitempty"`
	Action    string `json:"action,omitempty"` // BUY/SELL/SSHORT
	Exchange  string `json:"exchange,omitempty"`
	OpenClose int64  `json:"openClose,omitempty"`
	// for stock legs when doing short sale
	ShortSaleSlot      int64  `json:"shortSaleSlot,omitempty"` // 1 = clearing broker, 2 = third party
	DesignatedLocation string `json:"designatedLocation,omitempty"`
	ExemptCode         int64  `json:"exemptCode,omitempty"` // Default is -1
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

// MarshalJSON implements the json.Marshaler interface
func (c ComboLeg) MarshalJSON() ([]byte, error) {
	type Alias ComboLeg
	aux := &struct {
		*Alias
		ExemptCode *int64 `json:"exemptCode,omitempty"`
	}{
		Alias: (*Alias)(&c),
	}

	if c.ExemptCode != -1 {
		aux.ExemptCode = &c.ExemptCode
	}

	return json.Marshal(aux)
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (c *ComboLeg) UnmarshalJSON(data []byte) error {
	type Alias ComboLeg
	aux := &struct {
		*Alias
		ExemptCode *int64 `json:"exemptCode,omitempty"`
	}{
		Alias: (*Alias)(c),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.ExemptCode == nil {
		c.ExemptCode = -1
	} else {
		c.ExemptCode = *aux.ExemptCode
	}

	return nil
}

// DeltaNeutralContract .
type DeltaNeutralContract struct {
	ConID int64   `json:"conId,omitempty"`
	Delta float64 `json:"delta,omitempty"`
	Price float64 `json:"price,omitempty"`
}

func NewDeltaNeutralContract() DeltaNeutralContract {
	return DeltaNeutralContract{}
}

func (c DeltaNeutralContract) String() string {
	return fmt.Sprintf("%d, %f, %f", c.ConID, c.Delta, c.Price)
}

// Contract describes an instrument's definition.
type Contract struct {
	ConID                        int64                 `json:"conId,omitempty"`
	Symbol                       string                `json:"symbol,omitempty"`
	SecType                      string                `json:"secType,omitempty"`
	LastTradeDateOrContractMonth string                `json:"lastTradeDateOrContractMonth,omitempty"`
	LastTradeDate                string                `json:"lastTradeDate,omitempty"`
	Strike                       float64               `json:"strike"` // UNSET_FLOAT so no omitempty
	Right                        string                `json:"right,omitempty"`
	Multiplier                   string                `json:"multiplier,omitempty"`
	Exchange                     string                `json:"exchange,omitempty"`
	PrimaryExchange              string                `json:"primaryExchange,omitempty"` // pick an actual (ie non-aggregate) exchange that the contract trades on.  DO NOT SET TO SMART.
	Currency                     string                `json:"currency,omitempty"`
	LocalSymbol                  string                `json:"localSymbol,omitempty"`
	TradingClass                 string                `json:"tradingClass,omitempty"`
	IncludeExpired               bool                  `json:"includeExpired,omitempty"`
	SecIDType                    string                `json:"secIdType,omitempty"` // CUSIP;SEDOL;ISIN;RIC
	SecID                        string                `json:"secId,omitempty"`
	Description                  string                `json:"description,omitempty"`
	IssuerID                     string                `json:"issuerId,omitempty"`
	ComboLegsDescrip             string                `json:"comboLegsDescrip,omitempty"` // received in open order 14 and up for all combos
	ComboLegs                    []ComboLeg            `json:"comboLegs,omitempty"`
	DeltaNeutralContract         *DeltaNeutralContract `json:"deltaNeutralContract,omitempty"`
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
	if c.Symbol == other.Symbol &&
		c.SecType == other.SecType &&
		c.Exchange == other.Exchange &&
		c.Currency == other.Currency &&
		c.LastTradeDate == other.LastTradeDate &&
		c.Strike == other.Strike &&
		c.Right == other.Right {

		// Compare ComboLegs
		if len(c.ComboLegs) != len(other.ComboLegs) {
			return false
		}
		for i := range c.ComboLegs {
			if c.ComboLegs[i] != other.ComboLegs[i] {
				return false
			}
		}

		// Compare DeltaNeutralContract
		if (c.DeltaNeutralContract == nil) != (other.DeltaNeutralContract == nil) {
			return false
		}
		if c.DeltaNeutralContract != nil {
			return c.DeltaNeutralContract.ConID == other.DeltaNeutralContract.ConID &&
				c.DeltaNeutralContract.Delta == other.DeltaNeutralContract.Delta &&
				c.DeltaNeutralContract.Price == other.DeltaNeutralContract.Price
		}
		return true
	}
	return false
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

// MarshalJSON implements the json.Marshaler interface
func (c Contract) MarshalJSON() ([]byte, error) {
	type Alias Contract // prevent recursive call to MarshalJSON
	aux := &struct {
		*Alias
		Strike *float64 `json:"strike,omitempty"` // use pointer to handle UNSET_FLOAT
	}{
		Alias: (*Alias)(&c),
	}

	if c.Strike != UNSET_FLOAT {
		aux.Strike = &c.Strike
	}

	return json.Marshal(aux)
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (c *Contract) UnmarshalJSON(data []byte) error {
	type Alias Contract // prevent recursive call to UnmarshalJSON
	aux := &struct {
		*Alias
		Strike *float64 `json:"strike,omitempty"`
	}{
		Alias: (*Alias)(c),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.Strike == nil {
		c.Strike = UNSET_FLOAT
	} else {
		c.Strike = *aux.Strike
	}

	return nil
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
