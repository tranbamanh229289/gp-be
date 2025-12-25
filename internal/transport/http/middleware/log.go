package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (m *Middleware) LogMiddleware(engine *gin.Engine) {
	middleware := func(ctx *gin.Context) {
		startTime := time.Now()
		path := ctx.Request.URL.Path
		query := ctx.Request.URL.RawQuery

		if query != "" {
			path = path + "?" + query
		}

		method := ctx.Request.Method
		ctx.Next()

		endTime := time.Now()
		status := ctx.Writer.Status()

		m.logger.Info("HTTP request",
			zap.String("method", method),
			zap.String("path", path),
			zap.Int("status", status),
			zap.Time("start_time", startTime),
			zap.Time("end_time", endTime),
		)
	}
	engine.Use(middleware)
}
