package metric

import "testing"

func TestGaugeInfoJsonString(t *testing.T) {
	g := NewGaugeInfo()
	g.Set("key1", "value1")
	jsonString := g.Snapshot().Value().String()
	expected := `{"key1":"value1"}`
	if jsonString != expected {
		t.Errorf("Expected JSON string: %s, got: %s", expected, jsonString)
	}
}
