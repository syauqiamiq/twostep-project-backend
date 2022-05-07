package main

import (
	"log"
	"twostep-backend/auth"
	"twostep-backend/user"

	"twostep-backend/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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

	router := gin.Default()

	router.Use(cors.Default())

	v1 := router.Group("api/v1")

	routes.SetupV1(v1, userService, authService)

	router.Run("localhost:5000")
}
