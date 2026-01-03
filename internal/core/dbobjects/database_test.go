package dbobjects

import (
	"encoding/json"
	"testing"
)

func TestNewDatabase(t *testing.T) {
	db := NewDatabase("testdb", nil)

	if db.Name() != "testdb" {
		t.Errorf("expected name 'testdb', got %q", db.Name())
	}
	if db.Schemas() == nil {
		t.Error("expected schemas map to be initialized")
	}
	if len(db.Schemas()) != 0 {
		t.Errorf("expected empty schemas, got %d", len(db.Schemas()))
	}
}

func TestNewDatabaseWithSchemas(t *testing.T) {
	schema := NewSchema("public", "owner", nil)
	schemas := map[string]*Schema{
		"public": schema,
	}

	db := NewDatabase("testdb", schemas)

	if len(db.Schemas()) != 1 {
		t.Errorf("expected 1 schema, got %d", len(db.Schemas()))
	}
	if db.Schemas()["public"] != schema {
		t.Error("expected public schema to be set")
	}
}

func TestDatabaseAddSchema(t *testing.T) {
	db := NewDatabase("testdb", nil)
	schema1 := NewSchema("public", "owner1", nil)
	schema2 := NewSchema("private", "owner2", nil)

	db.AddSchema(schema1)
	db.AddSchema(schema2)

	schemas := db.Schemas()
	if len(schemas) != 2 {
		t.Fatalf("expected 2 schemas, got %d", len(schemas))
	}
	if schemas["public"] == nil {
		t.Error("expected public schema to exist")
	}
	if schemas["private"] == nil {
		t.Error("expected private schema to exist")
	}

	// Verify schema has database reference set
	if schema1.Database() != db {
		t.Error("expected schema to have database reference")
	}
}

func TestDatabaseMarshalJSON(t *testing.T) {
	db := NewDatabase("production", nil)
	schema := NewSchema("public", "admin", nil)
	db.AddSchema(schema)

	data, err := json.Marshal(db)
	if err != nil {
		t.Fatalf("failed to marshal database: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("failed to unmarshal json: %v", err)
	}

	if result["name"] != "production" {
		t.Errorf("expected name 'production', got %v", result["name"])
	}

	schemas := result["schemas"].(map[string]interface{})
	if len(schemas) != 1 {
		t.Errorf("expected 1 schema, got %d", len(schemas))
	}
	if schemas["public"] == nil {
		t.Error("expected public schema in JSON")
	}
}

func TestDatabaseMarshalJSONEmptySchemas(t *testing.T) {
	db := NewDatabase("empty_db", nil)

	data, err := json.Marshal(db)
	if err != nil {
		t.Fatalf("failed to marshal database: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("failed to unmarshal json: %v", err)
	}

	schemas := result["schemas"].(map[string]interface{})
	if len(schemas) != 0 {
		t.Errorf("expected empty schemas, got %d", len(schemas))
	}
}
