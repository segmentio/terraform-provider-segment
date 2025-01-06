package provider

import (
	"context"
	"fmt"
	"regexp"

	"github.com/segmentio/terraform-provider-segment/internal/provider/docs"
	"github.com/segmentio/terraform-provider-segment/internal/provider/models"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/segmentio/public-api-sdk-go/api"
)

var (
	_ resource.Resource                = &userGroupResource{}
	_ resource.ResourceWithConfigure   = &userGroupResource{}
	_ resource.ResourceWithImportState = &userGroupResource{}
)

func NewUserGroupResource() resource.Resource {
	return &userGroupResource{}
}

type userGroupResource struct {
	client      *api.APIClient
	authContext context.Context
}

func (r *userGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_group"
}

func (r *userGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Configures a User Group. For more information, visit the [Segment docs](https://segment.com/docs/segment-app/iam/concepts/#user-groups).\n\n" +
			docs.GenerateImportDocs("<id>", "segment_user_group"),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The id of the user group.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "A set of users with a set of shared permissions.",
			},
			"members": schema.SetAttribute{
				Required:    true,
				Description: "A list of emails that are members of this user group.",
				ElementType: types.StringType,
				Validators: []validator.Set{
					setvalidator.SizeAtMost(MaxPageSize),
				},
			},
			"permissions": schema.SetNestedAttribute{
				Description: "The permissions associated with this user. This field is currently limited to 200 items and must not be empty.",
				Required:    true,
				Validators: []validator.Set{
					setvalidator.SizeAtMost(MaxPageSize),
					setvalidator.SizeAtLeast(1),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"role_id": schema.StringAttribute{
							Description: "The id of the role associated with this permission.",
							Required:    true,
						},
						"resources": schema.SetNestedAttribute{
							Description: "The resources associated with this permission. This field is currently limited to 200 items.",
							Optional:    true,
							Validators: []validator.Set{
								setvalidator.SizeAtMost(MaxPageSize),
							},
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description: "The id of this resource.",
										Required:    true,
									},
									"type": schema.StringAttribute{
										Description: "The type for this resource.",
										Required:    true,
										Validators: []validator.String{
											stringvalidator.RegexMatches(regexp.MustCompile("^[A-Z_]+$"), "'type' must be in all uppercase"),
										},
									},
									"labels": schema.SetNestedAttribute{
										Description: "The labels that further refine access to this resource. Labels are exclusive to Workspace-level permissions. This field is currently limited to 200 items.",
										Required:    true,
										Validators: []validator.Set{
											setvalidator.SizeAtMost(MaxPageSize),
										},
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"key": schema.StringAttribute{
													Description: "The key that represents the name of this label.",
													Required:    true,
												},
												"value": schema.StringAttribute{
													Description: "The value associated with the key of this label.",
													Required:    true,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (r *userGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan models.UserGroupPlan
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	out, body, err := r.client.IAMGroupsAPI.CreateUserGroup(r.authContext).CreateUserGroupV1Input(api.CreateUserGroupV1Input{
		Name: plan.Name.ValueString(),
	}).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create User Group",
			getError(err, body),
		)

		return
	}

	resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(out.Data.GetUserGroup().Id))

	userGroup := out.Data.GetUserGroup()

	permissions, diags := models.GetPermissionsAPIValueFromPlan(ctx, plan.Permissions)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, body, err = r.client.IAMGroupsAPI.ReplacePermissionsForUserGroup(r.authContext, userGroup.Id).ReplacePermissionsForUserGroupV1Input(api.ReplacePermissionsForUserGroupV1Input{
		Permissions: models.PermissionsToPermissionsInput(permissions),
	}).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to add User Group permissions",
			getError(err, body),
		)

		return
	}

	members := []string{}
	for _, member := range plan.Members {
		members = append(members, member.ValueString())
	}
	if len(members) > 0 {
		_, body, err = r.client.IAMGroupsAPI.ReplaceUsersInUserGroup(r.authContext, userGroup.Id).ReplaceUsersInUserGroupV1Input(api.ReplaceUsersInUserGroupV1Input{
			Emails: members,
		}).Execute()
		if body != nil {
			defer body.Body.Close()
		}
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to add users/invites to User Group",
				getError(err, body),
			)

			return
		}
	}

	getOut, body, err := r.client.IAMGroupsAPI.GetUserGroup(r.authContext, userGroup.Id).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to read User Group (ID: %s)", userGroup.Id),
			getError(err, body),
		)

		return
	}

	var state models.UserGroupState
	err = state.Fill(getOut.Data.UserGroup, members)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to populate User Group state",
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

func (r *userGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var config models.UserGroupState
	diags := req.State.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	out, body, err := r.client.IAMGroupsAPI.GetUserGroup(r.authContext, config.ID.ValueString()).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to read User Group (ID: %s)", config.ID.ValueString()),
			getError(err, body),
		)

		return
	}
	userGroup := out.Data.GetUserGroup()

	usersOut, body, err := r.client.IAMGroupsAPI.ListUsersFromUserGroup(r.authContext, config.ID.ValueString()).Pagination(api.PaginationInput{Count: MaxPageSize}).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to read User Group members (ID: %s)", config.ID.ValueString()),
			getError(err, body),
		)

		return
	}

	invitesOut, body, err := r.client.IAMGroupsAPI.ListInvitesFromUserGroup(r.authContext, config.ID.ValueString()).Pagination(api.PaginationInput{Count: MaxPageSize}).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to read User Group members (ID: %s)", config.ID.ValueString()),
			getError(err, body),
		)

		return
	}

	members := []string{}
	for _, user := range usersOut.Data.Users {
		members = append(members, user.Email)
	}
	members = append(members, invitesOut.Data.Emails...)

	var state models.UserGroupState
	err = state.Fill(userGroup, members)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to populate User Group state",
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

func (r *userGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan models.UserGroupPlan
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var config models.UserGroupState
	diags = req.State.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, body, err := r.client.IAMGroupsAPI.UpdateUserGroup(r.authContext, config.ID.ValueString()).UpdateUserGroupV1Input(api.UpdateUserGroupV1Input{
		Name: plan.Name.ValueString(),
	}).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to update User Group (ID: %s)", plan.ID.ValueString()),
			getError(err, body),
		)

		return
	}

	permissions, diags := models.GetPermissionsAPIValueFromPlan(ctx, plan.Permissions)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, body, err = r.client.IAMGroupsAPI.ReplacePermissionsForUserGroup(r.authContext, config.ID.ValueString()).ReplacePermissionsForUserGroupV1Input(api.ReplacePermissionsForUserGroupV1Input{
		Permissions: models.PermissionsToPermissionsInput(permissions),
	}).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to add User Group permissions",
			getError(err, body),
		)

		return
	}

	members := []string{}
	for _, member := range plan.Members {
		members = append(members, member.ValueString())
	}
	_, body, err = r.client.IAMGroupsAPI.ReplaceUsersInUserGroup(r.authContext, config.ID.ValueString()).ReplaceUsersInUserGroupV1Input(api.ReplaceUsersInUserGroupV1Input{
		Emails: members,
	}).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to add users/invites to User Group",
			getError(err, body),
		)

		return
	}

	getOut, body, err := r.client.IAMGroupsAPI.GetUserGroup(r.authContext, config.ID.ValueString()).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to read User Group (ID: %s)", config.ID.ValueString()),
			getError(err, body),
		)

		return
	}

	var state models.UserGroupState
	err = state.Fill(getOut.Data.UserGroup, members)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to populate User Group state",
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

func (r *userGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var config models.UserGroupState
	diags := req.State.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, body, err := r.client.IAMGroupsAPI.DeleteUserGroup(r.authContext, config.ID.ValueString()).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to delete User Group (ID: %s)", config.ID.ValueString()),
			getError(err, body),
		)

		return
	}
}

func (r *userGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *userGroupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
