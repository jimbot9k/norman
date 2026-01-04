package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jimbot9k/norman/internal/adapters/database/mysql"
	"github.com/jimbot9k/norman/internal/adapters/database/postgres"
	"github.com/jimbot9k/norman/internal/adapters/reports"
	"github.com/jimbot9k/norman/internal/core"
	"github.com/jimbot9k/norman/internal/version"
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

	var showVersion = flag.Bool("version", false, "print version information and exit")
	var outputDir = flag.String("output-dir", "./norman/", "Directory to output reports to")
	var connStr = flag.String("conn", "", "Database connection string " + driverOptionHelperString(adapters) + " (required)")
	var reportCsv = flag.String("report-types", "all", "Comma-separated list of report types to generate " + reportOptionHelperString(reports))
	flag.Parse()

	if *showVersion {
		fmt.Printf("norman %s\ncommit: %s\nbuilt:  %s\n", version.Version, version.Commit, version.Date)
		return
	}

	runner := core.NewRunner(adapters, reports)
	err := runner.Run(connStr, outputDir, reportCsv)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func driverOptionHelperString(adapters []core.Adapter) string {
	var driverOptionsString = "(Drivers Available: "
	for i, a := range adapters {
		if i > 0 {
			driverOptionsString += ", "
		}
		driverOptionsString += a.UniqueSignature()
	}
	driverOptionsString += ")"
	return driverOptionsString
}


func reportOptionHelperString(reports []core.InventoryReportWriter) string {
	var reportOptionsString = "("
	for i, r := range reports {
		for _, key := range r.GetReportKeys() {

			if i > 0 {
				reportOptionsString += ", "
			}
			reportOptionsString += key
		}

		if i == len(reports) - 1 {
			reportOptionsString += ","
		}
	}



	reportOptionsString += " all)"
	return reportOptionsString
}
