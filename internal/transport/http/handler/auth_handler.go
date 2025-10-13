package handler

import (
	"be/internal/service/auth"
	"be/internal/transport/http/dtos"
	"net/http"

	"github.com/gin-gonic/gin"
)


type AuthHandler struct {
	authService auth.UserService
}

func NewUserHandler(as auth.UserService) *AuthHandler{
	return &AuthHandler{authService: as}
}


func (h *AuthHandler) Register(c *gin.Context) {
	var registerRequest dtos.RegisterRequest;
	if err := c.ShouldBindJSON(&registerRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	user, err := h.authService.Register()
	
	c.JSON(http.StatusOK, nil)
}

func (h *AuthHandler) Login(c *gin.Context) {
	user, err := h.authService.Login()
	c.JSON(http.StatusOK, nil)
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	user, err := h.authService.GetProfile(id)
	c.JSON(http.StatusOK, nil)
}

func (h *AuthHandler) GetAllUser(c *gin.Context) {
	users, err := h.authService.GetAllUsers()
	c.JSON(http.StatusOK, nil)
}