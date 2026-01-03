package dbobjects

type Domain struct {
	name       string
	definition string
}

func NewDomain(name string, definition string) *Domain {
	return &Domain{
		name:       name,
		definition: definition,
	}
}

func (d *Domain) Name() string {
	return d.name
}

func (d *Domain) Definition() string {
	return d.definition
}