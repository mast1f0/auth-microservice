package domain

import "testing"

func TestDomainErrors(t *testing.T) {
	if ErrUserNotFound == nil {
		t.Fatal("ErrUserNotFound must be defined")
	}
	if ErrInvalidCredentials == nil {
		t.Fatal("ErrInvalidCredentials must be defined")
	}
	if ErrLoginTooShort == nil {
		t.Fatal("ErrLoginTooShort must be defined")
	}
	if ErrPasswordTooShort == nil {
		t.Fatal("ErrPasswordTooShort must be defined")
	}
}
