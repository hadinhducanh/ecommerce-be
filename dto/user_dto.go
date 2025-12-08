package dto

type CreateUserByAdminRequest struct {
	Email    string  `json:"email" binding:"required,email"`
	Password string  `json:"password" binding:"required,min=6"`
	Name     string  `json:"name" binding:"required"`
	Phone    *string `json:"phone"`
	Role     *string `json:"role" binding:"omitempty,oneof=customer admin"`
	Avatar   *string `json:"avatar"`
}

type UpdateProfileRequest struct {
	Name    *string `json:"name"`
	Phone   *string `json:"phone"`
	Avatar  *string `json:"avatar"`
	Address *string `json:"address"`
	Gender  *string `json:"gender" binding:"omitempty,oneof=male female other"`
}

type UpdateUserByAdminRequest struct {
	Name            *string `json:"name"`
	Phone           *string `json:"phone"`
	Avatar          *string `json:"avatar"`
	Address         *string `json:"address"`
	Gender          *string `json:"gender" binding:"omitempty,oneof=male female other"`
	Role            *string `json:"role" binding:"omitempty,oneof=customer admin"`
	IsActive        *bool   `json:"isActive"`
	IsEmailVerified *bool   `json:"isEmailVerified"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=6"`
}

type SearchUserRequest struct {
	Email          *string   `json:"email"` // Search (partial match), không phải filter exact
	Name           *string   `json:"name"` // Search (partial match), không phải filter exact
	Role           interface{} `json:"role"` // string hoặc []string - hỗ trợ filter nhiều role
	IsActive       interface{} `json:"isActive"` // bool hoặc []bool - hỗ trợ filter nhiều trạng thái
	IsEmailVerified interface{} `json:"isEmailVerified"` // bool hoặc []bool - hỗ trợ filter nhiều trạng thái
	SortBy         *string   `json:"sortBy" binding:"omitempty,oneof=name email createdAt updatedAt"`
	SortOrder      *string   `json:"sortOrder" binding:"omitempty,oneof=ASC DESC"`
	Page           *int      `json:"page" binding:"omitempty,min=1"`
	Limit          *int      `json:"limit" binding:"omitempty,min=1,max=1000"`
}

type PaginationResponse struct {
	Data       []UserResponse `json:"data"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	TotalPages int            `json:"totalPages"`
}

type CreateUserResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	User    UserResponse `json:"user"`
}

type GetProfileResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Data    UserResponse `json:"data"`
}

type SearchUserResponse struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Data    []UserResponse    `json:"data"`
	Total   int64             `json:"total"`
	Page    int               `json:"page"`
	Limit   int               `json:"limit"`
	TotalPages int            `json:"totalPages"`
}

type ChangePasswordResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

