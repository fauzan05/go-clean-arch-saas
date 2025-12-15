package repository

import (
	"go-clean-arch-saas/internal/entity"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PlanRepository struct {
	Repository[entity.Plan]
	Log *logrus.Logger
}

func NewPlanRepository(log *logrus.Logger) *PlanRepository {
	return &PlanRepository{
		Log: log,
	}
}

func (r *PlanRepository) FindBySlug(db *gorm.DB, plan *entity.Plan, slug string) error {
	return db.Where("slug = ? AND is_active = ?", slug, true).First(plan).Error
}

func (r *PlanRepository) FindAllActive(db *gorm.DB) ([]entity.Plan, error) {
	var plans []entity.Plan
	err := db.Where("is_active = ?", true).Find(&plans).Error
	return plans, err
}
