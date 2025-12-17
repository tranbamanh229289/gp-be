package router

import (
	"be/internal/transport/http/handler"

	"github.com/gin-gonic/gin"
)

func SetupHolderRouter(apiGroup *gin.RouterGroup, holderHandler *handler.HolderHandler) {
	holderGroup := apiGroup.Group("holder")

	holderGroup.POST("request", holderHandler.CreateRequest)
	holderGroup.POST("generate-proof", holderHandler.GenerateProof)
	holderGroup.GET("credentials", holderHandler.GetCredentials)
	holderGroup.GET("credential/:id", holderHandler.GetCredentialByPublicId)
}
