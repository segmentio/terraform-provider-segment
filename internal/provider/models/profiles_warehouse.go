package models

import (
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/segmentio/public-api-sdk-go/api"
)

type ProfilesWarehouseState struct {
	ID         types.String         `tfsdk:"id"`
	SpaceID    types.String         `tfsdk:"space_id"`
	MetadataID types.String         `tfsdk:"metadata_id"`
	Name       types.String         `tfsdk:"name"`
	Enabled    types.Bool           `tfsdk:"enabled"`
	SchemaName types.String         `tfsdk:"schema_name"`
	Settings   jsontypes.Normalized `tfsdk:"settings"`
}

func (w *ProfilesWarehouseState) Fill(warehouse api.ProfilesWarehouse1) error {
	w.ID = types.StringValue(warehouse.Id)
	w.SpaceID = types.StringValue(warehouse.SpaceId)
	w.MetadataID = types.StringValue(warehouse.Metadata.Id)
	if warehouse.Settings.IsSet() {
		name := warehouse.Settings.Get().Get()["name"]
		if name != nil {
			stringName, ok := name.(string)
			if ok {
				w.Name = types.StringValue(stringName)
			}
		}
	}
	w.Enabled = types.BoolValue(warehouse.Enabled)
	w.SchemaName = types.StringPointerValue(warehouse.SchemaName)
	settings, err := GetSettings(warehouse.Settings)
	if err != nil {
		return err
	}
	w.Settings = settings

	return nil
}
