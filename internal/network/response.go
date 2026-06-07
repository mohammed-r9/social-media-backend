package network

import (
	"errors"
	"net/http"
	"social-media-backend/internal/apperrors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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
	// request body validation error
	if ve, ok := errors.AsType[validator.ValidationErrors](err); ok {
		errStr := ""

		for _, fe := range ve {
			field := fe.Field()
			switch fe.Tag() {
			case "required":
				errStr = field + " " + "is required"
			case "email":
				errStr = field + " " + "must be a valid email"
			case "min":
				errStr = field + " " + "is too short"
			case "max":
				errStr = field + " " + "is too long"
			default:
				errStr = field + " " + "is invalid"
			}
		}

		return http.StatusBadRequest, Response{
			Success: false,
			Error: &ErrorInfo{
				Code:    "validation_error",
				Message: errStr,
			}}
	}

	var (
		authErr    apperrors.AuthError
		userErr    apperrors.UserError
		sessionErr apperrors.SessionError
		dbErr      apperrors.DatabaseError
		networkErr apperrors.NetworkError
	)

	switch {
	case errors.As(err, &authErr):
		status := http.StatusUnauthorized
		if authErr == apperrors.UnverifiedUserEmail {
			status = http.StatusForbidden
		}

		return status, Response{
			Success: false,
			Error: &ErrorInfo{
				Code:    authErr.Code(),
				Message: authErr.Error(),
			},
		}

	case errors.As(err, &userErr):
		status := http.StatusBadRequest
		if userErr == apperrors.UserNotFound {
			status = http.StatusNotFound
		}

		if userErr == apperrors.EmailAlreadyTaken {
			status = http.StatusConflict
		}

		return status, Response{
			Success: false,
			Error: &ErrorInfo{
				Code:    userErr.Code(),
				Message: userErr.Error(),
			},
		}

	case errors.As(err, &sessionErr):
		status := http.StatusUnauthorized
		if sessionErr == apperrors.SessionAlreadyExists {
			status = http.StatusConflict
		}

		if sessionErr == apperrors.SessionNotFound {
			status = http.StatusNotFound
		}

		return status, Response{
			Success: false,
			Error: &ErrorInfo{
				Code:    sessionErr.Code(),
				Message: sessionErr.Error(),
			},
		}

	case errors.As(err, &dbErr):
		status := http.StatusInternalServerError
		if dbErr == apperrors.NoRowsAffected {
			status = http.StatusNotFound
		}

		return status, Response{
			Success: false,
			Error: &ErrorInfo{
				Code:    dbErr.Code(),
				Message: dbErr.Error(),
			},
		}

	case errors.As(err, &networkErr):
		return http.StatusBadRequest, Response{
			Success: false,
			Error: &ErrorInfo{
				Code:    networkErr.Code(),
				Message: networkErr.Error(),
			},
		}

	default:
		return http.StatusInternalServerError, Response{
			Success: false,
			Error:   &ErrorInfo{Code: "internal_error", Message: "internal server error"},
		}
	}
}
