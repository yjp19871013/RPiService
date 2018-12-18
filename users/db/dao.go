package db

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model

	Email    string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Token    string `gorm:"unique;not null"`
	Roles    []Role `gorm:"many2many:user_roles"`
}

type Role struct {
	ID          uint   `gorm:"primary_key"`
	Name        string `gorm:"unique;not null"`
	Description string
	Permissions []Permission `gorm:"many2many:role_permissions"`
}

type Permission struct {
	ID          uint   `gorm:"primary_key"`
	Name        string `gorm:"unique;not null"`
	Description string
}
