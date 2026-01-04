package dbobjects

import (
	"encoding/json"
	"testing"
)

func TestNewIndex(t *testing.T) {
	table := NewTable("users", nil)
	col := NewColumn("email", "varchar", false)
	columns := []*Column{col}

	idx := NewIndex("idx_users_email", table, columns, true)

	if idx.Name() != "idx_users_email" {
		t.Errorf("expected name 'idx_users_email', got %q", idx.Name())
	}
	if idx.Table() != table {
		t.Error("expected table to be set")
	}
	if len(idx.Columns()) != 1 {
		t.Errorf("expected 1 column, got %d", len(idx.Columns()))
	}
	if !idx.IsUnique() {
		t.Error("expected isUnique to be true")
	}
	if idx.IsPrimary() {
		t.Error("expected isPrimary to be false by default")
	}
	if idx.IndexType() != IndexTypeBTree {
		t.Errorf("expected default index type %q, got %q", IndexTypeBTree, idx.IndexType())
	}
}

func TestNewIndexWithNilColumns(t *testing.T) {
	table := NewTable("users", nil)
	idx := NewIndex("idx_test", table, nil, false)

	if idx.Columns() != nil {
		t.Error("expected nil columns")
	}
}

func TestIndexTable(t *testing.T) {
	idx := NewIndex("idx_test", nil, nil, false)

	if idx.Table() != nil {
		t.Error("expected nil table initially")
	}

	table := NewTable("users", nil)
	idx.SetTable(table)

	if idx.Table() == nil {
		t.Fatal("expected table to be set")
	}
	if idx.Table().Name() != "users" {
		t.Errorf("expected table name 'users', got %q", idx.Table().Name())
	}
}

func TestIndexColumns(t *testing.T) {
	idx := NewIndex("idx_test", nil, []*Column{}, false)
	col1 := NewColumn("first_name", "varchar", false)
	col2 := NewColumn("last_name", "varchar", false)

	idx.AddColumn(col1)
	idx.AddColumn(col2)

	columns := idx.Columns()
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

func TestIndexUnique(t *testing.T) {
	idx := NewIndex("idx_test", nil, nil, false)

	if idx.IsUnique() {
		t.Error("expected IsUnique to be false")
	}

	idx.SetUnique(true)

	if !idx.IsUnique() {
		t.Error("expected IsUnique to be true after setting")
	}
}

func TestIndexPrimary(t *testing.T) {
	idx := NewIndex("idx_test", nil, nil, false)

	if idx.IsPrimary() {
		t.Error("expected IsPrimary to be false")
	}

	idx.SetPrimary(true)

	if !idx.IsPrimary() {
		t.Error("expected IsPrimary to be true after setting")
	}
}

func TestIndexType(t *testing.T) {
	tests := []struct {
		name      string
		indexType IndexType
	}{
		{"btree", IndexTypeBTree},
		{"hash", IndexTypeHash},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idx := NewIndex("idx_test", nil, nil, false)
			idx.SetIndexType(tt.indexType)

			if idx.IndexType() != tt.indexType {
				t.Errorf("expected index type %q, got %q", tt.indexType, idx.IndexType())
			}
		})
	}
}

func TestIndexMarshalJSON(t *testing.T) {
	table := NewTable("users", nil)
	col1 := NewColumn("email", "varchar", false)
	col2 := NewColumn("name", "varchar", false)

	idx := NewIndex("idx_users_email_name", table, []*Column{col1, col2}, true)
	idx.SetPrimary(false)
	idx.SetIndexType(IndexTypeBTree)

	data, err := json.Marshal(idx)
	if err != nil {
		t.Fatalf("failed to marshal index: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("failed to unmarshal json: %v", err)
	}

	if result["name"] != "idx_users_email_name" {
		t.Errorf("expected name 'idx_users_email_name', got %v", result["name"])
	}
	if result["isUnique"] != true {
		t.Errorf("expected isUnique true, got %v", result["isUnique"])
	}
	if result["isPrimary"] != false {
		t.Errorf("expected isPrimary false, got %v", result["isPrimary"])
	}
	if result["indexType"] != string(IndexTypeBTree) {
		t.Errorf("expected indexType %q, got %v", IndexTypeBTree, result["indexType"])
	}

	columns := result["columns"].([]interface{})
	if len(columns) != 2 {
		t.Fatalf("expected 2 columns, got %d", len(columns))
	}
	if columns[0] != "email" {
		t.Errorf("expected first column 'email', got %v", columns[0])
	}
	if columns[1] != "name" {
		t.Errorf("expected second column 'name', got %v", columns[1])
	}
}

func TestIndexMarshalJSONEmptyColumns(t *testing.T) {
	idx := NewIndex("idx_test", nil, []*Column{}, false)

	data, err := json.Marshal(idx)
	if err != nil {
		t.Fatalf("failed to marshal index: %v", err)
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
