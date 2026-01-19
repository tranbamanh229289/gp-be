package handler

import (
	"be/config"
	"be/internal/service"
	"be/internal/shared/constant"
	"be/internal/shared/helper"
	"be/internal/transport/http/dto"
	"be/pkg/logger"
	"fmt"

	"github.com/gin-gonic/gin"
)

type CircuitHandler struct {
	circuitService service.ICircuitService
	config         *config.Config
	logger         *logger.ZapLogger
}

func NewCircuitHandler(
	circuitService service.ICircuitService,
	config *config.Config,
	logger *logger.ZapLogger) *CircuitHandler {
	return &CircuitHandler{
		circuitService: circuitService,
		config:         config,
		logger:         logger,
	}
}

func (h *CircuitHandler) GenerateCredentialAtomicQueryV3Input(c *gin.Context) {
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

	var request dto.CredentialAtomicQueryV3InputRequestDto
	if err := c.ShouldBindJSON(&request); err != nil {
		helper.RespondError(c, err)
		return
	}
	res, err := h.circuitService.GetCredentialAtomicQueryV3Input(c.Request.Context(), &request, claims)
	if err != nil {
		fmt.Println(err)
		helper.RespondError(c, err)
		return
	}
	helper.RespondSuccess(c, res)
}
