package entity

// System role constants for platform-level access control
const (
	SystemRoleUser       = "user"        // Regular tenant user (default)
	SystemRoleSupport    = "support"     // Customer support access
	SystemRoleAdmin      = "admin"       // Platform admin access
	SystemRoleSuperAdmin = "super_admin" // Platform owner (full access)
)

// User is a struct that represents a user entity
type User struct {
	ID                    string       `gorm:"column:id;primaryKey"`
	Name                  string       `gorm:"column:name"`
	Email                 string       `gorm:"column:email;unique"`
	Password              string       `gorm:"column:password"`
	SystemRole            string       `gorm:"column:system_role;default:user;index:idx_users_system_role"`
	EmailVerified         bool         `gorm:"column:email_verified;default:0"`
	EmailVerifiedAt       *int64       `gorm:"column:email_verified_at"`
	VerificationToken     *string      `gorm:"column:verification_token;index:idx_users_verification_token"`
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

// IsSystemAdmin checks if user has platform admin access
func (u *User) IsSystemAdmin() bool {
	return u.SystemRole == SystemRoleAdmin || u.SystemRole == SystemRoleSuperAdmin
}

// IsSuperAdmin checks if user is platform owner
func (u *User) IsSuperAdmin() bool {
	return u.SystemRole == SystemRoleSuperAdmin
}

// IsSupport checks if user has support access
func (u *User) IsSupport() bool {
	return u.SystemRole == SystemRoleSupport || u.IsSystemAdmin()
}
