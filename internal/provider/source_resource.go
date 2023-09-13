package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"

	"terraform-provider-segment/internal/provider/models"

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
	_ resource.Resource                = &sourceResource{}
	_ resource.ResourceWithConfigure   = &sourceResource{}
	_ resource.ResourceWithImportState = &sourceResource{}
)

func NewSourceResource() resource.Resource {
	return &sourceResource{}
}

type sourceResource struct {
	client      *api.APIClient
	authContext context.Context
}

func (r *sourceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_source"
}

func (r *sourceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The id of the Source.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"slug": schema.StringAttribute{
				Required:    true,
				Description: "The slug used to identify the Source in the Segment app.",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The name of the Source.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"metadata": schema.SingleNestedAttribute{
				Description: "The metadata for the Source.",
				Required:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Required: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
						Description: "The id for this Source metadata in the Segment catalog.",
					},
					"name": schema.StringAttribute{
						Computed:    true,
						Description: "The user-friendly name of this Source.",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"slug": schema.StringAttribute{
						Computed:    true,
						Description: "The slug that identifies this Source in the Segment app.",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"description": schema.StringAttribute{
						Computed:    true,
						Description: "The description of this Source.",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"logos": schema.SingleNestedAttribute{
						Description: "The logos for this Source.",
						Computed:    true,
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
						Description: "Options for this Source.",
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
									Description: "Defines the type for this option in the schema. Types are most commonly strings, but may also represent other primitive types, such as booleans, and numbers, as well as complex types, such as objects and arrays.",
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
					"categories": schema.ListAttribute{
						Computed:    true,
						ElementType: types.StringType,
						PlanModifiers: []planmodifier.List{
							listplanmodifier.UseStateForUnknown(),
						},
						Description: "A list of categories this Source belongs to.",
					},
					"is_cloud_event_source": schema.BoolAttribute{
						Computed:    true,
						Description: "True if this is a Cloud Event Source.",
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
			"settings": schema.StringAttribute{
				Required:    true,
				Description: "The settings associated with the Source.",
				CustomType:  jsontypes.NormalizedType{},
			},
			"workspace_id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "The id of the Workspace that owns the Source.",
			},
			"enabled": schema.BoolAttribute{
				Required:    true,
				Description: "Enable to receive data from the Source.",
			},
			"write_keys": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
				Description: "The write keys used to send data from the Source. This field is left empty when the current token does not have the 'source admin' permission.",
			},
			"labels": schema.ListNestedAttribute{
				Computed:    true,
				Description: "A list of labels applied to the Source.",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"key": schema.StringAttribute{
							Computed:    true,
							Description: "The key that represents the name of this label.",
						},
						"value": schema.StringAttribute{
							Computed:    true,
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

func (r *sourceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan models.SourcePlan
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	wrappedMetadataId, err := plan.Metadata.Attributes()["id"].ToTerraformValue(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to decode metadata id",
			err.Error(),
		)
		return
	}

	var metadataId string
	err = wrappedMetadataId.As(&metadataId)
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

	out, body, err := r.client.SourcesApi.CreateSource(r.authContext).CreateSourceV1Input(api.CreateSourceV1Input{
		Slug:       plan.Slug.ValueString(),
		Enabled:    plan.Enabled.ValueBool(),
		MetadataId: metadataId,
		Settings:   *api.NewNullableModelMap(modelMap),
	}).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create Source",
			getError(err, body),
		)
		return
	}

	source := out.Data.Source

	if !plan.Name.IsNull() && !plan.Name.IsUnknown() && plan.Name.ValueString() != "" {
		// This is a workaround for the fact that "name" is allowed to be provided during update but not create
		updateOut, body, err := r.client.SourcesApi.UpdateSource(r.authContext, out.Data.Source.Id).UpdateSourceV1Input(api.UpdateSourceV1Input{
			Name: plan.Name.ValueStringPointer(),
		}).Execute()
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to create Source",
				getError(err, body),
			)
			return
		}

		source.Name = updateOut.Data.Source.Name
	}

	var state models.SourceState
	err = state.Fill(api.Source4(source))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create Source",
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

func (r *sourceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var config models.SourceState
	diags := req.State.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	out, body, err := r.client.SourcesApi.GetSource(r.authContext, config.ID.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Source",
			getError(err, body),
		)
		return
	}

	source := out.Data.Source

	var state models.SourceState
	err = state.Fill(source)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Source",
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

func (r *sourceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan models.SourcePlan
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state models.SourceState
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

	var name *string
	if !plan.Name.IsNull() && !plan.Name.IsUnknown() && plan.Name.ValueString() != "" {
		name = plan.Name.ValueStringPointer()
	}

	// The default behavior of updating settings is to upsert. However, to eliminate settings that are no longer necessary, nil is assigned to fields that are no longer found in the resource.
	existingSource, body, err := r.client.SourcesApi.GetSource(r.authContext, state.ID.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to update Source",
			getError(err, body),
		)
		return
	}
	existingSettings := existingSource.Data.GetSource().Settings.Get().Get()

	for key := range existingSettings {
		if settings[key] == nil {
			settings[key] = nil
		}
	}

	out, body, err := r.client.SourcesApi.UpdateSource(r.authContext, state.ID.ValueString()).UpdateSourceV1Input(api.UpdateSourceV1Input{
		Slug:     plan.Slug.ValueStringPointer(),
		Enabled:  plan.Enabled.ValueBoolPointer(),
		Name:     name,
		Settings: *api.NewNullableModelMap(modelMap),
	}).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to update Source",
			getError(err, body),
		)
		return
	}

	source := out.Data.Source

	err = state.Fill(api.Source4(source))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to update Source",
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

func (r *sourceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var config models.SourceState
	diags := req.State.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, body, err := r.client.SourcesApi.DeleteSource(r.authContext, config.ID.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to delete Source",
			getError(err, body),
		)
		return
	}
}

func (r *sourceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *sourceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
