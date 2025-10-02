package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type AudienceState struct {
	ID          types.String `tfsdk:"id"`
	SpaceID     types.String `tfsdk:"space_id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Key         types.String `tfsdk:"key"`
	Enabled     types.Bool   `tfsdk:"enabled"`
	Definition  types.Map    `tfsdk:"definition"`
	Status      types.String `tfsdk:"status"`
	Options     types.Map    `tfsdk:"options"`
}
