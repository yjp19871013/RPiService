package handler

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/yjp19871013/RPiService/users/entities"

	"github.com/yjp19871013/RPiService/users/jwt_tools"

	"github.com/yjp19871013/RPiService/users/settings"

	"github.com/yjp19871013/RPiService/utils"

	"github.com/yjp19871013/RPiService/users/db"

	"github.com/gin-gonic/gin"
)

func CreateToken(c *gin.Context) {
	var request entities.CreateTokenRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusUnauthorized, entities.TokenResponse{Token: ""})
		return
	}

	encodePwd := utils.MD5(request.Password)

	var user = db.User{}
	db.GetInstance().Where("password = ?", encodePwd).First(&user)
	if user.ID == 0 {
		c.JSON(http.StatusUnauthorized, entities.TokenResponse{Token: ""})
		return
	}

	jwtCode, err := jwt_tools.NewJWT(settings.SecretKey, 24*time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entities.TokenResponse{Token: ""})
		return
	}

	db.GetInstance().Model(&user).Update("token", jwtCode)

	c.JSON(http.StatusOK, entities.TokenResponse{Token: jwtCode})
}

func DeleteToken(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if strings.Contains(token, "Bearer ") {
		token = token[len("Bearer "):]
	}

	var user = db.User{}
	db.GetInstance().Where("token = ?", token).First(&user)
	if user.ID == 0 {
		c.JSON(http.StatusUnauthorized, entities.TokenResponse{Token: ""})
		return
	}

	db.GetInstance().Model(&user).Update("token", "")
	c.JSON(http.StatusOK, entities.TokenResponse{Token: token})
}

func Register(c *gin.Context) {

}
