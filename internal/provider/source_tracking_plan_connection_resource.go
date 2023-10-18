package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/segmentio/public-api-sdk-go/api"
)

var (
	_ resource.Resource              = &sourceTrackingPlanConnectionResource{}
	_ resource.ResourceWithConfigure = &sourceTrackingPlanConnectionResource{}
)

func NewSourceTrackingPlanConnectionResource() resource.Resource {
	return &sourceTrackingPlanConnectionResource{}
}

type sourceTrackingPlanConnectionResource struct {
	client      *api.APIClient
	authContext context.Context
}

type sourceTrackingPlanConnectionState struct {
	SourceID       types.String `tfsdk:"source_id"`
	TrackingPlanID types.String `tfsdk:"tracking_plan_id"`
}

func (r *sourceTrackingPlanConnectionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_source_tracking_plan_connection"
}

func (r *sourceTrackingPlanConnectionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Represents a connection between a Source and a Tracking Plan",
		Attributes: map[string]schema.Attribute{
			"source_id": schema.StringAttribute{
				Required:    true,
				Description: "The id of the Source.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"tracking_plan_id": schema.StringAttribute{
				Required:    true,
				Description: "The id of the Tracking Plan.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *sourceTrackingPlanConnectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan sourceTrackingPlanConnectionState
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.TrackingPlanID.String() == "" || plan.SourceID.String() == "" {
		resp.Diagnostics.AddError("Unable to create connection between Source and Tracking Plan", "At least one ID is empty")

		return
	}

	_, body, err := r.client.TrackingPlansApi.AddSourceToTrackingPlan(r.authContext, plan.TrackingPlanID.ValueString()).AddSourceToTrackingPlanV1Input(api.AddSourceToTrackingPlanV1Input{
		SourceId: plan.SourceID.ValueString(),
	}).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create connection between Source and Tracking Plan",
			getError(err, body),
		)

		return
	}

	state := sourceTrackingPlanConnectionState{
		SourceID:       plan.SourceID,
		TrackingPlanID: plan.TrackingPlanID,
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *sourceTrackingPlanConnectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state sourceTrackingPlanConnectionState

	diags := req.State.Get(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paginationNext := "MA=="

	for paginationNext != "" {
		if state.SourceID.String() == "" {
			resp.Diagnostics.AddError("Unable to read Source-Tracking Plan connection", "At least one ID is empty")

			return
		}
		response, body, err := r.client.TrackingPlansApi.ListSourcesFromTrackingPlan(r.authContext, state.TrackingPlanID.ValueString()).Pagination(api.PaginationInput{
			Cursor: &paginationNext,
			Count:  MaxPageSize,
		}).Execute()
		if body != nil {
			defer body.Body.Close()
		}
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to read Source-Tracking Plan connection",
				getError(err, body),
			)

			return
		}

		for _, source := range response.Data.Sources {
			if source.Id == state.SourceID.ValueString() {
				diags = resp.State.Set(ctx, &state)
				resp.Diagnostics.Append(diags...)
				if resp.Diagnostics.HasError() {
					return
				}

				return
			}
		}

		if response.Data.Pagination.Next.IsSet() {
			paginationNext = *response.Data.Pagination.Next.Get()
		} else {
			paginationNext = ""
		}
	}

	diags = resp.State.Set(ctx, &sourceTrackingPlanConnectionState{
		SourceID:       types.StringValue("not_found"),
		TrackingPlanID: types.StringValue("not_found"),
	})
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *sourceTrackingPlanConnectionResource) Update(context.Context, resource.UpdateRequest, *resource.UpdateResponse) {
	// All fields force replacement
}

func (r *sourceTrackingPlanConnectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var config sourceTrackingPlanConnectionState
	diags := req.State.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.TrackingPlanID.String() == "" || config.SourceID.String() == "" {
		resp.Diagnostics.AddError("Unable to remove Source-Tracking Plan connection", "At least one ID is empty")

		return
	}

	_, body, err := r.client.TrackingPlansApi.RemoveSourceFromTrackingPlan(r.authContext, config.TrackingPlanID.ValueString()).SourceId(config.SourceID.ValueString()).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to remove Source-Tracking Plan connection",
			getError(err, body),
		)

		return
	}
}

func (r *sourceTrackingPlanConnectionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
