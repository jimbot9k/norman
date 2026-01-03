package dbobjects

import (
	"encoding/json"
	"testing"
)

func TestNewTrigger(t *testing.T) {
	tr := NewTrigger("trg_audit", "EXECUTE FUNCTION audit_log()")

	if tr.Name() != "trg_audit" {
		t.Errorf("expected name 'trg_audit', got %q", tr.Name())
	}
	if tr.Definition() != "EXECUTE FUNCTION audit_log()" {
		t.Errorf("expected definition, got %q", tr.Definition())
	}
	if len(tr.Events()) != 0 {
		t.Errorf("expected empty events, got %d", len(tr.Events()))
	}
	if tr.ForEach() != "ROW" {
		t.Errorf("expected default forEach 'ROW', got %q", tr.ForEach())
	}
}

func TestTriggerTable(t *testing.T) {
	tr := NewTrigger("trg_test", "test")

	if tr.Table() != nil {
		t.Error("expected nil table initially")
	}

	table := NewTable("users", nil)
	tr.SetTable(table)

	if tr.Table() == nil {
		t.Fatal("expected table to be set")
	}
	if tr.Table().Name() != "users" {
		t.Errorf("expected table name 'users', got %q", tr.Table().Name())
	}
}

func TestTriggerDefinition(t *testing.T) {
	tr := NewTrigger("trg_test", "original")

	tr.SetDefinition("updated definition")

	if tr.Definition() != "updated definition" {
		t.Errorf("expected 'updated definition', got %q", tr.Definition())
	}
}

func TestTriggerTiming(t *testing.T) {
	tests := []struct {
		name   string
		timing TriggerTiming
	}{
		{"before", TriggerTimingBefore},
		{"after", TriggerTimingAfter},
		{"instead of", TriggerTimingInsteadOf},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := NewTrigger("trg_test", "test")
			tr.SetTiming(tt.timing)

			if tr.Timing() != tt.timing {
				t.Errorf("expected timing %q, got %q", tt.timing, tr.Timing())
			}
		})
	}
}

func TestTriggerEvents(t *testing.T) {
	tr := NewTrigger("trg_all", "test")

	tr.AddEvent(TriggerEventInsert)
	tr.AddEvent(TriggerEventUpdate)
	tr.AddEvent(TriggerEventDelete)

	events := tr.Events()
	if len(events) != 3 {
		t.Fatalf("expected 3 events, got %d", len(events))
	}
	if events[0] != TriggerEventInsert {
		t.Errorf("expected first event INSERT, got %q", events[0])
	}
	if events[1] != TriggerEventUpdate {
		t.Errorf("expected second event UPDATE, got %q", events[1])
	}
	if events[2] != TriggerEventDelete {
		t.Errorf("expected third event DELETE, got %q", events[2])
	}
}

func TestTriggerFunction(t *testing.T) {
	tr := NewTrigger("trg_test", "test")

	if tr.Function() != nil {
		t.Error("expected nil function initially")
	}

	fn := NewFunction("audit_fn", "INSERT INTO audit_log VALUES (...)")
	tr.SetFunction(fn)

	if tr.Function() == nil {
		t.Fatal("expected function to be set")
	}
	if tr.Function().Name() != "audit_fn" {
		t.Errorf("expected function name 'audit_fn', got %q", tr.Function().Name())
	}
}

func TestTriggerForEach(t *testing.T) {
	tr := NewTrigger("trg_test", "test")

	tr.SetForEach("STATEMENT")

	if tr.ForEach() != "STATEMENT" {
		t.Errorf("expected forEach 'STATEMENT', got %q", tr.ForEach())
	}
}

func TestTriggerMarshalJSONWithFunction(t *testing.T) {
	tr := NewTrigger("trg_audit", "EXECUTE FUNCTION audit_log()")
	tr.SetTiming(TriggerTimingAfter)
	tr.AddEvent(TriggerEventInsert)
	tr.AddEvent(TriggerEventUpdate)
	tr.SetForEach("ROW")

	fn := NewFunction("audit_log", "INSERT INTO audit ...")
	tr.SetFunction(fn)

	data, err := json.Marshal(tr)
	if err != nil {
		t.Fatalf("failed to marshal trigger: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("failed to unmarshal json: %v", err)
	}

	if result["name"] != "trg_audit" {
		t.Errorf("expected name 'trg_audit', got %v", result["name"])
	}
	if result["timing"] != string(TriggerTimingAfter) {
		t.Errorf("expected timing %q, got %v", TriggerTimingAfter, result["timing"])
	}
	if result["forEach"] != "ROW" {
		t.Errorf("expected forEach 'ROW', got %v", result["forEach"])
	}
	if result["function"] != "audit_log" {
		t.Errorf("expected function 'audit_log', got %v", result["function"])
	}

	events := result["events"].([]interface{})
	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}
	if events[0] != string(TriggerEventInsert) {
		t.Errorf("expected first event INSERT, got %v", events[0])
	}
}

func TestTriggerMarshalJSONWithoutFunction(t *testing.T) {
	tr := NewTrigger("trg_simple", "test definition")
	tr.SetTiming(TriggerTimingBefore)
	tr.AddEvent(TriggerEventDelete)

	data, err := json.Marshal(tr)
	if err != nil {
		t.Fatalf("failed to marshal trigger: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("failed to unmarshal json: %v", err)
	}

	// Function should be omitted when nil (empty string marshals but omitempty applies)
	if fn, exists := result["function"]; exists && fn != "" {
		t.Errorf("expected function to be omitted or empty, got %v", fn)
	}
}

func TestTriggerMarshalJSONEmptyEvents(t *testing.T) {
	tr := NewTrigger("trg_test", "test")

	data, err := json.Marshal(tr)
	if err != nil {
		t.Fatalf("failed to marshal trigger: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("failed to unmarshal json: %v", err)
	}

	events := result["events"].([]interface{})
	if len(events) != 0 {
		t.Errorf("expected empty events, got %v", events)
	}
}
