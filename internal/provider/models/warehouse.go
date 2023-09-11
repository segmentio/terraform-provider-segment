package models

import (
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/segmentio/public-api-sdk-go/api"
)

type WarehouseState struct {
	ID          types.String            `tfsdk:"id"`
	Metadata    *WarehouseMetadataState `tfsdk:"metadata"`
	Name        types.String            `tfsdk:"name"`
	WorkspaceId types.String            `tfsdk:"workspace_id"`
	Enabled     types.Bool              `tfsdk:"enabled"`
	Settings    jsontypes.Normalized    `tfsdk:"settings"`
}

type WarehousePlan struct {
	ID          types.String         `tfsdk:"id"`
	Metadata    types.Object         `tfsdk:"metadata"`
	Name        types.String         `tfsdk:"name"`
	WorkspaceId types.String         `tfsdk:"workspace_id"`
	Enabled     types.Bool           `tfsdk:"enabled"`
	Settings    jsontypes.Normalized `tfsdk:"settings"`
}

func (w *WarehouseState) Fill(warehouse api.Warehouse) error {
	warehouseMetadata := WarehouseMetadataState{}
	err := warehouseMetadata.Fill(warehouse.Metadata)
	if err != nil {
		return err
	}

	w.ID = types.StringValue(warehouse.Id)
	w.WorkspaceId = types.StringValue(warehouse.WorkspaceId)
	w.Enabled = types.BoolValue(warehouse.Enabled)
	w.Metadata = &warehouseMetadata
	settings, err := w.getSettings(warehouse.Settings)
	if err != nil {
		return err
	}
	w.Settings = settings
	if warehouse.Settings.IsSet() {
		name := warehouse.Settings.Get().Get()["name"]
		if name != nil {
			stringName, ok := name.(string)
			if ok {
				w.Name = types.StringValue(stringName)
			}
		}
	}

	return nil
}

func (s *WarehouseState) getSettings(settings api.NullableModelMap) (jsontypes.Normalized, error) {
	if !settings.IsSet() {
		return jsontypes.NewNormalizedNull(), nil
	}

	unwrappedSettings := settings.Get().Get()

	// We remove "name" from the returned settings to surface it as a top level attribute
	settingsWithoutName := make(map[string]interface{})
	for k, v := range unwrappedSettings {
		if k != "name" {
			settingsWithoutName[k] = v
		}
	}
	jsonSettingsString, err := json.Marshal(settingsWithoutName)
	if err != nil {
		return jsontypes.NewNormalizedNull(), err
	}

	return jsontypes.NewNormalizedValue(string(jsonSettingsString)), nil
}
