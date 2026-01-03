package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/jimbot9k/norman/adapters/postgres"
	"github.com/jimbot9k/norman/core"
)

func main() {
	adapters := []core.Adapter{
		&postgres.PostgresAdapter{},
	}
	adapterManager := core.NewAdapterManager(adapters)
	activeAdapter, err := adapterManager.Connect("postgres://Jimbobby:Jimbobby@localhost:5432/TESTING_DATABASE?sslmode=disable")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	db, errs := activeAdapter.MapDatabase()
	if len(errs) > 0 {
		for _, e := range errs {
			fmt.Fprintf(os.Stderr, "warning: %v\n", e)
		}
	}

	fmt.Printf("Database: %s\n", db.Name())
	for _, schema := range db.Schemas() {
		fmt.Printf("  Schema: %s (owner: %s)\n", schema.Name(), schema.Owner())
		for _, table := range schema.Tables() {
			fmt.Printf("    Table: %s\n", table.Name())
			for _, col := range table.Columns() {
				defVal := ""
				if col.DefaultValue() != nil {
					defVal = fmt.Sprintf(", default: %s", *col.DefaultValue())
				}
				fmt.Printf("      Column: %s (%s, nullable: %v%s)\n", col.Name(), col.DataType(), col.IsNullable(), defVal)
			}
			if pk := table.PrimaryKey(); pk != nil {
				var pkCols []string
				for _, c := range pk.Columns() {
					pkCols = append(pkCols, c.Name())
				}
				fmt.Printf("      PrimaryKey: %s (%v)\n", pk.Name(), pkCols)
			}
			for _, idx := range table.Indexes() {
				var idxCols []string
				for _, c := range idx.Columns() {
					idxCols = append(idxCols, c.Name())
				}
				fmt.Printf("      Index: %s (unique: %v, type: %s, columns: %v)\n", idx.Name(), idx.IsUnique(), idx.IndexType(), idxCols)
			}
			for _, fk := range table.ForeignKeys() {
				var fkCols, refCols []string
				for _, c := range fk.Columns() {
					fkCols = append(fkCols, c.Name())
				}
				for _, c := range fk.ReferencedColumns() {
					refCols = append(refCols, c.Name())
				}
				fmt.Printf("      ForeignKey: %s (%v) -> %s.%s (%v) [ON DELETE %s, ON UPDATE %s]\n",
					fk.Name(), fkCols, fk.ReferencedSchema(), fk.ReferencedTable(), refCols, fk.OnDelete(), fk.OnUpdate())
			}
			for _, c := range table.Constraints() {
				var conCols []string
				for _, col := range c.Columns() {
					conCols = append(conCols, col.Name())
				}
				checkExpr := ""
				if c.CheckExpression() != "" {
					checkExpr = fmt.Sprintf(": %s", c.CheckExpression())
				}
				fmt.Printf("      Constraint: %s (%s, columns: %v%s)\n", c.Name(), c.Type(), conCols, checkExpr)
			}
			for _, t := range table.Triggers() {
				fmt.Printf("      Trigger: %s (%s %v, forEach: %s)\n", t.Name(), t.Timing(), t.Events(), t.ForEach())
			}
		}
		for _, view := range schema.Views() {
			fmt.Printf("    View: %s\n", view.Name())
		}
		for _, fn := range schema.Functions() {
			fmt.Printf("    Function: %s (returns: %s, language: %s)\n", fn.Name(), fn.ReturnType(), fn.Language())
		}
		for _, proc := range schema.Procedures() {
			fmt.Printf("    Procedure: %s (language: %s)\n", proc.Name(), proc.Language())
		}
		for _, seq := range schema.Sequences() {
			fmt.Printf("    Sequence: %s (start: %d, increment: %d, min: %d, max: %d, cycle: %v)\n",
				seq.Name(), seq.StartValue(), seq.Increment(), seq.MinValue(), seq.MaxValue(), seq.Cycle())
		}

		
		
	}

	json, err := json.MarshalIndent(db, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error marshaling db to JSON: %v\n", err)
		return
	}
	os.WriteFile(fmt.Sprintf("%s_schema.json", db.Name()), json, 0644)
	
}
