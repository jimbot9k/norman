package reports

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	dbo "github.com/jimbot9k/norman/internal/core/dbobjects"
)

// =============================================================================
// WriteInventoryReport Tests
// =============================================================================

func TestJSONReportWriter_WriteInventoryReport(t *testing.T) {
	t.Run("writes JSON file successfully", func(t *testing.T) {
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "test.json")

		db := dbo.NewDatabase("testdb", nil)
		schema := dbo.NewSchema("public", "owner", nil)
		db.AddSchema(schema)

		writer := &JSONReportWriter{}
		err := writer.WriteInventoryReport(filePath, db)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		content, err := os.ReadFile(filePath)
		if err != nil {
			t.Fatalf("failed to read file: %v", err)
		}

		// Verify it's valid JSON
		var result map[string]interface{}
		if err := json.Unmarshal(content, &result); err != nil {
			t.Errorf("expected valid JSON, got error: %v", err)
		}

		// Check content has name
		if result["name"] != "testdb" {
			t.Errorf("expected name 'testdb', got %v", result["name"])
		}
	})

	t.Run("returns error for invalid path", func(t *testing.T) {
		writer := &JSONReportWriter{}
		db := dbo.NewDatabase("testdb", nil)

		err := writer.WriteInventoryReport("/nonexistent/path/file.json", db)

		if err == nil {
			t.Error("expected error for invalid path")
		}
	})

	t.Run("writes indented JSON", func(t *testing.T) {
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "test.json")

		db := dbo.NewDatabase("testdb", nil)

		writer := &JSONReportWriter{}
		err := writer.WriteInventoryReport(filePath, db)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		content, err := os.ReadFile(filePath)
		if err != nil {
			t.Fatalf("failed to read file: %v", err)
		}

		// Indented JSON should contain newlines
		if !strings.Contains(string(content), "\n") {
			t.Error("expected indented JSON with newlines")
		}
	})
}

// =============================================================================
// columnToJSON Tests
// =============================================================================

func TestColumnToJSON(t *testing.T) {
	t.Run("basic column", func(t *testing.T) {
		col := dbo.NewColumn("id", "integer", false)

		result := columnToJSON(col)

		if result.Name != "id" {
			t.Errorf("expected name 'id', got %s", result.Name)
		}
		if result.DataType != "integer" {
			t.Errorf("expected dataType 'integer', got %s", result.DataType)
		}
		if result.Nullable != false {
			t.Error("expected nullable false")
		}
	})

	t.Run("nullable column", func(t *testing.T) {
		col := dbo.NewColumn("email", "varchar", true)

		result := columnToJSON(col)

		if !result.Nullable {
			t.Error("expected nullable true")
		}
	})

	t.Run("column with default value", func(t *testing.T) {
		col := dbo.NewColumn("status", "varchar", false)
		col.SetDefaultValue("active")

		result := columnToJSON(col)

		if result.DefaultValue == nil || *result.DefaultValue != "active" {
			t.Error("expected default value 'active'")
		}
	})

	t.Run("column with ordinal position", func(t *testing.T) {
		col := dbo.NewColumn("name", "varchar", false)
		col.SetOrdinalPosition(5)

		result := columnToJSON(col)

		if result.OrdinalPosition != 5 {
			t.Errorf("expected ordinal position 5, got %d", result.OrdinalPosition)
		}
	})

	t.Run("column with char max length", func(t *testing.T) {
		col := dbo.NewColumn("name", "varchar", false)
		col.SetCharMaxLength(255)

		result := columnToJSON(col)

		if result.CharMaxLength == nil || *result.CharMaxLength != 255 {
			t.Error("expected char max length 255")
		}
	})

	t.Run("column with numeric precision and scale", func(t *testing.T) {
		col := dbo.NewColumn("price", "decimal", false)
		col.SetNumericPrecision(10)
		col.SetNumericScale(2)

		result := columnToJSON(col)

		if result.NumericPrecision == nil || *result.NumericPrecision != 10 {
			t.Error("expected numeric precision 10")
		}
		if result.NumericScale == nil || *result.NumericScale != 2 {
			t.Error("expected numeric scale 2")
		}
	})
}

// =============================================================================
// constraintToJSON Tests
// =============================================================================

func TestConstraintToJSON(t *testing.T) {
	t.Run("check constraint", func(t *testing.T) {
		constraint := dbo.NewConstraint("price_check", dbo.ConstraintTypeCheck)
		constraint.SetCheckExpression("price > 0")

		col := dbo.NewColumn("price", "decimal", false)
		constraint.AddColumn(col)

		result := constraintToJSON(constraint)

		if result.Name != "price_check" {
			t.Errorf("expected name 'price_check', got %s", result.Name)
		}
		if result.Type != dbo.ConstraintTypeCheck {
			t.Errorf("expected type CHECK, got %v", result.Type)
		}
		if result.CheckExpression != "price > 0" {
			t.Errorf("expected expression 'price > 0', got %s", result.CheckExpression)
		}
		if len(result.Columns) != 1 || result.Columns[0] != "price" {
			t.Error("expected columns [price]")
		}
	})

	t.Run("unique constraint", func(t *testing.T) {
		constraint := dbo.NewConstraint("email_unique", dbo.ConstraintTypeUnique)

		col := dbo.NewColumn("email", "varchar", false)
		constraint.AddColumn(col)

		result := constraintToJSON(constraint)

		if result.Type != dbo.ConstraintTypeUnique {
			t.Errorf("expected type UNIQUE, got %v", result.Type)
		}
	})

	t.Run("not null constraint", func(t *testing.T) {
		constraint := dbo.NewConstraint("name_not_null", dbo.ConstraintTypeNotNull)

		col := dbo.NewColumn("name", "varchar", false)
		constraint.AddColumn(col)

		result := constraintToJSON(constraint)

		if result.Type != dbo.ConstraintTypeNotNull {
			t.Errorf("expected type NOT NULL, got %v", result.Type)
		}
	})

	t.Run("constraint with multiple columns", func(t *testing.T) {
		constraint := dbo.NewConstraint("composite_unique", dbo.ConstraintTypeUnique)

		col1 := dbo.NewColumn("first_name", "varchar", false)
		col2 := dbo.NewColumn("last_name", "varchar", false)
		constraint.AddColumn(col1)
		constraint.AddColumn(col2)

		result := constraintToJSON(constraint)

		if len(result.Columns) != 2 {
			t.Errorf("expected 2 columns, got %d", len(result.Columns))
		}
		if result.Columns[0] != "first_name" || result.Columns[1] != "last_name" {
			t.Error("expected columns [first_name, last_name]")
		}
	})
}

// =============================================================================
// foreignKeyToJSON Tests
// =============================================================================

func TestForeignKeyToJSON(t *testing.T) {
	t.Run("basic foreign key", func(t *testing.T) {
		childCol := dbo.NewColumn("department_id", "integer", false)
		parentCol := dbo.NewColumn("id", "integer", false)

		fk := dbo.NewForeignKey("fk_dept", "departments")
		fk.AddColumn(childCol)
		fk.AddReferencedColumn(parentCol)
		fk.SetReferencedSchema("public")

		result := foreignKeyToJSON(fk)

		if result.Name != "fk_dept" {
			t.Errorf("expected name 'fk_dept', got %s", result.Name)
		}
		if result.ReferencedTable != "departments" {
			t.Errorf("expected referenced table 'departments', got %s", result.ReferencedTable)
		}
		if result.ReferencedSchema != "public" {
			t.Errorf("expected referenced schema 'public', got %s", result.ReferencedSchema)
		}
		if len(result.Columns) != 1 || result.Columns[0] != "department_id" {
			t.Error("expected columns [department_id]")
		}
		if len(result.ReferencedColumns) != 1 || result.ReferencedColumns[0] != "id" {
			t.Error("expected referenced columns [id]")
		}
	})

	t.Run("foreign key with actions", func(t *testing.T) {
		childCol := dbo.NewColumn("user_id", "integer", false)
		parentCol := dbo.NewColumn("id", "integer", false)

		fk := dbo.NewForeignKey("fk_user", "users")
		fk.AddColumn(childCol)
		fk.AddReferencedColumn(parentCol)
		fk.SetOnDelete(dbo.ActionCascade)
		fk.SetOnUpdate(dbo.ActionSetNull)

		result := foreignKeyToJSON(fk)

		if result.OnDelete != dbo.ActionCascade {
			t.Errorf("expected onDelete CASCADE, got %v", result.OnDelete)
		}
		if result.OnUpdate != dbo.ActionSetNull {
			t.Errorf("expected onUpdate SET NULL, got %v", result.OnUpdate)
		}
	})

	t.Run("composite foreign key", func(t *testing.T) {
		col1 := dbo.NewColumn("tenant_id", "integer", false)
		col2 := dbo.NewColumn("user_id", "integer", false)
		refCol1 := dbo.NewColumn("tenant_id", "integer", false)
		refCol2 := dbo.NewColumn("id", "integer", false)

		fk := dbo.NewForeignKey("fk_tenant_user", "users")
		fk.AddColumn(col1)
		fk.AddColumn(col2)
		fk.AddReferencedColumn(refCol1)
		fk.AddReferencedColumn(refCol2)

		result := foreignKeyToJSON(fk)

		if len(result.Columns) != 2 {
			t.Errorf("expected 2 columns, got %d", len(result.Columns))
		}
		if len(result.ReferencedColumns) != 2 {
			t.Errorf("expected 2 referenced columns, got %d", len(result.ReferencedColumns))
		}
	})
}

// =============================================================================
// functionParameterToJSON Tests
// =============================================================================

func TestFunctionParameterToJSON(t *testing.T) {
	t.Run("IN parameter", func(t *testing.T) {
		param := dbo.NewFunctionParameter("user_id", "integer", dbo.ParameterModeIn)

		result := functionParameterToJSON(param)

		if result.Name != "user_id" {
			t.Errorf("expected name 'user_id', got %s", result.Name)
		}
		if result.DataType != "integer" {
			t.Errorf("expected dataType 'integer', got %s", result.DataType)
		}
		if result.Mode != dbo.ParameterModeIn {
			t.Errorf("expected mode IN, got %v", result.Mode)
		}
	})

	t.Run("OUT parameter", func(t *testing.T) {
		param := dbo.NewFunctionParameter("result", "text", dbo.ParameterModeOut)

		result := functionParameterToJSON(param)

		if result.Mode != dbo.ParameterModeOut {
			t.Errorf("expected mode OUT, got %v", result.Mode)
		}
	})

	t.Run("INOUT parameter", func(t *testing.T) {
		param := dbo.NewFunctionParameter("counter", "integer", dbo.ParameterModeInOut)

		result := functionParameterToJSON(param)

		if result.Mode != dbo.ParameterModeInOut {
			t.Errorf("expected mode INOUT, got %v", result.Mode)
		}
	})
}

// =============================================================================
// functionToJSON Tests
// =============================================================================

func TestFunctionToJSON(t *testing.T) {
	t.Run("basic function", func(t *testing.T) {
		fn := dbo.NewFunction("get_user", "SELECT * FROM users WHERE id = $1")
		fn.SetReturnType("setof users")
		fn.SetLanguage("sql")

		result := functionToJSON(fn)

		if result.Name != "get_user" {
			t.Errorf("expected name 'get_user', got %s", result.Name)
		}
		if result.Definition != "SELECT * FROM users WHERE id = $1" {
			t.Errorf("expected definition, got %s", result.Definition)
		}
		if result.ReturnType != "setof users" {
			t.Errorf("expected return type 'setof users', got %s", result.ReturnType)
		}
		if result.Language != "sql" {
			t.Errorf("expected language 'sql', got %s", result.Language)
		}
	})

	t.Run("function with parameters", func(t *testing.T) {
		fn := dbo.NewFunction("add_numbers", "SELECT $1 + $2")
		fn.SetReturnType("integer")
		fn.AddParameter(dbo.NewFunctionParameter("a", "integer", dbo.ParameterModeIn))
		fn.AddParameter(dbo.NewFunctionParameter("b", "integer", dbo.ParameterModeIn))

		result := functionToJSON(fn)

		if len(result.Parameters) != 2 {
			t.Errorf("expected 2 parameters, got %d", len(result.Parameters))
		}
		if result.Parameters[0].Name != "a" || result.Parameters[1].Name != "b" {
			t.Error("expected parameters [a, b]")
		}
	})

	t.Run("function with no parameters", func(t *testing.T) {
		fn := dbo.NewFunction("now_utc", "SELECT NOW() AT TIME ZONE 'UTC'")
		fn.SetReturnType("timestamp")

		result := functionToJSON(fn)

		if len(result.Parameters) != 0 {
			t.Errorf("expected 0 parameters, got %d", len(result.Parameters))
		}
	})

	t.Run("plpgsql function", func(t *testing.T) {
		fn := dbo.NewFunction("increment", "BEGIN RETURN i + 1; END;")
		fn.SetLanguage("plpgsql")
		fn.SetReturnType("integer")

		result := functionToJSON(fn)

		if result.Language != "plpgsql" {
			t.Errorf("expected language 'plpgsql', got %s", result.Language)
		}
	})
}

// =============================================================================
// indexToJSON Tests
// =============================================================================

func TestIndexToJSON(t *testing.T) {
	t.Run("unique index", func(t *testing.T) {
		table := dbo.NewTable("users", nil)
		col := dbo.NewColumn("email", "varchar", false)
		table.AddColumn(col)

		idx := dbo.NewIndex("idx_users_email", table, []*dbo.Column{col}, true)

		result := indexToJSON(idx)

		if result.Name != "idx_users_email" {
			t.Errorf("expected name 'idx_users_email', got %s", result.Name)
		}
		if len(result.Columns) != 1 || result.Columns[0] != "email" {
			t.Error("expected columns [email]")
		}
		if !result.IsUnique {
			t.Error("expected isUnique true")
		}
		if result.IndexType != dbo.IndexTypeBTree {
			t.Errorf("expected index type btree, got %v", result.IndexType)
		}
	})

	t.Run("non-unique index", func(t *testing.T) {
		table := dbo.NewTable("users", nil)
		col := dbo.NewColumn("name", "varchar", false)
		table.AddColumn(col)

		idx := dbo.NewIndex("idx_users_name", table, []*dbo.Column{col}, false)

		result := indexToJSON(idx)

		if result.IsUnique {
			t.Error("expected isUnique false")
		}
	})

	t.Run("primary index", func(t *testing.T) {
		table := dbo.NewTable("users", nil)
		col := dbo.NewColumn("id", "integer", false)
		table.AddColumn(col)

		idx := dbo.NewIndex("users_pkey", table, []*dbo.Column{col}, true)
		idx.SetPrimary(true)

		result := indexToJSON(idx)

		if !result.IsPrimary {
			t.Error("expected isPrimary true")
		}
	})

	t.Run("composite index", func(t *testing.T) {
		table := dbo.NewTable("orders", nil)
		col1 := dbo.NewColumn("user_id", "integer", false)
		col2 := dbo.NewColumn("created_at", "timestamp", false)
		table.AddColumn(col1)
		table.AddColumn(col2)

		idx := dbo.NewIndex("idx_orders_user_date", table, []*dbo.Column{col1, col2}, false)

		result := indexToJSON(idx)

		if len(result.Columns) != 2 {
			t.Errorf("expected 2 columns, got %d", len(result.Columns))
		}
	})

	t.Run("hash index", func(t *testing.T) {
		table := dbo.NewTable("users", nil)
		col := dbo.NewColumn("email", "varchar", false)
		table.AddColumn(col)

		idx := dbo.NewIndex("idx_users_email_hash", table, []*dbo.Column{col}, false)
		idx.SetIndexType(dbo.IndexTypeHash)

		result := indexToJSON(idx)

		if result.IndexType != dbo.IndexTypeHash {
			t.Errorf("expected index type hash, got %v", result.IndexType)
		}
	})
}

// =============================================================================
// primaryKeyToJSON Tests
// =============================================================================

func TestPrimaryKeyToJSON(t *testing.T) {
	t.Run("nil primary key", func(t *testing.T) {
		result := primaryKeyToJSON(nil)

		if result != nil {
			t.Error("expected nil result for nil input")
		}
	})

	t.Run("single column primary key", func(t *testing.T) {
		table := dbo.NewTable("users", nil)
		col := dbo.NewColumn("id", "integer", false)
		table.AddColumn(col)

		pk := dbo.NewPrimaryKey("users_pkey", table, []*dbo.Column{col})

		result := primaryKeyToJSON(pk)

		if result == nil {
			t.Fatal("expected non-nil result")
		}
		if result.Name != "users_pkey" {
			t.Errorf("expected name 'users_pkey', got %s", result.Name)
		}
		if len(result.Columns) != 1 || result.Columns[0] != "id" {
			t.Error("expected columns [id]")
		}
	})

	t.Run("composite primary key", func(t *testing.T) {
		table := dbo.NewTable("order_items", nil)
		col1 := dbo.NewColumn("order_id", "integer", false)
		col2 := dbo.NewColumn("item_id", "integer", false)
		table.AddColumn(col1)
		table.AddColumn(col2)

		pk := dbo.NewPrimaryKey("order_items_pkey", table, []*dbo.Column{col1, col2})

		result := primaryKeyToJSON(pk)

		if len(result.Columns) != 2 {
			t.Errorf("expected 2 columns, got %d", len(result.Columns))
		}
		if result.Columns[0] != "order_id" || result.Columns[1] != "item_id" {
			t.Error("expected columns [order_id, item_id]")
		}
	})
}

// =============================================================================
// procedureToJSON Tests
// =============================================================================

func TestProcedureToJSON(t *testing.T) {
	t.Run("basic procedure", func(t *testing.T) {
		proc := dbo.NewProcedure("update_stats", "UPDATE stats SET count = count + 1")
		proc.SetLanguage("sql")

		result := procedureToJSON(proc)

		if result.Name != "update_stats" {
			t.Errorf("expected name 'update_stats', got %s", result.Name)
		}
		if result.Definition != "UPDATE stats SET count = count + 1" {
			t.Errorf("expected definition, got %s", result.Definition)
		}
		if result.Language != "sql" {
			t.Errorf("expected language 'sql', got %s", result.Language)
		}
	})

	t.Run("procedure with parameters", func(t *testing.T) {
		proc := dbo.NewProcedure("transfer_funds", "BEGIN ... END;")
		proc.SetLanguage("plpgsql")
		proc.AddParameter(dbo.NewFunctionParameter("from_account", "integer", dbo.ParameterModeIn))
		proc.AddParameter(dbo.NewFunctionParameter("to_account", "integer", dbo.ParameterModeIn))
		proc.AddParameter(dbo.NewFunctionParameter("amount", "decimal", dbo.ParameterModeIn))

		result := procedureToJSON(proc)

		if len(result.Parameters) != 3 {
			t.Errorf("expected 3 parameters, got %d", len(result.Parameters))
		}
	})

	t.Run("procedure with no parameters", func(t *testing.T) {
		proc := dbo.NewProcedure("cleanup", "DELETE FROM temp_data")

		result := procedureToJSON(proc)

		if len(result.Parameters) != 0 {
			t.Errorf("expected 0 parameters, got %d", len(result.Parameters))
		}
	})
}

// =============================================================================
// sequenceToJSON Tests
// =============================================================================

func TestSequenceToJSON(t *testing.T) {
	t.Run("basic sequence", func(t *testing.T) {
		seq := dbo.NewSequence("users_id_seq", 1, 1)

		result := sequenceToJSON(seq)

		if result.Name != "users_id_seq" {
			t.Errorf("expected name 'users_id_seq', got %s", result.Name)
		}
		if result.StartValue != 1 {
			t.Errorf("expected start value 1, got %d", result.StartValue)
		}
		if result.Increment != 1 {
			t.Errorf("expected increment 1, got %d", result.Increment)
		}
	})

	t.Run("sequence with custom values", func(t *testing.T) {
		seq := dbo.NewSequence("batch_id_seq", 100, 10)
		seq.SetMinValue(100)
		seq.SetMaxValue(1000000)
		seq.SetCache(20)
		seq.SetCycle(true)

		result := sequenceToJSON(seq)

		if result.StartValue != 100 {
			t.Errorf("expected start value 100, got %d", result.StartValue)
		}
		if result.Increment != 10 {
			t.Errorf("expected increment 10, got %d", result.Increment)
		}
		if result.MinValue != 100 {
			t.Errorf("expected min value 100, got %d", result.MinValue)
		}
		if result.MaxValue != 1000000 {
			t.Errorf("expected max value 1000000, got %d", result.MaxValue)
		}
		if result.Cache != 20 {
			t.Errorf("expected cache 20, got %d", result.Cache)
		}
		if !result.Cycle {
			t.Error("expected cycle true")
		}
	})

	t.Run("sequence with default max value", func(t *testing.T) {
		seq := dbo.NewSequence("auto_seq", 1, 1)

		result := sequenceToJSON(seq)

		// Default max value is max int64
		if result.MaxValue != 9223372036854775807 {
			t.Errorf("expected default max value, got %d", result.MaxValue)
		}
	})
}

// =============================================================================
// triggerToJSON Tests
// =============================================================================

func TestTriggerToJSON(t *testing.T) {
	t.Run("trigger without function", func(t *testing.T) {
		trigger := dbo.NewTrigger("trg_audit", "EXECUTE PROCEDURE audit_fn()")
		trigger.SetTiming(dbo.TriggerTimingAfter)
		trigger.AddEvent(dbo.TriggerEventInsert)
		trigger.AddEvent(dbo.TriggerEventUpdate)
		trigger.SetForEach("ROW")

		result := triggerToJSON(trigger)

		if result.Name != "trg_audit" {
			t.Errorf("expected name 'trg_audit', got %s", result.Name)
		}
		if result.Definition != "EXECUTE PROCEDURE audit_fn()" {
			t.Errorf("expected definition, got %s", result.Definition)
		}
		if result.Timing != dbo.TriggerTimingAfter {
			t.Errorf("expected timing AFTER, got %v", result.Timing)
		}
		if len(result.Events) != 2 {
			t.Errorf("expected 2 events, got %d", len(result.Events))
		}
		if result.ForEach != "ROW" {
			t.Errorf("expected forEach 'ROW', got %s", result.ForEach)
		}
		if result.Function != "" {
			t.Errorf("expected empty function, got %s", result.Function)
		}
	})

	t.Run("trigger with function", func(t *testing.T) {
		fn := dbo.NewFunction("audit_trigger_fn", "BEGIN ... END;")
		fn.SetLanguage("plpgsql")
		fn.SetReturnType("trigger")

		trigger := dbo.NewTrigger("trg_audit", "EXECUTE FUNCTION audit_trigger_fn()")
		trigger.SetTiming(dbo.TriggerTimingBefore)
		trigger.AddEvent(dbo.TriggerEventDelete)
		trigger.SetFunction(fn)

		result := triggerToJSON(trigger)

		if result.Function != "audit_trigger_fn" {
			t.Errorf("expected function 'audit_trigger_fn', got %s", result.Function)
		}
	})

	t.Run("instead of trigger", func(t *testing.T) {
		trigger := dbo.NewTrigger("trg_view_insert", "EXECUTE PROCEDURE handle_insert()")
		trigger.SetTiming(dbo.TriggerTimingInsteadOf)
		trigger.AddEvent(dbo.TriggerEventInsert)
		trigger.SetForEach("ROW")

		result := triggerToJSON(trigger)

		if result.Timing != dbo.TriggerTimingInsteadOf {
			t.Errorf("expected timing INSTEAD OF, got %v", result.Timing)
		}
	})

	t.Run("statement level trigger", func(t *testing.T) {
		trigger := dbo.NewTrigger("trg_table_modified", "NOTIFY table_changed")
		trigger.SetTiming(dbo.TriggerTimingAfter)
		trigger.AddEvent(dbo.TriggerEventUpdate)
		trigger.SetForEach("STATEMENT")

		result := triggerToJSON(trigger)

		if result.ForEach != "STATEMENT" {
			t.Errorf("expected forEach 'STATEMENT', got %s", result.ForEach)
		}
	})
}

// =============================================================================
// viewToJSON Tests
// =============================================================================

func TestViewToJSON(t *testing.T) {
	t.Run("basic view", func(t *testing.T) {
		view := dbo.NewView("active_users", "SELECT * FROM users WHERE active = true")

		result := viewToJSON(view)

		if result.Name != "active_users" {
			t.Errorf("expected name 'active_users', got %s", result.Name)
		}
		if result.Definition != "SELECT * FROM users WHERE active = true" {
			t.Errorf("expected definition, got %s", result.Definition)
		}
	})

	t.Run("view with columns", func(t *testing.T) {
		view := dbo.NewView("user_summary", "SELECT id, name FROM users")
		view.AddColumn(dbo.NewColumn("id", "integer", false))
		view.AddColumn(dbo.NewColumn("name", "varchar", false))

		result := viewToJSON(view)

		if len(result.Columns) != 2 {
			t.Errorf("expected 2 columns, got %d", len(result.Columns))
		}
		if result.Columns[0].Name != "id" || result.Columns[1].Name != "name" {
			t.Error("expected columns [id, name]")
		}
	})

	t.Run("view without columns", func(t *testing.T) {
		view := dbo.NewView("empty_view", "SELECT 1")

		result := viewToJSON(view)

		if len(result.Columns) != 0 {
			t.Errorf("expected 0 columns, got %d", len(result.Columns))
		}
	})
}

// =============================================================================
// tableToJSON Tests
// =============================================================================

func TestTableToJSON(t *testing.T) {
	t.Run("basic table with columns", func(t *testing.T) {
		table := dbo.NewTable("users", nil)
		table.AddColumn(dbo.NewColumn("id", "integer", false))
		table.AddColumn(dbo.NewColumn("name", "varchar", false))

		result := tableToJSON(table)

		if result.Name != "users" {
			t.Errorf("expected name 'users', got %s", result.Name)
		}
		if len(result.Columns) != 2 {
			t.Errorf("expected 2 columns, got %d", len(result.Columns))
		}
	})

	t.Run("table with primary key", func(t *testing.T) {
		table := dbo.NewTable("users", nil)
		col := dbo.NewColumn("id", "integer", false)
		table.AddColumn(col)
		pk := dbo.NewPrimaryKey("users_pkey", table, []*dbo.Column{col})
		table.SetPrimaryKey(pk)

		result := tableToJSON(table)

		if result.PrimaryKey == nil {
			t.Fatal("expected primary key")
		}
		if result.PrimaryKey.Name != "users_pkey" {
			t.Errorf("expected pk name 'users_pkey', got %s", result.PrimaryKey.Name)
		}
	})

	t.Run("table without primary key", func(t *testing.T) {
		table := dbo.NewTable("logs", nil)
		table.AddColumn(dbo.NewColumn("message", "text", false))

		result := tableToJSON(table)

		if result.PrimaryKey != nil {
			t.Error("expected no primary key")
		}
	})

	t.Run("table with foreign keys", func(t *testing.T) {
		table := dbo.NewTable("orders", nil)
		userCol := dbo.NewColumn("user_id", "integer", false)
		refCol := dbo.NewColumn("id", "integer", false)
		table.AddColumn(userCol)

		fk := dbo.NewForeignKey("fk_user", "users")
		fk.SetTable(table)
		fk.AddColumn(userCol)
		fk.AddReferencedColumn(refCol)
		table.AddForeignKey(fk)

		result := tableToJSON(table)

		if len(result.ForeignKeys) != 1 {
			t.Errorf("expected 1 foreign key, got %d", len(result.ForeignKeys))
		}
	})

	t.Run("table with indexes", func(t *testing.T) {
		table := dbo.NewTable("users", nil)
		col := dbo.NewColumn("email", "varchar", false)
		table.AddColumn(col)

		idx := dbo.NewIndex("idx_email", table, []*dbo.Column{col}, true)
		table.AddIndex(idx)

		result := tableToJSON(table)

		if len(result.Indexes) != 1 {
			t.Errorf("expected 1 index, got %d", len(result.Indexes))
		}
	})

	t.Run("table with constraints", func(t *testing.T) {
		table := dbo.NewTable("products", nil)
		priceCol := dbo.NewColumn("price", "decimal", false)
		table.AddColumn(priceCol)

		constraint := dbo.NewConstraint("price_positive", dbo.ConstraintTypeCheck)
		constraint.SetCheckExpression("price > 0")
		constraint.AddColumn(priceCol)
		table.AddConstraint(constraint)

		result := tableToJSON(table)

		if len(result.Constraints) != 1 {
			t.Errorf("expected 1 constraint, got %d", len(result.Constraints))
		}
	})

	t.Run("table with triggers", func(t *testing.T) {
		table := dbo.NewTable("orders", nil)
		table.AddColumn(dbo.NewColumn("id", "integer", false))

		trigger := dbo.NewTrigger("trg_order_created", "NOTIFY order_created")
		trigger.SetTiming(dbo.TriggerTimingAfter)
		trigger.AddEvent(dbo.TriggerEventInsert)
		table.AddTrigger(trigger)

		result := tableToJSON(table)

		if len(result.Triggers) != 1 {
			t.Errorf("expected 1 trigger, got %d", len(result.Triggers))
		}
	})

	t.Run("full table with all components", func(t *testing.T) {
		table := dbo.NewTable("employees", nil)

		// Add columns
		idCol := dbo.NewColumn("id", "integer", false)
		nameCol := dbo.NewColumn("name", "varchar", false)
		deptCol := dbo.NewColumn("department_id", "integer", true)
		table.AddColumn(idCol)
		table.AddColumn(nameCol)
		table.AddColumn(deptCol)

		// Add primary key
		pk := dbo.NewPrimaryKey("employees_pkey", table, []*dbo.Column{idCol})
		table.SetPrimaryKey(pk)

		// Add foreign key
		refCol := dbo.NewColumn("id", "integer", false)
		fk := dbo.NewForeignKey("fk_dept", "departments")
		fk.SetTable(table)
		fk.AddColumn(deptCol)
		fk.AddReferencedColumn(refCol)
		table.AddForeignKey(fk)

		// Add index
		idx := dbo.NewIndex("idx_name", table, []*dbo.Column{nameCol}, false)
		table.AddIndex(idx)

		// Add constraint
		constraint := dbo.NewConstraint("name_not_empty", dbo.ConstraintTypeCheck)
		constraint.SetCheckExpression("name <> ''")
		constraint.AddColumn(nameCol)
		table.AddConstraint(constraint)

		// Add trigger
		trigger := dbo.NewTrigger("trg_audit", "EXECUTE PROCEDURE audit_fn()")
		trigger.SetTiming(dbo.TriggerTimingAfter)
		trigger.AddEvent(dbo.TriggerEventUpdate)
		table.AddTrigger(trigger)

		result := tableToJSON(table)

		if len(result.Columns) != 3 {
			t.Errorf("expected 3 columns, got %d", len(result.Columns))
		}
		if result.PrimaryKey == nil {
			t.Error("expected primary key")
		}
		if len(result.ForeignKeys) != 1 {
			t.Errorf("expected 1 foreign key, got %d", len(result.ForeignKeys))
		}
		if len(result.Indexes) != 1 {
			t.Errorf("expected 1 index, got %d", len(result.Indexes))
		}
		if len(result.Constraints) != 1 {
			t.Errorf("expected 1 constraint, got %d", len(result.Constraints))
		}
		if len(result.Triggers) != 1 {
			t.Errorf("expected 1 trigger, got %d", len(result.Triggers))
		}
	})
}

// =============================================================================
// schemaToJSON Tests
// =============================================================================

func TestSchemaToJSON(t *testing.T) {
	t.Run("empty schema", func(t *testing.T) {
		schema := dbo.NewSchema("public", "postgres", nil)

		result := schemaToJSON(schema)

		if result.Name != "public" {
			t.Errorf("expected name 'public', got %s", result.Name)
		}
		if result.Owner != "postgres" {
			t.Errorf("expected owner 'postgres', got %s", result.Owner)
		}
		if len(result.Tables) != 0 {
			t.Errorf("expected 0 tables, got %d", len(result.Tables))
		}
	})

	t.Run("schema with tables", func(t *testing.T) {
		schema := dbo.NewSchema("public", "postgres", nil)

		table1 := dbo.NewTable("users", nil)
		table1.AddColumn(dbo.NewColumn("id", "integer", false))
		schema.AddTable(table1)

		table2 := dbo.NewTable("orders", nil)
		table2.AddColumn(dbo.NewColumn("id", "integer", false))
		schema.AddTable(table2)

		result := schemaToJSON(schema)

		if len(result.Tables) != 2 {
			t.Errorf("expected 2 tables, got %d", len(result.Tables))
		}
	})

	t.Run("schema with views", func(t *testing.T) {
		schema := dbo.NewSchema("public", "postgres", nil)

		view := dbo.NewView("active_users", "SELECT * FROM users WHERE active")
		schema.AddView(view)

		result := schemaToJSON(schema)

		if len(result.Views) != 1 {
			t.Errorf("expected 1 view, got %d", len(result.Views))
		}
	})

	t.Run("schema with functions", func(t *testing.T) {
		schema := dbo.NewSchema("public", "postgres", nil)

		fn := dbo.NewFunction("get_user", "SELECT * FROM users WHERE id = $1")
		fn.SetReturnType("users")
		schema.AddFunction(fn)

		result := schemaToJSON(schema)

		if len(result.Functions) != 1 {
			t.Errorf("expected 1 function, got %d", len(result.Functions))
		}
	})

	t.Run("schema with procedures", func(t *testing.T) {
		schema := dbo.NewSchema("public", "postgres", nil)

		proc := dbo.NewProcedure("cleanup", "DELETE FROM temp")
		schema.AddProcedure(proc)

		result := schemaToJSON(schema)

		if len(result.Procedures) != 1 {
			t.Errorf("expected 1 procedure, got %d", len(result.Procedures))
		}
	})

	t.Run("schema with sequences", func(t *testing.T) {
		schema := dbo.NewSchema("public", "postgres", nil)

		seq := dbo.NewSequence("users_id_seq", 1, 1)
		schema.AddSequence(seq)

		result := schemaToJSON(schema)

		if len(result.Sequences) != 1 {
			t.Errorf("expected 1 sequence, got %d", len(result.Sequences))
		}
	})

	t.Run("schema with all components", func(t *testing.T) {
		schema := dbo.NewSchema("app", "app_owner", nil)

		// Add table
		table := dbo.NewTable("users", nil)
		table.AddColumn(dbo.NewColumn("id", "integer", false))
		schema.AddTable(table)

		// Add view
		view := dbo.NewView("user_list", "SELECT id FROM users")
		schema.AddView(view)

		// Add function
		fn := dbo.NewFunction("count_users", "SELECT COUNT(*) FROM users")
		fn.SetReturnType("integer")
		schema.AddFunction(fn)

		// Add procedure
		proc := dbo.NewProcedure("refresh_cache", "UPDATE cache SET updated = NOW()")
		schema.AddProcedure(proc)

		// Add sequence
		seq := dbo.NewSequence("users_id_seq", 1, 1)
		schema.AddSequence(seq)

		result := schemaToJSON(schema)

		if result.Name != "app" {
			t.Errorf("expected name 'app', got %s", result.Name)
		}
		if result.Owner != "app_owner" {
			t.Errorf("expected owner 'app_owner', got %s", result.Owner)
		}
		if len(result.Tables) != 1 {
			t.Errorf("expected 1 table, got %d", len(result.Tables))
		}
		if len(result.Views) != 1 {
			t.Errorf("expected 1 view, got %d", len(result.Views))
		}
		if len(result.Functions) != 1 {
			t.Errorf("expected 1 function, got %d", len(result.Functions))
		}
		if len(result.Procedures) != 1 {
			t.Errorf("expected 1 procedure, got %d", len(result.Procedures))
		}
		if len(result.Sequences) != 1 {
			t.Errorf("expected 1 sequence, got %d", len(result.Sequences))
		}
	})
}

// =============================================================================
// databaseToJSON Tests
// =============================================================================

func TestDatabaseToJSON(t *testing.T) {
	t.Run("empty database", func(t *testing.T) {
		db := dbo.NewDatabase("testdb", nil)

		result := databaseToJSON(db)

		if result.Name != "testdb" {
			t.Errorf("expected name 'testdb', got %s", result.Name)
		}
		if len(result.Schemas) != 0 {
			t.Errorf("expected 0 schemas, got %d", len(result.Schemas))
		}
	})

	t.Run("database with schemas", func(t *testing.T) {
		db := dbo.NewDatabase("testdb", nil)

		schema1 := dbo.NewSchema("public", "postgres", nil)
		db.AddSchema(schema1)

		schema2 := dbo.NewSchema("app", "app_owner", nil)
		db.AddSchema(schema2)

		result := databaseToJSON(db)

		if len(result.Schemas) != 2 {
			t.Errorf("expected 2 schemas, got %d", len(result.Schemas))
		}
	})

	t.Run("database with populated schema", func(t *testing.T) {
		db := dbo.NewDatabase("production", nil)

		schema := dbo.NewSchema("public", "admin", nil)
		table := dbo.NewTable("users", nil)
		table.AddColumn(dbo.NewColumn("id", "integer", false))
		schema.AddTable(table)
		db.AddSchema(schema)

		result := databaseToJSON(db)

		if result.Name != "production" {
			t.Errorf("expected name 'production', got %s", result.Name)
		}
		if len(result.Schemas) != 1 {
			t.Errorf("expected 1 schema, got %d", len(result.Schemas))
		}
		if len(result.Schemas[0].Tables) != 1 {
			t.Errorf("expected 1 table in schema, got %d", len(result.Schemas[0].Tables))
		}
	})
}

// =============================================================================
// marshalDatabaseIndent Tests
// =============================================================================

func TestMarshalDatabaseIndent(t *testing.T) {
	t.Run("marshal with default indent", func(t *testing.T) {
		db := dbo.NewDatabase("testdb", nil)

		data, err := marshalDatabaseIndent(db, "", "  ")

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Indented JSON should have newlines
		if !strings.Contains(string(data), "\n") {
			t.Error("expected indented JSON with newlines")
		}
	})

	t.Run("marshal with custom prefix and indent", func(t *testing.T) {
		db := dbo.NewDatabase("testdb", nil)

		data, err := marshalDatabaseIndent(db, ">>", "\t")

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Should have prefix and tabs
		if !strings.Contains(string(data), ">>") {
			t.Error("expected prefix in JSON")
		}
		if !strings.Contains(string(data), "\t") {
			t.Error("expected tab indent in JSON")
		}
	})

	t.Run("marshal complex database", func(t *testing.T) {
		db := dbo.NewDatabase("complex_db", nil)

		schema := dbo.NewSchema("public", "admin", nil)

		// Add table with various components
		table := dbo.NewTable("employees", nil)
		idCol := dbo.NewColumn("id", "serial", false)
		nameCol := dbo.NewColumn("name", "varchar", false)
		table.AddColumn(idCol)
		table.AddColumn(nameCol)

		pk := dbo.NewPrimaryKey("employees_pkey", table, []*dbo.Column{idCol})
		table.SetPrimaryKey(pk)

		schema.AddTable(table)

		// Add function
		fn := dbo.NewFunction("get_employee", "SELECT * FROM employees WHERE id = $1")
		fn.SetReturnType("employees")
		fn.AddParameter(dbo.NewFunctionParameter("emp_id", "integer", dbo.ParameterModeIn))
		schema.AddFunction(fn)

		// Add sequence
		seq := dbo.NewSequence("employees_id_seq", 1, 1)
		schema.AddSequence(seq)

		db.AddSchema(schema)

		data, err := marshalDatabaseIndent(db, "", "  ")

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Verify it's valid JSON
		var result map[string]interface{}
		if err := json.Unmarshal(data, &result); err != nil {
			t.Fatalf("invalid JSON: %v", err)
		}

		// Verify structure
		if result["name"] != "complex_db" {
			t.Errorf("expected name 'complex_db', got %v", result["name"])
		}

		schemas := result["schemas"].([]interface{})
		if len(schemas) != 1 {
			t.Errorf("expected 1 schema, got %d", len(schemas))
		}
	})
}

// =============================================================================
// Integration Tests
// =============================================================================

func TestIntegration_FullJSONReport(t *testing.T) {
	t.Run("complete database export", func(t *testing.T) {
		// Build a realistic database structure
		db := dbo.NewDatabase("ecommerce", nil)

		// Public schema
		publicSchema := dbo.NewSchema("public", "admin", nil)

		// Users table
		usersTable := dbo.NewTable("users", nil)
		userIDCol := dbo.NewColumn("id", "serial", false)
		userNameCol := dbo.NewColumn("name", "varchar", false)
		userEmailCol := dbo.NewColumn("email", "varchar", false)
		usersTable.AddColumn(userIDCol)
		usersTable.AddColumn(userNameCol)
		usersTable.AddColumn(userEmailCol)
		usersTable.SetPrimaryKey(dbo.NewPrimaryKey("users_pkey", usersTable, []*dbo.Column{userIDCol}))
		usersTable.AddIndex(dbo.NewIndex("idx_users_email", usersTable, []*dbo.Column{userEmailCol}, true))
		publicSchema.AddTable(usersTable)

		// Orders table with FK to users
		ordersTable := dbo.NewTable("orders", nil)
		orderIDCol := dbo.NewColumn("id", "serial", false)
		orderUserCol := dbo.NewColumn("user_id", "integer", false)
		orderTotalCol := dbo.NewColumn("total", "decimal", false)
		ordersTable.AddColumn(orderIDCol)
		ordersTable.AddColumn(orderUserCol)
		ordersTable.AddColumn(orderTotalCol)
		ordersTable.SetPrimaryKey(dbo.NewPrimaryKey("orders_pkey", ordersTable, []*dbo.Column{orderIDCol}))

		fk := dbo.NewForeignKey("fk_orders_user", "users")
		fk.SetTable(ordersTable)
		fk.AddColumn(orderUserCol)
		fk.AddReferencedColumn(userIDCol)
		fk.SetOnDelete(dbo.ActionCascade)
		ordersTable.AddForeignKey(fk)
		publicSchema.AddTable(ordersTable)

		// Add view
		view := dbo.NewView("order_summary", "SELECT u.name, COUNT(o.id) FROM users u JOIN orders o ON u.id = o.user_id GROUP BY u.name")
		publicSchema.AddView(view)

		// Add function
		fn := dbo.NewFunction("get_user_orders", "SELECT * FROM orders WHERE user_id = $1")
		fn.SetReturnType("SETOF orders")
		fn.SetLanguage("sql")
		fn.AddParameter(dbo.NewFunctionParameter("p_user_id", "integer", dbo.ParameterModeIn))
		publicSchema.AddFunction(fn)

		// Add sequence
		seq := dbo.NewSequence("orders_id_seq", 1, 1)
		publicSchema.AddSequence(seq)

		db.AddSchema(publicSchema)

		// Write report
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "report.json")

		writer := &JSONReportWriter{}
		err := writer.WriteInventoryReport(filePath, db)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Read and validate
		content, err := os.ReadFile(filePath)
		if err != nil {
			t.Fatalf("failed to read file: %v", err)
		}

		var result map[string]interface{}
		if err := json.Unmarshal(content, &result); err != nil {
			t.Fatalf("invalid JSON: %v", err)
		}

		// Validate structure
		if result["name"] != "ecommerce" {
			t.Errorf("expected name 'ecommerce', got %v", result["name"])
		}

		schemas := result["schemas"].([]interface{})
		if len(schemas) != 1 {
			t.Errorf("expected 1 schema, got %d", len(schemas))
		}

		schema := schemas[0].(map[string]interface{})
		if schema["name"] != "public" {
			t.Errorf("expected schema name 'public', got %v", schema["name"])
		}

		tables := schema["tables"].([]interface{})
		if len(tables) != 2 {
			t.Errorf("expected 2 tables, got %d", len(tables))
		}

		views := schema["views"].([]interface{})
		if len(views) != 1 {
			t.Errorf("expected 1 view, got %d", len(views))
		}

		functions := schema["functions"].([]interface{})
		if len(functions) != 1 {
			t.Errorf("expected 1 function, got %d", len(functions))
		}

		sequences := schema["sequences"].([]interface{})
		if len(sequences) != 1 {
			t.Errorf("expected 1 sequence, got %d", len(sequences))
		}
	})
}
