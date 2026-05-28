package dto

import (
	"errors"
)

type UpdateUserDTO struct {
	Role string `json:"role"`
}

func (u *UpdateUserDTO) Valid() error {
	if u.Role == "" {
		return errors.New("role is required")
	}
	switch u.Role {
	case "admin", "seller", "buyer":
		return nil
	default:
		return errors.New("invalid role")

	}
}
