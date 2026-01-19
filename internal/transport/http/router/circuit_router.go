package router

import (
	"be/internal/transport/http/handler"
	"be/internal/transport/http/middleware"

	"github.com/gin-gonic/gin"
)

func (r *Router) SetupCircuitRouter(apiGroup *gin.RouterGroup, circuitHandler *handler.CircuitHandler) {
	circuitGroup := apiGroup.Group("circuits")
	circuitGroup.Use(middleware.AuthenticateMiddleware(r.authZkService))
	circuitGroup.POST("credentialAtomicQueryV3", circuitHandler.GenerateCredentialAtomicQueryV3Input)
}
