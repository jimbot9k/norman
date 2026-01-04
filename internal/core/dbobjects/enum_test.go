package dbobjects

import "testing"

func TestNewEnum(t *testing.T) {
	values := []string{"pending", "approved", "rejected"}
	e := NewEnum("status_enum", values)

	if e.Name() != "status_enum" {
		t.Errorf("expected name 'status_enum', got %q", e.Name())
	}
	if len(e.Values()) != 3 {
		t.Errorf("expected 3 values, got %d", len(e.Values()))
	}
}

func TestEnumName(t *testing.T) {
	e := NewEnum("test_enum", []string{"a", "b"})

	if e.Name() != "test_enum" {
		t.Errorf("expected 'test_enum', got %q", e.Name())
	}
}

func TestEnumValues(t *testing.T) {
	values := []string{"low", "medium", "high"}
	e := NewEnum("priority_enum", values)

	enumValues := e.Values()
	if len(enumValues) != 3 {
		t.Fatalf("expected 3 values, got %d", len(enumValues))
	}
	if enumValues[0] != "low" {
		t.Errorf("expected first value 'low', got %q", enumValues[0])
	}
	if enumValues[1] != "medium" {
		t.Errorf("expected second value 'medium', got %q", enumValues[1])
	}
	if enumValues[2] != "high" {
		t.Errorf("expected third value 'high', got %q", enumValues[2])
	}
}

func TestEnumWithNilValues(t *testing.T) {
	e := NewEnum("empty_enum", nil)

	if e.Values() != nil {
		t.Error("expected nil values")
	}
}

func TestEnumWithEmptyValues(t *testing.T) {
	e := NewEnum("empty_enum", []string{})

	if len(e.Values()) != 0 {
		t.Errorf("expected empty values, got %d", len(e.Values()))
	}
}
