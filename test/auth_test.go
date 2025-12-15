package test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	CleanupDatabase(t)

	resp, err := MakeRequest("GET", "/health", "", "")
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	result := ParseResponse(t, resp)
	assert.Equal(t, "ok", result["status"])
}

func TestReadinessCheck(t *testing.T) {
	CleanupDatabase(t)

	resp, err := MakeRequest("GET", "/ready", "", "")
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	result := ParseResponse(t, resp)
	assert.Equal(t, "ok", result["status"])
	assert.Equal(t, "connected", result["database"])
}

func TestRegister_Success(t *testing.T) {
	CleanupDatabase(t)

	// Create test plan
	CreateTestPlan(t, "free", "Free Plan", 0)

	requestBody := `{
		"name": "John Doe",
		"email": "john@example.com",
		"password": "password123",
		"organization_name": "Acme Corp"
	}`

	resp, err := MakeRequest("POST", "/api/v1/auth/register", requestBody, "")
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	result := ParseResponse(t, resp)
	data := result["data"].(map[string]interface{})

	// Register returns user and organization, not tokens
	assert.NotNil(t, data["user"])
	assert.NotNil(t, data["organization"])

	user := data["user"].(map[string]interface{})
	assert.Equal(t, "John Doe", user["name"])
	assert.Equal(t, "john@example.com", user["email"])
	assert.NotEmpty(t, user["organization_id"])

	org := data["organization"].(map[string]interface{})
	assert.Equal(t, "Acme Corp", org["name"])
	assert.Equal(t, "acme-corp", org["slug"])
}

func TestRegister_DuplicateEmail(t *testing.T) {
	CleanupDatabase(t)

	// Create test plan
	CreateTestPlan(t, "free", "Free Plan", 0)

	requestBody := `{
		"name": "John Doe",
		"email": "john@example.com",
		"password": "password123",
		"organization_name": "Acme Corp"
	}`

	// First registration
	resp, err := MakeRequest("POST", "/api/v1/auth/register", requestBody, "")
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Second registration with same email
	resp, err = MakeRequest("POST", "/api/v1/auth/register", requestBody, "")
	assert.NoError(t, err)
	assert.Equal(t, 409, resp.StatusCode)
}

func TestRegister_InvalidRequest(t *testing.T) {
	CleanupDatabase(t)

	requestBody := `{
		"name": "",
		"email": "invalid-email",
		"password": "123"
	}`

	resp, err := MakeRequest("POST", "/api/v1/auth/register", requestBody, "")
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)
}

func TestLogin_Success(t *testing.T) {
	CleanupDatabase(t)

	// Create test plan and register user
	CreateTestPlan(t, "free", "Free Plan", 0)

	registerBody := `{
		"name": "Jane Doe",
		"email": "jane@example.com",
		"password": "password123",
		"organization_name": "Test Org"
	}`
	resp, err := MakeRequest("POST", "/api/v1/auth/register", registerBody, "")
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Login
	loginBody := `{
		"email": "jane@example.com",
		"password": "password123"
	}`

	resp, err = MakeRequest("POST", "/api/v1/auth/login", loginBody, "")
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	result := ParseResponse(t, resp)
	data := result["data"].(map[string]interface{})

	assert.NotEmpty(t, data["access_token"])
	assert.NotEmpty(t, data["refresh_token"])
	assert.Equal(t, "Bearer", data["token_type"])
}

func TestLogin_WrongPassword(t *testing.T) {
	CleanupDatabase(t)

	// Create test plan and register user
	CreateTestPlan(t, "free", "Free Plan", 0)

	registerBody := `{
		"name": "Jane Doe",
		"email": "jane@example.com",
		"password": "password123",
		"organization_name": "Test Org"
	}`
	resp, err := MakeRequest("POST", "/api/v1/auth/register", registerBody, "")
	assert.NoError(t, err)

	// Login with wrong password
	loginBody := `{
		"email": "jane@example.com",
		"password": "wrongpassword"
	}`

	resp, err = MakeRequest("POST", "/api/v1/auth/login", loginBody, "")
	assert.NoError(t, err)
	assert.Equal(t, 401, resp.StatusCode)
}

func TestLogin_UserNotFound(t *testing.T) {
	CleanupDatabase(t)

	loginBody := `{
		"email": "notfound@example.com",
		"password": "password123"
	}`

	resp, err := MakeRequest("POST", "/api/v1/auth/login", loginBody, "")
	assert.NoError(t, err)
	assert.Equal(t, 401, resp.StatusCode)
}

func TestRefreshToken_Success(t *testing.T) {
	CleanupDatabase(t)

	// Create test plan and register user
	CreateTestPlan(t, "free", "Free Plan", 0)

	registerBody := `{
		"name": "Test User",
		"email": "test@example.com",
		"password": "password123",
		"organization_name": "Test Org"
	}`

	resp, err := MakeRequest("POST", "/api/v1/auth/register", registerBody, "")
	assert.NoError(t, err)

	// Login to get refresh token
	loginBody := `{
		"email": "test@example.com",
		"password": "password123"
	}`

	resp, err = MakeRequest("POST", "/api/v1/auth/login", loginBody, "")
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	result := ParseResponse(t, resp)
	data := result["data"].(map[string]interface{})
	refreshToken := data["refresh_token"].(string)

	// Refresh token
	refreshBody := fmt.Sprintf(`{"refresh_token": "%s"}`, refreshToken)
	resp, err = MakeRequest("POST", "/api/v1/auth/refresh", refreshBody, "")
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	result = ParseResponse(t, resp)
	data = result["data"].(map[string]interface{})

	assert.NotEmpty(t, data["access_token"])
	assert.Equal(t, "Bearer", data["token_type"])
}

func TestRefreshToken_Invalid(t *testing.T) {
	CleanupDatabase(t)

	refreshBody := `{"refresh_token": "invalid-token"}`
	resp, err := MakeRequest("POST", "/api/v1/auth/refresh", refreshBody, "")
	assert.NoError(t, err)
	assert.Equal(t, 401, resp.StatusCode)
}

func TestLogout_Success(t *testing.T) {
	CleanupDatabase(t)

	token := GetAccessToken(t)

	resp, err := MakeRequest("DELETE", "/api/v1/auth/logout", "", token)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	result := ParseResponse(t, resp)
	data := result["data"].(string)
	assert.Equal(t, "Successfully logged out", data)
}

func TestLogout_Unauthorized(t *testing.T) {
	CleanupDatabase(t)

	resp, err := MakeRequest("DELETE", "/api/v1/auth/logout", "", "")
	assert.NoError(t, err)
	assert.Equal(t, 401, resp.StatusCode)
}
