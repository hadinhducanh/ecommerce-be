package handlers

import (
	"net/http"
	"strconv"

	"ecommerce-be/dto"
	"ecommerce-be/services"

	"github.com/gin-gonic/gin"
)

// AddToCart - Thêm sản phẩm vào giỏ hàng
func AddToCart(c *gin.Context) {
	// Lấy userID từ context (đã set bởi auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req dto.AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cartItem, err := services.AddToCart(userID.(uint), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Thêm sản phẩm vào giỏ hàng thành công",
		"data":    cartItem,
	})
}

// GetCart - Lấy giỏ hàng của user
func GetCart(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	cart, err := services.GetCartByUserID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    cart,
	})
}

// UpdateCartItem - Cập nhật số lượng sản phẩm trong giỏ
func UpdateCartItem(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Lấy cart item ID từ URL
	cartItemIDStr := c.Param("id")
	cartItemID, err := strconv.ParseUint(cartItemIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}

	var req dto.UpdateCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cartItem, err := services.UpdateCartItem(userID.(uint), uint(cartItemID), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Cập nhật giỏ hàng thành công",
		"data":    cartItem,
	})
}

// DeleteCartItem - Xóa sản phẩm khỏi giỏ hàng
func DeleteCartItem(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Lấy cart item ID từ URL
	cartItemIDStr := c.Param("id")
	cartItemID, err := strconv.ParseUint(cartItemIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}

	if err := services.DeleteCartItem(userID.(uint), uint(cartItemID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Xóa sản phẩm khỏi giỏ hàng thành công",
	})
}

// ClearCart - Xóa toàn bộ giỏ hàng
func ClearCart(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := services.ClearCart(userID.(uint)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Xóa toàn bộ giỏ hàng thành công",
	})
}

// GetCartItemCount - Lấy số lượng items trong giỏ hàng
func GetCartItemCount(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	count, err := services.GetCartItemCount(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"count":   count,
	})
}
