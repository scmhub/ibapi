package ibapi

import "fmt"

// var _ OrderCondition = (*TimeCondition)(nil)
// var _ ContractCondition = (*TimeCondition)(nil)
var _ OperatorCondition = (*TimeCondition)(nil)

type TimeCondition struct {
	*operatorCondition
	Time string
}

func newTimeCondition() *TimeCondition {
	return &TimeCondition{operatorCondition: newOperatorCondition(TimeOrderCondition)}
}

func (tc *TimeCondition) decode(msgBuf *MsgBuffer) {
	tc.operatorCondition.decode(msgBuf)
	tc.Time = msgBuf.decodeString()
}

func (tc TimeCondition) makeFields() []any {
	return append(tc.operatorCondition.makeFields(), tc.Time)
}

func (tc TimeCondition) String() string {
	return fmt.Sprintf("time is %s", tc.operatorCondition.stringWithOperator(tc.Time))
}
