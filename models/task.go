package models

import "gorm.io/gorm"

type Task struct {
	gorm.Model
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`

	// นี่คือ Foreign Key ที่ GORM จะรู้ทันทีว่ามันผูกกับตาราง User
	UserID uint `json:"user_id"`
}
