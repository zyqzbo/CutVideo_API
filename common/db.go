package common

import (
	"CutVido_api/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:zyq4836..@tcp(127.0.0.1:3306)/cut_video_db?charset=utf8&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		panic("err:" + err.Error())
	}
	db.AutoMigrate(&models.Video{})

	DB = db
	return DB
}

func GetDB() *gorm.DB {
	return DB
}
