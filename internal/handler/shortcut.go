package handler

import (
	"errors"
	"net/http"

	"github.com/chennupati97/go-links/internal/repository"
	"github.com/chennupati97/go-links/internal/service"
	"github.com/chennupati97/go-links/internal/validation"
	"github.com/gin-gonic/gin"
)

func (api *ShortcutAPI) Register(c *gin.Context) {
	var payload registerPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "malformed request body"})
		return
	}

	item, err := api.manager.Register(c.Request.Context(), service.RegisterInput{
		Alias:       payload.Alias,
		Destination: payload.Destination,
	})
	if err != nil {
		respondWithFailure(c, err)
		return
	}

	c.JSON(http.StatusCreated, mapShortcut(*item))
}

func (api *ShortcutAPI) Index(c *gin.Context) {
	items, err := api.manager.All(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not load shortcuts"})
		return
	}

	views := make([]shortcutView, 0, len(items))
	for _, item := range items {
		views = append(views, mapShortcut(item))
	}
	c.JSON(http.StatusOK, views)
}

func (api *ShortcutAPI) Follow(c *gin.Context) {
	item, err := api.manager.Lookup(c.Request.Context(), c.Param("alias"))
	if err != nil {
		respondWithFailure(c, err)
		return
	}
	c.Redirect(http.StatusFound, item.Destination)
}

func respondWithFailure(c *gin.Context, err error) {
	if _, ok := validation.IsInputError(err); ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if errors.Is(err, repository.ErrDuplicateAlias) {
		c.JSON(http.StatusConflict, gin.H{"error": "alias already taken"})
		return
	}
	if errors.Is(err, repository.ErrMissing) {
		c.JSON(http.StatusNotFound, gin.H{"error": "alias not found"})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": "unexpected server error"})
}
