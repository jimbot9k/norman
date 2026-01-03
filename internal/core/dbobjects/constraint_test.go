package dbobjects

import (
	"encoding/json"
	"testing"
)

func TestNewConstraint(t *testing.T) {
	tests := []struct {
		name           string
		constraintName string
		constraintType ConstraintType
	}{
		{"check constraint", "chk_positive", ConstraintTypeCheck},
		{"unique constraint", "uq_email", ConstraintTypeUnique},
		{"not null constraint", "nn_name", ConstraintTypeNotNull},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewConstraint(tt.constraintName, tt.constraintType)

			if c.Name() != tt.constraintName {
				t.Errorf("expected name %q, got %q", tt.constraintName, c.Name())
			}
			if c.Type() != tt.constraintType {
				t.Errorf("expected type %q, got %q", tt.constraintType, c.Type())
			}
			if len(c.Columns()) != 0 {
				t.Errorf("expected empty columns, got %d", len(c.Columns()))
			}
		})
	}
}

func TestConstraintSetType(t *testing.T) {
	c := NewConstraint("test", ConstraintTypeCheck)
	c.SetType(ConstraintTypeUnique)

	if c.Type() != ConstraintTypeUnique {
		t.Errorf("expected type %q, got %q", ConstraintTypeUnique, c.Type())
	}
}

func TestConstraintTable(t *testing.T) {
	c := NewConstraint("test", ConstraintTypeCheck)

	if c.Table() != nil {
		t.Error("expected nil table initially")
	}

	table := NewTable("users", nil)
	c.SetTable(table)

	if c.Table() == nil {
		t.Fatal("expected table to be set")
	}
	if c.Table().Name() != "users" {
		t.Errorf("expected table name 'users', got %q", c.Table().Name())
	}
}

func TestConstraintColumns(t *testing.T) {
	c := NewConstraint("uq_composite", ConstraintTypeUnique)
	col1 := NewColumn("first_name", "varchar", false)
	col2 := NewColumn("last_name", "varchar", false)

	c.AddColumn(col1)
	c.AddColumn(col2)

	columns := c.Columns()
	if len(columns) != 2 {
		t.Fatalf("expected 2 columns, got %d", len(columns))
	}
	if columns[0].Name() != "first_name" {
		t.Errorf("expected first column 'first_name', got %q", columns[0].Name())
	}
	if columns[1].Name() != "last_name" {
		t.Errorf("expected second column 'last_name', got %q", columns[1].Name())
	}
}

func TestConstraintCheckExpression(t *testing.T) {
	c := NewConstraint("chk_salary", ConstraintTypeCheck)

	if c.CheckExpression() != "" {
		t.Errorf("expected empty check expression, got %q", c.CheckExpression())
	}

	c.SetCheckExpression("salary >= 0")

	if c.CheckExpression() != "salary >= 0" {
		t.Errorf("expected 'salary >= 0', got %q", c.CheckExpression())
	}
}

func TestConstraintMarshalJSON(t *testing.T) {
	c := NewConstraint("chk_positive", ConstraintTypeCheck)
	col := NewColumn("amount", "decimal", false)
	c.AddColumn(col)
	c.SetCheckExpression("amount > 0")

	data, err := json.Marshal(c)
	if err != nil {
		t.Fatalf("failed to marshal constraint: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("failed to unmarshal json: %v", err)
	}

	if result["name"] != "chk_positive" {
		t.Errorf("expected name 'chk_positive', got %v", result["name"])
	}
	if result["type"] != string(ConstraintTypeCheck) {
		t.Errorf("expected type %q, got %v", ConstraintTypeCheck, result["type"])
	}
	if result["checkExpression"] != "amount > 0" {
		t.Errorf("expected checkExpression 'amount > 0', got %v", result["checkExpression"])
	}

	columns := result["columns"].([]interface{})
	if len(columns) != 1 || columns[0] != "amount" {
		t.Errorf("expected columns ['amount'], got %v", columns)
	}
}

func TestConstraintMarshalJSONOmitsEmptyCheckExpression(t *testing.T) {
	c := NewConstraint("uq_email", ConstraintTypeUnique)
	col := NewColumn("email", "varchar", false)
	c.AddColumn(col)

	data, err := json.Marshal(c)
	if err != nil {
		t.Fatalf("failed to marshal constraint: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("failed to unmarshal json: %v", err)
	}

	if _, exists := result["checkExpression"]; exists {
		t.Error("expected checkExpression to be omitted when empty")
	}
}

func TestConstraintMarshalJSONEmptyColumns(t *testing.T) {
	c := NewConstraint("test", ConstraintTypeNotNull)

	data, err := json.Marshal(c)
	if err != nil {
		t.Fatalf("failed to marshal constraint: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("failed to unmarshal json: %v", err)
	}

	columns := result["columns"].([]interface{})
	if len(columns) != 0 {
		t.Errorf("expected empty columns, got %v", columns)
	}
}
