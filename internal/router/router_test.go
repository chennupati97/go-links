package router

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Naveen-kumar525/go-links/internal/handler"
	"github.com/Naveen-kumar525/go-links/internal/model"
	"github.com/Naveen-kumar525/go-links/internal/repository"
	"github.com/Naveen-kumar525/go-links/internal/service"
	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

type stubRepo struct{}

func (stubRepo) Create(context.Context, *model.Link) error { return nil }
func (stubRepo) FindBySlug(context.Context, string) (*model.Link, error) {
	return nil, repository.ErrNotFound
}
func (stubRepo) List(context.Context) ([]model.Link, error) { return nil, nil }

func TestNew_Routes(t *testing.T) {
	h := handler.NewLinkHandler(service.NewLinkService(stubRepo{}))
	engine := New(h)

	tests := []struct {
		name       string
		method     string
		path       string
		wantStatus int
		wantJSON   map[string]string
	}{
		{
			name:       "health",
			method:     http.MethodGet,
			path:       "/health",
			wantStatus: http.StatusOK,
			wantJSON:   map[string]string{"status": "ok"},
		},
		{
			name:       "list links",
			method:     http.MethodGet,
			path:       "/api/links",
			wantStatus: http.StatusOK,
		},
		{
			name:       "redirect missing slug",
			method:     http.MethodGet,
			path:       "/go/missing",
			wantStatus: http.StatusNotFound,
			wantJSON:   map[string]string{"error": "shortcut not found"},
		},
		{
			name:       "unknown route",
			method:     http.MethodGet,
			path:       "/nope",
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()
			engine.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d; body=%s", w.Code, tt.wantStatus, w.Body.String())
			}
			if tt.wantJSON == nil {
				return
			}
			var got map[string]string
			if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			for k, v := range tt.wantJSON {
				if got[k] != v {
					t.Fatalf("json[%q] = %q, want %q", k, got[k], v)
				}
			}
		})
	}
}
