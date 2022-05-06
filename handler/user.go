package handler

import (
	"twostep-backend/user"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	service user.Service
}

func NewUserHandler(userService user.Service) *userHandler {
	return &userHandler{userService}
}

func (s *userHandler) RegisterUser(c *gin.Context) {

}
