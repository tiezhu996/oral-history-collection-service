package repository

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestUserRepoListBackingStableP802(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewUserRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `users`")).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(3))
	rows1 := sqlmock.NewRows([]string{"id", "username", "display_name", "role"}).
		AddRow(1, "a", "A", "interviewer").AddRow(2, "b", "B", "interviewer").AddRow(3, "c", "C", "interviewer")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` ORDER BY id DESC LIMIT ?")).
		WillReturnRows(rows1)
	first, _, err := repo.List(1, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `users`")).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	rows2 := sqlmock.NewRows([]string{"id", "username", "display_name", "role"}).
		AddRow(9, "z", "Z", "interviewer")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` ORDER BY id DESC LIMIT ?")).
		WillReturnRows(rows2)
	if _, _, err := repo.List(1, 10); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if first[0].Username != "a" {
		t.Fatalf("first snapshot corrupted: got %s", first[0].Username)
	}
}
