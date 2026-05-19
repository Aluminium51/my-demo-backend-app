package services

import (
	"my-go-backend/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 1. สร้าง MockUserRepository
type MockUserRepository struct {
	// เราสร้างฟังก์ชันจำลองเอาไว้ เพื่อให้ตอนเทสต์เราสามารถควบคุมผลลัพธ์ได้ตามใจชอบ
	MockCreate func(user *models.User) error
}

// 2. ทำให้MockUserRepository ทำงานได้ตาม UserRepository Interface
func (m *MockUserRepository) Create(user *models.User) error {
	// เรียกใช้ฟังก์ชันจำลองที่เราตั้งค่าไว้
	if m.MockCreate != nil {
		return m.MockCreate(user)
	}
	return nil
}

// (ฟังก์ชันอื่นๆ ของ Interface เราก็แค่เขียนหลอกไว้ให้คอมไพล์ผ่าน เพราะเทสต์นี้เราไม่ได้ใช้)
func (m *MockUserRepository) FindAll() ([]models.User, error)                      { return nil, nil }
func (m *MockUserRepository) FindByID(id uint) (*models.User, error)               { return nil, nil }
func (m *MockUserRepository) FindByEmail(email string) (*models.User, error)       { return nil, nil }
func (m *MockUserRepository) FindByGoogleID(googleID string) (*models.User, error) { return nil, nil }
func (m *MockUserRepository) Update(user *models.User, input models.User) error    { return nil }
func (m *MockUserRepository) Delete(user *models.User) error                       { return nil }

// ฟังก์ชันเทสต์ต้องขึ้นต้นด้วยคำว่า Test
func TestUserService_Register_Success(t *testing.T) {
	// Arrange (เตรียมของ)
	// สร้างโกดังปลอม และกำหนดให้ตอนสั่ง Create() ไม่ต้องคืนค่า Error (จำลองว่าเซฟสำเร็จ)
	mockRepo := &MockUserRepository{
		MockCreate: func(user *models.User) error {
			user.ID = 1 // จำลองว่า Database รันเลข ID ให้แล้ว
			return nil
		},
	}
	// ประกอบร่าง! เอาโกดังปลอมไปยัดใส่มือ UserService
	service := NewUserService(mockRepo)

	// Act (ลงมือทดสอบ)
	// สั่งให้ Service ทำการสมัครสมาชิกด้วยรหัสผ่านดิบๆ
	plainPassword := "mysecret123"
	result, err := service.Register("Ryu", "ryu@example.com", plainPassword, 20)

	// Assert (ตรวจสอบผลลัพธ์)
	assert.NoError(t, err)                           // 1. ต้องไม่มี Error เกิดขึ้น
	assert.NotNil(t, result)                         // 2. ต้องได้ข้อมูล User กลับมา
	assert.Equal(t, "Ryu", result.Name)              // 3. ชื่อต้องตรงกับที่ส่งไป
	assert.Equal(t, "ryu@example.com", result.Email) // 4. อีเมลต้องตรง

	// 5. ไฮไลต์สำคัญ! รหัสผ่านในผลลัพธ์ ต้อง "ไม่เท่ากับ" รหัสผ่านดิบ (พิสูจน์ว่าฟังก์ชัน Hash ทำงานจริง)
	assert.NotEqual(t, plainPassword, *result.Password)
}

// go get github.com/stretchr/testify
// go test ./services/ -v
