package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/segmentio/public-api-sdk-go/api"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &labelResource{}
	_ resource.ResourceWithConfigure = &labelResource{}
)

// NewLabelResource is a helper function to simplify the provider implementation.
func NewLabelResource() resource.Resource {
	return &labelResource{}
}

// labelResource is the resource implementation.
type labelResource struct {
	client      *api.APIClient
	authContext context.Context
}

type labelResourceModel struct {
	Id    types.String `tfsdk:"id"`
	Label labelModel   `tfsdk:"label"`
}

type labelModel struct {
	Key         types.String `tfsdk:"key"`
	Value       types.String `tfsdk:"value"`
	Description types.String `tfsdk:"description"`
}

// Metadata returns the resource type name.
func (r *labelResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_label"
}

// Schema defines the schema for the resource.
func (r *labelResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"label": schema.SingleNestedAttribute{
				Description: "A label associated with the current Workspace.",
				Required:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
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
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *labelResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan labelResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	label := api.Label{
		Key:   types.String.ValueString(plan.Label.Key),
		Value: types.String.ValueString(plan.Label.Value),
	}

	label.Description = types.String.ValueStringPointer(plan.Label.Description)

	// Generate API request body from plan
	response, _, err := r.client.LabelsApi.CreateLabel(r.authContext).CreateLabelV1Input(api.CreateLabelV1Input{
		Label: label,
	}).Execute()

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create a label",
			err.Error(),
		)
		return
	}

	plan.Label.Key = types.StringValue(response.Data.Label.Key)
	plan.Label.Value = types.StringValue(response.Data.Label.Value)
	plan.Id = types.StringValue("placeholder")

	if response.Data.Label.Description != nil {
		plan.Label.Description = types.StringPointerValue(response.Data.Label.Description)
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *labelResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state labelResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, _, err := r.client.LabelsApi.ListLabels(r.authContext).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Labels",
			err.Error(),
		)
		return
	}

	labels := response.Data.Labels

	label := api.LabelV1{}
	for _, l := range labels {
		if l.Key == types.String.ValueString(state.Label.Key) && l.Value == types.String.ValueString(state.Label.Value) {
			label = l
		}
	}

	state.Label.Key = types.StringValue(label.Key)
	state.Label.Value = types.StringValue(label.Value)
	state.Id = types.StringValue("placeholder")

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *labelResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Label does not have an update functionality, so added RequiresReplace to each attribute that can be configurable.
	// reference: https://developer.hashicorp.com/terraform/plugin/framework/resources/update#caveats
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *labelResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state labelResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, _, err := r.client.LabelsApi.DeleteLabel(ctx, state.Label.Key.ValueString(), state.Label.Value.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting a Label", "Could not delete a label, unexpected error: "+err.Error(),
		)
		return
	}
}

// Configure adds the provider configured client to the resource.
func (r *labelResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	config, ok := req.ProviderData.(*ClientInfo)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *hashicups.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = config.client
	r.authContext = config.authContext
}
