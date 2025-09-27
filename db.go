package main

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dsn := os.Getenv("DATABASE_DSN") // пример: postgres://user:pass@localhost:5432/mydb?sslmode=disable
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	DB = db

	err = DB.AutoMigrate(&User{}, &Post{})
	if err != nil {
		panic("failed to migrate: " + err.Error())
	}

	fmt.Println("✅ Connected to database and migrated successfully")
}
