package provider

import (
	"context"
	"fmt"

	"terraform-provider-segment/internal/provider/models"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/segmentio/public-api-sdk-go/api"
)

var (
	_ resource.Resource                = &trackingPlanResource{}
	_ resource.ResourceWithConfigure   = &trackingPlanResource{}
	_ resource.ResourceWithImportState = &trackingPlanResource{}
)

func NewTrackingPlanResource() resource.Resource {
	return &trackingPlanResource{}
}

type trackingPlanResource struct {
	client      *api.APIClient
	authContext context.Context
}

func (r *trackingPlanResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tracking_plan"
}

func (r *trackingPlanResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The Tracking Plan's identifier.",
			},
			"slug": schema.StringAttribute{
				Computed:    true,
				Description: "URL-friendly slug of this Tracking Plan.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The Tracking Plan's name.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The Tracking Plan's description.",
			},
			"type": schema.StringAttribute{
				Required:    true,
				Description: "The Tracking Plan's type.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"updated_at": schema.StringAttribute{
				Computed:    true,
				Description: "The timestamp of the last change to the Tracking Plan.",
			},
			"created_at": schema.StringAttribute{
				Computed:    true,
				Description: "The timestamp of this Tracking Plan's creation.",
			},
		},
	}
}

func (r *trackingPlanResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan models.TrackingPlanState
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var description *string
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() && plan.Description.ValueString() != "" {
		description = plan.Description.ValueStringPointer()
	}

	out, body, err := r.client.TrackingPlansApi.CreateTrackingPlan(r.authContext).CreateTrackingPlanV1Input(api.CreateTrackingPlanV1Input{
		Name:        plan.Name.ValueString(),
		Type:        plan.Type.ValueString(),
		Description: description,
	}).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create Tracking Plan",
			getError(err, body.Body),
		)
		return
	}

	trackingPlan := out.Data.GetTrackingPlan()

	var state models.TrackingPlanState
	err = state.Fill(api.TrackingPlan(trackingPlan))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create Tracking Plan",
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

func (r *trackingPlanResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var config models.TrackingPlanState
	diags := req.State.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := config.ID.ValueString()
	if id == "" {
		resp.Diagnostics.AddError("Unable to read Tracking Plan", "ID is empty")
		return
	}

	out, body, err := r.client.TrackingPlansApi.GetTrackingPlan(r.authContext, id).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Tracking Plan",
			getError(err, body.Body),
		)
		return
	}

	trackingPlan := out.Data.GetTrackingPlan()

	var state models.TrackingPlanState
	err = state.Fill(trackingPlan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Tracking Plan",
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

func (r *trackingPlanResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan models.TrackingPlanState
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var config models.TrackingPlanState
	diags = req.State.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var name *string
	if !plan.Name.IsNull() && !plan.Name.IsUnknown() && plan.Name.ValueString() != "" {
		name = plan.Name.ValueStringPointer()
	}

	var description *string
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() && plan.Description.ValueString() != "" {
		description = plan.Description.ValueStringPointer()
	}

	_, body, err := r.client.TrackingPlansApi.UpdateTrackingPlan(r.authContext, config.ID.ValueString()).UpdateTrackingPlanV1Input(api.UpdateTrackingPlanV1Input{
		Name:        name,
		Description: description,
	}).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to update Tracking Plan",
			getError(err, body.Body),
		)
		return
	}

	out, body, err := r.client.TrackingPlansApi.GetTrackingPlan(r.authContext, config.ID.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Tracking Plan",
			getError(err, body.Body),
		)
		return
	}

	trackingPlan := out.Data.GetTrackingPlan()

	var state models.TrackingPlanState
	err = state.Fill(trackingPlan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Tracking Plan",
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

func (r *trackingPlanResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var config models.TrackingPlanState
	diags := req.State.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, body, err := r.client.TrackingPlansApi.DeleteTrackingPlan(r.authContext, config.ID.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to delete Tracking Plan",
			getError(err, body.Body),
		)
		return
	}
}

func (r *trackingPlanResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *trackingPlanResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
