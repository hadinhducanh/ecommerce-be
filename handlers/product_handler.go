package handlers

import (
	"net/http"
	"strconv"

	"ecommerce-be/dto"
	"ecommerce-be/services"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productService   *services.ProductService
	cloudinaryService *services.CloudinaryService
}

func NewProductHandler() (*ProductHandler, error) {
	cloudinaryService, err := services.NewCloudinaryService()
	if err != nil {
		cloudinaryService = nil
	}

	return &ProductHandler{
		productService:   services.NewProductService(),
		cloudinaryService: cloudinaryService,
	}, nil
}

// UploadImage upload image cho product (Chỉ admin)
func (h *ProductHandler) UploadImage(c *gin.Context) {
	if h.cloudinaryService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Cloudinary service không khả dụng",
		})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "không có file được upload",
		})
		return
	}

	folder := c.PostForm("folder")
	if folder == "" {
		folder = "products"
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "không thể mở file",
		})
		return
	}
	defer src.Close()

	uploadResult, err := h.cloudinaryService.UploadImageWithResponse(
		src,
		file.Size,
		file.Filename,
		file.Header.Get("Content-Type"),
		folder,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, uploadResult)
}

// DeleteImage xóa image từ Cloudinary (Chỉ admin)
func (h *ProductHandler) DeleteImage(c *gin.Context) {
	if h.cloudinaryService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Cloudinary service không khả dụng",
		})
		return
	}

	var req dto.DeleteImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Dữ liệu không hợp lệ",
			"details": err.Error(),
		})
		return
	}

	deleteResult, err := h.cloudinaryService.DeleteImageWithResponse(req.URL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, deleteResult)
}

// Search tìm kiếm products (Public)
// Gộp GET all vào POST search - nếu body null/empty thì hiển thị tất cả
func (h *ProductHandler) Search(c *gin.Context) {
	language := c.DefaultQuery("language", "vi")

	var req dto.SearchProductRequest
	// Cho phép body null/empty - nếu không có body thì vẫn OK (sẽ lấy tất cả)
	if c.Request.ContentLength > 0 {
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "Dữ liệu không hợp lệ",
				"details": err.Error(),
			})
			return
		}
	}

	result, err := h.productService.Search(req, language)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	response := dto.SearchProductResponse{
		Success:    true,
		Message:    "Tìm kiếm sản phẩm thành công",
		Data:       result.Data,
		Total:      result.Total,
		Page:       result.Page,
		Limit:      result.Limit,
		TotalPages: result.TotalPages,
	}

	c.JSON(http.StatusOK, response)
}

// Create tạo product mới (Chỉ admin)
func (h *ProductHandler) Create(c *gin.Context) {
	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Dữ liệu không hợp lệ",
			"details": err.Error(),
		})
		return
	}

	product, err := h.productService.Create(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	var categoryResp *dto.CategoryResponse
	if product.Category.ID > 0 {
		categoryResp = &dto.CategoryResponse{
			ID:            product.Category.ID,
			Name:          product.Category.Name,
			NameEn:        product.Category.NameEn,
			Description:   product.Category.Description,
			DescriptionEn: product.Category.DescriptionEn,
			Image:         product.Category.Image,
			IsActive:      product.Category.IsActive,
			CreatedAt:     product.Category.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:     product.Category.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	response := dto.CreateProductResponse{
		Success: true,
		Message: "Tạo sản phẩm thành công",
		Data: dto.ProductResponse{
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
			Category:      categoryResp,
			CreatedAt:     product.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:     product.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
	}

	c.JSON(http.StatusCreated, response)
}

// FindOne lấy một product theo ID (Public)
func (h *ProductHandler) FindOne(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "ID không hợp lệ",
		})
		return
	}

	language := c.DefaultQuery("language", "vi")
	includeInactive := c.Query("includeInactive") == "true"

	product, err := h.productService.FindOne(uint(id), includeInactive, language)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	var categoryResp *dto.CategoryResponse
	if product.Category.ID > 0 {
		categoryResp = &dto.CategoryResponse{
			ID:            product.Category.ID,
			Name:          product.Category.Name,
			NameEn:        product.Category.NameEn,
			Description:   product.Category.Description,
			DescriptionEn: product.Category.DescriptionEn,
			Image:         product.Category.Image,
			IsActive:      product.Category.IsActive,
			CreatedAt:     product.Category.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:     product.Category.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	response := dto.GetProductResponse{
		Success: true,
		Message: "Lấy thông tin sản phẩm thành công",
		Data: dto.ProductResponse{
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
			Category:      categoryResp,
			CreatedAt:     product.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:     product.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
	}

	c.JSON(http.StatusOK, response)
}

// Update cập nhật một phần product (Chỉ admin - Partial Update)
func (h *ProductHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "ID không hợp lệ",
		})
		return
	}

	var req dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Dữ liệu không hợp lệ",
			"details": err.Error(),
		})
		return
	}

	product, err := h.productService.Update(uint(id), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	var categoryResp *dto.CategoryResponse
	if product.Category.ID > 0 {
		categoryResp = &dto.CategoryResponse{
			ID:            product.Category.ID,
			Name:          product.Category.Name,
			NameEn:        product.Category.NameEn,
			Description:   product.Category.Description,
			DescriptionEn: product.Category.DescriptionEn,
			Image:         product.Category.Image,
			IsActive:      product.Category.IsActive,
			CreatedAt:     product.Category.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:     product.Category.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	response := dto.GetProductResponse{
		Success: true,
		Message: "Cập nhật sản phẩm thành công",
		Data: dto.ProductResponse{
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
			Category:      categoryResp,
			CreatedAt:     product.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:     product.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
	}

	c.JSON(http.StatusOK, response)
}

// Replace thay thế toàn bộ product (Chỉ admin - Full Replacement)
func (h *ProductHandler) Replace(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "ID không hợp lệ",
		})
		return
	}

	var req dto.UpdateProductFullRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Dữ liệu không hợp lệ",
			"details": err.Error(),
		})
		return
	}

	product, err := h.productService.Update(uint(id), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	var categoryResp *dto.CategoryResponse
	if product.Category.ID > 0 {
		categoryResp = &dto.CategoryResponse{
			ID:            product.Category.ID,
			Name:          product.Category.Name,
			NameEn:        product.Category.NameEn,
			Description:   product.Category.Description,
			DescriptionEn: product.Category.DescriptionEn,
			Image:         product.Category.Image,
			IsActive:      product.Category.IsActive,
			CreatedAt:     product.Category.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:     product.Category.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	response := dto.GetProductResponse{
		Success: true,
		Message: "Cập nhật sản phẩm thành công",
		Data: dto.ProductResponse{
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
			Category:      categoryResp,
			CreatedAt:     product.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:     product.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
	}

	c.JSON(http.StatusOK, response)
}

// Remove xóa product (Chỉ admin - soft delete)
func (h *ProductHandler) Remove(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "ID không hợp lệ",
		})
		return
	}

	if err := h.productService.Remove(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	response := dto.DeleteProductResponse{
		Success: true,
		Message: "Sản phẩm đã được xóa thành công",
	}

	c.JSON(http.StatusOK, response)
}

// HardDelete xóa vĩnh viễn product (Chỉ admin)
func (h *ProductHandler) HardDelete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "ID không hợp lệ",
		})
		return
	}

	if err := h.productService.HardDelete(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	response := dto.DeleteProductResponse{
		Success: true,
		Message: "Sản phẩm đã được xóa thành công",
	}

	c.JSON(http.StatusOK, response)
}

// SearchSuggestions lấy danh sách gợi ý tìm kiếm dựa trên query (Public)
// @Summary Lấy gợi ý tìm kiếm
// @Description Lấy danh sách gợi ý tìm kiếm dựa trên query string
// @Tags products
// @Accept json
// @Produce json
// @Param query query string true "Từ khóa tìm kiếm"
// @Param language query string false "Ngôn ngữ (vi/en)" default(vi)
// @Param limit query int false "Số lượng suggestions" default(10)
// @Success 200 {object} dto.SearchSuggestionsResponse
// @Failure 400 {object} map[string]string
// @Router /api/v1/products/search-suggestions [get]
func (h *ProductHandler) SearchSuggestions(c *gin.Context) {
	query := c.Query("query")
	language := c.DefaultQuery("language", "vi")
	limitStr := c.DefaultQuery("limit", "10")

	limit := 10
	if l, err := strconv.Atoi(limitStr); err == nil {
		limit = l
	}

	suggestions, err := h.productService.GetSearchSuggestions(query, language, limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	response := dto.SearchSuggestionsResponse{
		Success: true,
		Data:    suggestions,
		Total:   len(suggestions),
	}

	c.JSON(http.StatusOK, response)
}

// PopularSearches lấy danh sách từ khóa tìm kiếm phổ biến (Public)
// @Summary Lấy từ khóa tìm kiếm phổ biến
// @Description Lấy danh sách từ khóa tìm kiếm phổ biến nhất
// @Tags products
// @Accept json
// @Produce json
// @Param language query string false "Ngôn ngữ (vi/en)" default(vi)
// @Param limit query int false "Số lượng searches" default(10)
// @Success 200 {object} dto.PopularSearchesResponse
// @Failure 400 {object} map[string]string
// @Router /api/v1/products/popular-searches [get]
func (h *ProductHandler) PopularSearches(c *gin.Context) {
	language := c.DefaultQuery("language", "vi")
	limitStr := c.DefaultQuery("limit", "10")

	limit := 10
	if l, err := strconv.Atoi(limitStr); err == nil {
		limit = l
	}

	searches, err := h.productService.GetPopularSearches(language, limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	response := dto.PopularSearchesResponse{
		Success: true,
		Data:    searches,
		Total:   len(searches),
	}

	c.JSON(http.StatusOK, response)
}
