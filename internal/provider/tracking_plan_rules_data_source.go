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

type trackingPlanRulesDataSource struct {
	client      *api.APIClient
	authContext context.Context
}

func NewTrackingPlanRulesDataSource() datasource.DataSource {
	return &trackingPlanRulesDataSource{}
}

func (d *trackingPlanRulesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *trackingPlanRulesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tracking_plan_rules"
}

func (d *trackingPlanRulesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"tracking_plan_id": schema.StringAttribute{
				Required:    true,
				Description: "The Tracking Plan's identifier.",
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

func (d *trackingPlanRulesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config models.TrackingPlanRulesDSState
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	out, body, err := d.client.TrackingPlansApi.ListRulesFromTrackingPlan(d.authContext, config.TrackingPlanID.ValueString()).Pagination(*api.NewPaginationInput(200)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to get Tracking Plan rules",
			getError(err, body),
		)
		return
	}

	var state models.TrackingPlanRulesDSState
	err = state.Fill(out.Data.GetRules(), config.TrackingPlanID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to get Tracking Plan rules",
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
