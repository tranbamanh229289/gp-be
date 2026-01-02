package router

import (
	"be/internal/transport/http/handler"

	"github.com/gin-gonic/gin"
)

func SetupProofRouter(apiGroup *gin.RouterGroup, proofHandler *handler.ProofHandler) {
	proofGroup := apiGroup.Group("proofs")
	proofRequestGroup := proofGroup.Group("request")
	proofResponseGroup := proofGroup.Group("response")

	proofRequestGroup.POST("", proofHandler.CreateProofRequest)
	proofRequestGroup.GET("", proofHandler.GetProofRequests)
	proofRequestGroup.PATCH("/:id", proofHandler.UpdateProofRequest)

	proofResponseGroup.POST("", proofHandler.CreateProofResponse)
	proofResponseGroup.GET("", proofHandler.GetProofResponses)
}
