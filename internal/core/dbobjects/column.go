package dbobjects

import "encoding/json"

type Column struct {
	name             string
	dataType         string
	nullable         bool
	defaultValue     *string
	ordinalPosition  int
	charMaxLength    *int
	numericPrecision *int
	numericScale     *int
	table            *Table
}

func (c *Column) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name             string  `json:"name"`
		DataType         string  `json:"dataType"`
		Nullable         bool    `json:"nullable"`
		DefaultValue     *string `json:"defaultValue,omitempty"`
		OrdinalPosition  int     `json:"ordinalPosition"`
		CharMaxLength    *int    `json:"charMaxLength,omitempty"`
		NumericPrecision *int    `json:"numericPrecision,omitempty"`
		NumericScale     *int    `json:"numericScale,omitempty"`
	}{
		Name:             c.name,
		DataType:         c.dataType,
		Nullable:         c.nullable,
		DefaultValue:     c.defaultValue,
		OrdinalPosition:  c.ordinalPosition,
		CharMaxLength:    c.charMaxLength,
		NumericPrecision: c.numericPrecision,
		NumericScale:     c.numericScale,
	})
}

func NewColumn(name string, dataType string, nullable bool) *Column {
	return &Column{
		name:     name,
		dataType: dataType,
		nullable: nullable,
	}
}

func (c *Column) Name() string {
	return c.name
}

func (c *Column) DataType() string {
	return c.dataType
}

func (c *Column) IsNullable() bool {
	return c.nullable
}

func (c *Column) DefaultValue() *string {
	return c.defaultValue
}

func (c *Column) SetDefaultValue(value string) {
	c.defaultValue = &value
}

func (c *Column) OrdinalPosition() int {
	return c.ordinalPosition
}

func (c *Column) SetOrdinalPosition(pos int) {
	c.ordinalPosition = pos
}

func (c *Column) CharMaxLength() *int {
	return c.charMaxLength
}

func (c *Column) SetCharMaxLength(length int) {
	c.charMaxLength = &length
}

func (c *Column) NumericPrecision() *int {
	return c.numericPrecision
}

func (c *Column) SetNumericPrecision(precision int) {
	c.numericPrecision = &precision
}

func (c *Column) NumericScale() *int {
	return c.numericScale
}

func (c *Column) SetNumericScale(scale int) {
	c.numericScale = &scale
}

func (c *Column) Table() *Table {
	return c.table
}

func (c *Column) SetTable(table *Table) {
	c.table = table
}
