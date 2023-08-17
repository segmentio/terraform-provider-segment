package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/segmentio/public-api-sdk-go/api"
)

type DestinationStateModel struct {
	Id       types.String                   `tfsdk:"id"`
	Name     types.String                   `tfsdk:"name"`
	Enabled  types.Bool                     `tfsdk:"enabled"`
	Metadata *destinationMetadataStateModel `tfsdk:"metadata"`
	SourceId types.String                   `tfsdk:"source_id"`
}

type DestinationPlanModel struct {
	Id       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	Enabled  types.Bool   `tfsdk:"enabled"`
	Metadata types.Object `tfsdk:"metadata"`
	SourceId types.String `tfsdk:"source_id"`
}

func GetDestinationMetadata(destinationMetadata api.Metadata) *destinationMetadataStateModel {
	var state destinationMetadataStateModel

	state.Id = types.StringValue(destinationMetadata.Id)
	state.Name = types.StringValue(destinationMetadata.Name)
	state.Description = types.StringValue(destinationMetadata.Description)
	state.Slug = types.StringValue(destinationMetadata.Slug)
	state.Logos = getLogosDestinationMetadata(destinationMetadata.Logos)
	state.Options = getOptions(destinationMetadata.Options)
	state.Actions = getActions(destinationMetadata.Actions)
	state.Categories = getCategories(destinationMetadata.Categories)
	state.Presets = getPresets(destinationMetadata.Presets)
	state.Contacts = getContacts(destinationMetadata.Contacts)
	state.PartnerOwned = getPartnerOwned(destinationMetadata.PartnerOwned)
	state.SupportedRegions = getSupportedRegions(destinationMetadata.SupportedRegions)
	state.RegionEndpoints = getRegionEndpoints(destinationMetadata.RegionEndpoints)
	state.Status = types.StringValue(destinationMetadata.Status)
	state.Website = types.StringValue(destinationMetadata.Website)
	state.Components = getComponents(destinationMetadata.Components)
	state.PreviousNames = getPreviousNames(destinationMetadata.PreviousNames)
	state.SupportedMethods = getSupportedMethods(destinationMetadata.SupportedMethods)
	state.SupportedFeatures = getSupportedFeatures(destinationMetadata.SupportedFeatures)
	state.SupportedPlatforms = getSupportedPlatforms(destinationMetadata.SupportedPlatforms)

	return &state
}
