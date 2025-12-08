package dto

type CreateCategoryRequest struct {
	Name          string  `json:"name" binding:"required"`
	NameEn        *string `json:"nameEn"`
	Description   *string `json:"description"`
	DescriptionEn *string `json:"descriptionEn"`
	Image         *string `json:"image"`
	IsActive      *bool   `json:"isActive"`
}

type UpdateCategoryRequest struct {
	Name          *string `json:"name"`
	NameEn        *string `json:"nameEn"`
	Description   *string `json:"description"`
	DescriptionEn *string `json:"descriptionEn"`
	Image         *string `json:"image"`
	IsActive      *bool   `json:"isActive"`
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
	Name      *string `json:"name"`
	IsActive  interface{} `json:"isActive"` // *bool hoặc []bool - true = active, false = inactive, nil = all, [true, false] = all
	SortBy    *string `json:"sortBy" binding:"omitempty,oneof=id name createdAt updatedAt"`
	SortOrder *string `json:"sortOrder" binding:"omitempty,oneof=ASC DESC"`
	Page      *int    `json:"page" binding:"omitempty,min=1"`
	Limit     *int    `json:"limit" binding:"omitempty,min=1,max=1000"`
}

type CategoryResponse struct {
	ID            uint    `json:"id"`
	Name          string  `json:"name"`
	NameEn        *string `json:"nameEn"`
	Description   *string `json:"description"`
	DescriptionEn *string `json:"descriptionEn"`
	Image         *string `json:"image"`
	IsActive      bool    `json:"isActive"`
	CreatedAt     string  `json:"createdAt"`
	UpdatedAt     string  `json:"updatedAt"`
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
