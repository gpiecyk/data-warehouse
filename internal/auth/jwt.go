package auth

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var (
	SecretKey = []byte("secret_golang_app")
)

type Claims struct {
	id int
	jwt.StandardClaims
}

func GenerateToken(id int) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = id
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	tokenString, err := token.SignedString(SecretKey)
	if err != nil {
		log.Fatal("Error in Generating key")
		return "", err
	}
	return tokenString, nil
}

func ParseToken(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		id := int(claims["id"].(float64))
		return id, nil
	} else {
		return -1, err
	}
}
