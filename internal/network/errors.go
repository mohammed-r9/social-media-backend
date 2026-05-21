package network

import (
	"errors"
	"net/http"
	"social-media-backend/internal/domain"

	"github.com/gin-gonic/gin"
)

type errorResponse struct {
	Status int    `json:"-"`
	Code   string `json:"code"`
	Msg    string `json:"error"`
}

// handleError translates domain and application errors into standardized HTTP responses
func handleError(c *gin.Context, err error) {
	res := mapError(err)

	c.JSON(res.Status, gin.H{
		"error": res.Msg,
		"code":  res.Code,
	})
}

func mapError(err error) errorResponse {
	switch {
	// user errors
	case errors.Is(err, domain.ErrUserNotFound):
		return errorResponse{http.StatusNotFound, "user_not_found", "user not found"}

	case errors.Is(err, domain.ErrInvalidCredentials):
		return errorResponse{http.StatusUnauthorized, "invalid_credentials", "invalid credentials"}

	case errors.Is(err, domain.ErrEmailAlreadyTaken):
		return errorResponse{http.StatusConflict, "email_taken", "email already taken"}

	case errors.Is(err, domain.ErrInvalidUserID):
		return errorResponse{http.StatusBadRequest, "invalid_user_id", "invalid user id"}

	case errors.Is(err, domain.ErrInvalidUserEmail):
		return errorResponse{http.StatusBadRequest, "invalid_user_email", "invalid user email"}

	case errors.Is(err, domain.ErrUnverifiedUserEmail):
		return errorResponse{http.StatusForbidden, "email_not_verified", "unverified user email"}

	case errors.Is(err, domain.ErrInvalidPassword):
		return errorResponse{http.StatusBadRequest, "invalid_password", "invalid password"}

	case errors.Is(err, domain.ErrInvalidOldPassword):
		return errorResponse{http.StatusBadRequest, "invalid_old_password", "invalid old password"}

	case errors.Is(err, domain.ErrInvalidUserName),
		errors.Is(err, domain.ErrInvalidUsername):
		return errorResponse{http.StatusBadRequest, "invalid_username", "invalid username"}

	case errors.Is(err, domain.ErrNoRowsAffected):
		return errorResponse{http.StatusNotFound, "no_rows_affected", "resource not found"}

	// token errors
	case errors.Is(err, domain.ErrTokenNotFound):
		return errorResponse{http.StatusUnauthorized, "token_not_found", "token not found"}

	case errors.Is(err, domain.ErrExpiredToken):
		return errorResponse{http.StatusUnauthorized, "expired_token", "token expired"}

	case errors.Is(err, domain.ErrInvalidToken):
		return errorResponse{http.StatusUnauthorized, "invalid_token", "token is invalid"}

	// session errors
	case errors.Is(err, domain.ErrSessionExpired):
		return errorResponse{http.StatusUnauthorized, "session_expired", "session is expired"}

	case errors.Is(err, domain.ErrSessionRevoked):
		return errorResponse{http.StatusUnauthorized, "session_revoked", "session is revoked"}

	case errors.Is(err, domain.ErrSessionAlreadyExists):
		return errorResponse{http.StatusConflict, "session_exists", "session already exists"}

	case errors.Is(err, domain.ErrSessionNotFound):
		return errorResponse{http.StatusNotFound, "session_not_found", "session not found"}

	default:
		return errorResponse{http.StatusInternalServerError, "internal_error", "internal server error"}
	}
}
