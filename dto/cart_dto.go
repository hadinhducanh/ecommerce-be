package dto

// AddToCartRequest - Request để thêm sản phẩm vào giỏ hàng
type AddToCartRequest struct {
	ProductID uint `json:"productId" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required,min=1"`
}

// UpdateCartItemRequest - Request để cập nhật số lượng
type UpdateCartItemRequest struct {
	Quantity int `json:"quantity" binding:"required,min=1"`
}

// CartItemResponse - Response cho một cart item
type CartItemResponse struct {
	ID        uint            `json:"id"`
	Quantity  int             `json:"quantity"`
	UserID    uint            `json:"userId"`
	ProductID uint            `json:"productId"`
	Product   ProductResponse `json:"product"`
	CreatedAt string          `json:"createdAt"`
	UpdatedAt string          `json:"updatedAt"`
}

// CartSummaryResponse - Tổng hợp thông tin giỏ hàng
type CartSummaryResponse struct {
	Items      []CartItemResponse `json:"items"`
	TotalItems int                `json:"totalItems"`
	TotalPrice float64            `json:"totalPrice"`
}
