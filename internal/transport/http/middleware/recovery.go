package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("panic: %v", r)
				c.JSON(http.StatusInternalServerError, gin.H{

				})
				c.Abort()
			}
		}()
		c.Next()
	}
}