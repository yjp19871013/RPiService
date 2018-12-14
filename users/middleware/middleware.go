package midware

import (
	"net/http"

	"github.com/yjp19871013/RPiService/users/model"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gin-gonic/gin"
)

const (
	secretKey = "rpi service jwt"
)

func JWTValidateMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := request.ParseFromRequest(c.Request, request.AuthorizationHeaderExtractor,
			func(token *jwt.Token) (interface{}, error) {
				return []byte(secretKey), nil
			})
		if err == nil {
			if token.Valid {
				c.Next()
			} else {
				c.JSON(http.StatusUnauthorized, model.LoginResponse{"无效的token"})
				c.Abort()
			}
		} else {
			c.JSON(http.StatusUnauthorized, model.LoginResponse{err.Error()})
			c.Abort()
		}
	}
}
