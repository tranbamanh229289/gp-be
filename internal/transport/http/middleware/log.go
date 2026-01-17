package middleware

import (
	"be/pkg/logger"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func LogMiddleware(logger *logger.ZapLogger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		startTime := time.Now().UTC()
		path := ctx.Request.URL.Path
		query := ctx.Request.URL.RawQuery

		if query != "" {
			path = path + "?" + query
		}

		method := ctx.Request.Method
		ctx.Next()

		endTime := time.Now().UTC()
		status := ctx.Writer.Status()

		logger.Info("HTTP request",
			zap.String("method", method),
			zap.String("path", path),
			zap.Int("status", status),
			zap.Time("start_time", startTime),
			zap.Time("end_time", endTime),
		)
	}
}
