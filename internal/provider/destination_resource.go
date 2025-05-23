package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"

	"github.com/segmentio/terraform-provider-segment/internal/provider/docs"
	"github.com/segmentio/terraform-provider-segment/internal/provider/models"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/segmentio/public-api-sdk-go/api"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &destinationResource{}
	_ resource.ResourceWithConfigure   = &destinationResource{}
	_ resource.ResourceWithImportState = &destinationResource{}
)

// NewDestinationResource is a helper function to simplify the provider implementation.
func NewDestinationResource() resource.Resource {
	return &destinationResource{}
}

// destinationResource is the resource implementation.
type destinationResource struct {
	client      *api.APIClient
	authContext context.Context
}

// Metadata returns the resource type name.
func (r *destinationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_destination"
}

// Schema defines the schema for the resource.
func (r *destinationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Configures a Destination. For more information, visit the [Segment docs](https://segment.com/docs/connections/destinations/).\n\n" +
			docs.GenerateImportDocs("<id>", "segment_destination"),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"enabled": schema.BoolAttribute{
				Required: true,
			},
			"source_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"metadata": schema.SingleNestedAttribute{
				Required:   true,
				Attributes: destinationMetadataResourceSchema(),
			},
			"settings": schema.StringAttribute{
				Required:    true,
				Description: "The settings associated with the Destination. Only settings included in the configuration will be managed by Terraform.",
				CustomType:  jsontypes.NormalizedType{},
			},
		},
	}
}

func destinationMetadataResourceSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The id of the Destination metadata. Config API note: analogous to `name`.",
			Required:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"name": schema.StringAttribute{
			Description: "The user-friendly name of the Destination. Config API note: equal to `displayName`.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"description": schema.StringAttribute{
			Description: "The description of the Destination.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"slug": schema.StringAttribute{
			Description: "The slug used to identify the Destination in the Segment app.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"logos": schema.SingleNestedAttribute{
			Description: "The Destination's logos.",
			Computed:    true,
			PlanModifiers: []planmodifier.Object{
				objectplanmodifier.UseStateForUnknown(),
			},
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
			PlanModifiers: []planmodifier.List{
				listplanmodifier.UseStateForUnknown(),
			},
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
						CustomType:  jsontypes.NormalizedType{},
						Description: "An optional default value for the field.",
						Computed:    true,
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
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"previous_names": schema.ListAttribute{
			ElementType: types.StringType,
			Description: "A list of names previously used by the Destination.",
			Computed:    true,
			PlanModifiers: []planmodifier.List{
				listplanmodifier.UseStateForUnknown(),
			},
		},
		"categories": schema.ListAttribute{
			ElementType: types.StringType,
			Description: "A list of categories with which the Destination is associated.",
			Computed:    true,
			PlanModifiers: []planmodifier.List{
				listplanmodifier.UseStateForUnknown(),
			},
		},
		"website": schema.StringAttribute{
			Description: "A website URL for this Destination.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"components": schema.ListNestedAttribute{
			Description: "A list of components this Destination provides.",
			Computed:    true,
			PlanModifiers: []planmodifier.List{
				listplanmodifier.UseStateForUnknown(),
			},
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
			PlanModifiers: []planmodifier.Object{
				objectplanmodifier.UseStateForUnknown(),
			},
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
			PlanModifiers: []planmodifier.Object{
				objectplanmodifier.UseStateForUnknown(),
			},
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
			PlanModifiers: []planmodifier.Object{
				objectplanmodifier.UseStateForUnknown(),
			},
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
			PlanModifiers: []planmodifier.List{
				listplanmodifier.UseStateForUnknown(),
			},
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
						Computed:    true,
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
									CustomType:  jsontypes.NormalizedType{},
									Description: "A default value that is saved the first time an action is created.",
									Computed:    true,
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
									CustomType:  jsontypes.NormalizedType{},
									Description: "A list of machine-readable value/label pairs to populate a static dropdown.",
									Computed:    true,
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
			PlanModifiers: []planmodifier.List{
				listplanmodifier.UseStateForUnknown(),
			},
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
						CustomType:  jsontypes.NormalizedType{},
						Computed:    true,
						Description: "The default settings for action fields.",
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
			Optional:    true,
			PlanModifiers: []planmodifier.List{
				listplanmodifier.UseStateForUnknown(),
			},
		},
		"partner_owned": schema.BoolAttribute{
			Description: "Partner Owned flag.",
			Computed:    true,
			Optional:    true,
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.UseStateForUnknown(),
			},
		},
		"supported_regions": schema.ListAttribute{
			ElementType: types.StringType,
			Description: "A list of supported regions for this Destination.",
			Computed:    true,
			Optional:    true,
			PlanModifiers: []planmodifier.List{
				listplanmodifier.UseStateForUnknown(),
			},
		},
		"region_endpoints": schema.ListAttribute{
			ElementType: types.StringType,
			Description: "The list of regional endpoints for this Destination.",
			Computed:    true,
			Optional:    true,
			PlanModifiers: []planmodifier.List{
				listplanmodifier.UseStateForUnknown(),
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *destinationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan models.DestinationPlan
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	wrappedMetadataID, err := plan.Metadata.Attributes()["id"].ToTerraformValue(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to decode metadata id",
			err.Error(),
		)

		return
	}

	var metadataID string
	err = wrappedMetadataID.As(&metadataID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to decode metadata id",
			err.Error(),
		)

		return
	}

	var settings map[string]interface{}
	diags = plan.Settings.Unmarshal(&settings)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := api.CreateDestinationV1Input{
		SourceId:   plan.SourceID.ValueString(),
		MetadataId: metadataID,
		Enabled:    plan.Enabled.ValueBoolPointer(),
		Name:       plan.Name.ValueStringPointer(),
		Settings:   settings,
	}

	// Generate API request body from plan
	out, body, err := r.client.DestinationsAPI.CreateDestination(r.authContext).CreateDestinationV1Input(input).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create Destination",
			getError(err, body),
		)

		return
	}
	resp.State.SetAttribute(ctx, path.Root("id"), out.Data.Destination.Id)

	var state models.DestinationState
	err = state.Fill(&out.Data.Destination)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to populate Destination state",
			err.Error(),
		)

		return
	}

	// This is to satisfy terraform requirements that the returned fields must match the input ones because new settings can be generated in the response
	state.Settings = plan.Settings

	// Set state to fully populated data
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *destinationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var previousState models.DestinationState
	diags := req.State.Get(ctx, &previousState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	out, body, err := r.client.DestinationsAPI.GetDestination(r.authContext, previousState.ID.ValueString()).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		if body.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)

			return
		}

		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to read Destination (ID: %s)", previousState.ID.ValueString()),
			getError(err, body),
		)

		return
	}

	destination := out.Data.Destination

	var state models.DestinationState
	err = state.Fill(&destination)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to populate Destination state",
			err.Error(),
		)

		return
	}

	// This is to satisfy terraform requirements that the returned fields must match the input ones because new settings can be generated in the response
	if !previousState.Settings.IsNull() && !previousState.Settings.IsUnknown() {
		state.Settings = previousState.Settings
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *destinationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan models.DestinationPlan
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var settings map[string]interface{}
	diags = plan.Settings.Unmarshal(&settings)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := api.UpdateDestinationV1Input{
		Name:     *api.NewNullableString(plan.Name.ValueStringPointer()),
		Enabled:  plan.Enabled.ValueBoolPointer(),
		Settings: settings,
	}

	// Generate API request body from plan
	out, body, err := r.client.DestinationsAPI.UpdateDestination(r.authContext, plan.ID.ValueString()).UpdateDestinationV1Input(input).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to update Destination (ID: %s)", plan.ID.ValueString()),
			getError(err, body),
		)

		return
	}

	var state models.DestinationState
	err = state.Fill(&out.Data.Destination)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to populate Destination state",
			err.Error(),
		)

		return
	}

	// This is to satisfy terraform requirements that the returned fields must match the input ones because new settings can be generated in the response
	state.Settings = plan.Settings

	// Set state to fully populated data
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *destinationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state models.DestinationState
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, body, err := r.client.DestinationsAPI.DeleteDestination(r.authContext, state.ID.ValueString()).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to delete Destination (ID: %s)", state.ID.ValueString()),
			getError(err, body),
		)

		return
	}
}

func (r *destinationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Configure adds the provider configured client to the resource.
func (r *destinationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	config, ok := req.ProviderData.(*ClientInfo)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected ClientInfo, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = config.client
	r.authContext = config.authContext
}
