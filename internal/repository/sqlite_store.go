package repository

import (
	"context"
	"errors"
	"strings"

	"github.com/chennupati97/go-links/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SQLiteShortcutStore struct {
	conn *gorm.DB
}

func NewSQLiteShortcutStore(conn *gorm.DB) *SQLiteShortcutStore {
	return &SQLiteShortcutStore{conn: conn}
}

func ConnectSQLite(filePath string) (*gorm.DB, error) {
	if filePath == "" {
		filePath = "jumpalias.db"
	}

	conn, err := gorm.Open(sqlite.Open(filePath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := conn.AutoMigrate(&model.Shortcut{}); err != nil {
		return nil, err
	}

	return conn, nil
}

func (s *SQLiteShortcutStore) Save(ctx context.Context, item *model.Shortcut) error {
	err := s.conn.WithContext(ctx).Create(item).Error
	if err == nil {
		return nil
	}
	if isDuplicateKey(err) {
		return ErrDuplicateAlias
	}
	return err
}

func (s *SQLiteShortcutStore) GetByAlias(ctx context.Context, alias string) (*model.Shortcut, error) {
	var item model.Shortcut
	err := s.conn.WithContext(ctx).Where("alias = ?", alias).First(&item).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrMissing
	}
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *SQLiteShortcutStore) All(ctx context.Context) ([]model.Shortcut, error) {
	var items []model.Shortcut
	if err := s.conn.WithContext(ctx).Order("alias ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func isDuplicateKey(err error) bool {
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "unique constraint") ||
		strings.Contains(msg, "duplicate")
}
