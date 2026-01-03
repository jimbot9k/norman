package dbobjects

import "testing"

func TestNewGrant(t *testing.T) {
	g := NewGrant("grant_select_users", "GRANT SELECT ON users TO reader_role")

	if g.Name() != "grant_select_users" {
		t.Errorf("expected name 'grant_select_users', got %q", g.Name())
	}
	if g.Definition() != "GRANT SELECT ON users TO reader_role" {
		t.Errorf("expected definition, got %q", g.Definition())
	}
}

func TestGrantName(t *testing.T) {
	g := NewGrant("test_grant", "definition")

	if g.Name() != "test_grant" {
		t.Errorf("expected 'test_grant', got %q", g.Name())
	}
}

func TestGrantDefinition(t *testing.T) {
	g := NewGrant("all_privileges", "GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO admin")

	if g.Definition() != "GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO admin" {
		t.Errorf("expected definition, got %q", g.Definition())
	}
}
