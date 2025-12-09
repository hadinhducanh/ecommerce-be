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

// Create tạo category mới (mặc định là danh mục cha - root category)
func (s *CategoryService) Create(req dto.CreateCategoryRequest) (*models.Category, error) {
	// Kiểm tra xem category với tên này đã tồn tại chưa (chỉ kiểm tra các record chưa bị soft delete)
	var existingCategory models.Category
	if err := database.DB.Where("name = ? AND deleted_at IS NULL", req.Name).First(&existingCategory).Error; err == nil {
		return nil, errors.New("danh mục với tên này đã tồn tại")
	}

	// Tạo category (mặc định là root category, không có parent)
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

	// Nếu có parentId, tự động thêm vào parent
	if req.ParentID != nil {
		// Validate parent tồn tại và active
		var parent models.Category
		if err := database.DB.Where("id = ?", *req.ParentID).First(&parent).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Xóa category vừa tạo nếu parent không tồn tại
				database.DB.Delete(&category)
				return nil, fmt.Errorf("không tìm thấy danh mục cha với ID %d", *req.ParentID)
			}
			database.DB.Delete(&category)
			return nil, errors.New("không thể kiểm tra danh mục cha")
		}
		if !parent.IsActive {
			database.DB.Delete(&category)
			return nil, errors.New("danh mục cha đã bị vô hiệu hóa")
		}

		// Không cho phép category là parent của chính nó
		if category.ID == *req.ParentID {
			database.DB.Delete(&category)
			return nil, errors.New("danh mục không thể là cha của chính nó")
		}

		// Tạo quan hệ parent-child
		categoryChild := models.CategoryChild{
			ParentID: *req.ParentID,
			ChildID:  category.ID,
		}

		if err := database.DB.Create(&categoryChild).Error; err != nil {
			// Xóa category vừa tạo nếu không thể tạo quan hệ
			database.DB.Delete(&category)
			if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "UNIQUE constraint") {
				return nil, errors.New("danh mục con này đã thuộc về một danh mục cha khác")
			}
			return nil, errors.New("không thể thêm danh mục vào danh mục cha")
		}

		// Invalidate cache cho parent
		s.invalidateCategoryCacheByID(*req.ParentID)
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

// GetParentCategories lấy danh sách tất cả parent categories (root categories) - dùng cho dropdown filter
func (s *CategoryService) GetParentCategories(includeInactive bool, language string) ([]dto.CategoryResponse, error) {
	// Lấy tất cả root categories (không có trong category_children như child)
	query := database.DB.Model(&models.Category{}).
		Where("NOT EXISTS (SELECT 1 FROM category_children WHERE category_children.child_id = categories.id AND category_children.deleted_at IS NULL)")

	if !includeInactive {
		query = query.Where("is_active = ?", true)
	}

	var categories []models.Category
	if err := query.Order("name ASC").Find(&categories).Error; err != nil {
		return nil, errors.New("không thể lấy danh sách danh mục cha")
	}

	// Transform với language nếu có
	if language == "en" || language == "vi" {
		for i := range categories {
			categories[i] = s.transformCategory(categories[i], language)
		}
	}

	// Load children cho từng parent
	categoryIDs := make([]uint, len(categories))
	for i, cat := range categories {
		categoryIDs[i] = cat.ID
	}

	// Lấy tất cả quan hệ parent-child
	var relations []models.CategoryChild
	if len(categoryIDs) > 0 {
		database.DB.Where("parent_id IN ?", categoryIDs).Find(&relations)
	}

	// Tạo map: parentID -> []childIDs
	childrenMap := make(map[uint][]uint)
	for _, rel := range relations {
		childrenMap[rel.ParentID] = append(childrenMap[rel.ParentID], rel.ChildID)
	}

	// Convert to response
	responses := make([]dto.CategoryResponse, len(categories))
	for i, cat := range categories {
		var childrenIDs []uint
		if cids, exists := childrenMap[cat.ID]; exists {
			childrenIDs = cids
		}

		responses[i] = dto.CategoryResponse{
			ID:            cat.ID,
			Name:          cat.Name,
			NameEn:        cat.NameEn,
			Description:   cat.Description,
			DescriptionEn: cat.DescriptionEn,
			Image:         cat.Image,
			IsActive:      cat.IsActive,
			ParentID:      nil, // Root category không có parent
			ChildrenIDs:   childrenIDs,
			CreatedAt:     cat.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:     cat.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	return responses, nil
}

// GetAllChildren lấy danh sách tất cả child categories - dùng cho dropdown filter
func (s *CategoryService) GetAllChildren(includeInactive bool, language string) ([]dto.CategoryResponse, error) {
	// Lấy tất cả child categories (có trong category_children như child)
	query := database.DB.Model(&models.Category{}).
		Where("EXISTS (SELECT 1 FROM category_children WHERE category_children.child_id = categories.id AND category_children.deleted_at IS NULL)")

	if !includeInactive {
		query = query.Where("is_active = ?", true)
	}

	var categories []models.Category
	if err := query.Order("name ASC").Find(&categories).Error; err != nil {
		return nil, fmt.Errorf("không thể lấy danh sách danh mục con: %v", err)
	}

	// Nếu không có child categories nào, trả về mảng rỗng
	if len(categories) == 0 {
		return []dto.CategoryResponse{}, nil
	}

	// Transform với language nếu có
	if language == "en" || language == "vi" {
		for i := range categories {
			categories[i] = s.transformCategory(categories[i], language)
		}
	}

	// Load parent và children cho từng child
	categoryIDs := make([]uint, len(categories))
	for i, cat := range categories {
		categoryIDs[i] = cat.ID
	}

	// Lấy tất cả quan hệ parent-child
	var relations []models.CategoryChild
	if err := database.DB.Where("parent_id IN ? OR child_id IN ?", categoryIDs, categoryIDs).Find(&relations).Error; err != nil {
		return nil, fmt.Errorf("không thể lấy quan hệ parent-child: %v", err)
	}

	// Tạo map: childID -> parentID và parentID -> []childIDs
	parentMap := make(map[uint]uint)     // childID -> parentID
	childrenMap := make(map[uint][]uint) // parentID -> []childIDs
	for _, rel := range relations {
		parentMap[rel.ChildID] = rel.ParentID
		childrenMap[rel.ParentID] = append(childrenMap[rel.ParentID], rel.ChildID)
	}

	// Convert to response
	responses := make([]dto.CategoryResponse, len(categories))
	for i, cat := range categories {
		var parentID *uint
		if pid, exists := parentMap[cat.ID]; exists {
			parentID = new(uint)
			*parentID = pid
		}

		var childrenIDs []uint
		if cids, exists := childrenMap[cat.ID]; exists {
			childrenIDs = cids
		}

		responses[i] = dto.CategoryResponse{
			ID:            cat.ID,
			Name:          cat.Name,
			NameEn:        cat.NameEn,
			Description:   cat.Description,
			DescriptionEn: cat.DescriptionEn,
			Image:         cat.Image,
			IsActive:      cat.IsActive,
			ParentID:      parentID,
			ChildrenIDs:   childrenIDs,
			CreatedAt:     cat.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:     cat.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	return responses, nil
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

	// Load parent và children relations
	var parentRelation models.CategoryChild
	if err := database.DB.Where("child_id = ?", id).First(&parentRelation).Error; err == nil {
		// Có parent, load parent vào relationship (nếu cần)
		category.ChildRelations = []models.CategoryChild{parentRelation}
	}

	var childrenRelations []models.CategoryChild
	database.DB.Where("parent_id = ?", id).Find(&childrenRelations)
	if len(childrenRelations) > 0 {
		category.ParentRelations = childrenRelations
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
	query := database.DB.Model(&models.Category{})

	// CHỈ LẤY ROOT CATEGORIES (parent categories)
	// Loại bỏ các categories có trong bảng category_children như child
	// Tức là chỉ lấy các categories KHÔNG có trong category_children.child_id
	// Sử dụng LEFT JOIN để xử lý trường hợp không có child nào
	query = query.Where("NOT EXISTS (SELECT 1 FROM category_children WHERE category_children.child_id = categories.id AND category_children.deleted_at IS NULL)")

	// Filter theo name (partial match, case-insensitive)
	if req.Name != nil && *req.Name != "" {
		normalizedSearch := s.normalizeForSearch(*req.Name)

		// Lấy tất cả ROOT categories để filter trong memory (vì PostgreSQL không có built-in function để xóa dấu tiếng Việt)
		// CHỈ lấy root categories (không có trong category_children như child)
		var allCategories []models.Category
		if err := database.DB.Select("id", "name", "name_en").
			Where("NOT EXISTS (SELECT 1 FROM category_children WHERE category_children.child_id = categories.id AND category_children.deleted_at IS NULL)").
			Find(&allCategories).Error; err != nil {
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

	// Load parent và children cho từng category
	categoryIDs := make([]uint, len(categories))
	for i, cat := range categories {
		categoryIDs[i] = cat.ID
	}

	// Lấy tất cả quan hệ parent-child
	var relations []models.CategoryChild
	if len(categoryIDs) > 0 {
		database.DB.Where("parent_id IN ? OR child_id IN ?", categoryIDs, categoryIDs).Find(&relations)
	}

	// Tạo map: childID -> parentID và parentID -> []childIDs
	parentMap := make(map[uint]uint)     // childID -> parentID
	childrenMap := make(map[uint][]uint) // parentID -> []childIDs
	for _, rel := range relations {
		parentMap[rel.ChildID] = rel.ParentID
		childrenMap[rel.ParentID] = append(childrenMap[rel.ParentID], rel.ChildID)
	}

	// Convert to response
	categoryResponses := make([]dto.CategoryResponse, len(categories))
	for i, cat := range categories {
		var parentID *uint
		if pid, exists := parentMap[cat.ID]; exists {
			pidCopy := pid // Copy value để tránh pointer đến loop variable
			parentID = &pidCopy
		}

		var childrenIDs []uint
		if cids, exists := childrenMap[cat.ID]; exists {
			childrenIDs = cids
		}

		categoryResponses[i] = dto.CategoryResponse{
			ID:            cat.ID,
			Name:          cat.Name,
			NameEn:        cat.NameEn,
			Description:   cat.Description,
			DescriptionEn: cat.DescriptionEn,
			Image:         cat.Image,
			IsActive:      cat.IsActive,
			ParentID:      parentID,
			ChildrenIDs:   childrenIDs,
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
	if err := database.DB.Where("id = ?", id).Preload("Products").First(&category).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("không tìm thấy danh mục với ID %d", id)
		}
		return nil, errors.New("không thể lấy danh mục")
	}

	// Nếu đổi tên, kiểm tra xem tên mới đã tồn tại chưa (chỉ kiểm tra các record chưa bị soft delete)
	if req.Name != nil && *req.Name != category.Name {
		var existingCategory models.Category
		if err := database.DB.Where("name = ? AND deleted_at IS NULL", *req.Name).First(&existingCategory).Error; err == nil {
			return nil, errors.New("danh mục với tên này đã tồn tại")
		}
	}

	// Cập nhật các fields (không thể đổi parent qua update, phải dùng API quản lý children)
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
		// Kiểm tra nếu đang inactive category
		if !*req.IsActive {
			// Kiểm tra nếu có products, không cho phép inactive
			productCount := len(category.Products)
			if productCount > 0 {
				return nil, fmt.Errorf("không thể vô hiệu hóa danh mục này vì có %d sản phẩm liên quan. Vui lòng xóa hoặc di chuyển các sản phẩm sang danh mục khác trước", productCount)
			}

			// Kiểm tra xem có category con nào đang active không (nếu là parent category)
			var activeChildrenCount int64
			var relations []models.CategoryChild
			if err := database.DB.Where("parent_id = ?", id).Find(&relations).Error; err == nil && len(relations) > 0 {
				// Lấy danh sách ID các children
				childIDs := make([]uint, len(relations))
				for i, rel := range relations {
					childIDs[i] = rel.ChildID
				}
				// Đếm số lượng children đang active
				if err := database.DB.Model(&models.Category{}).
					Where("id IN ? AND is_active = ?", childIDs, true).
					Count(&activeChildrenCount).Error; err != nil {
					return nil, errors.New("không thể kiểm tra trạng thái danh mục con")
				}
				// Nếu có ít nhất 1 child đang active, không cho phép inactive category cha
				if activeChildrenCount > 0 {
					return nil, fmt.Errorf("không thể vô hiệu hóa danh mục này vì có %d danh mục con đang hoạt động. Vui lòng vô hiệu hóa các danh mục con trước", activeChildrenCount)
				}
			}
		}
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

	// Kiểm tra nếu có category con nào đang active, không cho phép inactive category cha
	var activeChildrenCount int64
	var relations []models.CategoryChild
	if err := database.DB.Where("parent_id = ?", id).Find(&relations).Error; err == nil && len(relations) > 0 {
		// Lấy danh sách ID các children
		childIDs := make([]uint, len(relations))
		for i, rel := range relations {
			childIDs[i] = rel.ChildID
		}
		// Đếm số lượng children đang active
		if err := database.DB.Model(&models.Category{}).
			Where("id IN ? AND is_active = ?", childIDs, true).
			Count(&activeChildrenCount).Error; err != nil {
			return errors.New("không thể kiểm tra trạng thái danh mục con")
		}
		// Nếu có ít nhất 1 child đang active, không cho phép inactive category cha
		if activeChildrenCount > 0 {
			return fmt.Errorf("không thể vô hiệu hóa danh mục này vì có %d danh mục con đang hoạt động. Vui lòng vô hiệu hóa các danh mục con trước", activeChildrenCount)
		}
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

// isAncestorOf kiểm tra xem ancestorID có phải là ancestor của descendantID không
// bằng cách traverse lên parent chain từ descendantID
func (s *CategoryService) isAncestorOf(ancestorID, descendantID uint) bool {
	// Nếu cùng một ID thì không phải ancestor
	if ancestorID == descendantID {
		return false
	}

	// Traverse lên parent chain từ descendantID
	currentID := descendantID
	visited := make(map[uint]bool) // Để tránh infinite loop nếu có cycle (mặc dù không nên xảy ra)

	for {
		// Kiểm tra xem đã visit node này chưa (tránh infinite loop)
		if visited[currentID] {
			break
		}
		visited[currentID] = true

		// Tìm parent của currentID
		var relation models.CategoryChild
		if err := database.DB.Where("child_id = ?", currentID).First(&relation).Error; err != nil {
			// Không có parent nữa → đã đến root, không tìm thấy ancestor
			return false
		}

		// Nếu parent là ancestorID thì tìm thấy
		if relation.ParentID == ancestorID {
			return true
		}

		// Tiếp tục traverse lên parent
		currentID = relation.ParentID
	}

	return false
}

// AddChild thêm một category con vào category cha
func (s *CategoryService) AddChild(parentID, childID uint) error {
	// Kiểm tra parent và child có tồn tại không
	var parent, child models.Category
	if err := database.DB.Where("id = ?", parentID).First(&parent).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("không tìm thấy danh mục cha với ID %d", parentID)
		}
		return errors.New("không thể lấy danh mục cha")
	}

	if err := database.DB.Where("id = ?", childID).First(&child).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("không tìm thấy danh mục con với ID %d", childID)
		}
		return errors.New("không thể lấy danh mục con")
	}

	// Không cho phép category là parent của chính nó
	if parentID == childID {
		return errors.New("danh mục không thể là parent của chính nó")
	}

	// Kiểm tra child đã có parent chưa (một child chỉ có thể có một parent)
	var existingRelation models.CategoryChild
	if err := database.DB.Where("child_id = ?", childID).First(&existingRelation).Error; err == nil {
		return errors.New("danh mục con này đã thuộc về một danh mục cha khác")
	}

	// Kiểm tra circular reference: child không thể là parent hoặc ancestor của parent
	// Kiểm tra trực tiếp trước (A→B thì không cho B→A)
	var childAsParent models.CategoryChild
	if err := database.DB.Where("parent_id = ? AND child_id = ?", childID, parentID).First(&childAsParent).Error; err == nil {
		return errors.New("không thể tạo circular reference: danh mục con đã là parent của danh mục cha này")
	}

	// Kiểm tra transitive cycles: child không thể là ancestor của parent (A→B→C thì không cho C→A)
	if s.isAncestorOf(childID, parentID) {
		return errors.New("không thể tạo circular reference: danh mục con là ancestor của danh mục cha này")
	}

	// Tạo quan hệ parent-child
	categoryChild := models.CategoryChild{
		ParentID: parentID,
		ChildID:  childID,
	}

	if err := database.DB.Create(&categoryChild).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "UNIQUE constraint") {
			return errors.New("quan hệ parent-child này đã tồn tại")
		}
		return errors.New("không thể thêm danh mục con")
	}

	// Invalidate cache
	s.invalidateCategoryCache()
	s.invalidateCategoryCacheByID(parentID)
	s.invalidateCategoryCacheByID(childID)

	return nil
}

// RemoveChild xóa một category con khỏi category cha
func (s *CategoryService) RemoveChild(parentID, childID uint) error {
	// Kiểm tra quan hệ có tồn tại không
	var categoryChild models.CategoryChild
	if err := database.DB.Where("parent_id = ? AND child_id = ?", parentID, childID).First(&categoryChild).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("không tìm thấy quan hệ parent-child giữa danh mục %d và %d", parentID, childID)
		}
		return errors.New("không thể kiểm tra quan hệ parent-child")
	}

	// Xóa quan hệ
	if err := database.DB.Delete(&categoryChild).Error; err != nil {
		return errors.New("không thể xóa danh mục con")
	}

	// Invalidate cache
	s.invalidateCategoryCache()
	s.invalidateCategoryCacheByID(parentID)
	s.invalidateCategoryCacheByID(childID)

	return nil
}

// GetChildren lấy danh sách các category con của một category cha
func (s *CategoryService) GetChildren(parentID uint, language string) ([]models.Category, error) {
	var children []models.Category
	var relations []models.CategoryChild

	// Lấy tất cả quan hệ parent-child
	if err := database.DB.Where("parent_id = ?", parentID).Find(&relations).Error; err != nil {
		return nil, errors.New("không thể lấy danh sách danh mục con")
	}

	if len(relations) == 0 {
		return []models.Category{}, nil
	}

	// Lấy danh sách ID các children
	childIDs := make([]uint, len(relations))
	for i, rel := range relations {
		childIDs[i] = rel.ChildID
	}

	// Lấy thông tin các children
	if err := database.DB.Where("id IN ?", childIDs).Find(&children).Error; err != nil {
		return nil, errors.New("không thể lấy thông tin danh mục con")
	}

	// Transform với language nếu có
	if language == "en" || language == "vi" {
		for i := range children {
			children[i] = s.transformCategory(children[i], language)
		}
	}

	return children, nil
}

// GetParent lấy category cha của một category con
func (s *CategoryService) GetParent(childID uint) (*models.Category, error) {
	var relation models.CategoryChild
	if err := database.DB.Where("child_id = ?", childID).First(&relation).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Không có parent (root category)
		}
		return nil, errors.New("không thể lấy quan hệ parent-child")
	}

	var parent models.Category
	if err := database.DB.Where("id = ?", relation.ParentID).First(&parent).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("không tìm thấy danh mục cha với ID %d", relation.ParentID)
		}
		return nil, errors.New("không thể lấy danh mục cha")
	}

	return &parent, nil
}

// SearchChildren tìm kiếm và lọc category children
func (s *CategoryService) SearchChildren(req dto.SearchCategoryChildRequest, language string) (*dto.CategoryPaginationResponse, error) {
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

	// Bước 1: Lấy tất cả quan hệ parent-child từ bảng category_children
	var relations []models.CategoryChild
	relationQuery := database.DB.Model(&models.CategoryChild{})

	// Filter theo parentId nếu có
	if req.ParentID != nil {
		relationQuery = relationQuery.Where("parent_id = ?", *req.ParentID)
	}

	if err := relationQuery.Find(&relations).Error; err != nil {
		return nil, errors.New("không thể lấy danh sách quan hệ parent-child")
	}

	if len(relations) == 0 {
		// Không có quan hệ nào → trả về empty
		return &dto.CategoryPaginationResponse{
			Data:       []dto.CategoryResponse{},
			Total:      0,
			Page:       page,
			Limit:      limit,
			TotalPages: 0,
		}, nil
	}

	// Lấy danh sách ID các children
	childIDs := make([]uint, len(relations))
	for i, rel := range relations {
		childIDs[i] = rel.ChildID
	}

	// Bước 2: Query categories với các childIDs này
	query := database.DB.Model(&models.Category{}).Where("id IN ?", childIDs)

	// Filter theo name (partial match, case-insensitive)
	if req.Name != nil && *req.Name != "" {
		normalizedSearch := s.normalizeForSearch(*req.Name)

		// Lấy tất cả categories để filter trong memory
		var allCategories []models.Category
		if err := database.DB.Select("id", "name", "name_en").Where("id IN ?", childIDs).Find(&allCategories).Error; err != nil {
			return nil, errors.New("không thể tìm kiếm danh mục con")
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
			// Nếu có cả true và false (len == 2) → không filter, lấy tất cả
			if len(uniqueValues) == 1 {
				// Chỉ có 1 giá trị duy nhất → filter theo giá trị đó
				for value := range uniqueValues {
					query = query.Where("is_active = ?", value)
				}
			}
		case []bool:
			// Array of bools
			uniqueValues := make(map[bool]bool)
			for _, b := range v {
				uniqueValues[b] = true
			}
			// Nếu có cả true và false (len == 2) → không filter, lấy tất cả
			if len(uniqueValues) == 1 {
				// Chỉ có 1 giá trị duy nhất → filter theo giá trị đó
				for value := range uniqueValues {
					query = query.Where("is_active = ?", value)
				}
			}
		}
	}

	// Count total
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, errors.New("không thể đếm số lượng danh mục con")
	}

	// Sort - map field names to database columns
	sortByColumn := s.mapSortFieldToColumn(sortBy)
	orderBy := fmt.Sprintf("%s %s", sortByColumn, sortOrder)
	query = query.Order(orderBy)

	// Pagination
	offset := (page - 1) * limit
	var children []models.Category
	if err := query.Offset(offset).Limit(limit).Find(&children).Error; err != nil {
		return nil, errors.New("không thể lấy danh sách danh mục con")
	}

	// Transform với language nếu có
	if language == "en" || language == "vi" {
		for i := range children {
			children[i] = s.transformCategory(children[i], language)
		}
	}

	// Load parent và children cho từng category
	allChildIDs := make([]uint, len(children))
	for i, child := range children {
		allChildIDs[i] = child.ID
	}

	// Lấy tất cả quan hệ parent-child cho các children này
	var allRelations []models.CategoryChild
	if len(allChildIDs) > 0 {
		database.DB.Where("parent_id IN ? OR child_id IN ?", allChildIDs, allChildIDs).Find(&allRelations)
	}

	// Tạo map: childID -> parentID và parentID -> []childIDs
	parentMap := make(map[uint]uint)     // childID -> parentID
	childrenMap := make(map[uint][]uint) // parentID -> []childIDs
	for _, rel := range allRelations {
		parentMap[rel.ChildID] = rel.ParentID
		childrenMap[rel.ParentID] = append(childrenMap[rel.ParentID], rel.ChildID)
	}

	// Convert to response
	childrenResponses := make([]dto.CategoryResponse, len(children))
	for i, child := range children {
		var parentID *uint
		if pid, exists := parentMap[child.ID]; exists {
			pidCopy := pid // Copy value để tránh pointer đến loop variable
			parentID = &pidCopy
		}

		var childrenIDs []uint
		if cids, exists := childrenMap[child.ID]; exists {
			childrenIDs = cids
		}

		childrenResponses[i] = dto.CategoryResponse{
			ID:            child.ID,
			Name:          child.Name,
			NameEn:        child.NameEn,
			Description:   child.Description,
			DescriptionEn: child.DescriptionEn,
			Image:         child.Image,
			IsActive:      child.IsActive,
			ParentID:      parentID,
			ChildrenIDs:   childrenIDs,
			CreatedAt:     child.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:     child.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &dto.CategoryPaginationResponse{
		Data:       childrenResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

// ConvertCategoryToResponse convert category model sang response DTO với parent và children
func (s *CategoryService) ConvertCategoryToResponse(cat models.Category) dto.CategoryResponse {
	var parentID *uint
	var childrenIDs []uint

	// Load parent từ childRelations
	if len(cat.ChildRelations) > 0 {
		pidCopy := cat.ChildRelations[0].ParentID // Copy value để tránh pointer đến field trong struct value
		parentID = &pidCopy
	}

	// Load children từ parentRelations
	if len(cat.ParentRelations) > 0 {
		childrenIDs = make([]uint, len(cat.ParentRelations))
		for i, rel := range cat.ParentRelations {
			childrenIDs[i] = rel.ChildID
		}
	}

	return dto.CategoryResponse{
		ID:            cat.ID,
		Name:          cat.Name,
		NameEn:        cat.NameEn,
		Description:   cat.Description,
		DescriptionEn: cat.DescriptionEn,
		Image:         cat.Image,
		IsActive:      cat.IsActive,
		ParentID:      parentID,
		ChildrenIDs:   childrenIDs,
		CreatedAt:     cat.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:     cat.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
