package util

import (
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func expiredTokenForTest(secret string) string {
	claims := Claims{
		UserID:   1,
		Username: "alice",
		Role:     "interviewer",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "oralhistory",
		},
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, _ := tok.SignedString([]byte(secret))
	return s
}

func TestParseTokenExpiredChainP201(t *testing.T) {
	s := expiredTokenForTest("secret")
	_, err := ParseToken(s, "secret")
	if err == nil {
		t.Fatal("expected error for expired token")
	}
	if !errors.Is(err, jwt.ErrTokenExpired) {
		t.Fatalf("error chain broken: errors.Is(err, jwt.ErrTokenExpired) = false, err=%v", err)
	}
}

func TestParseTokenMalformedChainP202(t *testing.T) {
	_, err := ParseToken("not-a-valid-token", "secret")
	if err == nil {
		t.Fatal("expected error for malformed token")
	}
	if !errors.Is(err, jwt.ErrTokenMalformed) {
		t.Fatalf("error chain broken: errors.Is(err, jwt.ErrTokenMalformed) = false, err=%v", err)
	}
}
