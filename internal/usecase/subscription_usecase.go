package usecase

import (
	"context"
	"go-clean-arch-saas/internal/entity"
	"go-clean-arch-saas/internal/model"
	"go-clean-arch-saas/internal/model/converter"
	"go-clean-arch-saas/internal/repository"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type SubscriptionUseCase struct {
	DB                     *gorm.DB
	Log                    *logrus.Logger
	Validate               *validator.Validate
	SubscriptionRepository *repository.SubscriptionRepository
	PlanRepository         *repository.PlanRepository
}

func NewSubscriptionUseCase(
	db *gorm.DB,
	logger *logrus.Logger,
	validate *validator.Validate,
	subRepo *repository.SubscriptionRepository,
	planRepo *repository.PlanRepository,
) *SubscriptionUseCase {
	return &SubscriptionUseCase{
		DB:                     db,
		Log:                    logger,
		Validate:               validate,
		SubscriptionRepository: subRepo,
		PlanRepository:         planRepo,
	}
}

func (u *SubscriptionUseCase) GetCurrentSubscription(ctx context.Context, orgID string) (*model.SubscriptionResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	subscription := new(entity.Subscription)
	if err := u.SubscriptionRepository.FindActiveByOrganization(tx, subscription, orgID); err != nil {
		u.Log.Warnf("Failed to find subscription: %+v", err)
		return nil, fiber.ErrNotFound
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.Warnf("Failed to commit transaction: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.SubscriptionToResponse(subscription), nil
}

func (u *SubscriptionUseCase) Upgrade(ctx context.Context, request *model.UpgradeSubscriptionRequest) (*model.SubscriptionResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		u.Log.Warnf("Invalid request body: %+v", err)
		return nil, fiber.ErrBadRequest
	}

	// Get current subscription
	currentSub := new(entity.Subscription)
	if err := u.SubscriptionRepository.FindActiveByOrganization(tx, currentSub, request.OrganizationID); err != nil {
		u.Log.Warnf("Failed to find current subscription: %+v", err)
		return nil, fiber.ErrNotFound
	}

	// Get new plan
	newPlan := new(entity.Plan)
	if err := u.PlanRepository.FindById(tx, newPlan, request.PlanID); err != nil {
		u.Log.Warnf("Failed to find plan: %+v", err)
		return nil, fiber.ErrNotFound
	}

	// Cancel current subscription
	currentSub.Status = "cancelled"
	if err := u.SubscriptionRepository.Update(tx, currentSub); err != nil {
		u.Log.Warnf("Failed to cancel current subscription: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// Create new subscription
	newSub := &entity.Subscription{
		ID:                 uuid.New().String(),
		OrganizationID:     request.OrganizationID,
		PlanID:             request.PlanID,
		Status:             "active",
		CurrentPeriodStart: time.Now().UnixMilli(),
		CurrentPeriodEnd:   time.Now().AddDate(0, 1, 0).UnixMilli(), // 1 month
		CreatedAt:          time.Now().UnixMilli(),
		UpdatedAt:          time.Now().UnixMilli(),
	}
	newSub.Plan = *newPlan

	if err := u.SubscriptionRepository.Create(tx, newSub); err != nil {
		u.Log.Warnf("Failed to create new subscription: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.Warnf("Failed to commit transaction: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.SubscriptionToResponse(newSub), nil
}

func (u *SubscriptionUseCase) Cancel(ctx context.Context, request *model.CancelSubscriptionRequest) error {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		u.Log.Warnf("Invalid request body: %+v", err)
		return fiber.ErrBadRequest
	}

	subscription := new(entity.Subscription)
	if err := u.SubscriptionRepository.FindActiveByOrganization(tx, subscription, request.OrganizationID); err != nil {
		u.Log.Warnf("Failed to find subscription: %+v", err)
		return fiber.ErrNotFound
	}

	subscription.Status = "cancelled"
	if err := u.SubscriptionRepository.Update(tx, subscription); err != nil {
		u.Log.Warnf("Failed to cancel subscription: %+v", err)
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.Warnf("Failed to commit transaction: %+v", err)
		return fiber.ErrInternalServerError
	}

	return nil
}
