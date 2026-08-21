package service

import (
	"errors"
	"log/slog"
	"testing"

	"github.com/oralhistory/oralhistory/internal/model"
	"github.com/oralhistory/oralhistory/internal/repository"
	"github.com/oralhistory/oralhistory/internal/util"
)

type fakeMarkerRepo struct {
	markers []model.TimelineMarker
}

func (f *fakeMarkerRepo) Create(m *model.TimelineMarker) error { f.markers = append(f.markers, *m); return nil }
func (f *fakeMarkerRepo) FindByID(uint) (*model.TimelineMarker, error) { return nil, nil }
func (f *fakeMarkerRepo) ListByProject(uint) ([]model.TimelineMarker, error) { return f.markers, nil }
func (f *fakeMarkerRepo) ListByRecording(uint) ([]model.TimelineMarker, error) { return f.markers, nil }
func (f *fakeMarkerRepo) Update(*model.TimelineMarker) error { return nil }
func (f *fakeMarkerRepo) Delete(uint) error { return nil }

type fakeProjectRepoForMarker struct{ missing bool }

func (f *fakeProjectRepoForMarker) FindByID(uint) (*model.Project, error) {
	if f.missing {
		return nil, repository.ErrNotFound
	}
	return &model.Project{ID: 1}, nil
}
func (f *fakeProjectRepoForMarker) FindByIDForUpdate(uint) (*model.Project, error) { return f.FindByID(1) }
func (f *fakeProjectRepoForMarker) Create(*model.Project) error { return nil }
func (f *fakeProjectRepoForMarker) List(int, int, string) ([]model.Project, int64, error) { return nil, 0, nil }
func (f *fakeProjectRepoForMarker) ListByUser(uint, int, int) ([]model.Project, int64, error) { return nil, 0, nil }
func (f *fakeProjectRepoForMarker) Update(*model.Project) error { return nil }
func (f *fakeProjectRepoForMarker) UpdateStatus(*model.Project) error { return nil }
func (f *fakeProjectRepoForMarker) Delete(uint) error { return nil }
func (f *fakeProjectRepoForMarker) Count() (int64, error) { return 0, nil }

func TestMarkerGroupNoPanicP301(t *testing.T) {
	repo := &fakeMarkerRepo{markers: []model.TimelineMarker{
		{ID: 1, ProjectID: 1, RecordingID: 10, TimestampSecond: 5, Label: "a"},
		{ID: 2, ProjectID: 1, RecordingID: 10, TimestampSecond: 9, Label: "b"},
		{ID: 3, ProjectID: 1, RecordingID: 11, TimestampSecond: 3, Label: "c"},
	}}
	svc := NewTimelineMarkerService(repo, &fakeProjectRepoForMarker{}, nil, slog.Default())
	grouped, err := svc.GroupByRecording(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(grouped) != 2 {
		t.Fatalf("expected 2 recording groups, got %d", len(grouped))
	}
	if grouped[10][5].Label != "a" || grouped[10][9].Label != "b" || grouped[11][3].Label != "c" {
		t.Fatalf("grouping content wrong: %+v", grouped)
	}
}

func TestMarkerGroupEmptyMapP303(t *testing.T) {
	svc := NewTimelineMarkerService(&fakeMarkerRepo{}, &fakeProjectRepoForMarker{}, nil, slog.Default())
	grouped, err := svc.GroupByRecording(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if grouped == nil {
		t.Fatal("expected empty map, got nil")
	}
	if len(grouped) != 0 {
		t.Fatalf("expected empty map, got %d entries", len(grouped))
	}
}

func TestMarkerGroupMissingProjectP304(t *testing.T) {
	svc := NewTimelineMarkerService(&fakeMarkerRepo{}, &fakeProjectRepoForMarker{missing: true}, nil, slog.Default())
	_, err := svc.GroupByRecording(1)
	if err == nil {
		t.Fatal("expected error for missing project")
	}
	var appErr *util.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected app error, got %v", err)
	}
}

func TestMarkerListEmptySliceP305(t *testing.T) {
	svc := NewTimelineMarkerService(&fakeMarkerRepo{markers: nil}, &fakeProjectRepoForMarker{}, nil, slog.Default())
	markers, err := svc.List(1, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if markers == nil {
		t.Fatal("expected empty slice, got nil")
	}
}
