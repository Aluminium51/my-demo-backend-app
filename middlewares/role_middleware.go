package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ฟังก์ชันนี้จะรับชื่อยศที่อนุญาตให้ผ่านได้ (เช่น "admin")
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. ดึงข้อมูล role ออกมาจาก Context (ซึ่งด่านตรวจ JWT ปกติควรจะใส่มาให้แล้ว)
		userRoleContext, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "ไม่พบข้อมูลสิทธิ์การใช้งาน"})
			c.Abort()
			return
		}

		userRole := userRoleContext.(string)

		// 2. เช็คว่า role ของผู้ใช้คนนี้ ตรงกับยศที่อนุญาตไหม?
		isAllowed := false
		for _, role := range allowedRoles {
			if userRole == role {
				isAllowed = true
				break
			}
		}

		// 3. ถ้าไม่มีสิทธิ์ เตะกลับด้วย 403 Forbidden
		if !isAllowed {
			c.JSON(http.StatusForbidden, gin.H{"error": "คุณไม่มีสิทธิ์เข้าถึงส่วนนี้ (Forbidden)"})
			c.Abort()
			return
		}

		c.Next()
	}
}
