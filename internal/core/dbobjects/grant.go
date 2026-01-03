package dbobjects

type Grant struct {
	name       string
	definition string
}

func NewGrant(name string, definition string) *Grant {
	return &Grant{
		name:       name,
		definition: definition,
	}
}

func (g *Grant) Name() string {
	return g.name
}

func (g *Grant) Definition() string {
	return g.definition
}
