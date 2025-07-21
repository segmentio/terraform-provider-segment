package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/segmentio/public-api-sdk-go/api"
	"github.com/segmentio/terraform-provider-segment/internal/provider/models"
)

var (
	_ resource.Resource                = &audienceResource{}
	_ resource.ResourceWithConfigure   = &audienceResource{}
	_ resource.ResourceWithImportState = &audienceResource{}
)

func NewAudienceResource() resource.Resource {
	return &audienceResource{}
}

type audienceResource struct {
	client      *api.APIClient
	authContext context.Context
}

func (r *audienceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_audience"
}

func (r *audienceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The Audience ID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"space_id": schema.StringAttribute{
				Required:    true,
				Description: "The Space ID.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the Audience.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "The description of the Audience.",
			},
			"key": schema.StringAttribute{
				Optional:    true,
				Description: "The key of the Audience.",
			},
			"enabled": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether the Audience is enabled.",
			},
			"definition": schema.MapAttribute{
				ElementType: types.StringType,
				Required:    true,
				Description: "The definition of the Audience.",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "The status of the Audience.",
			},
			"options": schema.MapAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Additional options for the Audience.",
			},
		},
	}
}

// TODO: Implement Create, Read, Update, Delete, ImportState, Configure methods using the SDK

func (r *audienceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan models.AudienceState
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := api.CreateAudienceInput{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		Key:         plan.Key.ValueString(),
		Enabled:     plan.Enabled.ValueBool(),
		Options:     nil, // TODO: convert plan.Options to map[string]interface{}
		Definition:  nil, // TODO: convert plan.Definition to map[string]interface{}
	}
	// Convert plan.Definition and plan.Options from types.Map to map[string]interface{}
	if !plan.Definition.IsNull() && !plan.Definition.IsUnknown() {
		var def map[string]interface{}
		diags := plan.Definition.ElementsAs(ctx, &def, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		input.Definition = def
	}
	if !plan.Options.IsNull() && !plan.Options.IsUnknown() {
		var opts map[string]interface{}
		diags := plan.Options.ElementsAs(ctx, &opts, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		input.Options = opts
	}

	out, _, err := r.client.CreateAudience(r.authContext, plan.SpaceID.ValueString(), input)
	if err != nil {
		resp.Diagnostics.AddError("Unable to create Audience", err.Error())
		return
	}
	audience := out.Data.Audience
	plan.ID = types.StringValue(audience.Id)
	plan.Status = types.StringValue(audience.Status)
	// Set other fields as needed
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *audienceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.AudienceState
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	out, _, err := r.client.GetAudience(r.authContext, state.SpaceID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to read Audience", err.Error())
		return
	}
	audience := out.Data.Audience
	state.Name = types.StringValue(audience.Name)
	state.Description = types.StringValue(audience.Description)
	state.Key = types.StringValue(audience.Key)
	state.Enabled = types.BoolValue(audience.Enabled)
	state.Status = types.StringValue(audience.Status)
	// TODO: convert audience.Definition and audience.Options to types.Map
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *audienceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan models.AudienceState
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	input := api.UpdateAudienceInput{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		Key:         plan.Key.ValueString(),
		Enabled:     plan.Enabled.ValueBool(),
		Options:     nil, // TODO: convert plan.Options to map[string]interface{}
		Definition:  nil, // TODO: convert plan.Definition to map[string]interface{}
	}
	if !plan.Definition.IsNull() && !plan.Definition.IsUnknown() {
		var def map[string]interface{}
		diags := plan.Definition.ElementsAs(ctx, &def, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		input.Definition = def
	}
	if !plan.Options.IsNull() && !plan.Options.IsUnknown() {
		var opts map[string]interface{}
		diags := plan.Options.ElementsAs(ctx, &opts, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		input.Options = opts
	}
	out, _, err := r.client.UpdateAudience(r.authContext, plan.SpaceID.ValueString(), plan.ID.ValueString(), input)
	if err != nil {
		resp.Diagnostics.AddError("Unable to update Audience", err.Error())
		return
	}
	audience := out.Data.Audience
	plan.Status = types.StringValue(audience.Status)
	// Set other fields as needed
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *audienceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state models.AudienceState
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	_, err := r.client.DeleteAudience(r.authContext, state.SpaceID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to delete Audience", err.Error())
		return
	}
}

func (r *audienceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Expect import ID as <space_id>:<audience_id>
	idParts := strings.Split(req.ID, ":")
	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: <space_id>:<audience_id>. Got: %q", req.ID),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("space_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idParts[1])...)
}

func (r *audienceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
