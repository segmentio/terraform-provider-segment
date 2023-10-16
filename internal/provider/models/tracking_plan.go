package models

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/segmentio/public-api-sdk-go/api"
)

type TrackingPlanDSState struct {
	ID          types.String   `tfsdk:"id"`
	Name        types.String   `tfsdk:"name"`
	Slug        types.String   `tfsdk:"slug"`
	Description types.String   `tfsdk:"description"`
	Type        types.String   `tfsdk:"type"`
	UpdatedAt   types.String   `tfsdk:"updated_at"`
	CreatedAt   types.String   `tfsdk:"created_at"`
	Rules       []RulesDSState `tfsdk:"rules"`
}

type TrackingPlanState struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Slug        types.String `tfsdk:"slug"`
	Description types.String `tfsdk:"description"`
	Type        types.String `tfsdk:"type"`
	UpdatedAt   types.String `tfsdk:"updated_at"`
	CreatedAt   types.String `tfsdk:"created_at"`
	Rules       []RulesState `tfsdk:"rules"`
}

type TrackingPlanPlan struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Slug        types.String `tfsdk:"slug"`
	Description types.String `tfsdk:"description"`
	Type        types.String `tfsdk:"type"`
	UpdatedAt   types.String `tfsdk:"updated_at"`
	CreatedAt   types.String `tfsdk:"created_at"`
	Rules       types.Set    `tfsdk:"rules"`
}

func (t *TrackingPlanState) Fill(trackingPlan api.TrackingPlan, rules *[]api.RuleV1) error {
	t.ID = types.StringValue(trackingPlan.Id)
	t.Name = types.StringPointerValue(trackingPlan.Name)
	t.Slug = types.StringPointerValue(trackingPlan.Slug)
	t.Description = types.StringPointerValue(trackingPlan.Description)
	t.Type = types.StringValue(trackingPlan.Type)
	t.UpdatedAt = types.StringPointerValue(trackingPlan.UpdatedAt)
	t.CreatedAt = types.StringPointerValue(trackingPlan.CreatedAt)

	t.Rules = []RulesState{}
	if rules != nil {
		for _, rule := range *rules {
			r := RulesState{}

			r.Type = types.StringValue(rule.Type)
			if rule.Key != nil {
				r.Key = types.StringPointerValue(rule.Key)
			}

			jsonSchema, err := json.Marshal(rule.JsonSchema)
			if err != nil {
				return fmt.Errorf("could not marshal json: %w", err)
			}
			r.JSONSchema = jsontypes.NewNormalizedValue(string(jsonSchema))

			r.Version = types.Float64Value(float64(rule.Version))

			t.Rules = append(t.Rules, r)
		}
	}

	return nil
}

func (t *TrackingPlanDSState) Fill(trackingPlan api.TrackingPlan, rules *[]api.RuleV1) error {
	t.ID = types.StringValue(trackingPlan.Id)
	t.Name = types.StringPointerValue(trackingPlan.Name)
	t.Slug = types.StringPointerValue(trackingPlan.Slug)
	t.Description = types.StringPointerValue(trackingPlan.Description)
	t.Type = types.StringValue(trackingPlan.Type)
	t.UpdatedAt = types.StringPointerValue(trackingPlan.UpdatedAt)
	t.CreatedAt = types.StringPointerValue(trackingPlan.CreatedAt)

	t.Rules = []RulesDSState{}
	if rules != nil {
		for _, rule := range *rules {
			r := RulesDSState{}

			r.Type = types.StringValue(rule.Type)
			r.Key = types.StringPointerValue(rule.Key)
			jsonSchema, err := json.Marshal(rule.JsonSchema)
			if err != nil {
				return fmt.Errorf("could not marshal json: %w", err)
			}
			r.JSONSchema = jsontypes.NewNormalizedValue(string(jsonSchema))
			r.Version = types.Float64Value(float64(rule.Version))
			r.CreatedAt = types.StringPointerValue(rule.CreatedAt)
			r.UpdatedAt = types.StringPointerValue(rule.UpdatedAt)
			r.DeprecatedAt = types.StringPointerValue(rule.DeprecatedAt)

			t.Rules = append(t.Rules, r)
		}
	}

	return nil
}

type RulesState struct {
	Type       types.String         `tfsdk:"type"`
	Key        types.String         `tfsdk:"key"`
	JSONSchema jsontypes.Normalized `tfsdk:"json_schema"`
	Version    types.Float64        `tfsdk:"version"`
}

type RulesDSState struct {
	Type         types.String         `tfsdk:"type"`
	Key          types.String         `tfsdk:"key"`
	JSONSchema   jsontypes.Normalized `tfsdk:"json_schema"`
	Version      types.Float64        `tfsdk:"version"`
	CreatedAt    types.String         `tfsdk:"created_at"`
	UpdatedAt    types.String         `tfsdk:"updated_at"`
	DeprecatedAt types.String         `tfsdk:"deprecated_at"`
}

func (r *RulesState) ToAPIRule() (api.RuleV1, diag.Diagnostics) {
	var jsonSchema interface{}
	diags := r.JSONSchema.Unmarshal(&jsonSchema)
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

func (r *RulesState) ToAPIRuleInput() (api.RuleInputV1, diag.Diagnostics) {
	var jsonSchema interface{}
	diags := r.JSONSchema.Unmarshal(&jsonSchema)
	if diags.HasError() {
		return api.RuleInputV1{}, diags
	}

	return api.RuleInputV1{
		Type:       r.Type.ValueString(),
		Key:        r.Key.ValueStringPointer(),
		Version:    float32(r.Version.ValueFloat64()),
		JsonSchema: jsonSchema,
	}, diags
}
