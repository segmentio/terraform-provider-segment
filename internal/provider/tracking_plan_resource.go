package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/segmentio/terraform-provider-segment/internal/provider/docs"
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

var MaxRules = 2000

func (r *trackingPlanResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tracking_plan"
}

func (r *trackingPlanResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Configures a Tracking Plan. For more information, visit the [Segment docs](https://segment.com/docs/protocols/tracking-plan/create/).\n\n" +
			docs.GenerateImportDocs("<id>", "segment_tracking_plan"),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The Tracking Plan's identifier.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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
To see an exact representation of this Tracking Plan's rules, please use the data source.`,
				Validators: []validator.Set{
					setvalidator.SizeAtMost(MaxRules),
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

// trackingPlanRuleKey returns a comparable identity for a rule.
// Rules are identified by (type, key); key is optional and defaults to "".
func trackingPlanRuleKey(ruleType, key string) string {
	return ruleType + "\x00" + key
}

// ruleInputToUpsert converts a RuleInputV1 to the UpsertRuleV1 required by PATCH.
func ruleInputToUpsert(r api.RuleInputV1) api.UpsertRuleV1 {
	u := api.UpsertRuleV1{
		Type:       r.Type,
		JsonSchema: r.JsonSchema,
		Version:    r.Version,
	}
	if r.Key != nil {
		u.Key = r.Key
	}
	return u
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

	out, body, err := r.client.TrackingPlansAPI.CreateTrackingPlan(r.authContext).CreateTrackingPlanV1Input(api.CreateTrackingPlanV1Input{
		Name:        plan.Name.ValueString(),
		Type:        plan.Type.ValueString(),
		Description: description,
	}).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create Tracking Plan",
			getError(err, body),
		)

		return
	}

	trackingPlan := out.Data.GetTrackingPlan()

	resp.State.SetAttribute(ctx, path.Root("id"), trackingPlan.Id)

	var rules []models.RulesState
	plan.Rules.ElementsAs(ctx, &rules, false)

	upsertRules := []api.UpsertRuleV1{}
	rulesOut := []api.RuleV1{}
	for _, rule := range rules {
		apiRuleInput, diags := rule.ToAPIRuleInput()
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		upsertRules = append(upsertRules, ruleInputToUpsert(apiRuleInput))

		apiRule, diags := rule.ToAPIRule()
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		rulesOut = append(rulesOut, apiRule)
	}

	if len(upsertRules) > 0 {
		_, body, err = r.client.TrackingPlansAPI.UpdateRulesInTrackingPlan(r.authContext, out.Data.TrackingPlan.Id).UpdateRulesInTrackingPlanV1Input(api.UpdateRulesInTrackingPlanV1Input{
			Rules: upsertRules,
		}).Execute()
		if body != nil {
			defer body.Body.Close()
		}
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to create Tracking Plan rules",
				getError(err, body),
			)

			return
		}
	}

	var state models.TrackingPlanState
	err = state.Fill(trackingPlan, &rulesOut)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to populate Tracking Plan state",
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

	out, body, err := r.client.TrackingPlansAPI.GetTrackingPlan(r.authContext, id).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		if body.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)

			return
		}

		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to read Tracking Plan (ID: %s)", config.ID.ValueString()),
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
				"Unable to populate Tracking Plan state",
				err.Error(),
			)

			return
		}
		state.Rules = rules
	} else {
		outRules := []api.RuleV1{}

		out, body, err := r.client.TrackingPlansAPI.ListRulesFromTrackingPlan(r.authContext, id).Pagination(*api.NewPaginationInput(MaxPageSize)).Execute()
		if body != nil {
			defer body.Body.Close()
		}
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Unable to read Tracking Plan rules (ID: %s)", id),
				getError(err, body),
			)

			return
		}

		outRules = append(outRules, out.Data.GetRules()...)
		nextPointer := out.Data.GetPagination().Next.Get()

		for nextPointer != nil && len(outRules) < MaxRules {
			paginationInput := *api.NewPaginationInput(MaxPageSize)
			paginationInput.SetCursor(*nextPointer)

			out, body, err = r.client.TrackingPlansAPI.ListRulesFromTrackingPlan(r.authContext, id).Pagination(paginationInput).Execute()
			if body != nil {
				defer body.Body.Close()
			}
			if err != nil {
				resp.Diagnostics.AddError(
					fmt.Sprintf("Unable to read Tracking Plan rules (ID: %s)", id),
					getError(err, body),
				)

				return
			}

			outRules = append(outRules, out.Data.GetRules()...)
			nextPointer = out.Data.GetPagination().Next.Get()
		}

		// Limit the number of rules to MAX_RULES
		if len(outRules) > MaxRules {
			outRules = outRules[:MaxRules]
		}

		err = state.Fill(trackingPlan, &outRules)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to populate Tracking Plan state",
				err.Error(),
			)
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

	_, body, err := r.client.TrackingPlansAPI.UpdateTrackingPlan(r.authContext, config.ID.ValueString()).UpdateTrackingPlanV1Input(api.UpdateTrackingPlanV1Input{
		Name:        name,
		Description: description,
	}).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to update Tracking Plan (ID: %s)", plan.ID.ValueString()),
			getError(err, body),
		)

		return
	}

	out, body, err := r.client.TrackingPlansAPI.GetTrackingPlan(r.authContext, config.ID.ValueString()).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to read Tracking Plan (ID: %s)", config.ID.ValueString()),
			getError(err, body),
		)

		return
	}

	trackingPlan := out.Data.GetTrackingPlan()

	// Build desired rule set from plan.
	var desiredRules []models.RulesState
	plan.Rules.ElementsAs(ctx, &desiredRules, false)

	upsertRules := []api.UpsertRuleV1{}
	desiredKeys := make(map[string]struct{}, len(desiredRules))
	rulesOut := []api.RuleV1{}

	for _, rule := range desiredRules {
		key := ""
		if !rule.Key.IsNull() && !rule.Key.IsUnknown() {
			key = rule.Key.ValueString()
		}
		desiredKeys[trackingPlanRuleKey(rule.Type.ValueString(), key)] = struct{}{}

		apiRuleInput, diags := rule.ToAPIRuleInput()
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		upsertRules = append(upsertRules, ruleInputToUpsert(apiRuleInput))

		apiRule, diags := rule.ToAPIRule()
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		rulesOut = append(rulesOut, apiRule)
	}

	// PATCH: upsert all desired rules. Uses PATCH instead of PUT to avoid the
	// 200-rule replacement limit on PUT /tracking-plans/{id}/rules.
	_, body, err = r.client.TrackingPlansAPI.UpdateRulesInTrackingPlan(r.authContext, out.Data.TrackingPlan.Id).UpdateRulesInTrackingPlanV1Input(api.UpdateRulesInTrackingPlanV1Input{
		Rules: upsertRules,
	}).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to update Tracking Plan rules",
			getError(err, body),
		)

		return
	}

	// DELETE: remove rules that were in the previous state but are absent from the new plan.
	var previousRules []models.RulesState
	config.Rules.ElementsAs(ctx, &previousRules, false)

	deleteRules := []api.RemoveRuleV1{}
	for _, rule := range previousRules {
		key := ""
		if !rule.Key.IsNull() && !rule.Key.IsUnknown() {
			key = rule.Key.ValueString()
		}
		if _, exists := desiredKeys[trackingPlanRuleKey(rule.Type.ValueString(), key)]; !exists {
			removeRule := api.RemoveRuleV1{
				Type:    rule.Type.ValueString(),
				Version: float32(rule.Version.ValueFloat64()),
			}
			if key != "" {
				removeRule.Key = &key
			}
			deleteRules = append(deleteRules, removeRule)
		}
	}

	if len(deleteRules) > 0 {
		_, body, err = r.client.TrackingPlansAPI.RemoveRulesFromTrackingPlan(r.authContext, out.Data.TrackingPlan.Id).
			Rules(deleteRules).
			Execute()
		if body != nil {
			defer body.Body.Close()
		}
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to remove stale Tracking Plan rules",
				getError(err, body),
			)

			return
		}
	}

	var state models.TrackingPlanState
	err = state.Fill(trackingPlan, &rulesOut)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to populate Tracking Plan state",
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

	_, body, err := r.client.TrackingPlansAPI.DeleteTrackingPlan(r.authContext, config.ID.ValueString()).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to delete Tracking Plan (ID: %s)", config.ID.ValueString()),
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
