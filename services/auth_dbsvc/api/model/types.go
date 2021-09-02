package model

// User represents an user entity
type User struct {
	ID       int64  `json:"id,omitempty" db:"id"`
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
	RoleID   int64  `json:"role_id" db:"role_id"`
	Role     *Role  `pg:"rel:has-one"`
}

// Role represents an role entity
type Role struct {
	ID   int64  `json:"id,omitempty" db:"id"`
	Role string `json:"role" db:"role"`
}
