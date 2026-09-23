package main

import (
	"log"

	"github.com/Naveen-kumar525/go-links/internal/di"
)

func main() {
	app, err := di.New()
	if err != nil {
		log.Fatalf("failed to initialize application: %v", err)
	}

	log.Printf("Server started on %s", app.Config.Addr)
	if err := app.Router.Run(app.Config.Addr); err != nil {
		log.Fatal(err)
	}
}
