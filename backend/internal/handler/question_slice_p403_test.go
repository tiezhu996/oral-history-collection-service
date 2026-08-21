package handler

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/oralhistory/oralhistory/internal/dto"
	"github.com/oralhistory/oralhistory/internal/model"
)

type fakeQuestionSvc struct {
	questions []model.Question
}

func (s *fakeQuestionSvc) Create(*model.User, uint, *dto.CreateQuestionRequest) (*model.Question, error) { return nil, nil }
func (s *fakeQuestionSvc) ListByProject(uint) ([]model.Question, error) { return s.questions, nil }
func (s *fakeQuestionSvc) Update(*model.User, uint, *dto.UpdateQuestionRequest) (*model.Question, error) { return nil, nil }
func (s *fakeQuestionSvc) Delete(*model.User, uint) error { return nil }
func (s *fakeQuestionSvc) CountByProject(uint) (int64, error) { return 0, nil }

func TestQuestionHandlerListBackingStableP403(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &fakeQuestionSvc{questions: []model.Question{{ID: 1, Content: "a"}, {ID: 2, Content: "b"}, {ID: 3, Content: "c"}}}
	h := &QuestionHandler{questionSvc: svc}
	c1, _ := gin.CreateTestContext(httptest.NewRecorder())
	c1.Params = gin.Params{{Key: "id", Value: "1"}}
	h.ListByProject(c1)
	first := lastQuestions
	svc.questions = []model.Question{{ID: 9, Content: "z"}}
	c2, _ := gin.CreateTestContext(httptest.NewRecorder())
	c2.Params = gin.Params{{Key: "id", Value: "1"}}
	h.ListByProject(c2)
	if first[0].Content != "a" {
		t.Fatalf("first handler snapshot corrupted: got %s", first[0].Content)
	}
}
