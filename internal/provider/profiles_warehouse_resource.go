package provider

import (
	"context"
	"errors"
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
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/segmentio/public-api-sdk-go/api"
)

var (
	_ resource.Resource                = &profilesWarehouseResource{}
	_ resource.ResourceWithConfigure   = &profilesWarehouseResource{}
	_ resource.ResourceWithImportState = &profilesWarehouseResource{}
)

func NewProfilesWarehouseResource() resource.Resource {
	return &profilesWarehouseResource{}
}

type profilesWarehouseResource struct {
	client      *api.APIClient
	authContext context.Context
}

func (r *profilesWarehouseResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_profiles_warehouse"
}

func (r *profilesWarehouseResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Configures a Profiles Sync Warehouse. For more information, visit the [Segment docs](https://segment.com/docs/unify/profiles-sync/overview/).\n\n" +
			docs.GenerateImportDocs("<space_id>:<warehouse_id>", "segment_profiles_warehouse"),

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The id of the Warehouse.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"space_id": schema.StringAttribute{
				Required:    true,
				Description: "The Space id.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"metadata_id": schema.StringAttribute{
				Required:    true,
				Description: "The Warehouse metadata to use.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "An optional human-readable name for this Warehouse.",
			},
			"enabled": schema.BoolAttribute{
				Optional:    true,
				Description: "Enable to allow this Warehouse to receive data. Defaults to true.",
			},
			"schema_name": schema.StringAttribute{
				Optional:    true,
				Description: "The custom schema name that Segment uses on the Warehouse side. The space slug value is default otherwise.",
			},
			"settings": schema.StringAttribute{
				Required: true,
				Description: `A key-value object that contains instance-specific settings for a Warehouse.
			
			 Different kinds of Warehouses require different settings. The required and optional settings
			 for a Warehouse are described in the 'options' object of the associated Warehouse metadata.
			
			 You can find the full list of Warehouse metadata and related settings information in the
			 '/catalog/warehouses' endpoint.
			 
			 Only settings included in the configuration will be managed by Terraform.`,
				CustomType: jsontypes.NormalizedType{},
			},
		},
	}
}

func (r *profilesWarehouseResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan models.ProfilesWarehouseState
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

	out, body, err := r.client.ProfilesSyncAPI.CreateProfilesWarehouse(r.authContext, plan.SpaceID.ValueString()).CreateProfilesWarehouseAlphaInput(api.CreateProfilesWarehouseAlphaInput{
		Enabled:    plan.Enabled.ValueBoolPointer(),
		MetadataId: plan.MetadataID.ValueString(),
		Settings:   settings,
		Name:       plan.Name.ValueStringPointer(),
		SchemaName: plan.SchemaName.ValueStringPointer(),
	}).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create Profiles Warehouse",
			getError(err, body),
		)

		return
	}

	profilesWarehouse := out.Data.GetProfilesWarehouse()

	resp.State.SetAttribute(ctx, path.Root("id"), profilesWarehouse.Id)
	resp.State.SetAttribute(ctx, path.Root("space_id"), plan.SpaceID.ValueString())

	var state models.ProfilesWarehouseState
	err = state.Fill(profilesWarehouse)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to populate Profiles Warehouse state",
			err.Error(),
		)

		return
	}

	// This is to satisfy terraform requirements that the returned fields must match the input ones because new settings can be generated in the response.
	state.Settings = plan.Settings

	// Set state to fully populated data.
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *profilesWarehouseResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var previousState models.ProfilesWarehouseState

	diags := req.State.Get(ctx, &previousState)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	warehouse, err := findProfileWarehouse(r.authContext, r.client, previousState.ID.ValueString(), previousState.SpaceID.ValueString())
	if err != nil {
		resp.State.RemoveResource(ctx)

		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to read Profiles Warehouse (ID: %s)", previousState.ID.ValueString()),
			err.Error(),
		)

		return
	}

	if warehouse == nil {
		resp.State.RemoveResource(ctx)

		return
	}

	var state models.ProfilesWarehouseState

	err = state.Fill(*warehouse)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to populate Profiles Warehouse state",
			err.Error(),
		)

		return
	}

	if !previousState.Settings.IsNull() && !previousState.Settings.IsUnknown() {
		state.Settings = previousState.Settings
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *profilesWarehouseResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan models.ProfilesWarehouseState
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state models.ProfilesWarehouseState
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

	// Only send schemaName to API if it differs from the remote state.
	// This prevents API failures when the schema name already exists in the warehouse.
	// The Segment API fails if we send a schemaName that matches the current configuration,
	// even though it should be a no-op. This handles all cases:
	// 1. Both null/undefined: Equal() returns true, schemaName stays nil (not sent)
	// 2. Both have same value: Equal() returns true, schemaName stays nil (not sent)  
	// 3. One null, other has value: Equal() returns false, schemaName gets the plan value (sent)
	// 4. Both have different values: Equal() returns false, schemaName gets the plan value (sent)
	schemaName := determineSchemaNameForUpdate(plan.SchemaName, state.SchemaName)

	out, body, err := r.client.ProfilesSyncAPI.UpdateProfilesWarehouseForSpaceWarehouse(r.authContext, state.SpaceID.ValueString(), state.ID.ValueString()).UpdateProfilesWarehouseForSpaceWarehouseAlphaInput(api.UpdateProfilesWarehouseForSpaceWarehouseAlphaInput{
		Enabled:    plan.Enabled.ValueBoolPointer(),
		Settings:   settings,
		Name:       plan.Name.ValueStringPointer(),
		SchemaName: schemaName,
	}).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to update Profiles Warehouse (ID: %s)", plan.ID.ValueString()),
			getError(err, body),
		)

		return
	}

	warehouse := out.Data.GetProfilesWarehouse()

	err = state.Fill(warehouse)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to populate Profiles Warehouse state",
			err.Error(),
		)

		return
	}

	// This is to satisfy terraform requirements that the returned fields must match the input ones because new settings can be generated in the response.
	state.Settings = plan.Settings

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *profilesWarehouseResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var config models.ProfilesWarehouseState
	diags := req.State.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, body, err := r.client.ProfilesSyncAPI.RemoveProfilesWarehouseFromSpace(r.authContext, config.SpaceID.ValueString(), config.ID.ValueString()).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to delete Profiles Warehouse (ID: %s)", config.ID.ValueString()),
			getError(err, body),
		)

		return
	}
}

func (r *profilesWarehouseResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ":")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: <space_id>:<warehouse_id>. Got: %q", req.ID),
		)

		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("space_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idParts[1])...)
}

func (r *profilesWarehouseResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func findProfileWarehouse(authContext context.Context, client *api.APIClient, id string, spaceID string) (*api.ProfilesWarehouseAlpha, error) {
	var pageToken *string
	firstPageToken := "MA=="
	pageToken = &firstPageToken

	for pageToken != nil {
		out, body, err := client.ProfilesSyncAPI.ListProfilesWarehouseInSpace(authContext, spaceID).Pagination(api.PaginationInput{Count: MaxPageSize, Cursor: pageToken}).Execute()
		if body != nil {
			defer body.Body.Close()
		}
		if err != nil {
			return nil, errors.New(getError(err, body))
		}

		warehouses := out.Data.ProfilesWarehouses

		for _, warehouse := range warehouses {
			if warehouse.Id == id {
				return &warehouse, nil
			}
		}

		pageToken = out.Data.Pagination.Next.Get()
	}

	return nil, nil
}

// determineSchemaNameForUpdate determines whether schemaName should be sent to the API
// based on comparing the plan and state values. This prevents API failures when the
// schema name already exists in the warehouse configuration.
func determineSchemaNameForUpdate(planSchemaName, stateSchemaName types.String) *string {
	// Only send schemaName to API if it differs from the remote state.
	if !planSchemaName.Equal(stateSchemaName) {
		return planSchemaName.ValueStringPointer()
	}

	return nil
}
