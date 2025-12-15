package test

import (
	"encoding/json"
	"go-clean-arch-saas/internal/entity"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// CleanupDatabase clears all tables for testing
func CleanupDatabase(t *testing.T) {
	// Delete in correct order to respect foreign key constraints
	err := db.Exec("SET FOREIGN_KEY_CHECKS = 0").Error
	assert.NoError(t, err)

	err = db.Exec("TRUNCATE TABLE audit_logs").Error
	assert.NoError(t, err)

	err = db.Exec("TRUNCATE TABLE subscriptions").Error
	assert.NoError(t, err)

	err = db.Exec("TRUNCATE TABLE organization_members").Error
	assert.NoError(t, err)

	err = db.Exec("TRUNCATE TABLE users").Error
	assert.NoError(t, err)

	err = db.Exec("TRUNCATE TABLE organizations").Error
	assert.NoError(t, err)

	err = db.Exec("TRUNCATE TABLE plans").Error
	assert.NoError(t, err)

	err = db.Exec("SET FOREIGN_KEY_CHECKS = 1").Error
	assert.NoError(t, err)
}

// CreateTestPlan creates a test plan in the database
func CreateTestPlan(t *testing.T, slug string, name string, price float64) *entity.Plan {
	plan := &entity.Plan{
		ID:            uuid.New().String(),
		Slug:          slug,
		Name:          name,
		Price:         price,
		BillingPeriod: "monthly",
		Features:      `["feature1", "feature2"]`,
		Limits:        `{"users": 10, "projects": 5}`,
		IsActive:      true,
	}

	err := db.Create(plan).Error
	assert.NoError(t, err)

	return plan
}

// ParseResponse parses JSON response body into a map
func ParseResponse(t *testing.T, resp *http.Response) map[string]interface{} {
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	assert.NoError(t, err)

	return result
}

// MakeRequest makes an HTTP request for testing
func MakeRequest(method, url string, body string, token string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	return app.Test(req, -1)
}

// GetAccessToken registers a user and returns the access token
func GetAccessToken(t *testing.T) string {
	// Create a test plan first
	CreateTestPlan(t, "free", "Free Plan", 0)

	// Register user
	registerBody := `{
		"name": "Test User",
		"email": "test@example.com",
		"password": "password123",
		"organization_name": "Test Org"
	}`

	resp, err := MakeRequest("POST", "/api/v1/auth/register", registerBody, "")
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Login to get access token
	loginBody := `{
		"email": "test@example.com",
		"password": "password123"
	}`

	resp, err = MakeRequest("POST", "/api/v1/auth/login", loginBody, "")
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	result := ParseResponse(t, resp)
	data := result["data"].(map[string]interface{})

	return data["access_token"].(string)
}
