package db

import (
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
	superPermission := Permission{
		Name:        SuperPermissionName,
		Description: SuperPermissionDesc,
	}
	db.Save(&superPermission)

	commonPermission := Permission{
		Name:        CommonPermissionName,
		Description: CommonPermissionDesc,
	}
	db.Save(&commonPermission)
}

func createRoles() {
	var superPermission Permission
	db.Where("name = ?", SuperPermissionName).First(&superPermission)

	adminRole := Role{
		Name:        AdminRoleName,
		Description: AdminRoleDesc,
		Permissions: []Permission{
			superPermission,
		},
	}
	db.Save(&adminRole)
}

func createUsers() {
	var adminRole Role
	db.Where("name = ?", AdminRoleName).First(&adminRole)

	adminUser := User{
		Email:    AdminUserEmail,
		Password: utils.MD5(AdminUserPassword),
		Roles: []Role{
			adminRole,
		},
	}
	db.Save(&adminUser)
}
