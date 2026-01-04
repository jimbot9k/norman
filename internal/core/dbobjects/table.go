package dbobjects

import "encoding/json"

type Table struct {
	name        string
	schema      *Schema
	columns     map[string]*Column
	primaryKey  *PrimaryKey
	foreignKeys []*ForeignKey
	indexes     []*Index
	constraints []*Constraint
	triggers    []*Trigger
}

func (t *Table) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name        string             `json:"name"`
		Columns     map[string]*Column `json:"columns"`
		PrimaryKey  *PrimaryKey        `json:"primaryKey,omitempty"`
		ForeignKeys []*ForeignKey      `json:"foreignKeys,omitempty"`
		Indexes     []*Index           `json:"indexes,omitempty"`
		Constraints []*Constraint      `json:"constraints,omitempty"`
		Triggers    []*Trigger         `json:"triggers,omitempty"`
	}{
		Name:        t.name,
		Columns:     t.columns,
		PrimaryKey:  t.primaryKey,
		ForeignKeys: t.foreignKeys,
		Indexes:     t.indexes,
		Constraints: t.constraints,
		Triggers:    t.triggers,
	})
}

func NewTable(name string, columns map[string]*Column) *Table {
	if columns == nil {
		columns = make(map[string]*Column)
	}
	return &Table{
		name:        name,
		columns:     columns,
		foreignKeys: []*ForeignKey{},
		indexes:     []*Index{},
		constraints: []*Constraint{},
		triggers:    []*Trigger{},
	}
}

func (t *Table) Name() string {
	return t.name
}

func (t *Table) Schema() *Schema {
	return t.schema
}

func (t *Table) SetSchema(schema *Schema) {
	t.schema = schema
}

func (t *Table) Columns() map[string]*Column {
	return t.columns
}

func (t *Table) AddColumn(column *Column) {
	column.SetTable(t)
	t.columns[column.Name()] = column
}

func (t *Table) PrimaryKey() *PrimaryKey {
	return t.primaryKey
}

func (t *Table) SetPrimaryKey(pk *PrimaryKey) {
	t.primaryKey = pk
}

func (t *Table) ForeignKeys() []*ForeignKey {
	return t.foreignKeys
}

func (t *Table) AddForeignKey(fk *ForeignKey) {
	fk.SetTable(t)
	t.foreignKeys = append(t.foreignKeys, fk)
}

func (t *Table) Indexes() []*Index {
	return t.indexes
}

func (t *Table) AddIndex(index *Index) {
	index.SetTable(t)
	t.indexes = append(t.indexes, index)
}

func (t *Table) Constraints() []*Constraint {
	return t.constraints
}

func (t *Table) AddConstraint(constraint *Constraint) {
	constraint.SetTable(t)
	t.constraints = append(t.constraints, constraint)
}

func (t *Table) Triggers() []*Trigger {
	return t.triggers
}

func (t *Table) AddTrigger(trigger *Trigger) {
	trigger.SetTable(t)
	t.triggers = append(t.triggers, trigger)
}

// FullyQualifiedName returns schema.table format if schema is set
func (t *Table) FullyQualifiedName() string {
	if t.schema != nil {
		return t.schema.Name() + "." + t.name
	}
	return t.name
}
