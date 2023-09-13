package provider

import (
	"context"
	"fmt"

	"terraform-provider-segment/internal/provider/models"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/segmentio/public-api-sdk-go/api"
)

var (
	_ resource.Resource              = &trackingPlanRulesResource{}
	_ resource.ResourceWithConfigure = &trackingPlanRulesResource{}
)

func NewTrackingPlanRulesResource() resource.Resource {
	return &trackingPlanRulesResource{}
}

type trackingPlanRulesResource struct {
	client      *api.APIClient
	authContext context.Context
}

func (r *trackingPlanRulesResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tracking_plan_rules"
}

func (r *trackingPlanRulesResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"tracking_plan_id": schema.StringAttribute{
				Required:    true,
				Description: "The Tracking Plan's identifier.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"rules": schema.SetNestedAttribute{
				Required: true,
				Description: `The list of Tracking Plan rules. 
				
Due to Terraform resource limitations, this list might not show an exact representation of how the Tracking Plan interprets each rule.
To see an exact representation of this Tracking Plan's rules, please use the data source.

This field is currently limited to 200 items.`,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplace(),
				},
				Validators: []validator.Set{
					setvalidator.SizeAtMost(200),
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

func (r *trackingPlanRulesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan models.TrackingPlanRulesPlan
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var rules []models.RulesState
	plan.Rules.ElementsAs(ctx, &rules, false)

	replaceRules := []api.RuleV1{}
	for _, rule := range rules {
		apiRule, diags := rule.ToApiRule()
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		replaceRules = append(replaceRules, apiRule)
	}

	_, body, err := r.client.TrackingPlansApi.ReplaceRulesInTrackingPlan(r.authContext, plan.TrackingPlanID.ValueString()).ReplaceRulesInTrackingPlanV1Input(api.ReplaceRulesInTrackingPlanV1Input{
		Rules: replaceRules,
	}).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create Tracking Plan rules",
			getError(err, body),
		)
		return
	}

	var state models.TrackingPlanRulesState
	state.TrackingPlanID = plan.TrackingPlanID
	state.Rules = rules

	// Set state to fully populated data
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *trackingPlanRulesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var config models.TrackingPlanRulesPlan
	diags := req.State.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state models.TrackingPlanRulesState

	if !config.Rules.IsNull() && !config.Rules.IsUnknown() {
		state.TrackingPlanID = config.TrackingPlanID

		var rules []models.RulesState
		config.Rules.ElementsAs(ctx, &rules, false)
		state.Rules = rules
	} else {
		out, body, err := r.client.TrackingPlansApi.ListRulesFromTrackingPlan(r.authContext, config.TrackingPlanID.ValueString()).Pagination(*api.NewPaginationInput(200)).Execute()
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to get Tracking Plan rules",
				getError(err, body),
			)
			return
		}

		err = state.Fill(out.Data.GetRules(), config.TrackingPlanID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to get Tracking Plan rules",
				err.Error(),
			)
			return
		}
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *trackingPlanRulesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Everything forces replacement for simplicity sake
}

func (r *trackingPlanRulesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var config models.TrackingPlanRulesState
	diags := req.State.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, body, err := r.client.TrackingPlansApi.ReplaceRulesInTrackingPlan(r.authContext, config.TrackingPlanID.ValueString()).ReplaceRulesInTrackingPlanV1Input(api.ReplaceRulesInTrackingPlanV1Input{Rules: []api.RuleV1{}}).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to delete Tracking Plan rules",
			getError(err, body),
		)
		return
	}
}

func (r *trackingPlanRulesResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
