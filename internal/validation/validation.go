package validation

import (
	"errors"
	"net/url"
	"regexp"
	"strings"
)

var aliasPattern = regexp.MustCompile(`^[a-z0-9-]{1,50}$`)

// InputError is returned when client-provided data fails checks.
type InputError struct {
	Detail string
}

func (e *InputError) Error() string {
	return e.Detail
}

func CleanAlias(value string) string {
	return strings.TrimSpace(strings.ToLower(value))
}

func CheckAlias(value string) error {
	if value == "" {
		return &InputError{Detail: "alias is required"}
	}
	if !aliasPattern.MatchString(value) {
		return &InputError{Detail: "alias may contain only lowercase letters, numbers and hyphens"}
	}
	return nil
}

func CheckTarget(raw string) error {
	if raw == "" {
		return &InputError{Detail: "destination is required"}
	}

	parsed, err := url.ParseRequestURI(raw)
	if err != nil {
		return &InputError{Detail: "destination is not a valid URL"}
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return &InputError{Detail: "destination must use http or https"}
	}
	return nil
}

func IsInputError(err error) (*InputError, bool) {
	var inputErr *InputError
	if errors.As(err, &inputErr) {
		return inputErr, true
	}
	return nil, false
}
