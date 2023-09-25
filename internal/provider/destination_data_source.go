package provider

import (
	"context"
	"fmt"

	"github.com/segmentio/terraform-provider-segment/internal/provider/models"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

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

func (d *destinationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_destination"
}

func (d *destinationDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The destination",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of this instance of a Destination. Config API note: analogous to `name`.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of this instance of a Destination. Config API note: equal to `displayName`.",
				Computed:    true,
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
			"settings": schema.StringAttribute{
				Computed:    true,
				Description: "The settings associated with the Destination.",
				CustomType:  jsontypes.NormalizedType{},
			},
		},
	}
}

func (d *destinationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state models.DestinationState

	diags := req.Config.Get(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, body, err := d.client.DestinationsApi.GetDestination(d.authContext, state.ID.ValueString()).Execute()
	defer body.Body.Close()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Destination",
			getError(err, body),
		)

		return
	}

	destination := response.Data.GetDestination()

	err = state.Fill(&destination)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Destination",
			err.Error(),
		)

		return
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
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
