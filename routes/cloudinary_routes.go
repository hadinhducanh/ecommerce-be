package routes

import (
	"log"

	"ecommerce-be/handlers"
	"ecommerce-be/middleware"

	"github.com/gin-gonic/gin"
)

func SetupCloudinaryRoutes(api *gin.RouterGroup) {
	cloudinaryHandler, err := handlers.NewCloudinaryHandler()
	if err != nil {
		log.Printf("⚠️  Warning: Cloudinary không được khởi tạo: %v", err)
		log.Println("   Cloudinary endpoints sẽ không hoạt động. Vui lòng kiểm tra .env file.")
		return
	}

	cloudinary := api.Group("/cloudinary")
	cloudinary.Use(middleware.AuthMiddleware()) // Yêu cầu đăng nhập
	cloudinary.Use(middleware.RoleMiddleware("admin")) // Chỉ admin
	{
		cloudinary.POST("/upload-image", cloudinaryHandler.UploadImage)
		cloudinary.DELETE("/delete-image", cloudinaryHandler.DeleteImage)
	}
}

