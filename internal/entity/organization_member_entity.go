package entity

// OrganizationMember is a struct that represents an organization member entity
type OrganizationMember struct {
	OrganizationID string       `gorm:"column:organization_id;primaryKey"`
	UserID         string       `gorm:"column:user_id;primaryKey"`
	Role           string       `gorm:"column:role;default:member"`
	JoinedAt       int64        `gorm:"column:joined_at"`
	DeletedAt      *int64       `gorm:"column:deleted_at;index:idx_member_deleted"`
	Organization   Organization `gorm:"foreignKey:organization_id;references:id"`
	User           User         `gorm:"foreignKey:user_id;references:id"`
}

func (o *OrganizationMember) TableName() string {
	return "organization_members"
}
