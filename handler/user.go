package handler

import (
	"fmt"
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

func (s *userHandler) AuthMe(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(user.User)
	userId := currentUser.ID
	user, err := s.service.GetUserByID(userId)
	if err != nil {
		errorMessage := gin.H{
			"error": err.Error(),
		}
		response := helper.APIResponse("Authentication Failed", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
		"role":  user.Role,
	}
	response := helper.APIResponse("Authenticated", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}

func (h *userHandler) UploadAvatar(c *gin.Context) {
	file, err := c.FormFile("avatar")
	if err != nil {
		errorMessage := gin.H{
			"is_uploaded": false,
		}
		response := helper.APIResponse("Upload Failed", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)
	userID := currentUser.ID
	path := fmt.Sprintf("images/%d-%s", userID, file.Filename)

	//SAVE TO DIRECTORY
	err = c.SaveUploadedFile(file, path)
	if err != nil {
		errorMessage := gin.H{
			"is_uploaded": false,
		}
		response := helper.APIResponse("Upload Failed", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	//SAVE PATH TO DATABASE
	_, err = h.service.SaveAvatar(userID, path)
	if err != nil {
		errorMessage := gin.H{
			"is_uploaded": false,
		}
		response := helper.APIResponse("Upload Failed", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	data := gin.H{
		"is_uploaded": true,
	}

	//SUCCESSFULLY UPLOAD
	response := helper.APIResponse("Upload Success", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)

}
