package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/segmentio/public-api-sdk-go/api"
)

type UserGroupState struct {
	ID          types.String      `tfsdk:"id"`
	Name        types.String      `tfsdk:"name"`
	Members     []types.String    `tfsdk:"members"`
	Permissions []PermissionState `tfsdk:"permissions"`
}

type UserGroupPlan struct {
	ID          types.String   `tfsdk:"id"`
	Name        types.String   `tfsdk:"name"`
	Members     []types.String `tfsdk:"members"`
	Permissions types.Set      `tfsdk:"permissions"`
}

func (u *UserGroupState) Fill(userGroup api.UserGroup, members []string) error {
	u.ID = types.StringValue(userGroup.Id)
	u.Name = types.StringValue(userGroup.Name)

	u.Members = []types.String{}
	for _, m := range members {
		u.Members = append(u.Members, types.StringValue(m))
	}

	u.Permissions = []PermissionState{}
	for _, p := range userGroup.Permissions {
		var permission PermissionState
		err := permission.Fill(p)
		if err != nil {
			return err
		}
		u.Permissions = append(u.Permissions, permission)
	}

	return nil
}
