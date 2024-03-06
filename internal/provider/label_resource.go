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
	"github.com/segmentio/terraform-provider-segment/internal/provider/docs"
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

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("key"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("value"), idParts[1])...)
}

func (r *labelResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_label"
}

func (r *labelResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Configures a Label. For more information, visit the [Segment docs](https://segment.com/docs/segment-app/iam/labels/).\n\n" +
			docs.GenerateImportDocs("<key>:<value>", "segment_label"),
		Attributes: map[string]schema.Attribute{
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

	label := api.LabelV1{
		Key:   types.String.ValueString(plan.Key),
		Value: types.String.ValueString(plan.Value),
	}

	label.Description = types.String.ValueStringPointer(plan.Description)

	// Generate API request body from plan
	out, body, err := r.client.LabelsAPI.CreateLabel(r.authContext).CreateLabelV1Input(api.CreateLabelV1Input{
		Label: label,
	}).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create Label",
			getError(err, body),
		)

		return
	}

	outLabel := out.Data.Label
	var state models.LabelResourceState
	state.Fill(outLabel)

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

	response, body, err := r.client.LabelsAPI.ListLabels(r.authContext).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Labels",
			getError(err, body),
		)

		return
	}

	labels := response.Data.Labels

	var label *api.LabelV1
	for _, l := range labels {
		if l.Key == types.String.ValueString(state.Key) && l.Value == types.String.ValueString(state.Value) {
			label = &api.LabelV1{
				Key:         l.Key,
				Value:       l.Value,
				Description: l.Description,
			}
		}
	}

	if label == nil {
		resp.Diagnostics.AddError(
			"Unable to find Label",
			fmt.Sprintf("Unable to find Label with key: %q and value: %q", state.Key, state.Value),
		)

		return
	}

	state.Fill(*label)
	if label.Description != nil && *label.Description == "" {
		label.Description = nil
	}
	state.Description = types.StringPointerValue(label.Description)

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

	_, body, err := r.client.LabelsAPI.DeleteLabel(r.authContext, state.Key.ValueString(), state.Value.ValueString()).Execute()
	if body != nil {
		defer body.Body.Close()
	}
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
