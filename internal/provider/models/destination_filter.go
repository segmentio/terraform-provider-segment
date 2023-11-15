package models

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/segmentio/public-api-sdk-go/api"
)

type DestinationFilterState struct {
	ID            types.String                   `tfsdk:"id"`
	If            types.String                   `tfsdk:"if"`
	DestinationID types.String                   `tfsdk:"destination_id"`
	SourceID      types.String                   `tfsdk:"source_id"`
	Title         types.String                   `tfsdk:"title"`
	Description   types.String                   `tfsdk:"description"`
	Enabled       types.Bool                     `tfsdk:"enabled"`
	Actions       []DestinationFilterActionState `tfsdk:"actions"`
}

type DestinationFilterPlan struct {
	ID            types.String `tfsdk:"id"`
	If            types.String `tfsdk:"if"`
	DestinationID types.String `tfsdk:"destination_id"`
	SourceID      types.String `tfsdk:"source_id"`
	Title         types.String `tfsdk:"title"`
	Description   types.String `tfsdk:"description"`
	Enabled       types.Bool   `tfsdk:"enabled"`
	Actions       types.Set    `tfsdk:"actions"`
}

type DestinationFilterActionState struct {
	Type    types.String         `tfsdk:"type"`
	Percent types.Float64        `tfsdk:"percent"`
	Path    types.String         `tfsdk:"path"`
	Fields  jsontypes.Normalized `tfsdk:"fields"`
}

func ActionsPlanToAPIActions(ctx context.Context, actions types.Set) ([]api.DestinationFilterActionV1, diag.Diagnostics) {
	apiFilters := []api.DestinationFilterActionV1{}

	if !actions.IsNull() && !actions.IsUnknown() {
		stateActions := []DestinationFilterActionState{}
		diags := actions.ElementsAs(ctx, &stateActions, false)
		if diags.HasError() {
			return apiFilters, diags
		}
		for _, action := range stateActions {
			apiAction, diags := action.ToAPIValue()
			if diags.HasError() {
				return apiFilters, diags
			}
			apiFilters = append(apiFilters, apiAction)
		}
	}

	return apiFilters, diag.Diagnostics{}
}

func (d *DestinationFilterState) Fill(filter *api.DestinationFilterV1) error {
	d.ID = types.StringValue(filter.Id)
	d.If = types.StringValue(filter.If)
	d.DestinationID = types.StringValue(filter.DestinationId)
	d.SourceID = types.StringValue(filter.SourceId)
	d.Title = types.StringValue(filter.Title)
	d.Description = types.StringPointerValue(filter.Description)
	d.Enabled = types.BoolValue(filter.Enabled)
	actions, err := d.getActions(filter.Actions)
	if err != nil {
		return err
	}
	d.Actions = actions

	return nil
}

func (d *DestinationFilterState) getActions(actions []api.DestinationFilterActionV1) ([]DestinationFilterActionState, error) {
	var actionsToAdd []DestinationFilterActionState

	for _, action := range actions {
		actionToAdd := DestinationFilterActionState{}
		err := actionToAdd.Fill(action)
		if err != nil {
			return actionsToAdd, err
		}

		actionsToAdd = append(actionsToAdd, actionToAdd)
	}

	return actionsToAdd, nil
}

func (d *DestinationFilterActionState) Fill(action api.DestinationFilterActionV1) error {
	var percent *float64
	if action.Percent != nil {
		// Converts to float64 in a way that ensures equivalence with plan
		p, err := strconv.ParseFloat(fmt.Sprintf("%f", *action.Percent), 64)
		if err != nil {
			return fmt.Errorf("failed to parse action percent: %v", err)
		}
		percent = &p
	}

	var fields *string
	if action.Fields != nil {
		b, err := json.Marshal(action.Fields)
		if err != nil {
			return fmt.Errorf("failed to marshal action fields to JSON: %v", err)
		}
		fieldsString := string(b)
		fields = &fieldsString
	}

	d.Type = types.StringValue(action.Type)
	d.Percent = types.Float64PointerValue(percent)
	d.Path = types.StringPointerValue(action.Path)
	d.Fields = jsontypes.NewNormalizedPointerValue(fields)

	return nil
}

func (d *DestinationFilterActionState) ToAPIValue() (api.DestinationFilterActionV1, diag.Diagnostics) {
	var percent *float32
	if !d.Percent.IsNull() && !d.Percent.IsUnknown() {
		p := float32(d.Percent.ValueFloat64())
		percent = &p
	}

	var actions map[string]interface{}
	if !d.Fields.IsNull() && !d.Fields.IsUnknown() {
		diags := d.Fields.Unmarshal(&actions)
		if diags.HasError() {
			return api.DestinationFilterActionV1{}, diags
		}
	}

	return api.DestinationFilterActionV1{
		Type:    d.Type.ValueString(),
		Percent: percent,
		Path:    d.Path.ValueStringPointer(),
		Fields:  actions,
	}, diag.Diagnostics{}
}
