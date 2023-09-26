package provider

import (
	"context"
	"fmt"

	"github.com/segmentio/terraform-provider-segment/internal/provider/models"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/segmentio/public-api-sdk-go/api"
)

var (
	_ datasource.DataSource              = &roleDataSource{}
	_ datasource.DataSourceWithConfigure = &roleDataSource{}
)

type roleDataSource struct {
	client      *api.APIClient
	authContext context.Context
}

func NewRoleDataSource() datasource.DataSource {
	return &roleDataSource{}
}

func (d *roleDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *roleDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tracking_plan"
}

func (d *roleDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "The unique identifier of the role.",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "The human-readable name of the role.",
			},
			"description": schema.StringAttribute{
				Computed:    true,
				Description: "The human-readable description of the role.",
			},
		},
	}
}

func (d *roleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config models.RoleState
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

	out, body, err := d.client.IAMRolesApi.ListRoles(d.authContext).Pagination(*api.NewPaginationInput(200)).Execute()
	defer body.Body.Close()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Role",
			getError(err, body),
		)

		return
	}

	roles := out.Data.Roles
	var role *api.RoleV1
	for _, r := range roles {
		if r.Id == id {
			role = &r
			break
		}
	}

	if role == nil {
		resp.Diagnostics.AddError(
			"Unable to read Role",
			"Role not found",
		)

		return
	}

	var state models.RoleState
	err = state.Fill(*role)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Role",
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
