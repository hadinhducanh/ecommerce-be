package main

import (
	"log"

	"ecommerce-be/cache"
	"ecommerce-be/config"
	"ecommerce-be/database"
	"ecommerce-be/middleware"
	"ecommerce-be/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	if err := config.LoadConfig(); err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Connect to database
	if err := database.ConnectDB(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.CloseDB()

	// Connect to Redis (optional - nếu không có Redis thì vẫn chạy được)
	if err := cache.ConnectRedis(); err != nil {
		log.Printf("⚠️  Warning: Redis không kết nối được: %v", err)
		log.Println("   Ứng dụng sẽ chạy không có cache. Vui lòng kiểm tra Redis config.")
	} else {
		defer cache.CloseRedis()
	}

	// Setup Gin router
	r := gin.Default()

	// CORS middleware - cho phép frontend truy cập
	r.Use(middleware.CORSMiddleware())

	// Setup all routes
	routes.SetupRoutes(r)

	// Start server
	port := ":" + config.AppConfig.Port
	log.Printf("🚀 Server starting on port %s", port)
	if err := r.Run(port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
