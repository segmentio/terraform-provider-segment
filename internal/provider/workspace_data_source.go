package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/segmentio/public-api-sdk-go/api"
)

var (
	_ datasource.DataSource              = &workspaceDataSource{}
	_ datasource.DataSourceWithConfigure = &workspaceDataSource{}
)

func NewWorkspaceDataSource() datasource.DataSource {
	return &workspaceDataSource{}
}

type workspaceDataSource struct {
	client      *api.APIClient
	authContext context.Context
}

type workspaceDataSourceModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	Slug types.String `tfsdk:"slug"`
}

func (d *workspaceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workspace"
}

func (d *workspaceDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Reads the Workspace.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The unique identifier.",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "The human-readable name.",
			},
			"slug": schema.StringAttribute{
				Computed:    true,
				Description: "The URL-friendly slug.",
			},
		},
	}
}

func (d *workspaceDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state workspaceDataSourceModel

	workspace, body, err := d.client.WorkspacesAPI.GetWorkspace(d.authContext).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Workspace",
			getError(err, body),
		)

		return
	}

	state.ID = types.StringValue(workspace.Data.Workspace.Id)
	state.Name = types.StringValue(workspace.Data.Workspace.Name)
	state.Slug = types.StringValue(workspace.Data.Workspace.Slug)

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *workspaceDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
