package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupRoutes thiết lập tất cả routes cho ứng dụng
func SetupRoutes(r *gin.Engine) {
	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Server is running",
		})
	})

	// API routes
	api := r.Group("/api/v1")
	{
		api.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Ecommerce API v1",
			})
		})

		// Setup các route groups
		SetupAuthRoutes(api)
		SetupCloudinaryRoutes(api)
		SetupUserRoutes(api)
		SetupCategoryRoutes(api)
		SetupProductRoutes(api) // Added product routes
	}
}
