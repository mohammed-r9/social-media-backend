package network

// TODO: use google wire instead to manage dependency injection

import (
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	UserHandler *UserHandler
	AuthHandler *AuthHandler
}
type Router struct {
	router   *gin.RouterGroup
	handlers *Handlers
}

func NewRouter(r *gin.RouterGroup, h *Handlers) *Router {
	return &Router{
		router:   r,
		handlers: h,
	}
}

func (r *Router) Register(mw *Middlewares) {
	registerAuthRoutes(r.router, r.handlers.AuthHandler)
	registerUserRoutes(r.router, r.handlers.UserHandler, mw.Auth)
}

func registerUserRoutes(r *gin.RouterGroup, h *UserHandler, authMW gin.HandlerFunc) {
	users := r.Group("/users")

	{
		users.Use(authMW)
		users.PUT("me/password", h.UpdateSelfPassword)
		users.PUT("/me/avatar", h.UpdateSelfAvatar)
		users.GET("me", h.GetSelfUser)
	}
}

func registerAuthRoutes(r *gin.RouterGroup, h *AuthHandler) {
	auth := r.Group("/auth")

	{
		auth.POST("register", h.Register)
		auth.POST("login", h.Login)
		auth.POST("refresh", h.RefreshAccessToken)
		auth.POST("forgot-password", h.ForgotPassword)
		auth.POST("reset-password", h.ResetPassword)
		auth.GET("verify-email", h.VerifyUserEmail)
	}
}
