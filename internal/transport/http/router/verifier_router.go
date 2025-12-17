package router

import (
	"be/internal/transport/http/handler"

	"github.com/gin-gonic/gin"
)

func SetupVerifierRouter(apiGroup *gin.RouterGroup, verifierHandler *handler.VerifierHandler) {
	verifierGroup := apiGroup.Group("verifier")

	verifierGroup.POST("verify", verifierHandler.Verify)
}
