package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	db *gorm.DB
)

func GetInstance() *gorm.DB {
	return db
}

func InitDb() {
	var err error
	db, err = gorm.Open("mysql", "root:root@/rpi_file_station?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err.Error())
	}

	db.AutoMigrate(&DownloadTask{})

	db.LogMode(true)
}

func CloseDb() {
	db.Close()
}
