package service

import (
	"log/slog"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/oralhistory/oralhistory/internal/model"
	"github.com/oralhistory/oralhistory/internal/repository"
)

// statsFakeRepo 仅用于 Stats 相关测试，Count 可通过原子操作改写。
type statsFakeRepo struct {
	count int64
}

func (f *statsFakeRepo) Count() (int64, error) { return atomic.LoadInt64(&f.count), nil }
func (f *statsFakeRepo) Create(*model.Project) error { return nil }
func (f *statsFakeRepo) FindByID(uint) (*model.Project, error) { return nil, repository.ErrNotFound }
func (f *statsFakeRepo) List(int, int, string) ([]model.Project, int64, error) { return nil, 0, nil }
func (f *statsFakeRepo) ListByUser(uint, int, int) ([]model.Project, int64, error) { return nil, 0, nil }
func (f *statsFakeRepo) FindByIDForUpdate(uint) (*model.Project, error) { return nil, repository.ErrNotFound }
func (f *statsFakeRepo) Update(*model.Project) error { return nil }
func (f *statsFakeRepo) UpdateStatus(*model.Project) error { return nil }
func (f *statsFakeRepo) Delete(uint) error { return nil }

func TestStatsConcurrentNoRaceP101(t *testing.T) {
	repo := &statsFakeRepo{count: 42}
	svc := NewProjectService(repo, slog.Default())
	const workers = 16
	start := make(chan struct{})
	var wg sync.WaitGroup
	errCh := make(chan error, workers)
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			if _, err := svc.Stats(); err != nil {
				errCh <- err
			}
		}()
	}
	close(start)
	wg.Wait()
	close(errCh)
	for err := range errCh {
		t.Fatalf("unexpected stats error: %v", err)
	}
}

func TestStatsSnapshotStableP102(t *testing.T) {
	repo := &statsFakeRepo{count: 100}
	svc := NewProjectService(repo, slog.Default())
	first, err := svc.Stats()
	if err != nil {
		t.Fatalf("stats error: %v", err)
	}
	atomic.StoreInt64(&repo.count, 200)
	if _, err := svc.Stats(); err != nil {
		t.Fatalf("stats error: %v", err)
	}
	if got := first["project_total"]; got != int64(100) {
		t.Fatalf("first snapshot was mutated: got %v, want 100", got)
	}
}
