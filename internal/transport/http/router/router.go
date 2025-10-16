package router

import (
	"be/internal/transport/http/handler"

	"github.com/gin-gonic/gin"
)

type Router struct {
	authHandler *handler.AuthHandler
}

func NewRouter(authHandler *handler.AuthHandler) *Router {
	return &Router{
		authHandler: authHandler,
	}
}

func (r *Router) SetupRoutes(engine *gin.Engine) {
	apiGroup := engine.Group("api/v1")

	SetupAuthRouter(apiGroup, r.authHandler)
}