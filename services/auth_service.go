package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"ecommerce-be/database"
	"ecommerce-be/dto"
	"ecommerce-be/models"
	"ecommerce-be/utils"

	"gorm.io/gorm"
)

type AuthService struct {
	otpService   *OtpService
	emailService *EmailService
}

func NewAuthService() *AuthService {
	return &AuthService{
		otpService:   NewOtpService(),
		emailService: NewEmailService(),
	}
}

// Login xử lý đăng nhập
func (s *AuthService) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {
	// Normalize email
	email := strings.ToLower(strings.TrimSpace(req.Email))

	// Tìm user theo email
	var user models.User
	if err := database.DB.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("email hoặc mật khẩu không đúng")
		}
		return nil, err
	}

	// Kiểm tra password
	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("email hoặc mật khẩu không đúng")
	}

	// Kiểm tra user có active không
	if !user.IsActive {
		return nil, errors.New("tài khoản đã bị vô hiệu hóa")
	}

	// Nếu email chưa được verify → gửi OTP và yêu cầu verify
	if !user.IsEmailVerified {
		// Kiểm tra xem OTP đã từng được set chưa
		hasOtpBeenSet := user.OTP != nil && user.OTPExpiresAt != nil

		// Kiểm tra OTP hiện tại có hết hạn không (chỉ khi đã có OTP)
		isOtpExpired := hasOtpBeenSet && time.Now().After(*user.OTPExpiresAt)

		// Nếu OTP hết hạn hoặc không có OTP → gửi OTP mới
		if isOtpExpired || !hasOtpBeenSet {
			// Generate OTP mới
			otp := s.otpService.GenerateOtp()
			otpExpiresAt := s.otpService.GetOtpExpiry()
			now := time.Now()

			// Lưu OTP vào database
			user.OTP = &otp
			user.OTPExpiresAt = &otpExpiresAt
			user.LastOTPSentAt = &now
			if err := database.DB.Save(&user).Error; err != nil {
				return nil, errors.New("không thể cập nhật OTP")
			}

			// Gửi OTP qua email
			if err := s.emailService.SendOtpEmail(user.Email, otp, user.Name); err != nil {
				return nil, fmt.Errorf("không thể gửi email OTP: %v", err)
			}

			// Thông báo khác nhau tùy vào trường hợp
			var message string
			if user.IsFirstLogin {
				// First login: phân biệt giữa "chưa có OTP" và "OTP đã hết hạn"
				if !hasOtpBeenSet {
					message = "Đây là lần đăng nhập đầu tiên. OTP đã được gửi đến email của bạn. Vui lòng xác thực email để tiếp tục."
				} else if isOtpExpired {
					message = "OTP trước đó đã hết hạn. OTP mới đã được gửi đến email của bạn. Vui lòng xác thực email để tiếp tục."
				} else {
					message = "Email chưa được xác thực. OTP mới đã được gửi đến email của bạn. Vui lòng xác thực email để tiếp tục."
				}
			} else {
				// User tự đăng ký: phân biệt giữa "chưa có OTP" và "OTP đã hết hạn"
				if !hasOtpBeenSet {
					message = "Email chưa được xác thực. OTP đã được gửi đến email của bạn. Vui lòng xác thực email để tiếp tục."
				} else if isOtpExpired {
					message = "OTP trước đó đã hết hạn. OTP mới đã được gửi đến email của bạn. Vui lòng xác thực email để tiếp tục."
				} else {
					message = "Email chưa được xác thực. OTP mới đã được gửi đến email của bạn. Vui lòng xác thực email để tiếp tục."
				}
			}

			return &dto.LoginResponse{
				Success:     true,
				Message:     message,
				RequiresOtp: true,
			}, nil
		}

		// Nếu OTP còn hạn → thông báo sử dụng OTP hiện tại
		return &dto.LoginResponse{
			Success:     true,
			Message:     "Email chưa được xác thực. OTP đã được gửi trước đó và vẫn còn hiệu lực. Vui lòng kiểm tra email và nhập mã OTP.",
			RequiresOtp: true,
		}, nil
	}

	// Tạo tokens
	accessToken, refreshToken, err := utils.GenerateTokens(user.ID, user.Email)
	if err != nil {
		return nil, errors.New("không thể tạo token")
	}

	// Lưu refresh token vào database
	user.RefreshToken = &refreshToken
	if err := database.DB.Save(&user).Error; err != nil {
		return nil, errors.New("không thể lưu refresh token")
	}

	// Tạo response
	response := &dto.LoginResponse{
		Success:      true,
		Message:      "Đăng nhập thành công!",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: &dto.UserResponse{
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

	return response, nil
}

// Register xử lý đăng ký
func (s *AuthService) Register(req dto.RegisterRequest) (*dto.RegisterResponse, error) {
	// Normalize email
	email := strings.ToLower(strings.TrimSpace(req.Email))

	// Kiểm tra email đã tồn tại chưa
	var existingUser models.User
	if err := database.DB.Where("email = ?", email).First(&existingUser).Error; err == nil {
		// User đã tồn tại
		if existingUser.IsEmailVerified {
			return nil, errors.New("email đã được sử dụng")
		}
		// User tồn tại nhưng chưa verify - cập nhật OTP
		otp := s.otpService.GenerateOtp()
		otpExpiresAt := s.otpService.GetOtpExpiry()
		now := time.Now()

		// Hash password mới nếu có thay đổi
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			return nil, errors.New("không thể hash password")
		}

		// Cập nhật thông tin user và OTP
		existingUser.Password = hashedPassword
		existingUser.Name = strings.TrimSpace(req.Name)
		existingUser.OTP = &otp
		existingUser.OTPExpiresAt = &otpExpiresAt
		existingUser.LastOTPSentAt = &now

		if err := database.DB.Save(&existingUser).Error; err != nil {
			return nil, errors.New("không thể cập nhật tài khoản")
		}

		// Gửi OTP qua email
		if err := s.emailService.SendOtpEmail(email, otp, existingUser.Name); err != nil {
			return nil, fmt.Errorf("không thể gửi email OTP: %v", err)
		}

		return &dto.RegisterResponse{
			Success: true,
			Message: "OTP đã được gửi đến email của bạn. Vui lòng kiểm tra và nhập mã OTP để hoàn tất đăng ký.",
			Email:   email,
		}, nil
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("không thể hash password")
	}

	// Generate OTP
	otp := s.otpService.GenerateOtp()
	otpExpiresAt := s.otpService.GetOtpExpiry()
	now := time.Now()

	// Tạo user mới (chưa verified)
	user := models.User{
		Email:           email,
		Password:        hashedPassword,
		Name:            strings.TrimSpace(req.Name),
		Role:            "customer",
		IsEmailVerified: false,
		IsActive:        false, // Chưa active cho đến khi verify
		IsFirstLogin:    false,
		OTP:             &otp,
		OTPExpiresAt:    &otpExpiresAt,
		LastOTPSentAt:   &now,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "UNIQUE constraint") {
			return nil, errors.New("email đã được sử dụng")
		}
		return nil, errors.New("không thể tạo tài khoản")
	}

	// Gửi OTP qua email
	if err := s.emailService.SendOtpEmail(email, otp, user.Name); err != nil {
		// Log lỗi nhưng không fail registration
		// Có thể retry sau hoặc thông báo cho user
		return nil, fmt.Errorf("không thể gửi email OTP: %v", err)
	}

	return &dto.RegisterResponse{
		Success: true,
		Message: "OTP đã được gửi đến email của bạn. Vui lòng kiểm tra và nhập mã OTP để hoàn tất đăng ký.",
		Email:   email,
	}, nil
}

// RefreshToken làm mới access token bằng refresh token
func (s *AuthService) RefreshToken(refreshToken string) (*dto.RefreshTokenResponse, error) {
	// Validate refresh token
	claims, err := utils.ValidateToken(refreshToken)
	if err != nil {
		return nil, errors.New("refresh token không hợp lệ hoặc đã hết hạn")
	}

	// Tìm user theo refresh token trong database
	var user models.User
	if err := database.DB.Where("refresh_token = ? AND id = ?", refreshToken, claims.UserID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("refresh token không hợp lệ")
		}
		return nil, err
	}

	// Kiểm tra user có active không
	if !user.IsActive {
		return nil, errors.New("tài khoản đã bị vô hiệu hóa")
	}

	// Tạo tokens mới
	accessToken, newRefreshToken, err := utils.GenerateTokens(user.ID, user.Email)
	if err != nil {
		return nil, errors.New("không thể tạo token")
	}

	// Lưu refresh token mới vào database
	user.RefreshToken = &newRefreshToken
	if err := database.DB.Save(&user).Error; err != nil {
		return nil, errors.New("không thể lưu refresh token")
	}

	return &dto.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		User: &dto.UserResponse{
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
	}, nil
}

// VerifyOtp xác thực OTP và kích hoạt tài khoản
func (s *AuthService) VerifyOtp(req dto.VerifyOtpRequest) (*dto.VerifyOtpResponse, error) {
	// Normalize email
	email := strings.ToLower(strings.TrimSpace(req.Email))

	// Tìm user
	var user models.User
	if err := database.DB.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("email không tồn tại")
		}
		return nil, err
	}

	// Kiểm tra email đã được verify chưa
	if user.IsEmailVerified {
		return nil, errors.New("email đã được xác thực")
	}

	// Kiểm tra OTP
	if !s.otpService.IsOtpValid(req.OTP, *user.OTP, user.OTPExpiresAt) {
		return nil, errors.New("mã OTP không hợp lệ hoặc đã hết hạn")
	}

	// Xác thực email và kích hoạt tài khoản
	user.IsEmailVerified = true
	user.IsActive = true
	user.IsFirstLogin = false
	user.OTP = nil
	user.OTPExpiresAt = nil

	if err := database.DB.Save(&user).Error; err != nil {
		return nil, errors.New("không thể cập nhật tài khoản")
	}

	// Tạo tokens
	accessToken, refreshToken, err := utils.GenerateTokens(user.ID, user.Email)
	if err != nil {
		return nil, errors.New("không thể tạo token")
	}

	// Lưu refresh token
	user.RefreshToken = &refreshToken
	if err := database.DB.Save(&user).Error; err != nil {
		return nil, errors.New("không thể lưu refresh token")
	}

	return &dto.VerifyOtpResponse{
		Success:      true,
		Message:      "Đăng ký thành công!",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: &dto.UserResponse{
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
	}, nil
}

// ResendOtp gửi lại OTP
func (s *AuthService) ResendOtp(req dto.ResendOtpRequest) (*dto.ResendOtpResponse, error) {
	// Normalize email
	email := strings.ToLower(strings.TrimSpace(req.Email))

	// Tìm user
	var user models.User
	if err := database.DB.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("email không tồn tại")
		}
		return nil, err
	}

	// Chỉ cho phép resend OTP nếu chưa verified
	if user.IsEmailVerified {
		return nil, errors.New("email đã được xác thực. Không cần gửi lại OTP")
	}

	// Rate limiting: Kiểm tra cooldown 60 giây
	if user.LastOTPSentAt != nil {
		now := time.Now()
		secondsSinceLastSent := int(now.Sub(*user.LastOTPSentAt).Seconds())
		cooldownSeconds := 60

		if secondsSinceLastSent < cooldownSeconds {
			remainingSeconds := cooldownSeconds - secondsSinceLastSent
			return nil, fmt.Errorf("vui lòng đợi %d giây trước khi yêu cầu gửi lại OTP", remainingSeconds)
		}
	}

	// Generate OTP mới
	otp := s.otpService.GenerateOtp()
	otpExpiresAt := s.otpService.GetOtpExpiry()
	now := time.Now()

	// Cập nhật OTP
	user.OTP = &otp
	user.OTPExpiresAt = &otpExpiresAt
	user.LastOTPSentAt = &now

	if err := database.DB.Save(&user).Error; err != nil {
		return nil, errors.New("không thể cập nhật OTP")
	}

	// Gửi OTP qua email
	if err := s.emailService.SendOtpEmail(email, otp, user.Name); err != nil {
		return nil, fmt.Errorf("không thể gửi email OTP: %v", err)
	}

	return &dto.ResendOtpResponse{
		Success: true,
		Message: "OTP mới đã được gửi đến email của bạn. Vui lòng kiểm tra và nhập mã OTP.",
		Email:   email,
	}, nil
}
