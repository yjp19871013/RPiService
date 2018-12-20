package middleware

import (
	"net/http"

	"github.com/yjp19871013/RPiService/users/dto"

	"github.com/yjp19871013/RPiService/users/jwt_tools"

	"github.com/gin-gonic/gin"
)

func JWTValidateMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !jwt_tools.IsJWTValidate(c.Request) {
			c.JSON(http.StatusUnauthorized, dto.TokenResponse{"无效的token"})
			c.Abort()
		}

		c.Next()
	}
}
