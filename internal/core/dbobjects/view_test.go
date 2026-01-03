package dbobjects

import (
	"encoding/json"
	"testing"
)

func TestNewView(t *testing.T) {
	v := NewView("active_users", "SELECT * FROM users WHERE active = true")

	if v.Name() != "active_users" {
		t.Errorf("expected name 'active_users', got %q", v.Name())
	}
	if v.Definition() != "SELECT * FROM users WHERE active = true" {
		t.Errorf("expected definition, got %q", v.Definition())
	}
	if len(v.Columns()) != 0 {
		t.Errorf("expected empty columns, got %d", len(v.Columns()))
	}
}

func TestViewSchema(t *testing.T) {
	v := NewView("test_view", "SELECT 1")

	if v.Schema() != nil {
		t.Error("expected nil schema initially")
	}

	schema := NewSchema("public", "owner", nil)
	v.SetSchema(schema)

	if v.Schema() == nil {
		t.Fatal("expected schema to be set")
	}
	if v.Schema().Name() != "public" {
		t.Errorf("expected schema name 'public', got %q", v.Schema().Name())
	}
}

func TestViewDefinition(t *testing.T) {
	v := NewView("test_view", "original")

	v.SetDefinition("SELECT * FROM new_table")

	if v.Definition() != "SELECT * FROM new_table" {
		t.Errorf("expected updated definition, got %q", v.Definition())
	}
}

func TestViewColumns(t *testing.T) {
	v := NewView("user_summary", "SELECT id, name FROM users")
	col1 := NewColumn("id", "integer", false)
	col2 := NewColumn("name", "varchar", false)

	v.AddColumn(col1)
	v.AddColumn(col2)

	columns := v.Columns()
	if len(columns) != 2 {
		t.Fatalf("expected 2 columns, got %d", len(columns))
	}
	if columns[0].Name() != "id" {
		t.Errorf("expected first column 'id', got %q", columns[0].Name())
	}
	if columns[1].Name() != "name" {
		t.Errorf("expected second column 'name', got %q", columns[1].Name())
	}
}

func TestViewFullyQualifiedName(t *testing.T) {
	t.Run("without schema", func(t *testing.T) {
		v := NewView("my_view", "SELECT 1")

		if v.FullyQualifiedName() != "my_view" {
			t.Errorf("expected 'my_view', got %q", v.FullyQualifiedName())
		}
	})

	t.Run("with schema", func(t *testing.T) {
		v := NewView("my_view", "SELECT 1")
		schema := NewSchema("public", "owner", nil)
		v.SetSchema(schema)

		if v.FullyQualifiedName() != "public.my_view" {
			t.Errorf("expected 'public.my_view', got %q", v.FullyQualifiedName())
		}
	})
}

func TestViewMarshalJSON(t *testing.T) {
	v := NewView("employee_summary", "SELECT id, name FROM employees")
	col1 := NewColumn("id", "integer", false)
	col2 := NewColumn("name", "varchar", false)
	v.AddColumn(col1)
	v.AddColumn(col2)

	data, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("failed to marshal view: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("failed to unmarshal json: %v", err)
	}

	if result["name"] != "employee_summary" {
		t.Errorf("expected name 'employee_summary', got %v", result["name"])
	}
	if result["definition"] != "SELECT id, name FROM employees" {
		t.Errorf("expected definition, got %v", result["definition"])
	}

	columns := result["columns"].([]interface{})
	if len(columns) != 2 {
		t.Errorf("expected 2 columns, got %d", len(columns))
	}
}

func TestViewMarshalJSONOmitsEmptyColumns(t *testing.T) {
	v := NewView("simple_view", "SELECT 1")

	data, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("failed to marshal view: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("failed to unmarshal json: %v", err)
	}

	// Empty slices are marshaled as [] not omitted
	cols := result["columns"]
	if cols != nil {
		colSlice := cols.([]interface{})
		if len(colSlice) != 0 {
			t.Errorf("expected empty columns, got %v", cols)
		}
	}
}
