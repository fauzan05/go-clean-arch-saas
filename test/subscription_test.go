package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCurrentSubscription_Success(t *testing.T) {
	CleanupDatabase(t)

	token := GetAccessToken(t)

	resp, err := MakeRequest("GET", "/api/v1/subscriptions/current", "", token)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	result := ParseResponse(t, resp)
	data := result["data"].(map[string]interface{})

	assert.NotEmpty(t, data["id"])
	assert.NotEmpty(t, data["organization_id"])
	assert.Equal(t, "active", data["status"])

	// Check plan object (not plan_id, as plan is embedded)
	plan := data["plan"].(map[string]interface{})
	assert.NotEmpty(t, plan["id"])
	assert.Equal(t, "free", plan["slug"])
	assert.Equal(t, "Free Plan", plan["name"])
}

func TestGetCurrentSubscription_Unauthorized(t *testing.T) {
	CleanupDatabase(t)

	resp, err := MakeRequest("GET", "/api/v1/subscriptions/current", "", "")
	assert.NoError(t, err)
	assert.Equal(t, 401, resp.StatusCode)
}

func TestUpgradeSubscription_Success(t *testing.T) {
	CleanupDatabase(t)

	// Create pro plan
	proPlan := CreateTestPlan(t, "pro", "Pro Plan", 29.00)

	token := GetAccessToken(t)

	// Upgrade to pro plan using plan_id
	requestBody := `{
		"plan_id": "` + proPlan.ID + `"
	}`

	resp, err := MakeRequest("POST", "/api/v1/subscriptions/upgrade", requestBody, token)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	result := ParseResponse(t, resp)
	data := result["data"].(map[string]interface{})

	assert.Equal(t, "active", data["status"])
	assert.NotNil(t, data["plan"])

	plan := data["plan"].(map[string]interface{})
	assert.Equal(t, "pro", plan["slug"])
	assert.Equal(t, "Pro Plan", plan["name"])
	assert.Equal(t, float64(29), plan["price"])
}

func TestUpgradeSubscription_PlanNotFound(t *testing.T) {
	CleanupDatabase(t)

	token := GetAccessToken(t)

	// Use a non-existent UUID
	requestBody := `{
		"plan_id": "00000000-0000-0000-0000-000000000000"
	}`

	resp, err := MakeRequest("POST", "/api/v1/subscriptions/upgrade", requestBody, token)
	assert.NoError(t, err)
	assert.Equal(t, 404, resp.StatusCode)
}

func TestUpgradeSubscription_Unauthorized(t *testing.T) {
	CleanupDatabase(t)

	requestBody := `{
		"plan_id": "00000000-0000-0000-0000-000000000000"
	}`

	resp, err := MakeRequest("POST", "/api/v1/subscriptions/upgrade", requestBody, "")
	assert.NoError(t, err)
	assert.Equal(t, 401, resp.StatusCode)
}

func TestUpgradeSubscription_InvalidRequest(t *testing.T) {
	CleanupDatabase(t)

	token := GetAccessToken(t)

	// Empty plan_id should fail validation
	requestBody := `{
		"plan_id": ""
	}`

	resp, err := MakeRequest("POST", "/api/v1/subscriptions/upgrade", requestBody, token)
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)
}

func TestDowngradeSubscription_Success(t *testing.T) {
	CleanupDatabase(t)

	// Create both plans
	proPlan := CreateTestPlan(t, "pro", "Pro Plan", 29.00)

	token := GetAccessToken(t)

	// First upgrade to pro
	upgradeBody := `{
		"plan_id": "` + proPlan.ID + `"
	}`
	resp, err := MakeRequest("POST", "/api/v1/subscriptions/upgrade", upgradeBody, token)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Then downgrade back to free
	freePlan := CreateTestPlan(t, "free-downgrade", "Free Downgrade", 0)
	downgradeBody := `{
		"plan_id": "` + freePlan.ID + `"
	}`
	resp, err = MakeRequest("POST", "/api/v1/subscriptions/upgrade", downgradeBody, token)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	result := ParseResponse(t, resp)
	data := result["data"].(map[string]interface{})

	plan := data["plan"].(map[string]interface{})
	assert.Equal(t, "free-downgrade", plan["slug"])
}

func TestCancelSubscription_Success(t *testing.T) {
	CleanupDatabase(t)

	token := GetAccessToken(t)

	resp, err := MakeRequest("POST", "/api/v1/subscriptions/cancel", "", token)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	result := ParseResponse(t, resp)
	// Response data is a string message, not an object
	data := result["data"].(string)
	assert.Equal(t, "Subscription cancelled successfully", data)
}

func TestCancelSubscription_Unauthorized(t *testing.T) {
	CleanupDatabase(t)

	resp, err := MakeRequest("POST", "/api/v1/subscriptions/cancel", "", "")
	assert.NoError(t, err)
	assert.Equal(t, 401, resp.StatusCode)
}

func TestCancelSubscription_AlreadyCancelled(t *testing.T) {
	CleanupDatabase(t)

	token := GetAccessToken(t)

	// Cancel first time
	resp, err := MakeRequest("POST", "/api/v1/subscriptions/cancel", "", token)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Try to cancel again - should return 404 since no active subscription exists
	resp, err = MakeRequest("POST", "/api/v1/subscriptions/cancel", "", token)
	assert.NoError(t, err)
	assert.Equal(t, 404, resp.StatusCode) // No active subscription to cancel
}

func TestSubscriptionWorkflow_Complete(t *testing.T) {
	CleanupDatabase(t)

	// Create multiple plans
	basicPlan := CreateTestPlan(t, "basic", "Basic Plan", 9.00)
	proPlan := CreateTestPlan(t, "pro", "Pro Plan", 29.00)
	enterprisePlan := CreateTestPlan(t, "enterprise", "Enterprise Plan", 99.00)

	token := GetAccessToken(t)

	// 1. Check initial subscription (free)
	resp, err := MakeRequest("GET", "/api/v1/subscriptions/current", "", token)
	assert.NoError(t, err)
	result := ParseResponse(t, resp)
	data := result["data"].(map[string]interface{})
	plan := data["plan"].(map[string]interface{})
	assert.Equal(t, "free", plan["slug"])

	// 2. Upgrade to basic
	upgradeBody := `{"plan_id": "` + basicPlan.ID + `"}`
	resp, err = MakeRequest("POST", "/api/v1/subscriptions/upgrade", upgradeBody, token)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// 3. Upgrade to pro
	upgradeBody = `{"plan_id": "` + proPlan.ID + `"}`
	resp, err = MakeRequest("POST", "/api/v1/subscriptions/upgrade", upgradeBody, token)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	result = ParseResponse(t, resp)
	data = result["data"].(map[string]interface{})
	plan = data["plan"].(map[string]interface{})
	assert.Equal(t, "pro", plan["slug"])

	// 4. Upgrade to enterprise
	upgradeBody = `{"plan_id": "` + enterprisePlan.ID + `"}`
	resp, err = MakeRequest("POST", "/api/v1/subscriptions/upgrade", upgradeBody, token)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// 5. Downgrade to basic
	downgradeBody := `{"plan_id": "` + basicPlan.ID + `"}`
	resp, err = MakeRequest("POST", "/api/v1/subscriptions/upgrade", downgradeBody, token)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// 6. Cancel subscription
	resp, err = MakeRequest("POST", "/api/v1/subscriptions/cancel", "", token)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	result = ParseResponse(t, resp)
	message := result["data"].(string)
	assert.Equal(t, "Subscription cancelled successfully", message)
}
