package network

import (
	"errors"
	"net/http"
	"social-media-backend/internal/crypto/tokens"
	"social-media-backend/internal/domain"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool       `json:"success"`
	Data    any        `json:"data,omitempty"`
	Error   *ErrorInfo `json:"error,omitempty"`
	Meta    *Meta      `json:"meta,omitempty"`
}

type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Meta struct {
	Page       int `json:"page,omitempty"`
	PerPage    int `json:"per_page,omitempty"`
	Total      int `json:"total,omitempty"`
	TotalPages int `json:"total_pages,omitempty"`
}

// OK sends a success response.
func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    data,
	})
}

// Fail sends an error response.
func Fail(c *gin.Context, status int, code, message string) {
	c.Set("status_code", status)

	c.AbortWithStatusJSON(status, Response{
		Success: false,
		Error:   &ErrorInfo{Code: code, Message: message},
	})
}

func mapError(err error) (int, Response) {
	switch {
	// user errors
	case errors.Is(err, domain.ErrUserNotFound):
		return http.StatusNotFound, Response{
			Success: false,
			Error:   &ErrorInfo{Code: "user_not_found", Message: "user not found"},
		}

	case errors.Is(err, domain.ErrInvalidCredentials):
		return http.StatusUnauthorized, Response{
			Success: false,
			Error:   &ErrorInfo{Code: "invalid_credentials", Message: "invalid credentials"},
		}

	case errors.Is(err, domain.ErrEmailAlreadyTaken):
		return http.StatusConflict, Response{
			Success: false,
			Error:   &ErrorInfo{Code: "email_taken", Message: "email already taken"},
		}

	case errors.Is(err, domain.ErrInvalidUserID):
		return http.StatusBadRequest, Response{
			Success: false,
			Error:   &ErrorInfo{Code: "invalid_user_id", Message: "invalid user id"},
		}

	case errors.Is(err, domain.ErrInvalidUserEmail):
		return http.StatusBadRequest, Response{
			Success: false,
			Error:   &ErrorInfo{Code: "invalid_user_email", Message: "invalid user email"},
		}

	case errors.Is(err, domain.ErrUnverifiedUserEmail):
		return http.StatusForbidden, Response{
			Success: false,
			Error:   &ErrorInfo{Code: "email_not_verified", Message: "unverified user email"},
		}

	case errors.Is(err, domain.ErrInvalidPassword):
		return http.StatusBadRequest, Response{
			Success: false,
			Error:   &ErrorInfo{Code: "invalid_password", Message: "invalid password"},
		}

	case errors.Is(err, domain.ErrInvalidOldPassword):
		return http.StatusBadRequest, Response{
			Success: false,
			Error:   &ErrorInfo{Code: "invalid_old_password", Message: "invalid old password"},
		}

	case errors.Is(err, domain.ErrInvalidUserName), errors.Is(err, domain.ErrInvalidUsername):
		return http.StatusBadRequest, Response{
			Success: false,
			Error:   &ErrorInfo{Code: "invalid_username", Message: "invalid username"},
		}

	case errors.Is(err, domain.ErrNoRowsAffected):
		return http.StatusNotFound, Response{
			Success: false,
			Error:   &ErrorInfo{Code: "not_found", Message: "resource not found"},
		}

	// token errors
	case errors.Is(err, tokens.ErrTokenNotFound):
		return http.StatusUnauthorized, Response{
			Success: false,
			Error:   &ErrorInfo{Code: "token_not_found", Message: "token not found"},
		}

	case errors.Is(err, tokens.ErrExpiredToken):
		return http.StatusUnauthorized, Response{
			Success: false,
			Error:   &ErrorInfo{Code: "expired_token", Message: "token expired"},
		}

	case errors.Is(err, tokens.ErrInvalidToken):
		return http.StatusUnauthorized, Response{
			Success: false,
			Error:   &ErrorInfo{Code: "invalid_token", Message: "token is invalid"},
		}

	// session errors
	case errors.Is(err, domain.ErrSessionExpired):
		return http.StatusUnauthorized, Response{
			Success: false,
			Error:   &ErrorInfo{Code: "session_expired", Message: "session is expired"},
		}

	case errors.Is(err, domain.ErrSessionRevoked):
		return http.StatusUnauthorized, Response{
			Success: false,
			Error:   &ErrorInfo{Code: "session_revoked", Message: "session is revoked"},
		}

	case errors.Is(err, domain.ErrSessionAlreadyExists):
		return http.StatusConflict, Response{
			Success: false,
			Error:   &ErrorInfo{Code: "session_exists", Message: "session already exists"},
		}

	case errors.Is(err, domain.ErrSessionNotFound):
		return http.StatusNotFound, Response{
			Success: false,
			Error:   &ErrorInfo{Code: "session_not_found", Message: "session not found"},
		}

	// network layer errors
	case errors.Is(err, errMissingAuthModeHeader):
		return http.StatusBadRequest, Response{
			Success: false,
			Error:   &ErrorInfo{Code: "missing_auth_mode_header", Message: "missing auth mode header"},
		}
	case errors.Is(err, errInvalidAuthModeHeader):
		return http.StatusBadRequest, Response{
			Success: false,
			Error:   &ErrorInfo{Code: "invalid_auth_mode_header", Message: "invalid auth mode header"},
		}
	default:
		return http.StatusInternalServerError, Response{
			Success: false,
			Error:   &ErrorInfo{Code: "internal_error", Message: "internal server error"},
		}
	}
}
