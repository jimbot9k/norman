package dbobjects

import (
	"encoding/json"
	"testing"
)

func TestNewProcedure(t *testing.T) {
	p := NewProcedure("update_user", "UPDATE users SET name = $1 WHERE id = $2")

	if p.Name() != "update_user" {
		t.Errorf("expected name 'update_user', got %q", p.Name())
	}
	if p.Definition() != "UPDATE users SET name = $1 WHERE id = $2" {
		t.Errorf("expected definition, got %q", p.Definition())
	}
	if len(p.Parameters()) != 0 {
		t.Errorf("expected empty parameters, got %d", len(p.Parameters()))
	}
	if p.Language() != "sql" {
		t.Errorf("expected default language 'sql', got %q", p.Language())
	}
}

func TestProcedureSchema(t *testing.T) {
	p := NewProcedure("test_proc", "SELECT 1")

	if p.Schema() != nil {
		t.Error("expected nil schema initially")
	}

	schema := NewSchema("public", "owner", nil)
	p.SetSchema(schema)

	if p.Schema() == nil {
		t.Fatal("expected schema to be set")
	}
	if p.Schema().Name() != "public" {
		t.Errorf("expected schema name 'public', got %q", p.Schema().Name())
	}
}

func TestProcedureDefinition(t *testing.T) {
	p := NewProcedure("test_proc", "original")

	p.SetDefinition("updated definition")

	if p.Definition() != "updated definition" {
		t.Errorf("expected 'updated definition', got %q", p.Definition())
	}
}

func TestProcedureParameters(t *testing.T) {
	p := NewProcedure("update_records", "UPDATE table SET col = $1")
	param1 := NewFunctionParameter("new_value", "text", ParameterModeIn)
	param2 := NewFunctionParameter("record_id", "integer", ParameterModeIn)

	p.AddParameter(param1)
	p.AddParameter(param2)

	params := p.Parameters()
	if len(params) != 2 {
		t.Fatalf("expected 2 parameters, got %d", len(params))
	}
	if params[0].Name() != "new_value" {
		t.Errorf("expected first param 'new_value', got %q", params[0].Name())
	}
	if params[1].Name() != "record_id" {
		t.Errorf("expected second param 'record_id', got %q", params[1].Name())
	}
}

func TestProcedureLanguage(t *testing.T) {
	p := NewProcedure("test_proc", "BEGIN ... END")

	p.SetLanguage("plpgsql")

	if p.Language() != "plpgsql" {
		t.Errorf("expected 'plpgsql', got %q", p.Language())
	}
}

func TestProcedureFullyQualifiedName(t *testing.T) {
	t.Run("without schema", func(t *testing.T) {
		p := NewProcedure("my_proc", "SELECT 1")

		if p.FullyQualifiedName() != "my_proc" {
			t.Errorf("expected 'my_proc', got %q", p.FullyQualifiedName())
		}
	})

	t.Run("with schema", func(t *testing.T) {
		p := NewProcedure("my_proc", "SELECT 1")
		schema := NewSchema("public", "owner", nil)
		p.SetSchema(schema)

		if p.FullyQualifiedName() != "public.my_proc" {
			t.Errorf("expected 'public.my_proc', got %q", p.FullyQualifiedName())
		}
	})
}

func TestProcedureMarshalJSON(t *testing.T) {
	p := NewProcedure("refresh_data", "TRUNCATE table; INSERT INTO table SELECT ...")
	p.SetLanguage("sql")

	param := NewFunctionParameter("batch_size", "integer", ParameterModeIn)
	p.AddParameter(param)

	data, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("failed to marshal procedure: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("failed to unmarshal json: %v", err)
	}

	if result["name"] != "refresh_data" {
		t.Errorf("expected name 'refresh_data', got %v", result["name"])
	}
	if result["language"] != "sql" {
		t.Errorf("expected language 'sql', got %v", result["language"])
	}

	params := result["parameters"].([]interface{})
	if len(params) != 1 {
		t.Errorf("expected 1 parameter, got %d", len(params))
	}
}

func TestProcedureMarshalJSONEmptyParameters(t *testing.T) {
	p := NewProcedure("simple_proc", "SELECT 1")

	data, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("failed to marshal procedure: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("failed to unmarshal json: %v", err)
	}

	// Empty slices are marshaled as [] not omitted
	params := result["parameters"]
	if params != nil {
		paramSlice := params.([]interface{})
		if len(paramSlice) != 0 {
			t.Errorf("expected empty parameters, got %v", params)
		}
	}
}
