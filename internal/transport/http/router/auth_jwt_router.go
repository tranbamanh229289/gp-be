package router

import (
	"be/internal/transport/http/handler"

	"github.com/gin-gonic/gin"
)

func (r *Router) SetupAuthJWTRouter(apiGroup *gin.RouterGroup, authJWTHandler *handler.AuthJWTHandler) {
	authJWTGroup := apiGroup.Group("auth")

	authJWTGroup.GET("refresh-token", authJWTHandler.RefreshToken)
	authJWTGroup.POST("register", authJWTHandler.Register)
	authJWTGroup.POST("login", authJWTHandler.Login)

	adminGroup := apiGroup.Group("admin")
	adminGroup.GET("", authJWTHandler.GetAllUser)

	userGroup := apiGroup.Group("users")
	userGroup.GET("profile/:id", authJWTHandler.GetProfile)
	userGroup.PUT("profile/:id", authJWTHandler.UpdateProfile)
}
