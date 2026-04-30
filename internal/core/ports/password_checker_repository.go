package ports

type PasswordCheckerRepository interface {
	CheckPassword(pwd []byte, hashedPwd []byte) bool
	HashPassword(pwd string) []byte
}
