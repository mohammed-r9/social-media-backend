package network

import "github.com/gin-gonic/gin"

func RegisterUserRoutes(r *gin.RouterGroup, h *UserHandler) {
	users := r.Group("/users")
	{
		_ = users
	}
}

func RegisterAuthRoutes(r *gin.RouterGroup, h *AuthHandler) {
	auth := r.Group("/auth")
	{
		auth.POST("register", h.Register)
		auth.POST("login", h.Login)
	}
}
