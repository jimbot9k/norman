package dbobjects

type CompositeType struct {
	name       string
	definition string
}

func NewCompositeType(name string, definition string) *CompositeType {
	return &CompositeType{
		name:       name,
		definition: definition,
	}
}

func (ct *CompositeType) Name() string {
	return ct.name
}

func (ct *CompositeType) Definition() string {
	return ct.definition
}