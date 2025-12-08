package middleware

import (
	"net/http"
	"strings"

	"ecommerce-be/database"
	"ecommerce-be/models"
	"ecommerce-be/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware xác thực JWT token và lưu user vào context
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Lấy token từ header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Token không được cung cấp",
			})
			c.Abort()
			return
		}

		// Kiểm tra format "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Token không hợp lệ",
			})
			c.Abort()
			return
		}

		token := parts[1]

		// Validate token
		claims, err := utils.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Token không hợp lệ hoặc đã hết hạn",
			})
			c.Abort()
			return
		}

		// Tìm user trong database
		var user models.User
		if err := database.DB.Where("id = ?", claims.UserID).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "User không tồn tại trong hệ thống",
			})
			c.Abort()
			return
		}

		// Kiểm tra user có active không
		if !user.IsActive {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "Tài khoản đã bị vô hiệu hóa",
			})
			c.Abort()
			return
		}

		// Lưu user vào context
		c.Set("userID", user.ID)
		c.Set("userEmail", user.Email)
		c.Set("userRole", user.Role)
		c.Set("user", user)

		c.Next()
	}
}

// RoleMiddleware kiểm tra role của user
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("userRole")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Bạn cần đăng nhập để truy cập tài nguyên này",
			})
			c.Abort()
			return
		}

		roleStr := userRole.(string)
		hasRole := false
		for _, allowedRole := range allowedRoles {
			if strings.EqualFold(roleStr, allowedRole) {
				hasRole = true
				break
			}
		}

		if !hasRole {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "Bạn không có quyền truy cập. Yêu cầu role: " + strings.Join(allowedRoles, ", "),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
