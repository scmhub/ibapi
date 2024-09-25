package ibapi

import (
	"fmt"
	"strconv"
)

// Execution is the information of an order`s execution.
type Execution struct {
	ExecID               string
	Time                 string
	AcctNumber           string
	Exchange             string
	Side                 string
	Shares               Decimal
	Price                float64
	PermID               int64
	ClientID             int64
	OrderID              int64
	Liquidation          int64
	CumQty               Decimal
	AvgPrice             float64
	OrderRef             string
	EVRule               string
	EVMultiplier         float64
	ModelCode            string
	LastLiquidity        int64
	PendingPriceRevision bool
}

func (e Execution) String() string {
	return fmt.Sprintf("ExecId: %s, Time: %s, Account: %s, Exchange: %s, Side: %s, Shares: %s, Price: %s, PermId: %s, ClientId: %s, OrderId: %s, Liquidation: %s, CumQty: %s, AvgPrice: %s, OrderRef: %s, EvRule: %s, EvMultiplier: %s, ModelCode: %s, LastLiquidity: %s,  PendingPriceRevision: %s",
		e.ExecID, e.Time, e.AcctNumber, e.Exchange, e.Side, decimalMaxString(e.Shares), floatMaxString(e.Price), intMaxString(e.PermID), intMaxString(e.ClientID), intMaxString(e.OrderID), intMaxString(e.Liquidation), decimalMaxString(e.CumQty), floatMaxString(e.AvgPrice),
		e.OrderRef, e.EVRule, floatMaxString(e.EVMultiplier), e.ModelCode, intMaxString(e.LastLiquidity), strconv.FormatBool(e.PendingPriceRevision))
}

func NewExecution() *Execution {
	e := &Execution{}
	e.Shares = UNSET_DECIMAL
	e.CumQty = UNSET_DECIMAL
	return e
}

// ExecutionFilter .
type ExecutionFilter struct {
	ClientID int64
	AcctCode string
	Time     string
	Symbol   string
	SecType  string
	Exchange string
	Side     string
}
