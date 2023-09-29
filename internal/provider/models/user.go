package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/segmentio/public-api-sdk-go/api"
)

type UserState struct {
	ID          types.String       `tfsdk:"id"`
	Name        types.String       `tfsdk:"name"`
	Email       types.String       `tfsdk:"email"`
	IsInvite    types.Bool         `tfsdk:"is_invite"`
	Permissions []PermissionsState `tfsdk:"permissions"`
}

type PermissionsState struct {
	RoleName  types.String     `tfsdk:"role_name"`
	RoleID    types.String     `tfsdk:"role_id"`
	Resources []ResourcesState `tfsdk:"resources"`
	Labels    []LabelState     `tfsdk:"labels"`
}

type ResourcesState struct {
	ID     types.String `tfsdk:"id"`
	Type   types.String `tfsdk:"type"`
	Labels []LabelState `tfsdk:"labels"`
}

type UserPlan struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Email       types.String `tfsdk:"email"`
	Permissions types.Set    `tfsdk:"permissions"`
}

type PermissionsPlan struct {
	RoleName  types.String `tfsdk:"role_name"`
	RoleID    types.String `tfsdk:"role_id"`
	Resources types.Set    `tfsdk:"resources"`
	Labels    types.Set    `tfsdk:"labels"`
}

type ResourcesPlan struct {
	ID     types.String `tfsdk:"id"`
	Type   types.String `tfsdk:"type"`
	Labels types.Set    `tfsdk:"labels"`
}

func (u *UserState) Fill(user api.UserV1) error {
	u.ID = types.StringValue(user.Id)
	u.Name = types.StringValue(user.Name)
	u.Email = types.StringValue(user.Email)

	u.Permissions = []PermissionsState{}
	for _, p := range user.Permissions {
		var permission PermissionsState
		permission.Fill(p)
		u.Permissions = append(u.Permissions, permission)
	}

	return nil
}

func (p *PermissionsState) ToAPIValue() api.InvitePermissionV1 {
	labels := []api.AllowedLabelBeta{}
	for _, label := range p.Labels {
		labels = append(labels, label.ToAPIValue())
	}

	resources := []api.ResourceV1{}
	for _, resource := range p.Resources {
		resources = append(resources, resource.ToAPIValue())
	}

	return api.InvitePermissionV1{
		RoleId:    p.RoleID.ValueString(),
		Resources: resources,
		Labels:    labels,
	}
}

func (p *PermissionsState) Fill(permission api.PermissionV1) error {
	p.RoleID = types.StringValue(permission.RoleId)
	p.Resources = []ResourcesState{}
	for _, r := range permission.Resources {
		var resource ResourcesState
		resource.Fill(r)
		p.Resources = append(p.Resources, resource)
	}

	p.Labels = []LabelState{}
	for _, l := range permission.Labels {
		label := LabelState{}
		label.Fill(api.LabelV1(l))

		p.Labels = append(p.Labels, label)
	}

	return nil
}

func GetPermissionsAPIValue(permissions []PermissionsState) []api.InvitePermissionV1 {
	var apiPermissions []api.InvitePermissionV1

	for _, permission := range permissions {
		apiPermissions = append(apiPermissions, permission.ToAPIValue())
	}

	return apiPermissions
}

func GetPermissionsInputAPIValue(permissions []PermissionsState) []api.PermissionInputV1 {
	apiPermissions := GetPermissionsAPIValue(permissions)

	inputPermissions := []api.PermissionInputV1{}
	for _, permission := range apiPermissions {
		resources := []api.PermissionResourceV1{}
		for _, resource := range permission.Resources {
			resources = append(resources, api.PermissionResourceV1{
				Id:   resource.Id,
				Type: resource.Type,
			})
		}

		inputPermissions = append(inputPermissions, api.PermissionInputV1{
			RoleId:    permission.RoleId,
			Resources: resources,
		})
	}

	return inputPermissions
}

func (r *ResourcesState) ToAPIValue() api.ResourceV1 {
	return api.ResourceV1{
		Id:   r.ID.ValueString(),
		Type: r.Type.ValueString(),
	}
}

func (r *ResourcesState) Fill(resource api.PermissionResourceV1) {
	r.ID = types.StringValue(resource.Id)
	r.Type = types.StringValue(resource.Type)
	r.Labels = []LabelState{}
	for _, l := range resource.Labels {
		label := LabelState{}
		label.Fill(api.LabelV1(l))

		r.Labels = append(r.Labels, label)
	}
}
