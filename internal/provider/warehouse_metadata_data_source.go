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
	_ datasource.DataSource              = &warehouseMetadataDataSource{}
	_ datasource.DataSourceWithConfigure = &warehouseMetadataDataSource{}
)

// NewWarehouseMetadataDataSource is a helper function to simplify the provider implementation.
func NewWarehouseMetadataDataSource() datasource.DataSource {
	return &warehouseMetadataDataSource{}
}

// warehouseMetadataDataSource is the data source implementation.
type warehouseMetadataDataSource struct {
	client      *api.APIClient
	authContext context.Context
}

type warehouseMetadataDataSourceModel struct {
	Id          types.String        `tfsdk:"id"`
	Name        types.String        `tfsdk:"name"`
	Slug        types.String        `tfsdk:"slug"`
	Description types.String        `tfsdk:"description"`
	Logos       *Logos              `tfsdk:"logos"`
	Options     []IntegrationOption `tfsdk:"options"`
}

func warehouseMetadataSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed:    true,
			Description: "The id of this object.",
		},
		"name": schema.StringAttribute{
			Computed:    true,
			Description: "The name of this object.",
		},
		"slug": schema.StringAttribute{
			Computed:    true,
			Description: "A human-readable, unique identifier for object.",
		},
		"description": schema.StringAttribute{
			Computed:    true,
			Description: "A description, in English, of this object.",
		},
		"logos": schema.SingleNestedAttribute{
			Computed:    true,
			Description: "Logo information for this object.",
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
			Description: "The Integration options for this object.",
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
					//},
					"label": schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Description: "An optional label for this field.",
					},
				},
			},
		},
	}
}

// Metadata returns the data source type name.
func (d *warehouseMetadataDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_warehouse_metadata"
}

// Read refreshes the Terraform state with the latest data.
func (d *warehouseMetadataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state warehouseMetadataDataSourceModel

	diags := req.Config.Get(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, _, err := d.client.CatalogApi.GetWarehouseMetadata(d.authContext, state.Id.ValueString()).Execute()

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Warehouse metadata",
			err.Error(),
		)
		return
	}

	var warehouseMetadata = response.Data.WarehouseMetadata

	state.Id = types.StringValue(warehouseMetadata.Id)
	state.Name = types.StringValue(warehouseMetadata.Name)
	state.Description = types.StringValue(warehouseMetadata.Description)
	state.Slug = types.StringValue(warehouseMetadata.Slug)
	state.Logos = getLogos2(warehouseMetadata.Logos)
	state.Options = getOptions(warehouseMetadata.Options)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func getOptions(options []api.IntegrationOptionBeta) []IntegrationOption {
	var integrationOptions []IntegrationOption

	for _, opt := range options {
		integrationOption := IntegrationOption{
			Name:     types.StringValue(opt.Name),
			Type:     types.StringValue(opt.Type),
			Required: types.BoolValue(opt.Required),
		}

		if opt.Description != nil {
			integrationOption.Description = types.StringValue(*opt.Description)
		}

		if opt.Label != nil {
			integrationOption.Label = types.StringValue(*opt.Label)
		}

		//TODO: Handle defaultValue which is a field of type Any
		//if opt.DefaultValue != nil {
		//	integrationOption.Label = types.value(*opt.DefaultValue)
		//}

		integrationOptions = append(integrationOptions, integrationOption)
	}

	return integrationOptions
}

func getLogos2(logos api.Logos2) *Logos {
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

// Schema defines the schema for the data source.
func (d *warehouseMetadataDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The warehouse metadata",
		Attributes:  warehouseMetadataSchema(),
	}
}

// Configure adds the provider configured client to the data source.
func (d *warehouseMetadataDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
