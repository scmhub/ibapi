package ibapi

import (
	"testing"
)

func TestCreatePriceConditionProto(t *testing.T) {
	cond := NewPriceCondition(756733, "SMART", 0.01, LastTriggerMethod, false, true)

	proto := createPriceConditionProto(cond)

	// Type
	if proto.Type == nil || *proto.Type != int32(PriceOrderCondition) {
		t.Errorf("Type: got %v, want %d", proto.Type, PriceOrderCondition)
	}
	// IsConjunctionConnection
	if proto.IsConjunctionConnection == nil || *proto.IsConjunctionConnection != true {
		t.Errorf("IsConjunctionConnection: got %v, want true", proto.IsConjunctionConnection)
	}
	// IsMore
	if proto.IsMore == nil || *proto.IsMore != false {
		t.Errorf("IsMore: got %v, want false", proto.IsMore)
	}
	// ConId
	if proto.ConId == nil || *proto.ConId != 756733 {
		t.Errorf("ConId: got %v, want 756733", proto.ConId)
	}
	// Exchange
	if proto.Exchange == nil || *proto.Exchange != "SMART" {
		t.Errorf("Exchange: got %v, want SMART", proto.Exchange)
	}
	// Price
	if proto.Price == nil || *proto.Price != 0.01 {
		t.Errorf("Price: got %v, want 0.01", proto.Price)
	}
	// TriggerMethod
	if proto.TriggerMethod == nil || *proto.TriggerMethod != int32(LastTriggerMethod) {
		t.Errorf("TriggerMethod: got %v, want %d", proto.TriggerMethod, LastTriggerMethod)
	}
}

func TestCreatePriceConditionProtoDisjunction(t *testing.T) {
	cond := NewPriceCondition(12345, "NYSE", 100.50, DefaultTriggerMethod, true, false)

	proto := createPriceConditionProto(cond)

	if proto.IsConjunctionConnection == nil || *proto.IsConjunctionConnection != false {
		t.Errorf("IsConjunctionConnection: got %v, want false", proto.IsConjunctionConnection)
	}
	if proto.IsMore == nil || *proto.IsMore != true {
		t.Errorf("IsMore: got %v, want true", proto.IsMore)
	}
	if proto.ConId == nil || *proto.ConId != 12345 {
		t.Errorf("ConId: got %v, want 12345", proto.ConId)
	}
	if proto.Exchange == nil || *proto.Exchange != "NYSE" {
		t.Errorf("Exchange: got %v, want NYSE", proto.Exchange)
	}
	if proto.Price == nil || *proto.Price != 100.50 {
		t.Errorf("Price: got %v, want 100.50", proto.Price)
	}
	if proto.TriggerMethod == nil || *proto.TriggerMethod != int32(DefaultTriggerMethod) {
		t.Errorf("TriggerMethod: got %v, want %d", proto.TriggerMethod, DefaultTriggerMethod)
	}
}

func TestCreateVolumeConditionProto(t *testing.T) {
	cond := NewVolumeCondition(265598, "SMART", true, 10000, false)

	proto := createVolumeConditionProto(cond)

	// Type
	if proto.Type == nil || *proto.Type != int32(VolumeOrderCondition) {
		t.Errorf("Type: got %v, want %d", proto.Type, VolumeOrderCondition)
	}
	// IsConjunctionConnection
	if proto.IsConjunctionConnection == nil || *proto.IsConjunctionConnection != false {
		t.Errorf("IsConjunctionConnection: got %v, want false", proto.IsConjunctionConnection)
	}
	// IsMore
	if proto.IsMore == nil || *proto.IsMore != true {
		t.Errorf("IsMore: got %v, want true", proto.IsMore)
	}
	// ConId
	if proto.ConId == nil || *proto.ConId != 265598 {
		t.Errorf("ConId: got %v, want 265598", proto.ConId)
	}
	// Exchange
	if proto.Exchange == nil || *proto.Exchange != "SMART" {
		t.Errorf("Exchange: got %v, want SMART", proto.Exchange)
	}
	// Volume
	if proto.Volume == nil || *proto.Volume != 10000 {
		t.Errorf("Volume: got %v, want 10000", proto.Volume)
	}
}

func TestCreatePercentChangeConditionProto(t *testing.T) {
	cond := NewPercentageChangeCondition(5.5, 265598, "SMART", true, true)

	proto := createPercentChangeConditionProto(cond)

	// Type
	if proto.Type == nil || *proto.Type != int32(PercentChangeOrderCondition) {
		t.Errorf("Type: got %v, want %d", proto.Type, PercentChangeOrderCondition)
	}
	// IsConjunctionConnection
	if proto.IsConjunctionConnection == nil || *proto.IsConjunctionConnection != true {
		t.Errorf("IsConjunctionConnection: got %v, want true", proto.IsConjunctionConnection)
	}
	// IsMore
	if proto.IsMore == nil || *proto.IsMore != true {
		t.Errorf("IsMore: got %v, want true", proto.IsMore)
	}
	// ConId
	if proto.ConId == nil || *proto.ConId != 265598 {
		t.Errorf("ConId: got %v, want 265598", proto.ConId)
	}
	// Exchange
	if proto.Exchange == nil || *proto.Exchange != "SMART" {
		t.Errorf("Exchange: got %v, want SMART", proto.Exchange)
	}
	// ChangePercent
	if proto.ChangePercent == nil || *proto.ChangePercent != 5.5 {
		t.Errorf("ChangePercent: got %v, want 5.5", proto.ChangePercent)
	}
}

func TestCreateTimeConditionProto(t *testing.T) {
	cond := NewTimeCondition("20260101 10:00:00 US/Eastern", true, false)

	proto := createTimeConditionProto(cond)

	// Type
	if proto.Type == nil || *proto.Type != int32(TimeOrderCondition) {
		t.Errorf("Type: got %v, want %d", proto.Type, TimeOrderCondition)
	}
	// IsConjunctionConnection
	if proto.IsConjunctionConnection == nil || *proto.IsConjunctionConnection != false {
		t.Errorf("IsConjunctionConnection: got %v, want false", proto.IsConjunctionConnection)
	}
	// IsMore
	if proto.IsMore == nil || *proto.IsMore != true {
		t.Errorf("IsMore: got %v, want true", proto.IsMore)
	}
	// Time
	if proto.Time == nil || *proto.Time != "20260101 10:00:00 US/Eastern" {
		t.Errorf("Time: got %v, want '20260101 10:00:00 US/Eastern'", proto.Time)
	}
}

func TestCreateMarginConditionProto(t *testing.T) {
	cond := NewMarginCondition(30, false, true)

	proto := createMarginConditionProto(cond)

	// Type
	if proto.Type == nil || *proto.Type != int32(MarginOrderCondition) {
		t.Errorf("Type: got %v, want %d", proto.Type, MarginOrderCondition)
	}
	// IsConjunctionConnection
	if proto.IsConjunctionConnection == nil || *proto.IsConjunctionConnection != true {
		t.Errorf("IsConjunctionConnection: got %v, want true", proto.IsConjunctionConnection)
	}
	// IsMore
	if proto.IsMore == nil || *proto.IsMore != false {
		t.Errorf("IsMore: got %v, want false", proto.IsMore)
	}
	// Percent
	if proto.Percent == nil || *proto.Percent != 30 {
		t.Errorf("Percent: got %v, want 30", proto.Percent)
	}
}

func TestCreateExecutionConditionProto(t *testing.T) {
	cond := NewExecutionCondition("AAPL", "STK", "SMART", true)

	proto := createExecutionConditionProto(cond)

	// Type
	if proto.Type == nil || *proto.Type != int32(ExecutionOrderCondition) {
		t.Errorf("Type: got %v, want %d", proto.Type, ExecutionOrderCondition)
	}
	// IsConjunctionConnection
	if proto.IsConjunctionConnection == nil || *proto.IsConjunctionConnection != true {
		t.Errorf("IsConjunctionConnection: got %v, want true", proto.IsConjunctionConnection)
	}
	// SecType
	if proto.SecType == nil || *proto.SecType != "STK" {
		t.Errorf("SecType: got %v, want STK", proto.SecType)
	}
	// Exchange
	if proto.Exchange == nil || *proto.Exchange != "SMART" {
		t.Errorf("Exchange: got %v, want SMART", proto.Exchange)
	}
	// Symbol
	if proto.Symbol == nil || *proto.Symbol != "AAPL" {
		t.Errorf("Symbol: got %v, want AAPL", proto.Symbol)
	}
}

func TestCreateConditionsProto(t *testing.T) {
	order := NewOrder()
	order.Conditions = append(order.Conditions,
		NewPriceCondition(756733, "SMART", 100.0, LastTriggerMethod, true, true),
		NewVolumeCondition(265598, "SMART", false, 5000, false),
	)

	protos, err := createConditionsProto(order)
	if err != nil {
		t.Fatalf("createConditionsProto returned error: %v", err)
	}
	if len(protos) != 2 {
		t.Fatalf("expected 2 conditions, got %d", len(protos))
	}

	// First condition: PriceCondition
	p := protos[0]
	if p.Type == nil || *p.Type != int32(PriceOrderCondition) {
		t.Errorf("cond[0] Type: got %v, want %d", p.Type, PriceOrderCondition)
	}
	if p.ConId == nil || *p.ConId != 756733 {
		t.Errorf("cond[0] ConId: got %v, want 756733", p.ConId)
	}
	if p.Price == nil || *p.Price != 100.0 {
		t.Errorf("cond[0] Price: got %v, want 100.0", p.Price)
	}

	// Second condition: VolumeCondition
	v := protos[1]
	if v.Type == nil || *v.Type != int32(VolumeOrderCondition) {
		t.Errorf("cond[1] Type: got %v, want %d", v.Type, VolumeOrderCondition)
	}
	if v.ConId == nil || *v.ConId != 265598 {
		t.Errorf("cond[1] ConId: got %v, want 265598", v.ConId)
	}
	if v.Volume == nil || *v.Volume != 5000 {
		t.Errorf("cond[1] Volume: got %v, want 5000", v.Volume)
	}
}
