package ibapi

import (
	"fmt"
)

// Trigger Methods.
type TriggerMethod int64

const (
	DefaultTriggerMethod      TriggerMethod = 0
	DoubleBidAskTriggerMethod TriggerMethod = 1
	LastTriggerMethod         TriggerMethod = 2
	DoubleLastTriggerMethod   TriggerMethod = 3
	BidAskTriggerMethod       TriggerMethod = 4
	LastBidAskTriggerMethod   TriggerMethod = 7
	MidPointTriggerMethod     TriggerMethod = 8
)

func (tm TriggerMethod) String() string {
	switch tm {
	case DefaultTriggerMethod:
		return "default"
	case DoubleBidAskTriggerMethod:
		return "Double BidAsk"
	case LastTriggerMethod:
		return "last"
	case DoubleLastTriggerMethod:
		return "double Last"
	case BidAskTriggerMethod:
		return "Bid/ssk"
	case LastBidAskTriggerMethod:
		return "last of Bid/Ask"
	case MidPointTriggerMethod:
		return "mid-point"
	default:
		return fmt.Sprintf("Unknown Trigger Method %d", tm)
	}
}

// var _ OrderCondition = (*PriceCondition)(nil)
var _ ContractCondition = (*PriceCondition)(nil)

type PriceCondition struct {
	*contractCondition
	Price         float64
	TriggerMethod TriggerMethod
}

func newPriceCondition() *PriceCondition {
	return &PriceCondition{contractCondition: newContractCondition(PriceOrderCondition)}
}

func (pc *PriceCondition) decode(msgBuf *MsgBuffer) {
	pc.contractCondition.decode(msgBuf)
	pc.Price = msgBuf.decodeFloat64()
	pc.TriggerMethod = TriggerMethod(msgBuf.decodeInt64())
}

func (pc PriceCondition) makeFields() []any {
	return append(pc.contractCondition.makeFields(), pc.Price, pc.TriggerMethod)
}

func (pc PriceCondition) String() string {
	return fmt.Sprintf("%v %v", pc.TriggerMethod, pc.contractCondition)
}
