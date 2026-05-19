package main

import (
	"log"
	"os"

	"my-go-backend/config"
	"my-go-backend/routes"

	"github.com/joho/godotenv"
)

func main() {
	// 1. โหลดตัวแปรจากไฟล์ .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using OS environments")
	}

	// 2. เชื่อมต่อ Database + เรียกใช้ Redis
	config.ConnectDB()
	config.ConnectRedis()

	// 3. ตั้งค่า Route
	r := routes.SetupRouter()

	// 4. สั่งรัน Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)
}
