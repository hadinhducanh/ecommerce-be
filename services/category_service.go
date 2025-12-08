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

type CategoryService struct{}

func NewCategoryService() *CategoryService {
	return &CategoryService{}
}

// Create tạo category mới
func (s *CategoryService) Create(req dto.CreateCategoryRequest) (*models.Category, error) {
	// Kiểm tra xem category với tên này đã tồn tại chưa
	var existingCategory models.Category
	if err := database.DB.Where("name = ?", req.Name).First(&existingCategory).Error; err == nil {
		return nil, errors.New("danh mục với tên này đã tồn tại")
	}

	// Tạo category
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	category := models.Category{
		Name:          strings.TrimSpace(req.Name),
		NameEn:        req.NameEn,
		Description:   req.Description,
		DescriptionEn: req.DescriptionEn,
		Image:         req.Image,
		IsActive:      isActive,
	}

	if err := database.DB.Create(&category).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "UNIQUE constraint") {
			return nil, errors.New("danh mục với tên này đã tồn tại")
		}
		return nil, errors.New("không thể tạo danh mục")
	}

	// Invalidate cache
	s.invalidateCategoryCache()

	return &category, nil
}

// FindAll lấy tất cả categories
func (s *CategoryService) FindAll(includeInactive bool, language string) ([]models.Category, error) {
	// Tạo cache key
	cacheKey := fmt.Sprintf("%s:all:%v:%s", cache.CategoryListKey, includeInactive, language)

	// Thử lấy từ cache
	var categories []models.Category
	if cache.RedisClient != nil {
		if err := cache.Get(cacheKey, &categories); err == nil {
			return categories, nil
		}
	}

	// Nếu không có trong cache, lấy từ database
	query := database.DB.Order("created_at DESC")

	if !includeInactive {
		query = query.Where("is_active = ?", true)
	}

	if err := query.Find(&categories).Error; err != nil {
		return nil, errors.New("không thể lấy danh sách danh mục")
	}

	// Transform với language nếu có
	if language == "en" || language == "vi" {
		for i := range categories {
			categories[i] = s.transformCategory(categories[i], language)
		}
	}

	// Lưu vào cache (5 phút)
	if cache.RedisClient != nil {
		cache.Set(cacheKey, categories, 5*time.Minute)
	}

	return categories, nil
}

// FindOne lấy một category theo ID
func (s *CategoryService) FindOne(id uint, includeInactive bool, language string) (*models.Category, error) {
	// Tạo cache key
	cacheKey := fmt.Sprintf("%s%d:%v:%s", cache.CategoryKeyPrefix, id, includeInactive, language)

	// Thử lấy từ cache
	var category models.Category
	if cache.RedisClient != nil {
		if err := cache.Get(cacheKey, &category); err == nil {
			return &category, nil
		}
	}

	// Nếu không có trong cache, lấy từ database
	query := database.DB.Where("id = ?", id)

	if !includeInactive {
		query = query.Where("is_active = ?", true)
	}

	if err := query.First(&category).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("không tìm thấy danh mục với ID %d", id)
		}
		return nil, errors.New("không thể lấy danh mục")
	}

	// Transform với language nếu có
	if language == "en" || language == "vi" {
		category = s.transformCategory(category, language)
	}

	// Lưu vào cache (10 phút)
	if cache.RedisClient != nil {
		cache.Set(cacheKey, category, 10*time.Minute)
	}

	return &category, nil
}

// Search tìm kiếm và lọc categories
func (s *CategoryService) Search(req dto.SearchCategoryRequest, language string) (*dto.CategoryPaginationResponse, error) {
	// Default values
	sortBy := "createdAt"
	sortOrder := "DESC"
	page := 1
	limit := 50

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
	query := database.DB.Model(&models.Category{})

	// Filter theo name (partial match, case-insensitive)
	if req.Name != nil && *req.Name != "" {
		normalizedSearch := s.normalizeForSearch(*req.Name)

		// Lấy tất cả categories để filter trong memory (vì PostgreSQL không có built-in function để xóa dấu tiếng Việt)
		var allCategories []models.Category
		if err := database.DB.Select("id", "name", "name_en").Find(&allCategories).Error; err != nil {
			return nil, errors.New("không thể tìm kiếm danh mục")
		}

		// Filter categories có tên chứa search term
		var matchedIDs []uint
		for _, cat := range allCategories {
			nameVi := s.normalizeForSearch(cat.Name)
			nameEn := ""
			if cat.NameEn != nil {
				nameEn = s.normalizeForSearch(*cat.NameEn)
			}
			if strings.Contains(nameVi, normalizedSearch) || (nameEn != "" && strings.Contains(nameEn, normalizedSearch)) {
				matchedIDs = append(matchedIDs, cat.ID)
			}
		}

		if len(matchedIDs) > 0 {
			query = query.Where("id IN ?", matchedIDs)
		} else {
			// Nếu không tìm thấy → trả về empty result
			query = query.Where("1 = 0") // Always false condition
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

	// Count total
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, errors.New("không thể đếm số lượng danh mục")
	}

	// Sort - map field names to database columns (giống NestJS: category.${sortBy})
	sortByColumn := s.mapSortFieldToColumn(sortBy)
	orderBy := fmt.Sprintf("%s %s", sortByColumn, sortOrder)
	query = query.Order(orderBy)

	// Pagination
	offset := (page - 1) * limit
	var categories []models.Category
	if err := query.Offset(offset).Limit(limit).Find(&categories).Error; err != nil {
		return nil, errors.New("không thể lấy danh sách danh mục")
	}

	// Transform với language nếu có
	if language == "en" || language == "vi" {
		for i := range categories {
			categories[i] = s.transformCategory(categories[i], language)
		}
	}

	// Convert to response
	categoryResponses := make([]dto.CategoryResponse, len(categories))
	for i, cat := range categories {
		categoryResponses[i] = dto.CategoryResponse{
			ID:            cat.ID,
			Name:          cat.Name,
			NameEn:        cat.NameEn,
			Description:   cat.Description,
			DescriptionEn: cat.DescriptionEn,
			Image:         cat.Image,
			IsActive:      cat.IsActive,
			CreatedAt:     cat.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:     cat.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &dto.CategoryPaginationResponse{
		Data:       categoryResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

// Update cập nhật category
func (s *CategoryService) Update(id uint, req dto.UpdateCategoryRequest) (*models.Category, error) {
	var category models.Category
	if err := database.DB.Where("id = ?", id).First(&category).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("không tìm thấy danh mục với ID %d", id)
		}
		return nil, errors.New("không thể lấy danh mục")
	}

	// Nếu đổi tên, kiểm tra xem tên mới đã tồn tại chưa
	if req.Name != nil && *req.Name != category.Name {
		var existingCategory models.Category
		if err := database.DB.Where("name = ?", *req.Name).First(&existingCategory).Error; err == nil {
			return nil, errors.New("danh mục với tên này đã tồn tại")
		}
	}

	// Cập nhật các fields
	if req.Name != nil {
		category.Name = strings.TrimSpace(*req.Name)
	}
	if req.NameEn != nil {
		category.NameEn = req.NameEn
	}
	if req.Description != nil {
		category.Description = req.Description
	}
	if req.DescriptionEn != nil {
		category.DescriptionEn = req.DescriptionEn
	}
	if req.Image != nil {
		category.Image = req.Image
	}
	if req.IsActive != nil {
		category.IsActive = *req.IsActive
	}

	if err := database.DB.Save(&category).Error; err != nil {
		return nil, errors.New("không thể cập nhật danh mục")
	}

	// Invalidate cache
	s.invalidateCategoryCache()
	s.invalidateCategoryCacheByID(id)

	return &category, nil
}

// Remove xóa category (soft delete - set isActive = false)
func (s *CategoryService) Remove(id uint) error {
	var category models.Category
	if err := database.DB.Where("id = ?", id).Preload("Products").First(&category).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("không tìm thấy danh mục với ID %d", id)
		}
		return errors.New("không thể lấy danh mục")
	}

	if !category.IsActive {
		return errors.New("danh mục này đã bị vô hiệu hóa trước đó")
	}

	// Kiểm tra nếu có products, không cho phép xóa (giống NestJS: category.products?.length || 0)
	productCount := len(category.Products)
	if productCount > 0 {
		return fmt.Errorf("không thể xóa danh mục này vì có %d sản phẩm liên quan. Vui lòng xóa hoặc di chuyển các sản phẩm sang danh mục khác trước", productCount)
	}

	// Soft delete - set isActive = false
	category.IsActive = false
	if err := database.DB.Save(&category).Error; err != nil {
		return errors.New("không thể xóa danh mục")
	}

	// Invalidate cache
	s.invalidateCategoryCache()
	s.invalidateCategoryCacheByID(id)

	return nil
}

// HardDelete xóa vĩnh viễn category
func (s *CategoryService) HardDelete(id uint) error {
	var category models.Category
	if err := database.DB.Where("id = ?", id).Preload("Products").First(&category).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("không tìm thấy danh mục với ID %d", id)
		}
		return errors.New("không thể lấy danh mục")
	}

	// Kiểm tra nếu có products, không cho phép xóa (giống NestJS: category.products?.length || 0)
	productCount := len(category.Products)
	if productCount > 0 {
		return fmt.Errorf("không thể xóa vĩnh viễn danh mục này vì có %d sản phẩm liên quan. Vui lòng xóa hoặc di chuyển các sản phẩm trước", productCount)
	}

	// Hard delete
	if err := database.DB.Delete(&category).Error; err != nil {
		return errors.New("không thể xóa danh mục")
	}

	// Invalidate cache
	s.invalidateCategoryCache()
	s.invalidateCategoryCacheByID(id)

	return nil
}

// Helper methods

// transformCategory transform category với language
func (s *CategoryService) transformCategory(cat models.Category, language string) models.Category {
	if language == "en" {
		if cat.NameEn != nil && *cat.NameEn != "" {
			cat.Name = *cat.NameEn
		}
		if cat.DescriptionEn != nil && *cat.DescriptionEn != "" {
			cat.Description = cat.DescriptionEn
		}
	}
	// language == "vi" hoặc không truyền → giữ nguyên (tiếng Việt)
	return cat
}

// normalizeForSearch normalize string để tìm kiếm (xóa dấu, lowercase)
func (s *CategoryService) normalizeForSearch(str string) string {
	// Đơn giản hóa: chỉ lowercase và trim
	// Trong production có thể dùng thư viện để xóa dấu tiếng Việt
	return strings.ToLower(strings.TrimSpace(str))
}

// invalidateCategoryCache xóa tất cả cache liên quan đến categories
func (s *CategoryService) invalidateCategoryCache() {
	if cache.RedisClient == nil {
		return
	}
	// Xóa cache list
	cache.DeletePattern("categories:*")
}

// invalidateCategoryCacheByID xóa cache của một category cụ thể
func (s *CategoryService) invalidateCategoryCacheByID(id uint) {
	if cache.RedisClient == nil {
		return
	}
	// Xóa cache của category này (tất cả các language và includeInactive variants)
	cache.DeletePattern(fmt.Sprintf("%s%d:*", cache.CategoryKeyPrefix, id))
}

// mapSortFieldToColumn map sort field name to database column name
func (s *CategoryService) mapSortFieldToColumn(field string) string {
	fieldMap := map[string]string{
		"id":        "id",
		"name":      "name",
		"createdAt": "created_at",
		"updatedAt": "updated_at",
	}
	if column, ok := fieldMap[field]; ok {
		return column
	}
	return "created_at" // default
}
