package ibapi

import (
	"fmt"
)

type OrderConditionType = int64

const (
	PriceOrderCondition         OrderConditionType = 1
	TimeOrderCondition          OrderConditionType = 3
	MarginOrderCondition        OrderConditionType = 4
	ExecutionOrderCondition     OrderConditionType = 5
	VolumeOrderCondition        OrderConditionType = 6
	PercentChangeOrderCondition OrderConditionType = 7
)

// Trigger Methods.
const (
	DefaultTriggerMethod = iota
	DoubleBidAskTriggerMethod
	LastTriggerMethod
	DoubleLastTriggerMethod
	BidAskTriggerMethod
	LastBidAskTriggerMethod
	MidPointTriggerMethod
)

type OrderCondition interface {
	Type() OrderConditionType
	decode(*MsgBuffer)
	makeFields() []any
}

var _ OrderCondition = (*orderCondition)(nil)

type orderCondition struct {
	condType                OrderConditionType
	IsConjunctionConnection bool
}

func (oc orderCondition) Type() OrderConditionType {
	return oc.condType
}

func (oc *orderCondition) decode(msgBuf *MsgBuffer) {
	connector := msgBuf.decodeString()
	oc.IsConjunctionConnection = connector == "a"
}

func (oc orderCondition) makeFields() []any {
	if oc.IsConjunctionConnection {
		return []any{"a"}
	}
	return []any{"o"}
}

func (oc orderCondition) String() string {
	if oc.IsConjunctionConnection {
		return "<AND>"
	}
	return "<OR>"
}

var _ OrderCondition = (*ExecutionCondition)(nil)

type ExecutionCondition struct {
	*orderCondition
	SecType  string
	Exchange string
	Symbol   string
}

func (ec *ExecutionCondition) decode(msgBuf *MsgBuffer) { // 4 fields
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

var _ OrderCondition = (*OperatorCondition)(nil)

type OperatorCondition struct {
	*orderCondition
	IsMore bool
}

func (oc *OperatorCondition) decode(msgBuf *MsgBuffer) { // 2 fields
	oc.orderCondition.decode(msgBuf)
	oc.IsMore = msgBuf.decodeBool()
}

func (oc OperatorCondition) makeFields() []any {
	return append(oc.orderCondition.makeFields(), oc.IsMore)
}

type MarginCondition struct {
	*OperatorCondition
	Percent int64
}

func (mc *MarginCondition) decode(msgBuf *MsgBuffer) { // 3 fields
	mc.OperatorCondition.decode(msgBuf)
	mc.Percent = msgBuf.decodeInt64()
}

func (mc MarginCondition) makeFields() []any {
	return append(mc.OperatorCondition.makeFields(), mc.Percent)
}

type TimeCondition struct {
	*OperatorCondition
	Time string
}

func (tc *TimeCondition) decode(msgBuf *MsgBuffer) { // 3 fields
	tc.OperatorCondition.decode(msgBuf)
	// tc.Time = decodeTime(fields[2], "20060102")
	tc.Time = msgBuf.decodeString()
}

func (tc TimeCondition) makeFields() []any {
	return append(tc.OperatorCondition.makeFields(), tc.Time)
}

type ContractCondition struct {
	*OperatorCondition
	ConID    int64
	Exchange string
}

func (cc *ContractCondition) decode(msgBuf *MsgBuffer) { // 4 fields
	cc.OperatorCondition.decode(msgBuf)
	cc.ConID = msgBuf.decodeInt64()
	cc.Exchange = msgBuf.decodeString()
}

func (cc ContractCondition) makeFields() []any {
	return append(cc.OperatorCondition.makeFields(), cc.ConID, cc.Exchange)
}

var _ OrderCondition = (*PriceCondition)(nil)

type PriceCondition struct {
	*ContractCondition
	Price         float64
	TriggerMethod int64
}

func (pc *PriceCondition) decode(msgBuf *MsgBuffer) { // 6 fields
	pc.ContractCondition.decode(msgBuf)
	pc.Price = msgBuf.decodeFloat64()
	pc.TriggerMethod = msgBuf.decodeInt64()
}

func (pc PriceCondition) makeFields() []any {
	return append(pc.ContractCondition.makeFields(), pc.Price, pc.TriggerMethod)
}

type VolumeCondition struct {
	*ContractCondition
	Volume int64
}

func (vc *VolumeCondition) decode(msgBuf *MsgBuffer) { // 5 fields
	vc.ContractCondition.decode(msgBuf)
	vc.Volume = msgBuf.decodeInt64()
}

func (vc VolumeCondition) makeFields() []any {
	return append(vc.ContractCondition.makeFields(), vc.Volume)
}

type PercentChangeCondition struct {
	*ContractCondition
	ChangePercent float64
}

func (pcc *PercentChangeCondition) decode(msgBuf *MsgBuffer) { // 5 fields
	pcc.ContractCondition.decode(msgBuf)
	pcc.ChangePercent = msgBuf.decodeFloat64()
}

func (pcc PercentChangeCondition) makeFields() []any {
	return append(pcc.ContractCondition.makeFields(), pcc.ChangePercent)
}

func CreateOrderCondition(condType OrderConditionType) OrderCondition {
	var cond OrderCondition
	switch condType {
	case PriceOrderCondition:
		cond = &PriceCondition{
			ContractCondition: &ContractCondition{
				OperatorCondition: &OperatorCondition{
					orderCondition: &orderCondition{condType: PriceOrderCondition},
				},
			},
		}
	case TimeOrderCondition:
		cond = &TimeCondition{
			OperatorCondition: &OperatorCondition{
				orderCondition: &orderCondition{condType: TimeOrderCondition},
			},
		}
	case MarginOrderCondition:
		cond = &MarginCondition{
			OperatorCondition: &OperatorCondition{
				orderCondition: &orderCondition{condType: MarginOrderCondition},
			},
		}
	case ExecutionOrderCondition:
		cond = &ExecutionCondition{orderCondition: &orderCondition{condType: ExecutionOrderCondition}}
	case VolumeOrderCondition:
		cond = &VolumeCondition{
			ContractCondition: &ContractCondition{
				OperatorCondition: &OperatorCondition{
					orderCondition: &orderCondition{condType: VolumeOrderCondition},
				},
			},
		}
	case PercentChangeOrderCondition:
		cond = &PercentChangeCondition{
			ContractCondition: &ContractCondition{
				OperatorCondition: &OperatorCondition{
					orderCondition: &orderCondition{condType: PercentChangeOrderCondition},
				},
			},
		}
	default:
		log.Panic().Msg("unknown OrderConditionType")
	}
	return cond
}
