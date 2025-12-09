package handlers

import (
	"net/http"
	"strconv"

	"ecommerce-be/dto"
	"ecommerce-be/models"
	"ecommerce-be/services"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	categoryService   *services.CategoryService
	cloudinaryService *services.CloudinaryService
}

func NewCategoryHandler() (*CategoryHandler, error) {
	cloudinaryService, err := services.NewCloudinaryService()
	if err != nil {
		// Nếu Cloudinary không khởi tạo được, vẫn tạo handler nhưng không có cloudinary
		cloudinaryService = nil
	}

	return &CategoryHandler{
		categoryService:   services.NewCategoryService(),
		cloudinaryService: cloudinaryService,
	}, nil
}

// UploadImage upload image cho category (Chỉ admin)
func (h *CategoryHandler) UploadImage(c *gin.Context) {
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
		folder = "categories"
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
func (h *CategoryHandler) DeleteImage(c *gin.Context) {
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

// Search tìm kiếm categories (Public)
// Gộp GET all vào POST search - nếu body null/empty thì hiển thị tất cả
func (h *CategoryHandler) Search(c *gin.Context) {
	language := c.DefaultQuery("language", "vi")

	var req dto.SearchCategoryRequest
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

	result, err := h.categoryService.Search(req, language)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	response := dto.SearchCategoryResponse{
		Success:    true,
		Message:    "Tìm kiếm danh mục thành công",
		Data:       result.Data,
		Total:      result.Total,
		Page:       result.Page,
		Limit:      result.Limit,
		TotalPages: result.TotalPages,
	}

	c.JSON(http.StatusOK, response)
}

// Create tạo category mới (Chỉ admin)
func (h *CategoryHandler) Create(c *gin.Context) {
	var req dto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Dữ liệu không hợp lệ",
			"details": err.Error(),
		})
		return
	}

	category, err := h.categoryService.Create(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// Load relations để có parent và children
	categoryWithRelations, err := h.categoryService.FindOne(category.ID, true, "vi")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Không thể tải thông tin danh mục sau khi tạo",
		})
		return
	}

	response := dto.CreateCategoryResponse{
		Success: true,
		Message: "Tạo danh mục thành công",
		Data:    h.categoryService.ConvertCategoryToResponse(*categoryWithRelations),
	}

	c.JSON(http.StatusCreated, response)
}

// FindOne lấy một category theo ID (Public)
func (h *CategoryHandler) FindOne(c *gin.Context) {
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
	// Mặc định cho phép view category inactive (để có thể kích hoạt lại)
	// Chỉ filter inactive nếu truyền includeInactive=false
	includeInactive := c.Query("includeInactive") != "false"

	category, err := h.categoryService.FindOne(uint(id), includeInactive, language)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	response := dto.GetCategoryResponse{
		Success: true,
		Message: "Lấy thông tin danh mục thành công",
		Data:    h.categoryService.ConvertCategoryToResponse(*category),
	}

	c.JSON(http.StatusOK, response)
}

// Update cập nhật một phần category (Chỉ admin - Partial Update)
func (h *CategoryHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "ID không hợp lệ",
		})
		return
	}

	var req dto.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Dữ liệu không hợp lệ",
			"details": err.Error(),
		})
		return
	}

	category, err := h.categoryService.Update(uint(id), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// Load relations để có parent và children
	categoryWithRelations, err := h.categoryService.FindOne(category.ID, true, "vi")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Không thể tải thông tin danh mục sau khi cập nhật",
		})
		return
	}

	response := dto.GetCategoryResponse{
		Success: true,
		Message: "Cập nhật danh mục thành công",
		Data:    h.categoryService.ConvertCategoryToResponse(*categoryWithRelations),
	}

	c.JSON(http.StatusOK, response)
}

// Replace thay thế toàn bộ category (Chỉ admin - Full Replacement)
func (h *CategoryHandler) Replace(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "ID không hợp lệ",
		})
		return
	}

	var req dto.UpdateCategoryFullRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Dữ liệu không hợp lệ",
			"details": err.Error(),
		})
		return
	}

	// Convert UpdateCategoryFullRequest to UpdateCategoryRequest
	updateReq := dto.UpdateCategoryRequest{
		Name:          &req.Name,
		NameEn:        req.NameEn,
		Description:   req.Description,
		DescriptionEn: req.DescriptionEn,
		Image:         req.Image,
		IsActive:      req.IsActive,
	}

	category, err := h.categoryService.Update(uint(id), updateReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// Load relations để có parent và children
	categoryWithRelations, err := h.categoryService.FindOne(category.ID, true, "vi")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Không thể tải thông tin danh mục sau khi thay thế",
		})
		return
	}

	response := dto.GetCategoryResponse{
		Success: true,
		Message: "Cập nhật danh mục thành công",
		Data:    h.categoryService.ConvertCategoryToResponse(*categoryWithRelations),
	}

	c.JSON(http.StatusOK, response)
}

// Remove xóa category (Chỉ admin - soft delete)
func (h *CategoryHandler) Remove(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "ID không hợp lệ",
		})
		return
	}

	if err := h.categoryService.Remove(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	response := dto.DeleteCategoryResponse{
		Success: true,
		Message: "Danh mục đã được xóa thành công",
	}

	c.JSON(http.StatusOK, response)
}

// HardDelete xóa vĩnh viễn category (Chỉ admin)
func (h *CategoryHandler) HardDelete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "ID không hợp lệ",
		})
		return
	}

	if err := h.categoryService.HardDelete(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	response := dto.DeleteCategoryResponse{
		Success: true,
		Message: "Danh mục đã được xóa thành công",
	}

	c.JSON(http.StatusOK, response)
}

// SearchChildren tìm kiếm category children (Public)
// Gộp GET all vào POST search - nếu body null/empty thì hiển thị tất cả children
func (h *CategoryHandler) SearchChildren(c *gin.Context) {
	language := c.DefaultQuery("language", "vi")

	var req dto.SearchCategoryChildRequest
	// Cho phép body null/empty - nếu không có body thì vẫn OK (sẽ lấy tất cả children)
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

	result, err := h.categoryService.SearchChildren(req, language)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	response := dto.SearchCategoryChildResponse{
		Success:    true,
		Message:    "Tìm kiếm danh mục con thành công",
		Data:       result.Data,
		Total:      result.Total,
		Page:       result.Page,
		Limit:      result.Limit,
		TotalPages: result.TotalPages,
	}

	c.JSON(http.StatusOK, response)
}

// GetParentCategories lấy danh sách tất cả parent categories (dùng cho dropdown filter) (Public)
func (h *CategoryHandler) GetParentCategories(c *gin.Context) {
	language := c.DefaultQuery("language", "vi")
	// Thống nhất behavior: mặc định include inactive (để có thể view và kích hoạt lại)
	// Chỉ filter inactive nếu truyền includeInactive=false
	includeInactive := c.Query("includeInactive") != "false"

	parents, err := h.categoryService.GetParentCategories(includeInactive, language)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.ListCategoryResponse{
		Success: true,
		Data:    parents,
	})
}

// GetAllChildren lấy danh sách tất cả child categories (dùng cho dropdown filter) (Public)
func (h *CategoryHandler) GetAllChildren(c *gin.Context) {
	language := c.DefaultQuery("language", "vi")
	// Thống nhất behavior: mặc định include inactive (để có thể view và kích hoạt lại)
	// Chỉ filter inactive nếu truyền includeInactive=false
	includeInactive := c.Query("includeInactive") != "false"

	children, err := h.categoryService.GetAllChildren(includeInactive, language)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.ListCategoryResponse{
		Success: true,
		Data:    children,
	})
}

// AddChild thêm một category con vào category cha (Admin only)
func (h *CategoryHandler) AddChild(c *gin.Context) {
	idStr := c.Param("id")
	parentID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "ID danh mục cha không hợp lệ",
		})
		return
	}

	var req dto.AddChildCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Dữ liệu không hợp lệ",
			"details": err.Error(),
		})
		return
	}

	if err := h.categoryService.AddChild(uint(parentID), req.ChildID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Đã thêm danh mục con thành công",
	})
}

// RemoveChild xóa một category con khỏi category cha (Admin only)
func (h *CategoryHandler) RemoveChild(c *gin.Context) {
	idStr := c.Param("id")
	parentID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "ID danh mục cha không hợp lệ",
		})
		return
	}

	var req dto.RemoveChildCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Dữ liệu không hợp lệ",
			"details": err.Error(),
		})
		return
	}

	if err := h.categoryService.RemoveChild(uint(parentID), req.ChildID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Đã xóa danh mục con thành công",
	})
}

// GetChildren lấy danh sách các category con của một category cha (Public)
func (h *CategoryHandler) GetChildren(c *gin.Context) {
	idStr := c.Param("id")
	parentID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "ID danh mục cha không hợp lệ",
		})
		return
	}

	language := c.DefaultQuery("language", "vi")
	includeInactive := c.Query("includeInactive") != "false"

	children, err := h.categoryService.GetChildren(uint(parentID), language)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// Filter inactive nếu cần
	if !includeInactive {
		var activeChildren []models.Category
		for _, child := range children {
			if child.IsActive {
				activeChildren = append(activeChildren, child)
			}
		}
		children = activeChildren
	}

	// Convert to response
	childrenResponses := make([]dto.CategoryResponse, len(children))
	for i, child := range children {
		// Load parent và children của child này để có đầy đủ thông tin
		childWithRelations, err := h.categoryService.FindOne(child.ID, includeInactive, language)
		if err != nil {
			// Nếu không load được, sử dụng child hiện tại (không có relations)
			childrenResponses[i] = h.categoryService.ConvertCategoryToResponse(child)
			continue
		}
		childrenResponses[i] = h.categoryService.ConvertCategoryToResponse(*childWithRelations)
	}

	response := dto.GetChildrenResponse{
		Success: true,
		Message: "Lấy danh sách danh mục con thành công",
		Data:    childrenResponses,
	}

	c.JSON(http.StatusOK, response)
}
