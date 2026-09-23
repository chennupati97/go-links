package validation

import (
	"errors"
	"testing"
)

func TestNormalizeSlug(t *testing.T) {
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
			got := NormalizeSlug(tt.in)
			if got != tt.want {
				t.Fatalf("NormalizeSlug(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestValidateSlug(t *testing.T) {
	tests := []struct {
		name    string
		slug    string
		wantErr string
	}{
		{name: "valid simple", slug: "docs"},
		{name: "valid with hyphens", slug: "my-go-link"},
		{name: "valid with numbers", slug: "link123"},
		{name: "valid single char", slug: "a"},
		{name: "valid max length", slug: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}, // 50
		{name: "empty", slug: "", wantErr: "slug is required"},
		{name: "uppercase", slug: "Docs", wantErr: "slug may contain only lowercase letters, numbers and hyphens"},
		{name: "underscore", slug: "my_link", wantErr: "slug may contain only lowercase letters, numbers and hyphens"},
		{name: "space", slug: "my link", wantErr: "slug may contain only lowercase letters, numbers and hyphens"},
		{name: "special chars", slug: "link!", wantErr: "slug may contain only lowercase letters, numbers and hyphens"},
		{name: "too long", slug: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", wantErr: "slug may contain only lowercase letters, numbers and hyphens"}, // 51
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSlug(tt.slug)
			assertValidationError(t, err, tt.wantErr)
		})
	}
}

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr string
	}{
		{name: "valid https", url: "https://example.com"},
		{name: "valid http", url: "http://example.com/path"},
		{name: "valid with query", url: "https://example.com/search?q=go"},
		{name: "empty", url: "", wantErr: "url is required"},
		{name: "no scheme", url: "example.com", wantErr: "invalid url"},
		{name: "ftp scheme", url: "ftp://example.com", wantErr: "url must begin with http or https"},
		{name: "javascript scheme", url: "javascript:alert(1)", wantErr: "url must begin with http or https"},
		{name: "relative path", url: "/relative", wantErr: "url must begin with http or https"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateURL(tt.url)
			assertValidationError(t, err, tt.wantErr)
		})
	}
}

func TestAsError(t *testing.T) {
	tests := []struct {
		name   string
		err    error
		wantOK bool
		wantMsg string
	}{
		{
			name:    "validation error",
			err:     &Error{Message: "slug is required"},
			wantOK:  true,
			wantMsg: "slug is required",
		},
		{
			name:   "wrapped validation error",
			err:    errors.Join(&Error{Message: "invalid url"}),
			wantOK: true,
			wantMsg: "invalid url",
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
			got, ok := AsError(tt.err)
			if ok != tt.wantOK {
				t.Fatalf("AsError() ok = %v, want %v", ok, tt.wantOK)
			}
			if !tt.wantOK {
				if got != nil {
					t.Fatalf("AsError() got = %#v, want nil", got)
				}
				return
			}
			if got == nil || got.Message != tt.wantMsg {
				t.Fatalf("AsError() = %#v, want Message %q", got, tt.wantMsg)
			}
		})
	}
}

func assertValidationError(t *testing.T, err error, wantMsg string) {
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
	vErr, ok := AsError(err)
	if !ok {
		t.Fatalf("expected *validation.Error, got %T: %v", err, err)
	}
	if vErr.Message != wantMsg {
		t.Fatalf("error message = %q, want %q", vErr.Message, wantMsg)
	}
}
