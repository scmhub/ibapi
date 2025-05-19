package ibapi

import (
	"encoding/json"
	"testing"
)

func TestContractComplexJSON(t *testing.T) {
	// Create a complex contract
	contract := NewContract()
	contract.Symbol = "SPREAD"
	contract.SecType = "BAG"
	contract.Currency = "USD"
	contract.Exchange = "SMART"

	// Add combo legs
	leg1 := NewComboLeg()
	leg1.ConID = 12345
	leg1.Ratio = 1
	leg1.Action = "BUY"
	leg1.Exchange = "SMART"
	leg1.OpenClose = int64(OPEN_POS)

	leg2 := NewComboLeg()
	leg2.ConID = 67890
	leg2.Ratio = 1
	leg2.Action = "SELL"
	leg2.Exchange = "SMART"
	leg2.OpenClose = int64(CLOSE_POS)

	contract.ComboLegs = []ComboLeg{leg1, leg2}

	// Add delta neutral contract
	contract.DeltaNeutralContract = &DeltaNeutralContract{
		ConID: 11111,
		Delta: 0.5,
		Price: 150.0,
	}

	// Marshal to JSON
	data, err := json.Marshal(contract)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	// Verify JSON structure
	var jsonMap map[string]any
	if err := json.Unmarshal(data, &jsonMap); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	// Check combo legs
	legs, ok := jsonMap["comboLegs"].([]any)
	if !ok || len(legs) != 2 {
		t.Error("Expected 2 combo legs in JSON")
	}

	// Check delta neutral
	delta, ok := jsonMap["deltaNeutralContract"].(map[string]any)
	if !ok {
		t.Error("Expected deltaNeutralContract in JSON")
	}
	if delta["conId"].(float64) != 11111 {
		t.Errorf("Expected deltaNeutralContract.conId = 11111, got %v", delta["conId"])
	}

	// Unmarshal back to contract
	var decoded Contract
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}
	// Equal check
	if !decoded.Equal(contract) {
		t.Errorf("Decoded contract not equal to original:\nOriginal: %+v\nDecoded: %+v",
			contract, decoded)
	}
	// Verify combo legs
	if len(decoded.ComboLegs) != 2 {
		t.Fatalf("Expected 2 combo legs, got %d", len(decoded.ComboLegs))
	}
	if decoded.ComboLegs[0].ConID != 12345 {
		t.Errorf("Expected first leg ConID 12345, got %d", decoded.ComboLegs[0].ConID)
	}
	if decoded.ComboLegs[1].ConID != 67890 {
		t.Errorf("Expected second leg ConID 67890, got %d", decoded.ComboLegs[1].ConID)
	}

	// Verify delta neutral contract
	if decoded.DeltaNeutralContract == nil {
		t.Fatal("Expected non-nil DeltaNeutralContract")
	}
	if decoded.DeltaNeutralContract.ConID != 11111 {
		t.Errorf("Expected DeltaNeutralContract.ConID 11111, got %d",
			decoded.DeltaNeutralContract.ConID)
	}
	if decoded.DeltaNeutralContract.Delta != 0.5 {
		t.Errorf("Expected DeltaNeutralContract.Delta 0.5, got %f",
			decoded.DeltaNeutralContract.Delta)
	}
}
