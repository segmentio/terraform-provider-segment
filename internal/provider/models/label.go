package models

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/segmentio/public-api-sdk-go/api"
)

type LabelState struct {
	Key         types.String `tfsdk:"key"`
	Value       types.String `tfsdk:"value"`
	Description types.String `tfsdk:"description"`
}

func (l *LabelState) ToAPIValue() api.AllowedLabelBeta {
	return api.AllowedLabelBeta{
		Key:   l.Key.ValueString(),
		Value: l.Value.ValueString(),
	}
}

func (l *LabelState) Fill(label api.LabelV1) {
	l.Key = types.StringValue(label.Key)
	l.Value = types.StringValue(label.Value)
	l.Description = types.StringPointerValue(label.Description)
}

func LabelsPlanToAPILabels(ctx context.Context, labels types.Set) ([]api.AllowedLabelBeta, diag.Diagnostics) {
	apiLabels := []api.AllowedLabelBeta{}

	if !labels.IsNull() && !labels.IsUnknown() {
		stateLabels := []LabelState{}
		diags := labels.ElementsAs(ctx, &stateLabels, false)
		if diags.HasError() {
			return apiLabels, diags
		}
		for _, label := range stateLabels {
			apiLabels = append(apiLabels, label.ToAPIValue())
		}
	}

	return apiLabels, diag.Diagnostics{}
}

type LabelAssignmentState struct {
	ResourceID   types.String `tfsdk:"resource_id"`
	ResourceType types.String `tfsdk:"resource_type"`
	Labels       LabelState   `tfsdk:"labels"`
}

type LabelAssignmentPlan struct {
	ResourceID   types.String `tfsdk:"resource_id"`
	ResourceType types.String `tfsdk:"resource_type"`
	Labels       types.Set    `tfsdk:"labels"`
}
