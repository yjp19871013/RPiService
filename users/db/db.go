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
	_ = db.Close()
}

func createPermissions() {
	superPermission, err := FindPermissionByName(SuperPermissionName)
	if err != nil {
		superPermission = &Permission{
			Name:        SuperPermissionName,
			Description: SuperPermissionDesc,
		}

		err = SavePermission(superPermission)
		if err != nil {
			panic("Create super permission error")
		}
	}

	commonPermission, err := FindPermissionByName(CommonPermissionName)
	if err != nil {
		commonPermission = &Permission{
			Name:        CommonPermissionName,
			Description: CommonPermissionDesc,
		}

		err = SavePermission(commonPermission)
		if err != nil {
			panic("Create common permission error")
		}
	}
}

func createRoles() {
	superPermission, err := FindPermissionByName(SuperPermissionName)
	if err != nil {
		return
	}

	adminRole, err := FindRoleByName(AdminRoleName)
	if err != nil {
		adminRole = &Role{
			Name:        AdminRoleName,
			Description: AdminRoleDesc,
			Permissions: []Permission{*superPermission},
		}

		err = SaveRole(adminRole)
		if err != nil {
			panic("Create admin role failed")
		}
	}

	commonPermission, err := FindPermissionByName(CommonPermissionName)
	if err != nil {
		return
	}

	commonRole, err := FindRoleByName(CommonRoleName)
	if err != nil {
		commonRole = &Role{
			Name:        CommonRoleName,
			Description: CommonRoleDesc,
			Permissions: []Permission{*commonPermission},
		}

		err = SaveRole(commonRole)
		if err != nil {
			panic("Create common role failed")
		}
	}
}

func createUsers() {
	adminRole, err := FindRoleByName(AdminRoleName)
	if err != nil {
		return
	}

	adminUser, err := FindUserByEmail(AdminUserEmail)
	if err != nil {
		adminUser = &User{
			Email:    AdminUserEmail,
			Password: utils.MD5(AdminUserPassword),
			Roles:    []Role{*adminRole},
		}

		err := SaveUser(adminUser)
		if err != nil {
			panic("Create admin user failed")
		}
	}
}
