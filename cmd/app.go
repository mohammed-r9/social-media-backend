package main

import (
	"database/sql"
	"log"
	"net/http"
	"social-media-backend/internal/adapters/sqlc/db"
	"social-media-backend/internal/resources/user"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type application struct {
	router *gin.Engine
	addr   string
	db     *sql.DB
}

func (a *application) mount() {
	a.router.Use(func(c *gin.Context) {
		id := uuid.New().String()
		c.Writer.Header().Set("X-Request-ID", id)
		c.Set("request_id", id)
		c.Next()
	})
	a.router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	a.router.GET("/health", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Status Is Available")
	})

	queries := db.New(a.db)
	user.InitModule(queries, &a.router.RouterGroup)
}

func (a *application) run() error {
	log.Printf("Server Has Started At Addr %s", a.addr)
	return a.router.Run(a.addr)
}
