package middleware

import (
	"github.com/gin-gonic/gin"
)

func (m *Middleware) CORSMiddleware(engine *gin.Engine) {
	allowedOrigins := m.config.App.AllowedOrigins
	middleware := func(c *gin.Context) {
		m.logger.Info(allowedOrigins)
		c.Writer.Header().Set("Access-Control-Allow-Origin", allowedOrigins)
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
	engine.Use(middleware)
}
