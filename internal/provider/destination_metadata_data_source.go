package provider

import (
	"context"
	"fmt"

	"terraform-provider-segment/internal/provider/models"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/segmentio/public-api-sdk-go/api"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &destinationMetadataDataSource{}
	_ datasource.DataSourceWithConfigure = &destinationMetadataDataSource{}
)

// NewDestinationMetadataDataSource is a helper function to simplify the provider implementation.
func NewDestinationMetadataDataSource() datasource.DataSource {
	return &destinationMetadataDataSource{}
}

// destinationMetadataDataSource is the data source implementation.
type destinationMetadataDataSource struct {
	client      *api.APIClient
	authContext context.Context
}

func destinationMetadataSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The id of the Destination metadata. Config API note: analogous to `name`.",
			Required:    true,
		},
		"name": schema.StringAttribute{
			Description: "The user-friendly name of the Destination. Config API note: equal to `displayName`.",
			Computed:    true,
		},
		"description": schema.StringAttribute{
			Description: "The description of the Destination.",
			Computed:    true,
		},
		"slug": schema.StringAttribute{
			Description: "The slug used to identify the Destination in the Segment app.",
			Computed:    true,
		},
		"logos": schema.SingleNestedAttribute{
			Description: "The Destination's logos.",
			Computed:    true,
			Attributes: map[string]schema.Attribute{
				"default": schema.StringAttribute{
					Computed: true,
				},
				"mark": schema.StringAttribute{
					Description: "The logo mark.",
					Computed:    true,
				},
				"alt": schema.StringAttribute{
					Description: "The alternative text for this logo.",
					Computed:    true,
				},
			},
		},
		"options": schema.ListNestedAttribute{
			Description: "Options configured for the Destination.",
			Computed:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Description: "The name identifying this option in the context of a Segment Integration.",
						Computed:    true,
					},
					"type": schema.StringAttribute{
						Description: "Defines the type for this option in the schema.",
						Computed:    true,
					},
					"required": schema.BoolAttribute{
						Description: "Whether this is a required option when setting up the Integration.",
						Computed:    true,
					},
					"description": schema.StringAttribute{
						Description: "An optional short text description of the field.",
						Computed:    true,
					},
					"default_value": schema.StringAttribute{
						Description: "An optional default value for the field.",
						Computed:    true,
						CustomType:  jsontypes.NormalizedType{},
					},
					"label": schema.StringAttribute{
						Description: "An optional label for this field.",
						Computed:    true,
					},
				},
			},
		},
		"status": schema.StringAttribute{
			Description: "Support status of the Destination.",
			Computed:    true,
		},
		"previous_names": schema.ListAttribute{
			ElementType: types.StringType,
			Description: "A list of names previously used by the Destination.",
			Computed:    true,
		},
		"categories": schema.ListAttribute{
			ElementType: types.StringType,
			Description: "A list of categories with which the Destination is associated.",
			Computed:    true,
		},
		"website": schema.StringAttribute{
			Description: "A website URL for this Destination.",
			Computed:    true,
		},
		"components": schema.ListNestedAttribute{
			Description: "A list of components this Destination provides.",
			Computed:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Description: "The component type.",
						Computed:    true,
					},
					"code": schema.StringAttribute{
						Description: "Link to the repository hosting the code for this component.",
						Computed:    true,
					},
					"owner": schema.StringAttribute{
						Description: "The owner of this component. Either 'SEGMENT' or 'PARTNER'.",
						Computed:    true,
					},
				},
			},
		},
		"supported_features": schema.SingleNestedAttribute{
			Description: "Features that this Destination supports.",
			Computed:    true,
			Attributes: map[string]schema.Attribute{
				"cloud_mode_instances": schema.StringAttribute{
					Description: "This Destination's support level for cloud mode instances.",
					Computed:    true,
				},
				"device_mode_instances": schema.StringAttribute{
					Description: "This Destination's support level for device mode instances.",
					Computed:    true,
				},
				"replay": schema.BoolAttribute{
					Description: "Whether this Destination supports replays.",
					Computed:    true,
				},
				"browser_unbundling": schema.BoolAttribute{
					Description: "Whether this Destination supports browser unbundling.",
					Computed:    true,
				},
				"browser_unbundling_public": schema.BoolAttribute{
					Description: "Whether this Destination supports public browser unbundling.",
					Computed:    true,
				},
			},
		},
		"supported_methods": schema.SingleNestedAttribute{
			Description: "Methods that this Destination supports.",
			Computed:    true,
			Attributes: map[string]schema.Attribute{
				"pageview": schema.BoolAttribute{
					Description: "Identifies if the Destination supports the `pageview` method.",
					Computed:    true,
				},
				"identify": schema.BoolAttribute{
					Description: "Identifies if the Destination supports the `identify` method.",
					Computed:    true,
				},
				"alias": schema.BoolAttribute{
					Description: "Identifies if the Destination supports the `alias` method.",
					Computed:    true,
				},
				"track": schema.BoolAttribute{
					Description: "Identifies if the Destination supports the `track` method.",
					Computed:    true,
				},
				"group": schema.BoolAttribute{
					Description: "Identifies if the Destination supports the `group` method.",
					Computed:    true,
				},
			},
		},
		"supported_platforms": schema.SingleNestedAttribute{
			Description: "Platforms from which the Destination receives events.",
			Computed:    true,
			Attributes: map[string]schema.Attribute{
				"browser": schema.BoolAttribute{
					Description: "Whether this Destination supports browser events.",
					Computed:    true,
				},
				"server": schema.BoolAttribute{
					Description: "Whether this Destination supports server events.",
					Computed:    true,
				},
				"mobile": schema.BoolAttribute{
					Description: "Whether this Destination supports mobile events.",
					Computed:    true,
				},
			},
		},
		"actions": schema.ListNestedAttribute{
			Description: "Actions available for the Destination.",
			Computed:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Description: "The primary key of the action.",
						Computed:    true,
					},
					"slug": schema.StringAttribute{
						Description: "A machine-readable key unique to the action definition.",
						Computed:    true,
					},
					"name": schema.StringAttribute{
						Description: "A human-readable name for the action.",
						Computed:    true,
					},
					"description": schema.StringAttribute{
						Description: "A human-readable description of the action. May include Markdown.",
						Computed:    true,
					},
					"platform": schema.StringAttribute{
						Description: "The platform on which this action runs.",
						Computed:    true,
					},
					"hidden": schema.BoolAttribute{
						Description: "Whether the action should be hidden.",
						Computed:    true,
					},
					"default_trigger": schema.StringAttribute{
						Description: "The default value used as the trigger when connecting this action.",
						Optional:    true,
					},
					"fields": schema.ListNestedAttribute{
						Description: "The fields expected in order to perform the action.",
						Computed:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									Description: "The primary key of the field.",
									Computed:    true,
								},
								"sort_order": schema.Float64Attribute{
									Description: "The order this particular field is (used in the UI for displaying the fields in a specified order).",
									Computed:    true,
								},
								"field_key": schema.StringAttribute{
									Description: "A unique machine-readable key for the field. Should ideally match the expected key in the action's API request.",
									Computed:    true,
								},
								"label": schema.StringAttribute{
									Description: "A human-readable label for this value.",
									Computed:    true,
								},
								"type": schema.StringAttribute{
									Description: "The data type for this value.",
									Computed:    true,
								},
								"description": schema.StringAttribute{
									Description: "A human-readable description of this value. You can use Markdown.",
									Computed:    true,
								},
								"placeholder": schema.StringAttribute{
									Description: "An example value displayed but not saved.",
									Computed:    true,
								},
								"default_value": schema.StringAttribute{
									Description: "A default value that is saved the first time an action is created.",
									Computed:    true,
									CustomType:  jsontypes.NormalizedType{},
								},
								"required": schema.BoolAttribute{
									Description: "Whether this field is required.",
									Computed:    true,
								},
								"multiple": schema.BoolAttribute{
									Description: "Whether a user can provide multiples of this field.",
									Computed:    true,
								},
								"choices": schema.StringAttribute{
									Description: "A list of machine-readable value/label pairs to populate a static dropdown.",
									Computed:    true,
									CustomType:  jsontypes.NormalizedType{},
								},
								"dynamic": schema.BoolAttribute{
									Description: "Whether this field should execute a dynamic request to fetch choices to populate a dropdown. When true, `choices` is ignored.",
									Computed:    true,
								},
								"allow_null": schema.BoolAttribute{
									Description: "Whether this field allows null values.",
									Computed:    true,
								},
							},
						},
					},
				},
			},
		},
		"presets": schema.ListNestedAttribute{
			Description: "Predefined Destination subscriptions that can optionally be applied when connecting a new instance of the Destination.",
			Computed:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"action_id": schema.StringAttribute{
						Description: "The unique identifier for the Destination Action to trigger.",
						Computed:    true,
					},
					"name": schema.StringAttribute{
						Description: "The name of the subscription.",
						Computed:    true,
					},
					"fields": schema.StringAttribute{
						Computed:    true,
						Description: "The default settings for action fields.",
						CustomType:  jsontypes.NormalizedType{},
					},
					"trigger": schema.StringAttribute{
						Description: "FQL string that describes what events should trigger an action. See https://segment.com/docs/config-api/fql/ for more information regarding Segment's Filter Query Language (FQL).",
						Computed:    true,
					},
				},
			},
		},
		"contacts": schema.ListNestedAttribute{
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Description: "Name of this contact.",
						Computed:    true,
					},
					"email": schema.StringAttribute{
						Description: "Email of this contact.",
						Computed:    true,
					},
					"role": schema.StringAttribute{
						Description: "Role of this contact.",
						Computed:    true,
					},
					"is_primary": schema.BoolAttribute{
						Description: "Whether this is a primary contact.",
						Computed:    true,
					},
				},
			},
			Description: "Contact info for Integration Owners.",
			Computed:    true,
		},
		"partner_owned": schema.BoolAttribute{
			Description: "Partner Owned flag.",
			Computed:    true,
		},
		"supported_regions": schema.ListAttribute{
			ElementType: types.StringType,
			Description: "A list of supported regions for this Destination.",
			Computed:    true,
		},
		"region_endpoints": schema.ListAttribute{
			ElementType: types.StringType,
			Description: "The list of regional endpoints for this Destination.",
			Computed:    true,
		},
	}
}

// Metadata returns the data source type name.
func (d *destinationMetadataDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_destination_metadata"
}

// Schema defines the schema for the data source.
func (d *destinationMetadataDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The destination metadata",
		Attributes:  destinationMetadataSchema(),
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *destinationMetadataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state models.DestinationMetadataState

	diags := req.Config.Get(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, body, err := d.client.CatalogApi.GetDestinationMetadata(d.authContext, state.ID.ValueString()).Execute()
	defer body.Body.Close()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Source metadata",
			getError(err, body),
		)

		return
	}

	destinationMetadata := response.Data.DestinationMetadata
	err = state.Fill(destinationMetadata)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Source metadata",
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

// Configure adds the provider configured client to the data source.
func (d *destinationMetadataDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
