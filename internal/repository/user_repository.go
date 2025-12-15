package repository

import (
	"go-clean-arch-saas/internal/entity"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserRepository struct {
	Repository[entity.User]
	Log *logrus.Logger
}

func NewUserRepository(log *logrus.Logger) *UserRepository {
	return &UserRepository{
		Log: log,
	}
}

func (r *UserRepository) FindByEmail(db *gorm.DB, user *entity.User, email string) error {
	return db.Where("email = ?", email).First(user).Error
}

func (r *UserRepository) FindByRefreshToken(db *gorm.DB, user *entity.User, refreshToken string) error {
	return db.Where("refresh_token = ?", refreshToken).First(user).Error
}

func (r *UserRepository) CountByEmail(db *gorm.DB, email string) (int64, error) {
	var count int64
	err := db.Model(&entity.User{}).Where("email = ?", email).Count(&count).Error
	return count, err
}

func (r *UserRepository) FindByVerificationToken(db *gorm.DB, user *entity.User, token string) error {
	return db.Where("verification_token = ?", token).First(user).Error
}
