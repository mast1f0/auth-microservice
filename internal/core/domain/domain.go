package domain

import "time"

type Role string

const (
	RoleAdmin  Role = "admin"
	RoleSeller Role = "seller"
	RoleBuyer  Role = "buyer"
)

type User struct {
	Id        int64     `json:"id,omitempty"`
	Login     string    `json:"login"`
	Role      Role      `json:"role"`
	HashedPwd []byte    `json:"password,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
}
