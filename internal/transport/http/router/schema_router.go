package router

import (
	"be/internal/shared/constant"
	"be/internal/transport/http/handler"
	"be/internal/transport/http/middleware"

	"github.com/gin-gonic/gin"
)

func (r *Router) SetupSchemaRouter(apiGroup *gin.RouterGroup, schemaHandler *handler.SchemaHandler) {
	schemaGroup := apiGroup.Group("schemas")
	schemaGroup.Use(middleware.AuthenticateMiddleware(r.authZkService))

	schemaGroup.GET("", schemaHandler.GetSchemas)
	schemaGroup.GET("/:id", schemaHandler.GetSchemaByPublicId)
	schemaGroup.GET("/attributes/:id", schemaHandler.GetSchemaAttributeByPublicId)
	schemaGroup.POST("", middleware.AuthorizeMiddleware([]constant.IdentityRole{constant.IdentityIssuerRole}), schemaHandler.CreateSchema)
	schemaGroup.PATCH("/:id", middleware.AuthorizeMiddleware([]constant.IdentityRole{constant.IdentityIssuerRole}), schemaHandler.RemoveSchema)
}
