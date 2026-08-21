package util

import (
	"testing"

	"github.com/oralhistory/oralhistory/internal/model"
)

func TestGroupAuditLogsNoPanicP901(t *testing.T) {
	logs := []model.AuditLog{
		{ID: 1, Username: "u", Action: "project.create"},
		{ID: 2, Username: "u", Action: "project.update"},
		{ID: 3, Username: "v", Action: "project.create"},
	}
	grouped := GroupAuditLogs(logs)
	if len(grouped) != 2 {
		t.Fatalf("expected 2 users, got %d", len(grouped))
	}
	if grouped["u"]["project.create"].ID != 1 || grouped["u"]["project.update"].ID != 2 || grouped["v"]["project.create"].ID != 3 {
		t.Fatalf("grouping content wrong: %+v", grouped)
	}
}

func TestGroupAuditLogsEmptyP902(t *testing.T) {
	grouped := GroupAuditLogs(nil)
	if grouped == nil {
		t.Fatalf("expected empty map, got nil")
	}
}
