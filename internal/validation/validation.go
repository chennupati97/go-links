package validation

import (
	"errors"
	"net/url"
	"regexp"
	"strings"
)

var slugRegex = regexp.MustCompile(`^[a-z0-9-]{1,50}$`)

// Error represents a client input validation failure.
type Error struct {
	Message string
}

func (e *Error) Error() string {
	return e.Message
}

func NormalizeSlug(slug string) string {
	return strings.TrimSpace(strings.ToLower(slug))
}

func ValidateSlug(slug string) error {
	if slug == "" {
		return &Error{Message: "slug is required"}
	}
	if !slugRegex.MatchString(slug) {
		return &Error{Message: "slug may contain only lowercase letters, numbers and hyphens"}
	}
	return nil
}

func ValidateURL(raw string) error {
	if raw == "" {
		return &Error{Message: "url is required"}
	}

	u, err := url.ParseRequestURI(raw)
	if err != nil {
		return &Error{Message: "invalid url"}
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return &Error{Message: "url must begin with http or https"}
	}
	return nil
}

func AsError(err error) (*Error, bool) {
	var v *Error
	if errors.As(err, &v) {
		return v, true
	}
	return nil, false
}
