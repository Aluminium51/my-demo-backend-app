package controllers

import (
	"my-go-backend/models"
	"my-go-backend/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// จะถูกเรียกใช้ใน main.go เพื่อสร้าง Route และรับ Request
type TaskController struct {
	service services.TaskService // เรียกใช้ service ของTask
}

func NewTaskController(service services.TaskService) *TaskController {
	return &TaskController{service: service}
}

func (ctrl *TaskController) CreateTask(c *gin.Context) {
	// 1. รับ Request
	userIDContext, _ := c.Get("user_id")
	loggedInUserID := uint(userIDContext.(float64))

	var input models.Task
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. ส่งออเดอร์ให้พ่อครัว (Service) จัดการ
	task, err := ctrl.service.CreateTask(input.Title, input.Description, loggedInUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถสร้างงานได้"})
		return
	}

	// 3. เสิร์ฟอาหารกลับไปให้ลูกค้า
	c.JSON(http.StatusCreated, gin.H{"message": "สร้างงานสำเร็จ", "data": task})
}

func (ctrl *TaskController) GetTasks(c *gin.Context) {
	// 1. ตรวจบัตรและดึง ID ของลูกค้า
	userIDContext, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ไม่พบข้อมูลผู้ใช้งานในระบบ"})
		return
	}
	loggedInUserID := uint(userIDContext.(float64))

	// 2. ส่งออเดอร์ให้พ่อครัว (Service) ไปจัดการ
	tasks, err := ctrl.service.GetTasksByUserID(loggedInUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถดึงข้อมูลงานได้"})
		return
	}

	// 3. เสิร์ฟอาหารกลับไปให้ลูกค้า
	c.JSON(http.StatusOK, gin.H{
		"message": "ดึงข้อมูลงานสำเร็จ",
		"data":    tasks,
	})
}
