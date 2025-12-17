package handler

import (
	"be/internal/service"
	response "be/internal/shared/helper"
	"be/pkg/logger"

	"github.com/gin-gonic/gin"
)

type AuthZkHandler struct {
	authService service.IAuthZkService
	logger      *logger.ZapLogger
}

func NewAuthZkHandler(as service.IAuthZkService, logger *logger.ZapLogger) *AuthZkHandler {
	return &AuthZkHandler{authService: as, logger: logger}
}

func (h *AuthZkHandler) Signin(c *gin.Context) {
	request, err := h.authService.CreateAuthRequest(c.Request.Context())
	if err != nil {
		response.RespondError(c, err)
	}
	response.RespondSuccess(c, request)
}

func (h *AuthZkHandler) Callback(c *gin.Context) {

}
