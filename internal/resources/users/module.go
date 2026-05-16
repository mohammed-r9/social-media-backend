package users

import (
	"net/http"
	"social-media-backend/internal/adapters/sqlc/db"
	handler "social-media-backend/internal/resources/users/internal/http"
	"social-media-backend/internal/resources/users/internal/repository/postgres"
	"social-media-backend/internal/resources/users/internal/service"

	"github.com/gin-gonic/gin"
)

func InitModule(q *db.Queries, r gin.IRouter) {
	repo := postgres.NewUserRepository(q)
	svc := service.NewService(repo)
	h := handler.NewHandler(svc)

	registerRoutes(h, r)
}

func registerRoutes(h *handler.Handler, r gin.IRouter) {
	// TODO: Add routes
	r.GET("users/health", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "user status is available")
	})
	r.POST("/users", h.Register)
}
