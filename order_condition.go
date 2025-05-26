package ibapi

/*
------------------------------------------------
------------------------------------------------
OrderCondition
------------------------------------------------
	OperatorCondition
		OrderCondition
------------------------------------------------
		ContractCondition
			OperatorCondition
------------------------------------------------
------------------------------------------------
		PriceOrderCondition
			ContractCondition
------------------------------------------------
	TimeOrderCondition
		OperatorCondition
------------------------------------------------
	MarginOrderCondition
		OperatorCondition
------------------------------------------------
ExecutionOrderCondition
	OrderCondition
------------------------------------------------
		VolumeOrderCondition
			ContractCondition
------------------------------------------------
		PercentChangeOrderCondition
			ContractCondition
------------------------------------------------
------------------------------------------------
*/

type OrderConditionType = int64

const (
	PriceOrderCondition         OrderConditionType = 1
	TimeOrderCondition          OrderConditionType = 3
	MarginOrderCondition        OrderConditionType = 4
	ExecutionOrderCondition     OrderConditionType = 5
	VolumeOrderCondition        OrderConditionType = 6
	PercentChangeOrderCondition OrderConditionType = 7
)

type OrderCondition interface {
	Type() OrderConditionType
	TypeName() string
	IsConjunctionConnection() bool
	SetIsConjunctionConnection(bool)
	decode(*MsgBuffer)
	makeFields() []any
}

var _ OrderCondition = (*orderCondition)(nil)

type orderCondition struct {
	condType                OrderConditionType
	isConjunctionConnection bool
}

func newOrderCondition(condType OrderConditionType) *orderCondition {
	return &orderCondition{
		condType:                condType,
		isConjunctionConnection: true,
	}
}

func (oc orderCondition) Type() OrderConditionType {
	return oc.condType
}

func (oc orderCondition) TypeName() string {
	switch oc.condType {
	case PriceOrderCondition:
		return "Price"
	case TimeOrderCondition:
		return "Time"
	case MarginOrderCondition:
		return "Margin"
	case ExecutionOrderCondition:
		return "Execution"
	case VolumeOrderCondition:
		return "Volume"
	case PercentChangeOrderCondition:
		return "PercentChange"
	default:
		return "Unknown"
	}
}

func (oc *orderCondition) IsConjunctionConnection() bool {
	return oc.isConjunctionConnection
}

func (oc *orderCondition) SetIsConjunctionConnection(isConjunctionConnection bool) {
	oc.isConjunctionConnection = isConjunctionConnection
}

func (oc *orderCondition) decode(msgBuf *MsgBuffer) {
	connector := msgBuf.decodeString()
	oc.isConjunctionConnection = connector == "a"
}

func (oc orderCondition) makeFields() []any {
	if oc.isConjunctionConnection {
		return []any{"a"}
	}
	return []any{"o"}
}

func (oc orderCondition) String() string {
	if oc.isConjunctionConnection {
		return "<AND>"
	}
	return "<OR>"
}

func CreateOrderCondition(condType OrderConditionType) OrderCondition {
	var cond OrderCondition
	switch condType {
	case PriceOrderCondition:
		cond = newPriceCondition()
	case TimeOrderCondition:
		cond = newTimeCondition()
	case MarginOrderCondition:
		cond = newMarginCondition()
	case ExecutionOrderCondition:
		cond = newExecutionCondition()
	case VolumeOrderCondition:
		cond = newVolumeCondition()
	case PercentChangeOrderCondition:
		cond = newPercentChangeCondition()
	default:
		log.Panic().Msg("unknown OrderConditionType")
	}
	return cond
}
