package router

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chennupati97/go-links/internal/handler"
	"github.com/chennupati97/go-links/internal/model"
	"github.com/chennupati97/go-links/internal/repository"
	"github.com/chennupati97/go-links/internal/service"
	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

type emptyStore struct{}

func (emptyStore) Save(context.Context, *model.Shortcut) error { return nil }
func (emptyStore) GetByAlias(context.Context, string) (*model.Shortcut, error) {
	return nil, repository.ErrMissing
}
func (emptyStore) All(context.Context) ([]model.Shortcut, error) { return nil, nil }

func TestBuildEngine_Routes(t *testing.T) {
	api := handler.NewShortcutAPI(service.NewAliasManager(emptyStore{}))
	engine := BuildEngine(api)

	tests := []struct {
		name       string
		method     string
		path       string
		wantStatus int
		wantJSON   map[string]any
	}{
		{
			name:       "ready",
			method:     http.MethodGet,
			path:       "/ready",
			wantStatus: http.StatusOK,
			wantJSON:   map[string]any{"ready": true},
		},
		{
			name:       "index shortcuts",
			method:     http.MethodGet,
			path:       "/api/shortcuts",
			wantStatus: http.StatusOK,
		},
		{
			name:       "follow missing alias",
			method:     http.MethodGet,
			path:       "/j/missing",
			wantStatus: http.StatusNotFound,
			wantJSON:   map[string]any{"error": "alias not found"},
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
			var got map[string]any
			if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			for k, v := range tt.wantJSON {
				if got[k] != v {
					t.Fatalf("json[%q] = %#v, want %#v", k, got[k], v)
				}
			}
		})
	}
}
