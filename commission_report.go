package ibapi

import "fmt"

// CommissionReport .
type CommissionReport struct {
	ExecID              string
	Commission          float64
	Currency            string
	RealizedPNL         float64
	Yield               float64
	YieldRedemptionDate int64 // YYYYMMDD format
}

func NewCommissionReport() CommissionReport {
	return CommissionReport{}
}

func (cr CommissionReport) String() string {
	return fmt.Sprintf("ExecId: %s, Commission:%f, Currency: %s, RealizedPnL: %f, Yield: %f, YieldRedemptionDate: %d",
		cr.ExecID,
		cr.Commission,
		cr.Currency,
		cr.RealizedPNL,
		cr.Yield,
		cr.YieldRedemptionDate)
}
