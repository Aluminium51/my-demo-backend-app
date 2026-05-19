package routes

import (
	"my-go-backend/controllers"
	"my-go-backend/middlewares" // อย่าลืม Import middleware เข้ามา
	"my-go-backend/repositories"
	"my-go-backend/services"

	swaggerDocs "my-go-backend/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Health Check Endpoint
	r.GET("/health", controllers.HealthCheck)

	// === ประกอบร่างวัตถุดิบ (Dependency Injection) ===
	// 1. สร้างโกดัง (มีแค่อันเดียว แต่ให้ Service หลายตัวยืมใช้ได้)
	userRepo := repositories.NewUserRepository()

	// 2. ประกอบร่างฝั่ง Auth (เข้าสู่ระบบ)
	authService := services.NewAuthService(userRepo)
	authController := controllers.NewAuthController(authService)

	// 3. ประกอบร่างฝั่ง User (จัดการข้อมูลโปรไฟล์)
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	// 4. ฝั่ง Task
	taskRepo := repositories.NewTaskRepository()
	taskService := services.NewTaskService(taskRepo)
	taskController := controllers.NewTaskController(taskService)

	api := r.Group("/api/v1")
	{
		// === โซนสาธารณะ (ใครๆ ก็เข้าได้) ===
		api.POST("/login", userController.Login)
		api.POST("/users", userController.CreateUser) // สมัครสมาชิก

		api.GET("/auth/google/login", authController.GoogleLogin)
		api.GET("/auth/google/callback", authController.GoogleCallback)

		// === โซน VIP (ต้องมี Token เท่านั้น) ===
		// การใช้ .Use() เป็นการบอกว่า ทุก Route ที่อยู่ใต้กลุ่มนี้ ต้องผ่าน AuthRequired() ก่อนเสมอ
		protected := api.Group("/")
		protected.Use(middlewares.AuthRequired())
		{
			// User Management
			protected.GET("/users", userController.GetUsers)
			protected.PUT("/users/:id", userController.UpdateUser)
			protected.DELETE("/users/:id", userController.DeleteUser)

			// Task Management
			protected.POST("/tasks", taskController.CreateTask)
			protected.GET("/tasks", taskController.GetTasks)
		}
	}

	// swagger
	swaggerDocs.SwaggerInfo.Title = "My Go Backend API"
	swaggerDocs.SwaggerInfo.Description = "นี่คือคู่มือ API สำหรับโปรเจกต์จัดการ Task"
	swaggerDocs.SwaggerInfo.Version = "1.0"
	swaggerDocs.SwaggerInfo.Host = "localhost:8081"
	swaggerDocs.SwaggerInfo.BasePath = "/"

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
