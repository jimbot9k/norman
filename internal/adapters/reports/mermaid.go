package reports

import (
	"os"
	"strings"

	dbo "github.com/jimbot9k/norman/internal/core/dbobjects"
)

// MermaidReportWriter generates Mermaid ERD diagrams from database schemas
type MermaidReportWriter struct{}

// GetReportKeys returns the report keys supported by this writer
func (w *MermaidReportWriter) GetReportKeys() []string {
	return []string{"mermaid"}
}

// GetReportFileExtension returns the file extension for Mermaid reports
func (w *MermaidReportWriter) GetReportFileExtension() string {
	return "mmd"
}

// GetReportName returns the name of the Mermaid report
func (w *MermaidReportWriter) GetReportName() string {
	return "Mermaid ERD"
}

// WriteInventoryReport writes a Mermaid ERD diagram to the specified file
func (w *MermaidReportWriter) WriteInventoryReport(filePath string, db *dbo.Database) error {
	mermaid := GenerateMermaidERD(db)
	return os.WriteFile(filePath, []byte(mermaid), 0600)
}

// GenerateMermaidERD generates a Mermaid ERD diagram string from a database
func GenerateMermaidERD(db *dbo.Database) string {
	var sb strings.Builder
	sb.WriteString("erDiagram\n")

	for _, schema := range db.Schemas() {
		writeSchemaEntities(&sb, schema)
	}

	return sb.String()
}

// writeSchemaEntities writes all table entities and relationships for a schema
func writeSchemaEntities(sb *strings.Builder, schema *dbo.Schema) {
	// First pass: write all table entities
	for _, table := range schema.Tables() {
		writeTableEntity(sb, table)
	}

	// Second pass: write relationships (after all entities are defined)
	relationships := collectRelationships(schema)
	for _, rel := range relationships {
		sb.WriteString(rel)
		sb.WriteString("\n")
	}
}

// writeTableEntity writes a single table entity with its columns
func writeTableEntity(sb *strings.Builder, table *dbo.Table) {
	tableName := sanitizeMermaidName(table.Name())
	sb.WriteString("    ")
	sb.WriteString(tableName)
	sb.WriteString(" {\n")

	for _, col := range table.Columns() {
		writeColumnDefinition(sb, table, col)
	}

	sb.WriteString("    }\n")
}

// writeColumnDefinition writes a single column definition line
func writeColumnDefinition(sb *strings.Builder, table *dbo.Table, col *dbo.Column) {
	dataType := normalizeMermaidDataType(col.DataType())
	colName := sanitizeMermaidName(col.Name())

	sb.WriteString("        ")
	sb.WriteString(dataType)
	sb.WriteString(" ")
	sb.WriteString(colName)

	// Add PK/FK markers (Mermaid only supports one marker per column, PK takes precedence)
	isPK := isPrimaryKeyColumn(table, col.Name())
	isFK := isForeignKeyColumn(table, col.Name())

	switch {
	case isPK && isFK:
		sb.WriteString(" PK \"FK\"")
	case isPK:
		sb.WriteString(" PK")
	case isFK:
		sb.WriteString(" FK")
	}

	sb.WriteString("\n")
}

// collectRelationships extracts unique relationships from foreign keys
func collectRelationships(schema *dbo.Schema) []string {
	seen := make(map[string]bool)
	var relationships []string

	for _, table := range schema.Tables() {
		for _, fk := range table.ForeignKeys() {
			rel := formatMermaidRelationship(table, fk)
			if !seen[rel] {
				seen[rel] = true
				relationships = append(relationships, rel)
			}
		}
	}

	return relationships
}

// formatMermaidRelationship formats a foreign key as a Mermaid relationship
// Cardinality notation:
//   - ||--|| : one-to-one
//   - ||--o{ : one-to-many (zero or more)
//   - }o--|| : many-to-one (the FK side can have many records pointing to one)
//   - ||--o| : one-to-zero-or-one
func formatMermaidRelationship(table *dbo.Table, fk *dbo.ForeignKey) string {
	fromTable := sanitizeMermaidName(table.Name())
	toTable := sanitizeMermaidName(fk.ReferencedTable())
	label := sanitizeMermaidName(fk.Name())

	// Many-to-one: the table with FK has many rows pointing to one referenced row
	var sb strings.Builder
	sb.WriteString("    ")
	sb.WriteString(fromTable)
	sb.WriteString(" }o--|| ")
	sb.WriteString(toTable)
	sb.WriteString(" : \"")
	sb.WriteString(label)
	sb.WriteString("\"")

	return sb.String()
}

// sanitizeMermaidName replaces characters that Mermaid doesn't handle well
func sanitizeMermaidName(name string) string {
	// Replace hyphens and spaces with underscores
	result := strings.ReplaceAll(name, "-", "_")
	result = strings.ReplaceAll(result, " ", "_")
	return result
}

// normalizeMermaidDataType simplifies database data types for ERD readability
func normalizeMermaidDataType(dataType string) string {
	dt := strings.ToLower(dataType)

	switch {
	case strings.HasPrefix(dt, "character varying"), strings.HasPrefix(dt, "varchar"):
		return "varchar"
	case strings.HasPrefix(dt, "integer"), dt == "int", dt == "int4":
		return "int"
	case strings.HasPrefix(dt, "bigint"), dt == "int8":
		return "bigint"
	case strings.HasPrefix(dt, "smallint"), dt == "int2":
		return "smallint"
	case strings.HasPrefix(dt, "numeric"), strings.HasPrefix(dt, "decimal"):
		return "decimal"
	case strings.HasPrefix(dt, "boolean"), dt == "bool":
		return "bool"
	case strings.HasPrefix(dt, "timestamp"):
		return "timestamp"
	case dt == "date":
		return "date"
	case dt == "time":
		return "time"
	case dt == "text":
		return "text"
	case strings.HasPrefix(dt, "uuid"):
		return "uuid"
	case strings.HasPrefix(dt, "json"):
		return "json"
	default:
		// Strip parenthetical suffixes like varchar(255)
		if idx := strings.Index(dataType, "("); idx != -1 {
			return strings.ToLower(dataType[:idx])
		}
		return strings.ToLower(dataType)
	}
}

// isPrimaryKeyColumn checks if a column is part of the primary key
func isPrimaryKeyColumn(table *dbo.Table, colName string) bool {
	pk := table.PrimaryKey()
	if pk == nil {
		return false
	}
	for _, pkCol := range pk.Columns() {
		if pkCol.Name() == colName {
			return true
		}
	}
	return false
}

// isForeignKeyColumn checks if a column is part of any foreign key
func isForeignKeyColumn(table *dbo.Table, colName string) bool {
	for _, fk := range table.ForeignKeys() {
		for _, fkCol := range fk.Columns() {
			if fkCol.Name() == colName {
				return true
			}
		}
	}
	return false
}
