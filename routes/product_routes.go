package routes

import (
	"log"

	"ecommerce-be/handlers"
	"ecommerce-be/middleware"

	"github.com/gin-gonic/gin"
)

func SetupProductRoutes(api *gin.RouterGroup) {
	productHandler, err := handlers.NewProductHandler()
	if err != nil {
		log.Printf("⚠️  Warning: Product handler không được khởi tạo: %v", err)
		return
	}

	products := api.Group("/products")
	{
		// Public routes (không yêu cầu auth)
		products.POST("/search", productHandler.Search)
		products.GET("/:id", productHandler.FindOne)

		// Admin only routes (yêu cầu auth + admin role)
		adminRoutes := products.Group("")
		adminRoutes.Use(middleware.AuthMiddleware())
		adminRoutes.Use(middleware.RoleMiddleware("admin"))
		{
			adminRoutes.POST("/upload-image", productHandler.UploadImage)
			adminRoutes.DELETE("/delete-image", productHandler.DeleteImage)
			adminRoutes.POST("", productHandler.Create)
			adminRoutes.PATCH("/:id", productHandler.Update)
			adminRoutes.PUT("/:id", productHandler.Replace)
			adminRoutes.DELETE("/:id", productHandler.Remove)
			adminRoutes.DELETE("/:id/hard", productHandler.HardDelete)
		}
	}
}

