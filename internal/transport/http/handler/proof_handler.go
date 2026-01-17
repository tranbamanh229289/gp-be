package handler

import (
	"be/internal/service"
	"be/internal/shared/constant"
	response "be/internal/shared/helper"
	"be/internal/transport/http/dto"

	"github.com/gin-gonic/gin"
	"github.com/iden3/iden3comm/v2/protocol"
)

type ProofHandler struct {
	proofService service.IProofService
}

func NewProofHandler(proofService service.IProofService) *ProofHandler {
	return &ProofHandler{
		proofService: proofService,
	}
}

func (h *ProofHandler) CreateProofRequest(c *gin.Context) {
	var request protocol.AuthorizationRequestMessage
	if err := c.ShouldBindJSON(&request); err != nil {
		response.RespondError(c, err)
		return
	}
	user, ok := c.Get("user")
	if !ok {
		response.RespondError(c, &constant.InternalServer)
		return
	}
	claims, ok := user.(*dto.ZKClaims)
	if !ok {
		response.RespondError(c, &constant.InternalServer)
		return
	}

	if claims.DID != request.From {
		response.RespondError(c, &constant.BadRequest)
	}

	res, err := h.proofService.CreateProofRequest(c.Request.Context(), &request)

	if err != nil {
		response.RespondError(c, err)
		return
	}

	response.RespondSuccess(c, res)
}

func (h *ProofHandler) GetProofRequests(c *gin.Context) {
	user, ok := c.Get("user")
	if !ok {
		response.RespondError(c, &constant.InternalServer)
		return
	}
	claims, ok := user.(*dto.ZKClaims)
	if !ok {
		response.RespondError(c, &constant.InternalServer)
		return
	}

	res, err := h.proofService.GetProofRequests(c.Request.Context(), claims)
	if err != nil {
		response.RespondError(c, err)
		return
	}

	response.RespondSuccess(c, res)
}

func (h *ProofHandler) UpdateProofRequest(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.RespondError(c, &constant.BadRequest)
		return
	}
	var request dto.ProofRequestUpdatedRequestDto
	if err := c.ShouldBindJSON(&request); err != nil {
		response.RespondError(c, err)
		return
	}
	err := h.proofService.UpdateProofRequest(c.Request.Context(), id, &request)
	if err != nil {
		response.RespondError(c, err)
		return
	}

	response.RespondSuccess(c, "")
}

func (h *ProofHandler) CreateProofResponse(c *gin.Context) {
	var request protocol.AuthorizationResponseMessage
	if err := c.ShouldBindJSON(&request); err != nil {
		response.RespondError(c, err)
		return
	}
	res, err := h.proofService.CreateProofResponse(c.Request.Context(), &request)
	if err != nil {
		response.RespondError(c, err)
		return
	}

	response.RespondSuccess(c, res)
}

func (h *ProofHandler) GetProofResponses(c *gin.Context) {
	res, err := h.proofService.GetProofResponses(c.Request.Context())
	if err != nil {
		response.RespondError(c, err)
		return
	}

	response.RespondSuccess(c, res)
}
