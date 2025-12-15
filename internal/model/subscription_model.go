package model

type PlanResponse struct {
	ID            string                 `json:"id"`
	Name          string                 `json:"name"`
	Slug          string                 `json:"slug"`
	Price         float64                `json:"price"`
	BillingPeriod string                 `json:"billing_period"`
	Features      map[string]interface{} `json:"features"`
	Limits        map[string]interface{} `json:"limits"`
	IsActive      bool                   `json:"is_active"`
	CreatedAt     int64                  `json:"created_at"`
	UpdatedAt     int64                  `json:"updated_at"`
}

type SubscriptionResponse struct {
	ID                 string       `json:"id"`
	OrganizationID     string       `json:"organization_id"`
	Plan               PlanResponse `json:"plan"`
	Status             string       `json:"status"`
	CurrentPeriodStart int64        `json:"current_period_start"`
	CurrentPeriodEnd   int64        `json:"current_period_end"`
	CreatedAt          int64        `json:"created_at"`
	UpdatedAt          int64        `json:"updated_at"`
}

type GetCurrentSubscriptionRequest struct {
	OrganizationID string `json:"-" validate:"required,max=100"`
}

type UpgradeSubscriptionRequest struct {
	OrganizationID string `json:"-" validate:"required,max=100"`
	PlanID         string `json:"plan_id" validate:"required,max=100"`
}

type CancelSubscriptionRequest struct {
	OrganizationID string `json:"-" validate:"required,max=100"`
}
