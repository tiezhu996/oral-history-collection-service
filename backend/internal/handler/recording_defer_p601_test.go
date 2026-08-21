package handler

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/oralhistory/oralhistory/internal/dto"
	"github.com/oralhistory/oralhistory/internal/middleware"
	"github.com/oralhistory/oralhistory/internal/model"
	"github.com/oralhistory/oralhistory/internal/util"
)

type fakeStorage struct {
	uploaded []string
	removed  []string
}

func (f *fakeStorage) Upload(_ context.Context, objectKey string, _ io.Reader, size int64, contentType string) error {
	_ = contentType
	f.uploaded = append(f.uploaded, objectKey)
	return nil
}
func (f *fakeStorage) Get(ctx context.Context, objectKey string) (*minio.Object, error) { return nil, nil }
func (f *fakeStorage) Remove(ctx context.Context, objectKey string) error {
	f.removed = append(f.removed, objectKey)
	return nil
}

type fakeRecordingSvc struct{ attachErr error }

func (f *fakeRecordingSvc) Create(*model.User, *dto.CreateRecordingRequest) (*model.Recording, error) { return nil, nil }
func (f *fakeRecordingSvc) Get(uint) (*model.Recording, error) { return &model.Recording{AudioKey: "k"}, nil }
func (f *fakeRecordingSvc) List(uint, uint) ([]model.Recording, error) { return nil, nil }
func (f *fakeRecordingSvc) Update(*model.User, uint, *dto.UpdateRecordingRequest) (*model.Recording, error) { return nil, nil }
func (f *fakeRecordingSvc) UpdateSummary(*model.User, uint, string) (*model.Recording, error) { return nil, nil }
func (f *fakeRecordingSvc) AttachAudio(*model.User, uint, string, int) (*model.Recording, error) { return nil, f.attachErr }
func (f *fakeRecordingSvc) Delete(*model.User, uint) error { return nil }
func (f *fakeRecordingSvc) CountByProject(uint) (int64, error) { return 0, nil }

func TestUploadAudioAttachFailCleanupP601(t *testing.T) {
	gin.SetMode(gin.TestMode)
	storage := &fakeStorage{}
	rec := &fakeRecordingSvc{attachErr: errors.New("db attach failed")}
	h := &RecordingHandler{recordingSvc: rec, storageSvc: storage, logger: slog.Default()}

	body := &bytes.Buffer{}
	mw := multipart.NewWriter(body)
	fw, _ := mw.CreateFormFile("file", "a.webm")
	_, _ = fw.Write([]byte("fake-audio"))
	_ = mw.Close()

	req := httptest.NewRequest("POST", "/recordings/1/audio", body)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = req
	c.Set(middleware.ContextKeyClaims, &util.Claims{UserID: 1, Username: "u", Role: "interviewer"})
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.UploadAudio(c)
	if len(storage.removed) == 0 {
		t.Fatalf("expected uploaded object cleanup after attach failure")
	}
}
