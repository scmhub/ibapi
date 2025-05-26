package ibapi

import "fmt"

type ContractCondition interface {
	OperatorCondition
}

var _ OperatorCondition = (*contractCondition)(nil)

type contractCondition struct {
	*operatorCondition
	ConID    int64
	Exchange string
}

func newContractCondition(condType OrderConditionType) *contractCondition {
	return &contractCondition{operatorCondition: newOperatorCondition(condType)}
}

func (cc *contractCondition) decode(msgBuf *MsgBuffer) {
	cc.operatorCondition.decode(msgBuf)
	cc.ConID = msgBuf.decodeInt64()
	cc.Exchange = msgBuf.decodeString()
}

func (cc contractCondition) makeFields() []any {
	return append(cc.operatorCondition.makeFields(), cc.ConID, cc.Exchange)
}

func (cc contractCondition) String() string {
	return fmt.Sprintf("%v of %v", cc.TypeName(), cc.ConID)
}
