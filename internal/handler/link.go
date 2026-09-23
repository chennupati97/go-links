package handler

import (
	"errors"
	"net/http"

	"github.com/Naveen-kumar525/go-links/internal/repository"
	"github.com/Naveen-kumar525/go-links/internal/service"
	"github.com/Naveen-kumar525/go-links/internal/validation"
	"github.com/gin-gonic/gin"
)

func (h *LinkHandler) CreateLink(c *gin.Context) {
	var req createLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	link, err := h.links.Create(c.Request.Context(), service.CreateLinkInput{
		Slug: req.Slug,
		URL:  req.URL,
	})
	if err != nil {
		writeServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, toLinkResponse(*link))
}

func (h *LinkHandler) ListLinks(c *gin.Context) {
	links, err := h.links.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch links"})
		return
	}

	response := make([]linkResponse, 0, len(links))
	for _, link := range links {
		response = append(response, toLinkResponse(link))
	}
	c.JSON(http.StatusOK, response)
}

func (h *LinkHandler) Redirect(c *gin.Context) {
	link, err := h.links.Resolve(c.Request.Context(), c.Param("slug"))
	if err != nil {
		writeServiceError(c, err)
		return
	}
	c.Redirect(http.StatusFound, link.URL)
}

func writeServiceError(c *gin.Context, err error) {
	if _, ok := validation.AsError(err); ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if errors.Is(err, repository.ErrConflict) {
		c.JSON(http.StatusConflict, gin.H{"error": "slug already exists"})
		return
	}
	if errors.Is(err, repository.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "shortcut not found"})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
}
