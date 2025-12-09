package routes

import (
	"log"

	"ecommerce-be/handlers"
	"ecommerce-be/middleware"

	"github.com/gin-gonic/gin"
)

func SetupCategoryRoutes(api *gin.RouterGroup) {
	categoryHandler, err := handlers.NewCategoryHandler()
	if err != nil {
		log.Printf("⚠️  Warning: Category handler không được khởi tạo: %v", err)
		return
	}

	categories := api.Group("/categories")
	{
		// Public routes (không yêu cầu auth)
		// Gộp GET all vào POST search - nếu body null/empty thì hiển thị tất cả
		categories.POST("/search", categoryHandler.Search)                  // Tìm kiếm parent categories
		categories.GET("/parents", categoryHandler.GetParentCategories)     // Lấy danh sách parent categories (cho dropdown filter)
		categories.GET("/children", categoryHandler.GetAllChildren)         // Lấy danh sách tất cả child categories (cho dropdown filter)
		categories.POST("/children/search", categoryHandler.SearchChildren) // Tìm kiếm tất cả children
		// Quan trọng: Route cụ thể hơn phải đăng ký trước route generic
		categories.GET("/:id/children", categoryHandler.GetChildren) // Lấy danh sách children của một parent
		categories.GET("/:id", categoryHandler.FindOne)              // Lấy một category theo ID

		// Admin only routes (yêu cầu auth + admin role)
		adminRoutes := categories.Group("")
		adminRoutes.Use(middleware.AuthMiddleware())
		adminRoutes.Use(middleware.RoleMiddleware("admin"))
		{
			adminRoutes.POST("/upload-image", categoryHandler.UploadImage)
			adminRoutes.DELETE("/delete-image", categoryHandler.DeleteImage)
			adminRoutes.POST("", categoryHandler.Create)
			adminRoutes.PATCH("/:id", categoryHandler.Update)
			adminRoutes.PUT("/:id", categoryHandler.Replace)
			adminRoutes.DELETE("/:id", categoryHandler.Remove)
			adminRoutes.DELETE("/:id/hard", categoryHandler.HardDelete)
			// Quản lý children
			adminRoutes.POST("/:id/children", categoryHandler.AddChild)
			adminRoutes.DELETE("/:id/children", categoryHandler.RemoveChild)
		}
	}
}
