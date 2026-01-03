package dbobjects

import (
	"encoding/json"
	"testing"
)

func TestNewPrimaryKey(t *testing.T) {
	table := NewTable("users", nil)
	col := NewColumn("id", "integer", false)
	columns := []*Column{col}

	pk := NewPrimaryKey("pk_users", table, columns)

	if pk.Name() != "pk_users" {
		t.Errorf("expected name 'pk_users', got %q", pk.Name())
	}
	if pk.Table() != table {
		t.Error("expected table to be set")
	}
	if len(pk.Columns()) != 1 {
		t.Errorf("expected 1 column, got %d", len(pk.Columns()))
	}
}

func TestNewPrimaryKeyWithNilColumns(t *testing.T) {
	pk := NewPrimaryKey("pk_test", nil, nil)

	if pk.Columns() == nil {
		t.Error("expected columns to be initialized to empty slice")
	}
	if len(pk.Columns()) != 0 {
		t.Errorf("expected 0 columns, got %d", len(pk.Columns()))
	}
}

func TestPrimaryKeyTable(t *testing.T) {
	pk := NewPrimaryKey("pk_test", nil, nil)

	if pk.Table() != nil {
		t.Error("expected nil table initially")
	}

	table := NewTable("users", nil)
	pk.SetTable(table)

	if pk.Table() == nil {
		t.Fatal("expected table to be set")
	}
	if pk.Table().Name() != "users" {
		t.Errorf("expected table name 'users', got %q", pk.Table().Name())
	}
}

func TestPrimaryKeyColumns(t *testing.T) {
	pk := NewPrimaryKey("pk_composite", nil, nil)
	col1 := NewColumn("user_id", "integer", false)
	col2 := NewColumn("order_id", "integer", false)

	pk.AddColumn(col1)
	pk.AddColumn(col2)

	columns := pk.Columns()
	if len(columns) != 2 {
		t.Fatalf("expected 2 columns, got %d", len(columns))
	}
	if columns[0].Name() != "user_id" {
		t.Errorf("expected first column 'user_id', got %q", columns[0].Name())
	}
	if columns[1].Name() != "order_id" {
		t.Errorf("expected second column 'order_id', got %q", columns[1].Name())
	}
}

func TestPrimaryKeyMarshalJSON(t *testing.T) {
	col1 := NewColumn("emp_id", "integer", false)
	col2 := NewColumn("dept_id", "integer", false)

	pk := NewPrimaryKey("pk_emp_dept", nil, []*Column{col1, col2})

	data, err := json.Marshal(pk)
	if err != nil {
		t.Fatalf("failed to marshal primary key: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("failed to unmarshal json: %v", err)
	}

	if result["name"] != "pk_emp_dept" {
		t.Errorf("expected name 'pk_emp_dept', got %v", result["name"])
	}

	columns := result["columns"].([]interface{})
	if len(columns) != 2 {
		t.Fatalf("expected 2 columns, got %d", len(columns))
	}
	if columns[0] != "emp_id" {
		t.Errorf("expected first column 'emp_id', got %v", columns[0])
	}
	if columns[1] != "dept_id" {
		t.Errorf("expected second column 'dept_id', got %v", columns[1])
	}
}

func TestPrimaryKeyMarshalJSONEmptyColumns(t *testing.T) {
	pk := NewPrimaryKey("pk_test", nil, nil)

	data, err := json.Marshal(pk)
	if err != nil {
		t.Fatalf("failed to marshal primary key: %v", err)
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
