package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/segmentio/public-api-sdk-go/api"
)

type SourcePlanModel struct {
	Enabled     types.Bool   `tfsdk:"enabled"`
	ID          types.String `tfsdk:"id"`
	Labels      types.List   `tfsdk:"labels"`
	Metadata    types.Object `tfsdk:"metadata"`
	Name        types.String `tfsdk:"name"`
	Slug        types.String `tfsdk:"slug"`
	WorkspaceID types.String `tfsdk:"workspace_id"`
	WriteKeys   types.List   `tfsdk:"write_keys"`
}

type SourceStateModel struct {
	Enabled     types.Bool                `tfsdk:"enabled"`
	ID          types.String              `tfsdk:"id"`
	Labels      []LabelStateModel         `tfsdk:"labels"`
	Metadata    *SourceMetadataStateModel `tfsdk:"metadata"`
	Name        types.String              `tfsdk:"name"`
	Slug        types.String              `tfsdk:"slug"`
	WorkspaceID types.String              `tfsdk:"workspace_id"`
	WriteKeys   []types.String            `tfsdk:"write_keys"`
}

type LabelStateModel struct {
	Description types.String `tfsdk:"description"`
	Key         types.String `tfsdk:"key"`
	Value       types.String `tfsdk:"value"`
}

type SourceMetadataStateModel struct {
	Categories         []types.String                `tfsdk:"categories"`
	Description        types.String                  `tfsdk:"description"`
	ID                 types.String                  `tfsdk:"id"`
	IsCloudEventSource types.Bool                    `tfsdk:"is_cloud_event_source"`
	Logos              *LogosStateModel              `tfsdk:"logos"`
	Name               types.String                  `tfsdk:"name"`
	Options            []IntegrationOptionStateModel `tfsdk:"options"`
	Slug               types.String                  `tfsdk:"slug"`
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

func (d *SourceStateModel) Get(source api.Source4) {
	d.ID = types.StringValue(source.Id)
	if source.Name != nil {
		d.Name = types.StringValue(*source.Name)
	}
	d.Slug = types.StringValue(source.Slug)
	d.Enabled = types.BoolValue(source.Enabled)
	d.WorkspaceID = types.StringValue(source.WorkspaceId)
	d.WriteKeys = d.getWriteKeys(source.WriteKeys)
	d.Labels = d.getLabels(source.Labels)
	d.Metadata = d.getMetadata(source.Metadata)
	// TODO: Populate settings
}

func (d *SourceStateModel) getLogos(logos api.Logos1) *LogosStateModel {
	logosToAdd := LogosStateModel{
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

func (d *SourceStateModel) getMetadata(metadata api.Metadata2) *SourceMetadataStateModel {
	metadataToAdd := SourceMetadataStateModel{
		ID:                 types.StringValue(metadata.Id),
		Description:        types.StringValue(metadata.Description),
		IsCloudEventSource: types.BoolValue(metadata.IsCloudEventSource),
		Logos:              d.getLogos(metadata.Logos),
		Name:               types.StringValue(metadata.Name),
		Slug:               types.StringValue(metadata.Slug),
	}

	for _, metadataCategory := range metadata.Categories {
		metadataToAdd.Categories = append(metadataToAdd.Categories, types.StringValue(metadataCategory))
	}

	for _, integrationOption := range metadata.Options {
		integrationOptionToAdd := IntegrationOptionStateModel{
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

func (d *SourceStateModel) getLabels(labels []api.LabelV1) []LabelStateModel {
	var labelsToAdd []LabelStateModel

	for _, label := range labels {
		labelToAdd := LabelStateModel{
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

func (d *SourceStateModel) getWriteKeys(writeKeys []string) []types.String {
	var writeKeysToAdd []types.String

	for _, writeKey := range writeKeys {
		writeKeysToAdd = append(writeKeysToAdd, types.StringValue(writeKey))
	}

	return writeKeysToAdd
}
