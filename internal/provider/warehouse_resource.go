package provider

import (
	"context"
	"fmt"

	"github.com/segmentio/terraform-provider-segment/internal/provider/docs"
	"github.com/segmentio/terraform-provider-segment/internal/provider/models"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/segmentio/public-api-sdk-go/api"
)

var (
	_ resource.Resource                = &warehouseResource{}
	_ resource.ResourceWithConfigure   = &warehouseResource{}
	_ resource.ResourceWithImportState = &warehouseResource{}
)

func NewWarehouseResource() resource.Resource {
	return &warehouseResource{}
}

type warehouseResource struct {
	client      *api.APIClient
	authContext context.Context
}

func (r *warehouseResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_warehouse"
}

func (r *warehouseResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Configures a Warehouse. For more information, visit the [Segment docs](https://segment.com/docs/connections/storage/).\n\n" +
			docs.GenerateImportDocs("<id>", "segment_warehouse"),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The id of the Warehouse.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"metadata": schema.SingleNestedAttribute{
				Description: "The metadata for the Warehouse.",
				Required:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Required:    true,
						Description: "The id of this object.",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"name": schema.StringAttribute{
						Computed:    true,
						Description: "The name of this object.",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"slug": schema.StringAttribute{
						Computed:    true,
						Description: "A human-readable, unique identifier for object.",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"description": schema.StringAttribute{
						Computed:    true,
						Description: "A description, in English, of this object.",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"logos": schema.SingleNestedAttribute{
						Computed:    true,
						Description: "Logo information for this object.",
						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.UseStateForUnknown(),
						},
						Attributes: map[string]schema.Attribute{
							"default": schema.StringAttribute{
								Computed:    true,
								Description: "The default URL for this logo.",
							},
							"mark": schema.StringAttribute{
								Computed:    true,
								Description: "The logo mark.",
							},
							"alt": schema.StringAttribute{
								Computed:    true,
								Description: "The alternative text for this logo.",
							},
						},
					},
					"options": schema.ListNestedAttribute{
						Computed:    true,
						Description: "The Integration options for this object.",
						PlanModifiers: []planmodifier.List{
							listplanmodifier.UseStateForUnknown(),
						},
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Computed:    true,
									Description: "The name identifying this option in the context of a Segment Integration.",
								},
								"type": schema.StringAttribute{
									Computed:    true,
									Description: "Defines the type for this option in the schema.",
								},
								"required": schema.BoolAttribute{
									Computed:    true,
									Description: "Whether this is a required option when setting up the Integration.",
								},
								"description": schema.StringAttribute{
									Computed:    true,
									Description: "An optional short text description of the field.",
								},
								"default_value": schema.StringAttribute{
									CustomType:  jsontypes.NormalizedType{},
									Computed:    true,
									Description: "An optional default value for the field.",
								},
								"label": schema.StringAttribute{
									Computed:    true,
									Description: "An optional label for this field.",
								},
							},
						},
					},
				},
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "An optional human-readable name for this Warehouse.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"workspace_id": schema.StringAttribute{
				Computed:    true,
				Description: "The id of the Workspace that owns this Warehouse.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "When set to true, this Warehouse receives data.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"settings": schema.StringAttribute{
				Required:    true,
				Description: "The settings associated with this Warehouse. Common settings are connection-related configuration used to connect to it, for example host, username, and port. Only settings included in the configuration will be managed by Terraform.",
				CustomType:  jsontypes.NormalizedType{},
			},
		},
	}
}

func (r *warehouseResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan models.WarehousePlan
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

	name := plan.Name.ValueStringPointer()
	if *name == "" {
		name = nil
	}

	disconnectAllSources := true

	out, body, err := r.client.WarehousesAPI.CreateWarehouse(r.authContext).CreateWarehouseV1Input(api.CreateWarehouseV1Input{
		Enabled:              plan.Enabled.ValueBoolPointer(),
		MetadataId:           metadataID,
		Settings:             settings,
		Name:                 name,
		DisconnectAllSources: &disconnectAllSources,
	}).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create Warehouse",
			getError(err, body),
		)

		return
	}

	warehouse := out.Data.GetWarehouse()

	resp.State.SetAttribute(ctx, path.Root("id"), warehouse.Id)

	var state models.WarehouseState
	err = state.Fill(warehouse)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to populate Warehouse state",
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

func (r *warehouseResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var previousState models.WarehouseState

	diags := req.State.Get(ctx, &previousState)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := previousState.ID.ValueString()
	if id == "" {
		resp.Diagnostics.AddError("Unable to read Warehouse", "ID is empty")

		return
	}

	response, body, err := r.client.WarehousesAPI.GetWarehouse(r.authContext, previousState.ID.ValueString()).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		if body.StatusCode == 404 {
			resp.State.RemoveResource(ctx)

			return
		}

		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to read Warehouse (ID: %s)", previousState.ID.ValueString()),
			getError(err, body),
		)

		return
	}

	var state models.WarehouseState

	warehouse := response.Data.GetWarehouse()
	err = state.Fill(warehouse)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to populate Warehouse state",
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

func (r *warehouseResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan models.WarehousePlan
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state models.WarehouseState
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

	// The default behavior of updating settings is to upsert. However, to eliminate settings that are no longer necessary, nil is assigned to fields that are no longer found in the resource.
	existingWarehouse, body, err := r.client.WarehousesAPI.GetWarehouse(r.authContext, state.ID.ValueString()).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to read Warehouse before update (ID: %s)", state.ID.ValueString()),
			getError(err, body),
		)

		return
	}
	existingSettings := existingWarehouse.Data.GetWarehouse().Settings

	for key := range existingSettings {
		if settings[key] == nil {
			settings[key] = nil
		}
	}

	out, body, err := r.client.WarehousesAPI.UpdateWarehouse(r.authContext, state.ID.ValueString()).UpdateWarehouseV1Input(api.UpdateWarehouseV1Input{
		Enabled:  plan.Enabled.ValueBoolPointer(),
		Settings: settings,
		Name:     *api.NewNullableString(plan.Name.ValueStringPointer()),
	}).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to update Warehouse (ID: %s)", plan.ID.ValueString()),
			getError(err, body),
		)

		return
	}

	warehouse := out.Data.GetWarehouse()

	err = state.Fill(warehouse)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to populate Warehouse state",
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

func (r *warehouseResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var config models.WarehouseState
	diags := req.State.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, body, err := r.client.WarehousesAPI.DeleteWarehouse(r.authContext, config.ID.ValueString()).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to delete Warehouse (ID: %s)", config.ID.ValueString()),
			getError(err, body),
		)

		return
	}
}

func (r *warehouseResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *warehouseResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
