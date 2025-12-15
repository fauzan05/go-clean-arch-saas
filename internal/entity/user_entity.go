package entity

// User is a struct that represents a user entity
type User struct {
	ID                    string       `gorm:"column:id;primaryKey"`
	Name                  string       `gorm:"column:name"`
	Email                 string       `gorm:"column:email;unique"`
	Password              string       `gorm:"column:password"`
	RefreshToken          string       `gorm:"column:refresh_token"`
	RefreshTokenExpiresAt int64        `gorm:"column:refresh_token_expires_at"`
	OrganizationID        string       `gorm:"column:organization_id"`
	CreatedAt             int64        `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt             int64        `gorm:"column:updated_at;autoCreateTime:milli;autoUpdateTime:milli"`
	DeletedAt             *int64       `gorm:"column:deleted_at;index:idx_users_deleted"`
	Organization          Organization `gorm:"foreignKey:organization_id;references:id"`
}

func (u *User) TableName() string {
	return "users"
}
