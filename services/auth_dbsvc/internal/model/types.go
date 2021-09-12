package model

// User represents an user entity
type User struct {
	Id       int64  `json:"id,omitempty" db:"id"`
	KongId   string `json:"kong_id,omitempty" db:"kong_id"`
	Username string `json:"username" db:"username" validate:"required"`
	Password string `json:"password" db:"password" validate:"required"`
	Role     string `json:"role" db:"role" validate:"required,oneof=Ship Station Command"`
}

//Role represents an role entity
type Role struct {
	Role string `json:"role" db:"role"`
}
