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

	"github.com/segmentio/public-api-sdk-go/api"
)

var (
	_ resource.Resource                = &reverseETLModelResource{}
	_ resource.ResourceWithConfigure   = &reverseETLModelResource{}
	_ resource.ResourceWithImportState = &reverseETLModelResource{}
)

func NewReverseETLModelResource() resource.Resource {
	return &reverseETLModelResource{}
}

type reverseETLModelResource struct {
	client      *api.APIClient
	authContext context.Context
}

func (r *reverseETLModelResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_reverse_etl_model"
}

func (r *reverseETLModelResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Configures a Reverse ETL Model. For more information, visit the [Segment docs](https://segment.com/docs/connections/reverse-etl/).\n\n" +
			docs.GenerateImportDocs("<id>", "segment_reverse_etl_model"),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The unique identifier for the model.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"source_id": schema.StringAttribute{
				Required:    true,
				Description: "Indicates which Source to attach this model to.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "A short, human-readable description of the Model.",
			},
			"description": schema.StringAttribute{
				Required:    true,
				Description: "A longer, more descriptive explanation of the Model.",
			},
			"enabled": schema.BoolAttribute{
				Required:    true,
				Description: "Indicates whether the Model should have syncs enabled. When disabled, no syncs will be triggered, regardless of the enabled status of the attached destinations/subscriptions.",
			},
			// "schedule_strategy": schema.StringAttribute{
			// 	Optional:           true,
			// 	DeprecationMessage: "Remove this attribute's configuration as it no longer is used and the attribute will be removed in the next major version of the provider. Please use `reverse_etl_schedule` in the destination_subscription resource instead.",
			// 	Description:        "Determines the strategy used for triggering syncs, which will be used in conjunction with scheduleConfig.",
			// },
			"query": schema.StringAttribute{
				Required:    true,
				Description: "The SQL query that will be executed to extract data from the connected Source.",
			},
			"query_identifier_column": schema.StringAttribute{
				Required:    true,
				Description: "Indicates the column named in `query` that should be used to uniquely identify the extracted records.",
			},
			// "schedule_config": schema.StringAttribute{
			// 	Optional:           true,
			// 	DeprecationMessage: "Remove this attribute's configuration as it no longer is used and the attribute will be removed in the next major version of the provider. Please use `reverse_etl_schedule` in the destination_subscription resource instead.",
			// 	Description:        "Depending on the chosen strategy, configures the schedule for this model.",
			// 	CustomType:         jsontypes.NormalizedType{},
			// },
		},
	}
}

func (r *reverseETLModelResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan models.ReverseETLModelState
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	out, body, err := r.client.ReverseETLAPI.CreateReverseEtlModel(r.authContext).CreateReverseEtlModelInput(api.CreateReverseEtlModelInput{
		Name:                  plan.Name.ValueString(),
		SourceId:              plan.SourceID.ValueString(),
		Description:           plan.Description.ValueString(),
		Enabled:               plan.Enabled.ValueBool(),
		Query:                 plan.Query.ValueString(),
		QueryIdentifierColumn: plan.QueryIdentifierColumn.ValueString(),
	}).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create Reverse ETL model",
			getError(err, body),
		)

		return
	}

	reverseETLModel := out.Data.ReverseEtlModel

	resp.State.SetAttribute(ctx, path.Root("id"), reverseETLModel.Id)

	var state models.ReverseETLModelState
	err = state.Fill(reverseETLModel)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to populate Reverse ETL model state",
			err.Error(),
		)

		return
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *reverseETLModelResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var previousState models.ReverseETLModelState

	diags := req.State.Get(ctx, &previousState)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	out, body, err := r.client.ReverseETLAPI.GetReverseEtlModel(r.authContext, previousState.ID.ValueString()).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to read Reverse ETL model (ID: %s)", previousState.ID.ValueString()),
			getError(err, body),
		)

		return
	}

	var state models.ReverseETLModelState

	err = state.Fill(out.Data.ReverseEtlModel)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to populate Reverse ETL model state",
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

func (r *reverseETLModelResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan models.ReverseETLModelState
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state models.ReverseETLModelState
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	out, body, err := r.client.ReverseETLAPI.UpdateReverseEtlModel(r.authContext, state.ID.ValueString()).UpdateReverseEtlModelInput(api.UpdateReverseEtlModelInput{
		Name:                  plan.Name.ValueStringPointer(),
		Description:           plan.Description.ValueStringPointer(),
		Enabled:               plan.Enabled.ValueBoolPointer(),
		Query:                 plan.Query.ValueStringPointer(),
		QueryIdentifierColumn: plan.QueryIdentifierColumn.ValueStringPointer(),
	}).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to update Reverse ETL model (ID: %s)", plan.ID.ValueString()),
			getError(err, body),
		)

		return
	}

	err = state.Fill(out.Data.ReverseEtlModel)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to populate Reverse ETL model state",
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

func (r *reverseETLModelResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var config models.ReverseETLModelState
	diags := req.State.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, body, err := r.client.ReverseETLAPI.DeleteReverseEtlModel(r.authContext, config.ID.ValueString()).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to delete Reverse ETL model (ID: %s)", config.ID.ValueString()),
			getError(err, body),
		)

		return
	}
}

func (r *reverseETLModelResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *reverseETLModelResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
