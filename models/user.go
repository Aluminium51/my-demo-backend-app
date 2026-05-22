package models

import "gorm.io/gorm"

// ใส่ gorm.Model เพื่อให้มีฟิลด์ ID, CreatedAt, UpdatedAt, DeletedAt อัตโนมัติ
type User struct {
	gorm.Model
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email" gorm:"unique"`
	Age   int    `json:"age"`
	// 1. เอา binding:"required" ออก เพราะคนล็อกอินผ่าน Google จะไม่มีรหัสผ่าน
	// 2. ใช้ *string (Pointer) เพื่อให้ GORM สามารถเก็บค่า NULL ลงใน Database ได้
	Password *string `json:"password,omitempty"`

	// 3. เพิ่ม GoogleID เพื่อใช้เชื่อมโยงกับบัญชี Google และป้องกันการสมัครซ้ำ, ใช้ *string เพื่อให้เก็บค่า NULL ได้ (สำหรับผู้ที่ไม่ได้ล็อกอินผ่าน Google)
	GoogleID *string `json:"google_id" gorm:"uniqueIndex"`

	// one-to-many relationship: หนึ่ง User มีหลาย Task
	Tasks []Task `json:"tasks" gorm:"foreignKey:UserID"`
	Role  string `json:"role" gorm:"default:'user'"`
}
