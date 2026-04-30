package crypto

import "testing"

func TestBcryptHasher_HashAndCheckPassword(t *testing.T) {
	hasher := NewBcryptHasher()
	password := "password123"

	hash := hasher.HashPassword(password)
	if len(hash) == 0 {
		t.Fatal("expected hash to be generated")
	}

	if !hasher.CheckPassword([]byte(password), hash) {
		t.Fatal("expected password to match hash")
	}

	if hasher.CheckPassword([]byte("wrong-password"), hash) {
		t.Fatal("expected password mismatch")
	}
}
