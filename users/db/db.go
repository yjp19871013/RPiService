package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	db *gorm.DB

	allPermissions = []*Permission{
		{
			Name:        "all",
			Description: "All permissions",
		},
	}

	allRoles = []*Role{
		{
			Name:        "admin_role",
			Description: "admin role",
		},
	}

	adminUser = &User{
		Email:    "admin@yjp.com",
		Password: "3400CD4574D4D14D29251E5EA620A925", // rpi_admin
	}
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

	db.AutoMigrate(&Permission{}, &Role{}, &User{})

	db.LogMode(true)

	createAdminUser()
}

func CloseDb() {
	db.Close()
}

func createAdminUser() {
	permissions := make([]*Permission, 0)
	db.Find(&permissions)
	if len(permissions) != len(allPermissions) {
		for _, permission := range allPermissions {
			db.Create(permission)
		}
	}

	roles := make([]*Role, 0)
	db.Find(&roles)
	if len(roles) != len(allRoles) {
		for _, role := range allRoles {
			db.Create(role)
		}
	}

	var user User
	db.Where("username = ?", adminUser.Email).First(&user)
	if user.ID == 0 {
		db.Create(adminUser)
	}
}
