package provider

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/segmentio/terraform-provider-segment/internal/provider/docs"
	"github.com/segmentio/terraform-provider-segment/internal/provider/models"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/segmentio/public-api-sdk-go/api"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &destinationFilterResource{}
	_ resource.ResourceWithConfigure   = &destinationFilterResource{}
	_ resource.ResourceWithImportState = &destinationFilterResource{}
)

// NewDestinationFilterResource is a helper function to simplify the provider implementation.
func NewDestinationFilterResource() resource.Resource {
	return &destinationFilterResource{}
}

// destinationFilterResource is the resource implementation.
type destinationFilterResource struct {
	client      *api.APIClient
	authContext context.Context
}

// Metadata returns the resource type name.
func (r *destinationFilterResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_destination_filter"
}

// Schema defines the schema for the resource.
func (r *destinationFilterResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Configures a filter for a destination. For more information, visit the [Segment docs](https://segment.com/docs/connections/destinations/destination-filters/).\n\n" +
			docs.GenerateImportDocs("<destination_id>:<filter_id>", "segment_destination_filter"),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The unique id of this filter.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"if": schema.StringAttribute{
				Required:    true,
				Description: "The filter's condition.",
			},
			"destination_id": schema.StringAttribute{
				Required:    true,
				Description: "The id of the Destination associated with this filter.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"source_id": schema.StringAttribute{
				Required:    true,
				Description: "The id of the Source associated with this filter.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"title": schema.StringAttribute{
				Required:    true,
				Description: "The title of the filter.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "The description of the filter.",
			},
			"enabled": schema.BoolAttribute{
				Required:    true,
				Description: "When set to true, the Destination filter is active.",
			},
			"actions": schema.SetNestedAttribute{
				Required:    true,
				Description: "Actions for the Destination filter.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Required: true,
							Description: `The kind of Transformation to apply to any matched properties.

								Enum: "ALLOW_PROPERTIES" "DROP" "DROP_PROPERTIES" "SAMPLE"`,
							Validators: []validator.String{
								stringvalidator.RegexMatches(regexp.MustCompile("^[A-Z_]+$"), "'type' must be in all uppercase"),
							},
						},
						"percent": schema.Float64Attribute{
							Optional:    true,
							Description: "A decimal between 0 and 1 used for 'sample' type events and influences the likelihood of sampling to occur.",
						},
						"path": schema.StringAttribute{
							Optional:    true,
							Description: "The JSON path to a property within a payload object from which Segment generates a deterministic sampling rate.",
						},
						"fields": schema.StringAttribute{
							Optional:    true,
							Description: "A dictionary of paths to object keys that this filter applies to. The literal string '' represents the top level of the object.",
							CustomType:  jsontypes.NormalizedType{},
						},
					},
				},
			},
		},
	}
}

func (r *destinationFilterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan models.DestinationFilterPlan
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	actions, diags := models.ActionsPlanToAPIActions(ctx, plan.Actions)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := validateActions(actions)
	if err != nil {
		resp.Diagnostics.AddError(
			"Destination Filter actions are invalid",
			err.Error(),
		)
	}

	// Generate API request body from plan
	out, body, err := r.client.DestinationFiltersAPI.CreateFilterForDestination(r.authContext, plan.DestinationID.ValueString()).CreateFilterForDestinationV1Input(api.CreateFilterForDestinationV1Input{
		SourceId:    plan.SourceID.ValueString(),
		If:          plan.If.ValueString(),
		Title:       plan.Title.ValueString(),
		Description: plan.Description.ValueStringPointer(),
		Enabled:     plan.Enabled.ValueBool(),
		Actions:     actions,
	}).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create Destination Filter",
			getError(err, body),
		)

		return
	}

	destinationfilter := out.Data.Filter
	resp.State.SetAttribute(ctx, path.Root("id"), destinationfilter.Id)
	resp.State.SetAttribute(ctx, path.Root("destination_id"), destinationfilter.DestinationId)

	var state models.DestinationFilterState
	err = state.Fill(&destinationfilter)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to populate Destination Filter state",
			err.Error(),
		)

		return
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *destinationFilterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var previousState models.DestinationFilterState
	diags := req.State.Get(ctx, &previousState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	out, body, err := r.client.DestinationFiltersAPI.GetFilterInDestination(r.authContext, previousState.DestinationID.ValueString(), previousState.ID.ValueString()).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Destination Filter",
			getError(err, body),
		)

		return
	}

	destinationFilter := out.Data.Filter

	var state models.DestinationFilterState
	err = state.Fill(&destinationFilter)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to populate Destination Filter state",
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

func (r *destinationFilterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan models.DestinationFilterPlan
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state models.DestinationFilterState
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	actions, diags := models.ActionsPlanToAPIActions(ctx, plan.Actions)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := validateActions(actions)
	if err != nil {
		resp.Diagnostics.AddError(
			"Destination Filter actions are invalid",
			err.Error(),
		)
	}

	// Generate API request body from plan
	out, body, err := r.client.DestinationFiltersAPI.UpdateFilterForDestination(r.authContext, state.DestinationID.ValueString(), state.ID.ValueString()).UpdateFilterForDestinationV1Input(api.UpdateFilterForDestinationV1Input{
		If:          plan.If.ValueStringPointer(),
		Title:       plan.Title.ValueStringPointer(),
		Enabled:     plan.Enabled.ValueBoolPointer(),
		Description: *api.NewNullableString(plan.Description.ValueStringPointer()),
		Actions:     actions,
	}).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to update Destination Filter",
			getError(err, body),
		)

		return
	}

	destinationFilter := out.Data.Filter

	err = state.Fill(&destinationFilter)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to populate Destination Filter state",
			err.Error(),
		)

		return
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *destinationFilterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state models.DestinationFilterState
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, body, err := r.client.DestinationFiltersAPI.RemoveFilterFromDestination(r.authContext, state.DestinationID.ValueString(), state.ID.ValueString()).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to delete Destination Filter",
			getError(err, body),
		)

		return
	}
}

func (r *destinationFilterResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ":")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: <destination_id>:<filter_id>. Got: %q", req.ID),
		)

		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("destination_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idParts[1])...)
}

func (r *destinationFilterResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func validateActions(actions []api.DestinationFilterActionV1) error {
	for _, action := range actions {
		if action.Type == "SAMPLE" {
			if action.Fields != nil {
				return fmt.Errorf("'fields' cannot be set for 'SAMPLE' action")
			}
		}
		if action.Type == "DROP" {
			if action.Fields != nil {
				return fmt.Errorf("'fields' cannot be set for 'DROP' action")
			}
			if action.Path != nil {
				return fmt.Errorf("'path' cannot be set for 'DROP' action")
			}
			if action.Percent != nil {
				return fmt.Errorf("'percent' cannot be set for 'DROP' action")
			}
		}
		if action.Type == "ALLOW_PROPERTIES" || action.Type == "DROP_PROPERTIES" {
			if action.Path != nil {
				return fmt.Errorf("'path' cannot be set for '%s' action", action.Type)
			}
			if action.Percent != nil {
				return fmt.Errorf("'percent' cannot be set for '%s' action", action.Type)
			}
		}
	}

	return nil
}
