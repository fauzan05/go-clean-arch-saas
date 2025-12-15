package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCurrentUser_Success(t *testing.T) {
	CleanupDatabase(t)

	token := GetAccessToken(t)

	resp, err := MakeRequest("GET", "/api/v1/users/current", "", token)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	result := ParseResponse(t, resp)
	data := result["data"].(map[string]interface{})

	assert.Equal(t, "Test User", data["name"])
	assert.Equal(t, "test@example.com", data["email"])
	assert.NotEmpty(t, data["id"])
	assert.NotEmpty(t, data["organization_id"])
}

func TestGetCurrentUser_Unauthorized(t *testing.T) {
	CleanupDatabase(t)

	resp, err := MakeRequest("GET", "/api/v1/users/current", "", "")
	assert.NoError(t, err)
	assert.Equal(t, 401, resp.StatusCode)
}

func TestGetCurrentUser_InvalidToken(t *testing.T) {
	CleanupDatabase(t)

	resp, err := MakeRequest("GET", "/api/v1/users/current", "", "invalid-token")
	assert.NoError(t, err)
	assert.Equal(t, 401, resp.StatusCode)
}

func TestUpdateUser_Success(t *testing.T) {
	CleanupDatabase(t)

	token := GetAccessToken(t)

	requestBody := `{
		"name": "Updated Name",
		"password": "newpassword123"
	}`

	resp, err := MakeRequest("PATCH", "/api/v1/users/current", requestBody, token)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	result := ParseResponse(t, resp)
	data := result["data"].(map[string]interface{})

	assert.Equal(t, "Updated Name", data["name"])
	assert.Equal(t, "test@example.com", data["email"])
}

func TestUpdateUser_OnlyName(t *testing.T) {
	CleanupDatabase(t)

	token := GetAccessToken(t)

	requestBody := `{
		"name": "New Name Only"
	}`

	resp, err := MakeRequest("PATCH", "/api/v1/users/current", requestBody, token)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	result := ParseResponse(t, resp)
	data := result["data"].(map[string]interface{})

	assert.Equal(t, "New Name Only", data["name"])
}

func TestUpdateUser_OnlyPassword(t *testing.T) {
	CleanupDatabase(t)

	token := GetAccessToken(t)

	requestBody := `{
		"password": "newpassword456"
	}`

	resp, err := MakeRequest("PATCH", "/api/v1/users/current", requestBody, token)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Try to login with new password
	loginBody := `{
		"email": "test@example.com",
		"password": "newpassword456"
	}`

	resp, err = MakeRequest("POST", "/api/v1/auth/login", loginBody, "")
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestUpdateUser_Unauthorized(t *testing.T) {
	CleanupDatabase(t)

	requestBody := `{
		"name": "Hacker"
	}`

	resp, err := MakeRequest("PATCH", "/api/v1/users/current", requestBody, "")
	assert.NoError(t, err)
	assert.Equal(t, 401, resp.StatusCode)
}

func TestUpdateUser_EmptyBody(t *testing.T) {
	CleanupDatabase(t)

	token := GetAccessToken(t)

	requestBody := `{}`

	resp, err := MakeRequest("PATCH", "/api/v1/users/current", requestBody, token)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Verify nothing changed
	resp, err = MakeRequest("GET", "/api/v1/users/current", "", token)
	assert.NoError(t, err)

	result := ParseResponse(t, resp)
	data := result["data"].(map[string]interface{})

	assert.Equal(t, "Test User", data["name"])
}
