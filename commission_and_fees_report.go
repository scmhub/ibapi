package ibapi

import "fmt"

// CommissionAndFeesReport .
type CommissionAndFeesReport struct {
	ExecID              string
	CommissionAndFees   float64
	Currency            string
	RealizedPNL         float64
	Yield               float64
	YieldRedemptionDate int64 // YYYYMMDD format
}

func NewCommissionAndFeesReport() CommissionAndFeesReport {
	return CommissionAndFeesReport{}
}

func (cr CommissionAndFeesReport) String() string {
	return fmt.Sprintf("ExecId: %s, CommissionAndFees: %f, Currency: %s, RealizedPnL: %f, Yield: %f, YieldRedemptionDate: %d",
		cr.ExecID,
		cr.CommissionAndFees,
		cr.Currency,
		cr.RealizedPNL,
		cr.Yield,
		cr.YieldRedemptionDate)
}
