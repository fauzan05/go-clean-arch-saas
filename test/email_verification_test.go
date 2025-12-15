package test

import (
	"encoding/json"
	"go-clean-arch-saas/internal/entity"
	"go-clean-arch-saas/internal/model"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmailVerification_RegisterUserEmailNotVerified(t *testing.T) {
	CleanupDatabase(t)
	CreateTestPlan(t, "free", "Free Plan", 0)

	// Register user
	registerBody := `{
		"name": "Test User",
		"email": "test@example.com",
		"password": "password123",
		"organization_name": "Test Org"
	}`

	resp, err := MakeRequest("POST", "/api/v1/auth/register", registerBody, "")
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var registerResponse model.WebResponse[model.RegisterResponse]
	json.NewDecoder(resp.Body).Decode(&registerResponse)

	// Check that email_verified is false
	assert.False(t, registerResponse.Data.User.EmailVerified)
}

func TestEmailVerification_VerifyEmail_Success(t *testing.T) {
	CleanupDatabase(t)
	CreateTestPlan(t, "free", "Free Plan", 0)

	// Register user
	registerBody := `{
		"name": "Test User",
		"email": "test@example.com",
		"password": "password123",
		"organization_name": "Test Org"
	}`

	resp, err := MakeRequest("POST", "/api/v1/auth/register", registerBody, "")
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var registerResponse model.WebResponse[model.RegisterResponse]
	json.NewDecoder(resp.Body).Decode(&registerResponse)

	// Get verification token from database
	var user entity.User
	err = db.Where("email = ?", "test@example.com").First(&user).Error
	assert.Nil(t, err)
	assert.NotNil(t, user.VerificationToken)
	assert.False(t, user.EmailVerified)

	// Verify email
	verifyBody := `{
		"token": "` + *user.VerificationToken + `"
	}`

	resp, err = MakeRequest("POST", "/api/v1/auth/verify-email", verifyBody, "")
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var verifyResponse model.WebResponse[model.VerifyEmailResponse]
	json.NewDecoder(resp.Body).Decode(&verifyResponse)
	assert.Equal(t, "Email verified successfully", verifyResponse.Data.Message)

	// Check database - email should be verified
	err = db.Where("email = ?", "test@example.com").First(&user).Error
	assert.Nil(t, err)
	assert.True(t, user.EmailVerified)
	assert.NotNil(t, user.EmailVerifiedAt)
	assert.Nil(t, user.VerificationToken) // Token should be cleared
}

func TestEmailVerification_VerifyEmail_InvalidToken(t *testing.T) {
	CleanupDatabase(t)

	verifyBody := `{
		"token": "invalid-token-12345"
	}`

	resp, err := MakeRequest("POST", "/api/v1/auth/verify-email", verifyBody, "")
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestEmailVerification_VerifyEmail_AlreadyVerified(t *testing.T) {
	CleanupDatabase(t)
	CreateTestPlan(t, "free", "Free Plan", 0)

	// Register user
	registerBody := `{
		"name": "Test User",
		"email": "test@example.com",
		"password": "password123",
		"organization_name": "Test Org"
	}`

	resp, err := MakeRequest("POST", "/api/v1/auth/register", registerBody, "")
	assert.Nil(t, err)

	// Get verification token
	var user entity.User
	err = db.Where("email = ?", "test@example.com").First(&user).Error
	assert.Nil(t, err)

	verifyBody := `{
		"token": "` + *user.VerificationToken + `"
	}`

	// Verify first time
	resp, err = MakeRequest("POST", "/api/v1/auth/verify-email", verifyBody, "")
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Try to verify again with same token (should fail because token is cleared)
	resp, err = MakeRequest("POST", "/api/v1/auth/verify-email", verifyBody, "")
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestEmailVerification_ResendVerification_Success(t *testing.T) {
	CleanupDatabase(t)
	CreateTestPlan(t, "free", "Free Plan", 0)

	// Register user
	registerBody := `{
		"name": "Test User",
		"email": "test@example.com",
		"password": "password123",
		"organization_name": "Test Org"
	}`

	resp, err := MakeRequest("POST", "/api/v1/auth/register", registerBody, "")
	assert.Nil(t, err)

	// Get old verification token
	var user entity.User
	err = db.Where("email = ?", "test@example.com").First(&user).Error
	assert.Nil(t, err)
	oldToken := *user.VerificationToken

	// Resend verification
	resendBody := `{
		"email": "test@example.com"
	}`

	resp, err = MakeRequest("POST", "/api/v1/auth/resend-verification", resendBody, "")
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var resendResponse model.WebResponse[model.ResendVerificationResponse]
	json.NewDecoder(resp.Body).Decode(&resendResponse)
	assert.Equal(t, "If the email exists, a verification link has been sent", resendResponse.Data.Message)

	// Check that token changed
	err = db.Where("email = ?", "test@example.com").First(&user).Error
	assert.Nil(t, err)
	assert.NotEqual(t, oldToken, *user.VerificationToken)
}

func TestEmailVerification_ResendVerification_EmailNotExists(t *testing.T) {
	CleanupDatabase(t)

	resendBody := `{
		"email": "nonexistent@example.com"
	}`

	resp, err := MakeRequest("POST", "/api/v1/auth/resend-verification", resendBody, "")
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Should return success for security (don't reveal if email exists)
	var resendResponse model.WebResponse[model.ResendVerificationResponse]
	json.NewDecoder(resp.Body).Decode(&resendResponse)
	assert.Equal(t, "If the email exists, a verification link has been sent", resendResponse.Data.Message)
}

func TestEmailVerification_ResendVerification_AlreadyVerified(t *testing.T) {
	CleanupDatabase(t)
	CreateTestPlan(t, "free", "Free Plan", 0)

	// Register and verify user
	registerBody := `{
		"name": "Test User",
		"email": "test@example.com",
		"password": "password123",
		"organization_name": "Test Org"
	}`

	MakeRequest("POST", "/api/v1/auth/register", registerBody, "")

	// Verify email
	var user entity.User
	db.Where("email = ?", "test@example.com").First(&user)
	verifyBody := `{"token": "` + *user.VerificationToken + `"}`
	MakeRequest("POST", "/api/v1/auth/verify-email", verifyBody, "")

	// Try to resend verification
	resendBody := `{
		"email": "test@example.com"
	}`

	resp, err := MakeRequest("POST", "/api/v1/auth/resend-verification", resendBody, "")
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var resendResponse model.WebResponse[model.ResendVerificationResponse]
	json.NewDecoder(resp.Body).Decode(&resendResponse)
	assert.Equal(t, "Email already verified", resendResponse.Data.Message)
}
