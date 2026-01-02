package router

import (
	"be/internal/transport/http/handler"

	"github.com/gin-gonic/gin"
)

func SetupAuthZkRouter(apiGroup *gin.RouterGroup, authZKHandler *handler.AuthZkHandler) {
	authZkGroup := apiGroup.Group("authzk")

	authZkGroup.GET("", authZKHandler.GetIdentityByRole)
	authZkGroup.GET("challenge", authZKHandler.Challenge)
	authZkGroup.POST("register", authZKHandler.Register)
	authZkGroup.POST("login", authZKHandler.Login)
}
