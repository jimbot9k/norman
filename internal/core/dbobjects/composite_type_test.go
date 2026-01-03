package dbobjects

import "testing"

func TestNewCompositeType(t *testing.T) {
	ct := NewCompositeType("address_type", "CREATE TYPE address_type AS (street text, city text)")

	if ct.Name() != "address_type" {
		t.Errorf("expected name 'address_type', got %q", ct.Name())
	}
	if ct.Definition() != "CREATE TYPE address_type AS (street text, city text)" {
		t.Errorf("expected definition, got %q", ct.Definition())
	}
}

func TestCompositeTypeName(t *testing.T) {
	ct := NewCompositeType("test_type", "definition")

	if ct.Name() != "test_type" {
		t.Errorf("expected 'test_type', got %q", ct.Name())
	}
}

func TestCompositeTypeDefinition(t *testing.T) {
	ct := NewCompositeType("test_type", "CREATE TYPE test AS (a int, b text)")

	if ct.Definition() != "CREATE TYPE test AS (a int, b text)" {
		t.Errorf("expected definition, got %q", ct.Definition())
	}
}
