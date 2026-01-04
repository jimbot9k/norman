package dbobjects

import "encoding/json"

type ConstraintType string

const (
	ConstraintTypeCheck   ConstraintType = "CHECK"
	ConstraintTypeUnique  ConstraintType = "UNIQUE"
	ConstraintTypeNotNull ConstraintType = "NOT NULL"
)

type Constraint struct {
	name            string
	constraintType  ConstraintType
	table           *Table
	columns         []*Column
	checkExpression string
}

func (c *Constraint) MarshalJSON() ([]byte, error) {
	columnNames := make([]string, len(c.columns))
	for i, col := range c.columns {
		columnNames[i] = col.Name()
	}
	return json.Marshal(struct {
		Name            string         `json:"name"`
		Type            ConstraintType `json:"type"`
		Columns         []string       `json:"columns"`
		CheckExpression string         `json:"checkExpression,omitempty"`
	}{
		Name:            c.name,
		Type:            c.constraintType,
		Columns:         columnNames,
		CheckExpression: c.checkExpression,
	})
}

func NewConstraint(name string, constraintType ConstraintType) *Constraint {
	return &Constraint{
		name:           name,
		constraintType: constraintType,
		columns:        []*Column{},
	}
}

func (c *Constraint) Name() string {
	return c.name
}

func (c *Constraint) Type() ConstraintType {
	return c.constraintType
}

func (c *Constraint) SetType(constraintType ConstraintType) {
	c.constraintType = constraintType
}

func (c *Constraint) Table() *Table {
	return c.table
}

func (c *Constraint) SetTable(table *Table) {
	c.table = table
}

func (c *Constraint) Columns() []*Column {
	return c.columns
}

func (c *Constraint) AddColumn(column *Column) {
	c.columns = append(c.columns, column)
}

func (c *Constraint) CheckExpression() string {
	return c.checkExpression
}

func (c *Constraint) SetCheckExpression(expression string) {
	c.checkExpression = expression
}
