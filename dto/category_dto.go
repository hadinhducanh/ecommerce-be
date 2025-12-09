package dto

type CreateCategoryRequest struct {
	Name          string  `json:"name" binding:"required"`
	NameEn        *string `json:"nameEn"`
	Description   *string `json:"description"`
	DescriptionEn *string `json:"descriptionEn"`
	Image         *string `json:"image"`
	IsActive      *bool   `json:"isActive"`
	ParentID      *uint   `json:"parentId"` // Optional: Nếu có thì tự động thêm vào parent, nếu không thì tạo như root category
}

type UpdateCategoryRequest struct {
	Name          *string `json:"name"`
	NameEn        *string `json:"nameEn"`
	Description   *string `json:"description"`
	DescriptionEn *string `json:"descriptionEn"`
	Image         *string `json:"image"`
	IsActive      *bool   `json:"isActive"`
	// Không thể đổi parent qua update, phải dùng API quản lý children
}

type UpdateCategoryFullRequest struct {
	Name          string  `json:"name" binding:"required"`
	NameEn        *string `json:"nameEn"`
	Description   *string `json:"description"`
	DescriptionEn *string `json:"descriptionEn"`
	Image         *string `json:"image"`
	IsActive      *bool   `json:"isActive"`
}

type SearchCategoryRequest struct {
	Name      *string     `json:"name"`     // Search (partial match), không phải filter exact
	IsActive  interface{} `json:"isActive"` // *bool hoặc []bool - true = active, false = inactive, nil = all, [true, false] = all
	SortBy    *string     `json:"sortBy" binding:"omitempty,oneof=id name createdAt updatedAt"`
	SortOrder *string     `json:"sortOrder" binding:"omitempty,oneof=ASC DESC"`
	Page      *int        `json:"page" binding:"omitempty,min=1"`
	Limit     *int        `json:"limit" binding:"omitempty,min=1,max=1000"`
}

type CategoryResponse struct {
	ID            uint    `json:"id"`
	Name          string  `json:"name"`
	NameEn        *string `json:"nameEn"`
	Description   *string `json:"description"`
	DescriptionEn *string `json:"descriptionEn"`
	Image         *string `json:"image"`
	IsActive      bool    `json:"isActive"`
	ParentID      *uint   `json:"parentId,omitempty"`    // ID của category cha (null nếu là root) - tính từ childRelations
	ChildrenIDs   []uint  `json:"childrenIds,omitempty"` // Danh sách ID các category con
	CreatedAt     string  `json:"createdAt"`
	UpdatedAt     string  `json:"updatedAt"`
}

// DTOs cho quản lý children
type AddChildCategoryRequest struct {
	ChildID uint `json:"childId" binding:"required"` // ID của category con cần thêm
}

type RemoveChildCategoryRequest struct {
	ChildID uint `json:"childId" binding:"required"` // ID của category con cần xóa
}

type GetChildrenResponse struct {
	Success bool               `json:"success"`
	Message string             `json:"message"`
	Data    []CategoryResponse `json:"data"`
}

// SearchCategoryChildRequest DTO cho tìm kiếm category children
type SearchCategoryChildRequest struct {
	ParentID  *uint       `json:"parentId"` // Filter theo parent ID (exact match)
	Name      *string     `json:"name"`     // Search (partial match), không phải filter exact
	IsActive  interface{} `json:"isActive"` // *bool hoặc []bool - true = active, false = inactive, nil = all, [true, false] = all
	SortBy    *string     `json:"sortBy" binding:"omitempty,oneof=id name createdAt updatedAt"`
	SortOrder *string     `json:"sortOrder" binding:"omitempty,oneof=ASC DESC"`
	Page      *int        `json:"page" binding:"omitempty,min=1"`
	Limit     *int        `json:"limit" binding:"omitempty,min=1,max=1000"`
}

type SearchCategoryChildResponse struct {
	Success    bool               `json:"success"`
	Message    string             `json:"message"`
	Data       []CategoryResponse `json:"data"`
	Total      int64              `json:"total"`
	Page       int                `json:"page"`
	Limit      int                `json:"limit"`
	TotalPages int                `json:"totalPages"`
}

type CreateCategoryResponse struct {
	Success bool             `json:"success"`
	Message string           `json:"message"`
	Data    CategoryResponse `json:"data"`
}

type GetCategoryResponse struct {
	Success bool             `json:"success"`
	Message string           `json:"message"`
	Data    CategoryResponse `json:"data"`
}

type ListCategoryResponse struct {
	Success bool               `json:"success"`
	Data    []CategoryResponse `json:"data"`
}

type CategoryPaginationResponse struct {
	Data       []CategoryResponse `json:"data"`
	Total      int64              `json:"total"`
	Page       int                `json:"page"`
	Limit      int                `json:"limit"`
	TotalPages int                `json:"totalPages"`
}

type SearchCategoryResponse struct {
	Success    bool               `json:"success"`
	Message    string             `json:"message"`
	Data       []CategoryResponse `json:"data"`
	Total      int64              `json:"total"`
	Page       int                `json:"page"`
	Limit      int                `json:"limit"`
	TotalPages int                `json:"totalPages"`
}

type DeleteCategoryResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
