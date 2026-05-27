package main

import (
	"log"
	"log/slog"
	"net/http"
	"social-media-backend/internal/adapters/sqlc/db"
	"social-media-backend/internal/logging"
	"social-media-backend/internal/network"
	"social-media-backend/internal/repo/postgres"
	rdrepo "social-media-backend/internal/repo/redis"
	"social-media-backend/internal/service"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type application struct {
	router *gin.Engine
	addr   string
	db     *pgxpool.Pool
	rdb    *redis.Client
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
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-CSRF-Token", "X-Auth-Mode"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	a.router.GET("/health", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Status Is Available")
	})

	queries := db.New(a.db)
	logger := logging.NewLogger(logging.LoggerConfig{
		Level:  slog.LevelDebug,
		Format: logging.JSON,

		// add logging to a file later
		// Output: file,
	})

	// repos
	postgresRepo := postgres.NewPostgresRepository(queries)
	redisRepo := rdrepo.NewRedisRepository(a.rdb)

	// services
	userSvc := service.NewUserService(postgresRepo)
	tokenSvc := service.NewTokenService(redisRepo)
	authSvc := service.NewAuthService(postgresRepo, postgresRepo, *tokenSvc, logger)

	// handlers
	userHandler := network.NewUserHandler(userSvc)
	authHandler := network.NewAuthHandler(authSvc)

	// v1
	{
		v1 := a.router.Group("/api/v1")
		v1.Use(network.ErrorHandler())
		v1.Use(network.Logger(logger))

		network.RegisterUserRoutes(v1, userHandler)
		network.RegisterAuthRoutes(v1, authHandler)
	}
}

func (a *application) run() error {
	log.Printf("Server Has Started At Addr %s", a.addr)
	return a.router.Run(a.addr)
}
