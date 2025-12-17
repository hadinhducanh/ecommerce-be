package dto

type CreateProductRequest struct {
	Name          string   `json:"name" binding:"required"`
	NameEn        *string  `json:"nameEn"`
	Description   *string  `json:"description"`
	DescriptionEn *string  `json:"descriptionEn"`
	Price         float64  `json:"price" binding:"required,min=0"`
	Stock         int      `json:"stock" binding:"min=0"`
	Image         *string  `json:"image"`
	Images        []string `json:"images"`
	CategoryID    uint     `json:"categoryId" binding:"required"`
	SKU           *string  `json:"sku"`
	IsActive      *bool    `json:"isActive"`
}

type UpdateProductRequest struct {
	Name          *string  `json:"name"`
	NameEn        *string  `json:"nameEn"`
	Description   *string  `json:"description"`
	DescriptionEn *string  `json:"descriptionEn"`
	Price          *float64 `json:"price" binding:"omitempty,min=0"`
	Stock          *int     `json:"stock" binding:"omitempty,min=0"`
	Image          *string  `json:"image"`
	Images         []string `json:"images"`
	CategoryID     *uint    `json:"categoryId"`
	SKU            *string  `json:"sku"`
	IsActive       *bool    `json:"isActive"`
}

type UpdateProductFullRequest struct {
	Name          string   `json:"name" binding:"required"`
	NameEn        *string  `json:"nameEn"`
	Description   *string  `json:"description"`
	DescriptionEn *string  `json:"descriptionEn"`
	Price         float64  `json:"price" binding:"required,min=0"`
	Stock         int      `json:"stock" binding:"min=0"`
	Image         *string  `json:"image"`
	Images        []string `json:"images"`
	CategoryID    uint     `json:"categoryId" binding:"required"`
	SKU           *string  `json:"sku"`
	IsActive      *bool    `json:"isActive"`
}

type SearchProductRequest struct {
	Name           *string     `json:"name"` // Search (partial match), không phải filter exact
	CategoryID      *uint       `json:"categoryId"` // Filter (exact match) - ưu tiên nếu có cả categoryId và parentCategoryId
	ParentCategoryID *uint      `json:"parentCategoryId"` // Filter theo danh mục cha - lấy tất cả sản phẩm của các danh mục con
	IsActive       interface{} `json:"isActive"` // *bool hoặc []bool - true = active, false = inactive, nil = all, [true, false] = all
	MinPrice       *float64    `json:"minPrice" binding:"omitempty,min=0"` // Filter (>=)
	MaxPrice       *float64    `json:"maxPrice" binding:"omitempty,min=0"` // Filter (<=)
	InStock        *bool       `json:"inStock"` // true = chỉ lấy sản phẩm còn hàng (stock > 0)
	SortBy         *string     `json:"sortBy" binding:"omitempty,oneof=id name price stock createdAt updatedAt"`
	SortOrder      *string     `json:"sortOrder" binding:"omitempty,oneof=ASC DESC"`
	Page           *int        `json:"page" binding:"omitempty,min=1"`
	Limit          *int        `json:"limit" binding:"omitempty,min=1,max=1000"`
}

type ProductResponse struct {
	ID            uint     `json:"id"`
	Name          string   `json:"name"`
	NameEn        *string  `json:"nameEn"`
	Description   *string  `json:"description"`
	DescriptionEn *string  `json:"descriptionEn"`
	Price         float64  `json:"price"`
	Stock         int      `json:"stock"`
	Image         *string  `json:"image"`
	Images        []string `json:"images,omitempty"`
	Sold          int      `json:"sold"`
	Rating        float64  `json:"rating"`
	ReviewCount   int      `json:"reviewCount"`
	IsActive      bool     `json:"isActive"`
	SKU           *string  `json:"sku"`
	CategoryID    uint     `json:"categoryId"`
	Category      *CategoryResponse `json:"category,omitempty"`
	CreatedAt     string   `json:"createdAt"`
	UpdatedAt     string   `json:"updatedAt"`
}

type CreateProductResponse struct {
	Success bool           `json:"success"`
	Message string         `json:"message"`
	Data    ProductResponse `json:"data"`
}

type GetProductResponse struct {
	Success bool           `json:"success"`
	Message string         `json:"message"`
	Data    ProductResponse `json:"data"`
}

type ProductPaginationResponse struct {
	Data       []ProductResponse `json:"data"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
	TotalPages int               `json:"totalPages"`
}

type SearchProductResponse struct {
	Success    bool                   `json:"success"`
	Message    string                 `json:"message"`
	Data       []ProductResponse      `json:"data"`
	Total      int64                  `json:"total"`
	Page       int                    `json:"page"`
	Limit      int                    `json:"limit"`
	TotalPages int                    `json:"totalPages"`
}

type DeleteProductResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

