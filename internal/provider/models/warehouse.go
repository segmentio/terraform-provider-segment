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
	WorkspaceId types.String            `tfsdk:"workspace_id"`
	Enabled     types.Bool              `tfsdk:"enabled"`
	Settings    jsontypes.Normalized    `tfsdk:"settings"`
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

	return nil
}

func (s *WarehouseState) getSettings(settings api.NullableModelMap) (jsontypes.Normalized, error) {
	if !settings.IsSet() {
		return jsontypes.NewNormalizedNull(), nil
	}

	jsonSettingsString, err := json.Marshal(settings.Get().Get())
	if err != nil {
		return jsontypes.NewNormalizedNull(), err
	}

	return jsontypes.NewNormalizedValue(string(jsonSettingsString)), nil
}
