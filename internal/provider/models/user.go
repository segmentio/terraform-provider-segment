package models

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/segmentio/public-api-sdk-go/api"
)

type UserState struct {
	ID          types.String      `tfsdk:"id"`
	Name        types.String      `tfsdk:"name"`
	Email       types.String      `tfsdk:"email"`
	IsInvite    types.Bool        `tfsdk:"is_invite"`
	Permissions []PermissionState `tfsdk:"permissions"`
}

func (u *UserState) Fill(user api.UserV1) error {
	u.ID = types.StringValue(user.Id)
	u.Name = types.StringValue(user.Name)
	u.Email = types.StringValue(user.Email)

	u.Permissions = []PermissionState{}
	for _, p := range user.Permissions {
		var permission PermissionState
		permission.Fill(p)
		u.Permissions = append(u.Permissions, permission)
	}

	return nil
}

type PermissionState struct {
	RoleID    types.String    `tfsdk:"role_id"`
	Resources []ResourceState `tfsdk:"resources"`
}

func (p *PermissionState) ToAPIValue() api.PermissionV1 {
	labels := []api.AllowedLabelBeta{}

	resources := []api.PermissionResourceV1{}
	for _, resource := range p.Resources {
		resources = append(resources, resource.ToAPIValue())
	}

	return api.PermissionV1{
		RoleId:    p.RoleID.ValueString(),
		Resources: resources,
		Labels:    labels,
	}
}

type ResourceState struct {
	ID     types.String `tfsdk:"id"`
	Type   types.String `tfsdk:"type"`
	Labels []LabelState `tfsdk:"labels"`
}

func (p *PermissionState) Fill(permission api.PermissionV1) error {
	p.RoleID = types.StringValue(permission.RoleId)
	p.Resources = []ResourceState{}
	for _, r := range permission.Resources {
		var resource ResourceState
		resource.Fill(r)
		p.Resources = append(p.Resources, resource)
	}

	return nil
}

func (r *ResourceState) ToAPIValue() api.PermissionResourceV1 {
	apiLabels := []api.AllowedLabelBeta{}

	for _, label := range r.Labels {
		apiLabels = append(apiLabels, label.ToAPIValue())
	}

	if len(apiLabels) == 0 {
		apiLabels = nil
	}

	return api.PermissionResourceV1{
		Id:     r.ID.ValueString(),
		Type:   r.Type.ValueString(),
		Labels: apiLabels,
	}
}

func (r *ResourceState) Fill(resource api.PermissionResourceV1) {
	r.ID = types.StringValue(resource.Id)
	r.Type = types.StringValue(resource.Type)
	r.Labels = []LabelState{}
	for _, l := range resource.Labels {
		label := LabelState{}
		label.Fill(api.LabelV1(l))

		r.Labels = append(r.Labels, label)
	}
}

type UserPlan struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Email       types.String `tfsdk:"email"`
	IsInvite    types.Bool   `tfsdk:"is_invite"`
	Permissions types.Set    `tfsdk:"permissions"`
}

type PermissionPlan struct {
	RoleID    types.String `tfsdk:"role_id"`
	Resources types.Set    `tfsdk:"resources"`
}

func (p *PermissionPlan) ToAPIValue(ctx context.Context) (api.PermissionV1, diag.Diagnostics) {
	var outDiags diag.Diagnostics

	resources := []ResourcePlan{}
	diags := p.Resources.ElementsAs(ctx, &resources, false)
	outDiags.Append(diags...)
	if outDiags.HasError() {
		return api.PermissionV1{}, outDiags
	}
	apiResources := []api.PermissionResourceV1{}
	for _, resource := range resources {
		apiResource, diags := resource.ToAPIValue(ctx)
		outDiags.Append(diags...)
		if outDiags.HasError() {
			return api.PermissionV1{}, outDiags
		}
		apiResources = append(apiResources, apiResource)
	}

	return api.PermissionV1{
		RoleId:    p.RoleID.ValueString(),
		Resources: apiResources,
	}, outDiags
}

type ResourcePlan struct {
	ID     types.String `tfsdk:"id"`
	Type   types.String `tfsdk:"type"`
	Labels types.Set    `tfsdk:"labels"`
}

func (r *ResourcePlan) ToAPIValue(ctx context.Context) (api.PermissionResourceV1, diag.Diagnostics) {
	apiLabels, diags := LabelsPlanToApiLabels(ctx, r.Labels)
	if diags.HasError() {
		return api.PermissionResourceV1{}, diags
	}

	return api.PermissionResourceV1{
		Id:     r.ID.ValueString(),
		Type:   r.Type.ValueString(),
		Labels: apiLabels,
	}, diags
}

func GetPermissionsAPIValueFromPlan(ctx context.Context, permissions types.Set) ([]api.PermissionV1, diag.Diagnostics) {
	var outDiags diag.Diagnostics

	var permissionsPlan = []PermissionPlan{}
	diags := permissions.ElementsAs(ctx, &permissionsPlan, false)
	outDiags.Append(diags...)
	if outDiags.HasError() {
		return []api.PermissionV1{}, outDiags
	}

	var apiPermissions []api.PermissionV1

	for _, permission := range permissionsPlan {
		p, diags := permission.ToAPIValue(ctx)
		outDiags.Append(diags...)
		if outDiags.HasError() {
			return []api.PermissionV1{}, outDiags
		}
		apiPermissions = append(apiPermissions, p)
	}

	return apiPermissions, outDiags
}

func GetPermissionsAPIValueFromState(permissions []PermissionState) []api.PermissionV1 {
	var apiPermissions []api.PermissionV1

	for _, permission := range permissions {
		apiPermissions = append(apiPermissions, permission.ToAPIValue())
	}

	return apiPermissions
}

func PermissionsToInvitePermissions(permissions []api.PermissionV1) []api.InvitePermissionV1 {
	var apiPermissions []api.InvitePermissionV1

	for _, permission := range permissions {
		resources := []api.ResourceV1{}
		for _, resource := range permission.Resources {
			resources = append(resources, api.ResourceV1{
				Id:   resource.Id,
				Type: resource.Type,
			})
		}
		apiPermissions = append(apiPermissions, api.InvitePermissionV1{
			RoleId:    permission.RoleId,
			Resources: resources,
			Labels:    permission.Labels,
		})
	}

	return apiPermissions
}

func PermissionsToPermissionsInput(permissions []api.PermissionV1) []api.PermissionInputV1 {
	var apiPermissions []api.PermissionInputV1

	for _, permission := range permissions {
		resources := []api.PermissionResourceV1{}
		for _, resource := range permission.Resources {
			resources = append(resources, api.PermissionResourceV1{
				Id:     resource.Id,
				Type:   resource.Type,
				Labels: resource.Labels,
			})
		}
		apiPermissions = append(apiPermissions, api.PermissionInputV1{
			RoleId:    permission.RoleId,
			Resources: resources,
		})
	}

	return apiPermissions
}
