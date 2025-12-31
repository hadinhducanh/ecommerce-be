package routes

import (
	"ecommerce-be/handlers"
	"ecommerce-be/middleware"

	"github.com/gin-gonic/gin"
)

// SetupCartRoutes - Thiết lập routes cho cart
func SetupCartRoutes(r *gin.RouterGroup) {
	cart := r.Group("/cart")
	cart.Use(middleware.AuthMiddleware()) // Yêu cầu đăng nhập

	{
		cart.POST("", handlers.AddToCart)              // Thêm sản phẩm vào giỏ hàng
		cart.GET("", handlers.GetCart)                 // Lấy giỏ hàng
		cart.GET("/count", handlers.GetCartItemCount)  // Lấy số lượng items
		cart.PUT("/:id", handlers.UpdateCartItem)      // Cập nhật số lượng
		cart.DELETE("/:id", handlers.DeleteCartItem)   // Xóa một sản phẩm
		cart.DELETE("", handlers.ClearCart)            // Xóa toàn bộ giỏ hàng
	}
}
