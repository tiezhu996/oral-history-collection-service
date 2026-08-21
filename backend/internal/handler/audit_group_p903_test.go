package handler

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/oralhistory/oralhistory/internal/model"
)

type emptyAuditSvc struct{}

func (s *emptyAuditSvc) Record(uint, string, string, string, string, uint, string, string, string) {}
func (s *emptyAuditSvc) RecordActive(context.Context, uint, string, string, string, string, uint, string, string, string) {}
func (s *emptyAuditSvc) List(int, int, string) ([]model.AuditLog, int64, error) { return nil, 0, nil }

func TestAuditHandlerGroupEmptyP903(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	h := &AuditHandler{auditSvc: &emptyAuditSvc{}}
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/?page=1&page_size=10", nil)
	h.Group(c)
	var body struct {
		Data struct {
			Groups map[string]map[string]model.AuditLog `json:"groups"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Data.Groups == nil {
		t.Fatal("expected empty groups map, got null")
	}
}

func TestAuditHandlerListEmptyP904(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	h := &AuditHandler{auditSvc: &emptyAuditSvc{}}
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/?page=1&page_size=10", nil)
	h.List(c)
	var body struct {
		Data struct {
			List []model.AuditLog `json:"list"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Data.List == nil {
		t.Fatal("expected empty list [], got null")
	}
}
