package entity

// Subscription is a struct that represents a subscription entity
type Subscription struct {
	ID                 string       `gorm:"column:id;primaryKey"`
	OrganizationID     string       `gorm:"column:organization_id"`
	PlanID             string       `gorm:"column:plan_id"`
	Status             string       `gorm:"column:status;default:active"`
	CurrentPeriodStart int64        `gorm:"column:current_period_start"`
	CurrentPeriodEnd   int64        `gorm:"column:current_period_end"`
	CreatedAt          int64        `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt          int64        `gorm:"column:updated_at;autoCreateTime:milli;autoUpdateTime:milli"`
	DeletedAt          *int64       `gorm:"column:deleted_at;index:idx_sub_deleted"`
	Organization       Organization `gorm:"foreignKey:organization_id;references:id"`
	Plan               Plan         `gorm:"foreignKey:plan_id;references:id"`
}

func (s *Subscription) TableName() string {
	return "subscriptions"
}
