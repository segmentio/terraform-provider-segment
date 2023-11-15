package models

import (
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/segmentio/public-api-sdk-go/api"
)

type DestinationState struct {
	ID       types.String              `tfsdk:"id"`
	Name     types.String              `tfsdk:"name"`
	Enabled  types.Bool                `tfsdk:"enabled"`
	Metadata *DestinationMetadataState `tfsdk:"metadata"`
	SourceID types.String              `tfsdk:"source_id"`
	Settings jsontypes.Normalized      `tfsdk:"settings"`
}

type DestinationPlan struct {
	ID       types.String         `tfsdk:"id"`
	Name     types.String         `tfsdk:"name"`
	Enabled  types.Bool           `tfsdk:"enabled"`
	Metadata types.Object         `tfsdk:"metadata"`
	SourceID types.String         `tfsdk:"source_id"`
	Settings jsontypes.Normalized `tfsdk:"settings"`
}

func (d *DestinationState) Fill(destination *api.DestinationV1) error {
	var destinationMetadata DestinationMetadataState
	err := destinationMetadata.Fill(destination.Metadata)
	if err != nil {
		return err
	}

	d.ID = types.StringValue(destination.Id)
	d.Name = types.StringPointerValue(destination.Name)
	d.SourceID = types.StringValue(destination.SourceId)
	d.Enabled = types.BoolValue(destination.Enabled)
	d.Metadata = &destinationMetadata
	settings, err := GetSettingsFromMap(destination.Settings)
	if err != nil {
		return err
	}
	d.Settings = settings

	return nil
}

func GetSettingsFromMap(settings map[string]interface{}) (jsontypes.Normalized, error) {
	jsonSettingsString, err := json.Marshal(settings)
	if err != nil {
		return jsontypes.NewNormalizedNull(), err
	}

	return jsontypes.NewNormalizedValue(string(jsonSettingsString)), nil
}
