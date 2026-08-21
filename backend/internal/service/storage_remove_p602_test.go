package service

import (
	"context"
	"testing"
)

func TestStorageRemoveEmptyKeyP602(t *testing.T) {
	svc := &storageService{client: nil, bucket: "b", logger: nil}
	err := svc.Remove(context.Background(), "")
	if err == nil {
		t.Fatalf("expected error for empty object key")
	}
}
