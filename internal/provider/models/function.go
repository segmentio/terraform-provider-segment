package models

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/segmentio/public-api-sdk-go/api"
)

type FunctionState struct {
	ID                types.String           `tfsdk:"id"`
	Code              types.String           `tfsdk:"code"`
	DisplayName       types.String           `tfsdk:"display_name"`
	LogoURL           types.String           `tfsdk:"logo_url"`
	ResourceType      types.String           `tfsdk:"resource_type"`
	Description       types.String           `tfsdk:"description"`
	PreviewWebhookURL types.String           `tfsdk:"preview_webhook_url"`
	CatalogID         types.String           `tfsdk:"catalog_id"`
	Settings          []FunctionSettingState `tfsdk:"settings"`
}

type FunctionSettingState struct {
	Name        types.String `tfsdk:"name"`
	Label       types.String `tfsdk:"label"`
	Description types.String `tfsdk:"description"`
	Type        types.String `tfsdk:"type"`
	Required    types.Bool   `tfsdk:"required"`
	Sensitive   types.Bool   `tfsdk:"sensitive"`
}

type FunctionPlan struct {
	ID                types.String `tfsdk:"id"`
	Code              types.String `tfsdk:"code"`
	DisplayName       types.String `tfsdk:"display_name"`
	LogoURL           types.String `tfsdk:"logo_url"`
	ResourceType      types.String `tfsdk:"resource_type"`
	Description       types.String `tfsdk:"description"`
	PreviewWebhookURL types.String `tfsdk:"preview_webhook_url"`
	CatalogID         types.String `tfsdk:"catalog_id"`
	Settings          types.Set    `tfsdk:"settings"`
}

func (f *FunctionState) Fill(function api.FunctionV1) {
	f.ID = types.StringPointerValue(function.Id)
	f.Code = types.StringPointerValue(function.Code)
	f.DisplayName = types.StringPointerValue(function.DisplayName)
	f.LogoURL = types.StringPointerValue(function.LogoUrl)
	f.ResourceType = types.StringPointerValue(function.ResourceType)
	f.Description = types.StringPointerValue(function.Description)
	f.PreviewWebhookURL = types.StringPointerValue(function.PreviewWebhookUrl)
	f.CatalogID = types.StringPointerValue(function.CatalogId)

	f.Settings = []FunctionSettingState{}
	for _, s := range function.Settings {
		var setting FunctionSettingState
		setting.Fill(s)
		f.Settings = append(f.Settings, setting)
	}
}

func (f *FunctionSettingState) Fill(setting api.FunctionSettingV1) {
	f.Name = types.StringValue(setting.Name)
	f.Label = types.StringValue(setting.Label)
	f.Description = types.StringValue(setting.Description)
	f.Type = types.StringValue(setting.Type)
	f.Required = types.BoolValue(setting.Required)
	f.Sensitive = types.BoolValue(setting.Sensitive)
}

func (f *FunctionSettingState) ToAPIValue() api.FunctionSettingV1 {
	return api.FunctionSettingV1{
		Name:        f.Name.ValueString(),
		Label:       f.Label.ValueString(),
		Description: f.Description.ValueString(),
		Type:        f.Type.ValueString(),
		Required:    f.Required.ValueBool(),
		Sensitive:   f.Sensitive.ValueBool(),
	}
}

func GetFunctionSettingAPIValueFromPlan(ctx context.Context, settings types.Set) ([]api.FunctionSettingV1, diag.Diagnostics) {
	var outDiags diag.Diagnostics

	var settingsState = []FunctionSettingState{}
	diags := settings.ElementsAs(ctx, &settingsState, false)
	outDiags.Append(diags...)
	if outDiags.HasError() {
		return []api.FunctionSettingV1{}, outDiags
	}

	apiSettings := []api.FunctionSettingV1{}

	for _, setting := range settingsState {
		s := setting.ToAPIValue()
		outDiags.Append(diags...)
		if outDiags.HasError() {
			return []api.FunctionSettingV1{}, outDiags
		}
		apiSettings = append(apiSettings, s)
	}

	return apiSettings, outDiags
}
