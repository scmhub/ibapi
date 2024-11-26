package ibapi

import (
	"fmt"
)

// OrderAllocation .
type OrderAllocation struct {
	Account         string
	Position        Decimal // UNSET_DECIMAL
	PositionDesired Decimal // UNSET_DECIMAL
	PositionAfter   Decimal // UNSET_DECIMAL
	DesiredAllocQty Decimal // UNSET_DECIMAL
	AllowedAllocQty Decimal // UNSET_DECIMAL
	IsMonetary      bool
}

func NewOrderAllocation() *OrderAllocation {
	oa := &OrderAllocation{}
	oa.Position = UNSET_DECIMAL
	oa.PositionDesired = UNSET_DECIMAL
	oa.PositionAfter = UNSET_DECIMAL
	oa.DesiredAllocQty = UNSET_DECIMAL
	oa.AllowedAllocQty = UNSET_DECIMAL
	return oa
}

func (oa OrderAllocation) String() string {
	return fmt.Sprint(
		"Account: ", oa.Account,
		", Position: ", DecimalMaxString(oa.Position),
		", PositionDesired: ", DecimalMaxString(oa.PositionDesired),
		", PositionAfter: ", DecimalMaxString(oa.PositionAfter),
		", DesiredAllocQty: ", DecimalMaxString(oa.DesiredAllocQty),
		", AllowedAllocQty: ", DecimalMaxString(oa.AllowedAllocQty),
		", IsMonetary: ", oa.IsMonetary,
	)
}

// OrderState .
type OrderState struct {
	Status string

	InitMarginBefore     string
	MaintMarginBefore    string
	EquityWithLoanBefore string
	InitMarginChange     string
	MaintMarginChange    string
	EquityWithLoanChange string
	InitMarginAfter      string
	MaintMarginAfter     string
	EquityWithLoanAfter  string

	Commission                     float64 // UNSET_FLOAT
	MinCommission                  float64 // UNSET_FLOAT
	MaxCommission                  float64 // UNSET_FLOAT
	CommissionCurrency             string
	MarginCurrency                 string
	InitMarginBeforeOutsideRTH     float64 // UNSET_FLOAT
	MaintMarginBeforeOutsideRTH    float64 // UNSET_FLOAT
	EquityWithLoanBeforeOutsideRTH float64 // UNSET_FLOAT
	InitMarginChangeOutsideRTH     float64 // UNSET_FLOAT
	MaintMarginChangeOutsideRTH    float64 // UNSET_FLOAT
	EquityWithLoanChangeOutsideRTH float64 // UNSET_FLOAT
	InitMarginAfterOutsideRTH      float64 // UNSET_FLOAT
	MaintMarginAfterOutsideRTH     float64 // UNSET_FLOAT
	EquityWithLoanAfterOutsideRTH  float64 // UNSET_FLOAT
	SuggestedSize                  Decimal // UNSET_DECIMAL
	RejectReason                   string
	OrderAllocations               []*OrderAllocation
	WarningText                    string

	CompletedTime   string
	CompletedStatus string
}

func NewOrderState() *OrderState {
	os := &OrderState{}
	os.Commission = UNSET_FLOAT
	os.MinCommission = UNSET_FLOAT
	os.MaxCommission = UNSET_FLOAT
	os.InitMarginBeforeOutsideRTH = UNSET_FLOAT
	os.MaintMarginBeforeOutsideRTH = UNSET_FLOAT
	os.EquityWithLoanBeforeOutsideRTH = UNSET_FLOAT
	os.InitMarginChangeOutsideRTH = UNSET_FLOAT
	os.MaintMarginChangeOutsideRTH = UNSET_FLOAT
	os.EquityWithLoanChangeOutsideRTH = UNSET_FLOAT
	os.InitMarginAfterOutsideRTH = UNSET_FLOAT
	os.MaintMarginAfterOutsideRTH = UNSET_FLOAT
	os.EquityWithLoanAfterOutsideRTH = UNSET_FLOAT
	os.SuggestedSize = UNSET_DECIMAL

	return os
}

func (os OrderState) String() string {
	s := fmt.Sprint(
		"Status: ", os.Status,
		", InitMarginBefore: ", os.InitMarginBefore,
		", MaintMarginBefore: ", os.MaintMarginBefore,
		", EquityWithLoanBefore: ", os.EquityWithLoanBefore,
		", InitMarginChange: ", os.InitMarginChange,
		", MaintMarginChange: ", os.MaintMarginChange,
		", EquityWithLoanChange: ", os.EquityWithLoanChange,
		", InitMarginAfter: ", os.InitMarginAfter,
		", MaintMarginAfter: ", os.MaintMarginAfter,
		", EquityWithLoanAfter: ", os.EquityWithLoanAfter,
		", Commission: ", FloatMaxString(os.Commission),
		", MinCommission: ", FloatMaxString(os.MinCommission),
		", MaxCommission: ", FloatMaxString(os.MaxCommission),
		", CommissionCurrency: ", os.CommissionCurrency,
		", MarginCurrency: ", os.MarginCurrency,
		", InitMarginBeforeOutsideRTH: ", FloatMaxString(os.InitMarginBeforeOutsideRTH),
		", MaintMarginBeforeOutsideRTH: ", FloatMaxString(os.MaintMarginBeforeOutsideRTH),
		", EquityWithLoanBeforeOutsideRTH: ", FloatMaxString(os.EquityWithLoanBeforeOutsideRTH),
		", InitMarginChangeOutsideRTH: ", FloatMaxString(os.InitMarginChangeOutsideRTH),
		", MaintMarginChangeOutsideRTH: ", FloatMaxString(os.MaintMarginChangeOutsideRTH),
		", EquityWithLoanChangeOutsideRTH: ", FloatMaxString(os.EquityWithLoanChangeOutsideRTH),
		", InitMarginAfterOutsideRTH: ", FloatMaxString(os.InitMarginAfterOutsideRTH),
		", MaintMarginAfterOutsideRTH: ", FloatMaxString(os.MaintMarginAfterOutsideRTH),
		", EquityWithLoanAfterOutsideRTH: ", FloatMaxString(os.EquityWithLoanAfterOutsideRTH),
		", SuggestedSize: ", DecimalMaxString(os.SuggestedSize),
		", RejectReason: ", os.RejectReason,
		", WarningText: ", os.WarningText,
		", CompletedTime: ", os.CompletedTime,
		", CompletedStatus: ", os.CompletedStatus,
	)

	if os.OrderAllocations != nil {
		s += ", OrderAllocations: "
		for _, oa := range os.OrderAllocations {
			s += "[" + oa.String() + "]"
		}
	}

	return s
}
