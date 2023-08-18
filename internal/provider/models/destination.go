package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/segmentio/public-api-sdk-go/api"
)

type DestinationState struct {
	ID       types.String              `tfsdk:"id"`
	Name     types.String              `tfsdk:"name"`
	Enabled  types.Bool                `tfsdk:"enabled"`
	Metadata *DestinationMetadataState `tfsdk:"metadata"`
	SourceId types.String              `tfsdk:"source_id"`
}

type DestinationPlan struct {
	ID       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	Enabled  types.Bool   `tfsdk:"enabled"`
	Metadata types.Object `tfsdk:"metadata"`
	SourceId types.String `tfsdk:"source_id"`
}

func (d *DestinationState) Fill(destination *api.Destination) {
	var destinationMetadata DestinationMetadataState
	destinationMetadata.Fill(api.DestinationMetadata(destination.Metadata))

	d.ID = types.StringValue(destination.Id)
	d.Name = types.StringPointerValue(destination.Name)
	d.SourceId = types.StringValue(destination.SourceId)
	d.Enabled = types.BoolValue(destination.Enabled)
	d.Metadata = &destinationMetadata
}
