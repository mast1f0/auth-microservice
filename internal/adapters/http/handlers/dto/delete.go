package dto

import "errors"

type DeleteUserRequest struct {
	UserID int64 `json:"user_id"`
}

func (u *DeleteUserRequest) Validate() error {
	if u.UserID < 1 {
		return errors.New("user_id cannot be less than 1")
	}
	return nil
}
