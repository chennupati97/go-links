package di

import (
	"testing"
)

func TestEnvOr(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		set      string
		fallback string
		want     string
	}{
		{name: "uses env when set", key: "GO_LINKS_TEST_ENV_OR", set: "from-env", fallback: "fallback", want: "from-env"},
		{name: "uses fallback when unset", key: "GO_LINKS_TEST_ENV_OR_MISSING", set: "", fallback: "fallback", want: "fallback"},
		{name: "uses fallback when empty", key: "GO_LINKS_TEST_ENV_OR_EMPTY", set: "", fallback: "fallback", want: "fallback"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv(tt.key, tt.set)
			if tt.set == "" {
				t.Setenv(tt.key, "")
			}
			got := envOr(tt.key, tt.fallback)
			if got != tt.want {
				t.Fatalf("envOr(%q, %q) = %q, want %q", tt.key, tt.fallback, got, tt.want)
			}
		})
	}
}

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name       string
		port       *string
		dbPath     *string
		wantAddr   string
		wantDBPath string
	}{
		{
			name:       "defaults",
			wantAddr:   ":8080",
			wantDBPath: "golinks.db",
		},
		{
			name:       "port without colon",
			port:       strPtr("9090"),
			wantAddr:   ":9090",
			wantDBPath: "golinks.db",
		},
		{
			name:       "port with colon",
			port:       strPtr(":7070"),
			wantAddr:   ":7070",
			wantDBPath: "golinks.db",
		},
		{
			name:       "custom db path",
			dbPath:     strPtr("/tmp/custom.db"),
			wantAddr:   ":8080",
			wantDBPath: "/tmp/custom.db",
		},
		{
			name:       "port and db path",
			port:       strPtr("3000"),
			dbPath:     strPtr("data/app.db"),
			wantAddr:   ":3000",
			wantDBPath: "data/app.db",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.port != nil {
				t.Setenv("PORT", *tt.port)
			} else {
				t.Setenv("PORT", "")
			}
			if tt.dbPath != nil {
				t.Setenv("DB_PATH", *tt.dbPath)
			} else {
				t.Setenv("DB_PATH", "")
			}

			cfg := LoadConfig()
			if cfg.Addr != tt.wantAddr {
				t.Fatalf("Addr = %q, want %q", cfg.Addr, tt.wantAddr)
			}
			if cfg.DBPath != tt.wantDBPath {
				t.Fatalf("DBPath = %q, want %q", cfg.DBPath, tt.wantDBPath)
			}
		})
	}
}

func strPtr(s string) *string {
	return &s
}
