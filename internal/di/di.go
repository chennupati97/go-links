package di

import (
	"fmt"
	"os"

	"github.com/Naveen-kumar525/go-links/internal/handler"
	"github.com/Naveen-kumar525/go-links/internal/repository"
	"github.com/Naveen-kumar525/go-links/internal/router"
	"github.com/Naveen-kumar525/go-links/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Config holds process-level settings loaded once by the DI container.
type Config struct {
	Addr   string
	DBPath string
}

func LoadConfig() Config {
	addr := envOr("PORT", "8080")
	if len(addr) > 0 && addr[0] != ':' {
		addr = ":" + addr
	}

	return Config{
		Addr:   addr,
		DBPath: envOr("DB_PATH", "golinks.db"),
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// Container is the single composition root that constructs and injects dependencies.
type Container struct {
	Config Config
	DB     *gorm.DB

	LinkRepo    repository.LinkRepository
	LinkService *service.LinkService
	LinkHandler *handler.LinkHandler

	Router *gin.Engine
}

// New builds the dependency graph and HTTP router.
func New() (*Container, error) {
	cfg := LoadConfig()

	db, err := repository.OpenDB(cfg.DBPath)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	linkRepo := repository.NewGormLinkRepository(db)
	linkService := service.NewLinkService(linkRepo)
	linkHandler := handler.NewLinkHandler(linkService)

	return &Container{
		Config:      cfg,
		DB:          db,
		LinkRepo:    linkRepo,
		LinkService: linkService,
		LinkHandler: linkHandler,
		Router:      router.New(linkHandler),
	}, nil
}
