package db

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/yjp19871013/RPiService/utils"
)

const (
	SuperPermissionName = "super"
	SuperPermissionDesc = "super permissions"

	CommonPermissionName = "common"
	CommonPermissionDesc = "common permissions"

	AdminRoleName = "admin_role"
	AdminRoleDesc = "admin role"

	CommonRoleName = "common_role"
	CommonRoleDesc = "common role"

	AdminUserEmail    = "admin@yjp.com"
	AdminUserPassword = "123456"
)

var (
	db *gorm.DB
)

func GetInstance() *gorm.DB {
	return db
}

func InitDb() {
	var err error
	db, err = gorm.Open("mysql", "root:root@/rpi_users?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err.Error())
	}

	db.AutoMigrate(&Permission{}, &Role{}, &User{}, &ValidateCode{})

	db.LogMode(true)

	createPermissions()
	createRoles()
	createUsers()
}

func CloseDb() {
	db.Close()
}

func createPermissions() {
	var superPermission Permission
	err := db.Where("name = ?", SuperPermissionName).First(&superPermission).Error
	if err != nil {
		superPermission.Name = SuperPermissionName
		superPermission.Description = SuperPermissionDesc
		db.Save(&superPermission)
	}

	var commonPermission Permission
	err = db.Where("name = ?", CommonPermissionName).First(&commonPermission).Error
	if err != nil {
		log.Println(err, commonPermission)
		commonPermission.Name = CommonPermissionName
		commonPermission.Description = CommonPermissionDesc
		db.Save(&commonPermission)
	}
}

func createRoles() {
	var superPermission Permission
	err := db.Where("name = ?", SuperPermissionName).First(&superPermission).Error
	if err != nil {
		return
	}

	var adminRole Role
	err = db.Where("name = ?", AdminRoleName).First(&adminRole).Error
	if err != nil {
		adminRole.Name = AdminRoleName
		adminRole.Description = AdminRoleDesc
		adminRole.Permissions = []Permission{superPermission}
		db.Save(&adminRole)
	}

	var commonPermission Permission
	err = db.Where("name = ?", CommonPermissionName).First(&commonPermission).Error
	if err != nil {
		return
	}

	var commonRole Role
	err = db.Where("name = ?", CommonRoleName).First(&commonRole).Error
	if err != nil {
		commonRole.Name = CommonRoleName
		commonRole.Description = CommonRoleDesc
		commonRole.Permissions = []Permission{commonPermission}
		db.Save(&commonRole)
	}
}

func createUsers() {
	var adminRole Role
	err := db.Where("name = ?", AdminRoleName).First(&adminRole).Error
	if err != nil {
		return
	}

	var adminUser User
	err = db.Where("email = ?", AdminUserEmail).First(&adminUser).Error
	if err != nil {
		adminUser.Email = AdminUserEmail
		adminUser.Password = utils.MD5(AdminUserPassword)
		adminUser.Roles = []Role{adminRole}
		db.Save(&adminUser)
	}
}
