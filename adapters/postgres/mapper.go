package postgres

import (
	"context"
	"fmt"

	dbo "github.com/jimbot9k/norman/core/dbobjects"
)

func (a *PostgresAdapter) MapDatabase() (*dbo.Database, []error) {
	var errors []error
	ctx := context.Background()

	// Get current database name
	var dbName string
	err := a.conn.QueryRow(ctx, "SELECT current_database()").Scan(&dbName)
	if err != nil {
		return nil, []error{fmt.Errorf("failed to get database name: %w", err)}
	}

	db := dbo.NewDatabase(dbName, nil)

	// Map schemas
	schemas, errs := a.mapSchemas(ctx)
	errors = append(errors, errs...)
	for _, schema := range schemas {
		db.AddSchema(schema)
	}

	// Map tables and columns for each schema
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

	// Map sequences
	for _, schema := range db.Schemas() {
		sequences, errs := a.mapSequences(ctx, schema.Name())
		errors = append(errors, errs...)
		for _, seq := range sequences {
			schema.AddSequence(seq)
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

func (a *PostgresAdapter) mapSchemas(ctx context.Context) ([]*dbo.Schema, []error) {
	query := `
		SELECT schema_name, schema_owner 
		FROM information_schema.schemata 
		WHERE schema_name NOT IN ('pg_catalog', 'information_schema', 'pg_toast')
		ORDER BY schema_name`

	rows, err := a.conn.Query(ctx, query)
	if err != nil {
		return nil, []error{fmt.Errorf("failed to query schemas: %w", err)}
	}
	defer rows.Close()

	var schemas []*dbo.Schema
	for rows.Next() {
		var name, owner string
		if err := rows.Scan(&name, &owner); err != nil {
			return schemas, []error{fmt.Errorf("failed to scan schema: %w", err)}
		}
		schemas = append(schemas, dbo.NewSchema(name, owner, nil))
	}
	return schemas, nil
}

func (a *PostgresAdapter) mapTables(ctx context.Context, schemaName string) ([]*dbo.Table, []error) {
	query := `
		SELECT table_name 
		FROM information_schema.tables 
		WHERE table_schema = $1 AND table_type = 'BASE TABLE'
		ORDER BY table_name`

	rows, err := a.conn.Query(ctx, query, schemaName)
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

func (a *PostgresAdapter) mapColumns(ctx context.Context, schemaName, tableName string) ([]*dbo.Column, []error) {
	query := `
		SELECT 
			column_name,
			data_type,
			is_nullable,
			column_default,
			ordinal_position,
			character_maximum_length,
			numeric_precision,
			numeric_scale
		FROM information_schema.columns 
		WHERE table_schema = $1 AND table_name = $2
		ORDER BY ordinal_position`

	rows, err := a.conn.Query(ctx, query, schemaName, tableName)
	if err != nil {
		return nil, []error{fmt.Errorf("failed to query columns for %s.%s: %w", schemaName, tableName, err)}
	}
	defer rows.Close()

	var columns []*dbo.Column
	for rows.Next() {
		var name, dataType, isNullable string
		var columnDefault *string
		var ordinalPosition int
		var charMaxLength, numericPrecision, numericScale *int

		if err := rows.Scan(&name, &dataType, &isNullable, &columnDefault, &ordinalPosition, &charMaxLength, &numericPrecision, &numericScale); err != nil {
			return columns, []error{fmt.Errorf("failed to scan column: %w", err)}
		}

		col := dbo.NewColumn(name, dataType, isNullable == "YES")
		col.SetOrdinalPosition(ordinalPosition)
		if columnDefault != nil {
			col.SetDefaultValue(*columnDefault)
		}
		if charMaxLength != nil {
			col.SetCharMaxLength(*charMaxLength)
		}
		if numericPrecision != nil {
			col.SetNumericPrecision(*numericPrecision)
		}
		if numericScale != nil {
			col.SetNumericScale(*numericScale)
		}
		columns = append(columns, col)
	}
	return columns, nil
}

func (a *PostgresAdapter) mapPrimaryKey(ctx context.Context, schemaName, tableName string, table *dbo.Table) (*dbo.PrimaryKey, []error) {
	query := `
		SELECT 
			tc.constraint_name,
			kcu.column_name
		FROM information_schema.table_constraints tc
		JOIN information_schema.key_column_usage kcu 
			ON tc.constraint_name = kcu.constraint_name 
			AND tc.table_schema = kcu.table_schema
		WHERE tc.constraint_type = 'PRIMARY KEY'
			AND tc.table_schema = $1 
			AND tc.table_name = $2
		ORDER BY kcu.ordinal_position`

	rows, err := a.conn.Query(ctx, query, schemaName, tableName)
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

		// Find the column in the table
		if col, exists := table.Columns()[columnName]; exists {
			pk.AddColumn(col)
		}
	}
	return pk, nil
}

func (a *PostgresAdapter) mapIndexes(ctx context.Context, schemaName, tableName string, table *dbo.Table) ([]*dbo.Index, []error) {
	query := `
		SELECT 
			i.relname AS index_name,
			am.amname AS index_type,
			ix.indisunique AS is_unique,
			ix.indisprimary AS is_primary,
			a.attname AS column_name
		FROM pg_index ix
		JOIN pg_class t ON t.oid = ix.indrelid
		JOIN pg_class i ON i.oid = ix.indexrelid
		JOIN pg_namespace n ON n.oid = t.relnamespace
		JOIN pg_am am ON am.oid = i.relam
		JOIN pg_attribute a ON a.attrelid = t.oid AND a.attnum = ANY(ix.indkey)
		WHERE n.nspname = $1 AND t.relname = $2
		ORDER BY i.relname, array_position(ix.indkey, a.attnum)`

	rows, err := a.conn.Query(ctx, query, schemaName, tableName)
	if err != nil {
		return nil, []error{fmt.Errorf("failed to query indexes for %s.%s: %w", schemaName, tableName, err)}
	}
	defer rows.Close()

	indexMap := make(map[string]*dbo.Index)
	var indexOrder []string

	for rows.Next() {
		var indexName, indexType, columnName string
		var isUnique, isPrimary bool
		if err := rows.Scan(&indexName, &indexType, &isUnique, &isPrimary, &columnName); err != nil {
			return nil, []error{fmt.Errorf("failed to scan index: %w", err)}
		}

		idx, exists := indexMap[indexName]
		if !exists {
			idx = dbo.NewIndex(indexName, table, nil, isUnique)
			idx.SetPrimary(isPrimary)
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

func (a *PostgresAdapter) mapForeignKeys(ctx context.Context, schemaName, tableName string, table *dbo.Table) ([]*dbo.ForeignKey, []error) {
	query := `
		SELECT 
			c.conname AS constraint_name,
			a.attname AS column_name,
			nf.nspname AS ref_schema,
			clf.relname AS ref_table,
			af.attname AS ref_column,
			CASE c.confdeltype
				WHEN 'a' THEN 'NO ACTION'
				WHEN 'r' THEN 'RESTRICT'
				WHEN 'c' THEN 'CASCADE'
				WHEN 'n' THEN 'SET NULL'
				WHEN 'd' THEN 'SET DEFAULT'
			END AS delete_rule,
			CASE c.confupdtype
				WHEN 'a' THEN 'NO ACTION'
				WHEN 'r' THEN 'RESTRICT'
				WHEN 'c' THEN 'CASCADE'
				WHEN 'n' THEN 'SET NULL'
				WHEN 'd' THEN 'SET DEFAULT'
			END AS update_rule
		FROM pg_constraint c
		JOIN pg_class cl ON cl.oid = c.conrelid
		JOIN pg_namespace n ON n.oid = cl.relnamespace
		JOIN pg_class clf ON clf.oid = c.confrelid
		JOIN pg_namespace nf ON nf.oid = clf.relnamespace
		JOIN pg_attribute a ON a.attnum = ANY(c.conkey) AND a.attrelid = c.conrelid
		JOIN pg_attribute af ON af.attnum = ANY(c.confkey) AND af.attrelid = c.confrelid
		WHERE c.contype = 'f'
			AND n.nspname = $1
			AND cl.relname = $2
		ORDER BY c.conname, array_position(c.conkey, a.attnum)`

	rows, err := a.conn.Query(ctx, query, schemaName, tableName)
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
		// Note: referenced columns are stored as names, not full Column objects
		// since we may not have loaded the referenced table yet
		refCol := dbo.NewColumn(refColumn, "", false)
		fk.AddReferencedColumn(refCol)
	}

	var fks []*dbo.ForeignKey
	for _, name := range fkOrder {
		fks = append(fks, fkMap[name])
	}
	return fks, nil
}

func (a *PostgresAdapter) mapViews(ctx context.Context, schemaName string) ([]*dbo.View, []error) {
	query := `
		SELECT table_name, view_definition 
		FROM information_schema.views 
		WHERE table_schema = $1
		ORDER BY table_name`

	rows, err := a.conn.Query(ctx, query, schemaName)
	if err != nil {
		return nil, []error{fmt.Errorf("failed to query views for schema %s: %w", schemaName, err)}
	}
	defer rows.Close()

	var views []*dbo.View
	for rows.Next() {
		var name string
		var definition *string
		if err := rows.Scan(&name, &definition); err != nil {
			return views, []error{fmt.Errorf("failed to scan view: %w", err)}
		}
		def := ""
		if definition != nil {
			def = *definition
		}
		views = append(views, dbo.NewView(name, def))
	}
	return views, nil
}

func (a *PostgresAdapter) mapSequences(ctx context.Context, schemaName string) ([]*dbo.Sequence, []error) {
	query := `
		SELECT 
			sequence_name,
			start_value::bigint,
			increment::bigint,
			minimum_value::bigint,
			maximum_value::bigint,
			cycle_option
		FROM information_schema.sequences 
		WHERE sequence_schema = $1
		ORDER BY sequence_name`

	rows, err := a.conn.Query(ctx, query, schemaName)
	if err != nil {
		return nil, []error{fmt.Errorf("failed to query sequences for schema %s: %w", schemaName, err)}
	}
	defer rows.Close()

	var sequences []*dbo.Sequence
	for rows.Next() {
		var name, cycleOption string
		var startValue, increment, minValue, maxValue int64
		if err := rows.Scan(&name, &startValue, &increment, &minValue, &maxValue, &cycleOption); err != nil {
			return sequences, []error{fmt.Errorf("failed to scan sequence: %w", err)}
		}
		seq := dbo.NewSequence(name, startValue, increment)
		seq.SetMinValue(minValue)
		seq.SetMaxValue(maxValue)
		seq.SetCycle(cycleOption == "YES")
		sequences = append(sequences, seq)
	}
	return sequences, nil
}

func (a *PostgresAdapter) mapFunctions(ctx context.Context, schemaName string) ([]*dbo.Function, []error) {
	query := `
		SELECT 
			p.proname AS function_name,
			pg_get_functiondef(p.oid) AS definition,
			pg_get_function_result(p.oid) AS return_type,
			l.lanname AS language
		FROM pg_proc p
		JOIN pg_namespace n ON n.oid = p.pronamespace
		JOIN pg_language l ON l.oid = p.prolang
		WHERE n.nspname = $1 
			AND p.prokind = 'f'
		ORDER BY p.proname`

	rows, err := a.conn.Query(ctx, query, schemaName)
	if err != nil {
		return nil, []error{fmt.Errorf("failed to query functions for schema %s: %w", schemaName, err)}
	}
	defer rows.Close()

	var functions []*dbo.Function
	for rows.Next() {
		var name, definition, returnType, language string
		if err := rows.Scan(&name, &definition, &returnType, &language); err != nil {
			return functions, []error{fmt.Errorf("failed to scan function: %w", err)}
		}
		fn := dbo.NewFunction(name, definition)
		fn.SetReturnType(returnType)
		fn.SetLanguage(language)
		functions = append(functions, fn)
	}
	return functions, nil
}

func (a *PostgresAdapter) mapProcedures(ctx context.Context, schemaName string) ([]*dbo.Procedure, []error) {
	query := `
		SELECT 
			p.proname AS procedure_name,
			pg_get_functiondef(p.oid) AS definition,
			l.lanname AS language
		FROM pg_proc p
		JOIN pg_namespace n ON n.oid = p.pronamespace
		JOIN pg_language l ON l.oid = p.prolang
		WHERE n.nspname = $1 
			AND p.prokind = 'p'
		ORDER BY p.proname`

	rows, err := a.conn.Query(ctx, query, schemaName)
	if err != nil {
		return nil, []error{fmt.Errorf("failed to query procedures for schema %s: %w", schemaName, err)}
	}
	defer rows.Close()

	var procedures []*dbo.Procedure
	for rows.Next() {
		var name, definition, language string
		if err := rows.Scan(&name, &definition, &language); err != nil {
			return procedures, []error{fmt.Errorf("failed to scan procedure: %w", err)}
		}
		proc := dbo.NewProcedure(name, definition)
		proc.SetLanguage(language)
		procedures = append(procedures, proc)
	}
	return procedures, nil
}

func (a *PostgresAdapter) mapTriggers(ctx context.Context, schemaName, tableName string) ([]*dbo.Trigger, []error) {
	query := `
		SELECT 
			trigger_name,
			action_timing,
			event_manipulation,
			action_statement
		FROM information_schema.triggers 
		WHERE trigger_schema = $1 AND event_object_table = $2
		ORDER BY trigger_name, event_manipulation`

	rows, err := a.conn.Query(ctx, query, schemaName, tableName)
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

func (a *PostgresAdapter) mapConstraints(ctx context.Context, schemaName, tableName string, table *dbo.Table) ([]*dbo.Constraint, []error) {
	query := `
		SELECT 
			c.conname AS constraint_name,
			CASE c.contype
				WHEN 'c' THEN 'CHECK'
				WHEN 'u' THEN 'UNIQUE'
				WHEN 'n' THEN 'NOT NULL'
			END AS constraint_type,
			pg_get_constraintdef(c.oid) AS check_clause,
			COALESCE(
				(SELECT array_agg(a.attname ORDER BY array_position(c.conkey, a.attnum))
				 FROM pg_attribute a 
				 WHERE a.attnum = ANY(c.conkey) AND a.attrelid = c.conrelid),
				'{}'
			) AS column_names
		FROM pg_constraint c
		JOIN pg_class cl ON cl.oid = c.conrelid
		JOIN pg_namespace n ON n.oid = cl.relnamespace
		WHERE c.contype IN ('c', 'u', 'n')
			AND n.nspname = $1
			AND cl.relname = $2
		ORDER BY c.conname`

	rows, err := a.conn.Query(ctx, query, schemaName, tableName)
	if err != nil {
		return nil, []error{fmt.Errorf("failed to query constraints for %s.%s: %w", schemaName, tableName, err)}
	}
	defer rows.Close()

	var constraints []*dbo.Constraint
	for rows.Next() {
		var constraintName, constraintType, checkClause string
		var columnNames []string
		if err := rows.Scan(&constraintName, &constraintType, &checkClause, &columnNames); err != nil {
			return nil, []error{fmt.Errorf("failed to scan constraint: %w", err)}
		}

		constraint := dbo.NewConstraint(constraintName, dbo.ConstraintType(constraintType))
		if constraintType == "CHECK" {
			constraint.SetCheckExpression(checkClause)
		}

		for _, colName := range columnNames {
			if col, exists := table.Columns()[colName]; exists {
				constraint.AddColumn(col)
			}
		}

		constraints = append(constraints, constraint)
	}
	return constraints, nil
}
