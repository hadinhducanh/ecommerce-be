package routes

import (
	"ecommerce-be/handlers"

	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(api *gin.RouterGroup) {
	authHandler := handlers.NewAuthHandler()
	auth := api.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/verify-otp", authHandler.VerifyOtp)
		auth.POST("/resend-otp", authHandler.ResendOtp)
		auth.POST("/refresh", authHandler.RefreshToken)
	}
}

