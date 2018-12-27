package jwt_tools

import (
	"fmt"
	"net/http"
	"time"

	"github.com/yjp19871013/RPiService/db"
	"github.com/yjp19871013/RPiService/users/settings"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
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

func IsJWTValidate(req *http.Request) bool {
	token, err := request.ParseFromRequest(req, request.AuthorizationHeaderExtractor,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(settings.SecretKey), nil
		})
	if err != nil || !token.Valid {
		return false
	}

	_, err = db.FindUserByToken(token.Raw)
	if err != nil {
		return false
	}

	return true
}

func GetJWTUser(req *http.Request) (*db.User, error) {
	token, err := request.ParseFromRequest(req, request.AuthorizationHeaderExtractor,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(settings.SecretKey), nil
		})
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("token parse error")
	}

	user, err := db.FindUserByToken(token.Raw)
	if err != nil {
		return nil, err
	}

	return user, nil
}
