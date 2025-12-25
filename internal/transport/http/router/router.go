package router

import (
	"be/internal/transport/http/handler"

	"github.com/gin-gonic/gin"
)

type Router struct {
	authJWTHandler  *handler.AuthJWTHandler
	authZkHandler   *handler.AuthZkHandler
	documentHandler *handler.DocumentHandler
}

func NewRouter(
	authJWTHandler *handler.AuthJWTHandler,
	authZkHandler *handler.AuthZkHandler,
	credentialHandler *handler.DocumentHandler) *Router {
	return &Router{
		authJWTHandler:  authJWTHandler,
		authZkHandler:   authZkHandler,
		documentHandler: credentialHandler,
	}
}

func (r *Router) SetupRoutes(engine *gin.Engine) {
	apiGroup := engine.Group("api/v1")

	SetupAuthJWTRouter(apiGroup, r.authJWTHandler)
	SetupAuthZkRouter(apiGroup, r.authZkHandler)
	SetupDocumentRouter(apiGroup, r.documentHandler)
}
