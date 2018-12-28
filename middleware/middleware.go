package middleware

import (
	"log"
	"net/http"

	"github.com/yjp19871013/RPiService/jwt_tools"

	"github.com/gin-gonic/gin"
)

const (
	ContextUserKey = "user"
)

func JWTValidateMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := jwt_tools.GetJWTUser(c.Request)
		if err != nil {
			log.Println("jwt validate error")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set(ContextUserKey, user)
		c.Next()
	}
}
