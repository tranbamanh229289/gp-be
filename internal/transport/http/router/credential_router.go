package router

import (
	"be/internal/infrastructure/database/postgres"
	"be/internal/shared/helper"
	"be/internal/transport/http/handler"
	"be/internal/transport/http/middleware"

	"github.com/gin-gonic/gin"
)

func (r *Router) SetupCredentialRouter(apiGroup *gin.RouterGroup, credentialHandler *handler.CredentialHandler, db *postgres.PostgresDB) {
	credentialGroup := apiGroup.Group("credentials")
	credentialGroup.Use(middleware.AuthenticateMiddleware(r.authZkService))

	verifiableGroup := credentialGroup.Group("verifiable")
	requestGroup := credentialGroup.Group("request")

	requestGroup.GET("", credentialHandler.GetCredentialRequests)
	requestGroup.POST("", credentialHandler.CreateCredentialRequest)
	requestGroup.PATCH("/:id", credentialHandler.UpdateCredentialRequest)

	verifiableGroup.GET("", credentialHandler.GetVerifiableCredentials)
	verifiableGroup.GET("/:id", credentialHandler.GetVerifiableCredentialById)
	verifiableGroup.PATCH("/:id", credentialHandler.UpdateVerifiableCredential)
	verifiableGroup.POST("/:id", helper.TxMiddleware(db.GetGormDB()), credentialHandler.IssueVerifiableCredential)
}
