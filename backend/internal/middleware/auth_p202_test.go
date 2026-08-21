package middleware

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/oralhistory/oralhistory/internal/util"
)

func expiredHeaderForTest(secret string) string {
	claims := util.Claims{
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
	return "Bearer " + s
}

func badSignatureHeaderForTest(secret string) string {
	claims := util.Claims{
		UserID:   1,
		Username: "alice",
		Role:     "interviewer",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "oralhistory",
		},
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, _ := tok.SignedString([]byte("another-secret"))
	return "Bearer " + s
}

func TestAuthSignatureInvalidDistinctMessageP204(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/probe", Auth("secret", slog.Default()), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/probe", nil)
	req.Header.Set("Authorization", badSignatureHeaderForTest("secret"))
	router.ServeHTTP(w, req)
	body := w.Body.String()
	if !strings.Contains(body, "令牌签名无效") {
		t.Fatalf("expected signature-specific message, got %s", body)
	}
}

func TestAuthExpiredTokenDistinctMessageP203(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/probe", Auth("secret", slog.Default()), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/probe", nil)
	req.Header.Set("Authorization", expiredHeaderForTest("secret"))
	router.ServeHTTP(w, req)
	body := w.Body.String()
	if !strings.Contains(body, "登录已过期") {
		t.Fatalf("expected expired-specific message, got %s", body)
	}
	if strings.Contains(body, "令牌无效") {
		t.Fatalf("expired token misclassified as invalid: %s", body)
	}
}
