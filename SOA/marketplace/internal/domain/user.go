package domain

import "time"

type UserRole string

const (
	RoleUser  UserRole = "user"
	RoleAdmin UserRole = "admin"
)

type User struct {
	ID           string
	Email        string
	PasswordHash string
	Role         UserRole
	CreatedAt    time.Time
}

func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}
