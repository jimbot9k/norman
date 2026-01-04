package dbobjects

import "encoding/json"

type Schema struct {
	name       string
	owner      string
	database   *Database
	tables     map[string]*Table
	views      map[string]*View
	functions  map[string]*Function
	procedures map[string]*Procedure
	sequences  map[string]*Sequence
}

func (s *Schema) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name       string                `json:"name"`
		Owner      string                `json:"owner"`
		Tables     map[string]*Table     `json:"tables"`
		Views      map[string]*View      `json:"views"`
		Functions  map[string]*Function  `json:"functions"`
		Procedures map[string]*Procedure `json:"procedures"`
		Sequences  map[string]*Sequence  `json:"sequences"`
	}{
		Name:       s.name,
		Owner:      s.owner,
		Tables:     s.tables,
		Views:      s.views,
		Functions:  s.functions,
		Procedures: s.procedures,
		Sequences:  s.sequences,
	})
}

func NewSchema(name string, owner string, tables map[string]*Table) *Schema {
	if tables == nil {
		tables = make(map[string]*Table)
	}
	return &Schema{
		name:       name,
		owner:      owner,
		tables:     tables,
		views:      make(map[string]*View),
		functions:  make(map[string]*Function),
		procedures: make(map[string]*Procedure),
		sequences:  make(map[string]*Sequence),
	}
}

func (s *Schema) Name() string {
	return s.name
}

func (s *Schema) Owner() string {
	return s.owner
}

func (s *Schema) Database() *Database {
	return s.database
}

func (s *Schema) SetDatabase(database *Database) {
	s.database = database
}

func (s *Schema) Tables() map[string]*Table {
	return s.tables
}

func (s *Schema) AddTable(table *Table) {
	table.SetSchema(s)
	s.tables[table.Name()] = table
}

func (s *Schema) Views() map[string]*View {
	return s.views
}

func (s *Schema) AddView(view *View) {
	view.SetSchema(s)
	s.views[view.Name()] = view
}

func (s *Schema) Functions() map[string]*Function {
	return s.functions
}

func (s *Schema) AddFunction(function *Function) {
	function.SetSchema(s)
	s.functions[function.Name()] = function
}

func (s *Schema) Procedures() map[string]*Procedure {
	return s.procedures
}

func (s *Schema) AddProcedure(procedure *Procedure) {
	procedure.SetSchema(s)
	s.procedures[procedure.Name()] = procedure
}

func (s *Schema) Sequences() map[string]*Sequence {
	return s.sequences
}

func (s *Schema) AddSequence(sequence *Sequence) {
	sequence.SetSchema(s)
	s.sequences[sequence.Name()] = sequence
}

// FullyQualifiedName returns database.schema format if database is set
func (s *Schema) FullyQualifiedName() string {
	if s.database != nil {
		return s.database.Name() + "." + s.name
	}
	return s.name
}
