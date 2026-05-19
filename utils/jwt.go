package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateToken จะรับ ID ของ User มาเพื่อฝังไว้ในบัตร
func GenerateToken(userID uint) (string, error) {
	// ดึงรหัสลับจากไฟล์ .env
	secretKey := os.Getenv("JWT_SECRET")

	// สร้าง Payload (ข้อมูลที่จะฝังไว้ในบัตร)
	claims := jwt.MapClaims{
		"user_id": userID,                                // ฝัง ID ลงไป เพื่อให้รู้ว่าบัตรนี้เป็นของใคร
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // บัตรมีอายุ 24 ชั่วโมง
	}

	// สร้างตัว Token โดยใช้อัลกอริทึม HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// ประทับตราด้วย Secret Key แล้วแปลงเป็น String ยาวๆ ส่งกลับไป
	return token.SignedString([]byte(secretKey))
}
