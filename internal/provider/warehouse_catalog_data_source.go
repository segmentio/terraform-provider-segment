package provider

import (
	"context"
	"fmt"
	"github.com/segmentio/public-api-sdk-go/api"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &warehouseCatalogDataSource{}
	_ datasource.DataSourceWithConfigure = &warehouseCatalogDataSource{}
)

// NewWarehouseCatalogDataSource is a helper function to simplify the provider implementation.
func NewWarehouseCatalogDataSource() datasource.DataSource {
	return &warehouseCatalogDataSource{}
}

// warehouseCatalogDataSource is the data source implementation.
type warehouseCatalogDataSource struct {
	client      *api.APIClient
	authContext context.Context
}

// Metadata returns the data source type name.
func (d *warehouseCatalogDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_warehouse_catalog"
}

// Read refreshes the Terraform state with the latest data.
func (d *warehouseCatalogDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
}

// Schema defines the schema for the data source.
func (d *warehouseCatalogDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The warehouse catalog",
		Attributes: map[string]schema.Attribute{
			"warehouse_metadatas": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: warehouseMetadataSchema(),
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *warehouseCatalogDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	clientInfo, ok := req.ProviderData.(*ClientInfo)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *ClientInfo, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = clientInfo.client
	d.authContext = clientInfo.authContext
}
