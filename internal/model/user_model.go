package model

type UserResponse struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	EmailVerified  bool   `json:"email_verified"`
	OrganizationID string `json:"organization_id,omitempty"`
	CreatedAt      int64  `json:"created_at"`
	UpdatedAt      int64  `json:"updated_at"`
}

type UpdateUserRequest struct {
	ID       string `json:"-" validate:"required,max=100"`
	Name     string `json:"name,omitempty" validate:"omitempty,max=100"`
	Password string `json:"password,omitempty" validate:"omitempty,min=8,max=100"`
}

type GetUserRequest struct {
	ID string `json:"id" validate:"required,max=100"`
}
