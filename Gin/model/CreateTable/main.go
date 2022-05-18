package main

import (
	"Douyin/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func main() {
	addr := "root:root@tcp(localhost:3306)/douyin_user?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(addr), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		return
	}
	_ = db.AutoMigrate(&model.User{}, &model.Video{}, &model.FavoriteVideo{}, &model.Relation{}, &model.Comment{})
}
