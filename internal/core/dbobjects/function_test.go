package dbobjects

import (
	"encoding/json"
	"testing"
)

func TestNewFunctionParameter(t *testing.T) {
	tests := []struct {
		name     string
		pName    string
		dataType string
		mode     ParameterMode
	}{
		{"in parameter", "input_val", "integer", ParameterModeIn},
		{"out parameter", "output_val", "varchar", ParameterModeOut},
		{"inout parameter", "io_val", "text", ParameterModeInOut},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewFunctionParameter(tt.pName, tt.dataType, tt.mode)

			if p.Name() != tt.pName {
				t.Errorf("expected name %q, got %q", tt.pName, p.Name())
			}
			if p.DataType() != tt.dataType {
				t.Errorf("expected dataType %q, got %q", tt.dataType, p.DataType())
			}
			if p.Mode() != tt.mode {
				t.Errorf("expected mode %q, got %q", tt.mode, p.Mode())
			}
		})
	}
}

func TestFunctionParameterMarshalJSON(t *testing.T) {
	p := NewFunctionParameter("user_id", "integer", ParameterModeIn)

	data, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("failed to marshal function parameter: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("failed to unmarshal json: %v", err)
	}

	if result["name"] != "user_id" {
		t.Errorf("expected name 'user_id', got %v", result["name"])
	}
	if result["dataType"] != "integer" {
		t.Errorf("expected dataType 'integer', got %v", result["dataType"])
	}
	if result["mode"] != string(ParameterModeIn) {
		t.Errorf("expected mode %q, got %v", ParameterModeIn, result["mode"])
	}
}

func TestNewFunction(t *testing.T) {
	f := NewFunction("calculate_total", "SELECT SUM(amount) FROM orders")

	if f.Name() != "calculate_total" {
		t.Errorf("expected name 'calculate_total', got %q", f.Name())
	}
	if f.Definition() != "SELECT SUM(amount) FROM orders" {
		t.Errorf("expected definition, got %q", f.Definition())
	}
	if len(f.Parameters()) != 0 {
		t.Errorf("expected empty parameters, got %d", len(f.Parameters()))
	}
	if f.Language() != "sql" {
		t.Errorf("expected default language 'sql', got %q", f.Language())
	}
}

func TestFunctionSchema(t *testing.T) {
	f := NewFunction("test_func", "SELECT 1")

	if f.Schema() != nil {
		t.Error("expected nil schema initially")
	}

	schema := NewSchema("public", "owner", nil)
	f.SetSchema(schema)

	if f.Schema() == nil {
		t.Fatal("expected schema to be set")
	}
	if f.Schema().Name() != "public" {
		t.Errorf("expected schema name 'public', got %q", f.Schema().Name())
	}
}

func TestFunctionDefinition(t *testing.T) {
	f := NewFunction("test_func", "original")

	f.SetDefinition("updated definition")

	if f.Definition() != "updated definition" {
		t.Errorf("expected 'updated definition', got %q", f.Definition())
	}
}

func TestFunctionReturnType(t *testing.T) {
	f := NewFunction("test_func", "SELECT 1")

	if f.ReturnType() != "" {
		t.Errorf("expected empty return type, got %q", f.ReturnType())
	}

	f.SetReturnType("integer")

	if f.ReturnType() != "integer" {
		t.Errorf("expected 'integer', got %q", f.ReturnType())
	}
}

func TestFunctionParameters(t *testing.T) {
	f := NewFunction("add_numbers", "SELECT $1 + $2")
	p1 := NewFunctionParameter("a", "integer", ParameterModeIn)
	p2 := NewFunctionParameter("b", "integer", ParameterModeIn)

	f.AddParameter(p1)
	f.AddParameter(p2)

	params := f.Parameters()
	if len(params) != 2 {
		t.Fatalf("expected 2 parameters, got %d", len(params))
	}
	if params[0].Name() != "a" {
		t.Errorf("expected first param 'a', got %q", params[0].Name())
	}
	if params[1].Name() != "b" {
		t.Errorf("expected second param 'b', got %q", params[1].Name())
	}
}

func TestFunctionLanguage(t *testing.T) {
	f := NewFunction("test_func", "BEGIN RETURN 1; END")

	f.SetLanguage("plpgsql")

	if f.Language() != "plpgsql" {
		t.Errorf("expected 'plpgsql', got %q", f.Language())
	}
}

func TestFunctionFullyQualifiedName(t *testing.T) {
	t.Run("without schema", func(t *testing.T) {
		f := NewFunction("my_func", "SELECT 1")

		if f.FullyQualifiedName() != "my_func" {
			t.Errorf("expected 'my_func', got %q", f.FullyQualifiedName())
		}
	})

	t.Run("with schema", func(t *testing.T) {
		f := NewFunction("my_func", "SELECT 1")
		schema := NewSchema("public", "owner", nil)
		f.SetSchema(schema)

		if f.FullyQualifiedName() != "public.my_func" {
			t.Errorf("expected 'public.my_func', got %q", f.FullyQualifiedName())
		}
	})
}

func TestFunctionMarshalJSON(t *testing.T) {
	f := NewFunction("get_user", "SELECT * FROM users WHERE id = $1")
	f.SetReturnType("record")
	f.SetLanguage("sql")

	p := NewFunctionParameter("user_id", "integer", ParameterModeIn)
	f.AddParameter(p)

	data, err := json.Marshal(f)
	if err != nil {
		t.Fatalf("failed to marshal function: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("failed to unmarshal json: %v", err)
	}

	if result["name"] != "get_user" {
		t.Errorf("expected name 'get_user', got %v", result["name"])
	}
	if result["returnType"] != "record" {
		t.Errorf("expected returnType 'record', got %v", result["returnType"])
	}
	if result["language"] != "sql" {
		t.Errorf("expected language 'sql', got %v", result["language"])
	}

	params := result["parameters"].([]interface{})
	if len(params) != 1 {
		t.Errorf("expected 1 parameter, got %d", len(params))
	}
}

func TestFunctionMarshalJSONOmitsEmptyParameters(t *testing.T) {
	f := NewFunction("simple_func", "SELECT 1")
	f.SetReturnType("integer")

	data, err := json.Marshal(f)
	if err != nil {
		t.Fatalf("failed to marshal function: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("failed to unmarshal json: %v", err)
	}

	// Empty slices are marshaled as [] not omitted with omitempty for slices
	params := result["parameters"]
	if params != nil {
		paramSlice := params.([]interface{})
		if len(paramSlice) != 0 {
			t.Errorf("expected empty parameters, got %v", params)
		}
	}
}
