package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/segmentio/public-api-sdk-go/api"
	"github.com/segmentio/terraform-provider-segment/internal/provider/models"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var (
	_ resource.Resource                = &labelResource{}
	_ resource.ResourceWithConfigure   = &labelResource{}
	_ resource.ResourceWithImportState = &labelResource{}
)

func NewLabelResource() resource.Resource {
	return &labelResource{}
}

type labelResource struct {
	client      *api.APIClient
	authContext context.Context
}

func (r *labelResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ":")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: <key>:<value>. Got: %q", req.ID),
		)

		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("key"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("value"), idParts[1])...)
}

func (r *labelResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_label"
}

func (r *labelResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A label associated with the current Workspace. To import a label into Terraform, use the following format: 'key:value'",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The unique identifier for this label.",
			},
			"key": schema.StringAttribute{
				Description: "The key that represents the name of this label.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"value": schema.StringAttribute{
				Description: "The value associated with the key of this label.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Description: "An optional description of the purpose of this label.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *labelResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan models.LabelResourceState
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	label := api.Label{
		Key:   types.String.ValueString(plan.Key),
		Value: types.String.ValueString(plan.Value),
	}

	label.Description = types.String.ValueStringPointer(plan.Description)

	// Generate API request body from plan
	out, body, err := r.client.LabelsApi.CreateLabel(r.authContext).CreateLabelV1Input(api.CreateLabelV1Input{
		Label: label,
	}).Execute()
	defer body.Body.Close()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create a label",
			getError(err, body),
		)

		return
	}

	outLabel := out.Data.Label
	plan.Key = types.StringValue(outLabel.Key)
	plan.Value = types.StringValue(outLabel.Value)
	plan.ID = types.StringValue(id(outLabel.Key, outLabel.Value))

	if outLabel.Description != nil {
		plan.Description = types.StringValue(*outLabel.Description)
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *labelResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.LabelResourceState
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, body, err := r.client.LabelsApi.ListLabels(r.authContext).Execute()
	defer body.Body.Close()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Labels",
			getError(err, body),
		)

		return
	}

	labels := response.Data.Labels

	label := api.LabelV1{}
	for _, l := range labels {
		if l.Key == types.String.ValueString(state.Key) && l.Value == types.String.ValueString(state.Value) {
			label = l
		}
	}

	state.Key = types.StringValue(label.Key)
	state.Value = types.StringValue(label.Value)
	if label.Description != nil {
		state.Description = types.StringValue(*label.Description)
	}
	state.ID = types.StringValue(id(label.Key, label.Value))

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *labelResource) Update(context.Context, resource.UpdateRequest, *resource.UpdateResponse) {
	// Label does not have an update functionality, so added RequiresReplace to each attribute that can be configurable.
	// reference: https://developer.hashicorp.com/terraform/plugin/framework/resources/update#caveats
}

func (r *labelResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state models.LabelResourceState
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, body, err := r.client.LabelsApi.DeleteLabel(r.authContext, state.Key.ValueString(), state.Value.ValueString()).Execute()
	defer body.Body.Close()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to delete Label",
			getError(err, body),
		)

		return
	}
}

func (r *labelResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	config, ok := req.ProviderData.(*ClientInfo)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected ClientInfo, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = config.client
	r.authContext = config.authContext
}

func id(key, value string) string {
	return fmt.Sprintf("%s:%s", key, value)
}
