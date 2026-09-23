package repository

import "errors"

var (
	ErrNotFound = errors.New("link not found")
	ErrConflict = errors.New("slug already exists")
)
