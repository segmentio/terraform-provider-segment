package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/segmentio/public-api-sdk-go/api"
)

type WarehouseMetadataState struct {
	ID          types.String             `tfsdk:"id"`
	Name        types.String             `tfsdk:"name"`
	Slug        types.String             `tfsdk:"slug"`
	Description types.String             `tfsdk:"description"`
	Logos       *LogosState              `tfsdk:"logos"`
	Options     []IntegrationOptionState `tfsdk:"options"`
}

func (w *WarehouseMetadataState) Fill(warehouseMetadata api.Metadata1) {
	w.ID = types.StringValue(warehouseMetadata.Id)
	w.Name = types.StringValue(warehouseMetadata.Name)
	w.Description = types.StringValue(warehouseMetadata.Description)
	w.Slug = types.StringValue(warehouseMetadata.Slug)
	w.Logos = getLogos(api.Logos(warehouseMetadata.Logos))
	w.Options = getOptions(warehouseMetadata.Options)
}