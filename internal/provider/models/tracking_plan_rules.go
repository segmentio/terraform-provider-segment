package models

import (
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/segmentio/public-api-sdk-go/api"
)

type TrackingPlanRulesState struct {
	TrackingPlanID types.String `tfsdk:"tracking_plan_id"`
	Rules          []RulesState `tfsdk:"rules"`
}

type TrackingPlanRulesPlan struct {
	TrackingPlanID types.String `tfsdk:"tracking_plan_id"`
	Rules          types.Set    `tfsdk:"rules"`
}

type RulesState struct {
	Type       types.String         `tfsdk:"type"`
	Key        types.String         `tfsdk:"key"`
	JsonSchema jsontypes.Normalized `tfsdk:"json_schema"`
	Version    types.Float64        `tfsdk:"version"`
}

type RulesDSState struct {
	Type         types.String         `tfsdk:"type"`
	Key          types.String         `tfsdk:"key"`
	JsonSchema   jsontypes.Normalized `tfsdk:"json_schema"`
	Version      types.Float64        `tfsdk:"version"`
	CreatedAt    types.String         `tfsdk:"created_at"`
	UpdatedAt    types.String         `tfsdk:"updated_at"`
	DeprecatedAt types.String         `tfsdk:"deprecated_at"`
}

type TrackingPlanRulesDSState struct {
	TrackingPlanID types.String   `tfsdk:"tracking_plan_id"`
	Rules          []RulesDSState `tfsdk:"rules"`
}

func (t *TrackingPlanRulesState) Fill(rules []api.RuleV1, trackingPlanID string) error {
	t.TrackingPlanID = types.StringValue(trackingPlanID)

	t.Rules = []RulesState{}
	for _, rule := range rules {
		r := RulesState{}

		r.Type = types.StringValue(rule.Type)
		if rule.Key != nil {
			r.Key = types.StringValue(*rule.Key)
		}

		jsonSchema, err := json.Marshal(rule.JsonSchema)
		if err != nil {
			return err
		}
		r.JsonSchema = jsontypes.NewNormalizedValue(string(jsonSchema))

		r.Version = types.Float64Value(float64(rule.Version))

		t.Rules = append(t.Rules, r)
	}

	return nil
}

func (t *TrackingPlanRulesDSState) Fill(rules []api.RuleV1, trackingPlanID string) error {
	t.TrackingPlanID = types.StringValue(trackingPlanID)

	t.Rules = []RulesDSState{}
	for _, rule := range rules {
		r := RulesDSState{}

		r.Type = types.StringValue(rule.Type)
		if rule.Key != nil {
			r.Key = types.StringValue(*rule.Key)
		}

		jsonSchema, err := json.Marshal(rule.JsonSchema)
		if err != nil {
			return err
		}
		r.JsonSchema = jsontypes.NewNormalizedValue(string(jsonSchema))

		r.Version = types.Float64Value(float64(rule.Version))

		if rule.CreatedAt != nil {
			r.CreatedAt = types.StringValue(*rule.CreatedAt)
		}

		if rule.UpdatedAt != nil {
			r.UpdatedAt = types.StringValue(*rule.UpdatedAt)
		}

		if rule.DeprecatedAt != nil {
			r.DeprecatedAt = types.StringValue(*rule.DeprecatedAt)
		}

		t.Rules = append(t.Rules, r)
	}

	return nil
}

func (r *RulesState) ToApiRule() (api.RuleV1, diag.Diagnostics) {
	var jsonSchema interface{}
	diags := r.JsonSchema.Unmarshal(&jsonSchema)
	if diags.HasError() {
		return api.RuleV1{}, diags
	}

	return api.RuleV1{
		Type:       r.Type.ValueString(),
		Key:        r.Key.ValueStringPointer(),
		Version:    float32(r.Version.ValueFloat64()),
		JsonSchema: jsonSchema,
	}, diags
}
