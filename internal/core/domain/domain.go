package domain

import "time"

type Role string

const (
	RoleAdmin  Role = "admin"
	RoleSeller Role = "user"
	RoleBuyer  Role = "buyer"
)

type User struct {
	Id        int64     `json:"id,omitempty"`
	Login     string    `json:"login"`
	Role      Role      `json:"role"`
	HashedPwd []byte    `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
}
