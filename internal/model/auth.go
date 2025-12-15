package model

type Auth struct {
	UserID         string
	Email          string
	OrganizationID string
}

// RegisterRequest represents user registration request
type RegisterRequest struct {
	Name             string `json:"name" validate:"required,max=100"`
	Email            string `json:"email" validate:"required,email,max=255"`
	Password         string `json:"password" validate:"required,min=8,max=100"`
	OrganizationName string `json:"organization_name" validate:"required,max=200"`
}

// RegisterResponse represents user registration response
type RegisterResponse struct {
	User         UserResponse         `json:"user"`
	Organization OrganizationResponse `json:"organization"`
}

// LoginRequest represents user login request
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse represents user login response with JWT tokens
type LoginResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	ExpiresIn    int          `json:"expires_in"`
	TokenType    string       `json:"token_type"`
	User         UserResponse `json:"user"`
}

// RefreshTokenRequest represents refresh token request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// RefreshTokenResponse represents refresh token response
type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

// VerifyUserRequest represents verify user request (for middleware)
type VerifyUserRequest struct {
	Token string `json:"token" validate:"required"`
}

// VerifyEmailRequest represents email verification request
type VerifyEmailRequest struct {
	Token string `json:"token" validate:"required"`
}

// VerifyEmailResponse represents email verification response
type VerifyEmailResponse struct {
	Message string `json:"message"`
}

// ResendVerificationRequest represents resend verification email request
type ResendVerificationRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// ResendVerificationResponse represents resend verification email response
type ResendVerificationResponse struct {
	Message string `json:"message"`
}
