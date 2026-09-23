package router

import (
	"github.com/chennupati97/go-links/internal/handler"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// BuildEngine constructs the Gin engine with middleware and routes.
func BuildEngine(api *handler.ShortcutAPI) *gin.Engine {
	engine := gin.Default()
	applyMiddleware(engine)
	mountRoutes(engine, api)
	return engine
}

func applyMiddleware(engine *gin.Engine) {
	engine.Use(cors.Default())
}

func mountRoutes(engine *gin.Engine, api *handler.ShortcutAPI) {
	engine.GET("/ready", func(c *gin.Context) {
		c.JSON(200, gin.H{"ready": true})
	})

	engine.POST("/api/shortcuts", api.Register)
	engine.GET("/api/shortcuts", api.Index)
	engine.GET("/j/:alias", api.Follow)
}
