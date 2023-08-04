package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/segmentio/public-api-sdk-go/api"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &destinationMetadataDataSource{}
	_ datasource.DataSourceWithConfigure = &destinationMetadataDataSource{}
)

// NewDestinationMetadataDataSource is a helper function to simplify the provider implementation.
func NewDestinationMetadataDataSource() datasource.DataSource {
	return &destinationMetadataDataSource{}
}

// destinationMetadataDataSource is the data source implementation.
type destinationMetadataDataSource struct {
	client      *api.APIClient
	authContext context.Context
}

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

type destinationMetadataDataSourceModel struct {
	Id                 types.String        `tfsdk:"id"`
	Name               types.String        `tfsdk:"name"`
	Slug               types.String        `tfsdk:"slug"`
	Description        types.String        `tfsdk:"description"`
	Logos              *Logos              `tfsdk:"logos"`
	Options            []IntegrationOption `tfsdk:"options"`
	Categories         []types.String      `tfsdk:"categories"`
	Website            types.String        `tfsdk:"website"`
	Components         []Component         `tfsdk:"components"`
	PreviousNames      []types.String      `tfsdk:"previous_names"`
	Status             types.String        `tfsdk:"status"`
	SupportedFeatures  *SupportedFeature   `tfsdk:"supported_features"`
	SupportedMethods   *SupportedMethod    `tfsdk:"supported_methods"`
	SupportedPlatforms *SupportedPlatform  `tfsdk:"supported_platforms"`
	Actions            []Action            `tfsdk:"actions"`
	Presets            []Preset            `tfsdk:"presets"`
	Contacts           []Contact           `tfsdk:"contacts"`
	PartnerOwned       types.Bool          `tfsdk:"partner_owned"`
	SupportedRegions   []types.String      `tfsdk:"supported_regions"`
	RegionEndpoints    []types.String      `tfsdk:"region_endpoints"`
}

func destinationMetadataSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The id of the Destination metadata. Config API note: analogous to `name`.",
			Computed:    true,
		},
		"name": schema.StringAttribute{
			Description: "The user-friendly name of the Destination. Config API note: equal to `displayName`.",
			Computed:    true,
		},
		"description": schema.StringAttribute{
			Description: "The description of the Destination.",
			Computed:    true,
		},
		"slug": schema.StringAttribute{
			Description: "The slug used to identify the Destination in the Segment app.",
			Computed:    true,
		},
		"logos": schema.SingleNestedAttribute{
			Description: "The Destination's logos.",
			Computed:    true,
			Attributes: map[string]schema.Attribute{
				"default": schema.StringAttribute{
					Required: true,
				},
				"mark": schema.StringAttribute{
					Description: "The logo mark.",
					Optional:    true,
				},
				"alt": schema.StringAttribute{
					Description: "The alternative text for this logo.",
					Optional:    true,
				},
			},
		},
		"options": schema.ListNestedAttribute{
			Description: "Options configured for the Destination.",
			Computed:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Description: "The name identifying this option in the context of a Segment Integration.",
						Computed:    true,
					},
					"type": schema.StringAttribute{
						Description: "Defines the type for this option in the schema.",
						Computed:    true,
					},
					"required": schema.BoolAttribute{
						Description: "Whether this is a required option when setting up the Integration.",
						Computed:    true,
					},
					"description": schema.StringAttribute{
						Description: "An optional short text description of the field.",
						Optional:    true,
					},
					//TODO: There is no equivalent of schema.AnyAttribute, therefore this field is ignored.
					//"default_value": schema.AnyAttribute{
					//	Description: "An optional default value for the field.",
					//	Optional:    true,
					//},
					"label": schema.StringAttribute{
						Description: "An optional label for this field.",
						Optional:    true,
					},
				},
			},
		},
		"status": schema.StringAttribute{
			Description: "Support status of the Destination.",
			Computed:    true,
		},
		"previous_names": schema.ListAttribute{
			ElementType: types.StringType,
			Description: "A list of names previously used by the Destination.",
			Computed:    true,
		},
		"categories": schema.ListAttribute{
			ElementType: types.StringType,
			Description: "A list of categories with which the Destination is associated.",
			Computed:    true,
		},
		"website": schema.StringAttribute{
			Description: "A website URL for this Destination.",
			Computed:    true,
		},
		"components": schema.ListNestedAttribute{
			Description: "A list of components this Destination provides.",
			Computed:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Description: "The component type.",
						Computed:    true,
					},
					"code": schema.StringAttribute{
						Description: "Link to the repository hosting the code for this component.",
						Computed:    true,
					},
					"owner": schema.StringAttribute{
						Description: "The owner of this component. Either 'SEGMENT' or 'PARTNER'.",
						Optional:    true,
					},
				},
			},
		},
		"supported_features": schema.SingleNestedAttribute{
			Description: "Features that this Destination supports.",
			Computed:    true,
			Attributes: map[string]schema.Attribute{
				"cloud_mode_instances": schema.StringAttribute{
					Description: "This Destination's support level for cloud mode instances.",
					Optional:    true,
				},
				"device_mode_instances": schema.StringAttribute{
					Description: "This Destination's support level for device mode instances.",
					Optional:    true,
				},
				"replay": schema.BoolAttribute{
					Description: "Whether this Destination supports replays.",
					Optional:    true,
				},
				"browser_unbundling": schema.BoolAttribute{
					Description: "Whether this Destination supports browser unbundling.",
					Optional:    true,
				},
				"browser_unbundling_public": schema.BoolAttribute{
					Description: "Whether this Destination supports public browser unbundling.",
					Optional:    true,
				},
			},
		},
		"supported_methods": schema.SingleNestedAttribute{
			Description: "Methods that this Destination supports.",
			Computed:    true,
			Attributes: map[string]schema.Attribute{
				"pageview": schema.BoolAttribute{
					Description: "Identifies if the Destination supports the `pageview` method.",
					Optional:    true,
				},
				"identify": schema.BoolAttribute{
					Description: "Identifies if the Destination supports the `identify` method.",
					Optional:    true,
				},
				"alias": schema.BoolAttribute{
					Description: "Identifies if the Destination supports the `alias` method.",
					Optional:    true,
				},
				"track": schema.BoolAttribute{
					Description: "Identifies if the Destination supports the `track` method.",
					Optional:    true,
				},
				"group": schema.BoolAttribute{
					Description: "Identifies if the Destination supports the `group` method.",
					Optional:    true,
				},
			},
		},
		"supported_platforms": schema.SingleNestedAttribute{
			Description: "Platforms from which the Destination receives events.",
			Computed:    true,
			Attributes: map[string]schema.Attribute{
				"browser": schema.BoolAttribute{
					Description: "Whether this Destination supports browser events.",
					Optional:    true,
				},
				"server": schema.BoolAttribute{
					Description: "Whether this Destination supports server events.",
					Optional:    true,
				},
				"mobile": schema.BoolAttribute{
					Description: "Whether this Destination supports mobile events.",
					Optional:    true,
				},
			},
		},
		"actions": schema.ListNestedAttribute{
			Description: "Actions available for the Destination.",
			Computed:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Description: "The primary key of the action.",
						Computed:    true,
					},
					"slug": schema.StringAttribute{
						Description: "A machine-readable key unique to the action definition.",
						Computed:    true,
					},
					"name": schema.StringAttribute{
						Description: "A human-readable name for the action.",
						Computed:    true,
					},
					"description": schema.StringAttribute{
						Description: "A human-readable description of the action. May include Markdown.",
						Computed:    true,
					},
					"platform": schema.StringAttribute{
						Description: "The platform on which this action runs.",
						Computed:    true,
					},
					"hidden": schema.BoolAttribute{
						Description: "Whether the action should be hidden.",
						Computed:    true,
					},
					"default_trigger": schema.StringAttribute{
						Description: "The default value used as the trigger when connecting this action.",
						Optional:    true,
					},
					"fields": schema.ListNestedAttribute{
						Description: "The fields expected in order to perform the action.",
						Computed:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									Description: "The primary key of the field.",
									Computed:    true,
								},
								"sort_order": schema.Float64Attribute{
									Description: "The order this particular field is (used in the UI for displaying the fields in a specified order).",
									Computed:    true,
								},
								"field_key": schema.StringAttribute{
									Description: "A unique machine-readable key for the field. Should ideally match the expected key in the action's API request.",
									Computed:    true,
								},
								"label": schema.StringAttribute{
									Description: "A human-readable label for this value.",
									Computed:    true,
								},
								"type": schema.StringAttribute{
									Description: "The data type for this value.",
									Computed:    true,
								},
								"description": schema.StringAttribute{
									Description: "A human-readable description of this value. You can use Markdown.",
									Computed:    true,
								},
								"placeholder": schema.StringAttribute{
									Description: "An example value displayed but not saved.",
									Optional:    true,
								},
								//TODO: There is no equivalent of schema.AnyAttribute, therefore this field is ignored.
								//"default_value": {
								//	Type:        schema.TypeAny,
								//	Description: "A default value that is saved the first time an action is created.",
								//	Optional:    true,
								//}
								"required": schema.BoolAttribute{
									Description: "Whether this field is required.",
									Computed:    true,
								},
								"multiple": schema.BoolAttribute{
									Description: "Whether a user can provide multiples of this field.",
									Computed:    true,
								},
								//TODO: This Map field has dynamic values and since there is no equivalent of type Any, this field is excluded.
								//"choices": schema.MapAttribute{
								//	ElementType: types.MapType{},
								//	Description: "A list of machine-readable value/label pairs to populate a static dropdown.",
								//	Optional:    true,
								//},
								"dynamic": schema.BoolAttribute{
									Description: "Whether this field should execute a dynamic request to fetch choices to populate a dropdown. When true, `choices` is ignored.",
									Computed:    true,
								},
								"allow_null": schema.BoolAttribute{
									Description: "Whether this field allows null values.",
									Computed:    true,
								},
							},
						},
					},
				}},
		},
		"presets": schema.ListNestedAttribute{
			Description: "Predefined Destination subscriptions that can optionally be applied when connecting a new instance of the Destination.",
			Computed:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"action_id": schema.StringAttribute{
						Description: "The unique identifier for the Destination Action to trigger.",
						Computed:    true,
					},
					"name": schema.StringAttribute{
						Description: "The name of the subscription.",
						Computed:    true,
					},
					//TODO: This Map field has dynamic values and since there is no equivalent of type Any, this field is excluded.
					//"fields": schema.MapAttribute{
					//	ElementType: types.MapType{},
					//	Computed:    true,
					//	Description: "The default settings for action fields.",
					//},
					"trigger": schema.StringAttribute{
						Description: "FQL string that describes what events should trigger an action. See https://segment.com/docs/config-api/fql/ for more information regarding Segment's Filter Query Language (FQL).",
						Computed:    true,
					},
				}},
		},
		"contacts": schema.ListNestedAttribute{
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Description: "Name of this contact.",
						Computed:    true,
					},
					"email": schema.StringAttribute{
						Description: "Email of this contact.",
						Computed:    true,
					},
					"role": schema.StringAttribute{
						Description: "Role of this contact.",
						Computed:    true,
					},
					"is_primary": schema.BoolAttribute{
						Description: "Whether this is a primary contact.",
						Computed:    true,
					},
				},
			},
			Description: "Contact info for Integration Owners.",
			Computed:    true,
		},
		"partner_owned": schema.BoolAttribute{
			Description: "Partner Owned flag.",
			Computed:    true,
		},
		"supported_regions": schema.ListAttribute{
			ElementType: types.StringType,
			Description: "A list of supported regions for this Destination.",
			Computed:    true,
		},
		"region_endpoints": schema.ListAttribute{
			ElementType: types.StringType,
			Description: "The list of regional endpoints for this Destination.",
			Computed:    true,
		},
	}
}

// Metadata returns the data source type name.
func (d *destinationMetadataDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_destination_metadata"
}

// Schema defines the schema for the data source.
func (d *destinationMetadataDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The destination metadata",
		Attributes:  destinationMetadataSchema(),
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *destinationMetadataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state destinationMetadataDataSourceModel

	diags := req.Config.Get(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, _, err := d.client.CatalogApi.GetDestinationMetadata(d.authContext, state.Id.ValueString()).Execute()

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Source metadata",
			err.Error(),
		)
		return
	}

	var destinationMetadata = response.Data.DestinationMetadata

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

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func getPartnerOwned(owned *bool) types.Bool {
	var partnerOwned types.Bool

	if owned != nil {
		partnerOwned = types.BoolValue(*owned)
	}

	return partnerOwned
}

func getSupportedPlatforms(platforms api.SupportedPlatforms) *SupportedPlatform {
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

func getSupportedFeatures(features api.SupportedFeatures) *SupportedFeature {
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

func getSupportedMethods(methods api.SupportedMethods) *SupportedMethod {
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

func getPreviousNames(names []string) []types.String {
	var previousNames []types.String

	for _, name := range names {
		previousNames = append(previousNames, types.StringValue(name))
	}

	return previousNames
}

func getComponents(components []api.DestinationMetadataComponentV1) []Component {
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

func getRegionEndpoints(endpoints []string) []types.String {
	var regionEndpoints []types.String

	for _, endpoint := range endpoints {
		regionEndpoints = append(regionEndpoints, types.StringValue(endpoint))
	}

	return regionEndpoints
}

func getSupportedRegions(regions []string) []types.String {
	var supportedRegionsToAdd []types.String

	for _, region := range regions {
		supportedRegionsToAdd = append(supportedRegionsToAdd, types.StringValue(region))
	}

	return supportedRegionsToAdd
}

func getContacts(contacts []api.Contact) []Contact {
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

func getPresets(presets []api.DestinationMetadataSubscriptionPresetV1) []Preset {
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

func getActions(actions []api.DestinationMetadataActionV1) []Action {
	var destinationActionsToAdd []Action

	for _, action := range actions {
		destinationMetadataAction := Action{
			Id:          types.StringValue(action.Id),
			Slug:        types.StringValue(action.Slug),
			Name:        types.StringValue(action.Name),
			Description: types.StringValue(action.Description),
			Platform:    types.StringValue(action.Platform),
			Hidden:      types.BoolValue(action.Hidden),
			Fields:      getFields(action.Fields),
		}

		if action.DefaultTrigger.IsSet() {
			destinationMetadataAction.DefaultTrigger = types.StringValue(*action.DefaultTrigger.Get())
		}

		destinationActionsToAdd = append(destinationActionsToAdd, destinationMetadataAction)
	}

	return destinationActionsToAdd
}

func getFields(fields []api.DestinationMetadataActionFieldV1) []Field {
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

// Configure adds the provider configured client to the data source.
func (d *destinationMetadataDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	clientInfo, ok := req.ProviderData.(*ClientInfo)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *ClientInfo, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = clientInfo.client
	d.authContext = clientInfo.authContext
}

func getLogosDestinationMetadata(logos api.Logos) *Logos {
	logosToAdd := Logos{
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
