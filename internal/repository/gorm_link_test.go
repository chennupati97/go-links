package repository

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	"github.com/Naveen-kumar525/go-links/internal/model"
)

func newTestRepo(t *testing.T) *GormLinkRepository {
	t.Helper()
	db, err := OpenDB(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("OpenDB: %v", err)
	}
	return NewGormLinkRepository(db)
}

func TestGormLinkRepository_Create(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(t *testing.T, repo *GormLinkRepository)
		link    *model.Link
		wantErr error
	}{
		{
			name: "success",
			link: &model.Link{Slug: "docs", URL: "https://example.com/docs"},
		},
		{
			name: "conflict on duplicate slug",
			setup: func(t *testing.T, repo *GormLinkRepository) {
				t.Helper()
				if err := repo.Create(context.Background(), &model.Link{
					Slug: "docs",
					URL:  "https://example.com/one",
				}); err != nil {
					t.Fatalf("setup create: %v", err)
				}
			},
			link:    &model.Link{Slug: "docs", URL: "https://example.com/two"},
			wantErr: ErrConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newTestRepo(t)
			if tt.setup != nil {
				tt.setup(t, repo)
			}

			err := repo.Create(context.Background(), tt.link)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("error = %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.link.ID == 0 {
				t.Fatal("expected ID to be set")
			}
		})
	}
}

func TestGormLinkRepository_FindBySlug(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(t *testing.T, repo *GormLinkRepository)
		slug     string
		wantSlug string
		wantErr  error
	}{
		{
			name: "found",
			setup: func(t *testing.T, repo *GormLinkRepository) {
				t.Helper()
				if err := repo.Create(context.Background(), &model.Link{
					Slug: "docs",
					URL:  "https://example.com/docs",
				}); err != nil {
					t.Fatalf("setup create: %v", err)
				}
			},
			slug:     "docs",
			wantSlug: "docs",
		},
		{
			name:    "not found",
			slug:    "missing",
			wantErr: ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newTestRepo(t)
			if tt.setup != nil {
				tt.setup(t, repo)
			}

			got, err := repo.FindBySlug(context.Background(), tt.slug)
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

func TestGormLinkRepository_List(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(t *testing.T, repo *GormLinkRepository)
		wantSlugs []string
	}{
		{
			name:      "empty",
			wantSlugs: nil,
		},
		{
			name: "ordered by slug ascending",
			setup: func(t *testing.T, repo *GormLinkRepository) {
				t.Helper()
				for _, link := range []model.Link{
					{Slug: "zeta", URL: "https://z.example"},
					{Slug: "alpha", URL: "https://a.example"},
					{Slug: "beta", URL: "https://b.example"},
				} {
					l := link
					if err := repo.Create(context.Background(), &l); err != nil {
						t.Fatalf("setup create %q: %v", link.Slug, err)
					}
				}
			},
			wantSlugs: []string{"alpha", "beta", "zeta"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newTestRepo(t)
			if tt.setup != nil {
				tt.setup(t, repo)
			}

			got, err := repo.List(context.Background())
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(got) != len(tt.wantSlugs) {
				t.Fatalf("len = %d, want %d", len(got), len(tt.wantSlugs))
			}
			for i, slug := range tt.wantSlugs {
				if got[i].Slug != slug {
					t.Fatalf("got[%d].Slug = %q, want %q", i, got[i].Slug, slug)
				}
			}
		})
	}
}

func TestIsUniqueViolation(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{name: "unique constraint", err: errors.New("UNIQUE constraint failed: links.slug"), want: true},
		{name: "duplicate", err: errors.New("duplicate key value"), want: true},
		{name: "case insensitive", err: errors.New("Unique Constraint Failed"), want: true},
		{name: "unrelated", err: errors.New("connection refused"), want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isUniqueViolation(tt.err)
			if got != tt.want {
				t.Fatalf("isUniqueViolation(%v) = %v, want %v", tt.err, got, tt.want)
			}
		})
	}
}

func TestOpenDB(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name: "custom path",
			path: filepath.Join(t.TempDir(), "custom.db"),
		},
		{
			name: "empty path uses default file in cwd",
			path: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.path
			if path == "" {
				// Keep default filename isolated to temp working dir.
				t.Chdir(t.TempDir())
			}
			db, err := OpenDB(path)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("OpenDB: %v", err)
			}
			sqlDB, err := db.DB()
			if err != nil {
				t.Fatalf("db.DB: %v", err)
			}
			_ = sqlDB.Close()
		})
	}
}
