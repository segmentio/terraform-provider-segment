package provider

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/segmentio/terraform-provider-segment/internal/provider/docs"
	"github.com/segmentio/terraform-provider-segment/internal/provider/models"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/segmentio/public-api-sdk-go/api"
)

/*──────────────────────────── helper ────────────────────────────*/

// Matches "<name> (workspace-slug)" and returns "<name>"
var wsSuffix = regexp.MustCompile(`^(.*) [^)]+$`)

func stripWorkspaceSuffix(name string) string {
	if m := wsSuffix.FindStringSubmatch(name); m != nil {
		return strings.TrimSpace(m[1])
	}
	return name
}

/*──────────────────────────── resource boilerplate ──────────────*/

var (
	_ resource.Resource                = &functionResource{}
	_ resource.ResourceWithConfigure   = &functionResource{}
	_ resource.ResourceWithImportState = &functionResource{}
)

func NewFunctionResource() resource.Resource {
	return &functionResource{}
}

type functionResource struct {
	client      *api.APIClient
	authContext context.Context
}

func (r *functionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_function"
}

/*──────────────────────────── schema ────────────────────────────*/

func (r *functionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Configures a Function. For more information, visit the [Segment docs](https://segment.com/docs/connections/functions/).\n\n" +
			docs.GenerateImportDocs("<id>", "segment_function"),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The id of the Function.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"code": schema.StringAttribute{
				Required:    true,
				Description: "The Function code.",
			},
			"display_name": schema.StringAttribute{
				Optional:    true,
				Description: "A display name for this Function (alphanumeric + spaces).",
			},
			"logo_url": schema.StringAttribute{
				Optional:    true,
				Description: "The URL of the logo for this Function.",
			},
			"resource_type": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "The Function type.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "A description for this function.",
			},
			"preview_webhook_url": schema.StringAttribute{
				Computed:    true,
				Description: "The preview webhook URL for this Function.",
			},
			"catalog_id": schema.StringAttribute{
				Computed:    true,
				Description: "The catalog id of this Function.",
			},
			"settings": schema.SetNestedAttribute{
				Optional:    true,
				Description: "Settings associated with this Function.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required:    true,
							Description: "The name of this Function setting.",
						},
						"label": schema.StringAttribute{
							Required:    true,
							Description: "The label for this Function setting.",
						},
						"description": schema.StringAttribute{
							Required:    true,
							Description: "A description of this Function setting.",
						},
						"type": schema.StringAttribute{
							Required:    true,
							Description: "The type of this Function setting.",
							Validators: []validator.String{
								stringvalidator.RegexMatches(regexp.MustCompile("^[A-Z_]+$"), "'type' must be in all uppercase"),
							},
						},
						"required": schema.BoolAttribute{
							Required:    true,
							Description: "Whether this Function setting is required.",
						},
						"sensitive": schema.BoolAttribute{
							Required:    true,
							Description: "Whether this Function setting contains sensitive information.",
						},
					},
				},
			},
		},
	}
}

/*──────────────────────────── CREATE ────────────────────────────*/

func (r *functionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan models.FunctionPlan
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	settings, diags := models.GetFunctionSettingAPIValueFromPlan(ctx, plan.Settings)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	out, body, err := r.client.FunctionsAPI.CreateFunction(r.authContext).CreateFunctionV1Input(api.CreateFunctionV1Input{
		Code:         plan.Code.ValueString(),
		Description:  plan.Description.ValueStringPointer(),
		DisplayName:  plan.DisplayName.ValueString(),
		LogoUrl:      plan.LogoURL.ValueStringPointer(),
		ResourceType: plan.ResourceType.ValueString(),
		Settings:     settings,
	}).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError("Unable to create Function", getError(err, body))
		return
	}

	function := out.Data.GetFunction()

	resp.State.SetAttribute(ctx, path.Root("id"), function.Id)

	var state models.FunctionState
	state.Fill(function)

	// Always normalise the name for workspace-scoped types
	state.DisplayName = types.StringValue(stripWorkspaceSuffix(state.DisplayName.ValueString()))

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

/*──────────────────────────── READ ───────────────────────────────*/

func (r *functionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var previousState models.FunctionState
	diags := req.State.Get(ctx, &previousState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, body, err := r.client.FunctionsAPI.GetFunction(r.authContext, previousState.ID.ValueString()).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		if body.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(fmt.Sprintf("Unable to read Function (ID: %s)", previousState.ID.ValueString()), getError(err, body))
		return
	}

	var state models.FunctionState
	state.Fill(response.Data.GetFunction())

	// normalise if necessary
	if rt := state.ResourceType.ValueString(); rt == "DESTINATION" || rt == "INSERT_DESTINATION" || rt == "INSERT_SOURCE" {
		state.DisplayName = types.StringValue(stripWorkspaceSuffix(state.DisplayName.ValueString()))
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

/*──────────────────────────── UPDATE ────────────────────────────*/

func (r *functionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan models.FunctionPlan
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state models.FunctionState
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	settings, diags := models.GetFunctionSettingAPIValueFromPlan(ctx, plan.Settings)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	out, body, err := r.client.FunctionsAPI.UpdateFunction(r.authContext, state.ID.ValueString()).UpdateFunctionV1Input(api.UpdateFunctionV1Input{
		Code:        plan.Code.ValueStringPointer(),
		Description: plan.Description.ValueStringPointer(),
		DisplayName: plan.DisplayName.ValueStringPointer(),
		LogoUrl:     plan.LogoURL.ValueStringPointer(),
		Settings:    settings,
	}).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Unable to update Function (ID: %s)", plan.ID.ValueString()), getError(err, body))
		return
	}

	state.Fill(out.Data.GetFunction())

	if rt := state.ResourceType.ValueString(); rt == "DESTINATION" || rt == "INSERT_DESTINATION" || rt == "INSERT_SOURCE" {
		state.DisplayName = types.StringValue(stripWorkspaceSuffix(state.DisplayName.ValueString()))
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

/*──────────────────────────── DELETE ────────────────────────────*/

func (r *functionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var config models.FunctionState
	diags := req.State.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, body, err := r.client.FunctionsAPI.DeleteFunction(r.authContext, config.ID.ValueString()).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Unable to delete Function (ID: %s)", config.ID.ValueString()), getError(err, body))
	}
}

/*──────────────────────────── IMPORT ────────────────────────────*/

func (r *functionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

/*──────────────────────────── CONFIGURE ─────────────────────────*/

func (r *functionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
