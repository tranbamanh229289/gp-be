package middleware

import (
	"be/internal/shared/constant"
	response "be/internal/shared/helper"
	"be/internal/transport/http/dto"

	"github.com/gin-gonic/gin"
)

func AuthorizeMiddleware(allowedRoles []constant.IdentityRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, existed := c.Get("user")
		if !existed {
			response.RespondError(c, &constant.InternalServer)
			c.Abort()
			return
		}

		claims, ok := user.(*dto.ZKClaims)
		if !ok {
			response.RespondError(c, &constant.InternalServer)
			c.Abort()
			return
		}

		hasRole := false

		for _, role := range allowedRoles {
			if role == claims.Role {
				hasRole = true
				break
			}
		}

		if !hasRole {
			response.RespondError(c, &constant.Forbidden)
			c.Abort()
			return
		}

		c.Next()
	}
}
