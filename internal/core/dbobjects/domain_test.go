package dbobjects

import "testing"

func TestNewDomain(t *testing.T) {
	d := NewDomain("email_domain", "CREATE DOMAIN email_domain AS VARCHAR(255) CHECK (VALUE ~* '^.+@.+$')")

	if d.Name() != "email_domain" {
		t.Errorf("expected name 'email_domain', got %q", d.Name())
	}
	if d.Definition() != "CREATE DOMAIN email_domain AS VARCHAR(255) CHECK (VALUE ~* '^.+@.+$')" {
		t.Errorf("expected definition, got %q", d.Definition())
	}
}

func TestDomainName(t *testing.T) {
	d := NewDomain("test_domain", "definition")

	if d.Name() != "test_domain" {
		t.Errorf("expected 'test_domain', got %q", d.Name())
	}
}

func TestDomainDefinition(t *testing.T) {
	d := NewDomain("positive_int", "CREATE DOMAIN positive_int AS INT CHECK (VALUE > 0)")

	if d.Definition() != "CREATE DOMAIN positive_int AS INT CHECK (VALUE > 0)" {
		t.Errorf("expected definition, got %q", d.Definition())
	}
}
