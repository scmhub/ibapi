package ibapi

import "fmt"

// var _ OrderCondition = (*MarginCondition)(nil)
// var _ ContractCondition = (*MarginCondition)(nil)
var _ OperatorCondition = (*MarginCondition)(nil)

type MarginCondition struct {
	*operatorCondition
	Percent int64
}

func newMarginCondition() *MarginCondition {
	return &MarginCondition{operatorCondition: newOperatorCondition(MarginOrderCondition)}
}

func (mc *MarginCondition) decode(msgBuf *MsgBuffer) {
	mc.operatorCondition.decode(msgBuf)
	mc.Percent = msgBuf.decodeInt64()
}

func (mc MarginCondition) makeFields() []any {
	return append(mc.operatorCondition.makeFields(), mc.Percent)
}

func (mc MarginCondition) String() string {
	percent := fmt.Sprintf("%d", mc.Percent)
	return fmt.Sprintf("the margin cushion persent %s", mc.operatorCondition.stringWithOperator(percent))
}
