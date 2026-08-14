package ctxdata

import (
	"testing"
	"time"
)

/**
 * TestTokenManager_parseTokenItIssued 覆盖签发后再用同一密钥解析。
 */
func TestTokenManager_parseTokenItIssued(t *testing.T) {
	issuer, err := NewTokenManager(JWTConfig{
		Secret:             "shared-secret",
		AccessTokenExpire:  time.Hour,
		RefreshTokenExpire: 24 * time.Hour,
	})
	if err != nil {
		t.Fatalf("NewTokenManager() error = %v", err)
	}

	token, err := issuer.GenerateToken("user-1", "13800000000")
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	claims, err := issuer.ParseToken(token)
	if err != nil {
		t.Fatalf("ParseToken() error = %v", err)
	}
	if claims.UserID != "user-1" {
		t.Fatalf("UserID = %q, want user-1", claims.UserID)
	}
}
