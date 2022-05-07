package main

import (
	"fmt"
	"log"
	"os"
	"twostep-backend/auth"
	"twostep-backend/user"

	"twostep-backend/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbDatabase := os.Getenv("DB_DATABASE")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	//MYSQL SETUP
	// dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUsername, dbPassword, dbHost, dbPort, dbDatabase)
	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	//POSTGRESQL SETUP PRODUCTION
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=require TimeZone=Asia/Shanghai", dbHost, dbUsername, dbPassword, dbDatabase, dbPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
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

	SERVER := os.Getenv("SERVER_PORT")

	router.Run(SERVER)
}
