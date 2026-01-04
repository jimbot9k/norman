package dbobjects

import "encoding/json"

type Database struct {
	name    string
	schemas map[string]*Schema
}

func (d *Database) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name    string             `json:"name"`
		Schemas map[string]*Schema `json:"schemas"`
	}{
		Name:    d.name,
		Schemas: d.schemas,
	})
}

func NewDatabase(name string, schemas map[string]*Schema) *Database {
	if schemas == nil {
		schemas = make(map[string]*Schema)
	}
	return &Database{
		name:    name,
		schemas: schemas,
	}
}

func (d *Database) Name() string {
	return d.name
}

func (d *Database) Schemas() map[string]*Schema {
	return d.schemas
}

func (d *Database) AddSchema(schema *Schema) {
	schema.SetDatabase(d)
	d.schemas[schema.Name()] = schema
}
