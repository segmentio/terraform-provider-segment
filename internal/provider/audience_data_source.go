package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/segmentio/public-api-sdk-go/api"
	"github.com/segmentio/terraform-provider-segment/internal/provider/models"
)

var _ datasource.DataSource = &audienceDataSource{}
var _ datasource.DataSourceWithConfigure = &audienceDataSource{}

func NewAudienceDataSource() datasource.DataSource {
	return &audienceDataSource{}
}

type audienceDataSource struct {
	client      *api.APIClient
	authContext context.Context
}

func (d *audienceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_audience"
}

func (d *audienceDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"space_id": schema.StringAttribute{
				Required:    true,
				Description: "The Space ID.",
			},
			"id": schema.StringAttribute{
				Required:    true,
				Description: "The Audience ID.",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "The name of the Audience.",
			},
			"description": schema.StringAttribute{
				Computed:    true,
				Description: "The description of the Audience.",
			},
			"key": schema.StringAttribute{
				Computed:    true,
				Description: "The key of the Audience.",
			},
			"enabled": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether the Audience is enabled.",
			},
			"definition": schema.MapAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "The definition of the Audience.",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "The status of the Audience.",
			},
			"options": schema.MapAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "Additional options for the Audience.",
			},
		},
	}
}

func (d *audienceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.AudienceState
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	out, _, err := d.client.GetAudience(d.authContext, data.SpaceID.ValueString(), data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to read Audience", err.Error())
		return
	}
	audience := out.Data.Audience
	data.Name = types.StringValue(audience.Name)
	data.Description = types.StringValue(audience.Description)
	data.Key = types.StringValue(audience.Key)
	data.Enabled = types.BoolValue(audience.Enabled)
	data.Status = types.StringValue(audience.Status)
	// TODO: convert audience.Definition and audience.Options to types.Map
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (d *audienceDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	config, ok := req.ProviderData.(*ClientInfo)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected DataSource Configure Type",
			fmt.Sprintf("Expected ClientInfo, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	d.client = config.client
	d.authContext = config.authContext
}
