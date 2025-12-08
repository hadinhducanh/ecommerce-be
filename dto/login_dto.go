package dto

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
}

type RegisterResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Email   string `json:"email"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RefreshTokenResponse struct {
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
	User         *UserResponse `json:"user"`
}

type VerifyOtpRequest struct {
	Email string `json:"email" binding:"required,email"`
	OTP   string `json:"otp" binding:"required,len=6"`
}

type VerifyOtpResponse struct {
	Success      bool          `json:"success"`
	Message      string        `json:"message"`
	AccessToken  string        `json:"access_token,omitempty"`
	RefreshToken string        `json:"refresh_token,omitempty"`
	User         *UserResponse `json:"user,omitempty"`
}

type ResendOtpRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResendOtpResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Email   string `json:"email"`
}

type LoginResponse struct {
	Success      bool          `json:"success"`
	Message      string        `json:"message"`
	AccessToken  string        `json:"access_token,omitempty"`
	RefreshToken string        `json:"refresh_token,omitempty"`
	RequiresOtp  bool          `json:"requiresOtp,omitempty"`
	User         *UserResponse `json:"user,omitempty"`
}

type UserResponse struct {
	ID              uint    `json:"id"`
	Email           string  `json:"email"`
	Name            string  `json:"name"`
	Role            string  `json:"role"`
	Phone           *string `json:"phone"`
	Avatar          *string `json:"avatar"`
	Address         *string `json:"address"`
	Gender          *string `json:"gender"`
	IsEmailVerified bool    `json:"isEmailVerified"`
	IsActive        bool    `json:"isActive"`
	IsFirstLogin    bool    `json:"isFirstLogin"`
	CreatedAt       string  `json:"createdAt"`
	UpdatedAt       string  `json:"updatedAt"`
}
