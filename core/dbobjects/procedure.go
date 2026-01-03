package dbobjects

import "encoding/json"

type Procedure struct {
	name       string
	schema     *Schema
	definition string
	parameters []*FunctionParameter
	language   string
}

func (p *Procedure) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name       string               `json:"name"`
		Definition string               `json:"definition"`
		Parameters []*FunctionParameter `json:"parameters,omitempty"`
		Language   string               `json:"language"`
	}{
		Name:       p.name,
		Definition: p.definition,
		Parameters: p.parameters,
		Language:   p.language,
	})
}

func NewProcedure(name string, definition string) *Procedure {
	return &Procedure{
		name:       name,
		definition: definition,
		parameters: []*FunctionParameter{},
		language:   "sql",
	}
}

func (p *Procedure) Name() string {
	return p.name
}

func (p *Procedure) Schema() *Schema {
	return p.schema
}

func (p *Procedure) SetSchema(schema *Schema) {
	p.schema = schema
}

func (p *Procedure) Definition() string {
	return p.definition
}

func (p *Procedure) SetDefinition(definition string) {
	p.definition = definition
}

func (p *Procedure) Parameters() []*FunctionParameter {
	return p.parameters
}

func (p *Procedure) AddParameter(param *FunctionParameter) {
	p.parameters = append(p.parameters, param)
}

func (p *Procedure) Language() string {
	return p.language
}

func (p *Procedure) SetLanguage(language string) {
	p.language = language
}

// FullyQualifiedName returns schema.procedure format if schema is set
func (p *Procedure) FullyQualifiedName() string {
	if p.schema != nil {
		return p.schema.Name() + "." + p.name
	}
	return p.name
}
