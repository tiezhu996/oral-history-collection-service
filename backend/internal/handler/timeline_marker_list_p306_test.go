package handler

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/oralhistory/oralhistory/internal/dto"
	"github.com/oralhistory/oralhistory/internal/model"
)

type nilListMarkerSvc struct{}

func (s *nilListMarkerSvc) Create(*model.User, *dto.CreateTimelineMarkerRequest) (*model.TimelineMarker, error) { return nil, nil }
func (s *nilListMarkerSvc) List(uint, uint) ([]model.TimelineMarker, error) { return nil, nil }
func (s *nilListMarkerSvc) GroupByRecording(uint) (map[uint]map[int]model.TimelineMarker, error) { return nil, nil }
func (s *nilListMarkerSvc) Update(*model.User, uint, *dto.UpdateTimelineMarkerRequest) (*model.TimelineMarker, error) { return nil, nil }
func (s *nilListMarkerSvc) Delete(*model.User, uint) error { return nil }

func TestMarkerListEmptySliceP306(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	h := &TimelineMarkerHandler{markerSvc: &nilListMarkerSvc{}}
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/?project_id=1", nil)
	h.List(c)
	var body struct {
		Data struct {
			List []model.TimelineMarker `json:"list"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Data.List == nil {
		t.Fatal("expected empty list [], got null")
	}
}
