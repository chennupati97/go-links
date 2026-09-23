package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/chennupati97/go-links/internal/model"
	"github.com/chennupati97/go-links/internal/repository"
	"github.com/chennupati97/go-links/internal/validation"
)

type stubStore struct {
	saveFn      func(ctx context.Context, item *model.Shortcut) error
	getByAliasFn func(ctx context.Context, alias string) (*model.Shortcut, error)
	allFn       func(ctx context.Context) ([]model.Shortcut, error)
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

func TestAliasManager_Register(t *testing.T) {
	fixed := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)

	tests := []struct {
		name       string
		input      RegisterInput
		store      *stubStore
		wantAlias  string
		wantDest   string
		wantErr    error
		wantInErr  string
		checkStore bool
	}{
		{
			name:  "success cleans alias",
			input: RegisterInput{Alias: "  Docs  ", Destination: "https://example.com/docs"},
			store: &stubStore{
				saveFn: func(_ context.Context, item *model.Shortcut) error {
					item.ID = 1
					item.CreatedAt = fixed
					return nil
				},
			},
			wantAlias:  "docs",
			wantDest:   "https://example.com/docs",
			checkStore: true,
		},
		{
			name:      "invalid alias",
			input:     RegisterInput{Alias: "Bad_Slug", Destination: "https://example.com"},
			store:     &stubStore{},
			wantInErr: "alias may contain only lowercase letters, numbers and hyphens",
		},
		{
			name:      "empty alias",
			input:     RegisterInput{Alias: "   ", Destination: "https://example.com"},
			store:     &stubStore{},
			wantInErr: "alias is required",
		},
		{
			name:      "invalid destination",
			input:     RegisterInput{Alias: "docs", Destination: "://bad"},
			store:     &stubStore{},
			wantInErr: "destination is not a valid URL",
		},
		{
			name:  "unsupported scheme",
			input: RegisterInput{Alias: "docs", Destination: "ftp://example.com"},
			store: &stubStore{},
			wantInErr: "destination must use http or https",
		},
		{
			name:  "duplicate alias",
			input: RegisterInput{Alias: "docs", Destination: "https://example.com"},
			store: &stubStore{
				saveFn: func(context.Context, *model.Shortcut) error {
					return repository.ErrDuplicateAlias
				},
			},
			wantErr: repository.ErrDuplicateAlias,
		},
		{
			name:  "store unexpected error",
			input: RegisterInput{Alias: "docs", Destination: "https://example.com"},
			store: &stubStore{
				saveFn: func(context.Context, *model.Shortcut) error {
					return errors.New("db down")
				},
			},
			wantErr: errors.New("db down"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mgr := NewAliasManager(tt.store)
			got, err := mgr.Register(context.Background(), tt.input)

			if tt.wantInErr != "" {
				inputErr, ok := validation.IsInputError(err)
				if !ok {
					t.Fatalf("expected input error, got %v", err)
				}
				if inputErr.Detail != tt.wantInErr {
					t.Fatalf("input error = %q, want %q", inputErr.Detail, tt.wantInErr)
				}
				return
			}

			if tt.wantErr != nil {
				if err == nil || err.Error() != tt.wantErr.Error() {
					t.Fatalf("error = %v, want %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.Alias != tt.wantAlias {
				t.Fatalf("Alias = %q, want %q", got.Alias, tt.wantAlias)
			}
			if got.Destination != tt.wantDest {
				t.Fatalf("Destination = %q, want %q", got.Destination, tt.wantDest)
			}
			if tt.checkStore && got.ID != 1 {
				t.Fatalf("ID = %d, want 1", got.ID)
			}
		})
	}
}

func TestAliasManager_All(t *testing.T) {
	sample := []model.Shortcut{
		{ID: 1, Alias: "a", Destination: "https://a.example"},
		{ID: 2, Alias: "b", Destination: "https://b.example"},
	}

	tests := []struct {
		name    string
		store   *stubStore
		want    []model.Shortcut
		wantErr error
	}{
		{
			name: "success",
			store: &stubStore{
				allFn: func(context.Context) ([]model.Shortcut, error) {
					return sample, nil
				},
			},
			want: sample,
		},
		{
			name: "empty",
			store: &stubStore{
				allFn: func(context.Context) ([]model.Shortcut, error) {
					return []model.Shortcut{}, nil
				},
			},
			want: []model.Shortcut{},
		},
		{
			name: "store error",
			store: &stubStore{
				allFn: func(context.Context) ([]model.Shortcut, error) {
					return nil, errors.New("list failed")
				},
			},
			wantErr: errors.New("list failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mgr := NewAliasManager(tt.store)
			got, err := mgr.All(context.Background())

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
		})
	}
}

func TestAliasManager_Lookup(t *testing.T) {
	item := &model.Shortcut{ID: 7, Alias: "docs", Destination: "https://example.com/docs"}

	tests := []struct {
		name      string
		alias     string
		store     *stubStore
		wantAlias string
		wantErr   error
	}{
		{
			name:  "success cleans alias",
			alias: "  Docs  ",
			store: &stubStore{
				getByAliasFn: func(_ context.Context, alias string) (*model.Shortcut, error) {
					if alias != "docs" {
						t.Fatalf("GetByAlias called with %q, want %q", alias, "docs")
					}
					return item, nil
				},
			},
			wantAlias: "docs",
		},
		{
			name:    "empty alias",
			alias:   "   ",
			store:   &stubStore{},
			wantErr: repository.ErrMissing,
		},
		{
			name:  "not found",
			alias: "missing",
			store: &stubStore{
				getByAliasFn: func(context.Context, string) (*model.Shortcut, error) {
					return nil, repository.ErrMissing
				},
			},
			wantErr: repository.ErrMissing,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mgr := NewAliasManager(tt.store)
			got, err := mgr.Lookup(context.Background(), tt.alias)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("error = %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.Alias != tt.wantAlias {
				t.Fatalf("Alias = %q, want %q", got.Alias, tt.wantAlias)
			}
		})
	}
}
