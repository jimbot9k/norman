package main

import (
	"fmt"
	"os"

	"github.com/jimbot9k/norman/internal/adapters/database/mysql"
	"github.com/jimbot9k/norman/internal/adapters/database/postgres"
	"github.com/jimbot9k/norman/internal/adapters/reports"
	"github.com/jimbot9k/norman/internal/core"
)

func main() {
	adapters := []core.Adapter{
		&postgres.PostgresAdapter{},
		&mysql.MySqlAdapter{},
	}
	reports := []core.InventoryReportWriter{
		&reports.JSONReportWriter{},
		&reports.MermaidReportWriter{},
	}

	runner := core.NewRunner(adapters, reports)
	err := runner.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}
