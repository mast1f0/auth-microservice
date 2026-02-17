package models

import "time"

type User struct {
	Id        uint   `json:"id"`
	Login     string `json:"login"`
	HashedPwd string `json:"password"`
	CreatedAt time.Time
}
