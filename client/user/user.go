package user

// User ...
type User struct {
	ID   string
	Name string
}

// New ...
func New() *User {
	return &User{}
}
