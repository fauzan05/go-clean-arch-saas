package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCurrentOrganization_Success(t *testing.T) {
	CleanupDatabase(t)

	token := GetAccessToken(t)

	resp, err := MakeRequest("GET", "/api/v1/organizations/current", "", token)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	result := ParseResponse(t, resp)
	data := result["data"].(map[string]interface{})

	assert.Equal(t, "Test Org", data["name"])
	assert.Equal(t, "test-org", data["slug"])
	assert.NotEmpty(t, data["id"])
}

func TestGetCurrentOrganization_Unauthorized(t *testing.T) {
	CleanupDatabase(t)

	resp, err := MakeRequest("GET", "/api/v1/organizations/current", "", "")
	assert.NoError(t, err)
	assert.Equal(t, 401, resp.StatusCode)
}

func TestUpdateOrganization_Success(t *testing.T) {
	CleanupDatabase(t)

	token := GetAccessToken(t)

	requestBody := `{
		"name": "Updated Organization Name"
	}`

	resp, err := MakeRequest("PATCH", "/api/v1/organizations/current", requestBody, token)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	result := ParseResponse(t, resp)
	data := result["data"].(map[string]interface{})

	assert.Equal(t, "Updated Organization Name", data["name"])
	assert.Equal(t, "test-org", data["slug"]) // Slug should not change
}

func TestUpdateOrganization_Unauthorized(t *testing.T) {
	CleanupDatabase(t)

	requestBody := `{
		"name": "Hacker Org"
	}`

	resp, err := MakeRequest("PATCH", "/api/v1/organizations/current", requestBody, "")
	assert.NoError(t, err)
	assert.Equal(t, 401, resp.StatusCode)
}

func TestUpdateOrganization_EmptyName(t *testing.T) {
	CleanupDatabase(t)

	token := GetAccessToken(t)

	requestBody := `{
		"name": ""
	}`

	resp, err := MakeRequest("PATCH", "/api/v1/organizations/current", requestBody, token)
	assert.NoError(t, err)
	// Empty name is allowed, it won't update (returns 200 OK)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestListOrganizationMembers_Success(t *testing.T) {
	CleanupDatabase(t)

	token := GetAccessToken(t)

	resp, err := MakeRequest("GET", "/api/v1/organizations/members", "", token)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	result := ParseResponse(t, resp)
	data := result["data"].([]interface{})

	// Should have at least one member (the owner)
	assert.GreaterOrEqual(t, len(data), 1)

	member := data[0].(map[string]interface{})
	assert.NotEmpty(t, member["user_id"])
	assert.Equal(t, "owner", member["role"])
	assert.NotEmpty(t, member["joined_at"])

	// Check user object (should be embedded in response)
	user := member["user"].(map[string]interface{})
	assert.NotEmpty(t, user["id"])
	assert.Equal(t, "Test User", user["name"])
	assert.Equal(t, "test@example.com", user["email"])
}

func TestListOrganizationMembers_Unauthorized(t *testing.T) {
	CleanupDatabase(t)

	resp, err := MakeRequest("GET", "/api/v1/organizations/members", "", "")
	assert.NoError(t, err)
	assert.Equal(t, 401, resp.StatusCode)
}

func TestRemoveOrganizationMember_CannotRemoveOwner(t *testing.T) {
	CleanupDatabase(t)

	token := GetAccessToken(t)

	// Get current user ID
	resp, err := MakeRequest("GET", "/api/v1/users/current", "", token)
	assert.NoError(t, err)

	result := ParseResponse(t, resp)
	userData := result["data"].(map[string]interface{})
	userID := userData["id"].(string)

	// Try to remove owner (should fail with 403 Forbidden)
	resp, err = MakeRequest("DELETE", "/api/v1/organizations/members/"+userID, "", token)
	assert.NoError(t, err)
	assert.Equal(t, 403, resp.StatusCode)
}

func TestRemoveOrganizationMember_UserNotFound(t *testing.T) {
	CleanupDatabase(t)

	token := GetAccessToken(t)

	resp, err := MakeRequest("DELETE", "/api/v1/organizations/members/non-existent-user-id", "", token)
	assert.NoError(t, err)
	assert.Equal(t, 404, resp.StatusCode)
}

func TestRemoveOrganizationMember_Unauthorized(t *testing.T) {
	CleanupDatabase(t)

	resp, err := MakeRequest("DELETE", "/api/v1/organizations/members/some-user-id", "", "")
	assert.NoError(t, err)
	assert.Equal(t, 401, resp.StatusCode)
}
