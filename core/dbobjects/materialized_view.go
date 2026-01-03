package dbobjects

type MaterializedView struct {
	name       string
	definition string
}

func NewMaterializedView(name string, definition string) *MaterializedView {
	return &MaterializedView{
		name:       name,
		definition: definition,
	}
}

func (mv *MaterializedView) Name() string {
	return mv.name
}

func (mv *MaterializedView) Definition() string {
	return mv.definition
}