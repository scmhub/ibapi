package ibapi

import "fmt"

var _ OrderCondition = (*ExecutionCondition)(nil)

type ExecutionCondition struct {
	*orderCondition
	SecType  string
	Exchange string
	Symbol   string
}

func newExecutionCondition() *ExecutionCondition {
	return &ExecutionCondition{orderCondition: newOrderCondition(ExecutionOrderCondition)}
}

func (ec *ExecutionCondition) decode(msgBuf *MsgBuffer) {
	ec.orderCondition.decode(msgBuf)
	ec.SecType = msgBuf.decodeString()
	ec.Exchange = msgBuf.decodeString()
	ec.Symbol = msgBuf.decodeString()
}

func (ec ExecutionCondition) makeFields() []any {
	return append(ec.orderCondition.makeFields(), ec.SecType, ec.Exchange, ec.Symbol)
}

func (ec ExecutionCondition) String() string {
	return fmt.Sprintf("trade occurs for %v symbol on %v exchange for %v security type", ec.Symbol, ec.Exchange, ec.SecType)
}
