package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Naveen-kumar525/go-links/internal/model"
	"github.com/Naveen-kumar525/go-links/internal/repository"
	"github.com/Naveen-kumar525/go-links/internal/validation"
)

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

func TestLinkService_Create(t *testing.T) {
	fixed := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)

	tests := []struct {
		name      string
		input     CreateLinkInput
		repo      *mockLinkRepo
		wantSlug  string
		wantURL   string
		wantErr   error
		wantVErr  string
		checkRepo bool
	}{
		{
			name:  "success normalizes slug",
			input: CreateLinkInput{Slug: "  Docs  ", URL: "https://example.com/docs"},
			repo: &mockLinkRepo{
				createFn: func(_ context.Context, link *model.Link) error {
					link.ID = 1
					link.CreatedAt = fixed
					return nil
				},
			},
			wantSlug:  "docs",
			wantURL:   "https://example.com/docs",
			checkRepo: true,
		},
		{
			name:     "invalid slug",
			input:    CreateLinkInput{Slug: "Bad_Slug", URL: "https://example.com"},
			repo:     &mockLinkRepo{},
			wantVErr: "slug may contain only lowercase letters, numbers and hyphens",
		},
		{
			name:     "empty slug",
			input:    CreateLinkInput{Slug: "   ", URL: "https://example.com"},
			repo:     &mockLinkRepo{},
			wantVErr: "slug is required",
		},
		{
			name:     "invalid url",
			input:    CreateLinkInput{Slug: "docs", URL: "://bad"},
			repo:     &mockLinkRepo{},
			wantVErr: "invalid url",
		},
		{
			name:     "unsupported url scheme",
			input:    CreateLinkInput{Slug: "docs", URL: "ftp://example.com"},
			repo:     &mockLinkRepo{},
			wantVErr: "url must begin with http or https",
		},
		{
			name:  "repo conflict",
			input: CreateLinkInput{Slug: "docs", URL: "https://example.com"},
			repo: &mockLinkRepo{
				createFn: func(context.Context, *model.Link) error {
					return repository.ErrConflict
				},
			},
			wantErr: repository.ErrConflict,
		},
		{
			name:  "repo unexpected error",
			input: CreateLinkInput{Slug: "docs", URL: "https://example.com"},
			repo: &mockLinkRepo{
				createFn: func(context.Context, *model.Link) error {
					return errors.New("db down")
				},
			},
			wantErr: errors.New("db down"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewLinkService(tt.repo)
			got, err := svc.Create(context.Background(), tt.input)

			if tt.wantVErr != "" {
				vErr, ok := validation.AsError(err)
				if !ok {
					t.Fatalf("expected validation error, got %v", err)
				}
				if vErr.Message != tt.wantVErr {
					t.Fatalf("validation error = %q, want %q", vErr.Message, tt.wantVErr)
				}
				if got != nil {
					t.Fatalf("expected nil link, got %#v", got)
				}
				return
			}

			if tt.wantErr != nil {
				if err == nil || err.Error() != tt.wantErr.Error() {
					t.Fatalf("error = %v, want %v", err, tt.wantErr)
				}
				if got != nil {
					t.Fatalf("expected nil link, got %#v", got)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.Slug != tt.wantSlug {
				t.Fatalf("Slug = %q, want %q", got.Slug, tt.wantSlug)
			}
			if got.URL != tt.wantURL {
				t.Fatalf("URL = %q, want %q", got.URL, tt.wantURL)
			}
			if tt.checkRepo && got.ID != 1 {
				t.Fatalf("ID = %d, want 1", got.ID)
			}
		})
	}
}

func TestLinkService_List(t *testing.T) {
	sample := []model.Link{
		{ID: 1, Slug: "a", URL: "https://a.example"},
		{ID: 2, Slug: "b", URL: "https://b.example"},
	}

	tests := []struct {
		name    string
		repo    *mockLinkRepo
		want    []model.Link
		wantErr error
	}{
		{
			name: "success",
			repo: &mockLinkRepo{
				listFn: func(context.Context) ([]model.Link, error) {
					return sample, nil
				},
			},
			want: sample,
		},
		{
			name: "empty",
			repo: &mockLinkRepo{
				listFn: func(context.Context) ([]model.Link, error) {
					return []model.Link{}, nil
				},
			},
			want: []model.Link{},
		},
		{
			name: "repo error",
			repo: &mockLinkRepo{
				listFn: func(context.Context) ([]model.Link, error) {
					return nil, errors.New("list failed")
				},
			},
			wantErr: errors.New("list failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewLinkService(tt.repo)
			got, err := svc.List(context.Background())

			if tt.wantErr != nil {
				if err == nil || err.Error() != tt.wantErr.Error() {
					t.Fatalf("error = %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(got) != len(tt.want) {
				t.Fatalf("len = %d, want %d", len(got), len(tt.want))
			}
			for i := range tt.want {
				if got[i].Slug != tt.want[i].Slug || got[i].URL != tt.want[i].URL {
					t.Fatalf("got[%d] = %#v, want %#v", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestLinkService_Resolve(t *testing.T) {
	link := &model.Link{ID: 7, Slug: "docs", URL: "https://example.com/docs"}

	tests := []struct {
		name     string
		slug     string
		repo     *mockLinkRepo
		wantSlug string
		wantErr  error
	}{
		{
			name: "success normalizes slug",
			slug: "  Docs  ",
			repo: &mockLinkRepo{
				findBySlugFn: func(_ context.Context, slug string) (*model.Link, error) {
					if slug != "docs" {
						t.Fatalf("FindBySlug called with %q, want %q", slug, "docs")
					}
					return link, nil
				},
			},
			wantSlug: "docs",
		},
		{
			name:    "empty slug",
			slug:    "   ",
			repo:    &mockLinkRepo{},
			wantErr: repository.ErrNotFound,
		},
		{
			name: "not found",
			slug: "missing",
			repo: &mockLinkRepo{
				findBySlugFn: func(context.Context, string) (*model.Link, error) {
					return nil, repository.ErrNotFound
				},
			},
			wantErr: repository.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewLinkService(tt.repo)
			got, err := svc.Resolve(context.Background(), tt.slug)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("error = %v, want %v", err, tt.wantErr)
				}
				if got != nil {
					t.Fatalf("expected nil link, got %#v", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.Slug != tt.wantSlug {
				t.Fatalf("Slug = %q, want %q", got.Slug, tt.wantSlug)
			}
		})
	}
}
