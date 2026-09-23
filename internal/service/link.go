package service

import (
	"context"

	"github.com/Naveen-kumar525/go-links/internal/model"
	"github.com/Naveen-kumar525/go-links/internal/repository"
	"github.com/Naveen-kumar525/go-links/internal/validation"
)

type CreateLinkInput struct {
	Slug string
	URL  string
}

// LinkService contains business rules for go links.
type LinkService struct {
	repo repository.LinkRepository
}

func NewLinkService(repo repository.LinkRepository) *LinkService {
	return &LinkService{repo: repo}
}

func (s *LinkService) Create(ctx context.Context, in CreateLinkInput) (*model.Link, error) {
	in.Slug = validation.NormalizeSlug(in.Slug)

	if err := validation.ValidateSlug(in.Slug); err != nil {
		return nil, err
	}
	if err := validation.ValidateURL(in.URL); err != nil {
		return nil, err
	}

	link := &model.Link{
		Slug: in.Slug,
		URL:  in.URL,
	}
	if err := s.repo.Create(ctx, link); err != nil {
		return nil, err
	}
	return link, nil
}

func (s *LinkService) List(ctx context.Context) ([]model.Link, error) {
	return s.repo.List(ctx)
}

func (s *LinkService) Resolve(ctx context.Context, slug string) (*model.Link, error) {
	slug = validation.NormalizeSlug(slug)
	if slug == "" {
		return nil, repository.ErrNotFound
	}
	return s.repo.FindBySlug(ctx, slug)
}
