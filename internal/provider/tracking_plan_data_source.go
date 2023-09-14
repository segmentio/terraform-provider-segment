package provider

import (
	"context"
	"fmt"

	"terraform-provider-segment/internal/provider/models"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/segmentio/public-api-sdk-go/api"
)

var (
	_ datasource.DataSource              = &trackingPlanDataSource{}
	_ datasource.DataSourceWithConfigure = &trackingPlanDataSource{}
)

type trackingPlanDataSource struct {
	client      *api.APIClient
	authContext context.Context
}

func NewTrackingPlanDataSource() datasource.DataSource {
	return &trackingPlanDataSource{}
}

func (d *trackingPlanDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	config, ok := req.ProviderData.(*ClientInfo)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected ClientInfo, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = config.client
	d.authContext = config.authContext
}

func (d *trackingPlanDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tracking_plan"
}

func (d *trackingPlanDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "The Tracking Plan's identifier.",
			},
			"slug": schema.StringAttribute{
				Computed:    true,
				Description: "URL-friendly slug of this Tracking Plan.",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "The Tracking Plan's name.",
			},
			"description": schema.StringAttribute{
				Computed:    true,
				Description: "The Tracking Plan's description.",
			},
			"type": schema.StringAttribute{
				Computed:    true,
				Description: "The Tracking Plan's type.",
			},
			"updated_at": schema.StringAttribute{
				Computed:    true,
				Description: "The timestamp of the last change to the Tracking Plan.",
			},
			"created_at": schema.StringAttribute{
				Computed:    true,
				Description: "The timestamp of this Tracking Plan's creation.",
			},
			"rules": schema.SetNestedAttribute{
				Computed:    true,
				Description: `The list of Tracking Plan rules. Currently limited to 200 rules.`,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Computed: true,
							Description: `The type for this Tracking Plan rule.

							Enum: "COMMON" "GROUP" "IDENTIFY" "PAGE" "SCREEN" "TRACK"`,
						},
						"key": schema.StringAttribute{
							Computed:    true,
							Description: "Key to this rule (free-form string like 'Button clicked').",
						},
						"json_schema": schema.StringAttribute{
							Computed:    true,
							Description: "JSON Schema of this rule.",
						},
						"version": schema.Float64Attribute{
							Computed:    true,
							Description: "Version of this rule.",
						},
						"created_at": schema.StringAttribute{
							Computed:    true,
							Description: "The timestamp of this rule's creation.",
						},
						"updated_at": schema.StringAttribute{
							Computed:    true,
							Description: "The timestamp of this rule's last change.",
						},
						"deprecated_at": schema.StringAttribute{
							Computed:    true,
							Description: "The timestamp of this rule's deprecation.",
						},
					},
				},
			},
		},
	}
}

func (d *trackingPlanDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config models.TrackingPlanDSState
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := config.ID.ValueString()
	if id == "" {
		resp.Diagnostics.AddError("Unable to read Tracking Plan", "ID is empty")
		return
	}

	out, body, err := d.client.TrackingPlansApi.GetTrackingPlan(d.authContext, id).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Tracking Plan",
			getError(err, body),
		)
		return
	}

	trackingPlan := out.Data.GetTrackingPlan()

	rulesOut, body, err := d.client.TrackingPlansApi.ListRulesFromTrackingPlan(d.authContext, id).Pagination(*api.NewPaginationInput(200)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to get Tracking Plan rules",
			getError(err, body),
		)
		return
	}

	var state models.TrackingPlanDSState
	err = state.Fill(trackingPlan, &rulesOut.Data.Rules)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Tracking Plan",
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
