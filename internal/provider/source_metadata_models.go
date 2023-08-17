package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/segmentio/public-api-sdk-go/api"
)

type SourceMetadataStateModel struct {
	ID                 types.String                  `tfsdk:"id"`
	Name               types.String                  `tfsdk:"name"`
	Slug               types.String                  `tfsdk:"slug"`
	Description        types.String                  `tfsdk:"description"`
	Logos              *LogosStateModel              `tfsdk:"logos"`
	Options            []IntegrationOptionStateModel `tfsdk:"options"`
	Categories         []types.String                `tfsdk:"categories"`
	IsCloudEventSource types.Bool                    `tfsdk:"is_cloud_event_source"`
}

type IntegrationOptionStateModel struct {
	// TODO: DefaultValue types.String `tfsdk:"default_value"`
	Description types.String `tfsdk:"description"`
	Label       types.String `tfsdk:"label"`
	Name        types.String `tfsdk:"name"`
	Required    types.Bool   `tfsdk:"required"`
	Type        types.String `tfsdk:"type"`
}

type LogosStateModel struct {
	Alt     types.String `tfsdk:"alt"`
	Default types.String `tfsdk:"default"`
	Mark    types.String `tfsdk:"mark"`
}

func (s *SourceMetadataStateModel) Fill(sourceMetadata api.SourceMetadata) {
	s.ID = types.StringValue(sourceMetadata.Id)
	s.Name = types.StringValue(sourceMetadata.Name)
	s.Description = types.StringValue(sourceMetadata.Description)
	s.Slug = types.StringValue(sourceMetadata.Slug)
	s.Logos = getLogosSourceMetadata(sourceMetadata.Logos)
	s.Options = getOptions(sourceMetadata.Options)
	s.IsCloudEventSource = types.BoolValue(sourceMetadata.IsCloudEventSource)
	s.Categories = getCategories(sourceMetadata.Categories)
}

func getCategories(categories []string) []types.String {
	var categoriesToAdd []types.String

	for _, cat := range categories {
		categoriesToAdd = append(categoriesToAdd, types.StringValue(cat))
	}

	return categoriesToAdd
}

func getLogosSourceMetadata(logos api.Logos1) *LogosStateModel {
	logosToAdd := LogosStateModel{
		Default: types.StringValue(logos.Default),
	}

	if logos.Mark.IsSet() {
		logosToAdd.Mark = types.StringValue(*logos.Mark.Get())
	}

	if logos.Alt.IsSet() {
		logosToAdd.Alt = types.StringValue(*logos.Alt.Get())
	}

	return &logosToAdd
}
