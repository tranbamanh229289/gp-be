package handler

import (
	"be/internal/service"
	response "be/internal/shared/helper"
	"be/internal/transport/http/dto"

	"github.com/gin-gonic/gin"
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
	var request dto.CredentialRequestCreatedRequestDto
	if err := c.ShouldBindJSON(&request); err != nil {
		response.RespondError(c, err)
	}
	credential, err := h.credentialService.CreateCredentialRequest(c.Request.Context(), &request)
	if err != nil {
		response.RespondError(c, err)
	}
	response.RespondSuccess(c, credential)
}

func (h *CredentialHandler) GetCredentialRequests(c *gin.Context) {
	credentials, err := h.credentialService.GetCredentialRequests(c.Request.Context())
	if err != nil {
		response.RespondError(c, err)
	}
	response.RespondSuccess(c, credentials)
}

func (h *CredentialHandler) IssueCredential(c *gin.Context) {

}
