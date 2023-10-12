package provider

import (
	"context"
	"fmt"

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
	_ resource.Resource                = &profilesWarehouseResource{}
	_ resource.ResourceWithConfigure   = &profilesWarehouseResource{}
	_ resource.ResourceWithImportState = &profilesWarehouseResource{}
)

func NewprofilesWarehouseResource() resource.Resource {
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
			},
			"metadata_id": schema.StringAttribute{
				Required:    true,
				Description: "The Warehouse metadata to use.",
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
			 '/catalog/warehouses' endpoint.`,
				CustomType: jsontypes.NormalizedType{},
			},
		},
	}
}

func (r *profilesWarehouseResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan models.ProfilesWarehousePlan
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	wrappedMetadataID, err := plan.Metadata.Attributes()["id"].ToTerraformValue(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to decode metadata id",
			err.Error(),
		)

		return
	}

	var metadataID string
	err = wrappedMetadataID.As(&metadataID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to decode metadata id",
			err.Error(),
		)

		return
	}

	var settings map[string]interface{}
	diags = plan.Settings.Unmarshal(&settings)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	modelMap := api.NewModelMap(settings)

	name := plan.Name.ValueStringPointer()
	if *name == "" {
		name = nil
	}

	out, body, err := r.client.ProfilesWarehousesApi.CreateProfilesWarehouse(r.authContext).CreateProfilesWarehouseV1Input(api.CreateProfilesWarehouseV1Input{
		Enabled:    plan.Enabled.ValueBoolPointer(),
		MetadataId: metadataID,
		Settings:   *api.NewNullableModelMap(modelMap),
		Name:       name,
	}).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create ProfilesWarehouse",
			getError(err, body),
		)

		return
	}

	profileswarehouse := out.Data.GetProfilesWarehouse()

	var state models.ProfilesWarehouseState
	err = state.Fill(api.ProfilesWarehouse(profileswarehouse))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create ProfilesWarehouse",
			err.Error(),
		)

		return
	}

	// This is to satisfy terraform requirements that the returned fields must match the input ones because new settings can be generated in the response
	state.Settings = plan.Settings

	// Set state to fully populated data
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

	id := previousState.ID.ValueString()
	if id == "" {
		resp.Diagnostics.AddError("Unable to read ProfilesWarehouse", "ID is empty")

		return
	}

	response, body, err := r.client.ProfilesWarehousesApi.GetProfilesWarehouse(r.authContext, previousState.ID.ValueString()).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read ProfilesWarehouse",
			getError(err, body),
		)

		return
	}

	var state models.ProfilesWarehouseState

	profileswarehouse := response.Data.GetProfilesWarehouse()
	err = state.Fill(profileswarehouse)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read ProfilesWarehouse",
			err.Error(),
		)

		return
	}

	// This is to satisfy terraform requirements that the returned fields must match the input ones because new settings can be generated in the response
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
	var plan models.ProfilesWarehousePlan
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
	modelMap := api.NewModelMap(settings)

	// The default behavior of updating settings is to upsert. However, to eliminate settings that are no longer necessary, nil is assigned to fields that are no longer found in the resource.
	existingProfilesWarehouse, body, err := r.client.ProfilesWarehousesApi.GetProfilesWarehouse(r.authContext, state.ID.ValueString()).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to update ProfilesWarehouse",
			getError(err, body),
		)

		return
	}
	existingSettings := existingProfilesWarehouse.Data.GetProfilesWarehouse().Settings.Get().Get()

	for key := range existingSettings {
		if settings[key] == nil {
			settings[key] = nil
		}
	}

	out, body, err := r.client.ProfilesWarehousesApi.UpdateProfilesWarehouse(r.authContext, state.ID.ValueString()).UpdateProfilesWarehouseV1Input(api.UpdateProfilesWarehouseV1Input{
		Enabled:  plan.Enabled.ValueBoolPointer(),
		Settings: *api.NewNullableModelMap(modelMap),
		Name:     *api.NewNullableString(plan.Name.ValueStringPointer()),
	}).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to update ProfilesWarehouse",
			getError(err, body),
		)

		return
	}

	profileswarehouse := out.Data.GetProfilesWarehouse()

	err = state.Fill(api.ProfilesWarehouse(profileswarehouse))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to update ProfilesWarehouse",
			err.Error(),
		)

		return
	}

	// This is to satisfy terraform requirements that the returned fields must match the input ones because new settings can be generated in the response
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

	_, body, err := r.client.ProfilesWarehousesApi.DeleteProfilesWarehouse(r.authContext, config.ID.ValueString()).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to delete ProfilesWarehouse",
			getError(err, body),
		)

		return
	}
}

func (r *profilesWarehouseResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
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
