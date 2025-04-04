package dto

type Role string //TODO common type for all services

const (
	RoleAdmin     = "Admin"
	RoleModerator = "Moderator"
	RoleUser      = "Creator"
)

type User struct { //TODO common type for all services
	ID        int    `json:"id,string,omitempty"`
	Login     string `json:"login,omitempty"`
	UserName  string `json:"user_name,omitempty"`
	Role      Role   `json:"role,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
}
