package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/segmentio/public-api-sdk-go/api"
)

type RoleState struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

func (r *RoleState) Fill(role api.RoleV1) error {
	r.ID = types.StringValue(role.Id)
	r.Name = types.StringValue(role.Name)
	r.Description = types.StringValue(role.Description)

	return nil
}
