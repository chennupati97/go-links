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

	"github.com/Naveen-kumar525/go-links/internal/model"
	"github.com/Naveen-kumar525/go-links/internal/repository"
	"github.com/Naveen-kumar525/go-links/internal/service"
	"github.com/Naveen-kumar525/go-links/internal/validation"
	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

type mockLinkRepo struct {
	createFn     func(ctx context.Context, link *model.Link) error
	findBySlugFn func(ctx context.Context, slug string) (*model.Link, error)
	listFn       func(ctx context.Context) ([]model.Link, error)
}

func (m *mockLinkRepo) Create(ctx context.Context, link *model.Link) error {
	if m.createFn != nil {
		return m.createFn(ctx, link)
	}
	return nil
}

func (m *mockLinkRepo) FindBySlug(ctx context.Context, slug string) (*model.Link, error) {
	if m.findBySlugFn != nil {
		return m.findBySlugFn(ctx, slug)
	}
	return nil, repository.ErrNotFound
}

func (m *mockLinkRepo) List(ctx context.Context) ([]model.Link, error) {
	if m.listFn != nil {
		return m.listFn(ctx)
	}
	return nil, nil
}

func newTestHandler(repo *mockLinkRepo) *LinkHandler {
	return NewLinkHandler(service.NewLinkService(repo))
}

func TestLinkHandler_CreateLink(t *testing.T) {
	createdAt := time.Date(2024, 5, 1, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name       string
		body       string
		repo       *mockLinkRepo
		wantStatus int
		wantSlug   string
		wantErrMsg string
	}{
		{
			name: "created",
			body: `{"slug":"docs","url":"https://example.com/docs"}`,
			repo: &mockLinkRepo{
				createFn: func(_ context.Context, link *model.Link) error {
					link.ID = 10
					link.CreatedAt = createdAt
					return nil
				},
			},
			wantStatus: http.StatusCreated,
			wantSlug:   "docs",
		},
		{
			name:       "invalid json",
			body:       `{not-json`,
			repo:       &mockLinkRepo{},
			wantStatus: http.StatusBadRequest,
			wantErrMsg: "invalid json",
		},
		{
			name:       "validation error",
			body:       `{"slug":"bad_slug","url":"https://example.com"}`,
			repo:       &mockLinkRepo{},
			wantStatus: http.StatusBadRequest,
			wantErrMsg: "slug may contain only lowercase letters, numbers and hyphens",
		},
		{
			name: "conflict",
			body: `{"slug":"docs","url":"https://example.com"}`,
			repo: &mockLinkRepo{
				createFn: func(context.Context, *model.Link) error {
					return repository.ErrConflict
				},
			},
			wantStatus: http.StatusConflict,
			wantErrMsg: "slug already exists",
		},
		{
			name: "internal error",
			body: `{"slug":"docs","url":"https://example.com"}`,
			repo: &mockLinkRepo{
				createFn: func(context.Context, *model.Link) error {
					return errors.New("boom")
				},
			},
			wantStatus: http.StatusInternalServerError,
			wantErrMsg: "internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newTestHandler(tt.repo)
			r := gin.New()
			r.POST("/api/links", h.CreateLink)

			req := httptest.NewRequest(http.MethodPost, "/api/links", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

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

			var resp linkResponse
			if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
				t.Fatalf("decode response: %v", err)
			}
			if resp.Slug != tt.wantSlug {
				t.Fatalf("slug = %q, want %q", resp.Slug, tt.wantSlug)
			}
			if resp.ID != 10 {
				t.Fatalf("id = %d, want 10", resp.ID)
			}
		})
	}
}

func TestLinkHandler_ListLinks(t *testing.T) {
	tests := []struct {
		name       string
		repo       *mockLinkRepo
		wantStatus int
		wantLen    int
		wantErrMsg string
	}{
		{
			name: "success",
			repo: &mockLinkRepo{
				listFn: func(context.Context) ([]model.Link, error) {
					return []model.Link{
						{ID: 1, Slug: "a", URL: "https://a.example"},
						{ID: 2, Slug: "b", URL: "https://b.example"},
					}, nil
				},
			},
			wantStatus: http.StatusOK,
			wantLen:    2,
		},
		{
			name: "empty list",
			repo: &mockLinkRepo{
				listFn: func(context.Context) ([]model.Link, error) {
					return []model.Link{}, nil
				},
			},
			wantStatus: http.StatusOK,
			wantLen:    0,
		},
		{
			name: "repo error",
			repo: &mockLinkRepo{
				listFn: func(context.Context) ([]model.Link, error) {
					return nil, errors.New("db error")
				},
			},
			wantStatus: http.StatusInternalServerError,
			wantErrMsg: "failed to fetch links",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newTestHandler(tt.repo)
			r := gin.New()
			r.GET("/api/links", h.ListLinks)

			req := httptest.NewRequest(http.MethodGet, "/api/links", nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

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

			var resp []linkResponse
			if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
				t.Fatalf("decode response: %v", err)
			}
			if len(resp) != tt.wantLen {
				t.Fatalf("len = %d, want %d", len(resp), tt.wantLen)
			}
		})
	}
}

func TestLinkHandler_Redirect(t *testing.T) {
	tests := []struct {
		name       string
		slug       string
		repo       *mockLinkRepo
		wantStatus int
		wantLoc    string
		wantErrMsg string
	}{
		{
			name: "found",
			slug: "docs",
			repo: &mockLinkRepo{
				findBySlugFn: func(_ context.Context, slug string) (*model.Link, error) {
					if slug != "docs" {
						t.Fatalf("slug = %q", slug)
					}
					return &model.Link{Slug: "docs", URL: "https://example.com/docs"}, nil
				},
			},
			wantStatus: http.StatusFound,
			wantLoc:    "https://example.com/docs",
		},
		{
			name: "not found",
			slug: "missing",
			repo: &mockLinkRepo{
				findBySlugFn: func(context.Context, string) (*model.Link, error) {
					return nil, repository.ErrNotFound
				},
			},
			wantStatus: http.StatusNotFound,
			wantErrMsg: "shortcut not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newTestHandler(tt.repo)
			r := gin.New()
			r.GET("/go/:slug", h.Redirect)

			req := httptest.NewRequest(http.MethodGet, "/go/"+tt.slug, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d; body=%s", w.Code, tt.wantStatus, w.Body.String())
			}
			if tt.wantLoc != "" {
				if loc := w.Header().Get("Location"); loc != tt.wantLoc {
					t.Fatalf("Location = %q, want %q", loc, tt.wantLoc)
				}
				return
			}
			var resp map[string]string
			if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
				t.Fatalf("decode error body: %v", err)
			}
			if resp["error"] != tt.wantErrMsg {
				t.Fatalf("error = %q, want %q", resp["error"], tt.wantErrMsg)
			}
		})
	}
}

func TestToLinkResponse(t *testing.T) {
	createdAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	tests := []struct {
		name string
		in   model.Link
		want linkResponse
	}{
		{
			name: "maps all fields",
			in: model.Link{
				ID:        42,
				Slug:      "home",
				URL:       "https://example.com",
				CreatedAt: createdAt,
			},
			want: linkResponse{
				ID:        42,
				Slug:      "home",
				URL:       "https://example.com",
				CreatedAt: createdAt,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toLinkResponse(tt.in)
			if got != tt.want {
				t.Fatalf("toLinkResponse() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestWriteServiceError(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantMsg    string
	}{
		{
			name:       "validation",
			err:        &validation.Error{Message: "slug is required"},
			wantStatus: http.StatusBadRequest,
			wantMsg:    "slug is required",
		},
		{
			name:       "conflict",
			err:        repository.ErrConflict,
			wantStatus: http.StatusConflict,
			wantMsg:    "slug already exists",
		},
		{
			name:       "not found",
			err:        repository.ErrNotFound,
			wantStatus: http.StatusNotFound,
			wantMsg:    "shortcut not found",
		},
		{
			name:       "unknown",
			err:        errors.New("surprise"),
			wantStatus: http.StatusInternalServerError,
			wantMsg:    "internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			writeServiceError(c, tt.err)

			if w.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d; body=%s", w.Code, tt.wantStatus, w.Body.String())
			}
			var resp map[string]string
			if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
				t.Fatalf("decode error body: %v", err)
			}
			if resp["error"] != tt.wantMsg {
				t.Fatalf("error = %q, want %q", resp["error"], tt.wantMsg)
			}
		})
	}
}
