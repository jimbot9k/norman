package dbobjects

import "encoding/json"

type View struct {
	name       string
	schema     *Schema
	definition string
	columns    []*Column
}

func (v *View) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name       string    `json:"name"`
		Definition string    `json:"definition"`
		Columns    []*Column `json:"columns,omitempty"`
	}{
		Name:       v.name,
		Definition: v.definition,
		Columns:    v.columns,
	})
}

func NewView(name string, definition string) *View {
	return &View{
		name:       name,
		definition: definition,
		columns:    []*Column{},
	}
}

func (v *View) Name() string {
	return v.name
}

func (v *View) Schema() *Schema {
	return v.schema
}

func (v *View) SetSchema(schema *Schema) {
	v.schema = schema
}

func (v *View) Definition() string {
	return v.definition
}

func (v *View) SetDefinition(definition string) {
	v.definition = definition
}

func (v *View) Columns() []*Column {
	return v.columns
}

func (v *View) AddColumn(column *Column) {
	v.columns = append(v.columns, column)
}

// FullyQualifiedName returns schema.view format if schema is set
func (v *View) FullyQualifiedName() string {
	if v.schema != nil {
		return v.schema.Name() + "." + v.name
	}
	return v.name
}
