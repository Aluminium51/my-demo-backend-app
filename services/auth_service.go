package services

import (
	"my-go-backend/models"
	"my-go-backend/repositories"
	"my-go-backend/utils"
)

// Interface or Abstract Class
// กำหนดว่า Service นี้ทำอะไรได้บ้าง (ดีมากสำหรับการทำ Unit Test)
type AuthService interface {
	HandleGoogleCallback(googleID, email, name string) (string, *models.User, error)
}

// Struct
// เป็นตัวแทนของ Service จริงๆ ที่จะถูกเรียกใช้ใน Controller
type authService struct {
	userRepo repositories.UserRepository
}

// Constructor Function
// ใช้สร้าง Instance ของ Service และ Inject Repository เข้าไป
func NewAuthService(userRepo repositories.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

// Method
// เช็คว่ามี User หรือยัง ถ้าไม่มีให้สร้างใหม่ แล้วออก JWT
func (s *authService) HandleGoogleCallback(googleID, email, name string) (string, *models.User, error) {
	user, err := s.userRepo.FindByGoogleID(googleID)

	if err != nil {
		// 1. ถ้ายังไม่มีบัญชี (หาไม่เจอ) ให้สร้าง User ใหม่
		user = &models.User{
			Name:     name,
			Email:    email,
			GoogleID: &googleID,
		}
		if err := s.userRepo.Create(user); err != nil {
			return "", nil, err
		}
	}

	// 2. return JWT token
	token, err := utils.GenerateToken(user.ID)
	return token, user, err
}
