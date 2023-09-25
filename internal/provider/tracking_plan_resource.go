package provider

import (
	"context"
	"fmt"

	"github.com/segmentio/terraform-provider-segment/internal/provider/models"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
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
			"rules": schema.SetNestedAttribute{
				Required: true,
				Description: `The list of Tracking Plan rules. 
				
Due to Terraform resource limitations, this list might not show an exact representation of how the Tracking Plan interprets each rule.
To see an exact representation of this Tracking Plan's rules, please use the data source.

This field is currently limited to 200 items.`,
				Validators: []validator.Set{
					setvalidator.SizeAtMost(MaxPageSize),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Required: true,
							Description: `The type for this Tracking Plan rule.

							Enum: "COMMON" "GROUP" "IDENTIFY" "PAGE" "SCREEN" "TRACK"`,
						},
						"key": schema.StringAttribute{
							Optional:    true,
							Description: "Key to this rule (free-form string like 'Button clicked').",
						},
						"json_schema": schema.StringAttribute{
							Required:    true,
							Description: "JSON Schema of this rule.",
							CustomType:  jsontypes.NormalizedType{},
						},
						"version": schema.Float64Attribute{
							Required:    true,
							Description: "Version of this rule.",
						},
					},
				},
			},
		},
	}
}

func (r *trackingPlanResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan models.TrackingPlanPlan
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
	defer body.Body.Close()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create Tracking Plan",
			getError(err, body),
		)

		return
	}

	trackingPlan := out.Data.GetTrackingPlan()

	var rules []models.RulesState
	plan.Rules.ElementsAs(ctx, &rules, false)

	replaceRules := []api.RuleV1{}
	for _, rule := range rules {
		apiRule, diags := rule.ToAPIRule()
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		replaceRules = append(replaceRules, apiRule)
	}

	_, body, err = r.client.TrackingPlansApi.ReplaceRulesInTrackingPlan(r.authContext, out.Data.TrackingPlan.Id).ReplaceRulesInTrackingPlanV1Input(api.ReplaceRulesInTrackingPlanV1Input{
		Rules: replaceRules,
	}).Execute()
	defer body.Body.Close()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create Tracking Plan rules",
			getError(err, body),
		)

		return
	}

	var state models.TrackingPlanState
	err = state.Fill(api.TrackingPlan(trackingPlan), &replaceRules)
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
	var config models.TrackingPlanPlan
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
	defer body.Body.Close()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Tracking Plan",
			getError(err, body),
		)

		return
	}

	trackingPlan := out.Data.GetTrackingPlan()

	var state models.TrackingPlanState

	if !config.Rules.IsNull() && !config.Rules.IsUnknown() {
		var rules []models.RulesState
		config.Rules.ElementsAs(ctx, &rules, false)
		err = state.Fill(trackingPlan, nil)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to read Tracking Plan",
				err.Error(),
			)

			return
		}
		state.Rules = rules
	} else {
		out, body, err := r.client.TrackingPlansApi.ListRulesFromTrackingPlan(r.authContext, id).Pagination(*api.NewPaginationInput(MaxPageSize)).Execute()
		defer body.Body.Close()
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to get Tracking Plan rules",
				getError(err, body),
			)

			return
		}

		outRules := out.Data.GetRules()
		err = state.Fill(trackingPlan, &outRules)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to get Tracking Plan rules",
				err.Error(),
			)

			return
		}
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *trackingPlanResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan models.TrackingPlanPlan
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
	defer body.Body.Close()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to update Tracking Plan",
			getError(err, body),
		)

		return
	}

	out, body, err := r.client.TrackingPlansApi.GetTrackingPlan(r.authContext, config.ID.ValueString()).Execute()
	defer body.Body.Close()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Tracking Plan",
			getError(err, body),
		)

		return
	}

	trackingPlan := out.Data.GetTrackingPlan()

	var rules []models.RulesState
	plan.Rules.ElementsAs(ctx, &rules, false)

	replaceRules := []api.RuleV1{}
	for _, rule := range rules {
		apiRule, diags := rule.ToAPIRule()
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		replaceRules = append(replaceRules, apiRule)
	}

	_, body, err = r.client.TrackingPlansApi.ReplaceRulesInTrackingPlan(r.authContext, out.Data.TrackingPlan.Id).ReplaceRulesInTrackingPlanV1Input(api.ReplaceRulesInTrackingPlanV1Input{
		Rules: replaceRules,
	}).Execute()
	defer body.Body.Close()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create Tracking Plan rules",
			getError(err, body),
		)

		return
	}

	var state models.TrackingPlanState
	err = state.Fill(trackingPlan, &replaceRules)
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
	defer body.Body.Close()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to delete Tracking Plan",
			getError(err, body),
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
