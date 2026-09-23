package repository

import (
	"context"

	"github.com/chennupati97/go-links/internal/model"
)

// ShortcutStore is the persistence port for shortcuts.
type ShortcutStore interface {
	Save(ctx context.Context, item *model.Shortcut) error
	GetByAlias(ctx context.Context, alias string) (*model.Shortcut, error)
	All(ctx context.Context) ([]model.Shortcut, error)
}
