package dto

import "errors"

type LoginReq struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (r *LoginReq) Validate() error {
	if r.Login == "" {
		return errors.New("login required")
	}
	if r.Password == "" {
		return errors.New("password required")
	}
	return nil
}
