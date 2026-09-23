package router

import (
	"github.com/Naveen-kumar525/go-links/internal/handler"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// New creates the HTTP engine, applies middleware, and registers routes.
func New(links *handler.LinkHandler) *gin.Engine {
	r := gin.Default()
	configure(r)
	register(r, links)
	return r
}

func configure(r *gin.Engine) {
	r.Use(cors.Default())
}

func register(r *gin.Engine, links *handler.LinkHandler) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	r.POST("/api/links", links.CreateLink)
	r.GET("/api/links", links.ListLinks)
	r.GET("/go/:slug", links.Redirect)
}
