package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/segmentio/public-api-sdk-go/api"
)

type DestinationStateModel struct {
	Id       types.String                   `tfsdk:"id"`
	Name     types.String                   `tfsdk:"name"`
	Enabled  types.Bool                     `tfsdk:"enabled"`
	Metadata *DestinationMetadataStateModel `tfsdk:"metadata"`
	SourceId types.String                   `tfsdk:"source_id"`
}

type DestinationPlanModel struct {
	Id       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	Enabled  types.Bool   `tfsdk:"enabled"`
	Metadata types.Object `tfsdk:"metadata"`
	SourceId types.String `tfsdk:"source_id"`
}

func (d *DestinationStateModel) Fill(destination *api.Destination) {
	var destinationMetadata DestinationMetadataStateModel
	destinationMetadata.Fill(api.DestinationMetadata(destination.Metadata))

	d.Id = types.StringValue(destination.Id)
	d.Name = types.StringPointerValue(destination.Name)
	d.SourceId = types.StringValue(destination.SourceId)
	d.Enabled = types.BoolValue(destination.Enabled)
	d.Metadata = &destinationMetadata
}
