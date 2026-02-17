package models

type User struct {
	Id        uint   `json:"id"`
	Login     string `json:"login"`
	HashedPwd []byte `json:"password"`
}
