package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/segmentio/terraform-provider-segment/internal/provider/models"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/segmentio/public-api-sdk-go/api"
)

var (
	_ resource.Resource                = &destinationSubscriptionResource{}
	_ resource.ResourceWithConfigure   = &destinationSubscriptionResource{}
	_ resource.ResourceWithImportState = &destinationSubscriptionResource{}
)

func NewDestinationSubscriptionResource() resource.Resource {
	return &destinationSubscriptionResource{}
}

type destinationSubscriptionResource struct {
	client      *api.APIClient
	authContext context.Context
}

func (r *destinationSubscriptionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_destination_subscription"
}

func (r *destinationSubscriptionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The unique identifier for the subscription.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"destination_id": schema.StringAttribute{
				Required:    true,
				Description: "The associated Destination instance id.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the subscription.",
			},
			"enabled": schema.BoolAttribute{
				Required:    true,
				Description: "Is the subscription enabled.",
			},
			"action_id": schema.StringAttribute{
				Required:    true,
				Description: "The unique identifier for the Destination action to trigger.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"action_slug": schema.StringAttribute{
				Computed:    true,
				Description: "The URL-friendly key for the associated Destination action.",
			},
			"trigger": schema.StringAttribute{
				Required:    true,
				Description: "FQL string that describes what events should trigger a Destination action.",
			},
			"model_id": schema.StringAttribute{
				Optional:    true,
				Description: "The unique identifier for the linked ReverseETLModel, if this part of a Reverse ETL connection.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"settings": schema.StringAttribute{
				Required:    true,
				Description: `The customer settings for action fields. Only settings included in the configuration will be managed by Terraform.`,
				CustomType:  jsontypes.NormalizedType{},
			},
		},
	}
}

func (r *destinationSubscriptionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan models.DestinationSubscriptionState
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

	out, body, err := r.client.DestinationsAPI.CreateDestinationSubscription(r.authContext, plan.DestinationID.ValueString()).CreateDestinationSubscriptionAlphaInput(api.CreateDestinationSubscriptionAlphaInput{
		Name:     plan.Name.ValueString(),
		ActionId: plan.ActionID.ValueString(),
		Trigger:  plan.Trigger.ValueString(),
		Enabled:  plan.Enabled.ValueBool(),
		ModelId:  plan.ModelID.ValueStringPointer(),
		Settings: settings,
	}).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create Destination subscription",
			getError(err, body),
		)

		return
	}

	destinationSubscription := out.Data.GetDestinationSubscription()

	var state models.DestinationSubscriptionState
	err = state.Fill(destinationSubscription)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to populate Destination subscription state",
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

func (r *destinationSubscriptionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var previousState models.DestinationSubscriptionState

	diags := req.State.Get(ctx, &previousState)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	out, body, err := r.client.DestinationsAPI.GetSubscriptionFromDestination(r.authContext, previousState.DestinationID.ValueString(), previousState.ID.ValueString()).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Destination subscription",
			getError(err, body),
		)

		return
	}

	var state models.DestinationSubscriptionState

	err = state.Fill(api.DestinationSubscription(out.Data.GetSubscription()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to populate Destination subscription state",
			err.Error(),
		)

		return
	}

	if !previousState.Settings.IsNull() && !previousState.Settings.IsUnknown() {
		state.Settings = previousState.Settings
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *destinationSubscriptionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan models.DestinationSubscriptionState
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state models.DestinationSubscriptionState
	diags = req.State.Get(ctx, &state)
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

	out, body, err := r.client.DestinationsAPI.UpdateSubscriptionForDestination(r.authContext, state.DestinationID.ValueString(), state.ID.ValueString()).UpdateSubscriptionForDestinationAlphaInput(api.UpdateSubscriptionForDestinationAlphaInput{
		Input: api.DestinationSubscriptionUpdateInput{
			Name:     plan.Name.ValueStringPointer(),
			Trigger:  plan.Trigger.ValueStringPointer(),
			Enabled:  plan.Enabled.ValueBoolPointer(),
			Settings: settings,
		},
	}).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to update Destination subscription",
			getError(err, body),
		)

		return
	}

	err = state.Fill(api.DestinationSubscription(out.Data.GetSubscription()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to populate Destination subscription state",
			err.Error(),
		)

		return
	}

	// This is to satisfy terraform requirements that the returned fields must match the input ones because new settings can be generated in the response
	state.Settings = plan.Settings

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *destinationSubscriptionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var config models.DestinationSubscriptionState
	diags := req.State.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, body, err := r.client.DestinationsAPI.RemoveSubscriptionFromDestination(r.authContext, config.DestinationID.ValueString(), config.ID.ValueString()).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to delete Destination subscription",
			getError(err, body),
		)

		return
	}
}

func (r *destinationSubscriptionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ":")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: <destination_id>:<subscription_id>. Got: %q", req.ID),
		)

		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("destination_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idParts[1])...)
}

func (r *destinationSubscriptionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
