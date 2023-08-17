package provider

import (
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
	Id          types.String  `tfsdk:"id"`
	SortOrder   types.Float64 `tfsdk:"sort_order"`
	FieldKey    types.String  `tfsdk:"field_key"`
	Label       types.String  `tfsdk:"label"`
	Type        types.String  `tfsdk:"type"`
	Description types.String  `tfsdk:"description"`
	Placeholder types.String  `tfsdk:"placeholder"`
	Required    types.Bool    `tfsdk:"required"`
	Multiple    types.Bool    `tfsdk:"multiple"`
	Dynamic     types.Bool    `tfsdk:"dynamic"`
	AllowNull   types.Bool    `tfsdk:"allow_null"`
}

type Action struct {
	Id             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Slug           types.String `tfsdk:"slug"`
	Description    types.String `tfsdk:"description"`
	Platform       types.String `tfsdk:"platform"`
	Hidden         types.Bool   `tfsdk:"hidden"`
	DefaultTrigger types.String `tfsdk:"default_trigger"`
	Fields         []Field      `tfsdk:"fields"`
}

type Preset struct {
	ActionId types.String `tfsdk:"action_id"`
	Name     types.String `tfsdk:"name"`
	Trigger  types.String `tfsdk:"trigger"`
}

type Contact struct {
	Name      types.String `tfsdk:"name"`
	Email     types.String `tfsdk:"email"`
	Role      types.String `tfsdk:"role"`
	IsPrimary types.Bool   `tfsdk:"is_primary"`
}

type DestinationMetadataStateModel struct {
	Id                 types.String                  `tfsdk:"id"`
	Name               types.String                  `tfsdk:"name"`
	Slug               types.String                  `tfsdk:"slug"`
	Description        types.String                  `tfsdk:"description"`
	Logos              *LogosStateModel              `tfsdk:"logos"`
	Options            []IntegrationOptionStateModel `tfsdk:"options"`
	Categories         []types.String                `tfsdk:"categories"`
	Website            types.String                  `tfsdk:"website"`
	Components         []Component                   `tfsdk:"components"`
	PreviousNames      []types.String                `tfsdk:"previous_names"`
	Status             types.String                  `tfsdk:"status"`
	SupportedFeatures  *SupportedFeature             `tfsdk:"supported_features"`
	SupportedMethods   *SupportedMethod              `tfsdk:"supported_methods"`
	SupportedPlatforms *SupportedPlatform            `tfsdk:"supported_platforms"`
	Actions            []Action                      `tfsdk:"actions"`
	Presets            []Preset                      `tfsdk:"presets"`
	Contacts           []Contact                     `tfsdk:"contacts"`
	PartnerOwned       types.Bool                    `tfsdk:"partner_owned"`
	SupportedRegions   []types.String                `tfsdk:"supported_regions"`
	RegionEndpoints    []types.String                `tfsdk:"region_endpoints"`
}

func (d *DestinationMetadataStateModel) getPartnerOwned(owned *bool) types.Bool {
	var partnerOwned types.Bool

	if owned != nil {
		partnerOwned = types.BoolValue(*owned)
	}

	return partnerOwned
}

func (d *DestinationMetadataStateModel) getSupportedPlatforms(platforms api.SupportedPlatforms) *SupportedPlatform {
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

func (d *DestinationMetadataStateModel) getSupportedFeatures(features api.SupportedFeatures) *SupportedFeature {
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

func (d *DestinationMetadataStateModel) getSupportedMethods(methods api.SupportedMethods) *SupportedMethod {
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

func (d *DestinationMetadataStateModel) getPreviousNames(names []string) []types.String {
	var previousNames []types.String

	for _, name := range names {
		previousNames = append(previousNames, types.StringValue(name))
	}

	return previousNames
}

func (d *DestinationMetadataStateModel) getComponents(components []api.DestinationMetadataComponentV1) []Component {
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

func (d *DestinationMetadataStateModel) getRegionEndpoints(endpoints []string) []types.String {
	var regionEndpoints []types.String

	for _, endpoint := range endpoints {
		regionEndpoints = append(regionEndpoints, types.StringValue(endpoint))
	}

	return regionEndpoints
}

func (d *DestinationMetadataStateModel) getSupportedRegions(regions []string) []types.String {
	var supportedRegionsToAdd []types.String

	for _, region := range regions {
		supportedRegionsToAdd = append(supportedRegionsToAdd, types.StringValue(region))
	}

	return supportedRegionsToAdd
}

func (d *DestinationMetadataStateModel) getContacts(contacts []api.Contact) []Contact {
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

func (d *DestinationMetadataStateModel) getPresets(presets []api.DestinationMetadataSubscriptionPresetV1) []Preset {
	var presetsToAdd []Preset

	for _, preset := range presets {
		presetToAdd := Preset{
			ActionId: types.StringValue(preset.ActionId),
			Name:     types.StringValue(preset.Name),
			Trigger:  types.StringValue(preset.Trigger),
		}

		presetsToAdd = append(presetsToAdd, presetToAdd)
	}

	return presetsToAdd
}

func (d *DestinationMetadataStateModel) getActions(actions []api.DestinationMetadataActionV1) []Action {
	var destinationActionsToAdd []Action

	for _, action := range actions {
		destinationMetadataAction := Action{
			Id:          types.StringValue(action.Id),
			Slug:        types.StringValue(action.Slug),
			Name:        types.StringValue(action.Name),
			Description: types.StringValue(action.Description),
			Platform:    types.StringValue(action.Platform),
			Hidden:      types.BoolValue(action.Hidden),
			Fields:      d.getFields(action.Fields),
		}

		if action.DefaultTrigger.IsSet() {
			destinationMetadataAction.DefaultTrigger = types.StringValue(*action.DefaultTrigger.Get())
		}

		destinationActionsToAdd = append(destinationActionsToAdd, destinationMetadataAction)
	}

	return destinationActionsToAdd
}

func (d *DestinationMetadataStateModel) getFields(fields []api.DestinationMetadataActionFieldV1) []Field {
	var fieldsToAdd []Field

	for _, f := range fields {
		fieldToAdd := Field{
			Id:          types.StringValue(f.Id),
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

		fieldsToAdd = append(fieldsToAdd, fieldToAdd)
	}

	return fieldsToAdd
}

func (d *DestinationMetadataStateModel) getLogosDestinationMetadata(logos api.Logos) *LogosStateModel {
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

func (d *DestinationMetadataStateModel) Fill(destinationMetadata api.DestinationMetadata) {
	d.Id = types.StringValue(destinationMetadata.Id)
	d.Name = types.StringValue(destinationMetadata.Name)
	d.Description = types.StringValue(destinationMetadata.Description)
	d.Slug = types.StringValue(destinationMetadata.Slug)
	d.Logos = d.getLogosDestinationMetadata(destinationMetadata.Logos)
	d.Options = getOptions(destinationMetadata.Options)
	d.Actions = d.getActions(destinationMetadata.Actions)
	d.Categories = getCategories(destinationMetadata.Categories)
	d.Presets = d.getPresets(destinationMetadata.Presets)
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
}
