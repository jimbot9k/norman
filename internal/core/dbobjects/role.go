package dbobjects

type Role struct {
	name          string
	isSuperuser   bool
	canLogin      bool
	canCreateDB   bool
	canCreateRole bool
	memberOf      []*Role
}

func NewRole(name string) *Role {
	return &Role{
		name:     name,
		memberOf: []*Role{},
	}
}

func (r *Role) Name() string {
	return r.name
}

func (r *Role) IsSuperuser() bool {
	return r.isSuperuser
}

func (r *Role) SetSuperuser(isSuperuser bool) {
	r.isSuperuser = isSuperuser
}

func (r *Role) CanLogin() bool {
	return r.canLogin
}

func (r *Role) SetCanLogin(canLogin bool) {
	r.canLogin = canLogin
}

func (r *Role) CanCreateDB() bool {
	return r.canCreateDB
}

func (r *Role) SetCanCreateDB(canCreateDB bool) {
	r.canCreateDB = canCreateDB
}

func (r *Role) CanCreateRole() bool {
	return r.canCreateRole
}

func (r *Role) SetCanCreateRole(canCreateRole bool) {
	r.canCreateRole = canCreateRole
}

func (r *Role) MemberOf() []*Role {
	return r.memberOf
}

func (r *Role) AddMemberOf(role *Role) {
	r.memberOf = append(r.memberOf, role)
}
