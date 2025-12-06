package handler

import (
	"be/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthZkHandler struct {
	authService service.IAuthZkService
}

func NewAuthZkHandler(as service.IAuthZkService) *AuthZkHandler {
	return &AuthZkHandler{authService: as}
}

func (h *AuthZkHandler) Signin(c *gin.Context) {

}

func (h *AuthZkHandler) Callback(c *gin.Context) {

}
