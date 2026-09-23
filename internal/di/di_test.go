package di

import (
	"testing"
)

func TestLookupEnv(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		set      string
		fallback string
		want     string
	}{
		{name: "uses env when set", key: "JUMPALIAS_TEST_LOOKUP", set: "from-env", fallback: "fallback", want: "from-env"},
		{name: "uses fallback when empty", key: "JUMPALIAS_TEST_LOOKUP_EMPTY", set: "", fallback: "fallback", want: "fallback"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv(tt.key, tt.set)
			got := lookupEnv(tt.key, tt.fallback)
			if got != tt.want {
				t.Fatalf("lookupEnv(%q, %q) = %q, want %q", tt.key, tt.fallback, got, tt.want)
			}
		})
	}
}

func TestReadSettings(t *testing.T) {
	tests := []struct {
		name             string
		port             *string
		dbFile           *string
		wantListenAddr   string
		wantDatabaseFile string
	}{
		{
			name:             "defaults",
			wantListenAddr:   ":8080",
			wantDatabaseFile: "jumpalias.db",
		},
		{
			name:             "port without colon",
			port:             strPtr("9090"),
			wantListenAddr:   ":9090",
			wantDatabaseFile: "jumpalias.db",
		},
		{
			name:             "port with colon",
			port:             strPtr(":7070"),
			wantListenAddr:   ":7070",
			wantDatabaseFile: "jumpalias.db",
		},
		{
			name:             "custom sqlite file",
			dbFile:           strPtr("/tmp/custom.db"),
			wantListenAddr:   ":8080",
			wantDatabaseFile: "/tmp/custom.db",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.port != nil {
				t.Setenv("LISTEN_PORT", *tt.port)
			} else {
				t.Setenv("LISTEN_PORT", "")
			}
			if tt.dbFile != nil {
				t.Setenv("SQLITE_FILE", *tt.dbFile)
			} else {
				t.Setenv("SQLITE_FILE", "")
			}

			settings := ReadSettings()
			if settings.ListenAddr != tt.wantListenAddr {
				t.Fatalf("ListenAddr = %q, want %q", settings.ListenAddr, tt.wantListenAddr)
			}
			if settings.DatabaseFile != tt.wantDatabaseFile {
				t.Fatalf("DatabaseFile = %q, want %q", settings.DatabaseFile, tt.wantDatabaseFile)
			}
		})
	}
}

func strPtr(s string) *string {
	return &s
}
