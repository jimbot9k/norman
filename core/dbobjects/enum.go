package dbobjects

type Enum struct {
	name   string
	values []string
}

func NewEnum(name string, values []string) *Enum {
	return &Enum{
		name:   name,
		values: values,
	}
}

func (e *Enum) Name() string {
	return e.name
}

func (e *Enum) Values() []string {
	return e.values
}