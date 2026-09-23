package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/chennupati97/go-links/internal/model"
	"github.com/chennupati97/go-links/internal/repository"
	"github.com/chennupati97/go-links/internal/service"
	"github.com/chennupati97/go-links/internal/validation"
	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

type stubStore struct {
	saveFn       func(ctx context.Context, item *model.Shortcut) error
	getByAliasFn func(ctx context.Context, alias string) (*model.Shortcut, error)
	allFn        func(ctx context.Context) ([]model.Shortcut, error)
}

func (s *stubStore) Save(ctx context.Context, item *model.Shortcut) error {
	if s.saveFn != nil {
		return s.saveFn(ctx, item)
	}
	return nil
}

func (s *stubStore) GetByAlias(ctx context.Context, alias string) (*model.Shortcut, error) {
	if s.getByAliasFn != nil {
		return s.getByAliasFn(ctx, alias)
	}
	return nil, repository.ErrMissing
}

func (s *stubStore) All(ctx context.Context) ([]model.Shortcut, error) {
	if s.allFn != nil {
		return s.allFn(ctx)
	}
	return nil, nil
}

func newTestAPI(store *stubStore) *ShortcutAPI {
	return NewShortcutAPI(service.NewAliasManager(store))
}

func TestShortcutAPI_Register(t *testing.T) {
	createdAt := time.Date(2024, 5, 1, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name       string
		body       string
		store      *stubStore
		wantStatus int
		wantAlias  string
		wantErrMsg string
	}{
		{
			name: "created",
			body: `{"alias":"docs","destination":"https://example.com/docs"}`,
			store: &stubStore{
				saveFn: func(_ context.Context, item *model.Shortcut) error {
					item.ID = 10
					item.CreatedAt = createdAt
					return nil
				},
			},
			wantStatus: http.StatusCreated,
			wantAlias:  "docs",
		},
		{
			name:       "malformed json",
			body:       `{not-json`,
			store:      &stubStore{},
			wantStatus: http.StatusBadRequest,
			wantErrMsg: "malformed request body",
		},
		{
			name:       "validation error",
			body:       `{"alias":"bad_slug","destination":"https://example.com"}`,
			store:      &stubStore{},
			wantStatus: http.StatusBadRequest,
			wantErrMsg: "alias may contain only lowercase letters, numbers and hyphens",
		},
		{
			name: "duplicate",
			body: `{"alias":"docs","destination":"https://example.com"}`,
			store: &stubStore{
				saveFn: func(context.Context, *model.Shortcut) error {
					return repository.ErrDuplicateAlias
				},
			},
			wantStatus: http.StatusConflict,
			wantErrMsg: "alias already taken",
		},
		{
			name: "internal error",
			body: `{"alias":"docs","destination":"https://example.com"}`,
			store: &stubStore{
				saveFn: func(context.Context, *model.Shortcut) error {
					return errors.New("boom")
				},
			},
			wantStatus: http.StatusInternalServerError,
			wantErrMsg: "unexpected server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := newTestAPI(tt.store)
			engine := gin.New()
			engine.POST("/api/shortcuts", api.Register)

			req := httptest.NewRequest(http.MethodPost, "/api/shortcuts", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			engine.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d; body=%s", w.Code, tt.wantStatus, w.Body.String())
			}
			if tt.wantErrMsg != "" {
				var resp map[string]string
				if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
					t.Fatalf("decode error body: %v", err)
				}
				if resp["error"] != tt.wantErrMsg {
					t.Fatalf("error = %q, want %q", resp["error"], tt.wantErrMsg)
				}
				return
			}

			var resp shortcutView
			if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
				t.Fatalf("decode response: %v", err)
			}
			if resp.Alias != tt.wantAlias {
				t.Fatalf("alias = %q, want %q", resp.Alias, tt.wantAlias)
			}
		})
	}
}

func TestShortcutAPI_Index(t *testing.T) {
	tests := []struct {
		name       string
		store      *stubStore
		wantStatus int
		wantLen    int
		wantErrMsg string
	}{
		{
			name: "success",
			store: &stubStore{
				allFn: func(context.Context) ([]model.Shortcut, error) {
					return []model.Shortcut{
						{ID: 1, Alias: "a", Destination: "https://a.example"},
						{ID: 2, Alias: "b", Destination: "https://b.example"},
					}, nil
				},
			},
			wantStatus: http.StatusOK,
			wantLen:    2,
		},
		{
			name: "store error",
			store: &stubStore{
				allFn: func(context.Context) ([]model.Shortcut, error) {
					return nil, errors.New("db error")
				},
			},
			wantStatus: http.StatusInternalServerError,
			wantErrMsg: "could not load shortcuts",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := newTestAPI(tt.store)
			engine := gin.New()
			engine.GET("/api/shortcuts", api.Index)

			req := httptest.NewRequest(http.MethodGet, "/api/shortcuts", nil)
			w := httptest.NewRecorder()
			engine.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", w.Code, tt.wantStatus)
			}
			if tt.wantErrMsg != "" {
				var resp map[string]string
				_ = json.Unmarshal(w.Body.Bytes(), &resp)
				if resp["error"] != tt.wantErrMsg {
					t.Fatalf("error = %q, want %q", resp["error"], tt.wantErrMsg)
				}
				return
			}
			var resp []shortcutView
			_ = json.Unmarshal(w.Body.Bytes(), &resp)
			if len(resp) != tt.wantLen {
				t.Fatalf("len = %d, want %d", len(resp), tt.wantLen)
			}
		})
	}
}

func TestShortcutAPI_Follow(t *testing.T) {
	tests := []struct {
		name       string
		alias      string
		store      *stubStore
		wantStatus int
		wantLoc    string
		wantErrMsg string
	}{
		{
			name:  "found",
			alias: "docs",
			store: &stubStore{
				getByAliasFn: func(context.Context, string) (*model.Shortcut, error) {
					return &model.Shortcut{Alias: "docs", Destination: "https://example.com/docs"}, nil
				},
			},
			wantStatus: http.StatusFound,
			wantLoc:    "https://example.com/docs",
		},
		{
			name:  "missing",
			alias: "missing",
			store: &stubStore{
				getByAliasFn: func(context.Context, string) (*model.Shortcut, error) {
					return nil, repository.ErrMissing
				},
			},
			wantStatus: http.StatusNotFound,
			wantErrMsg: "alias not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := newTestAPI(tt.store)
			engine := gin.New()
			engine.GET("/j/:alias", api.Follow)

			req := httptest.NewRequest(http.MethodGet, "/j/"+tt.alias, nil)
			w := httptest.NewRecorder()
			engine.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", w.Code, tt.wantStatus)
			}
			if tt.wantLoc != "" {
				if loc := w.Header().Get("Location"); loc != tt.wantLoc {
					t.Fatalf("Location = %q, want %q", loc, tt.wantLoc)
				}
				return
			}
			var resp map[string]string
			_ = json.Unmarshal(w.Body.Bytes(), &resp)
			if resp["error"] != tt.wantErrMsg {
				t.Fatalf("error = %q, want %q", resp["error"], tt.wantErrMsg)
			}
		})
	}
}

func TestMapShortcut(t *testing.T) {
	createdAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	got := mapShortcut(model.Shortcut{
		ID: 42, Alias: "home", Destination: "https://example.com", CreatedAt: createdAt,
	})
	want := shortcutView{ID: 42, Alias: "home", Destination: "https://example.com", RegisteredAt: createdAt}
	if got != want {
		t.Fatalf("mapShortcut() = %#v, want %#v", got, want)
	}
}

func TestRespondWithFailure(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantMsg    string
	}{
		{name: "input", err: &validation.InputError{Detail: "alias is required"}, wantStatus: http.StatusBadRequest, wantMsg: "alias is required"},
		{name: "duplicate", err: repository.ErrDuplicateAlias, wantStatus: http.StatusConflict, wantMsg: "alias already taken"},
		{name: "missing", err: repository.ErrMissing, wantStatus: http.StatusNotFound, wantMsg: "alias not found"},
		{name: "unknown", err: errors.New("surprise"), wantStatus: http.StatusInternalServerError, wantMsg: "unexpected server error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			respondWithFailure(c, tt.err)
			if w.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", w.Code, tt.wantStatus)
			}
			var resp map[string]string
			_ = json.Unmarshal(w.Body.Bytes(), &resp)
			if resp["error"] != tt.wantMsg {
				t.Fatalf("error = %q, want %q", resp["error"], tt.wantMsg)
			}
		})
	}
}
