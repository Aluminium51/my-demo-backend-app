package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateToken จะรับ ID ของ User มาเพื่อฝังไว้ในบัตร
func GenerateToken(userID uint, role string) (string, error) {

	// สร้าง Payload (ข้อมูลที่จะฝังไว้ในบัตร)
	claims := jwt.MapClaims{
		"user_id": userID, // ฝัง ID ลงไป เพื่อให้รู้ว่าบัตรนี้เป็นของใคร
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Has a 24-hour expiration time
	}

	// สร้างตัว Token โดยใช้อัลกอริทึม HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey := os.Getenv("JWT_SECRET")
	return token.SignedString([]byte(secretKey))
}
