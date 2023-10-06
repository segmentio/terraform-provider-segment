package provider

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/segmentio/public-api-sdk-go/api"
	"github.com/segmentio/terraform-provider-segment/internal/provider/models"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var (
	_ resource.Resource                = &userResource{}
	_ resource.ResourceWithConfigure   = &userResource{}
	_ resource.ResourceWithImportState = &userResource{}
)

func NewUserResource() resource.Resource {
	return &userResource{}
}

type userResource struct {
	client      *api.APIClient
	authContext context.Context
}

func (r *userResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("is_invite"), false)...)
}

func (r *userResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (r *userResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A user or invite belonging to a Segment Workspace. Only users may be imported.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The unique identifier for this user, or the user's email if the invite has not been accepted.",
			},
			"name": schema.StringAttribute{
				Description: "The human-readable name of this user, or the user's email if the invite has not been accepted.",
				Computed:    true,
			},
			"email": schema.StringAttribute{
				Description: "The email address associated with this user.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"is_invite": schema.BoolAttribute{
				Description: "Whether or not this user is an invite.",
				Computed:    true,
			},
			"permissions": schema.SetNestedAttribute{
				Description: "The permissions associated with this user. This field is currently limited to 200 items.",
				Required:    true,
				Validators: []validator.Set{
					setvalidator.SizeAtMost(MaxPageSize),
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
											stringvalidator.RegexMatches(regexp.MustCompile("^[A-Z]+$"), "'type' must be in all uppercase"),
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
												"description": schema.StringAttribute{
													Description: "An optional description of the purpose of this label.",
													Computed:    true,
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

func (r *userResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan models.UserPlan
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiPermissions, diags := models.GetPermissionsAPIValueFromPlan(ctx, plan.Permissions)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	_, body, err := r.client.IAMUsersApi.CreateInvites(r.authContext).CreateInvitesV1Input(api.CreateInvitesV1Input{
		Invites: []api.InviteV1{
			{
				Email:       plan.Email.ValueString(),
				Permissions: models.PermissionsToInvitePermissions(apiPermissions),
			},
		},
	}).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to invite user",
			getError(err, body),
		)

		return
	}

	user, err := findUser(r.authContext, r.client, plan.Email.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to find user",
			err.Error(),
		)

		return
	}

	if user == nil { // Handle invite
		inviteUser := &api.UserV1{
			Email:       plan.Email.ValueString(),
			Id:          plan.Email.ValueString(),
			Name:        plan.Email.ValueString(),
			Permissions: apiPermissions,
		}
		state := models.UserState{}
		err := state.Fill(*inviteUser)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to populate invite state",
				err.Error(),
			)

			return
		}
		state.IsInvite = types.BoolValue(true)
		diags = resp.State.Set(ctx, state)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	} else { // Handle user
		state := models.UserState{}
		err := state.Fill(*user)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to populate user state",
				err.Error(),
			)

			return
		}
		state.IsInvite = types.BoolValue(false)
		diags = resp.State.Set(ctx, state)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
}

func (r *userResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.UserState
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var user api.UserV1

	if state.IsInvite.ValueBool() { // Handle potential invite
		foundUser, err := findUser(r.authContext, r.client, state.Email.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to find user",
				err.Error(),
			)

			return
		}

		if foundUser == nil { // Handle invite
			state.IsInvite = types.BoolValue(true)
			diags := resp.State.Set(ctx, state)
			resp.Diagnostics.Append(diags...)

			return
		}

		user = *foundUser
	} else { // Handle user
		out, body, err := r.client.IAMUsersApi.GetUser(r.authContext, state.ID.ValueString()).Execute()
		if body != nil {
			defer body.Body.Close()
		}
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to read user",
				getError(err, body),
			)

			return
		}

		user = api.UserV1(out.Data.User)
	}

	state = models.UserState{}
	err := state.Fill(user)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to populate user state",
			err.Error(),
		)

		return
	}
	state.IsInvite = types.BoolValue(false)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *userResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state models.UserState
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var plan models.UserPlan
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	permissions, diags := models.GetPermissionsAPIValueFromPlan(ctx, plan.Permissions)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var userID string

	if state.IsInvite.ValueBool() { // Handle potential invite
		foundUser, err := findUser(r.authContext, r.client, state.Email.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to find user",
				err.Error(),
			)

			return
		}

		if foundUser == nil { // Handle invite
			_, body, err := r.client.IAMUsersApi.DeleteInvites(r.authContext).Emails([]string{state.Email.ValueString()}).Execute()
			if body != nil {
				defer body.Body.Close()
			}
			if err != nil {
				resp.Diagnostics.AddError(
					"Unable to delete user",
					getError(err, body),
				)

				return
			}

			_, body, err = r.client.IAMUsersApi.CreateInvites(r.authContext).CreateInvitesV1Input(api.CreateInvitesV1Input{
				Invites: []api.InviteV1{
					{
						Email:       plan.Email.ValueString(),
						Permissions: models.PermissionsToInvitePermissions(permissions),
					},
				},
			}).Execute()
			if body != nil {
				defer body.Body.Close()
			}
			if err != nil {
				resp.Diagnostics.AddError(
					"Unable to invite user",
					getError(err, body),
				)

				return
			}

			inviteUser := &api.UserV1{
				Email:       plan.Email.ValueString(),
				Id:          plan.Email.ValueString(),
				Name:        plan.Email.ValueString(),
				Permissions: permissions,
			}
			state := models.UserState{}
			err = state.Fill(*inviteUser)
			if err != nil {
				resp.Diagnostics.AddError(
					"Unable to populate invite state",
					err.Error(),
				)

				return
			}
			state.IsInvite = types.BoolValue(true)
			diags = resp.State.Set(ctx, state)
			resp.Diagnostics.Append(diags...)

			return
		}

		// Handle user that was previously an invite
		userID = foundUser.Id
	} else { // Handle user
		userID = state.ID.ValueString()
	}

	_, body, err := r.client.IAMUsersApi.ReplacePermissionsForUser(r.authContext, userID).ReplacePermissionsForUserV1Input(api.ReplacePermissionsForUserV1Input{
		Permissions: models.PermissionsToPermissionsInput(permissions),
	}).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to update user",
			getError(err, body),
		)

		return
	}

	out, body, err := r.client.IAMUsersApi.GetUser(r.authContext, userID).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read user",
			getError(err, body),
		)

		return
	}

	state = models.UserState{}
	err = state.Fill(api.UserV1(out.Data.User))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to populate user state",
			err.Error(),
		)

		return
	}
	state.IsInvite = types.BoolValue(false)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *userResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state models.UserState
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var userID string

	if state.IsInvite.ValueBool() { // Handle potential invite
		foundUser, err := findUser(r.authContext, r.client, state.Email.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to find user",
				err.Error(),
			)

			return
		}

		if foundUser == nil { // Handle invite
			_, body, err := r.client.IAMUsersApi.DeleteInvites(r.authContext).Emails([]string{state.Email.ValueString()}).Execute()
			if body != nil {
				defer body.Body.Close()
			}
			if err != nil {
				resp.Diagnostics.AddError(
					"Unable to delete user",
					getError(err, body),
				)

				return
			}
		} else {
			userID = foundUser.Id
		}
	} else { // Handle user
		userID = state.ID.ValueString()
	}

	_, body, err := r.client.IAMUsersApi.DeleteUsers(r.authContext).UserIds([]string{userID}).Execute()
	if body != nil {
		defer body.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to delete user",
			getError(err, body),
		)

		return
	}
}

func (r *userResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = config.client
	r.authContext = config.authContext
}

func findUser(authContext context.Context, client *api.APIClient, email string) (*api.UserV1, error) {
	var pageToken api.NullableString
	firstPageToken := "MA=="
	pageToken.Set(&firstPageToken)

	for pageToken.IsSet() {
		out, body, err := client.IAMUsersApi.ListUsers(authContext).Pagination(api.PaginationInput{
			Count:  MaxPageSize,
			Cursor: pageToken.Get(),
		}).Execute()
		if body != nil {
			defer body.Body.Close()
		}

		if err != nil {
			return nil, errors.New(getError(err, body))
		}

		users := out.Data.Users

		for _, user := range users {
			if user.Email == email {
				return &user, nil
			}
		}

		pageToken = out.Data.Pagination.Next
	}

	return nil, nil
}
