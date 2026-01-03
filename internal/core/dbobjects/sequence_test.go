package dbobjects

import (
	"encoding/json"
	"testing"
)

func TestNewSequence(t *testing.T) {
	s := NewSequence("user_id_seq", 1, 1)

	if s.Name() != "user_id_seq" {
		t.Errorf("expected name 'user_id_seq', got %q", s.Name())
	}
	if s.StartValue() != 1 {
		t.Errorf("expected start value 1, got %d", s.StartValue())
	}
	if s.Increment() != 1 {
		t.Errorf("expected increment 1, got %d", s.Increment())
	}
	if s.MinValue() != 1 {
		t.Errorf("expected min value 1, got %d", s.MinValue())
	}
	if s.MaxValue() != 9223372036854775807 {
		t.Errorf("expected max value int64 max, got %d", s.MaxValue())
	}
	if s.Cache() != 1 {
		t.Errorf("expected cache 1, got %d", s.Cache())
	}
	if s.Cycle() {
		t.Error("expected cycle to be false by default")
	}
}

func TestSequenceSchema(t *testing.T) {
	s := NewSequence("test_seq", 1, 1)

	if s.Schema() != nil {
		t.Error("expected nil schema initially")
	}

	schema := NewSchema("public", "owner", nil)
	s.SetSchema(schema)

	if s.Schema() == nil {
		t.Fatal("expected schema to be set")
	}
	if s.Schema().Name() != "public" {
		t.Errorf("expected schema name 'public', got %q", s.Schema().Name())
	}
}

func TestSequenceStartValue(t *testing.T) {
	s := NewSequence("test_seq", 1, 1)

	s.SetStartValue(100)

	if s.StartValue() != 100 {
		t.Errorf("expected start value 100, got %d", s.StartValue())
	}
}

func TestSequenceIncrement(t *testing.T) {
	s := NewSequence("test_seq", 1, 1)

	s.SetIncrement(5)

	if s.Increment() != 5 {
		t.Errorf("expected increment 5, got %d", s.Increment())
	}
}

func TestSequenceMinValue(t *testing.T) {
	s := NewSequence("test_seq", 1, 1)

	s.SetMinValue(10)

	if s.MinValue() != 10 {
		t.Errorf("expected min value 10, got %d", s.MinValue())
	}
}

func TestSequenceMaxValue(t *testing.T) {
	s := NewSequence("test_seq", 1, 1)

	s.SetMaxValue(1000)

	if s.MaxValue() != 1000 {
		t.Errorf("expected max value 1000, got %d", s.MaxValue())
	}
}

func TestSequenceCache(t *testing.T) {
	s := NewSequence("test_seq", 1, 1)

	s.SetCache(20)

	if s.Cache() != 20 {
		t.Errorf("expected cache 20, got %d", s.Cache())
	}
}

func TestSequenceCycle(t *testing.T) {
	s := NewSequence("test_seq", 1, 1)

	s.SetCycle(true)

	if !s.Cycle() {
		t.Error("expected cycle to be true after setting")
	}

	s.SetCycle(false)

	if s.Cycle() {
		t.Error("expected cycle to be false after setting")
	}
}

func TestSequenceFullyQualifiedName(t *testing.T) {
	t.Run("without schema", func(t *testing.T) {
		s := NewSequence("my_seq", 1, 1)

		if s.FullyQualifiedName() != "my_seq" {
			t.Errorf("expected 'my_seq', got %q", s.FullyQualifiedName())
		}
	})

	t.Run("with schema", func(t *testing.T) {
		s := NewSequence("my_seq", 1, 1)
		schema := NewSchema("public", "owner", nil)
		s.SetSchema(schema)

		if s.FullyQualifiedName() != "public.my_seq" {
			t.Errorf("expected 'public.my_seq', got %q", s.FullyQualifiedName())
		}
	})
}

func TestSequenceMarshalJSON(t *testing.T) {
	s := NewSequence("order_id_seq", 1000, 10)
	s.SetMinValue(1)
	s.SetMaxValue(1000000)
	s.SetCache(50)
	s.SetCycle(true)

	data, err := json.Marshal(s)
	if err != nil {
		t.Fatalf("failed to marshal sequence: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("failed to unmarshal json: %v", err)
	}

	if result["name"] != "order_id_seq" {
		t.Errorf("expected name 'order_id_seq', got %v", result["name"])
	}
	if result["startValue"] != float64(1000) {
		t.Errorf("expected startValue 1000, got %v", result["startValue"])
	}
	if result["increment"] != float64(10) {
		t.Errorf("expected increment 10, got %v", result["increment"])
	}
	if result["minValue"] != float64(1) {
		t.Errorf("expected minValue 1, got %v", result["minValue"])
	}
	if result["maxValue"] != float64(1000000) {
		t.Errorf("expected maxValue 1000000, got %v", result["maxValue"])
	}
	if result["cache"] != float64(50) {
		t.Errorf("expected cache 50, got %v", result["cache"])
	}
	if result["cycle"] != true {
		t.Errorf("expected cycle true, got %v", result["cycle"])
	}
}
