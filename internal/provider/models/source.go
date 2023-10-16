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
	Labels      types.Set            `tfsdk:"labels"`
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

func (s *SourceState) Fill(source api.Source4) error {
	s.ID = types.StringValue(source.Id)
	s.Name = types.StringPointerValue(source.Name)
	s.Slug = types.StringValue(source.Slug)
	s.Enabled = types.BoolValue(source.Enabled)
	s.WorkspaceID = types.StringValue(source.WorkspaceId)
	s.WriteKeys = s.getWriteKeys(source.WriteKeys)
	s.Labels = s.getLabels(source.Labels)
	s.Metadata = &SourceMetadataState{}
	err := s.Metadata.Fill(api.SourceMetadata(source.Metadata))
	if err != nil {
		return err
	}
	settings, err := GetSettings(source.Settings)
	if err != nil {
		return err
	}
	s.Settings = settings

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

func GetSettings(settings api.NullableModelMap) (jsontypes.Normalized, error) {
	if !settings.IsSet() {
		return jsontypes.NewNormalizedNull(), nil
	}

	jsonSettingsString, err := json.Marshal(settings.Get().Get())
	if err != nil {
		return jsontypes.NewNormalizedNull(), err
	}

	return jsontypes.NewNormalizedValue(string(jsonSettingsString)), nil
}
