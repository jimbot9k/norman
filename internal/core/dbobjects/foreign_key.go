package dbobjects

import "encoding/json"

type ReferentialAction string

const (
	ActionNoAction   ReferentialAction = "NO ACTION"
	ActionRestrict   ReferentialAction = "RESTRICT"
	ActionCascade    ReferentialAction = "CASCADE"
	ActionSetNull    ReferentialAction = "SET NULL"
	ActionSetDefault ReferentialAction = "SET DEFAULT"
)

type ForeignKey struct {
	name              string
	table             *Table
	columns           []*Column
	referencedSchema  string
	referencedTable   string
	referencedColumns []*Column
	onDelete          ReferentialAction
	onUpdate          ReferentialAction
}

func (fk *ForeignKey) MarshalJSON() ([]byte, error) {
	columnNames := make([]string, len(fk.columns))
	for i, col := range fk.columns {
		columnNames[i] = col.Name()
	}
	refColumnNames := make([]string, len(fk.referencedColumns))
	for i, col := range fk.referencedColumns {
		refColumnNames[i] = col.Name()
	}
	return json.Marshal(struct {
		Name              string            `json:"name"`
		Columns           []string          `json:"columns"`
		ReferencedSchema  string            `json:"referencedSchema"`
		ReferencedTable   string            `json:"referencedTable"`
		ReferencedColumns []string          `json:"referencedColumns"`
		OnDelete          ReferentialAction `json:"onDelete"`
		OnUpdate          ReferentialAction `json:"onUpdate"`
	}{
		Name:              fk.name,
		Columns:           columnNames,
		ReferencedSchema:  fk.referencedSchema,
		ReferencedTable:   fk.referencedTable,
		ReferencedColumns: refColumnNames,
		OnDelete:          fk.onDelete,
		OnUpdate:          fk.onUpdate,
	})
}

func NewForeignKey(name string, referencedTable string) *ForeignKey {
	return &ForeignKey{
		name:              name,
		referencedTable:   referencedTable,
		columns:           []*Column{},
		referencedColumns: []*Column{},
		onDelete:          ActionNoAction,
		onUpdate:          ActionNoAction,
	}
}

func (fk *ForeignKey) Name() string {
	return fk.name
}

func (fk *ForeignKey) Table() *Table {
	return fk.table
}

func (fk *ForeignKey) SetTable(table *Table) {
	fk.table = table
}

func (fk *ForeignKey) Columns() []*Column {
	return fk.columns
}

func (fk *ForeignKey) AddColumn(column *Column) {
	fk.columns = append(fk.columns, column)
}

func (fk *ForeignKey) ReferencedSchema() string {
	return fk.referencedSchema
}

func (fk *ForeignKey) SetReferencedSchema(schema string) {
	fk.referencedSchema = schema
}

func (fk *ForeignKey) ReferencedTable() string {
	return fk.referencedTable
}

func (fk *ForeignKey) SetReferencedTable(table string) {
	fk.referencedTable = table
}

func (fk *ForeignKey) ReferencedColumns() []*Column {
	return fk.referencedColumns
}

func (fk *ForeignKey) AddReferencedColumn(column *Column) {
	fk.referencedColumns = append(fk.referencedColumns, column)
}

func (fk *ForeignKey) OnDelete() ReferentialAction {
	return fk.onDelete
}

func (fk *ForeignKey) SetOnDelete(action ReferentialAction) {
	fk.onDelete = action
}

func (fk *ForeignKey) OnUpdate() ReferentialAction {
	return fk.onUpdate
}

func (fk *ForeignKey) SetOnUpdate(action ReferentialAction) {
	fk.onUpdate = action
}
