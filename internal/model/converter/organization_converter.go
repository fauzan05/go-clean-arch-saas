package converter

import (
	"go-clean-arch-saas/internal/entity"
	"go-clean-arch-saas/internal/model"
)

func OrganizationToResponse(org *entity.Organization) *model.OrganizationResponse {
	return &model.OrganizationResponse{
		ID:        org.ID,
		Name:      org.Name,
		Slug:      org.Slug,
		CreatedAt: org.CreatedAt,
		UpdatedAt: org.UpdatedAt,
	}
}

func OrganizationMemberToResponse(member *entity.OrganizationMember) *model.OrganizationMemberResponse {
	response := &model.OrganizationMemberResponse{
		UserID:   member.UserID,
		Role:     member.Role,
		JoinedAt: member.JoinedAt,
	}

	if member.User.ID != "" {
		response.User = *UserToResponse(&member.User)
	}

	return response
}
