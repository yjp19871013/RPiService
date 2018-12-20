package api

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/yjp19871013/RPiService/users/dto"

	"github.com/yjp19871013/RPiService/users/jwt_tools"

	"github.com/yjp19871013/RPiService/users/settings"

	"github.com/yjp19871013/RPiService/utils"

	"github.com/yjp19871013/RPiService/users/db"

	"github.com/gin-gonic/gin"
)

func CreateToken(c *gin.Context) {
	var request dto.CreateTokenRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, dto.TokenResponse{Token: ""})
		return
	}

	encodePwd := utils.MD5(request.Password)

	var user = db.User{}
	err = db.GetInstance().Where("password = ?", encodePwd).First(&user).Error
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	jwtCode, err := jwt_tools.NewJWT(settings.SecretKey, 24*time.Hour)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	err = db.GetInstance().Model(&user).Update("token", jwtCode).Error
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, dto.TokenResponse{Token: jwtCode})
}

func DeleteToken(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if strings.Contains(token, "Bearer ") {
		token = token[len("Bearer "):]
	}

	var user = db.User{}
	err := db.GetInstance().Where("token = ?", token).First(&user).Error
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.TokenResponse{Token: ""})
		return
	}

	err = db.GetInstance().Model(&user).Update("token", "").Error
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, dto.TokenResponse{Token: token})
}

func GenerateValidateCode(c *gin.Context) {
	var request dto.ValidateCodeRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var user = db.User{}
	err = db.GetInstance().Where("email = ?", request.Email).First(&user).Error
	if err == nil {
		c.AbortWithStatus(http.StatusConflict)
		return
	}

	var validateCode = db.ValidateCode{}
	err = db.GetInstance().Where("email = ?", request.Email).First(&validateCode).Error
	if err != nil {
		validateCode.Email = request.Email
	}

	validateCode.ValidateCode = utils.GenerateValidateCode()
	err = db.GetInstance().Save(&validateCode).Error
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	//err = utils.SendEmail("RPiService验证码", "Your Code is: "+validateCode.ValidateCode, request.Email)
	//if err != nil {
	//	c.AbortWithStatus(http.StatusInternalServerError)
	//	return
	//}

	log.Println(validateCode.ValidateCode)

	c.AbortWithStatus(http.StatusOK)
}

func Register(c *gin.Context) {
	var request dto.RegisterRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var validateCode = db.ValidateCode{}
	err = db.GetInstance().Where("email = ?", request.Email).First(&validateCode).Error
	if err != nil || validateCode.ValidateCode != request.ValidateCode {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if request.Password1 != request.Password2 {
		c.AbortWithStatus(http.StatusNotAcceptable)
		return
	}

	var user db.User
	err = db.GetInstance().Where("email = ?", request.Email).First(&user).Error
	if err == nil {
		c.AbortWithStatus(http.StatusConflict)
		return
	}

	var commonRole db.Role
	err = db.GetInstance().Where("name = ?", db.CommonRoleName).First(&commonRole).Error
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var newUser = db.User{}
	newUser.Email = request.Email
	newUser.Password = utils.MD5(request.Password1)
	newUser.Roles = []db.Role{
		commonRole,
	}

	err = db.GetInstance().Save(&newUser).Error
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	db.GetInstance().Delete(&validateCode)

	c.AbortWithStatus(http.StatusOK)
}
