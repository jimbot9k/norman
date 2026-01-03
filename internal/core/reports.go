package core

import dbo "github.com/jimbot9k/norman/internal/core/dbobjects"

type InventoryReportWriter interface {
	WriteInventoryReport(filePath string, db *dbo.Database) error
	GetReportKeys() []string
	GetReportFileExtension() string
	GetReportName() string
}
