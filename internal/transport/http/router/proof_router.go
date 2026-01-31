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

	proofRequestGroup := proofGroup.Group("requests")
	proofSubmissionGroup := proofGroup.Group("submissions")

	proofRequestGroup.POST("", middleware.AuthorizeMiddleware([]constant.IdentityRole{constant.IdentityVerifierRole}), proofHandler.CreateProofRequest)
	proofRequestGroup.PATCH("/:id", proofHandler.UpdateProofRequest)
	proofRequestGroup.GET("", proofHandler.GetProofRequests)

	proofSubmissionGroup.POST("", proofHandler.CreateProofSubmission)
	proofSubmissionGroup.PATCH("/:id", proofHandler.VerifyZKProof)
	proofSubmissionGroup.GET("", proofHandler.GetProofSubmissions)
}
