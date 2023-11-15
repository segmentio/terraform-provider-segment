package models

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/segmentio/public-api-sdk-go/api"
)

type SourcePlan struct {
	Enabled        types.Bool           `tfsdk:"enabled"`
	ID             types.String         `tfsdk:"id"`
	Labels         types.Set            `tfsdk:"labels"`
	Metadata       types.Object         `tfsdk:"metadata"`
	Name           types.String         `tfsdk:"name"`
	Slug           types.String         `tfsdk:"slug"`
	WorkspaceID    types.String         `tfsdk:"workspace_id"`
	WriteKeys      types.List           `tfsdk:"write_keys"`
	Settings       jsontypes.Normalized `tfsdk:"settings"`
	SchemaSettings types.Object         `tfsdk:"schema_settings"`
}

type SourceState struct {
	Enabled        types.Bool           `tfsdk:"enabled"`
	ID             types.String         `tfsdk:"id"`
	Labels         []LabelState         `tfsdk:"labels"`
	Metadata       *SourceMetadataState `tfsdk:"metadata"`
	Name           types.String         `tfsdk:"name"`
	Slug           types.String         `tfsdk:"slug"`
	WorkspaceID    types.String         `tfsdk:"workspace_id"`
	WriteKeys      []types.String       `tfsdk:"write_keys"`
	Settings       jsontypes.Normalized `tfsdk:"settings"`
	SchemaSettings *SchemaSettingsState `tfsdk:"schema_settings"`
}

type SchemaSettingsState struct {
	Track                     *TrackSettings    `tfsdk:"track"`
	Identify                  *IdentifySettings `tfsdk:"identify"`
	Group                     *GroupSettings    `tfsdk:"group"`
	ForwardingViolationsTo    types.String      `tfsdk:"forwarding_violations_to"`
	ForwardingBlockedEventsTo types.String      `tfsdk:"forwarding_blocked_events_to"`
}

type SchemaSettingsPlan struct {
	Track                     types.Object `tfsdk:"track"`
	Identify                  types.Object `tfsdk:"identify"`
	Group                     types.Object `tfsdk:"group"`
	ForwardingViolationsTo    types.String `tfsdk:"forwarding_violations_to"`
	ForwardingBlockedEventsTo types.String `tfsdk:"forwarding_blocked_events_to"`
}

type TrackSettings struct {
	AllowUnplannedEvents          types.Bool   `tfsdk:"allow_unplanned_events"`
	AllowUnplannedEventProperties types.Bool   `tfsdk:"allow_unplanned_event_properties"`
	AllowEventOnViolations        types.Bool   `tfsdk:"allow_event_on_violations"`
	AllowPropertiesOnViolations   types.Bool   `tfsdk:"allow_properties_on_violations"`
	CommonEventOnViolations       types.String `tfsdk:"common_event_on_violations"`
}

type IdentifySettings struct {
	AllowUnplannedTraits    types.Bool   `tfsdk:"allow_unplanned_traits"`
	AllowTraitsOnViolations types.Bool   `tfsdk:"allow_traits_on_violations"`
	CommonEventOnViolations types.String `tfsdk:"common_event_on_violations"`
}

type GroupSettings struct {
	AllowUnplannedTraits    types.Bool   `tfsdk:"allow_unplanned_traits"`
	AllowTraitsOnViolations types.Bool   `tfsdk:"allow_traits_on_violations"`
	CommonEventOnViolations types.String `tfsdk:"common_event_on_violations"`
}

func (s *SourceState) Fill(source api.SourceV1, schemaSettings *api.SourceSettingsOutputV1) error {
	s.ID = types.StringValue(source.Id)
	s.Name = types.StringPointerValue(source.Name)
	s.Slug = types.StringValue(source.Slug)
	s.Enabled = types.BoolValue(source.Enabled)
	s.WorkspaceID = types.StringValue(source.WorkspaceId)
	s.WriteKeys = s.getWriteKeys(source.WriteKeys)
	s.Labels = s.getLabels(source.Labels)
	s.Metadata = &SourceMetadataState{}
	err := s.Metadata.Fill(source.Metadata)
	if err != nil {
		return err
	}
	settings, err := GetSettings(source.Settings)
	if err != nil {
		return err
	}
	s.Settings = settings

	if schemaSettings != nil {
		s.SchemaSettings = &SchemaSettingsState{}
		s.SchemaSettings.Fill(*schemaSettings)
	}

	return nil
}

func (s *SourceState) getLabels(labels []api.LabelV1) []LabelState {
	var labelsToAdd []LabelState

	for _, label := range labels {
		labelToAdd := LabelState{}
		labelToAdd.Fill(label)

		labelsToAdd = append(labelsToAdd, labelToAdd)
	}

	return labelsToAdd
}

func (s *SourceState) getWriteKeys(writeKeys []string) []types.String {
	var writeKeysToAdd []types.String

	for _, writeKey := range writeKeys {
		writeKeysToAdd = append(writeKeysToAdd, types.StringValue(writeKey))
	}

	return writeKeysToAdd
}

func (s *SchemaSettingsState) Fill(schemaSettings api.SourceSettingsOutputV1) {
	s.Track = &TrackSettings{}
	s.Track.Fill(schemaSettings.Track)

	s.Identify = &IdentifySettings{}
	s.Identify.Fill(schemaSettings.Identify)

	s.Group = &GroupSettings{}
	s.Group.Fill(schemaSettings.Group)
	s.ForwardingViolationsTo = types.StringPointerValue(schemaSettings.ForwardingViolationsTo)
	s.ForwardingBlockedEventsTo = types.StringPointerValue(schemaSettings.ForwardingBlockedEventsTo)
}

func (t *TrackSettings) Fill(trackSettings *api.TrackSourceSettingsV1) {
	if trackSettings == nil {
		return
	}

	t.AllowUnplannedEvents = types.BoolPointerValue(trackSettings.AllowUnplannedEvents)
	t.AllowUnplannedEventProperties = types.BoolPointerValue(trackSettings.AllowUnplannedEventProperties)
	t.AllowEventOnViolations = types.BoolPointerValue(trackSettings.AllowEventOnViolations)
	t.AllowPropertiesOnViolations = types.BoolPointerValue(trackSettings.AllowPropertiesOnViolations)
	t.CommonEventOnViolations = types.StringPointerValue(trackSettings.CommonEventOnViolations)
}

func (i *IdentifySettings) Fill(identifySettings *api.IdentifySourceSettingsV1) {
	if identifySettings == nil {
		return
	}

	i.AllowUnplannedTraits = types.BoolPointerValue(identifySettings.AllowUnplannedTraits)
	i.AllowTraitsOnViolations = types.BoolPointerValue(identifySettings.AllowTraitsOnViolations)
	i.CommonEventOnViolations = types.StringPointerValue(identifySettings.CommonEventOnViolations)
}

func (g *GroupSettings) Fill(groupSettings *api.GroupSourceSettingsV1) {
	if groupSettings == nil {
		return
	}

	g.AllowUnplannedTraits = types.BoolPointerValue(groupSettings.AllowUnplannedTraits)
	g.AllowTraitsOnViolations = types.BoolPointerValue(groupSettings.AllowTraitsOnViolations)
	g.CommonEventOnViolations = types.StringPointerValue(groupSettings.CommonEventOnViolations)
}

func GetSettings(settings map[string]interface{}) (jsontypes.Normalized, error) {
	if settings == nil {
		return jsontypes.NewNormalizedNull(), nil
	}

	jsonSettingsString, err := json.Marshal(settings)
	if err != nil {
		return jsontypes.NewNormalizedNull(), err
	}

	return jsontypes.NewNormalizedValue(string(jsonSettingsString)), nil
}

func GetSchemaSettingsFromPlan(ctx context.Context, settings types.Object) (*api.SourceSettingsOutputV1, diag.Diagnostics) {
	if settings.IsNull() || settings.IsUnknown() {
		return nil, nil
	}

	var schemaSettingsPlan SchemaSettingsPlan
	diags := settings.As(ctx, &schemaSettingsPlan, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil, diags
	}

	var apiTrackSettings *api.TrackSourceSettingsV1
	if !schemaSettingsPlan.Track.IsNull() && !schemaSettingsPlan.Track.IsUnknown() {
		var trackSettings TrackSettings
		diags = schemaSettingsPlan.Track.As(ctx, &trackSettings, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			return nil, diags
		}

		apiTrackSettings = &api.TrackSourceSettingsV1{
			AllowUnplannedEvents:          trackSettings.AllowUnplannedEvents.ValueBoolPointer(),
			AllowUnplannedEventProperties: trackSettings.AllowUnplannedEventProperties.ValueBoolPointer(),
			AllowEventOnViolations:        trackSettings.AllowEventOnViolations.ValueBoolPointer(),
			AllowPropertiesOnViolations:   trackSettings.AllowPropertiesOnViolations.ValueBoolPointer(),
			CommonEventOnViolations:       trackSettings.CommonEventOnViolations.ValueStringPointer(),
		}
	}

	var apiIdentifySettings *api.IdentifySourceSettingsV1
	if !schemaSettingsPlan.Identify.IsNull() && !schemaSettingsPlan.Identify.IsUnknown() {
		var identifySettings IdentifySettings
		diags = schemaSettingsPlan.Identify.As(ctx, &identifySettings, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			return nil, diags
		}

		apiIdentifySettings = &api.IdentifySourceSettingsV1{
			AllowUnplannedTraits:    identifySettings.AllowUnplannedTraits.ValueBoolPointer(),
			AllowTraitsOnViolations: identifySettings.AllowTraitsOnViolations.ValueBoolPointer(),
			CommonEventOnViolations: identifySettings.CommonEventOnViolations.ValueStringPointer(),
		}
	}

	var apiGroupSettings *api.GroupSourceSettingsV1
	if !schemaSettingsPlan.Group.IsNull() && !schemaSettingsPlan.Group.IsUnknown() {
		var groupSettings GroupSettings
		diags = schemaSettingsPlan.Group.As(ctx, &groupSettings, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			return nil, diags
		}

		apiGroupSettings = &api.GroupSourceSettingsV1{
			AllowUnplannedTraits:    groupSettings.AllowUnplannedTraits.ValueBoolPointer(),
			AllowTraitsOnViolations: groupSettings.AllowTraitsOnViolations.ValueBoolPointer(),
			CommonEventOnViolations: groupSettings.CommonEventOnViolations.ValueStringPointer(),
		}
	}

	return &api.SourceSettingsOutputV1{
		Track:                     apiTrackSettings,
		Identify:                  apiIdentifySettings,
		Group:                     apiGroupSettings,
		ForwardingViolationsTo:    schemaSettingsPlan.ForwardingViolationsTo.ValueStringPointer(),
		ForwardingBlockedEventsTo: schemaSettingsPlan.ForwardingBlockedEventsTo.ValueStringPointer(),
	}, nil
}

func SchemaSettingsPlanToState(ctx context.Context, settings types.Object) (*SchemaSettingsState, diag.Diagnostics) {
	if settings.IsNull() || settings.IsUnknown() {
		return nil, nil
	}

	var schemaSettingsPlan SchemaSettingsPlan
	diags := settings.As(ctx, &schemaSettingsPlan, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil, diags
	}

	var trackSettings *TrackSettings
	if !schemaSettingsPlan.Track.IsNull() && !schemaSettingsPlan.Track.IsUnknown() {
		var ts TrackSettings
		diags = schemaSettingsPlan.Track.As(ctx, &ts, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			return nil, diags
		}

		trackSettings = &ts
	}

	var identifySettings *IdentifySettings
	if !schemaSettingsPlan.Identify.IsNull() && !schemaSettingsPlan.Identify.IsUnknown() {
		var is IdentifySettings
		diags = schemaSettingsPlan.Identify.As(ctx, &is, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			return nil, diags
		}

		identifySettings = &is
	}

	var groupSettings *GroupSettings
	if !schemaSettingsPlan.Group.IsNull() && !schemaSettingsPlan.Group.IsUnknown() {
		var gs GroupSettings
		diags = schemaSettingsPlan.Group.As(ctx, &gs, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			return nil, diags
		}

		groupSettings = &gs
	}

	return &SchemaSettingsState{
		Track:                     trackSettings,
		Identify:                  identifySettings,
		Group:                     groupSettings,
		ForwardingViolationsTo:    schemaSettingsPlan.ForwardingViolationsTo,
		ForwardingBlockedEventsTo: schemaSettingsPlan.ForwardingBlockedEventsTo,
	}, nil
}
