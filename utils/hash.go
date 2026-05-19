package utils

import "golang.org/x/crypto/bcrypt"

// Convert normal password to hashed password
func HashPassword(password string) (string, error) {
	// ค่า 14 คือ Cost (ความยากในการแฮช) ยิ่งเลขเยอะยิ่งปลอดภัยแต่ก็จะใช้เวลาคำนวณนานขึ้น (มาตรฐานนิยมใช้ 10-14)
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash ใช้เปรียบเทียบรหัสที่รับมา กับรหัสใน Database ว่าตรงกันไหม (ใช้ตอนล็อกอิน)
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil // ถ้าไม่มี Error แปลว่ารหัสผ่านถูกต้อง
}
