package middleware

import (
	"github.com/gin-gonic/gin"
)

func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				
				c.Abort()
			}
		}()
		c.Next()
	}
}