package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	api "github.com/segmentio/public-api-sdk-go/api"
)

var (
	_ datasource.DataSource              = &sourceDataSource{}
	_ datasource.DataSourceWithConfigure = &sourceDataSource{}
)

func NewSourceDataSource() datasource.DataSource {
	return &sourceDataSource{}
}

type sourceDataSource struct {
	client      *api.APIClient
	authContext context.Context
}

type sourceDataSourceModel struct {
	Enabled     types.Bool      `tfsdk:"enabled"`
	ID          types.String    `tfsdk:"id"`
	Labels      []Label         `tfsdk:"labels"`
	Metadata    *SourceMetadata `tfsdk:"metadata"`
	Name        types.String    `tfsdk:"name"`
	Slug        types.String    `tfsdk:"slug"`
	WorkspaceID types.String    `tfsdk:"workspace_id"`
	WriteKeys   []types.String  `tfsdk:"write_keys"`
}

type Label struct {
	Description types.String `tfsdk:"description"`
	Key         types.String `tfsdk:"key"`
	Value       types.String `tfsdk:"value"`
}

type SourceMetadata struct {
	Categories         *[]types.String     `tfsdk:"categories"`
	Description        types.String        `tfsdk:"description"`
	ID                 types.String        `tfsdk:"id"`
	IsCloudEventSource types.Bool          `tfsdk:"is_cloud_event_source"`
	Logos              *Logos              `tfsdk:"logos"`
	Name               types.String        `tfsdk:"name"`
	Options            []IntegrationOption `tfsdk:"options"`
	Slug               types.String        `tfsdk:"slug"`
}

type IntegrationOption struct {
	DefaultValue interface{}  `tfsdk:"default_value"`
	Description  types.String `tfsdk:"description"`
	Label        types.String `tfsdk:"label"`
	Name         types.String `tfsdk:"name"`
	Required     types.Bool   `tfsdk:"required"`
	Type         types.String `tfsdk:"type"`
}

type Logos struct {
	Alt     types.String `tfsdk:"alt"`
	Default types.String `tfsdk:"default"`
	Mark    types.String `tfsdk:"mark"`
}

func (d *sourceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_source"
}

func (d *sourceDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "The id of the Source.",
			},
			"slug": schema.StringAttribute{
				Computed:    true,
				Description: "The slug used to identify the Source in the Segment app.",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The name of the Source.",
			},
			"metadata": schema.SingleNestedAttribute{
				Description: "The metadata for the Source.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed:    true,
						Description: "The id for this Source metadata in the Segment catalog.",
					},
					"name": schema.StringAttribute{
						Computed:    true,
						Description: "The user-friendly name of this Source.",
					},
					"slug": schema.StringAttribute{
						Computed:    true,
						Description: "The slug that identifies this Source in the Segment app.",
					},
					"description": schema.StringAttribute{
						Computed:    true,
						Description: "The description of this Source.",
					},
					"logos": schema.SingleNestedAttribute{
						Description: "The logos for this Source.",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"default": schema.StringAttribute{
								Computed:    true,
								Description: "The default URL for this logo.",
							},
							"mark": schema.StringAttribute{
								Computed:    true,
								Description: "The logo mark.",
							},
							"alt": schema.StringAttribute{
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
									Description: "Defines the type for this option in the schema. Types are most commonly strings, but may also represent other primitive types, such as booleans, and numbers, as well as complex types, such as objects and arrays.",
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
						Computed:    true,
						ElementType: types.StringType,
						Description: "A list of categories this Source belongs to.",
					},
					"is_cloud_event_source": schema.BoolAttribute{
						Computed:    true,
						Description: "True if this is a Cloud Event Source.",
					},
				},
			},
			// TODO: Support settings
			"workspace_id": schema.StringAttribute{
				Computed:    true,
				Description: "The id of the Workspace that owns the Source.",
			},
			"enabled": schema.BoolAttribute{
				Computed:    true,
				Description: "Enable to receive data from the Source.",
			},
			"write_keys": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "The write keys used to send data from the Source. This field is left empty when the current token does not have the 'source admin' permission.",
			},
			"labels": schema.ListNestedAttribute{
				Computed:    true,
				Description: "A list of labels applied to the Source.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"key": schema.StringAttribute{
							Computed:    true,
							Description: "The key that represents the name of this label.",
						},
						"value": schema.StringAttribute{
							Computed:    true,
							Description: "The value associated with the key of this label.",
						},
						"description": schema.BoolAttribute{
							Computed:    true,
							Description: "An optional description of the purpose of this label.",
						},
					},
				},
			},
		},
	}
}

func (d *sourceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config sourceDataSourceModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	out, _, err := d.client.SourcesApi.GetSource(d.authContext, config.ID.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Source",
			err.Error(),
		)
		return
	}

	source := out.Data.Source

	var state sourceDataSourceModel

	state.ID = types.StringValue(source.Id)
	state.Name = types.StringValue(*source.Name)
	state.Slug = types.StringValue(source.Slug)
	state.Enabled = types.BoolValue(source.Enabled)
	state.WorkspaceID = types.StringValue(source.WorkspaceId)

	// Populate write keys
	for _, writeKey := range source.WriteKeys {
		state.WriteKeys = append(state.WriteKeys, types.StringValue(writeKey))
	}

	// Populate labels
	for _, label := range source.Labels {
		state.Labels = append(state.Labels, Label{
			Key:         types.StringValue(label.Key),
			Value:       types.StringValue(label.Value),
			Description: types.StringValue(*label.Description),
		})
	}

	// TODO: Populate settings

	// Populate source metadata
	state.Metadata = &SourceMetadata{
		ID:                 types.StringValue(source.Metadata.Id),
		Description:        types.StringValue(source.Metadata.Description),
		IsCloudEventSource: types.BoolValue(source.Metadata.IsCloudEventSource),
		Logos: &Logos{
			Alt:     types.StringValue(*source.Metadata.Logos.Alt.Get()),
			Default: types.StringValue(source.Metadata.Logos.Default),
			Mark:    types.StringValue(*source.Metadata.Logos.Mark.Get()),
		},
		Name: types.StringValue(source.Metadata.Name),
		Slug: types.StringValue(source.Metadata.Slug),
	}

	for _, metadataCategory := range source.Metadata.Categories {
		*state.Metadata.Categories = append(*state.Metadata.Categories, types.StringValue(metadataCategory))
	}

	for _, integrationOption := range source.Metadata.Options {

		state.Metadata.Options = append(state.Metadata.Options, IntegrationOption{
			Name:         types.StringValue(integrationOption.Name),
			Type:         types.StringValue(integrationOption.Type),
			Required:     types.BoolValue(integrationOption.Required),
			Description:  types.StringValue(*integrationOption.Description),
			DefaultValue: types.StringValue(fmt.Sprintf("%v", integrationOption.DefaultValue)),
		})
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *sourceDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
