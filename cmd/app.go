package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	_ "social-media-backend/docs"
	"social-media-backend/internal/adapters/sqlc/db"
	"social-media-backend/internal/env"
	"social-media-backend/internal/logging"
	"social-media-backend/internal/network"
	"social-media-backend/internal/repo/postgres"
	rdrepo "social-media-backend/internal/repo/redis"
	"social-media-backend/internal/service"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type application struct {
	engine  *gin.Engine
	addr    string
	db      *pgxpool.Pool
	rdb     *redis.Client
	logFile *os.File
	envCfg  *env.Config
}

func (a *application) mount() {
	a.engine.Use(func(c *gin.Context) {
		id := uuid.New().String()
		c.Writer.Header().Set("X-Request-ID", id)
		c.Set("request_id", id)
		c.Next()
	})

	a.engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-CSRF-Token", "X-Auth-Mode"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	a.engine.GET("/health", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Status Is Available")
	})

	// docs.SwaggerInfo.BasePath = "/api/v1"
	a.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	queries := db.New(a.db)
	logger := logging.NewLogger(logging.LoggerConfig{
		Level:  slog.LevelDebug,
		Format: logging.JSON,

		Output: a.logFile,
	})

	// repos
	postgresRepo := postgres.NewPostgresRepository(a.db, queries)
	redisRepo := rdrepo.NewRedisRepository(a.rdb)

	// services
	userSvc := service.NewUserService(postgresRepo)
	tokenSvc := service.NewTokenService(redisRepo)
	authSvc := service.NewAuthService(postgresRepo, postgresRepo, *tokenSvc, logger, a.envCfg.JWTKey)

	// handlers
	userHandler := network.NewUserHandler(userSvc)
	authHandler := network.NewAuthHandler(authSvc)

	// network
	mw := network.InitMiddlwares(a.envCfg.JWTKey, logger)

	// v1
	{
		v1 := a.engine.Group("/api/v1")
		v1.Use(mw.Logger, mw.Error)

		router := network.NewRouter(v1, &network.Handlers{
			UserHandler: userHandler,
			AuthHandler: authHandler,
		})
		router.Register(mw)
	}
}

func (a *application) run() error {
	srv := &http.Server{
		Addr:    a.addr,
		Handler: a.engine,
	}

	log.Printf("Server started at %s", a.addr)

	serverErr := make(chan error, 1)
	go func() {
		serverErr <- srv.ListenAndServe()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErr:
		return err

	case sig := <-quit:
		log.Printf("shutdown signal received: %v", sig)
		return a.shutdown(srv)
	}
}

func (a *application) shutdown(srv *http.Server) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("shutting down server...")

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("http shutdown failed: %v", err)
		_ = srv.Close()
	}

	if a.db != nil {
		log.Println("closing postgres pool...")
		a.db.Close()
	}

	if a.rdb != nil {
		log.Println("closing redis client...")
		if err := a.rdb.Close(); err != nil {
			log.Printf("redis close error: %v", err)
		}
	}

	if a.logFile != nil {
		log.Println("closing log file...")
		_ = a.logFile.Sync()
		_ = a.logFile.Close()
	}

	log.Println("server stopped")
	return nil
}
