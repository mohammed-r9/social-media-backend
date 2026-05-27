package network

import (
	"log/slog"
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
	   - request_id
	   - user_id

	   MEDIUM PRIORITY:
	   - session_id
	   - error classification (4xx vs 5xx grouping)

	   LOW PRIORITY:
	   - request_size / response_size
	*/
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		errs := c.Errors.ByType(gin.ErrorTypeAny)

		attrs := []any{
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
