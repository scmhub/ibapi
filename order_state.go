package ibapi

import "fmt"

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

	Commission         float64
	MinCommission      float64
	MaxCommission      float64
	CommissionCurrency string

	WarningText string

	CompletedTime   string
	CompletedStatus string
}

func NewOrderState() *OrderState {
	os := &OrderState{}
	os.Commission = UNSET_FLOAT
	os.MinCommission = UNSET_FLOAT
	os.MaxCommission = UNSET_FLOAT

	return os
}

func (os OrderState) String() string {
	return fmt.Sprintf("Status: %s, Commission: %s, Commission currency %s, CompletedTime: %s, CompletedStatus: %s",
		os.Status,
		floatMaxString(os.Commission),
		os.CommissionCurrency,
		os.CompletedTime,
		os.CompletedStatus)
}
