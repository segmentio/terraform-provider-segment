package models

import (
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/segmentio/public-api-sdk-go/api"
)

type DestinationSubscriptionState struct {
	ID            types.String         `tfsdk:"id"`
	DestinationID types.String         `tfsdk:"destination_id"`
	Name          types.String         `tfsdk:"name"`
	Enabled       types.Bool           `tfsdk:"enabled"`
	ActionID      types.String         `tfsdk:"action_id"`
	ActionSlug    types.String         `tfsdk:"action_slug"`
	Trigger       types.String         `tfsdk:"trigger"`
	ModelID       types.String         `tfsdk:"model_id"`
	Settings      jsontypes.Normalized `tfsdk:"settings"`
}

func (d *DestinationSubscriptionState) Fill(subscription api.DestinationSubscription) error {
	d.ID = types.StringValue(subscription.Id)
	d.DestinationID = types.StringValue(subscription.DestinationId)
	d.Name = types.StringValue(subscription.Name)
	d.Enabled = types.BoolValue(subscription.Enabled)
	d.ActionID = types.StringValue(subscription.ActionId)
	d.ActionSlug = types.StringValue(subscription.ActionSlug)
	d.Trigger = types.StringValue(subscription.Trigger)
	if subscription.ModelId != nil && *subscription.ModelId == "" {
		subscription.ModelId = nil
	}
	d.ModelID = types.StringPointerValue(subscription.ModelId)
	settings, err := GetSettings(subscription.Settings)
	if err != nil {
		return err
	}
	d.Settings = settings

	return nil
}
