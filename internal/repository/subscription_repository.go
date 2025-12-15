package repository

import (
	"go-clean-arch-saas/internal/entity"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type SubscriptionRepository struct {
	Repository[entity.Subscription]
	Log *logrus.Logger
}

func NewSubscriptionRepository(log *logrus.Logger) *SubscriptionRepository {
	return &SubscriptionRepository{
		Log: log,
	}
}

func (r *SubscriptionRepository) FindByOrganization(db *gorm.DB, subscription *entity.Subscription, orgID string) error {
	return db.Where("organization_id = ?", orgID).Preload("Plan").First(subscription).Error
}

func (r *SubscriptionRepository) FindActiveByOrganization(db *gorm.DB, subscription *entity.Subscription, orgID string) error {
	return db.Where("organization_id = ? AND status = ?", orgID, "active").Preload("Plan").First(subscription).Error
}
