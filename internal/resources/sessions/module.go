package sessions

import (
	"net/http"
	"social-media-backend/internal/adapters/sqlc/db"
	handler "social-media-backend/internal/resources/sessions/internal/http"
	"social-media-backend/internal/resources/sessions/internal/repository/postgres"
	"social-media-backend/internal/resources/sessions/internal/service"

	"github.com/gin-gonic/gin"
)

func InitModule(q *db.Queries, r gin.IRouter) {
	repo := postgres.NewSessionRepository(q)
	svc := service.NewService(repo)
	h := handler.NewHandler(svc)

	registerRoutes(h, r)
}

func registerRoutes(h *handler.Handler, r gin.IRouter) {
	// TODO: Add routes
	r.GET("sessions/health", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "sessions status is available")
	})
}
