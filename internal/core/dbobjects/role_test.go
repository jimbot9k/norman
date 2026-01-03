package dbobjects

import "testing"

func TestNewRole(t *testing.T) {
	r := NewRole("admin")

	if r.Name() != "admin" {
		t.Errorf("expected name 'admin', got %q", r.Name())
	}
	if r.IsSuperuser() {
		t.Error("expected IsSuperuser to be false by default")
	}
	if r.CanLogin() {
		t.Error("expected CanLogin to be false by default")
	}
	if r.CanCreateDB() {
		t.Error("expected CanCreateDB to be false by default")
	}
	if r.CanCreateRole() {
		t.Error("expected CanCreateRole to be false by default")
	}
	if r.MemberOf() == nil {
		t.Error("expected MemberOf to be initialized")
	}
	if len(r.MemberOf()) != 0 {
		t.Errorf("expected empty MemberOf, got %d", len(r.MemberOf()))
	}
}

func TestRoleSuperuser(t *testing.T) {
	r := NewRole("superuser")

	r.SetSuperuser(true)
	if !r.IsSuperuser() {
		t.Error("expected IsSuperuser to be true")
	}

	r.SetSuperuser(false)
	if r.IsSuperuser() {
		t.Error("expected IsSuperuser to be false")
	}
}

func TestRoleCanLogin(t *testing.T) {
	r := NewRole("user")

	r.SetCanLogin(true)
	if !r.CanLogin() {
		t.Error("expected CanLogin to be true")
	}

	r.SetCanLogin(false)
	if r.CanLogin() {
		t.Error("expected CanLogin to be false")
	}
}

func TestRoleCanCreateDB(t *testing.T) {
	r := NewRole("dbcreator")

	r.SetCanCreateDB(true)
	if !r.CanCreateDB() {
		t.Error("expected CanCreateDB to be true")
	}

	r.SetCanCreateDB(false)
	if r.CanCreateDB() {
		t.Error("expected CanCreateDB to be false")
	}
}

func TestRoleCanCreateRole(t *testing.T) {
	r := NewRole("rolemanager")

	r.SetCanCreateRole(true)
	if !r.CanCreateRole() {
		t.Error("expected CanCreateRole to be true")
	}

	r.SetCanCreateRole(false)
	if r.CanCreateRole() {
		t.Error("expected CanCreateRole to be false")
	}
}

func TestRoleMemberOf(t *testing.T) {
	r := NewRole("developer")
	adminRole := NewRole("admin")
	readersRole := NewRole("readers")

	r.AddMemberOf(adminRole)
	r.AddMemberOf(readersRole)

	memberOf := r.MemberOf()
	if len(memberOf) != 2 {
		t.Fatalf("expected 2 memberships, got %d", len(memberOf))
	}
	if memberOf[0].Name() != "admin" {
		t.Errorf("expected first membership 'admin', got %q", memberOf[0].Name())
	}
	if memberOf[1].Name() != "readers" {
		t.Errorf("expected second membership 'readers', got %q", memberOf[1].Name())
	}
}

func TestRoleAllPermissions(t *testing.T) {
	r := NewRole("superadmin")
	r.SetSuperuser(true)
	r.SetCanLogin(true)
	r.SetCanCreateDB(true)
	r.SetCanCreateRole(true)

	if !r.IsSuperuser() {
		t.Error("expected IsSuperuser to be true")
	}
	if !r.CanLogin() {
		t.Error("expected CanLogin to be true")
	}
	if !r.CanCreateDB() {
		t.Error("expected CanCreateDB to be true")
	}
	if !r.CanCreateRole() {
		t.Error("expected CanCreateRole to be true")
	}
}
