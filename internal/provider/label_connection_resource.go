package provider

import (
	"context"
	"fmt"

	"github.com/segmentio/terraform-provider-segment/internal/provider/models"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/segmentio/public-api-sdk-go/api"
)

var (
	_ resource.Resource                = &labelAssignmentResource{}
	_ resource.ResourceWithConfigure   = &labelAssignmentResource{}
	_ resource.ResourceWithImportState = &labelAssignmentResource{}
)

func NewLabelAssignmentResource() resource.Resource {
	return &labelAssignmentResource{}
}

type labelAssignmentResource struct {
	client      *api.APIClient
	authContext context.Context
}

func (r *labelAssignmentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_label_assignment"
}

func (r *labelAssignmentResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Represents a resource-labels assignment",
		Attributes: map[string]schema.Attribute{
			"resource_id": schema.StringAttribute{
				Required:    true,
				Description: "The id of the resource to attach the label to.",
			},
			"resource_type": schema.StringAttribute{
				Required:    true,
				Description: "The type of the resource to attach the label to. Currently supports 'SOURCE'",
				Validators: []validator.String{
					stringvalidator.OneOf("SOURCE"),
				},
			},
			"labels": schema.SetNestedAttribute{
				Description: "The labels to attach to the resource.",
				Required:    true,
				Validators: []validator.Set{
					setvalidator.SizeAtMost(MaxPageSize),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"key": schema.StringAttribute{
							Required:    true,
							Description: "The key that represents the name of this label.",
						},
						"value": schema.StringAttribute{
							Required:    true,
							Description: "The value associated with the key of this label.",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "An optional description of the purpose of this label.",
						},
					},
				},
			},
		},
	}
}

func (r *labelAssignmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan models.LabelAssignmentPlan
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.ResourceType.ValueString() == "SOURCE" {
		_, body, err := r.client.SourcesApi.ReplaceLabelsInSource(r.authContext, plan.ResourceID.ValueString()).ReplaceLabelsInSourceV1Input(api.ReplaceLabelsInSourceV1Input{
			Labels: []api.LabelV1{},
		}).Execute()
		if body != nil {
			defer body.Body.Close()
		}
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to replace labels in Source",
				getError(err, body),
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

func (r *labelAssignmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var config models.LabelAssignmentPlan
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

	out, body, err := r.client.LabelAssignmentsApi.GetLabelAssignment(r.authContext, id).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Tracking Plan",
			getError(err, body),
		)

		return
	}

	labelAssignment := out.Data.GetLabelAssignment()

	var state models.LabelAssignmentState

	if !config.Rules.IsNull() && !config.Rules.IsUnknown() {
		var rules []models.RulesState
		config.Rules.ElementsAs(ctx, &rules, false)
		err = state.Fill(labelAssignment, nil)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to read Tracking Plan",
				err.Error(),
			)

			return
		}
		state.Rules = rules
	} else {
		out, body, err := r.client.LabelAssignmentsApi.ListRulesFromLabelAssignment(r.authContext, id).Pagination(*api.NewPaginationInput(MaxPageSize)).Execute()
		if body != nil {
			defer body.Body.Close()
		}
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to get Tracking Plan rules",
				getError(err, body),
			)

			return
		}

		outRules := out.Data.GetRules()
		err = state.Fill(labelAssignment, &outRules)
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

func (r *labelAssignmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan models.LabelAssignmentPlan
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var config models.LabelAssignmentState
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

	_, body, err := r.client.LabelAssignmentsApi.UpdateLabelAssignment(r.authContext, config.ID.ValueString()).UpdateLabelAssignmentV1Input(api.UpdateLabelAssignmentV1Input{
		Name:        name,
		Description: description,
	}).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to update Tracking Plan",
			getError(err, body),
		)

		return
	}

	out, body, err := r.client.LabelAssignmentsApi.GetLabelAssignment(r.authContext, config.ID.ValueString()).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Tracking Plan",
			getError(err, body),
		)

		return
	}

	labelAssignment := out.Data.GetLabelAssignment()

	var rules []models.RulesState
	plan.Rules.ElementsAs(ctx, &rules, false)

	replaceRules := []api.RuleInputV1{}
	for _, rule := range rules {
		apiRule, diags := rule.ToAPIRuleInput()
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		replaceRules = append(replaceRules, apiRule)
	}

	rulesOut := []api.RuleV1{}
	for _, rule := range rules {
		apiRule, diags := rule.ToAPIRule()
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		rulesOut = append(rulesOut, apiRule)
	}

	_, body, err = r.client.LabelAssignmentsApi.ReplaceRulesInLabelAssignment(r.authContext, out.Data.LabelAssignment.Id).ReplaceRulesInLabelAssignmentV1Input(api.ReplaceRulesInLabelAssignmentV1Input{
		Rules: replaceRules,
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

	var state models.LabelAssignmentState
	err = state.Fill(labelAssignment, &rulesOut)
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

func (r *labelAssignmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var config models.LabelAssignmentState
	diags := req.State.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, body, err := r.client.LabelAssignmentsApi.DeleteLabelAssignment(r.authContext, config.ID.ValueString()).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to delete Tracking Plan",
			getError(err, body),
		)

		return
	}
}

func (r *labelAssignmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *labelAssignmentResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
