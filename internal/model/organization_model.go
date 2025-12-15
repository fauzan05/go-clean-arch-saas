package model

type OrganizationResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

type CreateOrganizationRequest struct {
	Name string `json:"name" validate:"required,max=200"`
	Slug string `json:"slug" validate:"required,max=200"`
}

type UpdateOrganizationRequest struct {
	ID   string `json:"-" validate:"required,max=100"`
	Name string `json:"name,omitempty" validate:"omitempty,max=200"`
}

type OrganizationMemberResponse struct {
	UserID   string       `json:"user_id"`
	User     UserResponse `json:"user,omitempty"`
	Role     string       `json:"role"`
	JoinedAt int64        `json:"joined_at"`
}

type ListOrganizationMembersRequest struct {
	OrganizationID string `json:"-" validate:"required,max=100"`
}

type RemoveOrganizationMemberRequest struct {
	OrganizationID string `json:"-" validate:"required,max=100"`
	UserID         string `json:"-" validate:"required,max=100"`
}
