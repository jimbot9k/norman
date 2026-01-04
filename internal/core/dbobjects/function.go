package dbobjects

import "encoding/json"

type ParameterMode string

const (
	ParameterModeIn    ParameterMode = "IN"
	ParameterModeOut   ParameterMode = "OUT"
	ParameterModeInOut ParameterMode = "INOUT"
)

type FunctionParameter struct {
	name     string
	dataType string
	mode     ParameterMode
}

func (p *FunctionParameter) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name     string        `json:"name"`
		DataType string        `json:"dataType"`
		Mode     ParameterMode `json:"mode"`
	}{
		Name:     p.name,
		DataType: p.dataType,
		Mode:     p.mode,
	})
}

func NewFunctionParameter(name string, dataType string, mode ParameterMode) *FunctionParameter {
	return &FunctionParameter{
		name:     name,
		dataType: dataType,
		mode:     mode,
	}
}

func (p *FunctionParameter) Name() string {
	return p.name
}

func (p *FunctionParameter) DataType() string {
	return p.dataType
}

func (p *FunctionParameter) Mode() ParameterMode {
	return p.mode
}

type Function struct {
	name       string
	schema     *Schema
	definition string
	returnType string
	parameters []*FunctionParameter
	language   string
}

func (f *Function) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name       string               `json:"name"`
		Definition string               `json:"definition"`
		ReturnType string               `json:"returnType"`
		Parameters []*FunctionParameter `json:"parameters,omitempty"`
		Language   string               `json:"language"`
	}{
		Name:       f.name,
		Definition: f.definition,
		ReturnType: f.returnType,
		Parameters: f.parameters,
		Language:   f.language,
	})
}

func NewFunction(name string, definition string) *Function {
	return &Function{
		name:       name,
		definition: definition,
		parameters: []*FunctionParameter{},
		language:   "sql",
	}
}

func (f *Function) Name() string {
	return f.name
}

func (f *Function) Schema() *Schema {
	return f.schema
}

func (f *Function) SetSchema(schema *Schema) {
	f.schema = schema
}

func (f *Function) Definition() string {
	return f.definition
}

func (f *Function) SetDefinition(definition string) {
	f.definition = definition
}

func (f *Function) ReturnType() string {
	return f.returnType
}

func (f *Function) SetReturnType(returnType string) {
	f.returnType = returnType
}

func (f *Function) Parameters() []*FunctionParameter {
	return f.parameters
}

func (f *Function) AddParameter(param *FunctionParameter) {
	f.parameters = append(f.parameters, param)
}

func (f *Function) Language() string {
	return f.language
}

func (f *Function) SetLanguage(language string) {
	f.language = language
}

// FullyQualifiedName returns schema.function format if schema is set
func (f *Function) FullyQualifiedName() string {
	if f.schema != nil {
		return f.schema.Name() + "." + f.name
	}
	return f.name
}
