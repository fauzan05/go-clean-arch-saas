package usecase

import (
	"context"
	"go-clean-arch-saas/internal/entity"
	"go-clean-arch-saas/internal/model"
	"go-clean-arch-saas/internal/model/converter"
	"go-clean-arch-saas/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type OrganizationUseCase struct {
	DB                           *gorm.DB
	Log                          *logrus.Logger
	Validate                     *validator.Validate
	OrganizationRepository       *repository.OrganizationRepository
	OrganizationMemberRepository *repository.OrganizationMemberRepository
}

func NewOrganizationUseCase(
	db *gorm.DB,
	logger *logrus.Logger,
	validate *validator.Validate,
	orgRepo *repository.OrganizationRepository,
	orgMemberRepo *repository.OrganizationMemberRepository,
) *OrganizationUseCase {
	return &OrganizationUseCase{
		DB:                           db,
		Log:                          logger,
		Validate:                     validate,
		OrganizationRepository:       orgRepo,
		OrganizationMemberRepository: orgMemberRepo,
	}
}

func (u *OrganizationUseCase) GetByID(ctx context.Context, orgID string) (*model.OrganizationResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	org := new(entity.Organization)
	if err := u.OrganizationRepository.FindById(tx, org, orgID); err != nil {
		u.Log.Warnf("Failed to find organization: %+v", err)
		return nil, fiber.ErrNotFound
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.Warnf("Failed to commit transaction: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.OrganizationToResponse(org), nil
}

func (u *OrganizationUseCase) Update(ctx context.Context, request *model.UpdateOrganizationRequest) (*model.OrganizationResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		u.Log.Warnf("Invalid request body: %+v", err)
		return nil, fiber.ErrBadRequest
	}

	org := new(entity.Organization)
	if err := u.OrganizationRepository.FindById(tx, org, request.ID); err != nil {
		u.Log.Warnf("Failed to find organization: %+v", err)
		return nil, fiber.ErrNotFound
	}

	if request.Name != "" {
		org.Name = request.Name
	}

	if err := u.OrganizationRepository.Update(tx, org); err != nil {
		u.Log.Warnf("Failed to update organization: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.Warnf("Failed to commit transaction: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.OrganizationToResponse(org), nil
}

func (u *OrganizationUseCase) ListMembers(ctx context.Context, orgID string) ([]model.OrganizationMemberResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	members, err := u.OrganizationMemberRepository.ListByOrganization(tx, orgID)
	if err != nil {
		u.Log.Warnf("Failed to list organization members: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.Warnf("Failed to commit transaction: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	var responses []model.OrganizationMemberResponse
	for _, member := range members {
		responses = append(responses, *converter.OrganizationMemberToResponse(&member))
	}

	return responses, nil
}

func (u *OrganizationUseCase) RemoveMember(ctx context.Context, request *model.RemoveOrganizationMemberRequest) error {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		u.Log.Warnf("Invalid request body: %+v", err)
		return fiber.ErrBadRequest
	}

	// Check if member exists
	member := new(entity.OrganizationMember)
	if err := u.OrganizationMemberRepository.FindByOrgAndUser(tx, member, request.OrganizationID, request.UserID); err != nil {
		u.Log.Warnf("Member not found: %+v", err)
		return fiber.ErrNotFound
	}

	// Don't allow removing owner
	if member.Role == "owner" {
		u.Log.Warnf("Cannot remove owner from organization")
		return fiber.NewError(fiber.StatusForbidden, "Cannot remove owner from organization")
	}

	if err := u.OrganizationMemberRepository.DeleteByOrgAndUser(tx, request.OrganizationID, request.UserID); err != nil {
		u.Log.Warnf("Failed to remove member: %+v", err)
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.Warnf("Failed to commit transaction: %+v", err)
		return fiber.ErrInternalServerError
	}

	return nil
}
