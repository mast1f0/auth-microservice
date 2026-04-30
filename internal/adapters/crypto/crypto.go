package crypto

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

type BcryptHasher struct{}

func (*BcryptHasher) CheckPassword(pwd []byte, hashedPwd []byte) bool {
	err := bcrypt.CompareHashAndPassword(hashedPwd, pwd)
	return err == nil
}

func (*BcryptHasher) HashPassword(pwd string) []byte {
	password := []byte(pwd)
	bytes, err := bcrypt.GenerateFromPassword(password, 10)
	if err != nil {
		log.Printf("Error hash password: %e", err)
	}
	return bytes
}

func NewBcryptHasher() *BcryptHasher {
	return &BcryptHasher{}
}
