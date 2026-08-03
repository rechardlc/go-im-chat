package encrypt

import (
	"errors"
	"testing"
)

func TestHashPasswordAndCompare(t *testing.T) {
	plain := "P@ssw0rd123"
	hashed, err := HashPassword(plain)
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}

	if err := ValidatePassword(hashed, plain); err != nil {
		t.Fatalf("ComparePassword() same password error = %v", err)
	}

	if err := ValidatePassword(hashed, "wrong-password"); err == nil {
		t.Fatal("ComparePassword() expected error for wrong password")
	}
}

func TestHashPasswordTooLong(t *testing.T) {
	plain := make([]byte, MaxPasswordLen+1)
	for i := range plain {
		plain[i] = 'a'
	}

	if _, err := HashPassword(string(plain)); !errors.Is(err, ErrPasswordTooLong) {
		t.Fatalf("HashPassword() error = %v, want ErrPasswordTooLong", err)
	}
}
