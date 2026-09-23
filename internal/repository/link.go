package repository

import (
	"context"

	"github.com/Naveen-kumar525/go-links/internal/model"
)

// LinkRepository defines persistence operations for links.
type LinkRepository interface {
	Create(ctx context.Context, link *model.Link) error
	FindBySlug(ctx context.Context, slug string) (*model.Link, error)
	List(ctx context.Context) ([]model.Link, error)
}
