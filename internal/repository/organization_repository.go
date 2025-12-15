package repository

import (
	"go-clean-arch-saas/internal/entity"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type OrganizationRepository struct {
	Repository[entity.Organization]
	Log *logrus.Logger
}

func NewOrganizationRepository(log *logrus.Logger) *OrganizationRepository {
	return &OrganizationRepository{
		Log: log,
	}
}

func (r *OrganizationRepository) FindBySlug(db *gorm.DB, org *entity.Organization, slug string) error {
	return db.Where("slug = ?", slug).First(org).Error
}

func (r *OrganizationRepository) CountBySlug(db *gorm.DB, slug string) (int64, error) {
	var count int64
	err := db.Model(&entity.Organization{}).Where("slug = ?", slug).Count(&count).Error
	return count, err
}
