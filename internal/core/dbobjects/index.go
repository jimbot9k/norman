package dbobjects

import "encoding/json"

type IndexType string

const (
	IndexTypeBTree IndexType = "btree"
	IndexTypeHash  IndexType = "hash"
)

type Index struct {
	name      string
	table     *Table
	columns   []*Column // ordered slice for column order in index
	isUnique  bool
	isPrimary bool
	indexType IndexType
}

func (i *Index) MarshalJSON() ([]byte, error) {
	columnNames := make([]string, len(i.columns))
	for idx, col := range i.columns {
		columnNames[idx] = col.Name()
	}
	return json.Marshal(struct {
		Name      string    `json:"name"`
		Columns   []string  `json:"columns"`
		IsUnique  bool      `json:"isUnique"`
		IsPrimary bool      `json:"isPrimary"`
		IndexType IndexType `json:"indexType"`
	}{
		Name:      i.name,
		Columns:   columnNames,
		IsUnique:  i.isUnique,
		IsPrimary: i.isPrimary,
		IndexType: i.indexType,
	})
}

func NewIndex(name string, table *Table, columns []*Column, isUnique bool) *Index {
	return &Index{
		name:      name,
		table:     table,
		columns:   columns,
		isUnique:  isUnique,
		indexType: IndexTypeBTree,
	}
}

func (i *Index) Name() string {
	return i.name
}

func (i *Index) Table() *Table {
	return i.table
}

func (i *Index) SetTable(table *Table) {
	i.table = table
}

func (i *Index) Columns() []*Column {
	return i.columns
}

func (i *Index) AddColumn(column *Column) {
	i.columns = append(i.columns, column)
}

func (i *Index) IsUnique() bool {
	return i.isUnique
}

func (i *Index) SetUnique(unique bool) {
	i.isUnique = unique
}

func (i *Index) IsPrimary() bool {
	return i.isPrimary
}

func (i *Index) SetPrimary(primary bool) {
	i.isPrimary = primary
}

func (i *Index) IndexType() IndexType {
	return i.indexType
}

func (i *Index) SetIndexType(indexType IndexType) {
	i.indexType = indexType
}
