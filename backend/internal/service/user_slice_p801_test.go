package service

import (
	"log/slog"
	"testing"

	"github.com/oralhistory/oralhistory/internal/model"
)

type fakeUserRepoSlice struct {
	users []model.User
}

func (f *fakeUserRepoSlice) Create(u *model.User) error { f.users = append(f.users, *u); return nil }
func (f *fakeUserRepoSlice) FindByID(uint) (*model.User, error) { return nil, nil }
func (f *fakeUserRepoSlice) FindByUsername(string) (*model.User, error) { return nil, nil }
func (f *fakeUserRepoSlice) List(int, int) ([]model.User, int64, error) { return f.users, int64(len(f.users)), nil }
func (f *fakeUserRepoSlice) Update(*model.User) error { return nil }
func (f *fakeUserRepoSlice) Delete(uint) error { return nil }

func TestUserListBackingStableP801(t *testing.T) {
	repo := &fakeUserRepoSlice{users: []model.User{
		{ID: 1, Username: "a"}, {ID: 2, Username: "b"}, {ID: 3, Username: "c"},
	}}
	svc := &userService{userRepo: repo, cfg: nil, logger: slog.Default()}
	first, _, err := svc.List(1, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	repo.users = []model.User{{ID: 9, Username: "z"}}
	if _, _, err := svc.List(1, 10); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if first[0].Username != "a" {
		t.Fatalf("first snapshot corrupted: got %s", first[0].Username)
	}
}
