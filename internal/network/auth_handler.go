package network

import (
	"errors"
	"log"
	"net/http"
	"social-media-backend/internal/crypto/tokens"
	"social-media-backend/internal/domain"
	"social-media-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	svc *service.AuthService
}

func NewAuthHandler(svc *service.AuthService) *AuthHandler {
	return &AuthHandler{
		svc: svc,
	}
}

// createUserRequest represents the registration request body.
type createUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
}

// Register godoc
// @Summary      Register a new user
// @Description  Create a new user account with the provided details.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      createUserRequest  true  "Registration Details"
// @Success      201      {object}  domain.User
// @Failure      400      {object}  Response
// @Failure      409      {object}  Response
// @Failure      500      {object}  Response
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req createUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// need to update it later
		log.Printf("error binding request body: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	createdUser, err := h.svc.Register(c.Request.Context(), service.RegisterParams{
		Name:              req.Name,
		Email:             req.Email,
		PassowrdPlainText: req.Password,
		Username:          req.Username,
	})

	if err != nil {
		log.Printf("request failed: %v", err)
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, createdUser)
}

// VerifyUserEmail godoc
// @Summary      Verify user email
// @Description  Verify a user's email address using a token sent via email.
// @Tags         auth
// @Produce      json
// @Param        t    query     string  true  "Verification Token"
// @Success      200  {object}  Response
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /auth/verify-email [get]
func (h *AuthHandler) VerifyUserEmail(c *gin.Context) {
	token := c.Query("t")
	if token == "" {
		// need to update it later
		log.Println("invalid email verification token")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "bad request",
		})
	}

	err := h.svc.VerifyUserEmail(c.Request.Context(), token)

	if err != nil {
		log.Printf("request failed: %v", err)
		_ = c.Error(err)
		return
	}
}

// loginRequest represents the login request body.
type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Login godoc
// @Summary      Login user
// @Description  Authenticate a user and return access/refresh tokens.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        X-Auth-Mode  header    string        true  "Auth mode (cookie or token)"
// @Param        request      body      loginRequest  true  "Login Credentials"
// @Success      200          {object}  Response{data=domain.AuthTokens}
// @Failure      400          {object}  Response
// @Failure      401          {object}  Response
// @Failure      500          {object}  Response
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	authMode := c.GetHeader("X-Auth-Mode") // SHOULD ONLY BE USED TO FORMAT THE RESPONSE
	if authMode == "" {
		_ = c.Error(errMissingAuthModeHeader)
		return
	}

	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	sessionTokens, err := h.svc.Login(c.Request.Context(), service.LoginParams{
		Email:      req.Email,
		Password:   req.Password,
		DeviceName: "temp",
	})
	if err != nil {
		_ = c.Error(err)
		return
	}

	switch authMode {
	case "cookie":
		c.SetCookie(
			"refresh_token",
			sessionTokens.RefreshToken,
			int(tokens.REFRESH_TTL),
			"/",
			"",
			true,
			true,
		)
		c.SetCookie("csrf_token",
			sessionTokens.CsrfToken,
			int(tokens.REFRESH_TTL),
			"/",
			"",
			true,
			false,
		)
		c.SetCookie("session_id",
			sessionTokens.SessionID,
			int(tokens.REFRESH_TTL),
			"/",
			"",
			true,
			true,
		)

		OK(c, gin.H{
			"access_token": sessionTokens.AccessToken,
		})

	case "token":
		OK(c, gin.H{
			"access_token":  sessionTokens.AccessToken,
			"refresh_token": sessionTokens.RefreshToken,
			"session_id":    sessionTokens.SessionID,
		})
	default:
		_ = c.Error(errInvalidAuthModeHeader)
		return
	}
}

// RefreshAccessToken godoc
// @Summary      Refresh access token
// @Description  Get a new access token using a valid refresh token.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      refreshTokenReq  false  "Refresh tokens (if not using cookies)"
// @Success      200      {object}  Response{data=map[string]string}
// @Failure      401      {object}  Response
// @Failure      500      {object}  Response
// @Router       /auth/refresh [post]
func (h *AuthHandler) RefreshAccessToken(c *gin.Context) {
	var opTokens refreshTokens

	// web
	opTokens, err := refreshGetTokensFromCookies(c)
	if err != nil {
		// mobile
		opTokens, err = refreshGetTokensFromBody(c)
		if err != nil {
			_ = c.Error(err)
			return
		}
	}

	accessToken, err := h.svc.RefreshAccessToken(c.Request.Context(), service.RefreshParams{
		RefreshToken: opTokens.RefreshToken,
		SessionID:    opTokens.SessionID,
		CsrfToken:    opTokens.CsrfToken,
	})

	if err != nil {
		_ = c.Error(err)
		return
	}

	OK(c, gin.H{
		"access_token": accessToken,
	})
}

// forgotPassowrdRequest represents the forgot password request body.
type forgotPassowrdRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ForgotPassword godoc
// @Summary      Forgot password
// @Description  Send a password reset link to the user's email.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      forgotPassowrdRequest  true  "Email"
// @Success      200      {object}  Response{data=map[string]string}
// @Failure      400      {object}  Response
// @Failure      500      {object}  Response
// @Router       /auth/forgot-password [post]
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req forgotPassowrdRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	err := h.svc.AskForPasswordReset(c.Request.Context(), req.Email)
	if err != nil {
		if !errors.Is(err, domain.ErrUserNotFound) {
			_ = c.Error(err)
			return
		}
	}

	OK(c, gin.H{
		"message": "If the account exists, a password reset email has been sent.",
	})
}

// resetPasswordRequest represents the password reset request body.
type resetPasswordRequest struct {
	Password string `json:"password" binding:"required"`
	Token    string `json:"token" binding:"required"`
}

// ResetPassword godoc
// @Summary      Reset password
// @Description  Reset user password using a valid reset token.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      resetPasswordRequest  true  "New Password and Token"
// @Success      200      {object}  Response{data=map[string]string}
// @Failure      400      {object}  Response
// @Failure      500      {object}  Response
// @Router       /auth/reset-password [post]
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req resetPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	err := h.svc.ResetUserPassword(c.Request.Context(), service.ResetUserPasswordParams{
		Token:       req.Token,
		NewPassword: req.Password,
	})

	if err != nil {
		_ = c.Error(err)
		return
	}

	OK(c, gin.H{
		"message": "Password resetted successfully",
	})
}
