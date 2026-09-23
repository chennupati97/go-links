package di

import (
	"fmt"
	"os"

	"github.com/chennupati97/go-links/internal/handler"
	"github.com/chennupati97/go-links/internal/repository"
	"github.com/chennupati97/go-links/internal/router"
	"github.com/chennupati97/go-links/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Settings holds process-level runtime options.
type Settings struct {
	ListenAddr   string
	DatabaseFile string
}

func ReadSettings() Settings {
	addr := lookupEnv("LISTEN_PORT", "8080")
	if len(addr) > 0 && addr[0] != ':' {
		addr = ":" + addr
	}

	return Settings{
		ListenAddr:   addr,
		DatabaseFile: lookupEnv("SQLITE_FILE", "jumpalias.db"),
	}
}

func lookupEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// App is the composition root that wires dependencies together.
type App struct {
	Settings Settings
	DB       *gorm.DB

	Store   repository.ShortcutStore
	Manager *service.AliasManager
	API     *handler.ShortcutAPI

	Engine *gin.Engine
}

// Bootstrap constructs the full application graph.
func Bootstrap() (*App, error) {
	settings := ReadSettings()

	conn, err := repository.ConnectSQLite(settings.DatabaseFile)
	if err != nil {
		return nil, fmt.Errorf("connect sqlite: %w", err)
	}

	store := repository.NewSQLiteShortcutStore(conn)
	manager := service.NewAliasManager(store)
	api := handler.NewShortcutAPI(manager)

	return &App{
		Settings: settings,
		DB:       conn,
		Store:    store,
		Manager:  manager,
		API:      api,
		Engine:   router.BuildEngine(api),
	}, nil
}
