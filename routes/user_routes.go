package routes

import (
	"log"

	"ecommerce-be/handlers"
	"ecommerce-be/middleware"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(api *gin.RouterGroup) {
	userHandler, err := handlers.NewUserHandler()
	if err != nil {
		log.Printf("⚠️  Warning: User handler không được khởi tạo: %v", err)
		return
	}

	users := api.Group("/users")
	users.Use(middleware.AuthMiddleware()) // Tất cả routes đều yêu cầu auth
	{
		// Customer & Admin routes
		users.GET("/profile", userHandler.GetProfile)
		users.PATCH("/profile", userHandler.UpdateProfile)
		users.PATCH("/change-password", userHandler.ChangePassword)
		users.POST("/upload-avatar", userHandler.UploadAvatar)
		users.DELETE("/delete-avatar", userHandler.DeleteAvatar)

		// Admin only routes
		users.POST("", middleware.RoleMiddleware("admin"), userHandler.CreateUserByAdmin)
		users.POST("/search", middleware.RoleMiddleware("admin"), userHandler.Search)
		users.GET("/:id", middleware.RoleMiddleware("admin"), userHandler.GetUserById)
		users.PATCH("/:id", middleware.RoleMiddleware("admin"), userHandler.UpdateUserByAdmin)
		users.POST("/:id/upload-avatar", middleware.RoleMiddleware("admin"), userHandler.UploadAvatarByAdmin)
		users.PATCH("/:id/change-password", middleware.RoleMiddleware("admin"), userHandler.ChangePasswordByAdmin)
	}
}
