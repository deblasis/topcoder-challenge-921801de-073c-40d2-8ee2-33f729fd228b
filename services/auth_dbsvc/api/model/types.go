package model

// User represents an user entity
type User struct {
	ID       int64  `json:"id,omitempty" db:"id"`
	Username string `json:"username" db:"username" validate:"required"`
	Password string `json:"password" db:"password" validate:"required"`
	Role     string `json:"role" db:"role" validate:"required,oneof=Ship Station Command"`
}

// Role represents an role entity
type Role struct {
	ID   int64  `json:"id,omitempty" db:"id"`
	Role string `json:"role" db:"role"`
}
