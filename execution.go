package ibapi

import (
	"fmt"
	"strconv"
)

// Execution is the information of an order`s execution.
type Execution struct {
	ExecID                 string
	Time                   string
	AcctNumber             string
	Exchange               string
	Side                   string
	Shares                 Decimal //UNSET_DECIMAL
	Price                  float64
	PermID                 int64
	ClientID               int64
	OrderID                int64
	Liquidation            int64
	CumQty                 Decimal // UNSET_DECIMAL
	AvgPrice               float64
	OrderRef               string
	EVRule                 string
	EVMultiplier           float64
	ModelCode              string
	LastLiquidity          int64
	PendingPriceRevision   bool
	Submitter              string
	OptExerciseOrLapseType OptionExerciseType
}

func (e Execution) String() string {
	return fmt.Sprintf("ExecId: %s, Time: %s, Account: %s, Exchange: %s, Side: %s, Shares: %s, Price: %s, PermId: %s, ClientId: %s, OrderId: %s, Liquidation: %s, CumQty: %s, AvgPrice: %s, OrderRef: %s, EvRule: %s, EvMultiplier: %s, ModelCode: %s, LastLiquidity: %s,  PendingPriceRevision: %s, Submitter: %s, OptionExerciseType: %s",
		e.ExecID, e.Time, e.AcctNumber, e.Exchange, e.Side, DecimalMaxString(e.Shares), FloatMaxString(e.Price), LongMaxString(e.PermID), IntMaxString(e.ClientID), IntMaxString(e.OrderID), IntMaxString(e.Liquidation), DecimalMaxString(e.CumQty), FloatMaxString(e.AvgPrice),
		e.OrderRef, e.EVRule, FloatMaxString(e.EVMultiplier), e.ModelCode, IntMaxString(e.LastLiquidity), strconv.FormatBool(e.PendingPriceRevision), e.Submitter, e.OptExerciseOrLapseType.String())
}

func NewExecution() *Execution {
	e := &Execution{}
	e.Shares = UNSET_DECIMAL
	e.CumQty = UNSET_DECIMAL
	e.OptExerciseOrLapseType = OptionExerciseTypeNone
	return e
}

// ExecutionFilter .
type ExecutionFilter struct {
	ClientID      int64
	AcctCode      string
	Time          string
	Symbol        string
	SecType       string
	Exchange      string
	Side          string
	LastNDays     int64
	SpecificDates []int64
}

func NewExecutionFilter() *ExecutionFilter {
	ef := &ExecutionFilter{}
	ef.LastNDays = UNSET_INT
	return ef
}
