package entity

// Organization is a struct that represents an organization entity
type Organization struct {
	ID        string               `gorm:"column:id;primaryKey"`
	Name      string               `gorm:"column:name"`
	Slug      string               `gorm:"column:slug;unique"`
	CreatedAt int64                `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt int64                `gorm:"column:updated_at;autoCreateTime:milli;autoUpdateTime:milli"`
	Members   []OrganizationMember `gorm:"foreignKey:organization_id;references:id"`
	Users     []User               `gorm:"foreignKey:organization_id;references:id"`
}

func (o *Organization) TableName() string {
	return "organizations"
}
