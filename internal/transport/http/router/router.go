package router

import (
	"be/internal/infrastructure/database/postgres"
	"be/internal/service"
	"be/internal/transport/http/handler"

	"github.com/gin-gonic/gin"
)

type Router struct {
	db                *postgres.PostgresDB
	authJWTHandler    *handler.AuthJWTHandler
	authZkHandler     *handler.AuthZkHandler
	documentHandler   *handler.DocumentHandler
	credentialHandler *handler.CredentialHandler
	schemaHandler     *handler.SchemaHandler
	proofHandler      *handler.ProofHandler
	circuitHandler    *handler.CircuitHandler
	statisticHandler  *handler.StatisticHandler
	authZkService     service.IAuthZkService
}

func NewRouter(
	db *postgres.PostgresDB,
	authJWTHandler *handler.AuthJWTHandler,
	authZkHandler *handler.AuthZkHandler,
	documentHandler *handler.DocumentHandler,
	credentialHandler *handler.CredentialHandler,
	schemaHandler *handler.SchemaHandler,
	proofHandler *handler.ProofHandler,
	circuitHandler *handler.CircuitHandler,
	statisticHandler *handler.StatisticHandler,
	authZkService service.IAuthZkService,
) *Router {
	return &Router{
		db:                db,
		authJWTHandler:    authJWTHandler,
		authZkHandler:     authZkHandler,
		documentHandler:   documentHandler,
		credentialHandler: credentialHandler,
		schemaHandler:     schemaHandler,
		proofHandler:      proofHandler,
		circuitHandler:    circuitHandler,
		statisticHandler:  statisticHandler,
		authZkService:     authZkService,
	}
}

func (r *Router) SetupRoutes(engine *gin.Engine) {
	apiGroup := engine.Group("api/v1")
	r.SetupAuthJWTRouter(apiGroup, r.authJWTHandler)
	r.SetupAuthZkRouter(apiGroup, r.authZkHandler)
	r.SetupDocumentRouter(apiGroup, r.documentHandler)
	r.SetupCredentialRouter(apiGroup, r.credentialHandler, r.db)
	r.SetupSchemaRouter(apiGroup, r.schemaHandler, r.db)
	r.SetupProofRouter(apiGroup, r.proofHandler)
	r.SetupCircuitRouter(apiGroup, r.circuitHandler)
	r.SetupStatisticRouter(apiGroup, r.statisticHandler)
}
