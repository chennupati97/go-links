package validation

import (
	"errors"
	"testing"
)

func TestCleanAlias(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "lowercase unchanged", in: "docs", want: "docs"},
		{name: "uppercase to lowercase", in: "DOCS", want: "docs"},
		{name: "mixed case", in: "My-Link", want: "my-link"},
		{name: "trim spaces", in: "  docs  ", want: "docs"},
		{name: "trim and lowercase", in: "  Go-Link  ", want: "go-link"},
		{name: "empty", in: "", want: ""},
		{name: "only spaces", in: "   ", want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CleanAlias(tt.in)
			if got != tt.want {
				t.Fatalf("CleanAlias(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestCheckAlias(t *testing.T) {
	tests := []struct {
		name    string
		alias   string
		wantErr string
	}{
		{name: "valid simple", alias: "docs"},
		{name: "valid with hyphens", alias: "my-go-link"},
		{name: "valid with numbers", alias: "link123"},
		{name: "valid single char", alias: "a"},
		{name: "valid max length", alias: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
		{name: "empty", alias: "", wantErr: "alias is required"},
		{name: "uppercase", alias: "Docs", wantErr: "alias may contain only lowercase letters, numbers and hyphens"},
		{name: "underscore", alias: "my_link", wantErr: "alias may contain only lowercase letters, numbers and hyphens"},
		{name: "space", alias: "my link", wantErr: "alias may contain only lowercase letters, numbers and hyphens"},
		{name: "special chars", alias: "link!", wantErr: "alias may contain only lowercase letters, numbers and hyphens"},
		{name: "too long", alias: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", wantErr: "alias may contain only lowercase letters, numbers and hyphens"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckAlias(tt.alias)
			assertInputError(t, err, tt.wantErr)
		})
	}
}

func TestCheckTarget(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr string
	}{
		{name: "valid https", url: "https://example.com"},
		{name: "valid http", url: "http://example.com/path"},
		{name: "valid with query", url: "https://example.com/search?q=go"},
		{name: "empty", url: "", wantErr: "destination is required"},
		{name: "no scheme", url: "example.com", wantErr: "destination is not a valid URL"},
		{name: "ftp scheme", url: "ftp://example.com", wantErr: "destination must use http or https"},
		{name: "javascript scheme", url: "javascript:alert(1)", wantErr: "destination must use http or https"},
		{name: "relative path", url: "/relative", wantErr: "destination must use http or https"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckTarget(tt.url)
			assertInputError(t, err, tt.wantErr)
		})
	}
}

func TestIsInputError(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		wantOK  bool
		wantMsg string
	}{
		{
			name:    "input error",
			err:     &InputError{Detail: "alias is required"},
			wantOK:  true,
			wantMsg: "alias is required",
		},
		{
			name:    "joined input error",
			err:     errors.Join(&InputError{Detail: "destination is not a valid URL"}),
			wantOK:  true,
			wantMsg: "destination is not a valid URL",
		},
		{
			name:   "plain error",
			err:    errors.New("plain"),
			wantOK: false,
		},
		{
			name:   "nil",
			err:    nil,
			wantOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := IsInputError(tt.err)
			if ok != tt.wantOK {
				t.Fatalf("IsInputError() ok = %v, want %v", ok, tt.wantOK)
			}
			if !tt.wantOK {
				if got != nil {
					t.Fatalf("IsInputError() got = %#v, want nil", got)
				}
				return
			}
			if got == nil || got.Detail != tt.wantMsg {
				t.Fatalf("IsInputError() = %#v, want Detail %q", got, tt.wantMsg)
			}
		})
	}
}

func assertInputError(t *testing.T, err error, wantMsg string) {
	t.Helper()
	if wantMsg == "" {
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		return
	}
	if err == nil {
		t.Fatalf("expected error %q, got nil", wantMsg)
	}
	inputErr, ok := IsInputError(err)
	if !ok {
		t.Fatalf("expected *InputError, got %T: %v", err, err)
	}
	if inputErr.Detail != wantMsg {
		t.Fatalf("error detail = %q, want %q", inputErr.Detail, wantMsg)
	}
}
