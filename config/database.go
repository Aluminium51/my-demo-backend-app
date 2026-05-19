package config

import (
	"fmt"
	"log"
	"os"

	"my-go-backend/models" // เปลี่ยนชื่อ my-go-backend ตามชื่อ module ของคุณ

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	// ดึงค่าจาก .env
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Bangkok",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto Migrate: สั่งให้ GORM สร้าง/อัปเดตตารางให้ตรงกับ Struct อัตโนมัติ
	err = db.AutoMigrate(&models.User{}, &models.Task{})
	if err != nil {
		log.Fatal("Failed to auto-migrate database:", err)
		return
	}

	DB = db
	fmt.Println("✅Database Connected Successfully at : ", dsn)
}
