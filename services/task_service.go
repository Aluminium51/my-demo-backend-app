package services

import (
	"my-go-backend/models"
	"my-go-backend/repositories"
)

type TaskService interface {
	CreateTask(title string, description string, userID uint) (*models.Task, error)
	GetTasksByUserID(userID uint) ([]models.Task, error)
}

type taskService struct {
	repo repositories.TaskRepository // พ่อครัวต้องรู้จักพนักงานหลังร้าน
}

func NewTaskService(repo repositories.TaskRepository) TaskService {
	return &taskService{repo: repo}
}

// โค้ดส่วนนี้จัดการ Business Logic ล้วนๆ
func (s *taskService) CreateTask(title string, description string, userID uint) (*models.Task, error) {
	// ตัวอย่าง Business Logic: ถ้าไม่ได้ตั้งชื่องาน ให้ใส่ค่าเริ่มต้น
	if title == "" {
		title = "Untitled Task"
	}

	task := &models.Task{
		Title:       title,
		Description: description,
		UserID:      userID,
	}

	// สั่งพนักงานหลังร้านให้เอาไปเซฟลง DB
	err := s.repo.Create(task)
	return task, err
}

func (s *taskService) GetTasksByUserID(userID uint) ([]models.Task, error) {
	return s.repo.GetByUserID(userID)
}
