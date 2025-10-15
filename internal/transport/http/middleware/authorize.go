package middleware

import "github.com/gin-gonic/gin"

func AuthorizeMiddleware(allowdRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Next()
	}
}