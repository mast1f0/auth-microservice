package crypto

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func CheckPwd(pwd []byte, hashedPwd []byte) bool {
	err := bcrypt.CompareHashAndPassword(hashedPwd, []byte(pwd))
	return err == nil
}

func HashPassword(pwd []byte) []byte {
	bytes, err := bcrypt.GenerateFromPassword(pwd, 10)
	if err != nil {
		log.Printf("Error hash password: %e", err)
	}
	return bytes
}
