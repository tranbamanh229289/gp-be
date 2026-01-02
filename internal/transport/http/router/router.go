package router

import (
	"be/internal/transport/http/handler"

	"github.com/gin-gonic/gin"
)

type Router struct {
	authJWTHandler    *handler.AuthJWTHandler
	authZkHandler     *handler.AuthZkHandler
	documentHandler   *handler.DocumentHandler
	credentialHandler *handler.CredentialHandler
	schemaHandler     *handler.SchemaHandler
	proofHandler      *handler.ProofHandler
}

func NewRouter(
	authJWTHandler *handler.AuthJWTHandler,
	authZkHandler *handler.AuthZkHandler,
	documentHandler *handler.DocumentHandler,
	credentialHandler *handler.CredentialHandler,
	schemaHandler *handler.SchemaHandler,
	proofHandler *handler.ProofHandler,
) *Router {
	return &Router{
		authJWTHandler:    authJWTHandler,
		authZkHandler:     authZkHandler,
		documentHandler:   documentHandler,
		credentialHandler: credentialHandler,
		schemaHandler:     schemaHandler,
		proofHandler:      proofHandler,
	}
}

func (r *Router) SetupRoutes(engine *gin.Engine) {
	apiGroup := engine.Group("api/v1")

	SetupAuthJWTRouter(apiGroup, r.authJWTHandler)
	SetupAuthZkRouter(apiGroup, r.authZkHandler)
	SetupDocumentRouter(apiGroup, r.documentHandler)
	SetupCredentialRouter(apiGroup, r.credentialHandler)
	SetupSchemaRouter(apiGroup, r.schemaHandler)
	SetupProofRouter(apiGroup, r.proofHandler)
}
