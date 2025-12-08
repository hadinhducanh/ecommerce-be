package services

import (
	"errors"
	"fmt"
	"math"
	"strings"

	"ecommerce-be/database"
	"ecommerce-be/dto"
	"ecommerce-be/models"
	"ecommerce-be/utils"

	"gorm.io/gorm"
)

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

// CreateUserByAdmin tạo user bởi admin
func (s *UserService) CreateUserByAdmin(req dto.CreateUserByAdminRequest) (*models.User, error) {
	// Normalize email
	email := strings.ToLower(strings.TrimSpace(req.Email))

	// Kiểm tra email đã tồn tại chưa
	var existingUser models.User
	if err := database.DB.Where("email = ?", email).First(&existingUser).Error; err == nil {
		return nil, errors.New("email đã được sử dụng")
	}

	// Kiểm tra số điện thoại đã tồn tại chưa (nếu có)
	if req.Phone != nil && *req.Phone != "" {
		var existingUserByPhone models.User
		if err := database.DB.Where("phone = ?", *req.Phone).First(&existingUserByPhone).Error; err == nil {
			return nil, errors.New("số điện thoại đã được sử dụng")
		}
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("không thể hash password")
	}

	// Xác định role
	role := "customer"
	if req.Role != nil {
		role = *req.Role
	}

	// Tạo user
	user := models.User{
		Email:           email,
		Password:        hashedPassword,
		Name:            strings.TrimSpace(req.Name),
		Role:            role,
		Phone:           req.Phone,
		Avatar:          req.Avatar,
		IsEmailVerified: false,
		IsActive:        true, // Admin tạo → active ngay
		IsFirstLogin:    true, // Admin tạo → đây là first login
	}

	if err := database.DB.Create(&user).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "UNIQUE constraint") {
			return nil, errors.New("email hoặc số điện thoại đã được sử dụng")
		}
		return nil, errors.New("không thể tạo tài khoản")
	}

	return &user, nil
}

// GetProfile lấy thông tin profile của user
func (s *UserService) GetProfile(userID uint) (*models.User, error) {
	var user models.User
	if err := database.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("không tìm thấy người dùng")
		}
		return nil, err
	}
	return &user, nil
}

// UpdateProfile cập nhật profile của user
func (s *UserService) UpdateProfile(userID uint, req dto.UpdateProfileRequest) (*models.User, error) {
	var user models.User
	if err := database.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("không tìm thấy người dùng")
		}
		return nil, err
	}

	// Cập nhật các fields
	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Phone != nil {
		user.Phone = req.Phone
	}
	if req.Avatar != nil {
		user.Avatar = req.Avatar
	}
	if req.Address != nil {
		user.Address = req.Address
	}
	if req.Gender != nil {
		user.Gender = req.Gender
	}

	if err := database.DB.Save(&user).Error; err != nil {
		return nil, errors.New("không thể cập nhật thông tin")
	}

	return &user, nil
}

// UpdateUserByAdmin admin cập nhật thông tin user
func (s *UserService) UpdateUserByAdmin(userID uint, req dto.UpdateUserByAdminRequest, currentUserID uint) (*models.User, error) {
	var user models.User
	if err := database.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("không tìm thấy người dùng")
		}
		return nil, err
	}

	// Bảo vệ: Ngăn admin tự vô hiệu hóa hoặc đổi role của chính mình
	if currentUserID == userID {
		if req.IsActive != nil && !*req.IsActive {
			return nil, errors.New("bạn không thể vô hiệu hóa tài khoản của chính mình. Vui lòng nhờ admin khác thực hiện")
		}
		if req.Role != nil && *req.Role != user.Role {
			return nil, errors.New("bạn không thể thay đổi role của chính mình. Vui lòng nhờ admin khác thực hiện")
		}
	}

	// Kiểm tra số điện thoại đã tồn tại chưa (nếu có thay đổi)
	if req.Phone != nil {
		phoneChanged := false
		if user.Phone == nil {
			phoneChanged = *req.Phone != ""
		} else {
			phoneChanged = *req.Phone != *user.Phone
		}

		if phoneChanged && *req.Phone != "" {
			var existingUserByPhone models.User
			if err := database.DB.Where("phone = ? AND id != ?", *req.Phone, userID).First(&existingUserByPhone).Error; err == nil {
				return nil, errors.New("số điện thoại đã được sử dụng")
			}
		}
	}

	// Cập nhật các fields
	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Phone != nil {
		user.Phone = req.Phone
	}
	if req.Avatar != nil {
		user.Avatar = req.Avatar
	}
	if req.Address != nil {
		user.Address = req.Address
	}
	if req.Gender != nil {
		user.Gender = req.Gender
	}
	if req.Role != nil {
		user.Role = *req.Role
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}
	if req.IsEmailVerified != nil {
		user.IsEmailVerified = *req.IsEmailVerified
	}

	if err := database.DB.Save(&user).Error; err != nil {
		return nil, errors.New("không thể cập nhật thông tin")
	}

	return &user, nil
}

// ChangePassword đổi mật khẩu
func (s *UserService) ChangePassword(userID uint, req dto.ChangePasswordRequest) error {
	var user models.User
	if err := database.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("không tìm thấy người dùng")
		}
		return err
	}

	// Kiểm tra mật khẩu cũ có đúng không
	if !utils.CheckPassword(req.OldPassword, user.Password) {
		return errors.New("mật khẩu cũ không đúng")
	}

	// Kiểm tra mật khẩu mới có khác mật khẩu cũ không
	if utils.CheckPassword(req.NewPassword, user.Password) {
		return errors.New("mật khẩu mới phải khác mật khẩu cũ")
	}

	// Hash mật khẩu mới và lưu
	hashedNewPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return errors.New("không thể hash password")
	}

	user.Password = hashedNewPassword
	if err := database.DB.Save(&user).Error; err != nil {
		return errors.New("không thể cập nhật mật khẩu")
	}

	return nil
}

// Search tìm kiếm và lọc users
func (s *UserService) Search(req dto.SearchUserRequest) (*dto.PaginationResponse, error) {
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

	// Build query - chỉ select các fields cần thiết (không lấy password, refreshToken, otp)
	query := database.DB.Model(&models.User{}).
		Select("id", "email", "name", "role", "phone", "avatar", "address", "gender",
			"is_active", "is_email_verified", "is_first_login", "created_at", "updated_at")

	// Search theo email (partial match, case-insensitive) - không phải filter exact
	if req.Email != nil && *req.Email != "" {
		query = query.Where("LOWER(email) LIKE LOWER(?)", "%"+*req.Email+"%")
	}

	// Search theo name (partial match, case-insensitive) - không phải filter exact
	if req.Name != nil && *req.Name != "" {
		query = query.Where("LOWER(name) LIKE LOWER(?)", "%"+*req.Name+"%")
	}

	// Filter theo role (hỗ trợ string hoặc array)
	if req.Role != nil {
		switch v := req.Role.(type) {
		case string:
			// String đơn
			query = query.Where("role = ?", v)
		case []interface{}:
			// Array - extract strings
			var roles []string
			for _, item := range v {
				if str, ok := item.(string); ok {
					roles = append(roles, str)
				}
			}
			if len(roles) > 0 {
				query = query.Where("role IN ?", roles)
			}
		case []string:
			// Array of strings
			if len(v) > 0 {
				query = query.Where("role IN ?", v)
			}
		}
	}

	// Filter theo isActive (hỗ trợ boolean hoặc array)
	if req.IsActive != nil {
		switch v := req.IsActive.(type) {
		case bool:
			// Boolean đơn
			query = query.Where("is_active = ?", v)
		case []interface{}:
			// Array - extract booleans
			var values []bool
			for _, item := range v {
				if b, ok := item.(bool); ok {
					values = append(values, b)
				}
			}
			if len(values) > 0 {
				query = query.Where("is_active IN ?", values)
			}
		case []bool:
			// Array of bools
			if len(v) > 0 {
				query = query.Where("is_active IN ?", v)
			}
		}
	}

	// Filter theo isEmailVerified (hỗ trợ boolean hoặc array)
	if req.IsEmailVerified != nil {
		switch v := req.IsEmailVerified.(type) {
		case bool:
			// Boolean đơn
			query = query.Where("is_email_verified = ?", v)
		case []interface{}:
			// Array - extract booleans
			var values []bool
			for _, item := range v {
				if b, ok := item.(bool); ok {
					values = append(values, b)
				}
			}
			if len(values) > 0 {
				query = query.Where("is_email_verified IN ?", values)
			}
		case []bool:
			// Array of bools
			if len(v) > 0 {
				query = query.Where("is_email_verified IN ?", v)
			}
		}
	}

	// Count total
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, errors.New("không thể đếm số lượng users")
	}

	// Sort - map field names to database columns
	sortByColumn := s.mapSortFieldToColumn(sortBy)
	orderBy := fmt.Sprintf("%s %s", sortByColumn, sortOrder)
	query = query.Order(orderBy)

	// Pagination
	offset := (page - 1) * limit
	var users []models.User
	if err := query.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, errors.New("không thể lấy danh sách users")
	}

	// Convert to response
	userResponses := make([]dto.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = dto.UserResponse{
			ID:              user.ID,
			Email:           user.Email,
			Name:            user.Name,
			Role:            user.Role,
			Phone:           user.Phone,
			Avatar:          user.Avatar,
			Address:         user.Address,
			Gender:          user.Gender,
			IsEmailVerified: user.IsEmailVerified,
			IsActive:        user.IsActive,
			IsFirstLogin:    user.IsFirstLogin,
			CreatedAt:       user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:       user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &dto.PaginationResponse{
		Data:       userResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

// mapSortFieldToColumn map sort field name to database column name
func (s *UserService) mapSortFieldToColumn(field string) string {
	fieldMap := map[string]string{
		"name":      "name",
		"email":     "email",
		"createdAt": "created_at",
		"updatedAt": "updated_at",
	}
	if column, ok := fieldMap[field]; ok {
		return column
	}
	return "created_at" // default
}
