package handlers

import (
	"net/http"

	"ecommerce-be/dto"
	"ecommerce-be/services"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		authService: services.NewAuthService(),
	}
}

// Login xử lý request đăng nhập
// @Summary Đăng nhập
// @Description Đăng nhập với email và password
// @Tags auth
// @Accept json
// @Produce json
// @Param login body dto.LoginRequest true "Thông tin đăng nhập"
// @Success 200 {object} dto.LoginResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest

	// Validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Dữ liệu không hợp lệ",
			"details": err.Error(),
		})
		return
	}

	// Gọi service
	response, err := h.authService.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// Register xử lý request đăng ký
// @Summary Đăng ký
// @Description Đăng ký tài khoản mới
// @Tags auth
// @Accept json
// @Produce json
// @Param register body dto.RegisterRequest true "Thông tin đăng ký"
// @Success 201 {object} dto.RegisterResponse
// @Failure 400 {object} map[string]string
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest

	// Validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Dữ liệu không hợp lệ",
			"details": err.Error(),
		})
		return
	}

	// Gọi service
	response, err := h.authService.Register(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// RefreshToken xử lý request làm mới token
// @Summary Làm mới token
// @Description Làm mới access token bằng refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param refresh body dto.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} dto.RefreshTokenResponse
// @Failure 401 {object} map[string]string
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest

	// Validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Dữ liệu không hợp lệ",
			"details": err.Error(),
		})
		return
	}

	// Gọi service
	response, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// VerifyOtp xử lý request xác thực OTP
func (h *AuthHandler) VerifyOtp(c *gin.Context) {
	var req dto.VerifyOtpRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Dữ liệu không hợp lệ",
			"details": err.Error(),
		})
		return
	}

	response, err := h.authService.VerifyOtp(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// ResendOtp xử lý request gửi lại OTP
func (h *AuthHandler) ResendOtp(c *gin.Context) {
	var req dto.ResendOtpRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Dữ liệu không hợp lệ",
			"details": err.Error(),
		})
		return
	}

	response, err := h.authService.ResendOtp(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}
