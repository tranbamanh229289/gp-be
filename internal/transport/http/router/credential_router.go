package router

import (
	"be/internal/transport/http/handler"

	"github.com/gin-gonic/gin"
)

func SetupCredentialRouter(apiGroup *gin.RouterGroup, credentialHandler *handler.CredentialHandler) {
	credentialGroup := apiGroup.Group("credentials")
	credentialRequestGroup := credentialGroup.Group("request")
	verifiableCredentialGroup := credentialGroup.Group("verifiable")

	credentialRequestGroup.GET("", credentialHandler.GetCredentialRequests)
	credentialRequestGroup.POST("", credentialHandler.CreateCredentialRequest)
	credentialRequestGroup.PATCH("/:id", credentialHandler.UpdateCredentialRequest)

	verifiableCredentialGroup.GET("", credentialHandler.GetVerifiableCredentials)
	verifiableCredentialGroup.GET("/:id", credentialHandler.GetVerifiableCredential)
	verifiableCredentialGroup.PATCH("/:id", credentialHandler.UpdateVerifiableCredential)
	verifiableCredentialGroup.PATCH("sign/:id", credentialHandler.SignVerifiableCredential)
	verifiableCredentialGroup.POST("/:id", credentialHandler.IssueVerifiableCredential)
}
