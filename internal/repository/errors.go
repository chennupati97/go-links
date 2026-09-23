package repository

import "errors"

var (
	ErrMissing        = errors.New("shortcut not found")
	ErrDuplicateAlias = errors.New("alias already taken")
)
