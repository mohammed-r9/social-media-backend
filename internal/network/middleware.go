package network

import (
	"log/slog"
	"social-media-backend/internal/crypto/tokens"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err

		status, resp := mapError(err)

		Fail(c, status, resp.Error.Code, resp.Error.Message)
	}
}

func Logger(logger *slog.Logger) gin.HandlerFunc {
	/*
	   Planned logger improvements:

	   HIGH PRIORITY:
	   - user_id => How would I do that?

	   MEDIUM PRIORITY:
	   - session_id
	   - error classification (4xx vs 5xx grouping)

	   LOW PRIORITY:
	   - request_size / response_size
	*/
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		reqID, _ := c.Get("request_id")
		latency := time.Since(start)
		statusAny, exists := c.Get("status_code")
		status := 200
		if exists {
			status = statusAny.(int)
		}

		if v, ok := c.Get("status_code"); ok {
			status = v.(int)
		}

		errs := c.Errors.ByType(gin.ErrorTypeAny)

		attrs := []any{
			"request_id", reqID,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", status,
			"latency_ms", latency.Milliseconds(),
			"ip", c.ClientIP(),
		}

		if len(errs) > 0 {
			logger.Error("request failed",
				append(attrs, "error", errs.Last().Err.Error())...,
			)
			return
		}

		logger.Info("request", attrs...)
	}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			_ = c.Error(errMissingAuthorizationHeader)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		if token == "" {
			_ = c.Error(errMissingAccessToken)
			return
		}

		claims, err := tokens.VerifyAccessToken(token)
		if err != nil {
			_ = c.Error(errInvalidAccessToken)
			return
		}

		c.Set(ctxClaimKey, claims)

		c.Next()
	}
}
