package handler

import (
	"be/internal/service"
	"be/internal/transport/http/dtos"
	"be/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.IAuthService
	logger *logger.ZapLogger
}

func NewAuthHandler(as *service.IAuthService, logger *logger.ZapLogger) *AuthHandler{
	return &AuthHandler{authService: as, logger: logger}
}


func (h *AuthHandler) Register(c *gin.Context, email, password, name string) {
	var registerRequest dtos.RegisterRequest;
	if err := c.ShouldBindJSON(&registerRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	
	c.JSON(http.StatusOK, nil)
}
