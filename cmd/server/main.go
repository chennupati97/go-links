package main

import (
	"log"

	"github.com/chennupati97/go-links/internal/di"
)

func main() {
	app, err := di.Bootstrap()
	if err != nil {
		log.Fatalf("bootstrap failed: %v", err)
	}

	log.Printf("JumpAlias listening on %s", app.Settings.ListenAddr)
	if err := app.Engine.Run(app.Settings.ListenAddr); err != nil {
		log.Fatal(err)
	}
}
