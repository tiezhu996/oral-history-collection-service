package database

import (
	"errors"
	"log/slog"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestSeedAdminChainErrP1002(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	db, err := gorm.Open(mysql.New(mysql.Config{Conn: sqlDB, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `users` WHERE username = ?")).
		WithArgs("admin", 1).WillReturnError(errors.New("boom"))
	if err := SeedAdmin(db, slog.Default()); err == nil {
		t.Fatal("expected seed error")
	} else if !errors.Is(err, ErrSeed) {
		t.Fatalf("error chain broken: got %v", err)
	}
}

func TestDatabaseNewNilCfgChainErrP1003(t *testing.T) {
	_, err := New(nil, slog.Default())
	if err == nil {
		t.Fatal("expected error for nil config")
	}
	if !errors.Is(err, ErrDBConnect) {
		t.Fatalf("error chain broken: got %v", err)
	}
}
