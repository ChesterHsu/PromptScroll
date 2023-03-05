package config

import (
	"fmt"
	"log"
	"promptscroll/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ConnectDatabase 用來建立與資料庫的連線
func ConnectDatabase() (*gorm.DB, error) {
	// 取得資料庫連線設定
	dbConfig := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable TimeZone=Asia/Taipei",
		GetString("db_host"),
		GetString("db_port"),
		GetString("db_user"),
		GetString("db_name"),
		GetString("db_pass"))

	// 建立與資料庫的連線
	db, err := gorm.Open(postgres.Open(dbConfig), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
		return nil, err
	}

	db.AutoMigrate(&model.User{}, &model.Auth{})

	return db, nil
}
