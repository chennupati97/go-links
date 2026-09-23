package repository

import (
	"context"
	"errors"
	"strings"

	"github.com/Naveen-kumar525/go-links/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type GormLinkRepository struct {
	db *gorm.DB
}

func NewGormLinkRepository(db *gorm.DB) *GormLinkRepository {
	return &GormLinkRepository{db: db}
}

func OpenDB(path string) (*gorm.DB, error) {
	if path == "" {
		path = "golinks.db"
	}

	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&model.Link{}); err != nil {
		return nil, err
	}

	return db, nil
}

func (r *GormLinkRepository) Create(ctx context.Context, link *model.Link) error {
	err := r.db.WithContext(ctx).Create(link).Error
	if err == nil {
		return nil
	}
	if isUniqueViolation(err) {
		return ErrConflict
	}
	return err
}

func (r *GormLinkRepository) FindBySlug(ctx context.Context, slug string) (*model.Link, error) {
	var link model.Link
	err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&link).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &link, nil
}

func (r *GormLinkRepository) List(ctx context.Context) ([]model.Link, error) {
	var links []model.Link
	if err := r.db.WithContext(ctx).Order("slug ASC").Find(&links).Error; err != nil {
		return nil, err
	}
	return links, nil
}

func isUniqueViolation(err error) bool {
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "unique constraint") ||
		strings.Contains(msg, "duplicate")
}
