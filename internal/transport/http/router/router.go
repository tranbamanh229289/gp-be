package router

import (
	"be/internal/transport/http/handler"

	"github.com/gin-gonic/gin"
)

type Router struct {
	authJWTHandler *handler.AuthJWTHandler
	authZkHandler  *handler.AuthZkHandler
}

func NewRouter(authJWTHandler *handler.AuthJWTHandler, authZKHandler *handler.AuthZkHandler) *Router {
	return &Router{
		authJWTHandler: authJWTHandler,
		authZkHandler:  authZKHandler,
	}
}

func (r *Router) SetupRoutes(engine *gin.Engine) {
	apiGroup := engine.Group("api/v1")

	SetupAuthJWTRouter(apiGroup, r.authJWTHandler)
	SetupAuthZkRouter(apiGroup, r.authZkHandler)
}
