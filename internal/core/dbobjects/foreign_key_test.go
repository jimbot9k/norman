package dbobjects

import (
	"encoding/json"
	"testing"
)

func TestNewForeignKey(t *testing.T) {
	fk := NewForeignKey("fk_user_dept", "departments")

	if fk.Name() != "fk_user_dept" {
		t.Errorf("expected name 'fk_user_dept', got %q", fk.Name())
	}
	if fk.ReferencedTable() != "departments" {
		t.Errorf("expected referenced table 'departments', got %q", fk.ReferencedTable())
	}
	if len(fk.Columns()) != 0 {
		t.Errorf("expected empty columns, got %d", len(fk.Columns()))
	}
	if len(fk.ReferencedColumns()) != 0 {
		t.Errorf("expected empty referenced columns, got %d", len(fk.ReferencedColumns()))
	}
	if fk.OnDelete() != ActionNoAction {
		t.Errorf("expected default OnDelete %q, got %q", ActionNoAction, fk.OnDelete())
	}
	if fk.OnUpdate() != ActionNoAction {
		t.Errorf("expected default OnUpdate %q, got %q", ActionNoAction, fk.OnUpdate())
	}
}

func TestForeignKeyTable(t *testing.T) {
	fk := NewForeignKey("fk_test", "ref_table")

	if fk.Table() != nil {
		t.Error("expected nil table initially")
	}

	table := NewTable("users", nil)
	fk.SetTable(table)

	if fk.Table() == nil {
		t.Fatal("expected table to be set")
	}
	if fk.Table().Name() != "users" {
		t.Errorf("expected table name 'users', got %q", fk.Table().Name())
	}
}

func TestForeignKeyColumns(t *testing.T) {
	fk := NewForeignKey("fk_test", "ref_table")
	col1 := NewColumn("dept_id", "integer", false)
	col2 := NewColumn("org_id", "integer", false)

	fk.AddColumn(col1)
	fk.AddColumn(col2)

	columns := fk.Columns()
	if len(columns) != 2 {
		t.Fatalf("expected 2 columns, got %d", len(columns))
	}
	if columns[0].Name() != "dept_id" {
		t.Errorf("expected first column 'dept_id', got %q", columns[0].Name())
	}
	if columns[1].Name() != "org_id" {
		t.Errorf("expected second column 'org_id', got %q", columns[1].Name())
	}
}

func TestForeignKeyReferencedSchema(t *testing.T) {
	fk := NewForeignKey("fk_test", "ref_table")

	if fk.ReferencedSchema() != "" {
		t.Errorf("expected empty referenced schema, got %q", fk.ReferencedSchema())
	}

	fk.SetReferencedSchema("public")

	if fk.ReferencedSchema() != "public" {
		t.Errorf("expected 'public', got %q", fk.ReferencedSchema())
	}
}

func TestForeignKeyReferencedTable(t *testing.T) {
	fk := NewForeignKey("fk_test", "original_table")

	fk.SetReferencedTable("new_table")

	if fk.ReferencedTable() != "new_table" {
		t.Errorf("expected 'new_table', got %q", fk.ReferencedTable())
	}
}

func TestForeignKeyReferencedColumns(t *testing.T) {
	fk := NewForeignKey("fk_test", "ref_table")
	refCol1 := NewColumn("id", "integer", false)
	refCol2 := NewColumn("version", "integer", false)

	fk.AddReferencedColumn(refCol1)
	fk.AddReferencedColumn(refCol2)

	refCols := fk.ReferencedColumns()
	if len(refCols) != 2 {
		t.Fatalf("expected 2 referenced columns, got %d", len(refCols))
	}
	if refCols[0].Name() != "id" {
		t.Errorf("expected first ref column 'id', got %q", refCols[0].Name())
	}
	if refCols[1].Name() != "version" {
		t.Errorf("expected second ref column 'version', got %q", refCols[1].Name())
	}
}

func TestForeignKeyReferentialActions(t *testing.T) {
	tests := []struct {
		name     string
		action   ReferentialAction
		expected ReferentialAction
	}{
		{"no action", ActionNoAction, ActionNoAction},
		{"restrict", ActionRestrict, ActionRestrict},
		{"cascade", ActionCascade, ActionCascade},
		{"set null", ActionSetNull, ActionSetNull},
		{"set default", ActionSetDefault, ActionSetDefault},
	}

	for _, tt := range tests {
		t.Run(tt.name+" on delete", func(t *testing.T) {
			fk := NewForeignKey("fk_test", "ref_table")
			fk.SetOnDelete(tt.action)

			if fk.OnDelete() != tt.expected {
				t.Errorf("expected OnDelete %q, got %q", tt.expected, fk.OnDelete())
			}
		})

		t.Run(tt.name+" on update", func(t *testing.T) {
			fk := NewForeignKey("fk_test", "ref_table")
			fk.SetOnUpdate(tt.action)

			if fk.OnUpdate() != tt.expected {
				t.Errorf("expected OnUpdate %q, got %q", tt.expected, fk.OnUpdate())
			}
		})
	}
}

func TestForeignKeyMarshalJSON(t *testing.T) {
	fk := NewForeignKey("fk_emp_dept", "departments")
	fk.SetReferencedSchema("public")
	fk.SetOnDelete(ActionCascade)
	fk.SetOnUpdate(ActionRestrict)

	col := NewColumn("dept_id", "integer", false)
	fk.AddColumn(col)

	refCol := NewColumn("id", "integer", false)
	fk.AddReferencedColumn(refCol)

	data, err := json.Marshal(fk)
	if err != nil {
		t.Fatalf("failed to marshal foreign key: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("failed to unmarshal json: %v", err)
	}

	if result["name"] != "fk_emp_dept" {
		t.Errorf("expected name 'fk_emp_dept', got %v", result["name"])
	}
	if result["referencedSchema"] != "public" {
		t.Errorf("expected referencedSchema 'public', got %v", result["referencedSchema"])
	}
	if result["referencedTable"] != "departments" {
		t.Errorf("expected referencedTable 'departments', got %v", result["referencedTable"])
	}
	if result["onDelete"] != string(ActionCascade) {
		t.Errorf("expected onDelete %q, got %v", ActionCascade, result["onDelete"])
	}
	if result["onUpdate"] != string(ActionRestrict) {
		t.Errorf("expected onUpdate %q, got %v", ActionRestrict, result["onUpdate"])
	}

	columns := result["columns"].([]interface{})
	if len(columns) != 1 || columns[0] != "dept_id" {
		t.Errorf("expected columns ['dept_id'], got %v", columns)
	}

	refColumns := result["referencedColumns"].([]interface{})
	if len(refColumns) != 1 || refColumns[0] != "id" {
		t.Errorf("expected referencedColumns ['id'], got %v", refColumns)
	}
}

func TestForeignKeyMarshalJSONEmptyColumns(t *testing.T) {
	fk := NewForeignKey("fk_test", "ref_table")

	data, err := json.Marshal(fk)
	if err != nil {
		t.Fatalf("failed to marshal foreign key: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("failed to unmarshal json: %v", err)
	}

	columns := result["columns"].([]interface{})
	if len(columns) != 0 {
		t.Errorf("expected empty columns, got %v", columns)
	}

	refColumns := result["referencedColumns"].([]interface{})
	if len(refColumns) != 0 {
		t.Errorf("expected empty referencedColumns, got %v", refColumns)
	}
}
