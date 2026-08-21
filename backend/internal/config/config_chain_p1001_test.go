package config

import (
	"errors"
	"testing"
)

func TestConfigLoadInvalidEnvChainP1001(t *testing.T) {
	t.Setenv("MINIO_USE_SSL", "not-a-bool")
	_, err := Load()
	if err == nil {
		t.Fatal("expected config error")
	}
	if !errors.Is(err, ErrInvalidConfig) {
		t.Fatalf("error chain broken: errors.Is(err, ErrInvalidConfig) = false, err=%v", err)
	}
}

func TestConfigValidateMissingSecretP1005(t *testing.T) {
	cfg := &Config{}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected validation error for default secret")
	} else if !errors.Is(err, ErrInvalidConfig) {
		t.Fatalf("error chain broken: got %v", err)
	}
}
