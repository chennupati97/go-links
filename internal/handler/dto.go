package handler

import (
	"time"

	"github.com/chennupati97/go-links/internal/model"
	"github.com/chennupati97/go-links/internal/service"
)

type ShortcutAPI struct {
	manager *service.AliasManager
}

func NewShortcutAPI(manager *service.AliasManager) *ShortcutAPI {
	return &ShortcutAPI{manager: manager}
}

type registerPayload struct {
	Alias       string `json:"alias"`
	Destination string `json:"destination"`
}

type shortcutView struct {
	ID           uint      `json:"id"`
	Alias        string    `json:"alias"`
	Destination  string    `json:"destination"`
	RegisteredAt time.Time `json:"registeredAt"`
}

func mapShortcut(item model.Shortcut) shortcutView {
	return shortcutView{
		ID:           item.ID,
		Alias:        item.Alias,
		Destination:  item.Destination,
		RegisteredAt: item.CreatedAt,
	}
}
