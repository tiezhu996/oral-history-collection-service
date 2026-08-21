package dto

import (
	"testing"

	"github.com/gin-gonic/gin/binding"
)

func TestUpdateRecordingStatusAllowsRetryingP704(t *testing.T) {
	var req UpdateRecordingRequest
	body := []byte(`{"status":"retrying"}`)
	if err := binding.JSON.BindBody(body, &req); err != nil {
		t.Fatalf("retrying should be accepted: %v", err)
	}
	if req.Status != "retrying" {
		t.Fatalf("expected status retrying, got %s", req.Status)
	}
}
