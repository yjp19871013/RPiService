package db

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model

	Email    string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Token    string `gorm:"not null"`
	Roles    []Role `gorm:"many2many:user_roles"`

	DownloadTask []DownloadTask
	FileInfos    []FileInfo
}

func FindUserByEmail(email string) (*User, error) {
	var user = &User{}
	err := db.Where("email = ?", email).First(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func FindUserByToken(token string) (*User, error) {
	var user = &User{}
	err := db.Where("token = ?", token).First(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func GetUserRoles(user *User) ([]Role, error) {
	roles := make([]Role, 0)
	err := db.Model(user).Related(&roles, "Roles").Error
	if err != nil {
		return nil, err
	}

	return roles, nil
}

func GetAllUsers() ([]User, error) {
	users := make([]User, 0)
	err := db.Preload("Roles").Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}

func SaveUser(user *User) error {
	return db.Save(user).Error
}
