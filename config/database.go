package config

import (
	"fmt"
	"log"
	"os"
	"time"

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

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get database generic interface:", err)
	}
	// ตั้งค่า Connection Pool
	sqlDB.SetMaxIdleConns(10)           // จำนวน Connection ที่จะเก็บไว้ใน Pool สำหรับการใช้งานซ้ำ (ไม่ต้องสร้างใหม่ทุกครั้ง)
	sqlDB.SetMaxOpenConns(50)           // จำนวน Connection สูงสุดที่สามารถเปิดได้พร้อมกัน (ถ้าเกินจะรอจนกว่าจะมี Connection ว่าง)
	sqlDB.SetConnMaxLifetime(time.Hour) // ระยะเวลาที่ Connection จะถูกใช้งานก่อนที่จะถูกปิดและสร้างใหม่ (ช่วยป้องกัน Connection ที่ค้างอยู่เกินไป)

	DB = db
	DB.AutoMigrate(&models.User{}, &models.Task{})
	fmt.Println("✅Database Connected Successfully at : ", dsn)
}
