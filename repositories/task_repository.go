package repositories

import (
	"my-go-backend/config"
	"my-go-backend/models"
)

// สร้าง Interface เพื่อกำหนดว่าโกดังนี้ทำอะไรได้บ้าง (ดีมากสำหรับการทำ Unit Test)
type TaskRepository interface {
	Create(task *models.Task) error
	GetByUserID(userID uint) ([]models.Task, error)
}

type taskRepository struct{}

func NewTaskRepository() TaskRepository {
	return &taskRepository{}
}

// โค้ดส่วนนี้มีหน้าที่คุยกับ GORM อย่างเดียว
func (r *taskRepository) Create(task *models.Task) error {
	return config.DB.Create(task).Error
}

func (r *taskRepository) GetByUserID(userID uint) ([]models.Task, error) {
	var tasks []models.Task
	// ดึงเฉพาะงานที่ user_id ตรงกับที่ส่งเข้ามา
	err := config.DB.Where("user_id = ?", userID).Find(&tasks).Error
	return tasks, err
}
