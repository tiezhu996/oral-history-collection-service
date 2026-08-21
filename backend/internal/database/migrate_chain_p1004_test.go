package database

import (
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestMigrateChainErrP1004(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	db, err := gorm.Open(mysql.New(mysql.Config{Conn: sqlDB, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("CREATE TABLE")).WillReturnError(errors.New("boom"))
	mock.ExpectRollback()
	if err := Migrate(db); err == nil {
		t.Fatal("expected migrate error")
	} else if !errors.Is(err, ErrMigrate) {
		t.Fatalf("error chain broken: got %v", err)
	}
}
