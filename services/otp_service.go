package services

import (
	"fmt"
	"math/rand"
	"time"
)

type OtpService struct{}

func NewOtpService() *OtpService {
	return &OtpService{}
}

// GenerateOtp tạo mã OTP 6 chữ số
func (s *OtpService) GenerateOtp() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%06d", r.Intn(1000000))
}

// GetOtpExpiry trả về thời gian hết hạn OTP (5 phút)
func (s *OtpService) GetOtpExpiry() time.Time {
	return time.Now().Add(5 * time.Minute)
}

// IsOtpValid kiểm tra OTP có hợp lệ không
func (s *OtpService) IsOtpValid(otp, storedOtp string, expiresAt *time.Time) bool {
	if otp == "" || storedOtp == "" || expiresAt == nil {
		return false
	}

	// Kiểm tra OTP có khớp không
	if otp != storedOtp {
		return false
	}

	// Kiểm tra OTP còn hạn không
	if time.Now().After(*expiresAt) {
		return false
	}

	return true
}

