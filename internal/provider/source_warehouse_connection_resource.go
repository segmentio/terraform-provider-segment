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

func (r *sourceWarehouseConnectionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_source_warehouse_connection"
}

func (r *sourceWarehouseConnectionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Configures a connection between a source and a warehouse.\n\n" +
			`## Import
This resource is not intended to be imported. Instead, you can create a new connection between the Source and the Warehouse, and any existing connections will be handled automatically.`,
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

	if plan.WarehouseID.String() == "" || plan.SourceID.String() == "" {
		resp.Diagnostics.AddError("Unable to create connection between Source and Warehouse", "At least one ID is empty")

		return
	}

	_, body, err := r.client.WarehousesAPI.AddConnectionFromSourceToWarehouse(r.authContext, plan.WarehouseID.ValueString(), plan.SourceID.ValueString()).Execute()
	if body != nil {
		defer body.Body.Close()
	}
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

func (r *sourceWarehouseConnectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state sourceWarehouseConnectionState

	diags := req.State.Get(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	paginationNext := "MA=="

	for paginationNext != "" {
		if state.SourceID.String() == "" {
			resp.Diagnostics.AddError("Unable to read Source-Warehouse connection", "At least one ID is empty")

			return
		}
		response, body, err := r.client.SourcesAPI.ListConnectedWarehousesFromSource(r.authContext, state.SourceID.ValueString()).Pagination(api.PaginationInput{
			Cursor: &paginationNext,
			Count:  MaxPageSize,
		}).Execute()
		if body != nil {
			defer body.Body.Close()
		}
		if err != nil {
			if body.StatusCode == 404 {
				resp.State.RemoveResource(ctx)

				return
			}

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

func (r *sourceWarehouseConnectionResource) Update(context.Context, resource.UpdateRequest, *resource.UpdateResponse) {
	// All fields force replacement
}

func (r *sourceWarehouseConnectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var config sourceWarehouseConnectionState
	diags := req.State.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.WarehouseID.String() == "" || config.SourceID.String() == "" {
		resp.Diagnostics.AddError("Unable to remove Source-Warehouse connection", "At least one ID is empty")

		return
	}

	_, body, err := r.client.WarehousesAPI.RemoveSourceConnectionFromWarehouse(r.authContext, config.WarehouseID.ValueString(), config.SourceID.ValueString()).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to remove Source-Warehouse connection",
			getError(err, body),
		)

		return
	}
}

func (r *sourceWarehouseConnectionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
