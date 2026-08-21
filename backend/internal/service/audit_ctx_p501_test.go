package service

import (
	"context"
	"log/slog"
	"testing"

	"github.com/oralhistory/oralhistory/internal/model"
)

type fakeAuditRepo struct {
	created int
}

func (f *fakeAuditRepo) Create(*model.AuditLog) error { f.created++; return nil }
func (f *fakeAuditRepo) List(int, int, string) ([]model.AuditLog, int64, error) { return nil, 0, nil }

func TestAuditRecordActiveCancelledCtxP501(t *testing.T) {
	repo := &fakeAuditRepo{}
	svc := &auditService{auditRepo: repo, logger: slog.Default()}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	svc.RecordActive(ctx, 1, "u", "admin", "probe", "request", 0, "", "1.1.1.1", "rid")
	if repo.created != 0 {
		t.Fatalf("expected no write after ctx cancel, wrote %d", repo.created)
	}
}
