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

	cfg, err := env.New("../.env")
	if err != nil {
		log.Fatalf("error initializing env config: %v\n", err)
	}

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
		addr: fmt.Sprintf(":%d", port),
		db:   database.NewDb(cfg.POSTGRES_CONNECTION),
		rdb: database.NewRedisClient(database.RedisConfig{
			Addr:     cfg.REDIS_ADDR,
			Password: cfg.REDIS_PASSWORD,
		}),
		engine:  gin.Default(),
		logFile: logFile,
		envCfg:  cfg,
	}

	app.mount()
	if err := app.run(); err != nil {
		log.Fatalf("Failed To Run Application: %v", err)
	}
}
