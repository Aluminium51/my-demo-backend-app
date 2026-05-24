package services

import (
	"encoding/json"
	"errors"
	"my-go-backend/config"
	"my-go-backend/models"
	"my-go-backend/repositories"
	"my-go-backend/utils"
	"time"
)

type UserService interface {
	Register(name, email, password string, age int) (*models.User, error)
	Login(email, password string) (string, error)
	GetAll() ([]models.User, error)
	Update(id uint, loggedInUserID uint, input models.User) (*models.User, error)
	Delete(id uint, loggedInUserID uint) error
	GetProfile(id uint) (*models.User, error)
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) Register(name, email, password string, age int) (*models.User, error) {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Name:     name,
		Email:    email,
		Password: &hashedPassword,
		Age:      age,
	}

	// Transaction
	// GORM จะส่งตัวแปร `tx` (ซึ่งเป็น Database ชั่วคราว) มาให้เราใช้ด้านใน
	// err = config.DB.Transaction(func(tx *gorm.DB) error {

	// 	// Step 1: save user to temp DB (tx)
	// 	// ในการใช้งานจริงระดับ Advance เราจะส่ง `tx` เข้าไปใน repo ด้วย แต่เพื่อความเข้าใจง่ายในสเต็ปนี้ เราจะสั่งตรงผ่าน tx ก่อนครับ
	// 	if err := tx.Create(user).Error; err != nil {
	// 		// ถ้าพังตรงนี้ คืนค่า Error ออกไป GORM จะทำการ Rollback อัตโนมัติ
	// 		return err
	// 	}

	// 	// Step 2: Create new task for the user in the same transaction
	// 	defaultTask := &models.Task{
	// 		Title:       "กรอกโปรไฟล์ให้สมบูรณ์",
	// 		Description: "ยินดีต้อนรับ! กรุณาเข้าไปอัปเดตข้อมูลส่วนตัวของคุณให้ครบถ้วน",
	// 		UserID:      user.ID, // ได้เลข ID มาจากสเต็ปที่ 1 แล้ว
	// 	}

	// 	if err := tx.Create(defaultTask).Error; err != nil {
	// 		// สมมติตรงนี้พัง (เช่น พิมพ์ฟิลด์ผิด หรือ DB หลุด)
	// 		// พอคืนค่า Error ตรงนี้ GORM จะแอบไปลบ User ที่เพิ่งสร้างในสเต็ป 1 ทิ้งให้ทันที!
	// 		return err
	// 	}

	// 	// ถ้าทำงานมาถึงตรงนี้โดยไม่มี Error เลย
	// 	// GORM จะสั่ง Commit บันทึกทั้งคู่ลง DB จริงๆ
	// 	return nil
	// })

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) Login(email, password string) (string, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return "", errors.New("อีเมลหรือรหัสผ่านไม่ถูกต้อง")
	}

	if user.Password == nil || !utils.CheckPasswordHash(password, *user.Password) {
		return "", errors.New("อีเมลหรือรหัสผ่านไม่ถูกต้อง")
	}

	return utils.GenerateToken(user.ID, user.Role)
}

func (s *userService) GetAll() ([]models.User, error) {
	// Old Code (ก่อนเพิ่ม Redis Caching)
	// return s.repo.FindAll()

	// New Code (หลังเพิ่ม Redis Caching)
	var users []models.User
	cacheKey := "all_users_data" // ตั้งชื่อคีย์สำหรับเก็บข้อมูลใน Redis (ควรตั้งให้สื่อความหมายและไม่ซ้ำกับคีย์อื่นๆ)

	// 1. Check Redis Cache first
	cachedData, err := config.RDB.Get(config.Ctx, cacheKey).Result()

	if err == nil {
		// 🟢 Cache Hit!
		// ข้อมูลใน Redis จะเก็บเป็น String/JSON เราต้องแปลงกลับเป็น Struct ของ Go
		json.Unmarshal([]byte(cachedData), &users)
		return users, nil
	}

	// 2. 🔴 Cache Miss!
	// Expired or Not Found in Redis
	// we need to fetch from Database
	users, err = s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	// 3. แปลงข้อมูลเป็น JSON เพื่อเก็บใน Redis (Redis เก็บข้อมูลเป็น String)
	usersJSON, _ := json.Marshal(users)

	// 4. Save to Redis
	// พร้อมตั้งเวลาหมดอายุ (TTL - Time to Live) เช่น 1 นาที
	// ถ้าข้อมูลมีการอัปเดตบ่อย ก็ตั้งเวลาน้อยๆ
	// ถ้าข้อมูลไม่ค่อยเปลี่ยนแปลง ก็สามารถตั้งเวลานานขึ้นได้ เช่น 5 นาที หรือ 10 นาที
	config.RDB.Set(config.Ctx, cacheKey, usersJSON, 1*time.Minute)

	return users, nil
}

func (s *userService) Update(id uint, loggedInUserID uint, input models.User) (*models.User, error) {
	// ย้ายด่านเช็ค Authorization (ความเป็นเจ้าของบัญชี) มาไว้ที่ชั้น Service
	if id != loggedInUserID {
		return nil, errors.New("คุณไม่มีสิทธิ์แก้ไขข้อมูลของผู้ใช้งานคนอื่น")
	}

	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("ไม่พบผู้ใช้งาน")
	}

	if err := s.repo.Update(user, input); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) Delete(id uint, loggedInUserID uint) error {
	if id != loggedInUserID {
		return errors.New("คุณไม่มีสิทธิ์ลบข้อมูลของผู้ใช้งานคนอื่น")
	}

	user, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("ไม่พบผู้ใช้งาน")
	}

	return s.repo.Delete(user)
}

func (s *userService) GetProfile(id uint) (*models.User, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("User not Found!")
	}
	return user, nil
}
