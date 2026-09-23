package service

import (
	"context"

	"github.com/chennupati97/go-links/internal/model"
	"github.com/chennupati97/go-links/internal/repository"
	"github.com/chennupati97/go-links/internal/validation"
)

type RegisterInput struct {
	Alias       string
	Destination string
}

// AliasManager applies business rules for jump aliases.
type AliasManager struct {
	store repository.ShortcutStore
}

func NewAliasManager(store repository.ShortcutStore) *AliasManager {
	return &AliasManager{store: store}
}

func (m *AliasManager) Register(ctx context.Context, in RegisterInput) (*model.Shortcut, error) {
	in.Alias = validation.CleanAlias(in.Alias)

	if err := validation.CheckAlias(in.Alias); err != nil {
		return nil, err
	}
	if err := validation.CheckTarget(in.Destination); err != nil {
		return nil, err
	}

	item := &model.Shortcut{
		Alias:       in.Alias,
		Destination: in.Destination,
	}
	if err := m.store.Save(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}

func (m *AliasManager) All(ctx context.Context) ([]model.Shortcut, error) {
	return m.store.All(ctx)
}

func (m *AliasManager) Lookup(ctx context.Context, alias string) (*model.Shortcut, error) {
	alias = validation.CleanAlias(alias)
	if alias == "" {
		return nil, repository.ErrMissing
	}
	return m.store.GetByAlias(ctx, alias)
}
