package models

import (
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/segmentio/public-api-sdk-go/api"
)

type InsertFunctionInstanceState struct {
	ID            types.String         `tfsdk:"id"`
	FunctionID    types.String         `tfsdk:"function_id"`
	IntegrationID types.String         `tfsdk:"integration_id"`
	Name          types.String         `tfsdk:"name"`
	Enabled       types.Bool           `tfsdk:"enabled"`
	Settings      jsontypes.Normalized `tfsdk:"settings"`
}

func (i *InsertFunctionInstanceState) Fill(instance api.InsertFunctionInstanceAlpha) error {
	i.ID = types.StringValue(instance.Id)
	i.FunctionID = types.StringValue(instance.ClassId)
	i.IntegrationID = types.StringValue(instance.IntegrationId)
	i.Name = types.StringPointerValue(instance.Name)
	i.Enabled = types.BoolValue(instance.Enabled)
	settings, err := GetSettingsFromMap(instance.Settings)
	if err != nil {
		return err
	}
	i.Settings = settings

	return nil
}
