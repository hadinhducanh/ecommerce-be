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
		categories.POST("/search", categoryHandler.Search)
		categories.GET("/:id", categoryHandler.FindOne)

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
		}
	}
}
