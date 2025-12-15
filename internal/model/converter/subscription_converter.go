package converter

import (
	"encoding/json"
	"go-clean-arch-saas/internal/entity"
	"go-clean-arch-saas/internal/model"
)

func PlanToResponse(plan *entity.Plan) *model.PlanResponse {
	var features map[string]interface{}
	var limits map[string]interface{}

	if plan.Features != "" {
		json.Unmarshal([]byte(plan.Features), &features)
	}

	if plan.Limits != "" {
		json.Unmarshal([]byte(plan.Limits), &limits)
	}

	return &model.PlanResponse{
		ID:            plan.ID,
		Name:          plan.Name,
		Slug:          plan.Slug,
		Price:         plan.Price,
		BillingPeriod: plan.BillingPeriod,
		Features:      features,
		Limits:        limits,
		IsActive:      plan.IsActive,
		CreatedAt:     plan.CreatedAt,
		UpdatedAt:     plan.UpdatedAt,
	}
}

func SubscriptionToResponse(subscription *entity.Subscription) *model.SubscriptionResponse {
	response := &model.SubscriptionResponse{
		ID:                 subscription.ID,
		OrganizationID:     subscription.OrganizationID,
		Status:             subscription.Status,
		CurrentPeriodStart: subscription.CurrentPeriodStart,
		CurrentPeriodEnd:   subscription.CurrentPeriodEnd,
		CreatedAt:          subscription.CreatedAt,
		UpdatedAt:          subscription.UpdatedAt,
	}

	if subscription.Plan.ID != "" {
		response.Plan = *PlanToResponse(&subscription.Plan)
	}

	return response
}
