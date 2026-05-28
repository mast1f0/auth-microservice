package dto

import "errors"

type UpdateUserDTO struct {
	Role string `json:"role"`
}

func (u *UpdateUserDTO) Valid() error {
	if u.Role == "" {
		return errors.New("role is required")
	}
	return nil
}
