package network

import "github.com/gin-gonic/gin"

func RegisterUserRoutes(r *gin.RouterGroup, h *UserHandler) {
	users := r.Group("/users")
	{
		users.POST("", h.Register)
	}
}
