package handler

import (
	"time"

	"github.com/Naveen-kumar525/go-links/internal/model"
	"github.com/Naveen-kumar525/go-links/internal/service"
)

type LinkHandler struct {
	links *service.LinkService
}

func NewLinkHandler(links *service.LinkService) *LinkHandler {
	return &LinkHandler{links: links}
}

type createLinkRequest struct {
	Slug string `json:"slug"`
	URL  string `json:"url"`
}

type linkResponse struct {
	ID        uint      `json:"id"`
	Slug      string    `json:"slug"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"createdAt"`
}

func toLinkResponse(link model.Link) linkResponse {
	return linkResponse{
		ID:        link.ID,
		Slug:      link.Slug,
		URL:       link.URL,
		CreatedAt: link.CreatedAt,
	}
}
