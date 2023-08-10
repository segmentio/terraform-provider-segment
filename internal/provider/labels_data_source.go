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
	_ datasource.DataSource              = &labelsDataSource{}
	_ datasource.DataSourceWithConfigure = &labelsDataSource{}
)

func NewLabelsDataSource() datasource.DataSource {
	return &labelsDataSource{}
}

type labelsDataSource struct {
	client      *api.APIClient
	authContext context.Context
}

type labelsDataSourceModel struct {
	Id     types.String           `tfsdk:"id"`
	Labels []labelDataSourceModel `tfsdk:"labels"`
}

type labelDataSourceModel struct {
	Key         types.String `tfsdk:"key"`
	Value       types.String `tfsdk:"value"`
	Description types.String `tfsdk:"description"`
}

func (d *labelsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_labels"
}

func (d *labelsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "All labels associated with the current Workspace.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"labels": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"key": schema.StringAttribute{
							Description: "The key that represents the name of this label.",
							Computed:    true,
						},
						"value": schema.StringAttribute{
							Description: "The value associated with the key of this label.",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "An optional description of the purpose of this label.",
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

func (d *labelsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state labelsDataSourceModel

	state.Id = types.StringValue("placeholder")

	response, _, err := d.client.LabelsApi.ListLabels(d.authContext).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Labels",
			err.Error(),
		)
		return
	}

	labels := response.Data.Labels

	for _, label := range labels {
		newLabel := labelDataSourceModel{
			Key:   types.StringValue(label.Key),
			Value: types.StringValue(label.Value),
		}

		if label.Description != nil {
			newLabel.Description = types.StringValue(*label.Description)
		}

		state.Labels = append(state.Labels, newLabel)
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *labelsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
