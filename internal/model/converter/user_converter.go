package converter

import (
	"go-clean-arch-saas/internal/entity"
	"go-clean-arch-saas/internal/model"
)

func UserToResponse(user *entity.User) *model.UserResponse {
	return &model.UserResponse{
		ID:             user.ID,
		Name:           user.Name,
		Email:          user.Email,
		EmailVerified:  user.EmailVerified,
		OrganizationID: user.OrganizationID,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
	}
}
