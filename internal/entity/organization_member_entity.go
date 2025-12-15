package entity

// Organization role constants for organization-level access control
const (
	OrgRoleOwner  = "owner"  // Organization owner (1 per org, full control)
	OrgRoleAdmin  = "admin"  // Organization admin (can manage members & settings)
	OrgRoleMember = "member" // Regular member (limited permissions)
)

// OrganizationMember is a struct that represents an organization member entity
type OrganizationMember struct {
	OrganizationID string       `gorm:"column:organization_id;primaryKey"`
	UserID         string       `gorm:"column:user_id;primaryKey"`
	Role           string       `gorm:"column:role;default:member;index:idx_member_role"`
	JoinedAt       int64        `gorm:"column:joined_at"`
	DeletedAt      *int64       `gorm:"column:deleted_at;index:idx_member_deleted"`
	Organization   Organization `gorm:"foreignKey:organization_id;references:id"`
	User           User         `gorm:"foreignKey:user_id;references:id"`
}

func (o *OrganizationMember) TableName() string {
	return "organization_members"
}

// IsOwner checks if member is organization owner
func (o *OrganizationMember) IsOwner() bool {
	return o.Role == OrgRoleOwner
}

// IsAdmin checks if member is organization admin or owner
func (o *OrganizationMember) IsAdmin() bool {
	return o.Role == OrgRoleAdmin || o.Role == OrgRoleOwner
}

// IsMember checks if member has any valid role (owner, admin, or member)
func (o *OrganizationMember) IsMember() bool {
	return o.Role == OrgRoleMember || o.IsAdmin()
}
