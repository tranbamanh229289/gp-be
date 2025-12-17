package router

import (
	"be/internal/transport/http/handler"

	"github.com/gin-gonic/gin"
)

func SetupIssuerRouter(apiGroup *gin.RouterGroup, issuerHandler *handler.IssuerHandler) {
	issuerGroup := apiGroup.Group("issuer")

	issuerGroup.POST("issue-claim", issuerHandler.IssueClaim)
}
