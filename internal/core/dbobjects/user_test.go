package dbobjects

import "testing"

func TestNewUser(t *testing.T) {
	u := NewUser("john_doe")

	if u.Name() != "john_doe" {
		t.Errorf("expected name 'john_doe', got %q", u.Name())
	}
}

func TestUserName(t *testing.T) {
	tests := []struct {
		name     string
		username string
	}{
		{"simple name", "alice"},
		{"with underscore", "bob_smith"},
		{"with numbers", "user123"},
		{"mixed case", "JohnDoe"},
		{"empty string", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := NewUser(tt.username)

			if u.Name() != tt.username {
				t.Errorf("expected %q, got %q", tt.username, u.Name())
			}
		})
	}
}
