package routes

import (
	"my-go-backend/controllers"
	"my-go-backend/middlewares" // อย่าลืม Import middleware เข้ามา
	"my-go-backend/repositories"
	"my-go-backend/services"

	swaggerDocs "my-go-backend/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// CORS Middleware (อนุญาตให้ Frontend ที่อยู่คนละโดเมนเข้าถึง API ได้)
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // ในอนาคตถ้ามีเว็บหน้าบ้านจริง ค่อยเปลี่ยนจาก "*" เป็น "https://your-frontend.com"
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

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
		// 🟢 โซนทั่วไป: เข้าได้ทุกคน
		api.POST("/register", userController.CreateUser) // สมัครสมาชิก
		api.POST("/login", userController.Login)

		api.GET("/auth/google/login", authController.GoogleLogin)
		api.GET("/auth/google/callback", authController.GoogleCallback)

		// // 🟡 โซนสมาชิก: ต้องมีToken
		protected := api.Group("/")
		protected.Use(middlewares.AuthRequired())
		{
			// User Management
			protected.PUT("/users/:id", userController.UpdateUser)

			// Task Management
			protected.POST("/tasks", taskController.CreateTask)
			protected.GET("/tasks", taskController.GetTasks)

			// 🔴 โซน Admin Only: ต้องล็อกอิน และต้องมียศ admin
			adminOnly := protected.Group("/")
			adminOnly.Use(middlewares.RequireRole("admin"))
			{
				adminOnly.GET("/users", userController.GetUsers)          // ดึงรายชื่อทั้งหมด
				adminOnly.DELETE("/users/:id", userController.DeleteUser) // ลบบัญชีชาวบ้าน
			}
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
