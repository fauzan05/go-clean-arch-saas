package entity

// Plan is a struct that represents a subscription plan entity
type Plan struct {
	ID            string  `gorm:"column:id;primaryKey"`
	Name          string  `gorm:"column:name"`
	Slug          string  `gorm:"column:slug;unique"`
	Price         float64 `gorm:"column:price"`
	BillingPeriod string  `gorm:"column:billing_period"`
	Features      string  `gorm:"column:features;type:json"`
	Limits        string  `gorm:"column:limits;type:json"`
	IsActive      bool    `gorm:"column:is_active;default:true"`
	CreatedAt     int64   `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt     int64   `gorm:"column:updated_at;autoCreateTime:milli;autoUpdateTime:milli"`
	DeletedAt     *int64  `gorm:"column:deleted_at;index:idx_plan_deleted"`
}

func (p *Plan) TableName() string {
	return "plans"
}
