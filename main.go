package main

import (
	"log"
	"net/http"
	"strings"
	"twostep-backend/auth"
	"twostep-backend/handler"
	"twostep-backend/helper"
	"twostep-backend/user"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:@tcp(127.0.0.1:3306)/twostep_db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err.Error())
	}

	db.AutoMigrate(&user.User{})

	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)
	authService := auth.NewService()
	userHandler := handler.NewUserHandler(userService, authService)

	router := gin.Default()

	v1 := router.Group("api/v1")
	v1.POST("/users/auth/register", userHandler.RegisterUser)
	v1.POST("/users/auth/login", userHandler.LoginHandler)
	v1.GET("/users/auth/me", authMiddleware(userService, authService), userHandler.AuthMe)
	v1.POST("/users/avatars", authMiddleware(userService, authService), userHandler.UploadAvatar)

	router.Run("localhost:5000")
}

func authMiddleware(userService user.Service, authService auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if !strings.Contains(authHeader, "Bearer") {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		arrayToken := strings.Split(authHeader, " ")
		tokenString := ""
		//Berhasil Split

		if len(arrayToken) == 2 {
			tokenString = arrayToken[1]
		}
		validatedToken, err := authService.ValidateToken(tokenString)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		claim, ok := validatedToken.Claims.(jwt.MapClaims)
		if !ok || !validatedToken.Valid {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		userID := int(claim["user_id"].(float64))

		//Check User ada di database
		user, err := userService.GetUserByID(userID)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		c.Set("currentUser", user)
	}
}
