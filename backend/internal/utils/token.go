package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const jwtSecret = "secret"

func GenerateToken(role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"role":        role,
		"expire_date": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(jwtSecret))

	return tokenString, err

}

func GetRole(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("token claims error")
	}

	expireFloat, ok := claims["expire_date"].(float64)
	if !ok {
		return "", errors.New("invalid expire_date format")
	}

	if time.Now().After(time.Unix(int64(expireFloat), 0)) {
		return "", errors.New("token expired")
	}

	return claims["role"].(string), nil
}
