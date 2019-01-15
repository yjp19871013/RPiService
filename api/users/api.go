package users

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/yjp19871013/RPiService/api/users/dto"
	"github.com/yjp19871013/RPiService/settings"

	"github.com/yjp19871013/RPiService/db"

	"github.com/yjp19871013/RPiService/jwt_tools"

	"github.com/yjp19871013/RPiService/utils"

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

	user, err := db.FindUserByEmail(request.Email)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if user.Password != utils.MD5(request.Password) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	jwtCode, err := jwt_tools.NewJWT(settings.SecretKey, 24*time.Hour)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	user.Token = jwtCode
	err = db.SaveUser(user)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	roles, err := db.GetUserRoles(user)
	if err != nil {
		log.Println(err)
		roles = make([]db.Role, 0)
	}

	response := dto.TokenResponse{
		Token: jwtCode,
		Roles: make([]string, 0),
	}

	for _, role := range roles {
		response.Roles = append(response.Roles, role.Name)
	}

	c.JSON(http.StatusOK, response)
}

func GenerateValidateCode(c *gin.Context) {
	var request dto.ValidateCodeRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	_, err = db.FindUserByEmail(request.Email)
	if err == nil {
		c.AbortWithStatus(http.StatusConflict)
		return
	}

	validateCode, err := db.FindValidateCodeByEmail(request.Email)
	if err != nil {
		validateCode = &db.ValidateCode{
			Email: request.Email,
		}
	}

	validateCode.ValidateCode = utils.GenerateValidateCode()
	err = db.SaveValidateCode(validateCode)
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

	validateCode, err := db.FindValidateCodeByEmail(request.Email)
	if err != nil || validateCode.ValidateCode != request.ValidateCode {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if request.Password1 != request.Password2 {
		c.AbortWithStatus(http.StatusNotAcceptable)
		return
	}

	_, err = db.FindUserByEmail(request.Email)
	if err == nil {
		c.AbortWithStatus(http.StatusConflict)
		return
	}

	commonRole, err := db.FindRoleByName(db.CommonRoleName)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var newUser = &db.User{}
	newUser.Email = request.Email
	newUser.Password = utils.MD5(request.Password1)
	newUser.Roles = []db.Role{
		*commonRole,
	}

	err = db.SaveUser(newUser)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	err = db.DeleteValidateCode(validateCode)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	_ = os.Mkdir(newUser.Email, os.ModeDir|os.ModePerm)

	c.AbortWithStatus(http.StatusOK)
}

func GetAllRoles(c *gin.Context) {
	roles, err := db.GetAllRoles()
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	response := dto.GetRolesResponse{
		Roles: make([]string, 0),
	}

	for _, role := range roles {
		response.Roles = append(response.Roles, role.Name)
	}

	c.JSON(http.StatusOK, response)
}

func GetAllUsers(c *gin.Context) {
	users, err := db.GetAllUsers()
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	response := dto.GetUsersResponse{
		UserInfos: make([]*dto.UserInfo, 0),
	}

	for _, user := range users {
		userInfo := &dto.UserInfo{}
		userInfo.Email = user.Email

		for _, role := range user.Roles {
			userInfo.Roles = append(userInfo.Roles, role.Name)
		}

		userInfo.CreateDate = user.Model.CreatedAt.Format("2006-1-2 15:04:05")
		userInfo.UpdateDate = user.Model.UpdatedAt.Format("2006-1-2 15:04:05")

		response.UserInfos = append(response.UserInfos, userInfo)
	}

	c.JSON(http.StatusOK, response)
}
