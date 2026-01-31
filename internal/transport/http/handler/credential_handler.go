package handler

import (
	"be/internal/service"
	"be/internal/shared/constant"
	"be/internal/shared/helper"
	"be/internal/transport/http/dto"
	"fmt"

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
		helper.RespondError(c, err)
		return
	}
	credential, err := h.credentialService.CreateCredentialRequest(c.Request.Context(), &request)
	if err != nil {
		helper.RespondError(c, err)
		return
	}
	helper.RespondSuccess(c, credential)
}

func (h *CredentialHandler) GetCredentialRequests(c *gin.Context) {
	user, ok := c.Get("user")
	if !ok {
		helper.RespondError(c, &constant.InternalServer)
		return
	}

	claims, ok := user.(*dto.ZKClaims)
	if !ok {
		helper.RespondError(c, &constant.InternalServer)
		return
	}

	credentials, err := h.credentialService.GetCredentialRequests(c.Request.Context(), claims)
	if err != nil {
		helper.RespondError(c, err)
		return
	}
	helper.RespondSuccess(c, credentials)
}

func (h *CredentialHandler) UpdateCredentialRequest(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		helper.RespondError(c, &constant.BadRequest)
		return
	}
	var request dto.CredentialRequestUpdatedRequestDto
	if err := c.ShouldBindJSON(&request); err != nil {
		helper.RespondError(c, err)
		return
	}

	err := h.credentialService.UpdateCredentialRequest(c.Request.Context(), id, &request)
	if err != nil {
		helper.RespondError(c, err)
		return
	}
	helper.RespondSuccess(c, "")
}

func (h *CredentialHandler) GetVerifiableCredentials(c *gin.Context) {
	user, ok := c.Get("user")
	if !ok {
		helper.RespondError(c, &constant.InternalServer)
		return
	}

	claims, ok := user.(*dto.ZKClaims)
	if !ok {
		helper.RespondError(c, &constant.InternalServer)
		return
	}

	res, err := h.credentialService.GetVerifiableCredentials(c.Request.Context(), claims)
	if err != nil {
		helper.RespondError(c, err)
		return
	}

	helper.RespondSuccess(c, res)
}

func (h *CredentialHandler) GetVerifiableCredentialById(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		helper.RespondError(c, &constant.BadRequest)
		return
	}
	res, err := h.credentialService.GetVerifiableCredentialById(c.Request.Context(), id)
	if err != nil {
		helper.RespondError(c, err)
		return
	}
	helper.RespondSuccess(c, res)
}

func (h *CredentialHandler) UpdateVerifiableCredential(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		helper.RespondError(c, &constant.BadRequest)
		return
	}
	var request dto.VerifiableUpdatedRequestDto
	if err := c.ShouldBindJSON(&request); err != nil {
		helper.RespondError(c, &constant.BadRequest)
		return
	}
	err := h.credentialService.UpdateVerifiableCredential(c.Request.Context(), id, &request)
	if err != nil {
		helper.RespondError(c, err)
		return
	}
	helper.RespondSuccess(c, "")
}

func (h *CredentialHandler) IssueVerifiableCredential(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		helper.RespondError(c, &constant.BadRequest)
		return
	}

	var request dto.IssueVerifiableCredentialRequestDto
	if err := c.ShouldBindJSON(&request); err != nil {
		helper.RespondError(c, err)
		return
	}

	res, err := h.credentialService.IssueVerifiableCredential(c.Request.Context(), id, &request)
	if err != nil {
		fmt.Println(err)
		helper.RespondError(c, err)
		return
	}
	helper.RespondSuccess(c, res)

}
