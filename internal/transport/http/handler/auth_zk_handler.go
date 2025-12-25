package handler

import (
	"be/internal/service"
	response "be/internal/shared/helper"
	"be/internal/transport/http/dto"
	"be/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/iden3/iden3comm/v2/protocol"
)

type AuthZkHandler struct {
	authZkService service.IAuthZkService
	logger        *logger.ZapLogger
}

func NewAuthZkHandler(logger *logger.ZapLogger, authZkService service.IAuthZkService) *AuthZkHandler {
	return &AuthZkHandler{authZkService: authZkService, logger: logger}
}

func (h *AuthZkHandler) Register(c *gin.Context) {
	var request dto.IdentityCreatedRequestDto

	if err := c.ShouldBindJSON(&request); err != nil {
		response.RespondError(c, err)
	}
	identity, err := h.authZkService.Register(c.Request.Context(), &request)
	if err != nil {
		response.RespondError(c, err)
	}
	response.RespondSuccess(c, identity)
}

func (h *AuthZkHandler) Login(c *gin.Context) {
	var authResponse protocol.AuthorizationResponseMessage
	if err := c.ShouldBindJSON(&authResponse); err != nil {
		response.RespondError(c, err)
	}

	identity, err := h.authZkService.Login(c.Request.Context(), &authResponse)
	if err != nil {
		response.RespondError(c, err)
	}
	response.RespondSuccess(c, identity)
}

func (h *AuthZkHandler) Challenge(c *gin.Context) {
	res, err := h.authZkService.Challenge(c.Request.Context())
	if err != nil {
		response.RespondError(c, err)
	}
	response.RespondSuccess(c, res)
}
