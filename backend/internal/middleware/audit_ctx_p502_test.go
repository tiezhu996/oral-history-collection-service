package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/oralhistory/oralhistory/internal/model"
)

type fakeAuditSvc struct{ called int }

func (f *fakeAuditSvc) Record(uint, string, string, string, string, uint, string, string, string) {}
func (f *fakeAuditSvc) RecordActive(context.Context, uint, string, string, string, string, uint, string, string, string) {
	f.called++
}
func (f *fakeAuditSvc) List(int, int, string) ([]model.AuditLog, int64, error) { return nil, 0, nil }

func TestAuditMiddlewareCancelledCtxP502(t *testing.T) {
	gin.SetMode(gin.TestMode)
	fake := &fakeAuditSvc{}
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set(ContextKeyUserID, uint(7))
		c.Set(ContextKeyUsername, "u")
		c.Set(ContextKeyRole, "admin")
		c.Next()
	})
	router.Use(Audit(fake, slog.Default()))
	router.POST("/probe", func(c *gin.Context) { c.Status(http.StatusOK) })

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	req := httptest.NewRequest(http.MethodPost, "/probe", nil).WithContext(ctx)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if fake.called != 0 {
		t.Fatalf("expected no audit write after ctx cancel, got %d", fake.called)
	}
}
