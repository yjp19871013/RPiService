package midware

import (
	"net/http"

	"github.com/yjp19871013/RPiService/users/entities"

	"github.com/yjp19871013/RPiService/users/db"

	"github.com/yjp19871013/RPiService/users/settings"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gin-gonic/gin"
)

func JWTValidateMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := request.ParseFromRequest(c.Request, request.AuthorizationHeaderExtractor,
			func(token *jwt.Token) (interface{}, error) {
				return []byte(settings.SecretKey), nil
			})
		if err == nil {
			if token.Valid {
				var user db.User
				db.GetInstance().Where("token = ?", token.Raw).First(&user)
				if user.ID == 0 {
					c.JSON(http.StatusUnauthorized, entities.TokenResponse{"无效的token"})
					c.Abort()
				}

				c.Next()
			} else {
				c.JSON(http.StatusUnauthorized, entities.TokenResponse{"无效的token"})
				c.Abort()
			}
		} else {
			c.JSON(http.StatusUnauthorized, entities.TokenResponse{err.Error()})
			c.Abort()
		}
	}
}
