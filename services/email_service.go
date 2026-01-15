package services

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/gomail.v2"
)

type EmailService struct {
	smtpHost     string
	smtpPort     int
	smtpUser     string
	smtpPassword string
	fromEmail    string
}

func NewEmailService() *EmailService {
	return &EmailService{
		smtpHost:     getEnv("SMTP_HOST", "smtp.gmail.com"),
		smtpPort:     getEnvInt("SMTP_PORT", 587),
		smtpUser:     getEnv("SMTP_USER", ""),
		smtpPassword: getEnv("SMTP_PASS", ""),
		fromEmail:    getEnv("SMTP_USER", "noreply@ecommerce.com"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	var port int
	fmt.Sscanf(value, "%d", &port)
	if port == 0 {
		return defaultValue
	}
	return port
}

// SendOtpEmail gửi OTP qua email
func (s *EmailService) SendOtpEmail(email, otp, name string) error {
	// Log OTP ra terminal để test
	fmt.Printf("\n📧 ===== OTP EMAIL =====\n")
	fmt.Printf("📨 Gửi đến: %s\n", email)
	fmt.Printf("👤 Người nhận: %s\n", name)
	fmt.Printf("🔑 Mã OTP: %s\n", otp)
	fmt.Printf("⏰ Thời gian: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Printf("⏳ Hết hạn sau: 5 phút\n")
	fmt.Printf("========================\n\n")

	// Tạo email message
	m := gomail.NewMessage()
	m.SetHeader("From", s.fromEmail)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Mã OTP xác thực đăng ký tài khoản")

	// HTML template
	htmlBody := fmt.Sprintf(`
		<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
			<h2 style="color: #333;">Xác thực đăng ký tài khoản</h2>
			<p>Xin chào <strong>%s</strong>,</p>
			<p>Cảm ơn bạn đã đăng ký tài khoản. Vui lòng sử dụng mã OTP sau để xác thực:</p>
			<div style="background-color: #f4f4f4; padding: 20px; text-align: center; margin: 20px 0;">
				<h1 style="color: #007bff; font-size: 32px; margin: 0;">%s</h1>
			</div>
			<p>Mã OTP này sẽ hết hạn sau <strong>5 phút</strong>.</p>
			<p>Nếu bạn không yêu cầu mã này, vui lòng bỏ qua email này.</p>
			<hr style="border: none; border-top: 1px solid #eee; margin: 20px 0;">
			<p style="color: #666; font-size: 12px;">Đây là email tự động, vui lòng không trả lời.</p>
		</div>
	`, name, otp)

	m.SetBody("text/html", htmlBody)

	// Tạo dialer
	d := gomail.NewDialer(s.smtpHost, s.smtpPort, s.smtpUser, s.smtpPassword)

	// Gửi email
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("không thể gửi email: %w", err)
	}

	return nil
}

