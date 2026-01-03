package dbobjects

type User struct {
	name string
}

func NewUser(name string) *User {
	return &User{
		name: name,
	}
}

func (u *User) Name() string {
	return u.name
}