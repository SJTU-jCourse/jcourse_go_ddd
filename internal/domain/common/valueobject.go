package common

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

// SystemUser represents the system user for internal operations
var SystemUser = &User{
	UserID: 0, // System user ID
	Role:   RoleAdmin,
}
