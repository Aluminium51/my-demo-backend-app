package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"my-go-backend/services"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// จะถูกเรียกใช้ใน main.go เพื่อสร้าง Route สำหรับ Google OAuth
type AuthController struct {
	authService services.AuthService // เรียกใช้ service ของ Auth
}

func NewAuthController(authService services.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

// ตั้งค่า Google OAuth (เอาไว้ระดับ Package เหมือนเดิมได้ครับ)
var googleOauthConfig = &oauth2.Config{
	RedirectURL:  "http://localhost:8081/api/v1/auth/google/callback",
	ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
	Endpoint:     google.Endpoint,
}

// 1. ส่งผู้ใช้ไปหน้าเลือก Account ของ Google
func (ctrl *AuthController) GoogleLogin(c *gin.Context) {
	url := googleOauthConfig.AuthCodeURL("random_state_string")
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// 2. รับข้อมูลกลับจาก Google
func (ctrl *AuthController) GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถแลกเปลี่ยน Token ได้"})
		return
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถดึงข้อมูลจาก Google ได้"})
		return
	}
	defer resp.Body.Close()

	var googleUser struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	json.NewDecoder(resp.Body).Decode(&googleUser)

	// ---> จุดที่เปลี่ยนไป: โยนข้อมูลให้ AuthService จัดการต่อ <---
	jwtToken, user, err := ctrl.authService.HandleGoogleCallback(googleUser.ID, googleUser.Email, googleUser.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "สร้าง Token ไม่สำเร็จ: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "เข้าสู่ระบบผ่าน Google สำเร็จ",
		"token":   jwtToken,
		"user":    user,
	})
}
