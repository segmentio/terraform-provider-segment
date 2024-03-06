package provider

import (
	"context"
	"fmt"

	"github.com/segmentio/terraform-provider-segment/internal/provider/models"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/segmentio/public-api-sdk-go/api"
)

var (
	_ datasource.DataSource              = &sourceDataSource{}
	_ datasource.DataSourceWithConfigure = &sourceDataSource{}
)

type sourceDataSource struct {
	client      *api.APIClient
	authContext context.Context
}

func NewSourceDataSource() datasource.DataSource {
	return &sourceDataSource{}
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

func (d *sourceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_source"
}

func (d *sourceDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Reads a Source. For more information, visit the [Segment docs](https://segment.com/docs/connections/sources/).",
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
			"settings": schema.StringAttribute{
				Computed:    true,
				Description: "The settings associated with the Source.",
				CustomType:  jsontypes.NormalizedType{},
			},
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
					},
				},
			},
			"schema_settings": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "The schema settings associated with the Source.",
				Attributes: map[string]schema.Attribute{
					"track": schema.SingleNestedAttribute{
						Computed:    true,
						Description: "Track settings.",
						Attributes: map[string]schema.Attribute{
							"allow_unplanned_events": schema.BoolAttribute{
								Computed:    true,
								Description: "Enable to allow unplanned track events.",
							},
							"allow_unplanned_event_properties": schema.BoolAttribute{
								Computed:    true,
								Description: "Enable to allow unplanned track event properties.",
							},
							"allow_event_on_violations": schema.BoolAttribute{
								Computed:    true,
								Description: "Allow track event on violations.",
							},
							"allow_properties_on_violations": schema.BoolAttribute{
								Computed:    true,
								Description: "Enable to allow track properties on violations.",
							},
							"common_event_on_violations": schema.StringAttribute{
								Computed:    true,
								Description: "The common track event on violations.",
							},
						},
					},
					"identify": schema.SingleNestedAttribute{
						Computed:    true,
						Description: "Identify settings.",
						Attributes: map[string]schema.Attribute{
							"allow_unplanned_traits": schema.BoolAttribute{
								Computed:    true,
								Description: "Enable to allow unplanned identify traits.",
							},
							"allow_traits_on_violations": schema.BoolAttribute{
								Computed:    true,
								Description: "Enable to allow identify traits on violations.",
							},
							"common_event_on_violations": schema.StringAttribute{
								Computed:    true,
								Description: "The common identify event on violations.",
							},
						},
					},
					"group": schema.SingleNestedAttribute{
						Computed:    true,
						Description: "Group settings.",
						Attributes: map[string]schema.Attribute{
							"allow_unplanned_traits": schema.BoolAttribute{
								Computed:    true,
								Description: "Enable to allow unplanned group traits.",
							},
							"allow_traits_on_violations": schema.BoolAttribute{
								Computed:    true,
								Description: "Enable to allow group traits on violations.",
							},
							"common_event_on_violations": schema.StringAttribute{
								Computed:    true,
								Description: "The common group event on violations.",
							},
						},
					},
					"forwarding_violations_to": schema.StringAttribute{
						Computed:    true,
						Description: "Source id to forward violations to.",
					},
					"forwarding_blocked_events_to": schema.StringAttribute{
						Computed:    true,
						Description: "Source id to forward blocked events to.",
					},
				},
			},
		},
	}
}

func (d *sourceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config models.SourceDataSourceState
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := config.ID.ValueString()
	if id == "" {
		resp.Diagnostics.AddError("Unable to read Source", "ID is empty")

		return
	}

	out, body, err := d.client.SourcesAPI.GetSource(d.authContext, id).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Source",
			getError(err, body),
		)

		return
	}

	source := out.Data.Source

	var schemaSettings *api.SourceSettingsOutputV1

	if out.Data.TrackingPlanId.IsSet() && out.Data.TrackingPlanId.Get() != nil {
		settingsOut, body, err := d.client.SourcesAPI.ListSchemaSettingsInSource(d.authContext, source.Id).Execute()
		if body != nil {
			defer body.Body.Close()
		}
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to read Source schema settings",
				getError(err, body),
			)

			return
		}

		schemaSettings = &settingsOut.Data.Settings
	}

	var state models.SourceDataSourceState
	err = state.Fill(source, schemaSettings)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to populate Source state",
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
