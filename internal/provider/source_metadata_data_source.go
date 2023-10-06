package provider

import (
	"context"
	"fmt"

	"github.com/segmentio/terraform-provider-segment/internal/provider/models"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/segmentio/public-api-sdk-go/api"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &sourceMetadataDataSource{}
	_ datasource.DataSourceWithConfigure = &sourceMetadataDataSource{}
)

// NewSourceMetadataDataSource is a helper function to simplify the provider implementation.
func NewSourceMetadataDataSource() datasource.DataSource {
	return &sourceMetadataDataSource{}
}

// sourceMetadataDataSource is the data source implementation.
type sourceMetadataDataSource struct {
	client      *api.APIClient
	authContext context.Context
}

// Metadata returns the data source type name.
func (d *sourceMetadataDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_source_metadata"
}

// Read refreshes the Terraform state with the latest data.
func (d *sourceMetadataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state models.SourceMetadataState

	diags := req.Config.Get(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	if id == "" {
		resp.Diagnostics.AddError("Unable to read Source Metadata", "ID is empty")

		return
	}

	response, body, err := d.client.CatalogApi.GetSourceMetadata(d.authContext, id).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Source metadata",
			getError(err, body),
		)

		return
	}

	sourceMetadata := response.Data.SourceMetadata
	err = state.Fill(sourceMetadata)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Source Metadata",
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

func sourceMetadataSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Required:    true,
			Description: "The id for this Source metadata in the Segment catalog. Config API note: analogous to `name`.",
		},
		"name": schema.StringAttribute{
			Computed:    true,
			Description: "The user-friendly name of this Source. Config API note: equal to `displayName`.",
		},
		"slug": schema.StringAttribute{
			Computed:    true,
			Description: "The slug that identifies this Source in the Segment app. Config API note: equal to `name`.",
		},
		"description": schema.StringAttribute{
			Computed:    true,
			Description: "The description of this Source.",
		},
		"logos": schema.SingleNestedAttribute{
			Computed:    true,
			Description: "The logos for this Source.",
			Attributes: map[string]schema.Attribute{
				"default": schema.StringAttribute{
					Computed:    true,
					Description: "The default URL for this logo.",
				},
				"mark": schema.StringAttribute{
					Optional:    true,
					Computed:    true,
					Description: "The logo mark.",
				},
				"alt": schema.StringAttribute{
					Optional:    true,
					Computed:    true,
					Description: "The alternative text for this logo.",
				},
			},
		},
		"options": schema.ListNestedAttribute{
			Computed:    true,
			Description: "Options for this Source.",
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Computed:    true,
						Description: "The name identifying this option in the context of a Segment Integration.",
					},
					"type": schema.StringAttribute{
						Computed:    true,
						Description: "Defines the type for this option in the schema.",
					},
					"required": schema.BoolAttribute{
						Computed:    true,
						Description: "Whether this is a required option when setting up the Integration.",
					},
					"description": schema.StringAttribute{
						Computed:    true,
						Description: "An optional short text description of the field.",
					},
					"default_value": schema.StringAttribute{
						CustomType:  jsontypes.NormalizedType{},
						Computed:    true,
						Description: "An optional default value for the field.",
					},
					"label": schema.StringAttribute{
						Computed:    true,
						Description: "An optional label for this field.",
					},
				},
			},
		},

		"categories": schema.ListAttribute{
			ElementType: types.StringType,
			Computed:    true,
			Description: "A list of categories this Source belongs to.",
		},
		"is_cloud_event_source": schema.BoolAttribute{
			Computed:    true,
			Description: "True if this is a Cloud Event Source.",
		},
	}
}

// Schema defines the schema for the data source.
func (d *sourceMetadataDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The source metadata",
		Attributes:  sourceMetadataSchema(),
	}
}

// Configure adds the provider configured client to the data source.
func (d *sourceMetadataDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
