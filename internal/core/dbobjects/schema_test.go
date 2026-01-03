package dbobjects

import (
	"encoding/json"
	"testing"
)

func TestNewSchema(t *testing.T) {
	s := NewSchema("public", "admin", nil)

	if s.Name() != "public" {
		t.Errorf("expected name 'public', got %q", s.Name())
	}
	if s.Owner() != "admin" {
		t.Errorf("expected owner 'admin', got %q", s.Owner())
	}
	if s.Tables() == nil {
		t.Error("expected tables map to be initialized")
	}
	if s.Views() == nil {
		t.Error("expected views map to be initialized")
	}
	if s.Functions() == nil {
		t.Error("expected functions map to be initialized")
	}
	if s.Procedures() == nil {
		t.Error("expected procedures map to be initialized")
	}
	if s.Sequences() == nil {
		t.Error("expected sequences map to be initialized")
	}
}

func TestNewSchemaWithTables(t *testing.T) {
	table := NewTable("users", nil)
	tables := map[string]*Table{
		"users": table,
	}

	s := NewSchema("public", "admin", tables)

	if len(s.Tables()) != 1 {
		t.Errorf("expected 1 table, got %d", len(s.Tables()))
	}
	if s.Tables()["users"] != table {
		t.Error("expected users table to be set")
	}
}

func TestSchemaDatabase(t *testing.T) {
	s := NewSchema("public", "admin", nil)

	if s.Database() != nil {
		t.Error("expected nil database initially")
	}

	db := NewDatabase("testdb", nil)
	s.SetDatabase(db)

	if s.Database() == nil {
		t.Fatal("expected database to be set")
	}
	if s.Database().Name() != "testdb" {
		t.Errorf("expected database name 'testdb', got %q", s.Database().Name())
	}
}

func TestSchemaAddTable(t *testing.T) {
	s := NewSchema("public", "admin", nil)
	table1 := NewTable("users", nil)
	table2 := NewTable("orders", nil)

	s.AddTable(table1)
	s.AddTable(table2)

	tables := s.Tables()
	if len(tables) != 2 {
		t.Fatalf("expected 2 tables, got %d", len(tables))
	}
	if tables["users"] == nil {
		t.Error("expected users table to exist")
	}
	if tables["orders"] == nil {
		t.Error("expected orders table to exist")
	}

	// Verify table has schema reference set
	if table1.Schema() != s {
		t.Error("expected table to have schema reference")
	}
}

func TestSchemaAddView(t *testing.T) {
	s := NewSchema("public", "admin", nil)
	view1 := NewView("active_users", "SELECT * FROM users WHERE active")
	view2 := NewView("recent_orders", "SELECT * FROM orders WHERE created_at > now() - interval '1 day'")

	s.AddView(view1)
	s.AddView(view2)

	views := s.Views()
	if len(views) != 2 {
		t.Fatalf("expected 2 views, got %d", len(views))
	}
	if views["active_users"] == nil {
		t.Error("expected active_users view to exist")
	}

	// Verify view has schema reference set
	if view1.Schema() != s {
		t.Error("expected view to have schema reference")
	}
}

func TestSchemaAddFunction(t *testing.T) {
	s := NewSchema("public", "admin", nil)
	fn1 := NewFunction("get_user", "SELECT * FROM users WHERE id = $1")
	fn2 := NewFunction("calculate_total", "SELECT SUM(amount) FROM orders")

	s.AddFunction(fn1)
	s.AddFunction(fn2)

	functions := s.Functions()
	if len(functions) != 2 {
		t.Fatalf("expected 2 functions, got %d", len(functions))
	}
	if functions["get_user"] == nil {
		t.Error("expected get_user function to exist")
	}

	// Verify function has schema reference set
	if fn1.Schema() != s {
		t.Error("expected function to have schema reference")
	}
}

func TestSchemaAddProcedure(t *testing.T) {
	s := NewSchema("public", "admin", nil)
	proc1 := NewProcedure("refresh_cache", "TRUNCATE cache; INSERT INTO cache ...")
	proc2 := NewProcedure("cleanup_old_data", "DELETE FROM logs WHERE created_at < ...")

	s.AddProcedure(proc1)
	s.AddProcedure(proc2)

	procedures := s.Procedures()
	if len(procedures) != 2 {
		t.Fatalf("expected 2 procedures, got %d", len(procedures))
	}
	if procedures["refresh_cache"] == nil {
		t.Error("expected refresh_cache procedure to exist")
	}

	// Verify procedure has schema reference set
	if proc1.Schema() != s {
		t.Error("expected procedure to have schema reference")
	}
}

func TestSchemaAddSequence(t *testing.T) {
	s := NewSchema("public", "admin", nil)
	seq1 := NewSequence("user_id_seq", 1, 1)
	seq2 := NewSequence("order_id_seq", 1000, 1)

	s.AddSequence(seq1)
	s.AddSequence(seq2)

	sequences := s.Sequences()
	if len(sequences) != 2 {
		t.Fatalf("expected 2 sequences, got %d", len(sequences))
	}
	if sequences["user_id_seq"] == nil {
		t.Error("expected user_id_seq sequence to exist")
	}

	// Verify sequence has schema reference set
	if seq1.Schema() != s {
		t.Error("expected sequence to have schema reference")
	}
}

func TestSchemaFullyQualifiedName(t *testing.T) {
	t.Run("without database", func(t *testing.T) {
		s := NewSchema("public", "admin", nil)

		if s.FullyQualifiedName() != "public" {
			t.Errorf("expected 'public', got %q", s.FullyQualifiedName())
		}
	})

	t.Run("with database", func(t *testing.T) {
		s := NewSchema("public", "admin", nil)
		db := NewDatabase("testdb", nil)
		s.SetDatabase(db)

		if s.FullyQualifiedName() != "testdb.public" {
			t.Errorf("expected 'testdb.public', got %q", s.FullyQualifiedName())
		}
	})
}

func TestSchemaMarshalJSON(t *testing.T) {
	s := NewSchema("public", "admin", nil)

	table := NewTable("users", nil)
	s.AddTable(table)

	view := NewView("active_users", "SELECT * FROM users")
	s.AddView(view)

	fn := NewFunction("get_user", "SELECT * FROM users WHERE id = $1")
	s.AddFunction(fn)

	proc := NewProcedure("refresh", "TRUNCATE ...")
	s.AddProcedure(proc)

	seq := NewSequence("user_id_seq", 1, 1)
	s.AddSequence(seq)

	data, err := json.Marshal(s)
	if err != nil {
		t.Fatalf("failed to marshal schema: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("failed to unmarshal json: %v", err)
	}

	if result["name"] != "public" {
		t.Errorf("expected name 'public', got %v", result["name"])
	}
	if result["owner"] != "admin" {
		t.Errorf("expected owner 'admin', got %v", result["owner"])
	}

	tables := result["tables"].(map[string]interface{})
	if len(tables) != 1 {
		t.Errorf("expected 1 table, got %d", len(tables))
	}

	views := result["views"].(map[string]interface{})
	if len(views) != 1 {
		t.Errorf("expected 1 view, got %d", len(views))
	}

	functions := result["functions"].(map[string]interface{})
	if len(functions) != 1 {
		t.Errorf("expected 1 function, got %d", len(functions))
	}

	procedures := result["procedures"].(map[string]interface{})
	if len(procedures) != 1 {
		t.Errorf("expected 1 procedure, got %d", len(procedures))
	}

	sequences := result["sequences"].(map[string]interface{})
	if len(sequences) != 1 {
		t.Errorf("expected 1 sequence, got %d", len(sequences))
	}
}

func TestSchemaMarshalJSONEmpty(t *testing.T) {
	s := NewSchema("empty_schema", "admin", nil)

	data, err := json.Marshal(s)
	if err != nil {
		t.Fatalf("failed to marshal schema: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("failed to unmarshal json: %v", err)
	}

	tables := result["tables"].(map[string]interface{})
	if len(tables) != 0 {
		t.Errorf("expected empty tables, got %d", len(tables))
	}
}
