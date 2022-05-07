package handler

import (
	"net/http"
	"twostep-backend/auth"
	"twostep-backend/helper"
	"twostep-backend/user"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	service     user.Service
	authService auth.Service
}

func NewUserHandler(userService user.Service, authService auth.Service) *userHandler {
	return &userHandler{userService, authService}
}

func (h *userHandler) RegisterUser(c *gin.Context) {
	var input user.RegisterUserInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errorMessage := gin.H{
			"error": "INPUT FORMAT WRONG",
		}
		response := helper.APIResponse("Input Format Wrong", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	newUser, err := h.service.RegisterUser(input)
	if err != nil {
		errorMessage := gin.H{
			"error": err.Error(),
		}
		response := helper.APIResponse("Register Failed", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	//Generate Token
	userID := newUser.ID
	token, err := h.authService.GenerateToken(userID)
	if err != nil {
		errorMessage := gin.H{
			"error": err.Error(),
		}
		response := helper.APIResponse("Register Failed", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H{
		"id":    newUser.ID,
		"name":  newUser.Name,
		"email": newUser.Email,
		"token": token,
	}
	response := helper.APIResponse("Register Success", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)

}

func (h *userHandler) LoginHandler(c *gin.Context) {
	var input user.LoginInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errorMessage := gin.H{
			"error": "INPUT FORMAT WRONG",
		}
		response := helper.APIResponse("Input Format Wrong", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	loggedInUser, err := h.service.LoginUser(input)
	if err != nil {
		errorMessage := gin.H{
			"error": err.Error(),
		}
		response := helper.APIResponse("Input Format Wrong", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	//Generate Token
	userID := loggedInUser.ID
	token, err := h.authService.GenerateToken(userID)
	if err != nil {
		errorMessage := gin.H{
			"error": err.Error(),
		}
		response := helper.APIResponse("Register Failed", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H{
		"id":    loggedInUser.ID,
		"name":  loggedInUser.Name,
		"email": loggedInUser.Email,
		"role":  loggedInUser.Role,
		"token": token,
	}

	response := helper.APIResponse("Login Success", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}
