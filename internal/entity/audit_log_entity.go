package entity

// AuditLog is a struct that represents an audit log entity
type AuditLog struct {
	ID             string       `gorm:"column:id;primaryKey"`
	UserID         string       `gorm:"column:user_id"`
	OrganizationID string       `gorm:"column:organization_id"`
	Action         string       `gorm:"column:action"`
	Resource       string       `gorm:"column:resource"`
	ResourceID     string       `gorm:"column:resource_id"`
	Details        string       `gorm:"column:details;type:json"`
	IPAddress      string       `gorm:"column:ip_address"`
	UserAgent      string       `gorm:"column:user_agent"`
	CreatedAt      int64        `gorm:"column:created_at;autoCreateTime:milli"`
	DeletedAt      *int64       `gorm:"column:deleted_at;index:idx_audit_deleted"`
	User           User         `gorm:"foreignKey:user_id;references:id"`
	Organization   Organization `gorm:"foreignKey:organization_id;references:id"`
}

func (a *AuditLog) TableName() string {
	return "audit_logs"
}
