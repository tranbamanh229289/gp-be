package router

import (
	"be/internal/shared/constant"
	"be/internal/transport/http/handler"
	"be/internal/transport/http/middleware"

	"github.com/gin-gonic/gin"
)

func SetupAuthRouter(apiGroup *gin.RouterGroup, authHandler *handler.AuthHandler) {
	authGroup := apiGroup.Group("auth")
	
	authGroup.GET("refresh-token", authHandler.RefreshToken)
	authGroup.POST("register", authHandler.Register)
	authGroup.POST("login", authHandler.Login)
	
	adminGroup := apiGroup.Group("admin")
	adminGroup.GET("", authHandler.GetAllUser)

	userGroup := apiGroup.Group("users")
	userGroup.Use(middleware.AuthenticateMiddleware(authHandler.GetAuthService()), middleware.AuthorizeMiddleware(constant.UserRoleUser))
	userGroup.GET("profile/:id", authHandler.GetProfile)
	userGroup.PUT("profile/:id", authHandler.UpdateProfile)
}