package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/segmentio/public-api-sdk-go/api"
)

type TrackingPlanState struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Slug        types.String `tfsdk:"slug"`
	Description types.String `tfsdk:"description"`
	Type        types.String `tfsdk:"type"`
	UpdatedAt   types.String `tfsdk:"updated_at"`
	CreatedAt   types.String `tfsdk:"created_at"`
}

func (t *TrackingPlanState) Fill(trackingPlan api.TrackingPlan) error {
	t.ID = types.StringValue(trackingPlan.Id)
	if trackingPlan.Name != nil {
		t.Name = types.StringValue(*trackingPlan.Name)
	}
	if trackingPlan.Slug != nil {
		t.Slug = types.StringValue(*trackingPlan.Slug)
	}
	if trackingPlan.Description != nil {
		t.Description = types.StringValue(*trackingPlan.Description)
	}
	t.Type = types.StringValue(trackingPlan.Type)
	if trackingPlan.UpdatedAt != nil {
		t.UpdatedAt = types.StringValue(*trackingPlan.UpdatedAt)
	}
	if trackingPlan.CreatedAt != nil {
		t.CreatedAt = types.StringValue(*trackingPlan.CreatedAt)
	}

	return nil
}
