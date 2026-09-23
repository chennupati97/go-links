package repository

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	"github.com/chennupati97/go-links/internal/model"
)

func newTestStore(t *testing.T) *SQLiteShortcutStore {
	t.Helper()
	conn, err := ConnectSQLite(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("ConnectSQLite: %v", err)
	}
	return NewSQLiteShortcutStore(conn)
}

func TestSQLiteShortcutStore_Save(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(t *testing.T, store *SQLiteShortcutStore)
		item    *model.Shortcut
		wantErr error
	}{
		{
			name: "success",
			item: &model.Shortcut{Alias: "docs", Destination: "https://example.com/docs"},
		},
		{
			name: "duplicate alias",
			setup: func(t *testing.T, store *SQLiteShortcutStore) {
				t.Helper()
				if err := store.Save(context.Background(), &model.Shortcut{
					Alias:       "docs",
					Destination: "https://example.com/one",
				}); err != nil {
					t.Fatalf("setup save: %v", err)
				}
			},
			item:    &model.Shortcut{Alias: "docs", Destination: "https://example.com/two"},
			wantErr: ErrDuplicateAlias,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := newTestStore(t)
			if tt.setup != nil {
				tt.setup(t, store)
			}

			err := store.Save(context.Background(), tt.item)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("error = %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.item.ID == 0 {
				t.Fatal("expected ID to be set")
			}
		})
	}
}

func TestSQLiteShortcutStore_GetByAlias(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(t *testing.T, store *SQLiteShortcutStore)
		alias     string
		wantAlias string
		wantErr   error
	}{
		{
			name: "found",
			setup: func(t *testing.T, store *SQLiteShortcutStore) {
				t.Helper()
				if err := store.Save(context.Background(), &model.Shortcut{
					Alias:       "docs",
					Destination: "https://example.com/docs",
				}); err != nil {
					t.Fatalf("setup save: %v", err)
				}
			},
			alias:     "docs",
			wantAlias: "docs",
		},
		{
			name:    "missing",
			alias:   "missing",
			wantErr: ErrMissing,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := newTestStore(t)
			if tt.setup != nil {
				tt.setup(t, store)
			}

			got, err := store.GetByAlias(context.Background(), tt.alias)
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

func TestSQLiteShortcutStore_All(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(t *testing.T, store *SQLiteShortcutStore)
		wantAliases []string
	}{
		{
			name:        "empty",
			wantAliases: nil,
		},
		{
			name: "ordered by alias ascending",
			setup: func(t *testing.T, store *SQLiteShortcutStore) {
				t.Helper()
				for _, item := range []model.Shortcut{
					{Alias: "zeta", Destination: "https://z.example"},
					{Alias: "alpha", Destination: "https://a.example"},
					{Alias: "beta", Destination: "https://b.example"},
				} {
					copyItem := item
					if err := store.Save(context.Background(), &copyItem); err != nil {
						t.Fatalf("setup save %q: %v", item.Alias, err)
					}
				}
			},
			wantAliases: []string{"alpha", "beta", "zeta"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := newTestStore(t)
			if tt.setup != nil {
				tt.setup(t, store)
			}

			got, err := store.All(context.Background())
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(got) != len(tt.wantAliases) {
				t.Fatalf("len = %d, want %d", len(got), len(tt.wantAliases))
			}
			for i, alias := range tt.wantAliases {
				if got[i].Alias != alias {
					t.Fatalf("got[%d].Alias = %q, want %q", i, got[i].Alias, alias)
				}
			}
		})
	}
}

func TestIsDuplicateKey(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{name: "unique constraint", err: errors.New("UNIQUE constraint failed: shortcuts.alias"), want: true},
		{name: "duplicate", err: errors.New("duplicate key value"), want: true},
		{name: "unrelated", err: errors.New("connection refused"), want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isDuplicateKey(tt.err)
			if got != tt.want {
				t.Fatalf("isDuplicateKey(%v) = %v, want %v", tt.err, got, tt.want)
			}
		})
	}
}

func TestConnectSQLite(t *testing.T) {
	tests := []struct {
		name string
		path string
	}{
		{name: "custom path", path: filepath.Join(t.TempDir(), "custom.db")},
		{name: "empty path uses default in cwd", path: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.path
			if path == "" {
				t.Chdir(t.TempDir())
			}
			conn, err := ConnectSQLite(path)
			if err != nil {
				t.Fatalf("ConnectSQLite: %v", err)
			}
			sqlDB, err := conn.DB()
			if err != nil {
				t.Fatalf("conn.DB: %v", err)
			}
			_ = sqlDB.Close()
		})
	}
}
