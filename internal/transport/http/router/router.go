package router

import (
	"be/internal/transport/http/handler"

	"github.com/gin-gonic/gin"
)

type Router struct {
	authJWTHandler  *handler.AuthJWTHandler
	authZkHandler   *handler.AuthZkHandler
	issuerHandler   *handler.IssuerHandler
	holderHandler   *handler.HolderHandler
	verifierHandler *handler.VerifierHandler
	documentHandler *handler.DocumentHandler
}

func NewRouter(
	authJWTHandler *handler.AuthJWTHandler,
	authZkHandler *handler.AuthZkHandler,
	issuerHandler *handler.IssuerHandler,
	holderHandler *handler.HolderHandler,
	verifierHandler *handler.VerifierHandler,
	credentialHandler *handler.DocumentHandler) *Router {
	return &Router{
		authJWTHandler:  authJWTHandler,
		authZkHandler:   authZkHandler,
		issuerHandler:   issuerHandler,
		holderHandler:   holderHandler,
		verifierHandler: verifierHandler,
		documentHandler: credentialHandler,
	}
}

func (r *Router) SetupRoutes(engine *gin.Engine) {
	apiGroup := engine.Group("api/v1")

	SetupAuthJWTRouter(apiGroup, r.authJWTHandler)
	SetupAuthZkRouter(apiGroup, r.authZkHandler)
	SetupIssuerRouter(apiGroup, r.issuerHandler)
	SetupHolderRouter(apiGroup, r.holderHandler)
	SetupVerifierRouter(apiGroup, r.verifierHandler)
	SetupCredentialRouter(apiGroup, r.documentHandler)
}
