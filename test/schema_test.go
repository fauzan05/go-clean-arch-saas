package test

import (
	"go-clean-arch-saas/internal/entity"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDatabaseSchemaUUID(t *testing.T) {
	t.Run("should have CHAR(36) for organization ID", func(t *testing.T) {
		var result string
		err := db.Raw("SELECT COLUMN_TYPE FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'organizations' AND COLUMN_NAME = 'id' AND TABLE_SCHEMA = ?", viperConfig.GetString("database.name")).Scan(&result).Error
		assert.Nil(t, err)
		assert.Equal(t, "char(36)", result, "Organization ID should be CHAR(36) for UUID")
	})

	t.Run("should have CHAR(36) for user ID", func(t *testing.T) {
		var result string
		err := db.Raw("SELECT COLUMN_TYPE FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'users' AND COLUMN_NAME = 'id' AND TABLE_SCHEMA = ?", viperConfig.GetString("database.name")).Scan(&result).Error
		assert.Nil(t, err)
		assert.Equal(t, "char(36)", result, "User ID should be CHAR(36) for UUID")
	})

	t.Run("should have CHAR(36) for plan ID", func(t *testing.T) {
		var result string
		err := db.Raw("SELECT COLUMN_TYPE FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'plans' AND COLUMN_NAME = 'id' AND TABLE_SCHEMA = ?", viperConfig.GetString("database.name")).Scan(&result).Error
		assert.Nil(t, err)
		assert.Equal(t, "char(36)", result, "Plan ID should be CHAR(36) for UUID")
	})

	t.Run("should have CHAR(36) for subscription ID", func(t *testing.T) {
		var result string
		err := db.Raw("SELECT COLUMN_TYPE FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'subscriptions' AND COLUMN_NAME = 'id' AND TABLE_SCHEMA = ?", viperConfig.GetString("database.name")).Scan(&result).Error
		assert.Nil(t, err)
		assert.Equal(t, "char(36)", result, "Subscription ID should be CHAR(36) for UUID")
	})
}

func TestDatabaseSchemaSoftDelete(t *testing.T) {
	t.Run("should have deleted_at column in organizations", func(t *testing.T) {
		var result string
		err := db.Raw("SELECT COLUMN_NAME FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'organizations' AND COLUMN_NAME = 'deleted_at' AND TABLE_SCHEMA = ?", viperConfig.GetString("database.name")).Scan(&result).Error
		assert.Nil(t, err)
		assert.Equal(t, "deleted_at", result, "Organizations should have deleted_at column")
	})

	t.Run("should have deleted_at column in users", func(t *testing.T) {
		var result string
		err := db.Raw("SELECT COLUMN_NAME FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'users' AND COLUMN_NAME = 'deleted_at' AND TABLE_SCHEMA = ?", viperConfig.GetString("database.name")).Scan(&result).Error
		assert.Nil(t, err)
		assert.Equal(t, "deleted_at", result, "Users should have deleted_at column")
	})

	t.Run("should have deleted_at column in plans", func(t *testing.T) {
		var result string
		err := db.Raw("SELECT COLUMN_NAME FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'plans' AND COLUMN_NAME = 'deleted_at' AND TABLE_SCHEMA = ?", viperConfig.GetString("database.name")).Scan(&result).Error
		assert.Nil(t, err)
		assert.Equal(t, "deleted_at", result, "Plans should have deleted_at column")
	})

	t.Run("should have deleted_at column in subscriptions", func(t *testing.T) {
		var result string
		err := db.Raw("SELECT COLUMN_NAME FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'subscriptions' AND COLUMN_NAME = 'deleted_at' AND TABLE_SCHEMA = ?", viperConfig.GetString("database.name")).Scan(&result).Error
		assert.Nil(t, err)
		assert.Equal(t, "deleted_at", result, "Subscriptions should have deleted_at column")
	})

	t.Run("should have deleted_at column in organization_members", func(t *testing.T) {
		var result string
		err := db.Raw("SELECT COLUMN_NAME FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'organization_members' AND COLUMN_NAME = 'deleted_at' AND TABLE_SCHEMA = ?", viperConfig.GetString("database.name")).Scan(&result).Error
		assert.Nil(t, err)
		assert.Equal(t, "deleted_at", result, "Organization members should have deleted_at column")
	})

	t.Run("should have deleted_at column in audit_logs", func(t *testing.T) {
		var result string
		err := db.Raw("SELECT COLUMN_NAME FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'audit_logs' AND COLUMN_NAME = 'deleted_at' AND TABLE_SCHEMA = ?", viperConfig.GetString("database.name")).Scan(&result).Error
		assert.Nil(t, err)
		assert.Equal(t, "deleted_at", result, "Audit logs should have deleted_at column")
	})
}

func TestSoftDeleteFunctionality(t *testing.T) {
	CleanupDatabase(t)

	t.Run("should soft delete organization and exclude from normal queries", func(t *testing.T) {
		org := &entity.Organization{
			ID:   "test-org-soft-delete-1",
			Name: "Test Org Soft Delete",
			Slug: "test-org-soft-delete",
		}
		err := db.Create(org).Error
		assert.Nil(t, err)

		// Verify organization exists
		var foundOrg entity.Organization
		err = db.Where("id = ?", org.ID).First(&foundOrg).Error
		assert.Nil(t, err)
		assert.Equal(t, "Test Org Soft Delete", foundOrg.Name)
		assert.Nil(t, foundOrg.DeletedAt, "DeletedAt should be nil for active record")

		// Soft delete by setting deleted_at timestamp
		now := time.Now().UnixMilli()
		err = db.Model(&entity.Organization{}).Where("id = ?", org.ID).Update("deleted_at", now).Error
		assert.Nil(t, err)

		// Normal query without filter should still find it (because we don't have automatic soft delete in GORM)
		err = db.Where("id = ?", org.ID).First(&foundOrg).Error
		assert.Nil(t, err)
		assert.NotNil(t, foundOrg.DeletedAt, "DeletedAt should be set")
		assert.Equal(t, now, *foundOrg.DeletedAt, "DeletedAt should match timestamp")

		// Query with explicit deleted_at IS NULL filter should NOT find it
		err = db.Where("id = ? AND deleted_at IS NULL", org.ID).First(&foundOrg).Error
		assert.NotNil(t, err, "Should not find soft deleted organization when filtering by deleted_at IS NULL")
		assert.Equal(t, "record not found", err.Error())
	})

	t.Run("should soft delete user and maintain referential integrity", func(t *testing.T) {
		org := &entity.Organization{
			ID:   "test-org-soft-delete-2",
			Name: "Test Org 2",
			Slug: "test-org-2",
		}
		db.Create(org)

		user := &entity.User{
			ID:             "test-user-soft-delete-1",
			Name:           "Test User",
			Email:          "soft-delete@example.com",
			Password:       "password",
			OrganizationID: org.ID,
		}
		err := db.Create(user).Error
		assert.Nil(t, err)

		// Verify user exists
		var foundUser entity.User
		err = db.Where("id = ?", user.ID).First(&foundUser).Error
		assert.Nil(t, err)
		assert.Equal(t, "Test User", foundUser.Name)

		// Soft delete user
		now := time.Now().UnixMilli()
		err = db.Model(&entity.User{}).Where("id = ?", user.ID).Update("deleted_at", now).Error
		assert.Nil(t, err)

		// Query with filter should not find soft deleted user
		err = db.Where("id = ? AND deleted_at IS NULL", user.ID).First(&foundUser).Error
		assert.NotNil(t, err, "Should not find soft deleted user")

		// Verify organization still exists (referential integrity maintained)
		var checkOrg entity.Organization
		err = db.Where("id = ?", org.ID).First(&checkOrg).Error
		assert.Nil(t, err, "Organization should still exist")
	})

	t.Run("should restore soft deleted record by setting deleted_at to NULL", func(t *testing.T) {
		org := &entity.Organization{
			ID:   "test-org-restore-1",
			Name: "Test Org Restore",
			Slug: "test-org-restore",
		}
		db.Create(org)

		// Soft delete
		now := time.Now().UnixMilli()
		db.Model(&entity.Organization{}).Where("id = ?", org.ID).Update("deleted_at", now)

		// Verify it's soft deleted
		var foundOrg entity.Organization
		err := db.Where("id = ? AND deleted_at IS NULL", org.ID).First(&foundOrg).Error
		assert.NotNil(t, err, "Should not find soft deleted organization")

		// Restore by setting deleted_at to NULL
		err = db.Model(&entity.Organization{}).Where("id = ?", org.ID).Update("deleted_at", nil).Error
		assert.Nil(t, err)

		// Verify it's restored
		err = db.Where("id = ? AND deleted_at IS NULL", org.ID).First(&foundOrg).Error
		assert.Nil(t, err, "Should find restored organization")
		assert.Nil(t, foundOrg.DeletedAt, "DeletedAt should be nil after restore")
		assert.Equal(t, "Test Org Restore", foundOrg.Name)
	})

	t.Run("should count only active records when filtering by deleted_at", func(t *testing.T) {
		// Create 3 organizations
		orgs := []entity.Organization{
			{ID: "test-org-count-1", Name: "Org 1", Slug: "org-1"},
			{ID: "test-org-count-2", Name: "Org 2", Slug: "org-2"},
			{ID: "test-org-count-3", Name: "Org 3", Slug: "org-3"},
		}
		for _, org := range orgs {
			db.Create(&org)
		}

		// Count all
		var totalCount int64
		db.Model(&entity.Organization{}).Where("id LIKE ?", "test-org-count-%").Count(&totalCount)
		assert.Equal(t, int64(3), totalCount, "Should have 3 organizations")

		// Soft delete one
		now := time.Now().UnixMilli()
		db.Model(&entity.Organization{}).Where("id = ?", "test-org-count-2").Update("deleted_at", now)

		// Count active only
		var activeCount int64
		db.Model(&entity.Organization{}).Where("id LIKE ? AND deleted_at IS NULL", "test-org-count-%").Count(&activeCount)
		assert.Equal(t, int64(2), activeCount, "Should have 2 active organizations")

		// Count deleted only
		var deletedCount int64
		db.Model(&entity.Organization{}).Where("id LIKE ? AND deleted_at IS NOT NULL", "test-org-count-%").Count(&deletedCount)
		assert.Equal(t, int64(1), deletedCount, "Should have 1 deleted organization")
	})

	t.Run("should handle subscription soft delete", func(t *testing.T) {
		org := &entity.Organization{
			ID:   "test-org-sub-delete-1",
			Name: "Test Org Sub Delete",
			Slug: "test-org-sub-delete",
		}
		err := db.Create(org).Error
		assert.Nil(t, err)

		plan := &entity.Plan{
			ID:            "test-plan-delete-1",
			Name:          "Test Plan",
			Slug:          "test-plan-delete",
			Price:         9.99,
			BillingPeriod: "monthly",
			Features:      "{}",
			Limits:        "{}",
			IsActive:      true,
		}
		err = db.Create(plan).Error
		assert.Nil(t, err)

		subscription := &entity.Subscription{
			ID:                 "test-sub-delete-1",
			OrganizationID:     org.ID,
			PlanID:             plan.ID,
			Status:             "active",
			CurrentPeriodStart: time.Now().UnixMilli(),
			CurrentPeriodEnd:   time.Now().Add(30 * 24 * time.Hour).UnixMilli(),
		}
		err = db.Create(subscription).Error
		assert.Nil(t, err)

		// Soft delete subscription
		now := time.Now().UnixMilli()
		err = db.Model(&entity.Subscription{}).Where("id = ?", subscription.ID).Update("deleted_at", now).Error
		assert.Nil(t, err)

		// Verify subscription is soft deleted
		var foundSub entity.Subscription
		err = db.Where("id = ? AND deleted_at IS NULL", subscription.ID).First(&foundSub).Error
		assert.NotNil(t, err, "Should not find soft deleted subscription")

		// Verify organization and plan still exist (referential integrity maintained)
		var checkOrg entity.Organization
		var checkPlan entity.Plan
		err = db.Where("id = ?", org.ID).First(&checkOrg).Error
		assert.Nil(t, err, "Should find organization")
		assert.Equal(t, org.ID, checkOrg.ID, "Organization should still exist")

		err = db.Where("id = ?", plan.ID).First(&checkPlan).Error
		assert.Nil(t, err, "Should find plan")
		assert.Equal(t, plan.ID, checkPlan.ID, "Plan should still exist")
	})
}
