package jwt

import (
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	SecretKey = []byte("secret")
)

func GenerateAccessToken(nickname, user_id string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user_id
	claims["nickname"] = nickname
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	tokenString, err := token.SignedString(SecretKey)
	if err != nil {
		log.Fatal("Error in Generating key")
		return "", err
	}
	return tokenString, nil
}

func GenerateRefreshToken(nickname, email, user_id string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["nickname"] = nickname
	claims["email"] = email
	claims["user_id"] = user_id
	claims["exp"] = time.Now().Add(time.Hour * 24 * 7 * 4).Unix()
	tokenString, err := token.SignedString(SecretKey)
	if err != nil {
		log.Fatal("Error in Generating key")
		return "", err
	}
	return tokenString, nil
}

func ParseToken(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username := claims["user_id"].(string)
		return username, nil
	}
	return "", err
}
