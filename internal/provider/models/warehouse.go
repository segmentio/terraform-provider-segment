package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/segmentio/public-api-sdk-go/api"
)

type WarehouseState struct {
	ID          types.String            `tfsdk:"id"`
	Metadata    *WarehouseMetadataState `tfsdk:"metadata"`
	WorkspaceId types.String            `tfsdk:"workspace_id"`
	Enabled     types.Bool              `tfsdk:"enabled"`
	// TODO: Add settings
}

func (w *WarehouseState) Fill(warehouse api.Warehouse) {
	warehouseMetadata := WarehouseMetadataState{}
	warehouseMetadata.Fill(warehouse.Metadata)

	w.ID = types.StringValue(warehouse.Id)
	w.WorkspaceId = types.StringValue(warehouse.WorkspaceId)
	w.Enabled = types.BoolValue(warehouse.Enabled)
	w.Metadata = &warehouseMetadata
}
