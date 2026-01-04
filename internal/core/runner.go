package core

import (
	"flag"
	"fmt"
	"os"
	"strings"
)


type Runner struct {
	adapterManager   *AdapterManager
	connectionString string
	reportOutputDir  string
	inventoryReportWriterRegistry map[string]*InventoryReportWriter
}

func NewRunner(adapters []Adapter, reports []InventoryReportWriter) *Runner {

	var reportOptionsMap = map[string]*InventoryReportWriter{}
	for _, r := range reports {
		for _, key := range r.GetReportKeys() {
			reportOptionsMap[key] = &r
		}
	}

	var r = &Runner{
		adapterManager: NewAdapterManager(adapters),
		inventoryReportWriterRegistry: reportOptionsMap,
	}
	return r
}

func (r *Runner) Run() error {

	var outputDir = flag.String("output-dir", "./norman/", "Directory to output reports to")
	var connStr = flag.String("conn", "", "Database connection string")
	var reportCsv *string = flag.String("report-types", "all", "Comma-separated list of report types to generate " + r.reportOptionHelperString())
	flag.Parse()

	if *connStr == "" {
		return fmt.Errorf("connection string is required")
	}

	if *outputDir != "" && !strings.HasSuffix(*outputDir, "/") {
		*outputDir += "/"
	}

	r.connectionString = *connStr
	r.reportOutputDir = *outputDir
	selectedReports := r.parseReportArgument(*reportCsv)

	activeAdapter, err := r.adapterManager.Connect(r.connectionString)
	if err != nil {
		return err
	}

	fmt.Printf("Connected using adapter: %s\n", activeAdapter.UniqueSignature())
	fmt.Println("Mapping database...")
	db, errs := activeAdapter.MapDatabase()
	if len(errs) > 0 {
		for _, e := range errs {
			// Log warnings but continue
			fmt.Fprintf(os.Stderr, "warning: %v\n", e)
		}
	}

	fmt.Printf("Mapped Database: %s\n", db.Name())
	if len(selectedReports) == 0 {
		fmt.Println("No report types specified, skipping report generation.")
	} else {
		os.MkdirAll(r.reportOutputDir, os.ModePerm)
	}
	for writer := range selectedReports {
		fmt.Printf("Generating report: %s\n", (*writer).GetReportName())
		err := (*writer).WriteInventoryReport(strings.ReplaceAll(fmt.Sprintf("%s%s_%s.%s", r.reportOutputDir, db.Name(), (*writer).GetReportName(), (*writer).GetReportFileExtension()), " ", "_"), db)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error generating report %s: %v\n", (*writer).GetReportName(), err)
		} else {
			fmt.Printf("Report %s generated successfully.\n", (*writer).GetReportName())
		}
	}
	return nil
}

func (r *Runner) parseReportArgument(reportArg string) map[*InventoryReportWriter]struct{} {
	var selectedReports = map[*InventoryReportWriter]struct{}{}
	for _, reportKey := range strings.Split(reportArg, ",") {

		reportKey = strings.TrimSpace(reportKey)
		if reportKey == "all"{
			for _, writer := range r.inventoryReportWriterRegistry {
				selectedReports[writer] = struct{}{}
			}
			return selectedReports
		}

		if writer, exists := r.inventoryReportWriterRegistry[reportKey]; exists {
			selectedReports[writer] = struct{}{}
		} else {
			fmt.Fprintf(os.Stderr, "warning: unknown report type '%s' specified, ignoring\n", reportKey)
		}
	}
	return selectedReports
}

func (r *Runner) reportOptionHelperString() string {
	var reportOptionsString = "("
	for key := range r.inventoryReportWriterRegistry {
		if reportOptionsString != "(" {
			reportOptionsString += ", "
		}
		reportOptionsString += key
	}
	if reportOptionsString != "(" {
		reportOptionsString += ","
	}
	reportOptionsString += " all)"
	return reportOptionsString
}
