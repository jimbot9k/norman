package dbobjects

import (
	"encoding/json"
	"testing"
)

func TestNewTable(t *testing.T) {
	tbl := NewTable("users", nil)

	if tbl.Name() != "users" {
		t.Errorf("expected name 'users', got %q", tbl.Name())
	}
	if tbl.Columns() == nil {
		t.Error("expected columns map to be initialized")
	}
	if tbl.ForeignKeys() == nil {
		t.Error("expected foreign keys to be initialized")
	}
	if tbl.Indexes() == nil {
		t.Error("expected indexes to be initialized")
	}
	if tbl.Constraints() == nil {
		t.Error("expected constraints to be initialized")
	}
	if tbl.Triggers() == nil {
		t.Error("expected triggers to be initialized")
	}
}

func TestNewTableWithColumns(t *testing.T) {
	col := NewColumn("id", "integer", false)
	columns := map[string]*Column{
		"id": col,
	}

	tbl := NewTable("users", columns)

	if len(tbl.Columns()) != 1 {
		t.Errorf("expected 1 column, got %d", len(tbl.Columns()))
	}
	if tbl.Columns()["id"] != col {
		t.Error("expected id column to be set")
	}
}

func TestTableSchema(t *testing.T) {
	tbl := NewTable("users", nil)

	if tbl.Schema() != nil {
		t.Error("expected nil schema initially")
	}

	schema := NewSchema("public", "admin", nil)
	tbl.SetSchema(schema)

	if tbl.Schema() == nil {
		t.Fatal("expected schema to be set")
	}
	if tbl.Schema().Name() != "public" {
		t.Errorf("expected schema name 'public', got %q", tbl.Schema().Name())
	}
}

func TestTableAddColumn(t *testing.T) {
	tbl := NewTable("users", nil)
	col1 := NewColumn("id", "integer", false)
	col2 := NewColumn("name", "varchar", false)

	tbl.AddColumn(col1)
	tbl.AddColumn(col2)

	columns := tbl.Columns()
	if len(columns) != 2 {
		t.Fatalf("expected 2 columns, got %d", len(columns))
	}
	if columns["id"] == nil {
		t.Error("expected id column to exist")
	}
	if columns["name"] == nil {
		t.Error("expected name column to exist")
	}

	// Verify column has table reference set
	if col1.Table() != tbl {
		t.Error("expected column to have table reference")
	}
}

func TestTablePrimaryKey(t *testing.T) {
	tbl := NewTable("users", nil)

	if tbl.PrimaryKey() != nil {
		t.Error("expected nil primary key initially")
	}

	col := NewColumn("id", "integer", false)
	pk := NewPrimaryKey("pk_users", tbl, []*Column{col})
	tbl.SetPrimaryKey(pk)

	if tbl.PrimaryKey() == nil {
		t.Fatal("expected primary key to be set")
	}
	if tbl.PrimaryKey().Name() != "pk_users" {
		t.Errorf("expected pk name 'pk_users', got %q", tbl.PrimaryKey().Name())
	}
}

func TestTableAddForeignKey(t *testing.T) {
	tbl := NewTable("orders", nil)
	fk1 := NewForeignKey("fk_user", "users")
	fk2 := NewForeignKey("fk_product", "products")

	tbl.AddForeignKey(fk1)
	tbl.AddForeignKey(fk2)

	fks := tbl.ForeignKeys()
	if len(fks) != 2 {
		t.Fatalf("expected 2 foreign keys, got %d", len(fks))
	}
	if fks[0].Name() != "fk_user" {
		t.Errorf("expected first fk 'fk_user', got %q", fks[0].Name())
	}

	// Verify foreign key has table reference set
	if fk1.Table() != tbl {
		t.Error("expected foreign key to have table reference")
	}
}

func TestTableAddIndex(t *testing.T) {
	tbl := NewTable("users", nil)
	col := NewColumn("email", "varchar", false)
	idx1 := NewIndex("idx_email", tbl, []*Column{col}, true)
	idx2 := NewIndex("idx_name", tbl, nil, false)

	tbl.AddIndex(idx1)
	tbl.AddIndex(idx2)

	indexes := tbl.Indexes()
	if len(indexes) != 2 {
		t.Fatalf("expected 2 indexes, got %d", len(indexes))
	}
	if indexes[0].Name() != "idx_email" {
		t.Errorf("expected first index 'idx_email', got %q", indexes[0].Name())
	}

	// Verify index has table reference set
	if idx1.Table() != tbl {
		t.Error("expected index to have table reference")
	}
}

func TestTableAddConstraint(t *testing.T) {
	tbl := NewTable("users", nil)
	c1 := NewConstraint("chk_email", ConstraintTypeCheck)
	c2 := NewConstraint("uq_username", ConstraintTypeUnique)

	tbl.AddConstraint(c1)
	tbl.AddConstraint(c2)

	constraints := tbl.Constraints()
	if len(constraints) != 2 {
		t.Fatalf("expected 2 constraints, got %d", len(constraints))
	}
	if constraints[0].Name() != "chk_email" {
		t.Errorf("expected first constraint 'chk_email', got %q", constraints[0].Name())
	}

	// Verify constraint has table reference set
	if c1.Table() != tbl {
		t.Error("expected constraint to have table reference")
	}
}

func TestTableAddTrigger(t *testing.T) {
	tbl := NewTable("users", nil)
	tr1 := NewTrigger("trg_audit", "EXECUTE audit_fn()")
	tr2 := NewTrigger("trg_updated", "EXECUTE update_timestamp()")

	tbl.AddTrigger(tr1)
	tbl.AddTrigger(tr2)

	triggers := tbl.Triggers()
	if len(triggers) != 2 {
		t.Fatalf("expected 2 triggers, got %d", len(triggers))
	}
	if triggers[0].Name() != "trg_audit" {
		t.Errorf("expected first trigger 'trg_audit', got %q", triggers[0].Name())
	}

	// Verify trigger has table reference set
	if tr1.Table() != tbl {
		t.Error("expected trigger to have table reference")
	}
}

func TestTableFullyQualifiedName(t *testing.T) {
	t.Run("without schema", func(t *testing.T) {
		tbl := NewTable("users", nil)

		if tbl.FullyQualifiedName() != "users" {
			t.Errorf("expected 'users', got %q", tbl.FullyQualifiedName())
		}
	})

	t.Run("with schema", func(t *testing.T) {
		tbl := NewTable("users", nil)
		schema := NewSchema("public", "admin", nil)
		tbl.SetSchema(schema)

		if tbl.FullyQualifiedName() != "public.users" {
			t.Errorf("expected 'public.users', got %q", tbl.FullyQualifiedName())
		}
	})
}

func TestTableMarshalJSON(t *testing.T) {
	tbl := NewTable("users", nil)

	col := NewColumn("id", "integer", false)
	tbl.AddColumn(col)

	pk := NewPrimaryKey("pk_users", tbl, []*Column{col})
	tbl.SetPrimaryKey(pk)

	fk := NewForeignKey("fk_dept", "departments")
	tbl.AddForeignKey(fk)

	idx := NewIndex("idx_id", tbl, []*Column{col}, true)
	tbl.AddIndex(idx)

	constraint := NewConstraint("chk_id", ConstraintTypeCheck)
	tbl.AddConstraint(constraint)

	trigger := NewTrigger("trg_audit", "EXECUTE audit()")
	tbl.AddTrigger(trigger)

	data, err := json.Marshal(tbl)
	if err != nil {
		t.Fatalf("failed to marshal table: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("failed to unmarshal json: %v", err)
	}

	if result["name"] != "users" {
		t.Errorf("expected name 'users', got %v", result["name"])
	}

	columns := result["columns"].(map[string]interface{})
	if len(columns) != 1 {
		t.Errorf("expected 1 column, got %d", len(columns))
	}

	if result["primaryKey"] == nil {
		t.Error("expected primaryKey to be present")
	}

	fks := result["foreignKeys"].([]interface{})
	if len(fks) != 1 {
		t.Errorf("expected 1 foreign key, got %d", len(fks))
	}

	indexes := result["indexes"].([]interface{})
	if len(indexes) != 1 {
		t.Errorf("expected 1 index, got %d", len(indexes))
	}

	constraints := result["constraints"].([]interface{})
	if len(constraints) != 1 {
		t.Errorf("expected 1 constraint, got %d", len(constraints))
	}

	triggers := result["triggers"].([]interface{})
	if len(triggers) != 1 {
		t.Errorf("expected 1 trigger, got %d", len(triggers))
	}
}

func TestTableMarshalJSONOmitsEmpty(t *testing.T) {
	tbl := NewTable("empty_table", nil)

	data, err := json.Marshal(tbl)
	if err != nil {
		t.Fatalf("failed to marshal table: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("failed to unmarshal json: %v", err)
	}

	// primaryKey should be omitted when nil
	if _, exists := result["primaryKey"]; exists {
		t.Error("expected primaryKey to be omitted")
	}

	// Empty slices should be omitted with omitempty
	if fks, exists := result["foreignKeys"]; exists {
		fkSlice := fks.([]interface{})
		if len(fkSlice) != 0 {
			t.Errorf("expected foreignKeys to be empty or omitted, got %v", fks)
		}
	}
}
