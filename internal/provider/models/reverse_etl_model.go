package models

import (
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/segmentio/public-api-sdk-go/api"
)

type ReverseETLModelState struct {
	ID                    types.String         `tfsdk:"id"`
	SourceID              types.String         `tfsdk:"source_id"`
	Name                  types.String         `tfsdk:"name"`
	Description           types.String         `tfsdk:"description"`
	Enabled               types.Bool           `tfsdk:"enabled"`
	ScheduleStrategy      types.String         `tfsdk:"schedule_strategy"`
	Query                 types.String         `tfsdk:"query"`
	QueryIdentifierColumn types.String         `tfsdk:"query_identifier_column"`
	ScheduleConfig        jsontypes.Normalized `tfsdk:"schedule_config"`
}

func (r *ReverseETLModelState) Fill(model api.ReverseEtlModel) error {
	r.ID = types.StringValue(model.Id)
	r.SourceID = types.StringValue(model.SourceId)
	r.Name = types.StringValue(model.Name)
	r.Description = types.StringValue(model.Description)
	r.Enabled = types.BoolValue(model.Enabled)
	r.ScheduleStrategy = types.StringValue(model.ScheduleStrategy)
	r.Query = types.StringValue(model.Query)
	r.QueryIdentifierColumn = types.StringValue(model.QueryIdentifierColumn)
	scheduleConfig, err := GetScheduleConfig(model.ScheduleConfig)
	if err != nil {
		return err
	}
	r.ScheduleConfig = scheduleConfig
	if r.ScheduleConfig.IsNull() {
		empty := "{}"
		r.ScheduleConfig = jsontypes.NewNormalizedPointerValue(&empty)
	}

	return nil
}

func GetScheduleConfig(scheduleConfig api.NullableScheduleConfig) (jsontypes.Normalized, error) {
	if !scheduleConfig.IsSet() {
		return jsontypes.NewNormalizedNull(), nil
	}

	jsonScheduleConfigString, err := scheduleConfig.Get().MarshalJSON()
	if err != nil {
		return jsontypes.NewNormalizedNull(), err
	}

	if jsonScheduleConfigString == nil {
		return jsontypes.NewNormalizedValue("{}"), nil
	}

	return jsontypes.NewNormalizedValue(string(jsonScheduleConfigString)), nil
}
