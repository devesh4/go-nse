package handler

import (
	as "github.com/devesh44/gin-poc/service/auth"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	as *as.AuthS
}

func NewAuthHandler() *AuthHandler {
	as := as.NewAuthService()
	return &AuthHandler{as: as}
}

func (a *AuthHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {}
}
