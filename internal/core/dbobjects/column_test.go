package dbobjects

import (
	"encoding/json"
	"testing"
)

func TestNewColumn(t *testing.T) {
	col := NewColumn("id", "integer", false)

	if col.Name() != "id" {
		t.Errorf("expected name 'id', got '%s'", col.Name())
	}
	if col.DataType() != "integer" {
		t.Errorf("expected dataType 'integer', got '%s'", col.DataType())
	}
	if col.IsNullable() != false {
		t.Errorf("expected nullable false, got true")
	}
}

func TestColumnSetDefaultValue(t *testing.T) {
	col := NewColumn("status", "varchar", true)
	col.SetDefaultValue("active")

	if col.DefaultValue() == nil {
		t.Fatal("expected default value to be set")
	}
	if *col.DefaultValue() != "active" {
		t.Errorf("expected default value 'active', got '%s'", *col.DefaultValue())
	}
}

func TestColumnSetOrdinalPosition(t *testing.T) {
	col := NewColumn("name", "text", true)
	col.SetOrdinalPosition(3)

	if col.OrdinalPosition() != 3 {
		t.Errorf("expected ordinal position 3, got %d", col.OrdinalPosition())
	}
}

func TestColumnSetCharMaxLength(t *testing.T) {
	col := NewColumn("description", "varchar", true)
	col.SetCharMaxLength(255)

	if col.CharMaxLength() == nil {
		t.Fatal("expected char max length to be set")
	}
	if *col.CharMaxLength() != 255 {
		t.Errorf("expected char max length 255, got %d", *col.CharMaxLength())
	}
}

func TestColumnSetNumericPrecision(t *testing.T) {
	col := NewColumn("price", "decimal", false)
	col.SetNumericPrecision(10)

	if col.NumericPrecision() == nil {
		t.Fatal("expected numeric precision to be set")
	}
	if *col.NumericPrecision() != 10 {
		t.Errorf("expected numeric precision 10, got %d", *col.NumericPrecision())
	}
}

func TestColumnSetNumericScale(t *testing.T) {
	col := NewColumn("amount", "decimal", false)
	col.SetNumericScale(2)

	if col.NumericScale() == nil {
		t.Fatal("expected numeric scale to be set")
	}
	if *col.NumericScale() != 2 {
		t.Errorf("expected numeric scale 2, got %d", *col.NumericScale())
	}
}

func TestColumnSetTable(t *testing.T) {
	col := NewColumn("user_id", "integer", false)
	table := NewTable("users", make(map[string]*Column))
	col.SetTable(table)
	col.Table().AddColumn(col)

	if col.Table() == nil {
		t.Fatal("expected table to be set")
	}
	if col.Table().Name() != "users" {
		t.Errorf("expected table name 'users', got '%s'", col.Table().Name())
	}
}

func TestColumnMarshalJSON(t *testing.T) {
	col := NewColumn("email", "varchar", false)
	col.SetOrdinalPosition(2)
	col.SetCharMaxLength(100)
	defaultVal := "test@example.com"
	col.SetDefaultValue(defaultVal)

	data, err := json.Marshal(col)
	if err != nil {
		t.Fatalf("failed to marshal column: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("failed to unmarshal json: %v", err)
	}

	if result["name"] != "email" {
		t.Errorf("expected name 'email', got '%v'", result["name"])
	}
	if result["dataType"] != "varchar" {
		t.Errorf("expected dataType 'varchar', got '%v'", result["dataType"])
	}
	if result["nullable"] != false {
		t.Errorf("expected nullable false, got %v", result["nullable"])
	}
	if result["ordinalPosition"] != float64(2) {
		t.Errorf("expected ordinalPosition 2, got %v", result["ordinalPosition"])
	}
}

func TestColumnNullableFieldsOmitEmpty(t *testing.T) {
	col := NewColumn("id", "integer", false)
	col.SetOrdinalPosition(1)

	data, err := json.Marshal(col)
	if err != nil {
		t.Fatalf("failed to marshal column: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("failed to unmarshal json: %v", err)
	}

	if _, exists := result["defaultValue"]; exists {
		t.Error("expected defaultValue to be omitted")
	}
	if _, exists := result["charMaxLength"]; exists {
		t.Error("expected charMaxLength to be omitted")
	}
	if _, exists := result["numericPrecision"]; exists {
		t.Error("expected numericPrecision to be omitted")
	}
	if _, exists := result["numericScale"]; exists {
		t.Error("expected numericScale to be omitted")
	}
}