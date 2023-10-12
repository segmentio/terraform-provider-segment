package models

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/segmentio/public-api-sdk-go/api"
)

type Component struct {
	Type  types.String `tfsdk:"type"`
	Code  types.String `tfsdk:"code"`
	Owner types.String `tfsdk:"owner"`
}

type SupportedFeature struct {
	CloudModeInstances     types.String `tfsdk:"cloud_mode_instances"`
	DeviceModeInstances    types.String `tfsdk:"device_mode_instances"`
	Replay                 types.Bool   `tfsdk:"replay"`
	BrowseUnbundling       types.Bool   `tfsdk:"browser_unbundling"`
	BrowseUnbundlingPublic types.Bool   `tfsdk:"browser_unbundling_public"`
}

type SupportedMethod struct {
	PageView types.Bool `tfsdk:"pageview"`
	Identify types.Bool `tfsdk:"identify"`
	Alias    types.Bool `tfsdk:"alias"`
	Track    types.Bool `tfsdk:"track"`
	Group    types.Bool `tfsdk:"group"`
}

type SupportedPlatform struct {
	Browser types.Bool `tfsdk:"browser"`
	Server  types.Bool `tfsdk:"server"`
	Mobile  types.Bool `tfsdk:"mobile"`
}

type Field struct {
	ID           types.String         `tfsdk:"id"`
	SortOrder    types.Float64        `tfsdk:"sort_order"`
	FieldKey     types.String         `tfsdk:"field_key"`
	Label        types.String         `tfsdk:"label"`
	Type         types.String         `tfsdk:"type"`
	Description  types.String         `tfsdk:"description"`
	Placeholder  types.String         `tfsdk:"placeholder"`
	Required     types.Bool           `tfsdk:"required"`
	Multiple     types.Bool           `tfsdk:"multiple"`
	Dynamic      types.Bool           `tfsdk:"dynamic"`
	AllowNull    types.Bool           `tfsdk:"allow_null"`
	DefaultValue jsontypes.Normalized `tfsdk:"default_value"`
	Choices      jsontypes.Normalized `tfsdk:"choices"`
}

type Action struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Slug           types.String `tfsdk:"slug"`
	Description    types.String `tfsdk:"description"`
	Platform       types.String `tfsdk:"platform"`
	Hidden         types.Bool   `tfsdk:"hidden"`
	DefaultTrigger types.String `tfsdk:"default_trigger"`
	Fields         []Field      `tfsdk:"fields"`
}

type Preset struct {
	ActionID types.String         `tfsdk:"action_id"`
	Name     types.String         `tfsdk:"name"`
	Trigger  types.String         `tfsdk:"trigger"`
	Fields   jsontypes.Normalized `tfsdk:"fields"`
}

type Contact struct {
	Name      types.String `tfsdk:"name"`
	Email     types.String `tfsdk:"email"`
	Role      types.String `tfsdk:"role"`
	IsPrimary types.Bool   `tfsdk:"is_primary"`
}

type DestinationMetadataState struct {
	ID                 types.String             `tfsdk:"id"`
	Name               types.String             `tfsdk:"name"`
	Slug               types.String             `tfsdk:"slug"`
	Description        types.String             `tfsdk:"description"`
	Logos              *LogosState              `tfsdk:"logos"`
	Options            []IntegrationOptionState `tfsdk:"options"`
	Categories         []types.String           `tfsdk:"categories"`
	Website            types.String             `tfsdk:"website"`
	Components         []Component              `tfsdk:"components"`
	PreviousNames      []types.String           `tfsdk:"previous_names"`
	Status             types.String             `tfsdk:"status"`
	SupportedFeatures  *SupportedFeature        `tfsdk:"supported_features"`
	SupportedMethods   *SupportedMethod         `tfsdk:"supported_methods"`
	SupportedPlatforms *SupportedPlatform       `tfsdk:"supported_platforms"`
	Actions            []Action                 `tfsdk:"actions"`
	Presets            []Preset                 `tfsdk:"presets"`
	Contacts           []Contact                `tfsdk:"contacts"`
	PartnerOwned       types.Bool               `tfsdk:"partner_owned"`
	SupportedRegions   []types.String           `tfsdk:"supported_regions"`
	RegionEndpoints    []types.String           `tfsdk:"region_endpoints"`
}

func (d *DestinationMetadataState) getPartnerOwned(owned *bool) types.Bool {
	var partnerOwned types.Bool

	if owned != nil {
		partnerOwned = types.BoolValue(*owned)
	}

	return partnerOwned
}

func (d *DestinationMetadataState) getSupportedPlatforms(platforms api.SupportedPlatforms) *SupportedPlatform {
	var supportedPlatform SupportedPlatform

	if platforms.Browser != nil {
		supportedPlatform.Browser = types.BoolValue(*platforms.Browser)
	}

	if platforms.Server != nil {
		supportedPlatform.Server = types.BoolValue(*platforms.Server)
	}

	if platforms.Mobile != nil {
		supportedPlatform.Mobile = types.BoolValue(*platforms.Mobile)
	}

	return &supportedPlatform
}

func (d *DestinationMetadataState) getSupportedFeatures(features api.SupportedFeatures) *SupportedFeature {
	var supportedFeature SupportedFeature

	if features.CloudModeInstances != nil {
		supportedFeature.CloudModeInstances = types.StringValue(*features.CloudModeInstances)
	}

	if features.DeviceModeInstances != nil {
		supportedFeature.DeviceModeInstances = types.StringValue(*features.DeviceModeInstances)
	}

	if features.Replay != nil {
		supportedFeature.Replay = types.BoolValue(*features.Replay)
	}

	if features.BrowserUnbundling != nil {
		supportedFeature.BrowseUnbundling = types.BoolValue(*features.BrowserUnbundling)
	}

	if features.BrowserUnbundlingPublic != nil {
		supportedFeature.BrowseUnbundlingPublic = types.BoolValue(*features.BrowserUnbundlingPublic)
	}

	return &supportedFeature
}

func (d *DestinationMetadataState) getSupportedMethods(methods api.SupportedMethods) *SupportedMethod {
	var supportedMethod SupportedMethod

	if methods.Pageview != nil {
		supportedMethod.PageView = types.BoolValue(*methods.Pageview)
	}

	if methods.Identify != nil {
		supportedMethod.Identify = types.BoolValue(*methods.Identify)
	}

	if methods.Alias != nil {
		supportedMethod.Alias = types.BoolValue(*methods.Alias)
	}

	if methods.Track != nil {
		supportedMethod.Track = types.BoolValue(*methods.Track)
	}

	if methods.Group != nil {
		supportedMethod.Group = types.BoolValue(*methods.Group)
	}

	return &supportedMethod
}

func (d *DestinationMetadataState) getPreviousNames(names []string) []types.String {
	var previousNames []types.String

	for _, name := range names {
		previousNames = append(previousNames, types.StringValue(name))
	}

	return previousNames
}

func (d *DestinationMetadataState) getComponents(components []api.DestinationMetadataComponentV1) []Component {
	var componentsToAdd []Component

	for _, c := range components {
		componentToAdd := Component{
			Type: types.StringValue(c.Type),
			Code: types.StringValue(c.Code),
		}

		if c.Owner != nil {
			componentToAdd.Owner = types.StringValue(*c.Owner)
		}

		componentsToAdd = append(componentsToAdd, componentToAdd)
	}

	return componentsToAdd
}

func (d *DestinationMetadataState) getRegionEndpoints(endpoints []string) []types.String {
	var regionEndpoints []types.String

	for _, endpoint := range endpoints {
		regionEndpoints = append(regionEndpoints, types.StringValue(endpoint))
	}

	return regionEndpoints
}

func (d *DestinationMetadataState) getSupportedRegions(regions []string) []types.String {
	var supportedRegionsToAdd []types.String

	for _, region := range regions {
		supportedRegionsToAdd = append(supportedRegionsToAdd, types.StringValue(region))
	}

	return supportedRegionsToAdd
}

func (d *DestinationMetadataState) getContacts(contacts []api.Contact) []Contact {
	var contactsToAdd []Contact

	for _, c := range contacts {
		contactToAdd := Contact{
			Name:      types.StringValue(*c.Name),
			Email:     types.StringValue(c.Email),
			Role:      types.StringValue(*c.Role),
			IsPrimary: types.BoolValue(*c.IsPrimary),
		}

		contactsToAdd = append(contactsToAdd, contactToAdd)
	}

	return contactsToAdd
}

func (d *DestinationMetadataState) getPresets(presets []api.DestinationMetadataSubscriptionPresetV1) ([]Preset, error) {
	var presetsToAdd []Preset

	for _, preset := range presets {
		presetToAdd := Preset{
			ActionID: types.StringValue(preset.ActionId),
			Name:     types.StringValue(preset.Name),
			Trigger:  types.StringValue(preset.Trigger),
		}

		fields, err := json.Marshal(preset.Fields)
		if err != nil {
			return []Preset{}, fmt.Errorf("could not marshal json: %w", err)
		}
		presetToAdd.Fields = jsontypes.NewNormalizedValue(string(fields))

		presetsToAdd = append(presetsToAdd, presetToAdd)
	}

	return presetsToAdd, nil
}

func (d *DestinationMetadataState) getActions(actions []api.DestinationMetadataActionV1) ([]Action, error) {
	var destinationActionsToAdd []Action

	for _, action := range actions {
		destinationMetadataAction := Action{
			ID:          types.StringValue(action.Id),
			Slug:        types.StringValue(action.Slug),
			Name:        types.StringValue(action.Name),
			Description: types.StringValue(action.Description),
			Platform:    types.StringValue(action.Platform),
			Hidden:      types.BoolValue(action.Hidden),
		}

		fields, err := d.getFields(action.Fields)
		if err != nil {
			return []Action{}, err
		}
		destinationMetadataAction.Fields = fields

		if action.DefaultTrigger.IsSet() {
			destinationMetadataAction.DefaultTrigger = types.StringPointerValue(action.DefaultTrigger.Get())
		}

		destinationActionsToAdd = append(destinationActionsToAdd, destinationMetadataAction)
	}

	return destinationActionsToAdd, nil
}

func (d *DestinationMetadataState) getFields(fields []api.DestinationMetadataActionFieldV1) ([]Field, error) {
	var fieldsToAdd []Field

	for _, f := range fields {
		fieldToAdd := Field{
			ID:          types.StringValue(f.Id),
			SortOrder:   types.Float64Value(float64(f.SortOrder)),
			FieldKey:    types.StringValue(f.FieldKey),
			Label:       types.StringValue(f.Label),
			Type:        types.StringValue(f.Type),
			Description: types.StringValue(f.Description),
			Required:    types.BoolValue(f.Required),
			Multiple:    types.BoolValue(f.Multiple),
			Dynamic:     types.BoolValue(f.Dynamic),
			AllowNull:   types.BoolValue(f.AllowNull),
		}

		if f.Placeholder != nil {
			fieldToAdd.Placeholder = types.StringValue(*f.Placeholder)
		}

		if f.DefaultValue != nil {
			defaultValue, err := json.Marshal(f.DefaultValue)
			if err != nil {
				return []Field{}, fmt.Errorf("could not marshal json: %w", err)
			}
			fieldToAdd.DefaultValue = jsontypes.NewNormalizedValue(string(defaultValue))
		}

		choices, err := json.Marshal(f.Choices)
		if err != nil {
			return []Field{}, fmt.Errorf("could not marshal json: %w", err)
		}
		fieldToAdd.Choices = jsontypes.NewNormalizedValue(string(choices))

		fieldsToAdd = append(fieldsToAdd, fieldToAdd)
	}

	return fieldsToAdd, nil
}

func (d *DestinationMetadataState) getLogosDestinationMetadata(logos api.Logos) *LogosState {
	logosToAdd := LogosState{
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

func getOptions(options []api.IntegrationOptionBeta) ([]IntegrationOptionState, error) {
	var integrationOptions []IntegrationOptionState

	for _, opt := range options {
		integrationOption := IntegrationOptionState{
			Name:     types.StringValue(opt.Name),
			Type:     types.StringValue(opt.Type),
			Required: types.BoolValue(opt.Required),
		}

		if opt.Description != nil {
			integrationOption.Description = types.StringValue(*opt.Description)
		}

		if opt.Label != nil {
			integrationOption.Label = types.StringValue(*opt.Label)
		}

		if opt.DefaultValue != nil {
			defaultValue, err := json.Marshal(opt.DefaultValue)
			if err != nil {
				return []IntegrationOptionState{}, fmt.Errorf("could not marshal json: %w", err)
			}
			integrationOption.DefaultValue = jsontypes.NewNormalizedValue(string(defaultValue))
		}

		integrationOptions = append(integrationOptions, integrationOption)
	}

	return integrationOptions, nil
}

func (d *DestinationMetadataState) Fill(destinationMetadata api.DestinationMetadata) error {
	d.ID = types.StringValue(destinationMetadata.Id)
	d.Name = types.StringValue(destinationMetadata.Name)
	d.Description = types.StringValue(destinationMetadata.Description)
	d.Slug = types.StringValue(destinationMetadata.Slug)
	d.Logos = d.getLogosDestinationMetadata(destinationMetadata.Logos)
	options, err := getOptions(destinationMetadata.Options)
	if err != nil {
		return err
	}
	d.Options = options
	actions, err := d.getActions(destinationMetadata.Actions)
	if err != nil {
		return err
	}
	d.Actions = actions
	d.Categories = getCategories(destinationMetadata.Categories)
	presets, err := d.getPresets(destinationMetadata.Presets)
	if err != nil {
		return err
	}
	d.Presets = presets
	d.Contacts = d.getContacts(destinationMetadata.Contacts)
	d.PartnerOwned = d.getPartnerOwned(destinationMetadata.PartnerOwned)
	d.SupportedRegions = d.getSupportedRegions(destinationMetadata.SupportedRegions)
	d.RegionEndpoints = d.getRegionEndpoints(destinationMetadata.RegionEndpoints)
	d.Status = types.StringValue(destinationMetadata.Status)
	d.Website = types.StringValue(destinationMetadata.Website)
	d.Components = d.getComponents(destinationMetadata.Components)
	d.PreviousNames = d.getPreviousNames(destinationMetadata.PreviousNames)
	d.SupportedMethods = d.getSupportedMethods(destinationMetadata.SupportedMethods)
	d.SupportedFeatures = d.getSupportedFeatures(destinationMetadata.SupportedFeatures)
	d.SupportedPlatforms = d.getSupportedPlatforms(destinationMetadata.SupportedPlatforms)

	return nil
}
