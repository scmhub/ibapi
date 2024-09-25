package ibapi

import "strings"

// AccountSummaryTags .
type AccountSummaryTags = string

const (
	AccountType                 AccountSummaryTags = "AccountType"
	NetLiquidation              AccountSummaryTags = "NetLiquidation"
	TotalCashValue              AccountSummaryTags = "TotalCashValue"
	SettledCash                 AccountSummaryTags = "SettledCash"
	AccruedCash                 AccountSummaryTags = "AccruedCash"
	BuyingPower                 AccountSummaryTags = "BuyingPower"
	EquityWithLoanValue         AccountSummaryTags = "EquityWithLoanValue"
	PreviousEquityWithLoanValue AccountSummaryTags = "PreviousEquityWithLoanValue"
	GrossPositionValue          AccountSummaryTags = "GrossPositionValue"
	ReqTEquity                  AccountSummaryTags = "ReqTEquity"
	ReqTMargin                  AccountSummaryTags = "ReqTMargin"
	SMA                         AccountSummaryTags = "SMA"
	InitMarginReq               AccountSummaryTags = "InitMarginReq"
	MaintMarginReq              AccountSummaryTags = "MaintMarginReq"
	AvailableFunds              AccountSummaryTags = "AvailableFunds"
	ExcessLiquidity             AccountSummaryTags = "ExcessLiquidity"
	Cushion                     AccountSummaryTags = "Cushion"
	FullInitMarginReq           AccountSummaryTags = "FullInitMarginReq"
	FullMaintMarginReq          AccountSummaryTags = "FullMaintMarginReq"
	FullAvailableFunds          AccountSummaryTags = "FullAvailableFunds"
	FullExcessLiquidity         AccountSummaryTags = "FullExcessLiquidity"
	LookAheadNextChange         AccountSummaryTags = "LookAheadNextChange"
	LookAheadInitMarginReq      AccountSummaryTags = "LookAheadInitMarginReq"
	LookAheadMaintMarginReq     AccountSummaryTags = "LookAheadMaintMarginReq"
	LookAheadAvailableFunds     AccountSummaryTags = "LookAheadAvailableFunds"
	LookAheadExcessLiquidity    AccountSummaryTags = "LookAheadExcessLiquidity"
	HighestSeverity             AccountSummaryTags = "HighestSeverity"
	DayTradesRemaining          AccountSummaryTags = "DayTradesRemaining"
	Leverage                    AccountSummaryTags = "Leverage"
)

func GetAllTags() string {
	tags := []AccountSummaryTags{
		AccountType,
		NetLiquidation,
		TotalCashValue,
		SettledCash,
		AccruedCash,
		BuyingPower,
		EquityWithLoanValue,
		PreviousEquityWithLoanValue,
		GrossPositionValue,
		ReqTEquity,
		ReqTMargin,
		SMA,
		InitMarginReq,
		MaintMarginReq,
		AvailableFunds,
		ExcessLiquidity,
		Cushion,
		FullInitMarginReq,
		FullMaintMarginReq,
		FullAvailableFunds,
		FullExcessLiquidity,
		LookAheadNextChange,
		LookAheadInitMarginReq,
		LookAheadMaintMarginReq,
		LookAheadAvailableFunds,
		LookAheadExcessLiquidity,
		HighestSeverity,
		DayTradesRemaining,
		Leverage,
	}
	return strings.Join(tags, ",")
}
