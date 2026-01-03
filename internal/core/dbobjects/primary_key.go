package dbobjects

import "encoding/json"

type PrimaryKey struct {
	name    string
	table   *Table
	columns []*Column // ordered slice for composite key column order
}

func (pk *PrimaryKey) MarshalJSON() ([]byte, error) {
	columnNames := make([]string, len(pk.columns))
	for i, col := range pk.columns {
		columnNames[i] = col.Name()
	}
	return json.Marshal(struct {
		Name    string   `json:"name"`
		Columns []string `json:"columns"`
	}{
		Name:    pk.name,
		Columns: columnNames,
	})
}

func NewPrimaryKey(name string, table *Table, columns []*Column) *PrimaryKey {
	if columns == nil {
		columns = []*Column{}
	}
	return &PrimaryKey{
		name:    name,
		table:   table,
		columns: columns,
	}
}

func (pk *PrimaryKey) Name() string {
	return pk.name
}

func (pk *PrimaryKey) Table() *Table {
	return pk.table
}

func (pk *PrimaryKey) SetTable(table *Table) {
	pk.table = table
}

func (pk *PrimaryKey) Columns() []*Column {
	return pk.columns
}

func (pk *PrimaryKey) AddColumn(column *Column) {
	pk.columns = append(pk.columns, column)
}
