package main

import (
	"flag"
	"fmt"
	"log"
	"social-media-backend/internal/database"
	"social-media-backend/internal/shared/env"

	"github.com/gin-gonic/gin"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 8080, "Sets The Server Port")
	flag.Parse()

	env.LoadEnv("./.env")
	app := application{
		addr:   fmt.Sprintf(":%d", port),
		db:     database.NewDb(),
		router: gin.Default(),
	}
	app.mount()
	if err := app.run(); err != nil {
		log.Fatalf("Failed To Run Application: %v", err)
	}
}
