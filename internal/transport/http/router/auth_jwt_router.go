package router

import (
	"be/internal/shared/constant"
	"be/internal/transport/http/handler"
	"be/internal/transport/http/middleware"

	"github.com/gin-gonic/gin"
)

func SetupAuthJWTRouter(apiGroup *gin.RouterGroup, authJWTHandler *handler.AuthJWTHandler) {
	authJWTGroup := apiGroup.Group("auth")

	authJWTGroup.GET("refresh-token", authJWTHandler.RefreshToken)
	authJWTGroup.POST("register", authJWTHandler.Register)
	authJWTGroup.POST("login", authJWTHandler.Login)

	adminGroup := apiGroup.Group("admin")
	adminGroup.GET("", authJWTHandler.GetAllUser)

	userGroup := apiGroup.Group("users")
	userGroup.Use(middleware.AuthenticateMiddleware(authJWTHandler.GetAuthService()), middleware.AuthorizeMiddleware(constant.UserRoleUser))
	userGroup.GET("profile/:id", authJWTHandler.GetProfile)
	userGroup.PUT("profile/:id", authJWTHandler.UpdateProfile)
}
