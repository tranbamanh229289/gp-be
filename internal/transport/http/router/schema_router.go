package router

import (
	"be/internal/transport/http/handler"

	"github.com/gin-gonic/gin"
)

func SetupSchemaRouter(apiGroup *gin.RouterGroup, schemaHandler *handler.SchemaHandler) {
	schemaGroup := apiGroup.Group("schemas")

	schemaGroup.POST("", schemaHandler.CreateSchema)
	schemaGroup.PATCH("", schemaHandler.RemoveSchema)
}
