package repository

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestQuestionRepoListBackingStableP402(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewQuestionRepository(db)

	rows1 := sqlmock.NewRows([]string{"id", "project_id", "content", "sort_order"}).
		AddRow(1, 1, "a", 1).AddRow(2, 1, "b", 2).AddRow(3, 1, "c", 3)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `questions` WHERE project_id = ? ORDER BY sort_order ASC, id ASC")).
		WithArgs(1).WillReturnRows(rows1)
	first, err := repo.ListByProject(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	rows2 := sqlmock.NewRows([]string{"id", "project_id", "content", "sort_order"}).
		AddRow(9, 1, "z", 9)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `questions` WHERE project_id = ? ORDER BY sort_order ASC, id ASC")).
		WithArgs(1).WillReturnRows(rows2)
	if _, err := repo.ListByProject(1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if first[0].Content != "a" {
		t.Fatalf("first snapshot corrupted: got %s", first[0].Content)
	}
}
