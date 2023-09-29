package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/segmentio/public-api-sdk-go/api"
)

type LabelResourceState struct {
	ID          types.String `tfsdk:"id"`
	Key         types.String `tfsdk:"key"`
	Value       types.String `tfsdk:"value"`
	Description types.String `tfsdk:"description"`
}

type LabelState struct {
	Description types.String `tfsdk:"description"`
	Key         types.String `tfsdk:"key"`
	Value       types.String `tfsdk:"value"`
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

	if label.Description != nil {
		l.Description = types.StringValue(*label.Description)
	}
}
