package jwt_tools

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

func NewJWT(SecretKey string, exp time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(exp).Unix()
	claims["iat"] = time.Now().Unix()
	token.Claims = claims

	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
