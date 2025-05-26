package ibapi

import (
	"fmt"
	"strconv"
)

// var _ OrderCondition = (*VolumeCondition)(nil)
var _ ContractCondition = (*VolumeCondition)(nil)

type VolumeCondition struct {
	*contractCondition
	Volume int64
}

func newVolumeCondition() *VolumeCondition {
	return &VolumeCondition{contractCondition: newContractCondition(VolumeOrderCondition)}
}

func (vc *VolumeCondition) decode(msgBuf *MsgBuffer) {
	vc.contractCondition.decode(msgBuf)
	vc.Volume = msgBuf.decodeInt64()
}

func (vc VolumeCondition) makeFields() []any {
	return append(vc.contractCondition.makeFields(), vc.Volume)
}

func (vc VolumeCondition) String() string {
	volume := strconv.FormatInt(vc.Volume, 10)
	return fmt.Sprintf("%s %s", vc.contractCondition, vc.operatorCondition.stringWithOperator(volume))
}
