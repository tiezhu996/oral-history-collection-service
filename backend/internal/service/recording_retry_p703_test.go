package service

import (
	"log/slog"
	"testing"

	"github.com/oralhistory/oralhistory/internal/constants"
	"github.com/oralhistory/oralhistory/internal/model"
)

type fakeRecordingRepo struct {
	rec *model.Recording
}

func (f *fakeRecordingRepo) Create(*model.Recording) error { return nil }
func (f *fakeRecordingRepo) FindByID(uint) (*model.Recording, error) { return f.rec, nil }
func (f *fakeRecordingRepo) ListByProject(uint) ([]model.Recording, error) { return nil, nil }
func (f *fakeRecordingRepo) ListByQuestion(uint) ([]model.Recording, error) { return nil, nil }
func (f *fakeRecordingRepo) FindByIDForUpdate(uint) (*model.Recording, error) { return f.rec, nil }
func (f *fakeRecordingRepo) Update(*model.Recording) error { return nil }
func (f *fakeRecordingRepo) UpdateStatus(*model.Recording) error { return nil }
func (f *fakeRecordingRepo) Delete(uint) error { return nil }
func (f *fakeRecordingRepo) CountByProject(uint) (int64, error) { return 0, nil }
func (f *fakeRecordingRepo) ListByStatus(string) ([]model.Recording, error) { return nil, nil }

func TestRecordingRetryStatusP703(t *testing.T) {
	repo := &fakeRecordingRepo{rec: &model.Recording{ID: 1, Status: constants.RecordingStatusFailed}}
	svc := NewRecordingService(repo, nil, nil, slog.Default())
	got, err := svc.Retry(&model.User{ID: 1, Username: "u", Role: "interviewer"}, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Status != constants.RecordingStatusRetrying {
		t.Fatalf("expected retrying status, got %s", got.Status)
	}
}

func TestRecordingAttachAudioRetryingToReadyP705(t *testing.T) {
	repo := &fakeRecordingRepo{rec: &model.Recording{ID: 1, Status: constants.RecordingStatusRetrying}}
	svc := NewRecordingService(repo, nil, nil, slog.Default())
	got, err := svc.AttachAudio(&model.User{ID: 1, Username: "u", Role: "interviewer"}, 1, "k", 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Status != constants.RecordingStatusReady {
		t.Fatalf("expected ready status after attach, got %s", got.Status)
	}
}
