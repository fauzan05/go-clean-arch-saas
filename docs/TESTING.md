# Testing Guide

This project includes comprehensive unit tests covering all API endpoints.

## Test Files

- `test/auth_test.go` - Authentication tests (register, login, refresh, logout)
- `test/user_test.go` - User management tests (get current, update)
- `test/organization_test.go` - Organization management tests (get, update, members)
- `test/subscription_test.go` - Subscription management tests (get, upgrade, cancel)
- `test/helper_test.go` - Test utilities and helpers
- `test/init.go` - Test initialization and configuration

## Prerequisites

1. MySQL database must be running and accessible
2. Database `go_clean_arch_saas` must exist
3. Migrations must be applied

## Setup Database for Testing

```bash
# Create database
mysql -u root -p -e "CREATE DATABASE IF NOT EXISTS go_clean_arch_saas"

# Or using Docker
docker-compose up -d mysql

# Run migrations
make migrate-up
```

## Running Tests

### Run All Tests

```bash
make test
```

Or manually:

```bash
go test -v ./test/
```

### Run Specific Test File

```bash
# Auth tests only
go test -v ./test/ -run TestRegister
go test -v ./test/ -run TestLogin

# User tests only
go test -v ./test/ -run TestGetCurrentUser
go test -v ./test/ -run TestUpdateUser

# Organization tests only
go test -v ./test/ -run TestGetCurrentOrganization
go test -v ./test/ -run TestListOrganizationMembers

# Subscription tests only
go test -v ./test/ -run TestGetCurrentSubscription
go test -v ./test/ -run TestUpgradeSubscription
```

### Run with Coverage

```bash
make test-coverage
```

Or manually:

```bash
go test -v ./test/ -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

## Test Coverage

The test suite covers:

### Health Checks (2 tests)
- ✅ Health endpoint
- ✅ Readiness endpoint with database check

### Authentication (10 tests)
- ✅ Register new user with organization
- ✅ Register with duplicate email (409 error)
- ✅ Register with invalid request (400 error)
- ✅ Login with valid credentials
- ✅ Login with wrong password (401 error)
- ✅ Login with non-existent user (401 error)
- ✅ Refresh token with valid token
- ✅ Refresh token with invalid token (401 error)
- ✅ Logout successfully
- ✅ Logout without authentication (401 error)

### User Management (8 tests)
- ✅ Get current user
- ✅ Get current user without authentication (401 error)
- ✅ Get current user with invalid token (401 error)
- ✅ Update user (name and password)
- ✅ Update user (name only)
- ✅ Update user (password only)
- ✅ Update user without authentication (401 error)
- ✅ Update user with empty body

### Organization Management (8 tests)
- ✅ Get current organization
- ✅ Get current organization without authentication (401 error)
- ✅ Update organization name
- ✅ Update organization without authentication (401 error)
- ✅ Update organization with empty name (400 error)
- ✅ List organization members
- ✅ List organization members without authentication (401 error)
- ✅ Cannot remove owner from organization
- ✅ Remove member with non-existent user ID (404 error)
- ✅ Remove member without authentication (401 error)

### Subscription Management (11 tests)
- ✅ Get current subscription
- ✅ Get current subscription without authentication (401 error)
- ✅ Upgrade subscription to higher plan
- ✅ Upgrade subscription with non-existent plan (404 error)
- ✅ Upgrade subscription without authentication (401 error)
- ✅ Upgrade subscription with invalid request (400 error)
- ✅ Downgrade subscription to lower plan
- ✅ Cancel subscription
- ✅ Cancel subscription without authentication (401 error)
- ✅ Cancel already cancelled subscription
- ✅ Complete subscription workflow (free → basic → pro → enterprise → basic → cancelled)

## Test Utilities

### Helper Functions

```go
// CleanupDatabase - Clears all tables before each test
CleanupDatabase(t)

// CreateTestPlan - Creates a test subscription plan
plan := CreateTestPlan(t, "pro", "Pro Plan", 29.00)

// ParseResponse - Parses JSON response into map
result := ParseResponse(t, resp)

// MakeRequest - Makes HTTP request with optional auth
resp, err := MakeRequest("GET", "/api/v1/users/current", "", token)

// GetAccessToken - Registers user and returns access token
token := GetAccessToken(t)
```

## Test Data

Each test runs in isolation with:
- Fresh database (all tables truncated)
- New test user (test@example.com / password123)
- New test organization (Test Org)
- Free plan subscription

## CI/CD Integration

Tests can be integrated into CI/CD pipelines:

```yaml
# Example GitHub Actions
- name: Run tests
  run: |
    docker-compose up -d mysql
    sleep 10
    make migrate-up
    make test
```

## Troubleshooting

### Database Connection Error

```
failed to connect database: Error 1049 (42000): Unknown database
```

**Solution**: Create the database first:
```bash
mysql -u root -p -e "CREATE DATABASE go_clean_arch_saas"
```

### Migration Not Applied

```
Error 1146 (42S02): Table 'go_clean_arch_saas.users' doesn't exist
```

**Solution**: Run migrations:
```bash
make migrate-up
```

### Port Already in Use

```
bind: address already in use
```

**Solution**: Stop the running application before testing:
```bash
# Kill process on port 3000
lsof -ti:3000 | xargs kill -9
```

## Writing New Tests

When adding new features, follow this pattern:

```go
func TestYourFeature_Success(t *testing.T) {
    CleanupDatabase(t)
    token := GetAccessToken(t)

    requestBody := `{"field": "value"}`
    
    resp, err := MakeRequest("POST", "/api/v1/your/endpoint", requestBody, token)
    assert.NoError(t, err)
    assert.Equal(t, 200, resp.StatusCode)

    result := ParseResponse(t, resp)
    data := result["data"].(map[string]interface{})
    
    assert.Equal(t, "expected", data["field"])
}
```

## Best Practices

1. **Always cleanup database** before each test with `CleanupDatabase(t)`
2. **Test both success and error cases** for comprehensive coverage
3. **Use descriptive test names** following pattern: `Test{Feature}_{Scenario}`
4. **Assert HTTP status codes** to verify correct error handling
5. **Parse and validate response data** to ensure API contract
6. **Test authorization** by running tests with and without tokens
7. **Test edge cases** like empty inputs, invalid IDs, etc.

## Performance

Tests typically complete in:
- Individual test: 50-200ms
- Full test suite: 3-5 seconds
- With coverage: 5-8 seconds

## Next Steps

- [ ] Add integration tests with external services
- [ ] Add load/stress testing
- [ ] Add API contract testing
- [ ] Mock external dependencies
- [ ] Add benchmark tests
