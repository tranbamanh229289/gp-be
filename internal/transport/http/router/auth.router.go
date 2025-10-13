package router

import (
	"be/internal/transport/http/handler"
	"be/internal/transport/http/middleware"

	"github.com/gin-gonic/gin"
)

func SetupUserRouter(apiGroup *gin.RouterGroup, authHandler *handler.AuthHandler) {
	userGroup := apiGroup.Group("users")
	{
		userGroup.Use(middleware.AuthMiddleware())
		userGroup.GET("profile", authHandler.GetProfile)
		userGroup.GET("users", authHandler.GetAllUser)
		userGroup.POST("register", authHandler.Register)
		userGroup.POST("login", authHandler.Login)
	}
}