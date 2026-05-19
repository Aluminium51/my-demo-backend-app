package config

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client
var Ctx = context.Background()

func ConnectRedis() {

	// Read Host from env
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost" // เผื่อกรณีรันแบบไม่ใช้ Docker
	}

	// Read Port from env
	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		redisPort = "6379"
	}

	redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort)

	RDB = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})

	_, err := RDB.Ping(Ctx).Result()
	if err != nil {
		fmt.Println("❌ ไม่สามารถเชื่อมต่อ Redis ได้:", err)
		return
	}
	fmt.Println("✅Redis Connected Successfully at : ", redisAddr)
}

// go get github.com/redis/go-redis/v9
