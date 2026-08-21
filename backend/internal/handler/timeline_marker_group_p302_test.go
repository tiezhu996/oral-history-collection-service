package handler

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/oralhistory/oralhistory/internal/dto"
	"github.com/oralhistory/oralhistory/internal/model"
)

type emptyMarkerSvc struct{}

func (s *emptyMarkerSvc) Create(*model.User, *dto.CreateTimelineMarkerRequest) (*model.TimelineMarker, error) { return nil, nil }
func (s *emptyMarkerSvc) List(uint, uint) ([]model.TimelineMarker, error) { return nil, nil }
func (s *emptyMarkerSvc) GroupByRecording(uint) (map[uint]map[int]model.TimelineMarker, error) { return map[uint]map[int]model.TimelineMarker{}, nil }
func (s *emptyMarkerSvc) Update(*model.User, uint, *dto.UpdateTimelineMarkerRequest) (*model.TimelineMarker, error) { return nil, nil }
func (s *emptyMarkerSvc) Delete(*model.User, uint) error { return nil }

func TestMarkerGroupEmptyListP302(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	h := &TimelineMarkerHandler{markerSvc: &emptyMarkerSvc{}}
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	h.Group(c)
	var body struct {
		Data map[uint]map[int]model.TimelineMarker `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Data == nil {
		t.Fatalf("expected empty map ({}), got null")
	}
}
