package service

import (
	"log/slog"
	"testing"

	"github.com/oralhistory/oralhistory/internal/model"
)

type fakeQuestionRepo struct {
	questions []model.Question
}

func (f *fakeQuestionRepo) Create(q *model.Question) error { f.questions = append(f.questions, *q); return nil }
func (f *fakeQuestionRepo) FindByID(uint) (*model.Question, error) { return nil, nil }
func (f *fakeQuestionRepo) ListByProject(uint) ([]model.Question, error) { return f.questions, nil }
func (f *fakeQuestionRepo) Update(*model.Question) error { return nil }
func (f *fakeQuestionRepo) Delete(uint) error { return nil }
func (f *fakeQuestionRepo) CountByProject(uint) (int64, error) { return 0, nil }

func TestQuestionListBackingStableP401(t *testing.T) {
	repo := &fakeQuestionRepo{questions: []model.Question{
		{ID: 1, Content: "a"}, {ID: 2, Content: "b"}, {ID: 3, Content: "c"},
	}}
	svc := NewQuestionService(repo, nil, slog.Default())
	first, err := svc.ListByProject(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	repo.questions = []model.Question{{ID: 9, Content: "z"}}
	if _, err := svc.ListByProject(1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if first[0].Content != "a" {
		t.Fatalf("first snapshot corrupted: got %s", first[0].Content)
	}
}
