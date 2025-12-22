package router

import (
	"be/internal/transport/http/handler"

	"github.com/gin-gonic/gin"
)

func SetupAuthZkRouter(apiGroup *gin.RouterGroup, authZKHandler *handler.AuthZkHandler) {
	authZkGroup := apiGroup.Group("auth")

	authZkGroup.POST("login", authZKHandler.Login)
	authZkGroup.GET("register", authZKHandler.Register)
	authZkGroup.POST("callback", authZKHandler.Callback)
}
