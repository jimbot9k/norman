package dbobjects

import "testing"

func TestNewMaterializedView(t *testing.T) {
	mv := NewMaterializedView("sales_summary", "SELECT date, SUM(amount) FROM sales GROUP BY date")

	if mv.Name() != "sales_summary" {
		t.Errorf("expected name 'sales_summary', got %q", mv.Name())
	}
	if mv.Definition() != "SELECT date, SUM(amount) FROM sales GROUP BY date" {
		t.Errorf("expected definition, got %q", mv.Definition())
	}
}

func TestMaterializedViewName(t *testing.T) {
	mv := NewMaterializedView("test_mv", "SELECT 1")

	if mv.Name() != "test_mv" {
		t.Errorf("expected 'test_mv', got %q", mv.Name())
	}
}

func TestMaterializedViewDefinition(t *testing.T) {
	mv := NewMaterializedView("monthly_stats", "SELECT month, COUNT(*) FROM orders GROUP BY month")

	if mv.Definition() != "SELECT month, COUNT(*) FROM orders GROUP BY month" {
		t.Errorf("expected definition, got %q", mv.Definition())
	}
}
