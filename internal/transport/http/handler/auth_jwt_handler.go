package handler

import (
	"be/internal/service"
	"be/internal/shared/constant"
	response "be/internal/shared/helper"
	"be/internal/transport/http/dto"
	"be/pkg/logger"

	"github.com/gin-gonic/gin"
)

type AuthJWTHandler struct {
	authService service.IAuthJWTService
	logger      *logger.ZapLogger
}

func NewAuthJWTHandler(as service.IAuthJWTService, logger *logger.ZapLogger) *AuthJWTHandler {
	return &AuthJWTHandler{authService: as, logger: logger}
}

func (h *AuthJWTHandler) GetAllUser(c *gin.Context) {
	users, err := h.authService.GetAllUsers(c.Request.Context())
	if err != nil {
		response.RespondError(c, err)
	}
	response.RespondSuccess(c, users)
}

func (h *AuthJWTHandler) GetProfile(c *gin.Context) {
	id := c.Param("id")
	user, err := h.authService.GetProfile(c.Request.Context(), id)
	if err != nil {
		response.RespondError(c, err)
	}
	response.RespondSuccess(c, user)
}

func (h *AuthJWTHandler) UpdateProfile(c *gin.Context) {
	id := c.Param("id")
	var userRequest dto.UserRequest
	if err := c.ShouldBindJSON(userRequest); err != nil {
		response.RespondError(c, &constant.BadRequest)
	}

	userResponse, err := h.authService.UpdateProfile(c.Request.Context(), id, &userRequest)

	if err != nil {
		response.RespondError(c, err)
	}
	response.RespondSuccess(c, userResponse)
}

func (h *AuthJWTHandler) Register(c *gin.Context) {
	var registerRequest dto.RegisterRequest
	if err := c.ShouldBindJSON(&registerRequest); err != nil {
		response.RespondError(c, &constant.BadRequest)
	}

	accessToken, refreshToken, err := h.authService.Register(c.Request.Context(), registerRequest.Email, registerRequest.Password, registerRequest.Name)
	if err != nil {
		response.RespondError(c, err)
	}
	response.RespondSuccess(c, dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func (h *AuthJWTHandler) Login(c *gin.Context) {
	var loginRequest dto.LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		response.RespondError(c, &constant.BadRequest)
	}

	accessToken, refreshToken, err := h.authService.Login(c.Request.Context(), loginRequest.Email, loginRequest.Password)
	if err != nil {
		response.RespondError(c, err)
	}
	response.RespondSuccess(c, dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func (h *AuthJWTHandler) RefreshToken(c *gin.Context) {
	tokenString := c.GetString("token")
	if tokenString == "" {
		response.RespondError(c, &constant.InvalidToken)
	}
	accessToken, refreshToken, err := h.authService.RefreshToken(c.Request.Context(), tokenString)
	if err != nil {
		response.RespondError(c, err)
	}
	response.RespondSuccess(c, dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func (h *AuthJWTHandler) GetAuthService() service.IAuthJWTService {
	return h.authService
}
