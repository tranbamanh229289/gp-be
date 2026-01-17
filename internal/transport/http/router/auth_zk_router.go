package router

import (
	"be/internal/transport/http/handler"
	"be/internal/transport/http/middleware"

	"github.com/gin-gonic/gin"
)

func (r *Router) SetupAuthZkRouter(apiGroup *gin.RouterGroup, authZKHandler *handler.AuthZkHandler) {
	authZkGroup := apiGroup.Group("authzk")

	authZkGroup.GET("", authZKHandler.GetIdentityByRole)
	authZkGroup.GET("/:did", authZKHandler.GetIdentityByDID)
	authZkGroup.GET("logout", middleware.AuthenticateMiddleware(r.authZkService), authZKHandler.Logout)
	authZkGroup.GET("challenge", authZKHandler.Challenge)
	authZkGroup.GET("refresh-token", authZKHandler.RefreshZKToken)
	authZkGroup.POST("register", authZKHandler.Register)
	authZkGroup.POST("login", authZKHandler.Login)

}
