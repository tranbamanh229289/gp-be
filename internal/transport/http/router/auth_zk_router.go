package router

import (
	"be/internal/transport/http/handler"

	"github.com/gin-gonic/gin"
)

func SetupAuthZkRouter(apiGroup *gin.RouterGroup, authZKHandler *handler.AuthZkHandler) {
	authZkGroup := apiGroup.Group("auth")

	authZkGroup.GET("sign-in", authZKHandler.Signin)
	authZkGroup.POST("call-back", authZKHandler.Callback)
}
