package user

import (
	"net"

	"github.com/google/uuid"
)

// User ...
type User struct {
	ID   string
	Name string
	Conn net.Conn
}

// New ...
func New(conn net.Conn) *User {
	return &User{
		ID:   uuid.New().String(),
		Conn: conn,
	}
}
