package provider

import (
	"context"
	"fmt"

	"github.com/segmentio/terraform-provider-segment/internal/provider/docs"
	"github.com/segmentio/terraform-provider-segment/internal/provider/models"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/segmentio/public-api-sdk-go/api"
)

var (
	_ resource.Resource                = &transformationResource{}
	_ resource.ResourceWithConfigure   = &transformationResource{}
	_ resource.ResourceWithImportState = &transformationResource{}
)

func NewTransformationResource() resource.Resource {
	return &transformationResource{}
}

type transformationResource struct {
	client      *api.APIClient
	authContext context.Context
}

func (r *transformationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_transformation"
}

func (r *transformationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Configures a Transformation. For more information, visit the [Segment docs](https://segment.com/docs/protocols/transform/).\n\n" +
			docs.GenerateImportDocs("<id>", "segment_transformation"),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The id of the Transformation.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"source_id": schema.StringAttribute{
				Required:    true,
				Description: "The Source associated with the Transformation.",
			},
			"destination_metadata_id": schema.StringAttribute{
				Optional:    true,
				Description: "The optional Destination metadata associated with the Transformation.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the Transformation.",
			},
			"enabled": schema.BoolAttribute{
				Required:    true,
				Description: "If the Transformation is enabled.",
			},
			"if": schema.StringAttribute{
				Required: true,
				Description: `If statement (FQL) to match events.

				For standard event matchers, use the following: Track -> "event='EVENT_NAME'" Identify -> "type='identify'" Group -> "type='group'"`,
			},
			"new_event_name": schema.StringAttribute{
				Optional:    true,
				Description: "Optional new event name for renaming events. Works only for 'track' event type.",
			},
			"property_renames": schema.SetNestedAttribute{
				Required:    true,
				Description: "Optional array for renaming properties collected by your events.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"old_name": schema.StringAttribute{
							Required:    true,
							Description: "The old name of the property.",
						},
						"new_name": schema.StringAttribute{
							Required:    true,
							Description: "The new name to rename the property.",
						},
					},
				},
			},
			"property_value_transformations": schema.SetNestedAttribute{
				Required:    true,
				Description: "Optional array for transforming properties and values collected by your events. Limited to 10 properties.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"property_paths": schema.SetAttribute{
							Required:    true,
							Description: "The property paths. The maximum number of paths is 10.",
							ElementType: types.StringType,
						},
						"property_value": schema.StringAttribute{
							Required:    true,
							Description: "The new value of the property paths.",
						},
					},
				},
			},
			"fql_defined_properties": schema.SetNestedAttribute{
				Required:    true,
				Description: "Optional array for defining new properties in FQL. Currently limited to 1 property.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"fql": schema.StringAttribute{
							Required:    true,
							Description: "The FQL expression used to compute the property.",
						},
						"property_name": schema.StringAttribute{
							Required:    true,
							Description: "The new property name.",
						},
					},
				},
			},
		},
	}
}

func (r *transformationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan models.TransformationPlan
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	propertyRenames, diags := models.PropertyRenamesPlanToAPIValue(ctx, plan.PropertyRenames)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	propertyValueTransformations, diags := models.PropertyValueTransformationsPlanToAPIValue(ctx, plan.PropertyValueTransformations)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	fqlDefinedProperties, diags := models.FQLDefinedPropertiesPlanToAPIValue(ctx, plan.FQLDefinedProperties)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	out, body, err := r.client.TransformationsAPI.CreateTransformation(r.authContext).CreateTransformationV1Input(api.CreateTransformationV1Input{
		Name:                         plan.Name.ValueString(),
		SourceId:                     plan.SourceID.ValueString(),
		DestinationMetadataId:        plan.DestinationMetadataID.ValueStringPointer(),
		Enabled:                      plan.Enabled.ValueBool(),
		If:                           plan.If.ValueString(),
		NewEventName:                 plan.NewEventName.ValueStringPointer(),
		PropertyRenames:              propertyRenames,
		PropertyValueTransformations: propertyValueTransformations,
		FqlDefinedProperties:         fqlDefinedProperties,
	}).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create Transformation",
			getError(err, body),
		)

		return
	}

	transformation := out.Data.GetTransformation()

	resp.State.SetAttribute(ctx, path.Root("id"), transformation.Id)

	var state models.TransformationState
	state.Fill(transformation)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *transformationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var previousState models.TransformationState

	diags := req.State.Get(ctx, &previousState)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	out, body, err := r.client.TransformationsAPI.GetTransformation(r.authContext, previousState.ID.ValueString()).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to read Transformation (ID: %s)", previousState.ID.ValueString()),
			getError(err, body),
		)

		return
	}

	var state models.TransformationState

	state.Fill(out.Data.Transformation)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *transformationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan models.TransformationPlan
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state models.TransformationState
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	propertyRenames, diags := models.PropertyRenamesPlanToAPIValue(ctx, plan.PropertyRenames)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	propertyValueTransformations, diags := models.PropertyValueTransformationsPlanToAPIValue(ctx, plan.PropertyValueTransformations)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	fqlDefinedProperties, diags := models.FQLDefinedPropertiesPlanToAPIValue(ctx, plan.FQLDefinedProperties)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	out, body, err := r.client.TransformationsAPI.UpdateTransformation(r.authContext, state.ID.ValueString()).UpdateTransformationV1Input(api.UpdateTransformationV1Input{
		Name:                         plan.Name.ValueStringPointer(),
		Enabled:                      plan.Enabled.ValueBoolPointer(),
		If:                           plan.If.ValueStringPointer(),
		NewEventName:                 plan.NewEventName.ValueStringPointer(),
		SourceId:                     plan.SourceID.ValueStringPointer(),
		DestinationMetadataId:        plan.DestinationMetadataID.ValueStringPointer(),
		PropertyRenames:              propertyRenames,
		PropertyValueTransformations: propertyValueTransformations,
		FqlDefinedProperties:         fqlDefinedProperties,
	}).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to update Transformation (ID: %s)", plan.ID.ValueString()),
			getError(err, body),
		)

		return
	}

	state.Fill(out.Data.Transformation)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *transformationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var config models.TransformationState
	diags := req.State.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, body, err := r.client.TransformationsAPI.DeleteTransformation(r.authContext, config.ID.ValueString()).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to delete Transformation (ID: %s)", config.ID.ValueString()),
			getError(err, body),
		)

		return
	}
}

func (r *transformationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *transformationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
