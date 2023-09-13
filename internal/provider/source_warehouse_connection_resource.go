package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/segmentio/public-api-sdk-go/api"
)

var (
	_ resource.Resource              = &sourceWarehouseConnectionResource{}
	_ resource.ResourceWithConfigure = &sourceWarehouseConnectionResource{}
)

func NewSourceWarehouseConnectionResource() resource.Resource {
	return &sourceWarehouseConnectionResource{}
}

type sourceWarehouseConnectionResource struct {
	client      *api.APIClient
	authContext context.Context
}

type sourceWarehouseConnectionState struct {
	SourceID    types.String `tfsdk:"source_id"`
	WarehouseID types.String `tfsdk:"warehouse_id"`
}

func (d *sourceWarehouseConnectionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_source_warehouse_connection"
}

func (d *sourceWarehouseConnectionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Represents a connection between a source and a warehouse",
		Attributes: map[string]schema.Attribute{
			"source_id": schema.StringAttribute{
				Required:    true,
				Description: "The id of the Source.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"warehouse_id": schema.StringAttribute{
				Required:    true,
				Description: "The id of the Warehouse.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *sourceWarehouseConnectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan sourceWarehouseConnectionState
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, body, err := r.client.WarehousesApi.AddConnectionFromSourceToWarehouse(r.authContext, plan.WarehouseID.ValueString(), plan.SourceID.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create connection between Source and Warehouse",
			getError(err, body),
		)
		return
	}

	state := sourceWarehouseConnectionState{
		SourceID:    plan.SourceID,
		WarehouseID: plan.WarehouseID,
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *sourceWarehouseConnectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state sourceWarehouseConnectionState

	diags := req.State.Get(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paginationNext := "MA=="

	for paginationNext != "" {
		response, body, err := d.client.SourcesApi.ListConnectedWarehousesFromSource(d.authContext, state.SourceID.ValueString()).Pagination(api.PaginationInput{
			Cursor: &paginationNext,
			Count:  200,
		}).Execute()
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to read Source-Warehouse connection",
				getError(err, body),
			)
			return
		}

		for _, warehouse := range response.Data.GetWarehouses() {
			if warehouse.Id == state.WarehouseID.ValueString() {
				diags = resp.State.Set(ctx, &state)
				resp.Diagnostics.Append(diags...)
				if resp.Diagnostics.HasError() {
					return
				}

				return
			}
		}

		if response.Data.Pagination.Next.IsSet() {
			paginationNext = *response.Data.Pagination.Next.Get()
		} else {
			paginationNext = ""
		}
	}

	diags = resp.State.Set(ctx, &sourceWarehouseConnectionState{
		SourceID:    types.StringValue("not_found"),
		WarehouseID: types.StringValue("not_found"),
	})
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *sourceWarehouseConnectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// All fields force replacement
}

func (r *sourceWarehouseConnectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var config sourceWarehouseConnectionState
	diags := req.State.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, body, err := r.client.WarehousesApi.RemoveSourceConnectionFromWarehouse(r.authContext, config.WarehouseID.ValueString(), config.SourceID.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to remove Source-Warehouse connection",
			getError(err, body),
		)
		return
	}
}

func (d *sourceWarehouseConnectionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	d.client = config.client
	d.authContext = config.authContext
}
