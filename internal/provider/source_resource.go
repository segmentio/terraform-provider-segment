package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

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
		Description: "Configures a Source. For more information, visit the [Segment docs](https://segment.com/docs/connections/sources/).\n\n" +
			docs.GenerateImportDocs("<id>", "segment_source"),
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
									PlanModifiers: []planmodifier.String{
										stringplanmodifier.UseStateForUnknown(),
									},
								},
								"type": schema.StringAttribute{
									Computed:    true,
									Description: "Defines the type for this option in the schema. Types are most commonly strings, but may also represent other primitive types, such as booleans, and numbers, as well as complex types, such as objects and arrays.",
									PlanModifiers: []planmodifier.String{
										stringplanmodifier.UseStateForUnknown(),
									},
								},
								"required": schema.BoolAttribute{
									Computed:    true,
									Description: "Whether this is a required option when setting up the Integration.",
									PlanModifiers: []planmodifier.Bool{
										boolplanmodifier.UseStateForUnknown(),
									},
								},
								"description": schema.StringAttribute{
									Computed:    true,
									Description: "An optional short text description of the field.",
									PlanModifiers: []planmodifier.String{
										stringplanmodifier.UseStateForUnknown(),
									},
								},
								"default_value": schema.StringAttribute{
									CustomType:  jsontypes.NormalizedType{},
									Computed:    true,
									Description: "An optional default value for the field.",
									PlanModifiers: []planmodifier.String{
										stringplanmodifier.UseStateForUnknown(),
									},
								},
								"label": schema.StringAttribute{
									Computed:    true,
									Description: "An optional label for this field.",
									PlanModifiers: []planmodifier.String{
										stringplanmodifier.UseStateForUnknown(),
									},
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
				Description: "The settings associated with the Source. Only settings included in the configuration will be managed by Terraform.",
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
			"labels": schema.SetNestedAttribute{
				Optional:    true,
				Description: "A list of labels applied to the Source.",
				Validators: []validator.Set{
					setvalidator.SizeAtMost(MaxPageSize),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"key": schema.StringAttribute{
							Required:    true,
							Description: "The key that represents the name of this label.",
						},
						"value": schema.StringAttribute{
							Required:    true,
							Description: "The value associated with the key of this label.",
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

	disconnectAllWarehouses := true

	out, body, err := r.client.SourcesAPI.CreateSource(r.authContext).CreateSourceV1Input(api.CreateSourceV1Input{
		Slug:                    plan.Slug.ValueString(),
		Enabled:                 plan.Enabled.ValueBool(),
		MetadataId:              metadataID,
		Settings:                settings,
		DisconnectAllWarehouses: &disconnectAllWarehouses,
	}).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create Source",
			getError(err, body),
		)

		return
	}

	source := out.Data.Source
	resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(source.Id))

	if !plan.Name.IsNull() && !plan.Name.IsUnknown() && plan.Name.ValueString() != "" {
		// This is a workaround for the fact that "name" is allowed to be provided during update but not create
		updateOut, body, err := r.client.SourcesAPI.UpdateSource(r.authContext, out.Data.Source.Id).UpdateSourceV1Input(api.UpdateSourceV1Input{
			Name: plan.Name.ValueStringPointer(),
		}).Execute()
		if body != nil {
			defer body.Body.Close()
		}
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Unable to update Source after creation (ID: %s)", plan.ID.ValueString()),
				getError(err, body),
			)

			return
		}

		source.Name = updateOut.Data.Source.Name
	}

	labels, diags := models.LabelsPlanToAPILabels(ctx, plan.Labels)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if len(labels) > 0 {
		_, body, err := r.client.SourcesAPI.ReplaceLabelsInSource(r.authContext, source.Id).ReplaceLabelsInSourceV1Input(api.ReplaceLabelsInSourceV1Input{
			Labels: models.APILabelsToLabelsV1(labels),
		}).Execute()
		if body != nil {
			defer body.Body.Close()
		}
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to replace Source labels",
				getError(err, body),
			)

			return
		}

		source.Labels = models.APILabelsToLabelsV1(labels)
	}

	var state models.SourceState
	err = state.Fill(source)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to populate Source state",
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
	var previousState models.SourceState
	diags := req.State.Get(ctx, &previousState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := previousState.ID.ValueString()
	if id == "" {
		resp.Diagnostics.AddError("Unable to read Source", "ID is empty")

		return
	}

	out, body, err := r.client.SourcesAPI.GetSource(r.authContext, id).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		if body.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)

			return
		}

		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to read Source (ID: %s)", previousState.ID.ValueString()),
			getError(err, body),
		)

		return
	}

	source := out.Data.Source

	var state models.SourceState
	err = state.Fill(source)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to populate Source state",
			err.Error(),
		)

		return
	}

	// Merge settings: keep config-defined settings while ignoring backend-generated ones not in config
	if !previousState.Settings.IsNull() && !previousState.Settings.IsUnknown() {
		mergedSettings, err := mergeSettings(previousState.Settings, state.Settings, false)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to merge Source settings",
				err.Error(),
			)
			return
		}
		state.Settings = mergedSettings
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

	var name *string
	if !plan.Name.IsNull() && !plan.Name.IsUnknown() && plan.Name.ValueString() != "" {
		name = plan.Name.ValueStringPointer()
	}

	// The default behavior of updating settings is to upsert. However, to eliminate settings that are no longer necessary, nil is assigned to fields that are no longer found in the resource.
	existingSource, body, err := r.client.SourcesAPI.GetSource(r.authContext, state.ID.ValueString()).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to read Source before update (ID: %s)", plan.ID.ValueString()),
			getError(err, body),
		)

		return
	}
	existingSettings := existingSource.Data.GetSource().Settings

	for key := range existingSettings {
		if settings[key] == nil {
			settings[key] = nil
		}
	}

	out, body, err := r.client.SourcesAPI.UpdateSource(r.authContext, state.ID.ValueString()).UpdateSourceV1Input(api.UpdateSourceV1Input{
		Slug:     plan.Slug.ValueStringPointer(),
		Enabled:  plan.Enabled.ValueBoolPointer(),
		Name:     name,
		Settings: settings,
	}).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to update Source (ID: %s)", plan.ID.ValueString()),
			getError(err, body),
		)

		return
	}

	source := out.Data.Source

	labels, diags := models.LabelsPlanToAPILabels(ctx, plan.Labels)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if len(labels) > 0 {
		_, body, err := r.client.SourcesAPI.ReplaceLabelsInSource(r.authContext, source.Id).ReplaceLabelsInSourceV1Input(api.ReplaceLabelsInSourceV1Input{
			Labels: models.APILabelsToLabelsV1(labels),
		}).Execute()
		if body != nil {
			defer body.Body.Close()
		}
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to replace Source Labels",
				getError(err, body),
			)

			return
		}

		source.Labels = models.APILabelsToLabelsV1(labels)
	}

	err = state.Fill(source)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to populate Source state",
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

	_, body, err := r.client.SourcesAPI.DeleteSource(r.authContext, config.ID.ValueString()).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to delete Source (ID: %s)", config.ID.ValueString()),
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
