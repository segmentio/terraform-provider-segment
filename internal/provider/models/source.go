package models

import (
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/segmentio/public-api-sdk-go/api"
)

type SourcePlan struct {
	Enabled     types.Bool           `tfsdk:"enabled"`
	ID          types.String         `tfsdk:"id"`
	Labels      types.List           `tfsdk:"labels"`
	Metadata    types.Object         `tfsdk:"metadata"`
	Name        types.String         `tfsdk:"name"`
	Slug        types.String         `tfsdk:"slug"`
	WorkspaceID types.String         `tfsdk:"workspace_id"`
	WriteKeys   types.List           `tfsdk:"write_keys"`
	Settings    jsontypes.Normalized `tfsdk:"settings"`
}

type SourceState struct {
	Enabled     types.Bool           `tfsdk:"enabled"`
	ID          types.String         `tfsdk:"id"`
	Labels      []LabelState         `tfsdk:"labels"`
	Metadata    *SourceMetadataState `tfsdk:"metadata"`
	Name        types.String         `tfsdk:"name"`
	Slug        types.String         `tfsdk:"slug"`
	WorkspaceID types.String         `tfsdk:"workspace_id"`
	WriteKeys   []types.String       `tfsdk:"write_keys"`
	Settings    jsontypes.Normalized `tfsdk:"settings"`
}

type LabelState struct {
	Description types.String `tfsdk:"description"`
	Key         types.String `tfsdk:"key"`
	Value       types.String `tfsdk:"value"`
}

func (s *SourceState) Fill(source api.Source4) error {
	s.ID = types.StringValue(source.Id)
	if source.Name != nil {
		s.Name = types.StringValue(*source.Name)
	}
	s.Slug = types.StringValue(source.Slug)
	s.Enabled = types.BoolValue(source.Enabled)
	s.WorkspaceID = types.StringValue(source.WorkspaceId)
	s.WriteKeys = s.getWriteKeys(source.WriteKeys)
	s.Labels = s.getLabels(source.Labels)
	s.Metadata = s.getMetadata(source.Metadata)
	settings, err := s.getSettings(source.Settings)
	if err != nil {
		return err
	}
	s.Settings = settings

	return nil
}

func (s *SourceState) getLogos(logos api.Logos1) *LogosState {
	logosToAdd := LogosState{
		Default: types.StringValue(logos.Default),
	}
	if logos.Alt.IsSet() {
		logosToAdd.Alt = types.StringValue(*logos.Alt.Get())
	}
	if logos.Mark.IsSet() {
		logosToAdd.Mark = types.StringValue(*logos.Mark.Get())
	}

	return &logosToAdd
}

func (s *SourceState) getMetadata(metadata api.Metadata2) *SourceMetadataState {
	metadataToAdd := SourceMetadataState{
		ID:                 types.StringValue(metadata.Id),
		Description:        types.StringValue(metadata.Description),
		IsCloudEventSource: types.BoolValue(metadata.IsCloudEventSource),
		Logos:              s.getLogos(metadata.Logos),
		Name:               types.StringValue(metadata.Name),
		Slug:               types.StringValue(metadata.Slug),
	}

	for _, metadataCategory := range metadata.Categories {
		metadataToAdd.Categories = append(metadataToAdd.Categories, types.StringValue(metadataCategory))
	}

	for _, integrationOption := range metadata.Options {
		integrationOptionToAdd := IntegrationOptionState{
			Name:     types.StringValue(integrationOption.Name),
			Type:     types.StringValue(integrationOption.Type),
			Required: types.BoolValue(integrationOption.Required),
		}

		if integrationOption.Description != nil {
			integrationOptionToAdd.Description = types.StringValue(*integrationOption.Description)
		}

		// TODO handle integrationOption.DefaultValue (typed as interface{})

		metadataToAdd.Options = append(metadataToAdd.Options, integrationOptionToAdd)
	}

	return &metadataToAdd
}

func (s *SourceState) getLabels(labels []api.LabelV1) []LabelState {
	var labelsToAdd []LabelState

	for _, label := range labels {
		labelToAdd := LabelState{
			Key:   types.StringValue(label.Key),
			Value: types.StringValue(label.Value),
		}

		if label.Description != nil {
			labelToAdd.Description = types.StringValue(*label.Description)
		}

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

func (s *SourceState) getSettings(settings api.NullableModelMap) (jsontypes.Normalized, error) {
	if !settings.IsSet() {
		return jsontypes.NewNormalizedNull(), nil
	}

	jsonSettingsString, err := json.Marshal(settings.Get().Get())
	if err != nil {
		return jsontypes.NewNormalizedNull(), err
	}

	return jsontypes.NewNormalizedValue(string(jsonSettingsString)), nil
}
