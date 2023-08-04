package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/segmentio/public-api-sdk-go/api"
)

var (
	_ datasource.DataSource              = &destinationDataSource{}
	_ datasource.DataSourceWithConfigure = &destinationDataSource{}
)

func NewDestinationDataSource() datasource.DataSource {
	return &destinationDataSource{}
}

type destinationDataSource struct {
	client      *api.APIClient
	authContext context.Context
}

type destinationDataSourceModel struct {
	Id       types.String                        `tfsdk:"id"`
	Name     types.String                        `tfsdk:"name"`
	Enabled  types.Bool                          `tfsdk:"enabled"`
	Metadata *destinationMetadataDataSourceModel `tfsdk:"metadata"`
	SourceId types.String                        `tfsdk:"source_id"`
	//TODO: Add Settings
}

func (d *destinationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_destination"
}

func (d *destinationDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The destination",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of this instance of a Destination. Config API note: analogous to `name`.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of this instance of a Destination. Config API note: equal to `displayName`.",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether this instance of a Destination receives data.",
				Computed:    true,
			},
			"metadata": schema.SingleNestedAttribute{
				Description: "The metadata of the Destination of which this Destination is an instance of. For example, Google Analytics or Amplitude.",
				Computed:    true,
				Attributes:  destinationMetadataSchema(),
			},
			"source_id": schema.StringAttribute{
				Description: "The id of a Source connected to this instance of a Destination. Config API note: analogous to `parent`.",
				Computed:    true,
			},
			//TODO: Settings
		},
	}
}

func (d *destinationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state destinationDataSourceModel

	diags := req.Config.Get(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, _, err := d.client.DestinationsApi.GetDestination(d.authContext, state.Id.ValueString()).Execute()

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Destination",
			err.Error(),
		)
		return
	}

	var destination = response.Data.GetDestination()

	state.Id = types.StringValue(destination.Id)

	if destination.Name != nil {
		state.Name = types.StringValue(*destination.Name)
	}

	state.SourceId = types.StringValue(destination.SourceId)
	state.Enabled = types.BoolValue(destination.Enabled)
	state.Metadata = getDestinationMetadata(destination.Metadata)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func getDestinationMetadata(destinationMetadata api.Metadata) *destinationMetadataDataSourceModel {
	var state destinationMetadataDataSourceModel

	state.Id = types.StringValue(destinationMetadata.Id)
	state.Name = types.StringValue(destinationMetadata.Name)
	state.Description = types.StringValue(destinationMetadata.Description)
	state.Slug = types.StringValue(destinationMetadata.Slug)
	state.Logos = getLogosDestinationMetadata(destinationMetadata.Logos)
	state.Options = getOptions(destinationMetadata.Options)
	state.Actions = getActions(destinationMetadata.Actions)
	state.Categories = getCategories(destinationMetadata.Categories)
	state.Presets = getPresets(destinationMetadata.Presets)
	state.Contacts = getContacts(destinationMetadata.Contacts)
	state.PartnerOwned = getPartnerOwned(destinationMetadata.PartnerOwned)
	state.SupportedRegions = getSupportedRegions(destinationMetadata.SupportedRegions)
	state.RegionEndpoints = getRegionEndpoints(destinationMetadata.RegionEndpoints)
	state.Status = types.StringValue(destinationMetadata.Status)
	state.Website = types.StringValue(destinationMetadata.Website)
	state.Components = getComponents(destinationMetadata.Components)
	state.PreviousNames = getPreviousNames(destinationMetadata.PreviousNames)
	state.SupportedMethods = getSupportedMethods(destinationMetadata.SupportedMethods)
	state.SupportedFeatures = getSupportedFeatures(destinationMetadata.SupportedFeatures)
	state.SupportedPlatforms = getSupportedPlatforms(destinationMetadata.SupportedPlatforms)

	return &state
}

func (d *destinationDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	config, ok := req.ProviderData.(*ClientInfo)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected ClientInfo, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = config.client
	d.authContext = config.authContext
}
