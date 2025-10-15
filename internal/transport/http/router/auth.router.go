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
		userGroup.GET("/", authHandler.GetAllUser)
		userGroup.GET("profile/:id", authHandler.GetProfile)
		userGroup.PUT("profile/:id", authHandler.UpdateProfile)
		userGroup.GET("refresh-token", authHandler.RefreshToken)
		userGroup.POST("register", authHandler.Register)
		userGroup.POST("login", authHandler.Login)
	}
}