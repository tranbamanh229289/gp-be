package router

import (
	"be/internal/transport/http/handler"

	"github.com/gin-gonic/gin"
)

func SetupCredentialRouter(apiGroup *gin.RouterGroup, credentialHandler *handler.CredentialHandler) {
	credentialRequestGroup := apiGroup.Group("credential-requests")

	credentialRequestGroup.POST("", credentialHandler.CreateCredentialRequest)
	credentialRequestGroup.GET("", credentialHandler.GetCredentialRequests)
}
