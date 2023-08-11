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
	_ datasource.DataSource              = &warehouseDataSource{}
	_ datasource.DataSourceWithConfigure = &warehouseDataSource{}
)

func NewWarehouseDataSource() datasource.DataSource {
	return &warehouseDataSource{}
}

type warehouseDataSource struct {
	client      *api.APIClient
	authContext context.Context
}

type warehouseDataSourceModel struct {
	Id          types.String                      `tfsdk:"id"`
	Metadata    *warehouseMetadataDataSourceModel `tfsdk:"metadata"`
	WorkspaceId types.String                      `tfsdk:"workspace_id"`
	Enabled     types.Bool                        `tfsdk:"enabled"`
	// TODO: Add settings
}

func (d *warehouseDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_warehouse"
}

func (d *warehouseDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The warehouse",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The id of the Warehouse.",
			},
			"metadata": schema.SingleNestedAttribute{
				Description: "The metadata for the Warehouse.",
				Computed:    true,
				Attributes:  warehouseMetadataSchema(),
			},
			"workspace_id": schema.StringAttribute{
				Computed:    true,
				Description: "The id of the Workspace that owns this Warehouse.",
			},
			"enabled": schema.BoolAttribute{
				Computed:    true,
				Description: "When set to true, this Warehouse receives data.",
			},
			// TODO: Add settings
		},
	}
}

func (d *warehouseDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state warehouseDataSourceModel

	diags := req.Config.Get(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, _, err := d.client.WarehousesApi.GetWarehouse(d.authContext, state.Id.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Warehouse",
			err.Error(),
		)
		return
	}

	warehouse := response.Data.GetWarehouse()

	state.Id = types.StringValue(warehouse.Id)
	state.WorkspaceId = types.StringValue(warehouse.WorkspaceId)
	state.Enabled = types.BoolValue(warehouse.Enabled)
	state.Metadata = getWarehouseMetadata(warehouse.Metadata)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func getWarehouseMetadata(warehouseMetadata api.Metadata1) *warehouseMetadataDataSourceModel {
	var state warehouseMetadataDataSourceModel
	state.Id = types.StringValue(warehouseMetadata.Id)
	state.Name = types.StringValue(warehouseMetadata.Name)
	state.Description = types.StringValue(warehouseMetadata.Description)
	state.Slug = types.StringValue(warehouseMetadata.Slug)
	state.Logos = getLogos2(warehouseMetadata.Logos)
	state.Options = getOptions(warehouseMetadata.Options)
	return &state
}

func (d *warehouseDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
