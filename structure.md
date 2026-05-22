Main concept: 
- API Design: RESTful API
- Clean Architecture: แยกความรับผิดชอบของแต่ละชั้นอย่างชัดเจน ทำให้โค้ดมีความยืดหยุ่นและง่ายต่อการบำรุงรักษา
- Database: GORM (ORM สำหรับ Go)
- Authentication: JWT + OAuth2.0
- Role-Based Access Control
- Dependency Injection: ใช้ Dependency Injection เพื่อจัดการกับการสร้างและการจัดการของ
- Unit Testing: ใช้ Go's testing package เพื่อเขียน Unit Tests สำหรับแต่ละชั้นของแอปพลิเคชัน

=========================================================================================================================================

Main concept (Advanced):
- Data Validation
- Database Transactions
- Caching (Redis)
- API Documentation (Swagger/OpenAPI)
- Rate Limiting
- CORS 
- CI/CD (GitHub Actions, Jenkins)

=========================================================================================================================================

my-go-backend/
- controllers/   (ชั้นที่ 1)
- services/      (ชั้นที่ 2)
- repositories/  (ชั้นที่ 3)
- models/        (ข้อมูลตัวกลาง)

Explanation:
- controllers/: ชั้นที่ 1 รับคำขอจากผู้ใช้และส่งต่อไปยัง services ห้ามคิดBusiness logic ในชั้นนี้ ควรทำหน้าที่เพียงแค่รับคำขอและส่งต่อข้อมูลไปยัง services เท่านั้น
- services/: ชั้นที่ 2 จัดการ Business Logic 
- repositories/: ชั้นที่ 3 คุยกับ Database (GORM) หรือเชื่อมต่อกับ API ภายนอกเท่านั้น มีหน้าที่แค่ดึงข้อมูล, บันทึก, ลบ, หรืออัปเดต
- models/: เปรียบเสมือนSchema ทำให้เห็นหน้าตาของ Data ที่ส่งไปมาและไปได้ชัดเจนขึ้น ช่วยให้การทำงานกับข้อมูลมีความชัดเจนและง่ายขึ้น

=========================================================================================================================================

Interface
- บอกแค่ว่าต้องการคนทำอะไรได้บ้าง เช่น ดึงข้อมูลได้, เซฟข้อมูลได้
- ไม่บอกว่าต้องทำยังไง เช่น ดึงข้อมูลจากไหน, เซฟข้อมูลยังไง

Struct
- เปรียบเสมือนการสร้างBlueprint หรือ Template สำหรับสร้าง Object ที่มีลักษณะและคุณสมบัติที่กำหนดไว้
- ใช้ในการจัดเก็บข้อมูลและกำหนดพฤติกรรมของ Object ที่สร้างขึ้นจาก Struct นั้นๆ

=========================================================================================================================================

API documentation (Swagger/OpenAPI)
<!--1. โหลดตัว CLI สำหรับสร้างไฟล์คู่มือ-->
go install github.com/swaggo/swag/cmd/swag@latest

<!--2. โหลดตัวเชื่อม Swagger เข้ากับ Gin-->
go get -u github.com/swaggo/gin-swagger
go get -u github.com/swaggo/files
