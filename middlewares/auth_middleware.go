package middlewares

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// ValidateToken
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. ขอดูบัตรจาก Header ที่ชื่อ "Authorization"
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "กรุณาเข้าสู่ระบบ (ไม่พบ Token)"})
			return
		}

		// 2. รูปแบบมาตรฐานของ Token จะส่งมาเป็น "Bearer eyJhbGciOi..." เราต้องตัดคำว่า Bearer ออก
		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		// 3. เอาบัตรไปเข้าเครื่องตรวจสอบ (เช็คลายเซ็นและวันหมดอายุ)
		secretKey := os.Getenv("JWT_SECRET")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// ยืนยันว่าใช้อัลกอริทึมถูกต้อง
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token ไม่ถูกต้องหรือหมดอายุแล้ว"})
			return
		}

		// 4. JWT ถูกต้อง ใส่ user_id และ role ลงใน Context เผื่อให้ Controller เอาไปใช้ต่อ
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("user_id", claims["user_id"])
			if role, roleExists := claims["role"].(string); roleExists {
				c.Set("role", role)
			} else {
				// เผื่อกรณี Token เก่าไม่มี Role ฝังมา ให้มองเป็น user ธรรมดาไปก่อน
				c.Set("role", "user")
			}
		}

		// 5. ปล่อยให้เดินผ่านเข้าประตูไปได้!
		c.Next()
	}
}

// c คือ Context ของ Gin
// c.GetHeader("Authorization") คือการดึงข้อมูลจาก Header ที่ชื่อ "Authorization"
// c.AbortWithStatusJSON คือการบอกว่าให้หยุดการทำงานของ Route นี้ และส่ง Status Code กับ JSON กลับไปทันที
// jwt.Parse คือฟังก์ชันจากไลบรารี JWT ที่ใช้ในการตรวจสอบความถูกต้องของ Token

// สิ่งที่context เก็บมีดังนี้
// c.Set("user_id", claims["user_id"]) คือการเก็บข้อมูล user_id ไว้ใน Context เพื่อให้ Controller สามารถดึงไปใช้ได้ต่อไป
