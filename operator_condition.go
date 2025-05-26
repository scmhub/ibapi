package ibapi

import "fmt"

// OperatorCondition embeds OrderCondition and requires valueConverter
type OperatorCondition interface {
	OrderCondition
}

var _ OrderCondition = (*operatorCondition)(nil)

type operatorCondition struct {
	*orderCondition
	IsMore bool
}

func newOperatorCondition(orderConditionType OrderConditionType) *operatorCondition {
	return &operatorCondition{orderCondition: newOrderCondition(orderConditionType)}
}

func (oc *operatorCondition) decode(msgBuf *MsgBuffer) { // 2 fields
	oc.orderCondition.decode(msgBuf)
	oc.IsMore = msgBuf.decodeBool()
}

func (oc operatorCondition) makeFields() []any {
	return append(oc.orderCondition.makeFields(), oc.IsMore)
}

func (oc operatorCondition) stringWithOperator(value string) string {
	if oc.IsMore {
		return fmt.Sprintf("is >= %s", value)
	}
	return fmt.Sprintf("is <= %s", value)
}
