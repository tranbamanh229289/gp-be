package handler

import (
	"be/internal/service"
	"be/internal/shared/constant"
	response "be/internal/shared/helper"
	"be/internal/transport/http/dto"

	"github.com/gin-gonic/gin"
	"github.com/iden3/iden3comm/v2/protocol"
)

type CredentialHandler struct {
	credentialService service.ICredentialService
}

func NewCredentialHandler(credentialService service.ICredentialService) *CredentialHandler {
	return &CredentialHandler{
		credentialService: credentialService,
	}
}

func (h *CredentialHandler) CreateCredentialRequest(c *gin.Context) {
	var request protocol.CredentialIssuanceRequestMessage
	if err := c.ShouldBindJSON(&request); err != nil {
		response.RespondError(c, err)
		return
	}
	credential, err := h.credentialService.CreateCredentialRequest(c.Request.Context(), &request)
	if err != nil {
		response.RespondError(c, err)
		return
	}
	response.RespondSuccess(c, credential)
}

func (h *CredentialHandler) GetCredentialRequests(c *gin.Context) {
	credentials, err := h.credentialService.GetCredentialRequests(c.Request.Context())
	if err != nil {
		response.RespondError(c, err)
		return
	}
	response.RespondSuccess(c, credentials)
}

func (h *CredentialHandler) UpdateCredentialRequest(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.RespondError(c, &constant.BadRequest)
		return
	}
	var request dto.CredentialRequestUpdatedRequestDto
	if err := c.ShouldBindJSON(&request); err != nil {
		response.RespondError(c, err)
	}

	err := h.credentialService.UpdateCredentialRequest(c.Request.Context(), id, &request)
	if err != nil {
		response.RespondError(c, err)
		return
	}
	response.RespondSuccess(c, "")
}

func (h *CredentialHandler) GetVerifiableCredentials(c *gin.Context) {
	res, err := h.credentialService.GetVerifiableCredentials(c.Request.Context())
	if err != nil {
		response.RespondError(c, err)
		return
	}
	response.RespondSuccess(c, res)
}

func (h *CredentialHandler) GetVerifiableCredential(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.RespondError(c, &constant.BadRequest)
		return
	}
	res, err := h.credentialService.GetVerifiableCredential(c.Request.Context(), id)
	if err != nil {
		response.RespondError(c, err)
		return
	}
	response.RespondSuccess(c, res)
}

func (h *CredentialHandler) UpdateVerifiableCredential(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.RespondError(c, &constant.BadRequest)
		return
	}
	res, err := h.credentialService.GetVerifiableCredential(c.Request.Context(), id)
	if err != nil {
		response.RespondError(c, err)
		return
	}
	response.RespondSuccess(c, res)
}

func (h *CredentialHandler) SignVerifiableCredential(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.RespondError(c, &constant.BadRequest)
		return
	}
	var request *dto.SignCredentialRequestDto
	if err := c.ShouldBindJSON(&request); err != nil {
		response.RespondError(c, err)
		return
	}

	err := h.credentialService.SignVerifiableCredential(c.Request.Context(), id, request)
	if err != nil {
		response.RespondError(c, err)
		return
	}
	response.RespondSuccess(c, "")
}

func (h *CredentialHandler) IssueVerifiableCredential(c *gin.Context) {
	id := c.Param("id")

	var request dto.IssueVerifiableCredentialRequestDto
	if err := c.ShouldBindJSON(&request); err != nil {
		response.RespondError(c, err)
		return
	}
	res, err := h.credentialService.IssueVerifiableCredential(c.Request.Context(), id, &request)
	if err != nil {
		response.RespondError(c, err)
		return
	}
	response.RespondSuccess(c, res)

}
