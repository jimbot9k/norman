package mysql

import (
	"context"
	"database/sql"
	"fmt"

	dbo "github.com/jimbot9k/norman/internal/core/dbobjects"
)

func (a *MySqlAdapter) MapDatabase() (*dbo.Database, []error) {
	var errors []error
	ctx := context.Background()

	// Get current database name
	var dbName string
	err := a.db.QueryRowContext(ctx, "SELECT DATABASE()").Scan(&dbName)
	if err != nil {
		return nil, []error{fmt.Errorf("failed to get database name: %w", err)}
	}

	db := dbo.NewDatabase(dbName, nil)

	// Map schemas (databases in MySQL)
	schemas, errs := a.mapSchemas(ctx)
	errors = append(errors, errs...)
	for _, schema := range schemas {
		db.AddSchema(schema)
	}

	// Map tables for each schema
	for _, schema := range schemas {
		tables, errs := a.mapTables(ctx, schema.Name())
		errors = append(errors, errs...)
		for _, table := range tables {
			schema.AddTable(table)
		}
	}

	// Map columns for each table
	for _, schema := range db.Schemas() {
		for _, table := range schema.Tables() {
			columns, errs := a.mapColumns(ctx, schema.Name(), table.Name())
			errors = append(errors, errs...)
			for _, col := range columns {
				table.AddColumn(col)
			}
		}
	}

	// Map primary keys
	for _, schema := range db.Schemas() {
		for _, table := range schema.Tables() {
			pk, errs := a.mapPrimaryKey(ctx, schema.Name(), table.Name(), table)
			errors = append(errors, errs...)
			if pk != nil {
				table.SetPrimaryKey(pk)
			}
		}
	}

	// Map indexes
	for _, schema := range db.Schemas() {
		for _, table := range schema.Tables() {
			indexes, errs := a.mapIndexes(ctx, schema.Name(), table.Name(), table)
			errors = append(errors, errs...)
			for _, idx := range indexes {
				table.AddIndex(idx)
			}
		}
	}

	// Map foreign keys
	for _, schema := range db.Schemas() {
		for _, table := range schema.Tables() {
			fks, errs := a.mapForeignKeys(ctx, schema.Name(), table.Name(), table)
			errors = append(errors, errs...)
			for _, fk := range fks {
				table.AddForeignKey(fk)
			}
		}
	}

	// Map constraints (CHECK, UNIQUE, NOT NULL)
	for _, schema := range db.Schemas() {
		for _, table := range schema.Tables() {
			constraints, errs := a.mapConstraints(ctx, schema.Name(), table.Name(), table)
			errors = append(errors, errs...)
			for _, c := range constraints {
				table.AddConstraint(c)
			}
		}
	}

	// Map views
	for _, schema := range db.Schemas() {
		views, errs := a.mapViews(ctx, schema.Name())
		errors = append(errors, errs...)
		for _, view := range views {
			schema.AddView(view)
		}
	}

	// Map functions
	for _, schema := range db.Schemas() {
		functions, errs := a.mapFunctions(ctx, schema.Name())
		errors = append(errors, errs...)
		for _, fn := range functions {
			schema.AddFunction(fn)
		}
	}

	// Map procedures
	for _, schema := range db.Schemas() {
		procedures, errs := a.mapProcedures(ctx, schema.Name())
		errors = append(errors, errs...)
		for _, proc := range procedures {
			schema.AddProcedure(proc)
		}
	}

	// Map triggers
	for _, schema := range db.Schemas() {
		for _, table := range schema.Tables() {
			triggers, errs := a.mapTriggers(ctx, schema.Name(), table.Name())
			errors = append(errors, errs...)
			for _, trigger := range triggers {
				table.AddTrigger(trigger)
			}
		}
	}

	if len(errors) > 0 {
		return db, errors
	}
	return db, nil
}

func (a *MySqlAdapter) mapSchemas(ctx context.Context) ([]*dbo.Schema, []error) {
	// In MySQL, the current database is treated as a schema
	var dbName string
	err := a.db.QueryRowContext(ctx, "SELECT DATABASE()").Scan(&dbName)
	if err != nil {
		return nil, []error{fmt.Errorf("failed to get current database: %w", err)}
	}

	// MySQL doesn't have schema owners like PostgreSQL, using empty string
	schemas := []*dbo.Schema{dbo.NewSchema(dbName, "", nil)}
	return schemas, nil
}

func (a *MySqlAdapter) mapTables(ctx context.Context, schemaName string) ([]*dbo.Table, []error) {
	query := `
		SELECT TABLE_NAME 
		FROM information_schema.TABLES 
		WHERE TABLE_SCHEMA = ? AND TABLE_TYPE = 'BASE TABLE'
		ORDER BY TABLE_NAME`

	rows, err := a.db.QueryContext(ctx, query, schemaName)
	if err != nil {
		return nil, []error{fmt.Errorf("failed to query tables for schema %s: %w", schemaName, err)}
	}
	defer rows.Close()

	var tables []*dbo.Table
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return tables, []error{fmt.Errorf("failed to scan table: %w", err)}
		}
		tables = append(tables, dbo.NewTable(name, nil))
	}
	return tables, nil
}

func (a *MySqlAdapter) mapColumns(ctx context.Context, schemaName, tableName string) ([]*dbo.Column, []error) {
	query := `
		SELECT 
			COLUMN_NAME,
			DATA_TYPE,
			IS_NULLABLE,
			COLUMN_DEFAULT,
			ORDINAL_POSITION,
			CHARACTER_MAXIMUM_LENGTH,
			NUMERIC_PRECISION,
			NUMERIC_SCALE
		FROM information_schema.COLUMNS 
		WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?
		ORDER BY ORDINAL_POSITION`

	rows, err := a.db.QueryContext(ctx, query, schemaName, tableName)
	if err != nil {
		return nil, []error{fmt.Errorf("failed to query columns for %s.%s: %w", schemaName, tableName, err)}
	}
	defer rows.Close()

	var columns []*dbo.Column
	for rows.Next() {
		var name, dataType, isNullable string
		var columnDefault sql.NullString
		var ordinalPosition int
		var charMaxLength, numericPrecision, numericScale sql.NullInt64

		if err := rows.Scan(&name, &dataType, &isNullable, &columnDefault, &ordinalPosition, &charMaxLength, &numericPrecision, &numericScale); err != nil {
			return columns, []error{fmt.Errorf("failed to scan column: %w", err)}
		}

		col := dbo.NewColumn(name, dataType, isNullable == "YES")
		col.SetOrdinalPosition(ordinalPosition)
		if columnDefault.Valid {
			col.SetDefaultValue(columnDefault.String)
		}
		if charMaxLength.Valid {
			col.SetCharMaxLength(int(charMaxLength.Int64))
		}
		if numericPrecision.Valid {
			col.SetNumericPrecision(int(numericPrecision.Int64))
		}
		if numericScale.Valid {
			col.SetNumericScale(int(numericScale.Int64))
		}
		columns = append(columns, col)
	}
	return columns, nil
}

func (a *MySqlAdapter) mapPrimaryKey(ctx context.Context, schemaName, tableName string, table *dbo.Table) (*dbo.PrimaryKey, []error) {
	query := `
		SELECT 
			CONSTRAINT_NAME,
			COLUMN_NAME
		FROM information_schema.KEY_COLUMN_USAGE
		WHERE TABLE_SCHEMA = ? 
			AND TABLE_NAME = ?
			AND CONSTRAINT_NAME = 'PRIMARY'
		ORDER BY ORDINAL_POSITION`

	rows, err := a.db.QueryContext(ctx, query, schemaName, tableName)
	if err != nil {
		return nil, []error{fmt.Errorf("failed to query primary key for %s.%s: %w", schemaName, tableName, err)}
	}
	defer rows.Close()

	var pk *dbo.PrimaryKey
	for rows.Next() {
		var constraintName, columnName string
		if err := rows.Scan(&constraintName, &columnName); err != nil {
			return pk, []error{fmt.Errorf("failed to scan primary key: %w", err)}
		}

		if pk == nil {
			pk = dbo.NewPrimaryKey(constraintName, table, nil)
		}

		if col, exists := table.Columns()[columnName]; exists {
			pk.AddColumn(col)
		}
	}
	return pk, nil
}

func (a *MySqlAdapter) mapIndexes(ctx context.Context, schemaName, tableName string, table *dbo.Table) ([]*dbo.Index, []error) {
	query := `
		SELECT 
			INDEX_NAME,
			INDEX_TYPE,
			NON_UNIQUE,
			COLUMN_NAME
		FROM information_schema.STATISTICS
		WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?
		ORDER BY INDEX_NAME, SEQ_IN_INDEX`

	rows, err := a.db.QueryContext(ctx, query, schemaName, tableName)
	if err != nil {
		return nil, []error{fmt.Errorf("failed to query indexes for %s.%s: %w", schemaName, tableName, err)}
	}
	defer rows.Close()

	indexMap := make(map[string]*dbo.Index)
	var indexOrder []string

	for rows.Next() {
		var indexName, indexType, columnName string
		var nonUnique int
		if err := rows.Scan(&indexName, &indexType, &nonUnique, &columnName); err != nil {
			return nil, []error{fmt.Errorf("failed to scan index: %w", err)}
		}

		idx, exists := indexMap[indexName]
		if !exists {
			isUnique := nonUnique == 0
			idx = dbo.NewIndex(indexName, table, nil, isUnique)
			idx.SetPrimary(indexName == "PRIMARY")
			idx.SetIndexType(dbo.IndexType(indexType))
			indexMap[indexName] = idx
			indexOrder = append(indexOrder, indexName)
		}

		if col, colExists := table.Columns()[columnName]; colExists {
			idx.AddColumn(col)
		}
	}

	var indexes []*dbo.Index
	for _, name := range indexOrder {
		indexes = append(indexes, indexMap[name])
	}
	return indexes, nil
}

func (a *MySqlAdapter) mapForeignKeys(ctx context.Context, schemaName, tableName string, table *dbo.Table) ([]*dbo.ForeignKey, []error) {
	query := `
		SELECT 
			kcu.CONSTRAINT_NAME,
			kcu.COLUMN_NAME,
			kcu.REFERENCED_TABLE_SCHEMA,
			kcu.REFERENCED_TABLE_NAME,
			kcu.REFERENCED_COLUMN_NAME,
			rc.DELETE_RULE,
			rc.UPDATE_RULE
		FROM information_schema.KEY_COLUMN_USAGE kcu
		JOIN information_schema.REFERENTIAL_CONSTRAINTS rc
			ON kcu.CONSTRAINT_NAME = rc.CONSTRAINT_NAME
			AND kcu.TABLE_SCHEMA = rc.CONSTRAINT_SCHEMA
		WHERE kcu.TABLE_SCHEMA = ?
			AND kcu.TABLE_NAME = ?
			AND kcu.REFERENCED_TABLE_NAME IS NOT NULL
		ORDER BY kcu.CONSTRAINT_NAME, kcu.ORDINAL_POSITION`

	rows, err := a.db.QueryContext(ctx, query, schemaName, tableName)
	if err != nil {
		return nil, []error{fmt.Errorf("failed to query foreign keys for %s.%s: %w", schemaName, tableName, err)}
	}
	defer rows.Close()

	fkMap := make(map[string]*dbo.ForeignKey)
	var fkOrder []string

	for rows.Next() {
		var constraintName, columnName, refSchema, refTable, refColumn, deleteRule, updateRule string
		if err := rows.Scan(&constraintName, &columnName, &refSchema, &refTable, &refColumn, &deleteRule, &updateRule); err != nil {
			return nil, []error{fmt.Errorf("failed to scan foreign key: %w", err)}
		}

		fk, exists := fkMap[constraintName]
		if !exists {
			fk = dbo.NewForeignKey(constraintName, refTable)
			fk.SetReferencedSchema(refSchema)
			fk.SetOnDelete(dbo.ReferentialAction(deleteRule))
			fk.SetOnUpdate(dbo.ReferentialAction(updateRule))
			fkMap[constraintName] = fk
			fkOrder = append(fkOrder, constraintName)
		}

		if col, colExists := table.Columns()[columnName]; colExists {
			fk.AddColumn(col)
		}
		refCol := dbo.NewColumn(refColumn, "", false)
		fk.AddReferencedColumn(refCol)
	}

	var fks []*dbo.ForeignKey
	for _, name := range fkOrder {
		fks = append(fks, fkMap[name])
	}
	return fks, nil
}

func (a *MySqlAdapter) mapViews(ctx context.Context, schemaName string) ([]*dbo.View, []error) {
	query := `
		SELECT TABLE_NAME, VIEW_DEFINITION 
		FROM information_schema.VIEWS 
		WHERE TABLE_SCHEMA = ?
		ORDER BY TABLE_NAME`

	rows, err := a.db.QueryContext(ctx, query, schemaName)
	if err != nil {
		return nil, []error{fmt.Errorf("failed to query views for schema %s: %w", schemaName, err)}
	}
	defer rows.Close()

	var views []*dbo.View
	for rows.Next() {
		var name string
		var definition sql.NullString
		if err := rows.Scan(&name, &definition); err != nil {
			return views, []error{fmt.Errorf("failed to scan view: %w", err)}
		}
		def := ""
		if definition.Valid {
			def = definition.String
		}
		views = append(views, dbo.NewView(name, def))
	}
	return views, nil
}

func (a *MySqlAdapter) mapFunctions(ctx context.Context, schemaName string) ([]*dbo.Function, []error) {
	query := `
		SELECT 
			ROUTINE_NAME,
			ROUTINE_DEFINITION,
			DTD_IDENTIFIER
		FROM information_schema.ROUTINES 
		WHERE ROUTINE_SCHEMA = ? AND ROUTINE_TYPE = 'FUNCTION'
		ORDER BY ROUTINE_NAME`

	rows, err := a.db.QueryContext(ctx, query, schemaName)
	if err != nil {
		return nil, []error{fmt.Errorf("failed to query functions for schema %s: %w", schemaName, err)}
	}
	defer rows.Close()

	var functions []*dbo.Function
	for rows.Next() {
		var name string
		var definition, returnType sql.NullString
		if err := rows.Scan(&name, &definition, &returnType); err != nil {
			return functions, []error{fmt.Errorf("failed to scan function: %w", err)}
		}
		def := ""
		if definition.Valid {
			def = definition.String
		}
		fn := dbo.NewFunction(name, def)
		if returnType.Valid {
			fn.SetReturnType(returnType.String)
		}
		fn.SetLanguage("SQL")
		functions = append(functions, fn)
	}
	return functions, nil
}

func (a *MySqlAdapter) mapProcedures(ctx context.Context, schemaName string) ([]*dbo.Procedure, []error) {
	query := `
		SELECT 
			ROUTINE_NAME,
			ROUTINE_DEFINITION
		FROM information_schema.ROUTINES 
		WHERE ROUTINE_SCHEMA = ? AND ROUTINE_TYPE = 'PROCEDURE'
		ORDER BY ROUTINE_NAME`

	rows, err := a.db.QueryContext(ctx, query, schemaName)
	if err != nil {
		return nil, []error{fmt.Errorf("failed to query procedures for schema %s: %w", schemaName, err)}
	}
	defer rows.Close()

	var procedures []*dbo.Procedure
	for rows.Next() {
		var name string
		var definition sql.NullString
		if err := rows.Scan(&name, &definition); err != nil {
			return procedures, []error{fmt.Errorf("failed to scan procedure: %w", err)}
		}
		def := ""
		if definition.Valid {
			def = definition.String
		}
		proc := dbo.NewProcedure(name, def)
		proc.SetLanguage("SQL")
		procedures = append(procedures, proc)
	}
	return procedures, nil
}

func (a *MySqlAdapter) mapTriggers(ctx context.Context, schemaName, tableName string) ([]*dbo.Trigger, []error) {
	query := `
		SELECT 
			TRIGGER_NAME,
			ACTION_TIMING,
			EVENT_MANIPULATION,
			ACTION_STATEMENT
		FROM information_schema.TRIGGERS 
		WHERE TRIGGER_SCHEMA = ? AND EVENT_OBJECT_TABLE = ?
		ORDER BY TRIGGER_NAME, EVENT_MANIPULATION`

	rows, err := a.db.QueryContext(ctx, query, schemaName, tableName)
	if err != nil {
		return nil, []error{fmt.Errorf("failed to query triggers for %s.%s: %w", schemaName, tableName, err)}
	}
	defer rows.Close()

	triggerMap := make(map[string]*dbo.Trigger)
	var triggerOrder []string

	for rows.Next() {
		var name, timing, event, definition string
		if err := rows.Scan(&name, &timing, &event, &definition); err != nil {
			return nil, []error{fmt.Errorf("failed to scan trigger: %w", err)}
		}

		trigger, exists := triggerMap[name]
		if !exists {
			trigger = dbo.NewTrigger(name, definition)
			trigger.SetTiming(dbo.TriggerTiming(timing))
			triggerMap[name] = trigger
			triggerOrder = append(triggerOrder, name)
		}

		trigger.AddEvent(dbo.TriggerEvent(event))
	}

	var triggers []*dbo.Trigger
	for _, name := range triggerOrder {
		triggers = append(triggers, triggerMap[name])
	}
	return triggers, nil
}

func (a *MySqlAdapter) mapConstraints(ctx context.Context, schemaName, tableName string, table *dbo.Table) ([]*dbo.Constraint, []error) {
	var constraints []*dbo.Constraint

	// Map CHECK constraints (MySQL 8.0.16+)
	checkConstraints, errs := a.mapCheckConstraints(ctx, schemaName, tableName, table)
	constraints = append(constraints, checkConstraints...)

	// Map UNIQUE constraints
	uniqueConstraints, errs2 := a.mapUniqueConstraints(ctx, schemaName, tableName, table)
	constraints = append(constraints, uniqueConstraints...)
	errs = append(errs, errs2...)

	// Map NOT NULL constraints from column definitions
	notNullConstraints := a.mapNotNullConstraints(table)
	constraints = append(constraints, notNullConstraints...)

	return constraints, errs
}

func (a *MySqlAdapter) mapCheckConstraints(ctx context.Context, schemaName, tableName string, table *dbo.Table) ([]*dbo.Constraint, []error) {
	// CHECK_CONSTRAINTS table available in MySQL 8.0.16+
	query := `
		SELECT 
			cc.CONSTRAINT_NAME,
			cc.CHECK_CLAUSE
		FROM information_schema.CHECK_CONSTRAINTS cc
		JOIN information_schema.TABLE_CONSTRAINTS tc
			ON cc.CONSTRAINT_NAME = tc.CONSTRAINT_NAME
			AND cc.CONSTRAINT_SCHEMA = tc.CONSTRAINT_SCHEMA
		WHERE tc.TABLE_SCHEMA = ? AND tc.TABLE_NAME = ?
		ORDER BY cc.CONSTRAINT_NAME`

	rows, err := a.db.QueryContext(ctx, query, schemaName, tableName)
	if err != nil {
		// CHECK_CONSTRAINTS may not exist in older MySQL versions
		return nil, nil
	}
	defer rows.Close()

	var constraints []*dbo.Constraint
	for rows.Next() {
		var constraintName, checkClause string
		if err := rows.Scan(&constraintName, &checkClause); err != nil {
			return nil, []error{fmt.Errorf("failed to scan check constraint: %w", err)}
		}

		constraint := dbo.NewConstraint(constraintName, dbo.ConstraintType("CHECK"))
		constraint.SetCheckExpression(checkClause)

		// Try to find column from constraint name (MySQL often names them table_chk_N or column_name)
		for colName, col := range table.Columns() {
			if contains(checkClause, colName) {
				constraint.AddColumn(col)
			}
		}

		constraints = append(constraints, constraint)
	}
	return constraints, nil
}

func (a *MySqlAdapter) mapUniqueConstraints(ctx context.Context, schemaName, tableName string, table *dbo.Table) ([]*dbo.Constraint, []error) {
	query := `
		SELECT 
			tc.CONSTRAINT_NAME,
			kcu.COLUMN_NAME
		FROM information_schema.TABLE_CONSTRAINTS tc
		JOIN information_schema.KEY_COLUMN_USAGE kcu
			ON tc.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
			AND tc.TABLE_SCHEMA = kcu.TABLE_SCHEMA
			AND tc.TABLE_NAME = kcu.TABLE_NAME
		WHERE tc.TABLE_SCHEMA = ? 
			AND tc.TABLE_NAME = ?
			AND tc.CONSTRAINT_TYPE = 'UNIQUE'
		ORDER BY tc.CONSTRAINT_NAME, kcu.ORDINAL_POSITION`

	rows, err := a.db.QueryContext(ctx, query, schemaName, tableName)
	if err != nil {
		return nil, []error{fmt.Errorf("failed to query unique constraints for %s.%s: %w", schemaName, tableName, err)}
	}
	defer rows.Close()

	constraintMap := make(map[string]*dbo.Constraint)
	var constraintOrder []string

	for rows.Next() {
		var constraintName, columnName string
		if err := rows.Scan(&constraintName, &columnName); err != nil {
			return nil, []error{fmt.Errorf("failed to scan unique constraint: %w", err)}
		}

		constraint, exists := constraintMap[constraintName]
		if !exists {
			constraint = dbo.NewConstraint(constraintName, dbo.ConstraintType("UNIQUE"))
			constraintMap[constraintName] = constraint
			constraintOrder = append(constraintOrder, constraintName)
		}

		if col, colExists := table.Columns()[columnName]; colExists {
			constraint.AddColumn(col)
		}
	}

	var constraints []*dbo.Constraint
	for _, name := range constraintOrder {
		constraints = append(constraints, constraintMap[name])
	}
	return constraints, nil
}

func (a *MySqlAdapter) mapNotNullConstraints(table *dbo.Table) []*dbo.Constraint {
	var constraints []*dbo.Constraint
	for _, col := range table.Columns() {
		if !col.IsNullable() {
			constraintName := fmt.Sprintf("%s_%s_not_null", table.Name(), col.Name())
			constraint := dbo.NewConstraint(constraintName, dbo.ConstraintType("NOT NULL"))
			constraint.AddColumn(col)
			constraints = append(constraints, constraint)
		}
	}
	return constraints
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && searchString(s, substr)))
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
