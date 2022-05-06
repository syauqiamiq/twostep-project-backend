package auth

import (
	"log"
	"os"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
)

type Service interface {
	GenerateToken(userID int) (string, error)
}

type jwtService struct {
}

func NewService() *jwtService {
	return &jwtService{}
}

func (s *jwtService) GenerateToken(userID int) (string, error) {
	claim := jwt.MapClaims{
		"user_id": userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err.Error())
	}
	var secretKey = []byte(os.Getenv("SECRET_KEY"))
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return signedToken, err
	}

	return signedToken, nil

}
