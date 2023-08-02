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
	_ datasource.DataSource              = &sourceCatalogDataSource{}
	_ datasource.DataSourceWithConfigure = &sourceCatalogDataSource{}
)

// NewSourceCatalogDataSource is a helper function to simplify the provider implementation.
func NewSourceCatalogDataSource() datasource.DataSource {
	return &sourceCatalogDataSource{}
}

// sourceCatalogDataSource is the data source implementation.
type sourceCatalogDataSource struct {
	client      *api.APIClient
	authContext context.Context
}

// Metadata returns the data source type name.
func (d *sourceCatalogDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sourceCatalog"
}

// Read refreshes the Terraform state with the latest data.
func (d *sourceCatalogDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
}

// Schema defines the schema for the data source.
func (d *sourceCatalogDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The source catalog",
		Attributes: map[string]schema.Attribute{
			"source_metadatas": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: sourceMetadataSchema(),
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *sourceCatalogDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
