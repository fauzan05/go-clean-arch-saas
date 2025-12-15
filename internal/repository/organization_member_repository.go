package repository

import (
	"go-clean-arch-saas/internal/entity"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type OrganizationMemberRepository struct {
	Repository[entity.OrganizationMember]
	Log *logrus.Logger
}

func NewOrganizationMemberRepository(log *logrus.Logger) *OrganizationMemberRepository {
	return &OrganizationMemberRepository{
		Log: log,
	}
}

func (r *OrganizationMemberRepository) FindByOrgAndUser(db *gorm.DB, member *entity.OrganizationMember, orgID, userID string) error {
	return db.Where("organization_id = ? AND user_id = ?", orgID, userID).First(member).Error
}

func (r *OrganizationMemberRepository) ListByOrganization(db *gorm.DB, orgID string) ([]entity.OrganizationMember, error) {
	var members []entity.OrganizationMember
	err := db.Where("organization_id = ?", orgID).Preload("User").Find(&members).Error
	return members, err
}

func (r *OrganizationMemberRepository) DeleteByOrgAndUser(db *gorm.DB, orgID, userID string) error {
	return db.Where("organization_id = ? AND user_id = ?", orgID, userID).Delete(&entity.OrganizationMember{}).Error
}
