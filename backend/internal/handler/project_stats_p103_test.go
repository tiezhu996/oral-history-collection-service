package handler

import (
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/oralhistory/oralhistory/internal/dto"
	"github.com/oralhistory/oralhistory/internal/model"
)

// statsFakeProjectSvc 只服务于 Stats 并发测试，每次返回全新的 map。
type statsFakeProjectSvc struct{}

func (s *statsFakeProjectSvc) Create(*model.User, *dto.CreateProjectRequest) (*model.Project, error) { return nil, nil }
func (s *statsFakeProjectSvc) Get(uint) (*model.Project, error) { return nil, nil }
func (s *statsFakeProjectSvc) List(int, int, string) ([]model.Project, int64, error) { return nil, 0, nil }
func (s *statsFakeProjectSvc) ListMine(uint, int, int) ([]model.Project, int64, error) { return nil, 0, nil }
func (s *statsFakeProjectSvc) Update(*model.User, uint, *dto.UpdateProjectRequest) (*model.Project, error) { return nil, nil }
func (s *statsFakeProjectSvc) TransitionStatus(*model.User, uint, string) (*model.Project, error) { return nil, nil }
func (s *statsFakeProjectSvc) Delete(*model.User, uint) error { return nil }
func (s *statsFakeProjectSvc) Stats() (map[string]any, error) { return map[string]any{"project_total": int64(7)}, nil }

func TestStatsHandlerConcurrentNoRaceP103(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &ProjectHandler{projectSvc: &statsFakeProjectSvc{}}
	const workers = 16
	start := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			h.Stats(c)
		}()
	}
	close(start)
	wg.Wait()
}
