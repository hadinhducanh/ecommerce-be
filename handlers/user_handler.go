package handlers

import (
	"net/http"
	"strconv"

	"ecommerce-be/dto"
	"ecommerce-be/services"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService      *services.UserService
	cloudinaryService *services.CloudinaryService
}

func NewUserHandler() (*UserHandler, error) {
	cloudinaryService, err := services.NewCloudinaryService()
	if err != nil {
		// Nếu Cloudinary không khởi tạo được, vẫn tạo handler nhưng không có cloudinary
		cloudinaryService = nil
	}

	return &UserHandler{
		userService:      services.NewUserService(),
		cloudinaryService: cloudinaryService,
	}, nil
}

// CreateUserByAdmin tạo user bởi admin
func (h *UserHandler) CreateUserByAdmin(c *gin.Context) {
	var req dto.CreateUserByAdminRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Dữ liệu không hợp lệ",
			"details": err.Error(),
		})
		return
	}

	user, err := h.userService.CreateUserByAdmin(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	response := dto.CreateUserResponse{
		Success: true,
		Message: "Tạo tài khoản thành công. Người dùng sẽ nhận OTP khi đăng nhập lần đầu.",
		User: dto.UserResponse{
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
		},
	}

	c.JSON(http.StatusCreated, response)
}

// Search tìm kiếm users
// Search tìm kiếm users (Admin only)
// Gộp GET all vào POST search - nếu body null/empty thì hiển thị tất cả
func (h *UserHandler) Search(c *gin.Context) {
	var req dto.SearchUserRequest
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

	result, err := h.userService.Search(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	response := dto.SearchUserResponse{
		Success:    true,
		Message:    "Tìm kiếm người dùng thành công",
		Data:       result.Data,
		Total:      result.Total,
		Page:       result.Page,
		Limit:      result.Limit,
		TotalPages: result.TotalPages,
	}

	c.JSON(http.StatusOK, response)
}

// GetProfile xem profile của chính mình
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, _ := c.Get("userID")
	userIDUint := userID.(uint)

	user, err := h.userService.GetProfile(userIDUint)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	response := dto.GetProfileResponse{
		Success: true,
		Message: "Lấy thông tin cá nhân thành công",
		Data: dto.UserResponse{
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
		},
	}

	c.JSON(http.StatusOK, response)
}

// UpdateProfile cập nhật profile của chính mình
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, _ := c.Get("userID")
	userIDUint := userID.(uint)

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Dữ liệu không hợp lệ",
			"details": err.Error(),
		})
		return
	}

	user, err := h.userService.UpdateProfile(userIDUint, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	response := dto.GetProfileResponse{
		Success: true,
		Message: "Cập nhật thông tin cá nhân thành công",
		Data: dto.UserResponse{
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
		},
	}

	c.JSON(http.StatusOK, response)
}

// ChangePassword đổi mật khẩu
func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID, _ := c.Get("userID")
	userIDUint := userID.(uint)

	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Dữ liệu không hợp lệ",
			"details": err.Error(),
		})
		return
	}

	if err := h.userService.ChangePassword(userIDUint, req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	response := dto.ChangePasswordResponse{
		Success: true,
		Message: "Đổi mật khẩu thành công",
	}

	c.JSON(http.StatusOK, response)
}

// UploadAvatar upload avatar cho chính mình
func (h *UserHandler) UploadAvatar(c *gin.Context) {
	if h.cloudinaryService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Cloudinary service không khả dụng",
		})
		return
	}

	userID, _ := c.Get("userID")
	userIDUint := userID.(uint)

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
		folder = "users"
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

	// Tự động cập nhật avatar URL
	updateReq := dto.UpdateProfileRequest{
		Avatar: &uploadResult.Data.URL,
	}
	_, err = h.userService.UpdateProfile(userIDUint, updateReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "không thể cập nhật avatar",
		})
		return
	}

	c.JSON(http.StatusOK, uploadResult)
}

// DeleteAvatar xóa avatar
func (h *UserHandler) DeleteAvatar(c *gin.Context) {
	if h.cloudinaryService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Cloudinary service không khả dụng",
		})
		return
	}

	userID, _ := c.Get("userID")
	userIDUint := userID.(uint)

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

	// Tự động cập nhật avatar = null
	nullStr := ""
	updateReq := dto.UpdateProfileRequest{
		Avatar: &nullStr,
	}
	_, err = h.userService.UpdateProfile(userIDUint, updateReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "không thể cập nhật avatar",
		})
		return
	}

	c.JSON(http.StatusOK, deleteResult)
}

// GetUserById admin xem user bất kỳ
func (h *UserHandler) GetUserById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "ID không hợp lệ",
		})
		return
	}

	user, err := h.userService.GetProfile(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	response := dto.GetProfileResponse{
		Success: true,
		Message: "Lấy thông tin người dùng thành công",
		Data: dto.UserResponse{
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
		},
	}

	c.JSON(http.StatusOK, response)
}

// UpdateUserByAdmin admin cập nhật user bất kỳ
func (h *UserHandler) UpdateUserByAdmin(c *gin.Context) {
	currentUserID, _ := c.Get("userID")
	currentUserIDUint := currentUserID.(uint)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "ID không hợp lệ",
		})
		return
	}

	var req dto.UpdateUserByAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Dữ liệu không hợp lệ",
			"details": err.Error(),
		})
		return
	}

	user, err := h.userService.UpdateUserByAdmin(uint(id), req, currentUserIDUint)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	response := dto.GetProfileResponse{
		Success: true,
		Message: "Cập nhật thông tin người dùng thành công",
		Data: dto.UserResponse{
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
		},
	}

	c.JSON(http.StatusOK, response)
}

// UploadAvatarByAdmin admin upload avatar cho user bất kỳ
func (h *UserHandler) UploadAvatarByAdmin(c *gin.Context) {
	if h.cloudinaryService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Cloudinary service không khả dụng",
		})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "ID không hợp lệ",
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
		folder = "users"
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

	// Tự động cập nhật avatar URL
	updateReq := dto.UpdateUserByAdminRequest{
		Avatar: &uploadResult.Data.URL,
	}
	currentUserID, _ := c.Get("userID")
	_, err = h.userService.UpdateUserByAdmin(uint(id), updateReq, currentUserID.(uint))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "không thể cập nhật avatar",
		})
		return
	}

	c.JSON(http.StatusOK, uploadResult)
}

// ChangePasswordByAdmin admin đổi mật khẩu cho user
func (h *UserHandler) ChangePasswordByAdmin(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "ID không hợp lệ",
		})
		return
	}

	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Dữ liệu không hợp lệ",
			"details": err.Error(),
		})
		return
	}

	if err := h.userService.ChangePassword(uint(id), req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	response := dto.ChangePasswordResponse{
		Success: true,
		Message: "Đổi mật khẩu người dùng thành công",
	}

	c.JSON(http.StatusOK, response)
}

