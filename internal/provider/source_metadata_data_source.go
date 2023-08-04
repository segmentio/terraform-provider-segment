package provider

import (
	"context"
	"fmt"
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

type sourceMetadataDataSourceModel struct {
	Id                 types.String        `tfsdk:"id"`
	Name               types.String        `tfsdk:"name"`
	Slug               types.String        `tfsdk:"slug"`
	Description        types.String        `tfsdk:"description"`
	Logos              *Logos              `tfsdk:"logos"`
	Options            []IntegrationOption `tfsdk:"options"`
	Categories         []types.String      `tfsdk:"categories"`
	IsCloudEventSource types.Bool          `tfsdk:"is_cloud_event_source"`
}

// Metadata returns the data source type name.
func (d *sourceMetadataDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_source_metadata"
}

// Read refreshes the Terraform state with the latest data.
func (d *sourceMetadataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state sourceMetadataDataSourceModel

	diags := req.Config.Get(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, _, err := d.client.CatalogApi.GetSourceMetadata(d.authContext, state.Id.ValueString()).Execute()

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Source metadata",
			err.Error(),
		)
		return
	}

	var sourceMetadata = response.Data.SourceMetadata

	state.Id = types.StringValue(sourceMetadata.Id)
	state.Name = types.StringValue(sourceMetadata.Name)
	state.Description = types.StringValue(sourceMetadata.Description)
	state.Slug = types.StringValue(sourceMetadata.Slug)
	state.Logos = getLogosSourceMetadata(sourceMetadata.Logos)
	state.Options = getOptions(sourceMetadata.Options)
	state.IsCloudEventSource = types.BoolValue(sourceMetadata.IsCloudEventSource)
	state.Categories = getCategories(sourceMetadata.Categories)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func getCategories(categories []string) []types.String {
	var categoriesToAdd []types.String

	for _, cat := range categories {
		categoriesToAdd = append(categoriesToAdd, types.StringValue(cat))
	}

	return categoriesToAdd
}

func getLogosSourceMetadata(logos api.Logos1) *Logos {
	logosToAdd := Logos{
		Default: types.StringValue(logos.Default),
	}

	if logos.Mark.IsSet() {
		logosToAdd.Mark = types.StringValue(*logos.Mark.Get())
	}

	if logos.Alt.IsSet() {
		logosToAdd.Alt = types.StringValue(*logos.Alt.Get())
	}

	return &logosToAdd
}

func sourceMetadataSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed:    true,
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
						Optional:    true,
						Computed:    true,
						Description: "An optional short text description of the field.",
					},
					//TODO: There is no equivalent of schema.AnyAttribute, therefore this field is ignored.
					//"default_value": {
					//	Type:        schema.TypeAny,
					//	Optional:    true,
					//	Computed:    true,
					//	Description: "An optional default value for the field.",
					"label": schema.StringAttribute{
						Optional:    true,
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