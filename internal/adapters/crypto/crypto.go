package crypto

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func CheckPwd(pwd []byte, hashedPwd []byte) bool {
	err := bcrypt.CompareHashAndPassword(hashedPwd, []byte(pwd))
	return err == nil
}

func HashPassword(pwd string) []byte {
	pswrd := []byte(pwd)
	bytes, err := bcrypt.GenerateFromPassword(pswrd, 10)
	if err != nil {
		log.Printf("Error hash password: %e", err)
	}
	return bytes
}
