package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/segmentio/terraform-provider-segment/internal/provider/docs"
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
	_ resource.Resource                = &insertFunctionInstanceResource{}
	_ resource.ResourceWithConfigure   = &insertFunctionInstanceResource{}
	_ resource.ResourceWithImportState = &insertFunctionInstanceResource{}
)

func NewInsertFunctionInstanceResource() resource.Resource {
	return &insertFunctionInstanceResource{}
}

type insertFunctionInstanceResource struct {
	client      *api.APIClient
	authContext context.Context
}

func (r *insertFunctionInstanceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_insert_function_instance"
}

func (r *insertFunctionInstanceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Configures an Insert Function. For more information, visit the [Segment docs](https://segment.com/docs/connections/functions/insert-functions/).\n\n" +
			docs.GenerateImportDocs("<id>", "segment_insert_function_instance"),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The unique identifier for the insert function instance.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"function_id": schema.StringAttribute{
				Required:    true,
				Description: "Insert Function id to which this instance is associated.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"integration_id": schema.StringAttribute{
				Required:    true,
				Description: "The Source or Destination id to be connected.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Defines the display name of the insert Function instance.",
			},
			"enabled": schema.BoolAttribute{
				Required:    true,
				Description: "Whether this insert Function instance should be enabled for the Destination.",
			},
			"settings": schema.StringAttribute{
				Required:    true,
				Description: `An object that contains settings for this insert Function instance based on the settings present in the insert Function class. Only settings included in the configuration will be managed by Terraform.`,
				CustomType:  jsontypes.NormalizedType{},
			},
		},
	}
}

func (r *insertFunctionInstanceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan models.InsertFunctionInstanceState
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

	enabled := plan.Enabled.ValueBool()
	out, body, err := r.client.FunctionsAPI.CreateInsertFunctionInstance(r.authContext).CreateInsertFunctionInstanceAlphaInput(api.CreateInsertFunctionInstanceAlphaInput{
		Name:          plan.Name.ValueString(),
		FunctionId:    strings.TrimPrefix(plan.FunctionID.ValueString(), "ifnd_"),
		IntegrationId: plan.IntegrationID.ValueString(),
		Enabled:       &enabled,
		Settings:      settings,
	}).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create Insert Function instance",
			getError(err, body),
		)

		return
	}

	insertFunctionInstance := out.Data.InsertFunctionInstance

	resp.State.SetAttribute(ctx, path.Root("id"), insertFunctionInstance.Id)

	var state models.InsertFunctionInstanceState
	err = state.Fill(insertFunctionInstance)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to populate Insert Function instance state",
			err.Error(),
		)

		return
	}

	// This is to satisfy terraform requirements that the returned fields must match the input ones because new settings can be generated in the response
	state.Settings = plan.Settings

	// This is to satisfy terraform requirements that the input fields must match the returned ones. The input FunctionID can be prefixed with "ifnd_" and the returned one is not.
	state.FunctionID = plan.FunctionID

	// Set state to fully populated data
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *insertFunctionInstanceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var previousState models.InsertFunctionInstanceState

	diags := req.State.Get(ctx, &previousState)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	out, body, err := r.client.FunctionsAPI.GetInsertFunctionInstance(r.authContext, previousState.ID.ValueString()).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		if body.StatusCode == 404 {
			resp.State.RemoveResource(ctx)

			return
		}

		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to read Insert Function instance (ID: %s)", previousState.ID.ValueString()),
			getError(err, body),
		)

		return
	}

	var state models.InsertFunctionInstanceState

	err = state.Fill(out.Data.InsertFunctionInstance)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to populate Insert Function instance state",
			err.Error(),
		)

		return
	}

	if !previousState.Settings.IsNull() && !previousState.Settings.IsUnknown() {
		state.Settings = previousState.Settings
	}

	// This is to satisfy terraform requirements that the input fields must match the returned ones. The input FunctionID can be prefixed with "ifnd_" and the returned one is not.
	if !previousState.Settings.IsNull() && !previousState.Settings.IsUnknown() {
		state.FunctionID = previousState.FunctionID
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *insertFunctionInstanceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan models.InsertFunctionInstanceState
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state models.InsertFunctionInstanceState
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

	out, body, err := r.client.FunctionsAPI.UpdateInsertFunctionInstance(r.authContext, state.ID.ValueString()).UpdateInsertFunctionInstanceAlphaInput(api.UpdateInsertFunctionInstanceAlphaInput{
		Enabled:  plan.Enabled.ValueBoolPointer(),
		Name:     plan.Name.ValueStringPointer(),
		Settings: settings,
	}).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to update Insert Function instance (ID: %s)", plan.ID.ValueString()),
			getError(err, body),
		)

		return
	}

	err = state.Fill(out.Data.InsertFunctionInstance)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to populate Insert Function instance state",
			err.Error(),
		)

		return
	}

	// This is to satisfy terraform requirements that the returned fields must match the input ones because new settings can be generated in the response
	state.Settings = plan.Settings

	// This is to satisfy terraform requirements that the input fields must match the returned ones. The input FunctionID can be prefixed with "ifnd_" and the returned one is not.
	state.FunctionID = plan.FunctionID

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *insertFunctionInstanceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var config models.InsertFunctionInstanceState
	diags := req.State.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, body, err := r.client.FunctionsAPI.DeleteInsertFunctionInstance(r.authContext, config.ID.ValueString()).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to delete Insert Function instance (ID: %s)", config.ID.ValueString()),
			getError(err, body),
		)

		return
	}
}

func (r *insertFunctionInstanceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *insertFunctionInstanceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
