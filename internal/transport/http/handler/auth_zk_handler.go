package handler

import (
	"be/config"
	"be/internal/service"
	"be/internal/shared/constant"
	response "be/internal/shared/helper"
	"be/internal/transport/http/dto"
	"be/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/iden3/iden3comm/v2/protocol"
)

type AuthZkHandler struct {
	authZkService service.IAuthZkService
	config        *config.Config
	logger        *logger.ZapLogger
}

func NewAuthZkHandler(config *config.Config, logger *logger.ZapLogger, authZkService service.IAuthZkService) *AuthZkHandler {
	return &AuthZkHandler{config: config, logger: logger, authZkService: authZkService}
}

func (h *AuthZkHandler) Register(c *gin.Context) {
	var request dto.IdentityCreatedRequestDto

	if err := c.ShouldBindJSON(&request); err != nil {
		response.RespondError(c, err)
		return
	}
	identity, err := h.authZkService.Register(c.Request.Context(), &request)
	if err != nil {
		response.RespondError(c, err)
		return
	}
	response.RespondSuccess(c, identity)
}

func (h *AuthZkHandler) Challenge(c *gin.Context) {
	res, err := h.authZkService.Challenge(c.Request.Context())
	if err != nil {
		response.RespondError(c, err)
		return
	}
	response.RespondSuccess(c, res)
}

func (h *AuthZkHandler) Login(c *gin.Context) {
	var authResponse protocol.AuthorizationResponseMessage
	if err := c.ShouldBindJSON(&authResponse); err != nil {
		response.RespondError(c, err)
		return
	}

	res, err := h.authZkService.Login(c.Request.Context(), &authResponse)
	if err != nil {
		response.RespondError(c, err)
		return
	}
	age := h.config.JWT.AccessTokenTTL.Milliseconds()
	c.SetCookie("refreshToken", res.RefreshToken, int(age), "/", "", false, true)
	response.RespondSuccess(c, res)
}

func (h *AuthZkHandler) RefreshZKToken(c *gin.Context) {
	refreshToken, err := c.Cookie("refreshToken")
	res, err := h.authZkService.RefreshZKToken(c.Request.Context(), refreshToken)
	if err != nil {
		response.RespondError(c, err)
		return
	}
	age := h.config.JWT.AccessTokenTTL.Milliseconds()
	c.SetCookie("refreshToken", res.RefreshToken, int(age), "/", "", false, true)
	response.RespondSuccess(c, res)
}

func (h *AuthZkHandler) Logout(c *gin.Context) {
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

	err := h.authZkService.Logout(c.Request.Context(), claims.ID)
	if err != nil {
		response.RespondError(c, err)
		return
	}
	response.RespondSuccess(c, "")
}

func (h *AuthZkHandler) GetIdentityByDID(c *gin.Context) {
	did := c.Param("did")
	if did == "" {
		response.RespondError(c, &constant.BadRequest)
	}
	res, err := h.authZkService.GetIdentityByDID(c.Request.Context(), did)
	if err != nil {
		response.RespondError(c, err)
		return
	}
	response.RespondSuccess(c, res)
}

func (h *AuthZkHandler) GetIdentityByRole(c *gin.Context) {
	role := c.Query("role")
	if role == "" {
		response.RespondError(c, &constant.BadRequest)
	}
	res, err := h.authZkService.GetIdentityByRole(c.Request.Context(), role)
	if err != nil {
		response.RespondError(c, err)
		return
	}
	response.RespondSuccess(c, res)
}
