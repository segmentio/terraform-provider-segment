package models

import (
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/segmentio/public-api-sdk-go/api"
)

type SourceMetadataState struct {
	ID                 types.String             `tfsdk:"id"`
	Name               types.String             `tfsdk:"name"`
	Slug               types.String             `tfsdk:"slug"`
	Description        types.String             `tfsdk:"description"`
	Logos              *LogosState              `tfsdk:"logos"`
	Options            []IntegrationOptionState `tfsdk:"options"`
	Categories         []types.String           `tfsdk:"categories"`
	IsCloudEventSource types.Bool               `tfsdk:"is_cloud_event_source"`
}

type IntegrationOptionState struct {
	DefaultValue jsontypes.Normalized `tfsdk:"default_value"`
	Description  types.String         `tfsdk:"description"`
	Label        types.String         `tfsdk:"label"`
	Name         types.String         `tfsdk:"name"`
	Required     types.Bool           `tfsdk:"required"`
	Type         types.String         `tfsdk:"type"`
}

type LogosState struct {
	Alt     types.String `tfsdk:"alt"`
	Default types.String `tfsdk:"default"`
	Mark    types.String `tfsdk:"mark"`
}

func (s *SourceMetadataState) Fill(sourceMetadata api.SourceMetadataV1) error {
	s.ID = types.StringValue(sourceMetadata.Id)
	s.Name = types.StringValue(sourceMetadata.Name)
	s.Description = types.StringValue(sourceMetadata.Description)
	s.Slug = types.StringValue(sourceMetadata.Slug)
	s.Logos = getLogos(sourceMetadata.Logos)
	options, err := getOptions(sourceMetadata.Options)
	if err != nil {
		return err
	}
	s.Options = options
	s.IsCloudEventSource = types.BoolValue(sourceMetadata.IsCloudEventSource)
	s.Categories = getCategories(sourceMetadata.Categories)

	return nil
}

func getCategories(categories []string) []types.String {
	var categoriesToAdd []types.String

	for _, cat := range categories {
		categoriesToAdd = append(categoriesToAdd, types.StringValue(cat))
	}

	return categoriesToAdd
}

func getLogos(logos api.LogosBeta) *LogosState {
	logosToAdd := LogosState{
		Default: types.StringValue(logos.Default),
	}

	if logos.Mark.IsSet() {
		logosToAdd.Mark = types.StringPointerValue(logos.Mark.Get())
	}

	if logos.Alt.IsSet() {
		logosToAdd.Alt = types.StringPointerValue(logos.Alt.Get())
	}

	return &logosToAdd
}
