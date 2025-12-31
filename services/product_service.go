package services

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"ecommerce-be/cache"
	"ecommerce-be/database"
	"ecommerce-be/dto"
	"ecommerce-be/models"

	"gorm.io/gorm"
)

type ProductService struct{}

func NewProductService() *ProductService {
	return &ProductService{}
}

// Create tạo product mới
func (s *ProductService) Create(req dto.CreateProductRequest) (*models.Product, error) {
	// Kiểm tra category có tồn tại không
	var category models.Category
	if err := database.DB.Where("id = ?", req.CategoryID).First(&category).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("không tìm thấy danh mục với ID %d", req.CategoryID)
		}
		return nil, errors.New("không thể kiểm tra danh mục")
	}

	if !category.IsActive {
		return nil, errors.New("không thể tạo sản phẩm trong danh mục đã bị vô hiệu hóa")
	}

	// Kiểm tra SKU nếu có (SKU phải unique)
	if req.SKU != nil && *req.SKU != "" {
		var existingProduct models.Product
		if err := database.DB.Where("sku = ?", *req.SKU).First(&existingProduct).Error; err == nil {
			return nil, errors.New("SKU đã tồn tại")
		}
	}

	// Tạo product
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	product := models.Product{
		Name:          strings.TrimSpace(req.Name),
		NameEn:        req.NameEn,
		Description:   req.Description,
		DescriptionEn: req.DescriptionEn,
		Price:         req.Price,
		Stock:         req.Stock,
		Image:         req.Image,
		Images:        req.Images,
		CategoryID:    req.CategoryID,
		SKU:           req.SKU,
		IsActive:      isActive,
		Sold:          0,
		Rating:        0,
		ReviewCount:   0,
	}

	if err := database.DB.Create(&product).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "UNIQUE constraint") {
			return nil, errors.New("SKU đã tồn tại")
		}
		return nil, errors.New("không thể tạo sản phẩm")
	}

	// Invalidate cache
	s.invalidateProductCache()

	return &product, nil
}

// Search tìm kiếm và lọc products
func (s *ProductService) Search(req dto.SearchProductRequest, language string) (*dto.ProductPaginationResponse, error) {
	// Default values
	sortBy := "createdAt"
	sortOrder := "DESC"
	page := 1
	limit := 10

	if req.SortBy != nil {
		sortBy = *req.SortBy
	}
	if req.SortOrder != nil {
		sortOrder = *req.SortOrder
	}
	if req.Page != nil {
		page = *req.Page
	}
	if req.Limit != nil {
		limit = *req.Limit
	}

	// Build query
	query := database.DB.Model(&models.Product{}).Preload("Category")

	// Search theo name (partial match, case-insensitive, không phân biệt dấu)
	// Tìm kiếm cả tiếng Việt và tiếng Anh
	if req.Name != nil && *req.Name != "" {
		normalizedSearch := s.normalizeForSearch(*req.Name)

		// Lấy tất cả products để filter trong memory (vì PostgreSQL không có built-in function để xóa dấu tiếng Việt)
		var allProducts []models.Product
		if err := database.DB.Select("id", "name", "name_en").Find(&allProducts).Error; err != nil {
			return nil, errors.New("không thể tìm kiếm sản phẩm")
		}

		// Filter products có tên chứa search term
		var matchedIDs []uint
		for _, prod := range allProducts {
			nameVi := s.normalizeForSearch(prod.Name)
			nameEn := ""
			if prod.NameEn != nil {
				nameEn = s.normalizeForSearch(*prod.NameEn)
			}
			if strings.Contains(nameVi, normalizedSearch) || (nameEn != "" && strings.Contains(nameEn, normalizedSearch)) {
				matchedIDs = append(matchedIDs, prod.ID)
			}
		}

		if len(matchedIDs) > 0 {
			query = query.Where("id IN ?", matchedIDs)
		} else {
			query = query.Where("1 = 0") // Always false condition
		}
	}

	// Filter theo categoryId (exact match) - ưu tiên nếu có cả categoryId và parentCategoryId
	if req.CategoryID != nil {
		query = query.Where("category_id = ?", *req.CategoryID)
	} else if req.ParentCategoryID != nil {
		// Filter theo parent category - lấy tất cả sản phẩm của các danh mục con
		// Bước 1: Lấy tất cả child categories của parent category
		var relations []models.CategoryChild
		if err := database.DB.Where("parent_id = ? AND deleted_at IS NULL", *req.ParentCategoryID).Find(&relations).Error; err != nil {
			return nil, errors.New("không thể lấy danh sách danh mục con")
		}

		if len(relations) == 0 {
			// Không có child categories nào → không có sản phẩm nào
			query = query.Where("1 = 0") // Always false condition
		} else {
			// Lấy danh sách ID các child categories
			childIDs := make([]uint, len(relations))
			for i, rel := range relations {
				childIDs[i] = rel.ChildID
			}
			// Filter products theo các child categories
			query = query.Where("category_id IN ?", childIDs)
		}
	}

	// Filter theo isActive (hỗ trợ boolean hoặc array)
	// Logic giống NestJS:
	// - Không truyền isActive → Lấy tất cả (active + inactive)
	// - isActive = true hoặc [true] → Chỉ lấy active
	// - isActive = false hoặc [false] → Chỉ lấy inactive
	// - isActive = [true, false] → Lấy tất cả (không filter)
	if req.IsActive != nil {
		switch v := req.IsActive.(type) {
		case bool:
			// Boolean đơn
			query = query.Where("is_active = ?", v)
		case []interface{}:
			// Array - kiểm tra nếu có cả true và false thì không filter
			uniqueValues := make(map[bool]bool)
			for _, item := range v {
				if b, ok := item.(bool); ok {
					uniqueValues[b] = true
				}
			}
			// Nếu có cả true và false → không filter (lấy tất cả)
			if len(uniqueValues) == 1 {
				// Chỉ có 1 giá trị duy nhất → filter theo giá trị đó
				for value := range uniqueValues {
					query = query.Where("is_active = ?", value)
				}
			}
			// Nếu có cả true và false (len == 2) → không filter, lấy tất cả
		case []bool:
			// Array of bools
			uniqueValues := make(map[bool]bool)
			for _, b := range v {
				uniqueValues[b] = true
			}
			// Nếu có cả true và false → không filter (lấy tất cả)
			if len(uniqueValues) == 1 {
				// Chỉ có 1 giá trị duy nhất → filter theo giá trị đó
				for value := range uniqueValues {
					query = query.Where("is_active = ?", value)
				}
			}
			// Nếu có cả true và false (len == 2) → không filter, lấy tất cả
		}
	}

	// Filter theo khoảng giá
	if req.MinPrice != nil {
		query = query.Where("price >= ?", *req.MinPrice)
	}

	if req.MaxPrice != nil {
		query = query.Where("price <= ?", *req.MaxPrice)
	}

	// Filter theo số lượng tồn kho (chỉ sản phẩm còn hàng)
	if req.InStock != nil && *req.InStock {
		query = query.Where("stock > 0")
	}

	// Count total
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, errors.New("không thể đếm số lượng sản phẩm")
	}

	// Sort - map field names to database columns
	sortByColumn := s.mapSortFieldToColumn(sortBy)
	orderBy := fmt.Sprintf("%s %s", sortByColumn, sortOrder)
	query = query.Order(orderBy)

	// Pagination
	offset := (page - 1) * limit
	var products []models.Product
	if err := query.Offset(offset).Limit(limit).Find(&products).Error; err != nil {
		return nil, errors.New("không thể lấy danh sách sản phẩm")
	}

	// Transform với language nếu có
	if language == "en" || language == "vi" {
		for i := range products {
			products[i] = s.transformProduct(products[i], language)
		}
	}

	// Convert to response
	productResponses := make([]dto.ProductResponse, len(products))
	for i, prod := range products {
		var categoryResp *dto.CategoryResponse
		if prod.Category.ID > 0 {
			categoryResp = &dto.CategoryResponse{
				ID:            prod.Category.ID,
				Name:          prod.Category.Name,
				NameEn:        prod.Category.NameEn,
				Description:   prod.Category.Description,
				DescriptionEn: prod.Category.DescriptionEn,
				Image:         prod.Category.Image,
				IsActive:      prod.Category.IsActive,
				CreatedAt:     prod.Category.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
				UpdatedAt:     prod.Category.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
			}
		}

		productResponses[i] = dto.ProductResponse{
			ID:            prod.ID,
			Name:          prod.Name,
			NameEn:        prod.NameEn,
			Description:   prod.Description,
			DescriptionEn: prod.DescriptionEn,
			Price:         prod.Price,
			Stock:         prod.Stock,
			Image:         prod.Image,
			Images:        prod.Images,
			Sold:          prod.Sold,
			Rating:        prod.Rating,
			ReviewCount:   prod.ReviewCount,
			IsActive:      prod.IsActive,
			SKU:           prod.SKU,
			CategoryID:    prod.CategoryID,
			Category:      categoryResp,
			CreatedAt:     prod.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:     prod.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &dto.ProductPaginationResponse{
		Data:       productResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

// FindOne lấy một product theo ID
func (s *ProductService) FindOne(id uint, includeInactive bool, language string) (*models.Product, error) {
	// Tạo cache key
	cacheKey := fmt.Sprintf("%s%d:%v:%s", cache.ProductKeyPrefix, id, includeInactive, language)

	// Thử lấy từ cache
	var product models.Product
	if cache.RedisClient != nil {
		if err := cache.Get(cacheKey, &product); err == nil {
			return &product, nil
		}
	}

	// Nếu không có trong cache, lấy từ database
	query := database.DB.Where("id = ?", id).Preload("Category")

	if !includeInactive {
		query = query.Where("is_active = ?", true)
	}

	if err := query.First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("không tìm thấy sản phẩm với ID %d", id)
		}
		return nil, errors.New("không thể lấy sản phẩm")
	}

	// Transform với language nếu có
	if language == "en" || language == "vi" {
		product = s.transformProduct(product, language)
	}

	// Lưu vào cache (10 phút)
	if cache.RedisClient != nil {
		cache.Set(cacheKey, product, 10*time.Minute)
	}

	return &product, nil
}

// Update cập nhật product
func (s *ProductService) Update(id uint, req interface{}) (*models.Product, error) {
	var product models.Product
	if err := database.DB.Where("id = ?", id).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("không tìm thấy sản phẩm với ID %d", id)
		}
		return nil, errors.New("không thể lấy sản phẩm")
	}

	// Kiểm tra category nếu có thay đổi
	if updateReq, ok := req.(dto.UpdateProductRequest); ok && updateReq.CategoryID != nil && *updateReq.CategoryID != product.CategoryID {
		var category models.Category
		if err := database.DB.Where("id = ?", *updateReq.CategoryID).First(&category).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, fmt.Errorf("không tìm thấy danh mục với ID %d", *updateReq.CategoryID)
			}
			return nil, errors.New("không thể kiểm tra danh mục")
		}

		if !category.IsActive {
			return nil, errors.New("không thể chuyển sản phẩm sang danh mục đã bị vô hiệu hóa")
		}
	} else if updateReqFull, ok := req.(dto.UpdateProductFullRequest); ok && updateReqFull.CategoryID != product.CategoryID {
		var category models.Category
		if err := database.DB.Where("id = ?", updateReqFull.CategoryID).First(&category).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, fmt.Errorf("không tìm thấy danh mục với ID %d", updateReqFull.CategoryID)
			}
			return nil, errors.New("không thể kiểm tra danh mục")
		}

		if !category.IsActive {
			return nil, errors.New("không thể chuyển sản phẩm sang danh mục đã bị vô hiệu hóa")
		}
	}

	// Kiểm tra SKU nếu có thay đổi
	if updateReq, ok := req.(dto.UpdateProductRequest); ok && updateReq.SKU != nil && *updateReq.SKU != "" {
		if product.SKU == nil || *updateReq.SKU != *product.SKU {
			var existingProduct models.Product
			if err := database.DB.Where("sku = ?", *updateReq.SKU).First(&existingProduct).Error; err == nil {
				return nil, errors.New("SKU đã tồn tại")
			}
		}
	} else if updateReqFull, ok := req.(dto.UpdateProductFullRequest); ok && updateReqFull.SKU != nil && *updateReqFull.SKU != "" {
		if product.SKU == nil || *updateReqFull.SKU != *product.SKU {
			var existingProduct models.Product
			if err := database.DB.Where("sku = ?", *updateReqFull.SKU).First(&existingProduct).Error; err == nil {
				return nil, errors.New("SKU đã tồn tại")
			}
		}
	}

	// Cập nhật các fields
	if updateReq, ok := req.(dto.UpdateProductRequest); ok {
		if updateReq.Name != nil {
			product.Name = strings.TrimSpace(*updateReq.Name)
		}
		if updateReq.NameEn != nil {
			product.NameEn = updateReq.NameEn
		}
		if updateReq.Description != nil {
			product.Description = updateReq.Description
		}
		if updateReq.DescriptionEn != nil {
			product.DescriptionEn = updateReq.DescriptionEn
		}
		if updateReq.Price != nil {
			product.Price = *updateReq.Price
		}
		if updateReq.Stock != nil {
			product.Stock = *updateReq.Stock
		}
		if updateReq.Image != nil {
			product.Image = updateReq.Image
		}
		if updateReq.Images != nil {
			product.Images = updateReq.Images
		}
		if updateReq.CategoryID != nil {
			product.CategoryID = *updateReq.CategoryID
		}
		if updateReq.SKU != nil {
			product.SKU = updateReq.SKU
		}
		if updateReq.IsActive != nil {
			product.IsActive = *updateReq.IsActive
		}
	} else if updateReqFull, ok := req.(dto.UpdateProductFullRequest); ok {
		product.Name = strings.TrimSpace(updateReqFull.Name)
		product.NameEn = updateReqFull.NameEn
		product.Description = updateReqFull.Description
		product.DescriptionEn = updateReqFull.DescriptionEn
		product.Price = updateReqFull.Price
		product.Stock = updateReqFull.Stock
		product.Image = updateReqFull.Image
		product.Images = updateReqFull.Images
		product.CategoryID = updateReqFull.CategoryID
		product.SKU = updateReqFull.SKU
		if updateReqFull.IsActive != nil {
			product.IsActive = *updateReqFull.IsActive
		}
	}

	if err := database.DB.Save(&product).Error; err != nil {
		return nil, errors.New("không thể cập nhật sản phẩm")
	}

	// Invalidate cache
	s.invalidateProductCache()
	s.invalidateProductCacheByID(id)

	return &product, nil
}

// Remove xóa product (soft delete - set isActive = false)
func (s *ProductService) Remove(id uint) error {
	var product models.Product
	if err := database.DB.Where("id = ?", id).Preload("OrderItems").Preload("CartItems").First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("không tìm thấy sản phẩm với ID %d", id)
		}
		return errors.New("không thể lấy sản phẩm")
	}

	if !product.IsActive {
		return errors.New("sản phẩm này đã bị vô hiệu hóa trước đó")
	}

	// Trong NestJS, vẫn cho phép soft delete dù có order/cart items
	// Không cần kiểm tra orderItemCount và cartItemCount cho soft delete

	// Soft delete
	product.IsActive = false
	if err := database.DB.Save(&product).Error; err != nil {
		return errors.New("không thể xóa sản phẩm")
	}

	// Invalidate cache
	s.invalidateProductCache()
	s.invalidateProductCacheByID(id)

	return nil
}

// HardDelete xóa vĩnh viễn product
func (s *ProductService) HardDelete(id uint) error {
	var product models.Product
	if err := database.DB.Where("id = ?", id).Preload("OrderItems").Preload("CartItems").First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("không tìm thấy sản phẩm với ID %d", id)
		}
		return errors.New("không thể lấy sản phẩm")
	}

	// Kiểm tra nếu có order items, không cho phép xóa
	orderItemCount := len(product.OrderItems)
	if orderItemCount > 0 {
		return fmt.Errorf("không thể xóa vĩnh viễn sản phẩm này vì có %d đơn hàng liên quan. Vui lòng xử lý các đơn hàng trước", orderItemCount)
	}

	// Kiểm tra nếu có cart items, không cho phép xóa
	cartItemCount := len(product.CartItems)
	if cartItemCount > 0 {
		return fmt.Errorf("không thể xóa vĩnh viễn sản phẩm này vì có %d giỏ hàng liên quan. Vui lòng xử lý các giỏ hàng trước", cartItemCount)
	}

	// Hard delete
	if err := database.DB.Delete(&product).Error; err != nil {
		return errors.New("không thể xóa sản phẩm")
	}

	// Invalidate cache
	s.invalidateProductCache()
	s.invalidateProductCacheByID(id)

	return nil
}

// Helper methods

// transformProduct transform product với language
func (s *ProductService) transformProduct(prod models.Product, language string) models.Product {
	if language == "en" {
		if prod.NameEn != nil && *prod.NameEn != "" {
			prod.Name = *prod.NameEn
		}
		if prod.DescriptionEn != nil && *prod.DescriptionEn != "" {
			prod.Description = prod.DescriptionEn
		}
	}
	// language == "vi" hoặc không truyền → giữ nguyên (tiếng Việt)
	return prod
}

// normalizeForSearch normalize string để tìm kiếm
func (s *ProductService) normalizeForSearch(str string) string {
	return strings.ToLower(strings.TrimSpace(str))
}

// mapSortFieldToColumn map sort field name to database column name
func (s *ProductService) mapSortFieldToColumn(field string) string {
	fieldMap := map[string]string{
		"id":        "id",
		"name":      "name",
		"price":     "price",
		"stock":     "stock",
		"createdAt": "created_at",
		"updatedAt": "updated_at",
	}
	if column, ok := fieldMap[field]; ok {
		return column
	}
	return "created_at" // default
}

// invalidateProductCache xóa tất cả cache liên quan đến products
func (s *ProductService) invalidateProductCache() {
	if cache.RedisClient == nil {
		return
	}
	// Xóa tất cả keys bắt đầu bằng "product:"
	cache.DeletePattern(cache.ProductKeyPrefix + "*")
}

// invalidateProductCacheByID xóa cache của một product cụ thể
func (s *ProductService) invalidateProductCacheByID(id uint) {
	if cache.RedisClient == nil {
		return
	}
	// Xóa cache của product này (tất cả các language và includeInactive variants)
	cache.DeletePattern(fmt.Sprintf("%s%d:*", cache.ProductKeyPrefix, id))
}

// MapProductToResponse converts Product model to ProductResponse DTO
func MapProductToResponse(product *models.Product) *dto.ProductResponse {
	return &dto.ProductResponse{
		ID:            product.ID,
		Name:          product.Name,
		NameEn:        product.NameEn,
		Description:   product.Description,
		DescriptionEn: product.DescriptionEn,
		Price:         product.Price,
		Stock:         product.Stock,
		Image:         product.Image,
		Images:        product.Images,
		Sold:          product.Sold,
		Rating:        product.Rating,
		ReviewCount:   product.ReviewCount,
		IsActive:      product.IsActive,
		SKU:           product.SKU,
		CategoryID:    product.CategoryID,
		CreatedAt:     product.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:     product.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
