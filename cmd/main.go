package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"social-media-backend/internal/database"
	"social-media-backend/internal/env"

	"github.com/gin-gonic/gin"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 8080, "Sets The Server Port")
	flag.Parse()

	env.LoadEnv("./.env")

	logFile, err := os.OpenFile(
		"app.log",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		panic(err)
	}
	defer func() {
		err := logFile.Close()
		if err != nil {
			log.Printf("error closing log file: %v\n", err)
		}
	}()
	app := application{
		addr:    fmt.Sprintf(":%d", port),
		db:      database.NewDb(),
		rdb:     database.NewRedisClient(),
		router:  gin.Default(),
		logFile: logFile,
	}

	app.mount()
	if err := app.run(); err != nil {
		log.Fatalf("Failed To Run Application: %v", err)
	}
}
