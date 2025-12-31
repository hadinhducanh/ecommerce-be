package services

import (
	"errors"
	"fmt"

	"ecommerce-be/database"
	"ecommerce-be/dto"
	"ecommerce-be/models"

	"gorm.io/gorm"
)

// AddToCart - Thêm sản phẩm vào giỏ hàng
func AddToCart(userID uint, req dto.AddToCartRequest) (*dto.CartItemResponse, error) {
	// Kiểm tra sản phẩm có tồn tại không
	var product models.Product
	if err := database.DB.Where("id = ? AND is_active = ?", req.ProductID, true).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("sản phẩm không tồn tại hoặc không khả dụng")
		}
		return nil, err
	}

	// Kiểm tra stock
	if product.Stock < req.Quantity {
		return nil, fmt.Errorf("sản phẩm chỉ còn %d sản phẩm trong kho", product.Stock)
	}

	// Kiểm tra xem sản phẩm đã có trong giỏ hàng chưa
	var existingCartItem models.CartItem
	result := database.DB.Where("user_id = ? AND product_id = ?", userID, req.ProductID).First(&existingCartItem)

	if result.Error == nil {
		// Đã có trong giỏ → Cập nhật số lượng
		newQuantity := existingCartItem.Quantity + req.Quantity

		// Kiểm tra stock với số lượng mới
		if product.Stock < newQuantity {
			return nil, fmt.Errorf("sản phẩm chỉ còn %d sản phẩm trong kho", product.Stock)
		}

		existingCartItem.Quantity = newQuantity
		if err := database.DB.Save(&existingCartItem).Error; err != nil {
			return nil, err
		}

		// Preload product để trả về
		database.DB.Preload("Product").First(&existingCartItem, existingCartItem.ID)

		return mapCartItemToResponse(&existingCartItem), nil
	}

	// Chưa có trong giỏ → Tạo mới
	cartItem := models.CartItem{
		UserID:    userID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
	}

	if err := database.DB.Create(&cartItem).Error; err != nil {
		return nil, err
	}

	// Preload product để trả về
	database.DB.Preload("Product").First(&cartItem, cartItem.ID)

	return mapCartItemToResponse(&cartItem), nil
}

// GetCartByUserID - Lấy giỏ hàng của user
func GetCartByUserID(userID uint) (*dto.CartSummaryResponse, error) {
	var cartItems []models.CartItem

	// Lấy tất cả cart items của user với product info
	if err := database.DB.Where("user_id = ?", userID).
		Preload("Product").
		Find(&cartItems).Error; err != nil {
		return nil, err
	}

	// Tính tổng
	var totalItems int
	var totalPrice float64

	items := make([]dto.CartItemResponse, 0)
	for _, item := range cartItems {
		// Chỉ tính những sản phẩm còn active
		if item.Product.IsActive {
			items = append(items, *mapCartItemToResponse(&item))
			totalItems += item.Quantity
			totalPrice += item.Product.Price * float64(item.Quantity)
		}
	}

	return &dto.CartSummaryResponse{
		Items:      items,
		TotalItems: totalItems,
		TotalPrice: totalPrice,
	}, nil
}

// UpdateCartItem - Cập nhật số lượng sản phẩm trong giỏ
func UpdateCartItem(userID, cartItemID uint, req dto.UpdateCartItemRequest) (*dto.CartItemResponse, error) {
	var cartItem models.CartItem

	// Tìm cart item và kiểm tra ownership
	if err := database.DB.Where("id = ? AND user_id = ?", cartItemID, userID).
		Preload("Product").
		First(&cartItem).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("không tìm thấy sản phẩm trong giỏ hàng")
		}
		return nil, err
	}

	// Kiểm tra stock
	if cartItem.Product.Stock < req.Quantity {
		return nil, fmt.Errorf("sản phẩm chỉ còn %d sản phẩm trong kho", cartItem.Product.Stock)
	}

	// Cập nhật số lượng
	cartItem.Quantity = req.Quantity
	if err := database.DB.Save(&cartItem).Error; err != nil {
		return nil, err
	}

	return mapCartItemToResponse(&cartItem), nil
}

// DeleteCartItem - Xóa sản phẩm khỏi giỏ hàng
func DeleteCartItem(userID, cartItemID uint) error {
	// Kiểm tra cart item có tồn tại và thuộc về user không
	var cartItem models.CartItem
	if err := database.DB.Where("id = ? AND user_id = ?", cartItemID, userID).First(&cartItem).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("không tìm thấy sản phẩm trong giỏ hàng")
		}
		return err
	}

	// Xóa cart item
	if err := database.DB.Delete(&cartItem).Error; err != nil {
		return err
	}

	return nil
}

// ClearCart - Xóa toàn bộ giỏ hàng của user
func ClearCart(userID uint) error {
	if err := database.DB.Where("user_id = ?", userID).Delete(&models.CartItem{}).Error; err != nil {
		return err
	}
	return nil
}

// GetCartItemCount - Đếm số lượng items trong giỏ hàng
func GetCartItemCount(userID uint) (int64, error) {
	var count int64
	if err := database.DB.Model(&models.CartItem{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// Helper function để map CartItem sang CartItemResponse
func mapCartItemToResponse(cartItem *models.CartItem) *dto.CartItemResponse {
	return &dto.CartItemResponse{
		ID:        cartItem.ID,
		Quantity:  cartItem.Quantity,
		UserID:    cartItem.UserID,
		ProductID: cartItem.ProductID,
		Product:   *MapProductToResponse(&cartItem.Product),
		CreatedAt: cartItem.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: cartItem.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
