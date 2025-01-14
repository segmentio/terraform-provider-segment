package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/segmentio/public-api-sdk-go/api"
)

type ReverseETLModelState struct {
	ID                    types.String `tfsdk:"id"`
	SourceID              types.String `tfsdk:"source_id"`
	Name                  types.String `tfsdk:"name"`
	Description           types.String `tfsdk:"description"`
	Enabled               types.Bool   `tfsdk:"enabled"`
	Query                 types.String `tfsdk:"query"`
	QueryIdentifierColumn types.String `tfsdk:"query_identifier_column"`
}

func (r *ReverseETLModelState) Fill(model api.ReverseEtlModel) error {
	r.ID = types.StringValue(model.Id)
	r.SourceID = types.StringValue(model.SourceId)
	r.Name = types.StringValue(model.Name)
	r.Description = types.StringValue(model.Description)
	r.Enabled = types.BoolValue(model.Enabled)
	r.Query = types.StringValue(model.Query)
	r.QueryIdentifierColumn = types.StringValue(model.QueryIdentifierColumn)

	return nil
}
