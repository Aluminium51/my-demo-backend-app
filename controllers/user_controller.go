package controllers

import (
	"net/http"
	"strconv"

	"my-go-backend/models"
	"my-go-backend/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// จะถูกเรียกใช้ใน main.go เพื่อสร้าง Route และรับ Request สำหรับ User
type UserController struct {
	service services.UserService // เรียกใช้ service ของ User
}

func NewUserController(service services.UserService) *UserController {
	return &UserController{service: service}
}

type RegisterInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Age      int    `json:"age" binding:"required,gte=0"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (ctrl *UserController) CreateUser(c *gin.Context) {
	var input RegisterInput

	// Data Validation ด้วย Gin + go-playground/validator
	// 1. ตรวจสอบข้อมูลจาก JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		// 2. ดักจับว่า Error เกิดจาก Validation หรือเปล่า?
		if errs, ok := err.(validator.ValidationErrors); ok {
			var errorMessages []string
			for _, e := range errs {
				switch e.Tag() {
				case "required":
					errorMessages = append(errorMessages, "กรุณากรอกข้อมูลในช่อง "+e.Field())
				case "email":
					errorMessages = append(errorMessages, "รูปแบบอีเมลไม่ถูกต้อง")
				case "min":
					errorMessages = append(errorMessages, e.Field()+" ต้องมีความยาวอย่างน้อย "+e.Param()+" ตัวอักษร")
				case "gte":
					errorMessages = append(errorMessages, e.Field()+" ต้องมีค่ามากกว่าหรือเท่ากับ "+e.Param())
				default:
					errorMessages = append(errorMessages, "ข้อมูล "+e.Field()+" ไม่ถูกต้อง")
				}
			}
			// ตอบกลับเป็น Array ของข้อความแจ้งเตือนให้ Frontend เอาไปโชว์ง่ายๆ
			c.JSON(http.StatusBadRequest, gin.H{"errors": errorMessages})
			return
		}

		// ถ้าเป็น Error ปกติ (เช่น ส่ง JSON ผิดโครงสร้าง)
		c.JSON(http.StatusBadRequest, gin.H{"error": "รูปแบบข้อมูลไม่ถูกต้อง"})
		return
	}

	// เรียกใช้ UserService จาก package services
	user, err := ctrl.service.Register(input.Name, input.Email, input.Password, input.Age)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "สร้างบัญชีไม่สำเร็จ: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "สร้างบัญชีผู้ใช้สำเร็จ!", "data": user})
}

func (ctrl *UserController) Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ข้อมูลไม่ถูกต้อง"})
		return
	}

	token, err := ctrl.service.Login(input.Email, input.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "เข้าสู่ระบบสำเร็จ", "token": token})
}

func (ctrl *UserController) GetUsers(c *gin.Context) {
	users, err := ctrl.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": users})
}

func (ctrl *UserController) UpdateUser(c *gin.Context) {
	paramID, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	userIDContext, _ := c.Get("user_id")
	loggedInUserID := uint(userIDContext.(float64))

	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := ctrl.service.Update(uint(paramID), loggedInUserID, input)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func (ctrl *UserController) DeleteUser(c *gin.Context) {
	paramID, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	userIDContext, _ := c.Get("user_id")
	loggedInUserID := uint(userIDContext.(float64))

	err := ctrl.service.Delete(uint(paramID), loggedInUserID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": true})
}
