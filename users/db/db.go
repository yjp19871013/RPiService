package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/yjp19871013/RPiService/users/model"
)

var (
	db *gorm.DB

	allPermissions = []*model.Permission{
		{
			Name:        "all",
			Description: "All permissions",
		},
	}

	allRoles = []*model.Role{
		{
			Name:        "admin_role",
			Description: "admin role",
		},
	}

	adminUser = &model.User{
		Username: "admin",
		Password: "3400CD4574D4D14D29251E5EA620A925", // rpi_admin
	}
)

func getInstance() *gorm.DB {
	return db
}

func InitDb() {
	var err error
	db, err = gorm.Open("mysql", "root:root@/rpi_users?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err.Error())
	}

	db.AutoMigrate(&model.Permission{}, &model.Role{}, &model.User{})

	createAdminUser()
}

func CloseDb() {
	db.Close()
}

func createAdminUser() {
	permissions := make([]*model.Permission, 0)
	db.Find(&permissions)
	if len(permissions) != len(allPermissions) {
		for _, permission := range allPermissions {
			db.Create(permission)
		}
	}

	roles := make([]*model.Role, 0)
	db.Find(&roles)
	if len(roles) != len(allRoles) {
		for _, role := range allRoles {
			db.Create(role)
		}
	}

	var user model.User
	db.Where("username = ?", adminUser.Username).First(&user)
	if user.Username == "" {
		db.Create(adminUser)
	}
}
