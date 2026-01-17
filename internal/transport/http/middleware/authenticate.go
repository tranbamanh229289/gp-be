package middleware

import (
	"be/internal/service"
	"be/internal/shared/constant"
	response "be/internal/shared/helper"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthenticateMiddleware(authService service.IAuthZkService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			response.RespondError(c, &constant.InvalidAuthHeader)
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.RespondError(c, &constant.InvalidAuthHeader)
			c.Abort()
			return
		}

		tokenString := parts[1]

		claims, err := authService.VerifyZKToken(tokenString, constant.AccessToken)
		if err != nil {
			response.RespondError(c, &constant.InvalidToken)
			c.Abort()
			return
		}
		c.Set("user", claims)

		c.Next()
	}
}
