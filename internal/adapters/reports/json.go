// Package reports provides report generation for database inventory.
package reports

import (
	"encoding/json"
	"os"

	dbo "github.com/jimbot9k/norman/internal/core/dbobjects"
)

// JSONReportWriter generates JSON format inventory reports.
// It implements the ReportWriter interface for JSON output.
type JSONReportWriter struct{}

// GetReportKeys returns the report keys supported by this writer
func (w *JSONReportWriter) GetReportKeys() []string {
	return []string{"json"}
}

// GetReportFileExtension returns the file extension for JSON reports
func (w *JSONReportWriter) GetReportFileExtension() string {
	return "json"
}

// GetReportName returns the name of the JSON report
func (w *JSONReportWriter) GetReportName() string {
	return "JSON Report"
}

// WriteInventoryReport writes the database inventory to a JSON file.
// The output is formatted with indentation for readability.
func (w *JSONReportWriter) WriteInventoryReport(filePath string, db *dbo.Database) error {
	data, err := marshalDatabaseIndent(db, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, 0600)
}

// JSON representation structs for serialization (private)

// columnJSON represents a database column in JSON format.
type columnJSON struct {
	Name             string  `json:"name"`
	DataType         string  `json:"dataType"`
	Nullable         bool    `json:"nullable"`
	DefaultValue     *string `json:"defaultValue,omitempty"`
	OrdinalPosition  int     `json:"ordinalPosition"`
	CharMaxLength    *int    `json:"charMaxLength,omitempty"`
	NumericPrecision *int    `json:"numericPrecision,omitempty"`
	NumericScale     *int    `json:"numericScale,omitempty"`
}

// constraintJSON represents a table constraint in JSON format.
type constraintJSON struct {
	Name            string             `json:"name"`
	Type            dbo.ConstraintType `json:"type"`
	Columns         []string           `json:"columns"`
	CheckExpression string             `json:"checkExpression,omitempty"`
}

// foreignKeyJSON represents a foreign key relationship in JSON format.
type foreignKeyJSON struct {
	Name              string                `json:"name"`
	Columns           []string              `json:"columns"`
	ReferencedSchema  string                `json:"referencedSchema"`
	ReferencedTable   string                `json:"referencedTable"`
	ReferencedColumns []string              `json:"referencedColumns"`
	OnDelete          dbo.ReferentialAction `json:"onDelete"`
	OnUpdate          dbo.ReferentialAction `json:"onUpdate"`
}

// functionParameterJSON represents a function or procedure parameter in JSON format.
type functionParameterJSON struct {
	Name     string            `json:"name"`
	DataType string            `json:"dataType"`
	Mode     dbo.ParameterMode `json:"mode"`
}

// functionJSON represents a database function in JSON format.
type functionJSON struct {
	Name       string                  `json:"name"`
	Definition string                  `json:"definition"`
	ReturnType string                  `json:"returnType"`
	Parameters []functionParameterJSON `json:"parameters,omitempty"`
	Language   string                  `json:"language"`
}

// indexJSON represents a database index in JSON format.
type indexJSON struct {
	Name      string        `json:"name"`
	Columns   []string      `json:"columns"`
	IsUnique  bool          `json:"isUnique"`
	IsPrimary bool          `json:"isPrimary"`
	IndexType dbo.IndexType `json:"indexType"`
}

// primaryKeyJSON represents a primary key in JSON format.
type primaryKeyJSON struct {
	Name    string   `json:"name"`
	Columns []string `json:"columns"`
}

// procedureJSON represents a stored procedure in JSON format.
type procedureJSON struct {
	Name       string                  `json:"name"`
	Definition string                  `json:"definition"`
	Parameters []functionParameterJSON `json:"parameters,omitempty"`
	Language   string                  `json:"language"`
}

// sequenceJSON represents a database sequence in JSON format.
type sequenceJSON struct {
	Name       string `json:"name"`
	StartValue int64  `json:"startValue"`
	Increment  int64  `json:"increment"`
	MinValue   int64  `json:"minValue"`
	MaxValue   int64  `json:"maxValue"`
	Cache      int64  `json:"cache"`
	Cycle      bool   `json:"cycle"`
}

// triggerJSON represents a database trigger in JSON format.
type triggerJSON struct {
	Name       string             `json:"name"`
	Definition string             `json:"definition"`
	Timing     dbo.TriggerTiming  `json:"timing"`
	Events     []dbo.TriggerEvent `json:"events"`
	Function   string             `json:"function,omitempty"`
	ForEach    string             `json:"forEach"`
}

// viewJSON represents a database view in JSON format.
type viewJSON struct {
	Name       string       `json:"name"`
	Definition string       `json:"definition"`
	Columns    []columnJSON `json:"columns,omitempty"`
}

// tableJSON represents a database table in JSON format.
type tableJSON struct {
	Name        string           `json:"name"`
	Columns     []columnJSON     `json:"columns"`
	PrimaryKey  *primaryKeyJSON  `json:"primaryKey,omitempty"`
	ForeignKeys []foreignKeyJSON `json:"foreignKeys,omitempty"`
	Indexes     []indexJSON      `json:"indexes,omitempty"`
	Constraints []constraintJSON `json:"constraints,omitempty"`
	Triggers    []triggerJSON    `json:"triggers,omitempty"`
}

// schemaJSON represents a database schema in JSON format.
type schemaJSON struct {
	Name       string          `json:"name"`
	Owner      string          `json:"owner"`
	Tables     []tableJSON     `json:"tables"`
	Views      []viewJSON      `json:"views"`
	Functions  []functionJSON  `json:"functions"`
	Procedures []procedureJSON `json:"procedures"`
	Sequences  []sequenceJSON  `json:"sequences"`
}

// databaseJSON represents a database in JSON format.
type databaseJSON struct {
	Name    string       `json:"name"`
	Schemas []schemaJSON `json:"schemas"`
}

// Conversion functions from domain objects to JSON structs

// columnToJSON converts a Column domain object to its JSON representation.
func columnToJSON(c *dbo.Column) columnJSON {
	return columnJSON{
		Name:             c.Name(),
		DataType:         c.DataType(),
		Nullable:         c.IsNullable(),
		DefaultValue:     c.DefaultValue(),
		OrdinalPosition:  c.OrdinalPosition(),
		CharMaxLength:    c.CharMaxLength(),
		NumericPrecision: c.NumericPrecision(),
		NumericScale:     c.NumericScale(),
	}
}

// constraintToJSON converts a Constraint domain object to its JSON representation.
func constraintToJSON(c *dbo.Constraint) constraintJSON {
	columnNames := make([]string, len(c.Columns()))
	for i, col := range c.Columns() {
		columnNames[i] = col.Name()
	}
	return constraintJSON{
		Name:            c.Name(),
		Type:            c.Type(),
		Columns:         columnNames,
		CheckExpression: c.CheckExpression(),
	}
}

// foreignKeyToJSON converts a ForeignKey domain object to its JSON representation.
func foreignKeyToJSON(fk *dbo.ForeignKey) foreignKeyJSON {
	columnNames := make([]string, len(fk.Columns()))
	for i, col := range fk.Columns() {
		columnNames[i] = col.Name()
	}
	refColumnNames := make([]string, len(fk.ReferencedColumns()))
	for i, col := range fk.ReferencedColumns() {
		refColumnNames[i] = col.Name()
	}
	return foreignKeyJSON{
		Name:              fk.Name(),
		Columns:           columnNames,
		ReferencedSchema:  fk.ReferencedSchema(),
		ReferencedTable:   fk.ReferencedTable(),
		ReferencedColumns: refColumnNames,
		OnDelete:          fk.OnDelete(),
		OnUpdate:          fk.OnUpdate(),
	}
}

// functionParameterToJSON converts a FunctionParameter domain object to its JSON representation.
func functionParameterToJSON(p *dbo.FunctionParameter) functionParameterJSON {
	return functionParameterJSON{
		Name:     p.Name(),
		DataType: p.DataType(),
		Mode:     p.Mode(),
	}
}

// functionToJSON converts a Function domain object to its JSON representation.
func functionToJSON(f *dbo.Function) functionJSON {
	params := make([]functionParameterJSON, len(f.Parameters()))
	for i, p := range f.Parameters() {
		params[i] = functionParameterToJSON(p)
	}
	return functionJSON{
		Name:       f.Name(),
		Definition: f.Definition(),
		ReturnType: f.ReturnType(),
		Parameters: params,
		Language:   f.Language(),
	}
}

// indexToJSON converts an Index domain object to its JSON representation.
func indexToJSON(i *dbo.Index) indexJSON {
	columnNames := make([]string, len(i.Columns()))
	for idx, col := range i.Columns() {
		columnNames[idx] = col.Name()
	}
	return indexJSON{
		Name:      i.Name(),
		Columns:   columnNames,
		IsUnique:  i.IsUnique(),
		IsPrimary: i.IsPrimary(),
		IndexType: i.IndexType(),
	}
}

// primaryKeyToJSON converts a PrimaryKey domain object to its JSON representation.
// Returns nil if the input is nil.
func primaryKeyToJSON(pk *dbo.PrimaryKey) *primaryKeyJSON {
	if pk == nil {
		return nil
	}
	columnNames := make([]string, len(pk.Columns()))
	for i, col := range pk.Columns() {
		columnNames[i] = col.Name()
	}
	return &primaryKeyJSON{
		Name:    pk.Name(),
		Columns: columnNames,
	}
}

// procedureToJSON converts a Procedure domain object to its JSON representation.
func procedureToJSON(p *dbo.Procedure) procedureJSON {
	params := make([]functionParameterJSON, len(p.Parameters()))
	for i, param := range p.Parameters() {
		params[i] = functionParameterToJSON(param)
	}
	return procedureJSON{
		Name:       p.Name(),
		Definition: p.Definition(),
		Parameters: params,
		Language:   p.Language(),
	}
}

// sequenceToJSON converts a Sequence domain object to its JSON representation.
func sequenceToJSON(s *dbo.Sequence) sequenceJSON {
	return sequenceJSON{
		Name:       s.Name(),
		StartValue: s.StartValue(),
		Increment:  s.Increment(),
		MinValue:   s.MinValue(),
		MaxValue:   s.MaxValue(),
		Cache:      s.Cache(),
		Cycle:      s.Cycle(),
	}
}

// triggerToJSON converts a Trigger domain object to its JSON representation.
func triggerToJSON(t *dbo.Trigger) triggerJSON {
	var functionName string
	if t.Function() != nil {
		functionName = t.Function().Name()
	}
	return triggerJSON{
		Name:       t.Name(),
		Definition: t.Definition(),
		Timing:     t.Timing(),
		Events:     t.Events(),
		Function:   functionName,
		ForEach:    t.ForEach(),
	}
}

// viewToJSON converts a View domain object to its JSON representation.
func viewToJSON(v *dbo.View) viewJSON {
	columns := make([]columnJSON, len(v.Columns()))
	for i, col := range v.Columns() {
		columns[i] = columnToJSON(col)
	}
	return viewJSON{
		Name:       v.Name(),
		Definition: v.Definition(),
		Columns:    columns,
	}
}

// tableToJSON converts a Table domain object to its JSON representation.
func tableToJSON(t *dbo.Table) tableJSON {
	columns := make([]columnJSON, 0, len(t.Columns()))
	for _, col := range t.Columns() {
		columns = append(columns, columnToJSON(col))
	}

	foreignKeys := make([]foreignKeyJSON, len(t.ForeignKeys()))
	for i, fk := range t.ForeignKeys() {
		foreignKeys[i] = foreignKeyToJSON(fk)
	}

	indexes := make([]indexJSON, len(t.Indexes()))
	for i, idx := range t.Indexes() {
		indexes[i] = indexToJSON(idx)
	}

	constraints := make([]constraintJSON, len(t.Constraints()))
	for i, c := range t.Constraints() {
		constraints[i] = constraintToJSON(c)
	}

	triggers := make([]triggerJSON, len(t.Triggers()))
	for i, tr := range t.Triggers() {
		triggers[i] = triggerToJSON(tr)
	}

	return tableJSON{
		Name:        t.Name(),
		Columns:     columns,
		PrimaryKey:  primaryKeyToJSON(t.PrimaryKey()),
		ForeignKeys: foreignKeys,
		Indexes:     indexes,
		Constraints: constraints,
		Triggers:    triggers,
	}
}

// schemaToJSON converts a Schema domain object to its JSON representation.
func schemaToJSON(s *dbo.Schema) schemaJSON {
	tables := make([]tableJSON, 0, len(s.Tables()))
	for _, t := range s.Tables() {
		tables = append(tables, tableToJSON(t))
	}

	views := make([]viewJSON, 0, len(s.Views()))
	for _, v := range s.Views() {
		views = append(views, viewToJSON(v))
	}

	functions := make([]functionJSON, 0, len(s.Functions()))
	for _, f := range s.Functions() {
		functions = append(functions, functionToJSON(f))
	}

	procedures := make([]procedureJSON, 0, len(s.Procedures()))
	for _, p := range s.Procedures() {
		procedures = append(procedures, procedureToJSON(p))
	}

	sequences := make([]sequenceJSON, 0, len(s.Sequences()))
	for _, seq := range s.Sequences() {
		sequences = append(sequences, sequenceToJSON(seq))
	}

	return schemaJSON{
		Name:       s.Name(),
		Owner:      s.Owner(),
		Tables:     tables,
		Views:      views,
		Functions:  functions,
		Procedures: procedures,
		Sequences:  sequences,
	}
}

// databaseToJSON converts a Database domain object to its JSON representation.
func databaseToJSON(d *dbo.Database) databaseJSON {
	schemas := make([]schemaJSON, 0, len(d.Schemas()))
	for _, s := range d.Schemas() {
		schemas = append(schemas, schemaToJSON(s))
	}
	return databaseJSON{
		Name:    d.Name(),
		Schemas: schemas,
	}
}

// marshalDatabaseIndent serializes a Database to indented JSON bytes.
func marshalDatabaseIndent(db *dbo.Database, prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(databaseToJSON(db), prefix, indent)
}
