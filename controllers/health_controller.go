package controllers

import (
	"my-go-backend/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthCheck godoc
// @Summary ตรวจสอบสถานะระบบ (Health Check)
// @Description ใช้สำหรับให้ระบบ Monitor เช็คว่า API และ Database ยังทำงานปกติหรือไม่
// @Tags System
// @Produce json
// @Success 200 {object} map[string]interface{} "ระบบปกติ"
// @Failure 503 {object} map[string]interface{} "ระบบมีปัญหา"
// @Router /health [get]
func HealthCheck(c *gin.Context) {
	// 1. เช็คสถานะของ API
	apiStatus := "UP"

	// 2. เช็คสถานะการเชื่อมต่อ Database
	dbStatus := "UP"

	// ดึง object ฐานข้อมูลระดับล่าง (sql.DB) ออกมาจาก GORM เพื่อใช้คำสั่ง Ping
	sqlDB, err := config.DB.DB()

	// ถ้าดึงไม่ได้ หรือ Ping ไปหา Database แล้วไม่มีคนตอบรับ
	if err != nil || sqlDB.Ping() != nil {
		dbStatus = "DOWN"
	}

	// 3. สรุปผลลัพธ์
	// ถ้าระบบใดระบบหนึ่งล่ม เราจะส่ง HTTP 503 (Service Unavailable)
	if dbStatus == "DOWN" {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":   "ERROR",
			"api":      apiStatus,
			"database": dbStatus,
		})
		return
	}

	// ถ้าทุกอย่างปกติ ส่ง HTTP 200 (OK)
	c.JSON(http.StatusOK, gin.H{
		"status":   "OK",
		"api":      apiStatus,
		"database": dbStatus,
	})
}
