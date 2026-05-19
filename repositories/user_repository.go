package repositories

import (
	"my-go-backend/config"
	"my-go-backend/models"
)

type UserRepository interface {
	Create(user *models.User) error
	FindAll() ([]models.User, error)
	FindByID(id uint) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	FindByGoogleID(googleID string) (*models.User, error)
	Update(user *models.User, input models.User) error
	Delete(user *models.User) error
}

type userRepository struct{}

func NewUserRepository() UserRepository {
	return &userRepository{}
}

func (r *userRepository) Create(user *models.User) error {
	return config.DB.Create(user).Error
}

func (r *userRepository) FindAll() ([]models.User, error) {
	var users []models.User
	// ใส่ Preload ดึงงานพ่วงไปด้วยตามลอจิกเดิม
	err := config.DB.Preload("Tasks").Find(&users).Error
	return users, err
}

func (r *userRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	err := config.DB.First(&user, id).Error
	return &user, err
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := config.DB.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *userRepository) FindByGoogleID(googleID string) (*models.User, error) {
	var user models.User
	err := config.DB.Where("google_id = ?", googleID).First(&user).Error
	return &user, err
}

func (r *userRepository) Update(user *models.User, input models.User) error {
	return config.DB.Model(user).Updates(input).Error
}

func (r *userRepository) Delete(user *models.User) error {
	return config.DB.Delete(user).Error
}
