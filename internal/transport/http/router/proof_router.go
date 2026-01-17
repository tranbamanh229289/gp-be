package router

import (
	"be/internal/shared/constant"
	"be/internal/transport/http/handler"
	"be/internal/transport/http/middleware"

	"github.com/gin-gonic/gin"
)

func (r *Router) SetupProofRouter(apiGroup *gin.RouterGroup, proofHandler *handler.ProofHandler) {
	proofGroup := apiGroup.Group("proofs")
	proofGroup.Use(middleware.AuthenticateMiddleware(r.authZkService))

	proofRequestGroup := proofGroup.Group("request")
	proofResponseGroup := proofGroup.Group("response")

	proofRequestGroup.POST("", middleware.AuthorizeMiddleware([]constant.IdentityRole{constant.IdentityVerifierRole}), proofHandler.CreateProofRequest)
	proofRequestGroup.GET("", proofHandler.GetProofRequests)
	proofRequestGroup.PATCH("/:id", proofHandler.UpdateProofRequest)

	proofResponseGroup.POST("", proofHandler.CreateProofResponse)
	proofResponseGroup.GET("", proofHandler.GetProofResponses)
}
