package ibapi

import (
	"fmt"
)

// var _ OrderCondition = (*PercentChangeCondition)(nil)
var _ ContractCondition = (*PercentChangeCondition)(nil)

type PercentChangeCondition struct {
	*contractCondition
	ChangePercent float64
}

func newPercentChangeCondition() *PercentChangeCondition {
	return &PercentChangeCondition{contractCondition: newContractCondition(PercentChangeOrderCondition)}
}

func (pcc *PercentChangeCondition) decode(msgBuf *MsgBuffer) {
	pcc.ChangePercent = msgBuf.decodeFloat64()
}

func (pcc PercentChangeCondition) makeFields() []any {
	return append(pcc.contractCondition.makeFields(), pcc.ChangePercent)
}

func (pcc PercentChangeCondition) String() string {
	volume := fmt.Sprintf("%f", pcc.ChangePercent)
	return fmt.Sprintf("%s %s", pcc.contractCondition, pcc.operatorCondition.stringWithOperator(volume))
}
