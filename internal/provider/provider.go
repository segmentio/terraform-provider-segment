// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/segmentio/public-api-sdk-go/api"
)

// Ensure segmentProvider satisfies various provider interfaces.
var _ provider.Provider = &segmentProvider{}

// segmentProvider defines the provider implementation.
type segmentProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

type ClientInfo struct {
	client      *api.APIClient
	authContext context.Context
}

// segmentProviderModel describes the provider data model.
type segmentProviderModel struct {
	URL   types.String `tfsdk:"url"`
	Token types.String `tfsdk:"token"`
}

func (p *segmentProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "segment"
	resp.Version = p.version
}

func (p *segmentProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use the Segment provider to manage resources in your [Segment workspace](https://segment.com/docs/). This provider is built on top of Segment's [Public API](https://segment.com/docs/api/public-api/), so you must configure the provider with the proper Public API token before you can use it.",
		Attributes: map[string]schema.Attribute{
			"url": schema.StringAttribute{
				Optional:    true,
				Description: "The Public API url. Defaults to 'api.segmentapis.com', but can be overwritten by supplying it as an input to the provider or as a PUBLIC_API_URL environment variable.",
			},
			"token": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "The Public API token. If not set, the PUBLIC_API_TOKEN environment variable will be used.",
			},
		},
	}
}

func (p *segmentProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var config segmentProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.URL.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("url"),
			"Unknown Public API url",
			"The provider cannot create the Public API client as there is an unknown configuration value for the Public API url. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the PUBLIC_API_URL environment variable.",
		)
	}

	if config.Token.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Unknown Public API token",
			"The provider cannot create the Public API client as there is an unknown configuration value for the Public API token. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the PUBLIC_API_TOKEN environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	url := os.Getenv("PUBLIC_API_URL")
	token := os.Getenv("PUBLIC_API_TOKEN")

	if !config.URL.IsNull() {
		url = config.URL.ValueString()
	}

	if !config.Token.IsNull() && !config.Token.IsUnknown() {
		token = config.Token.ValueString()
	}

	if url == "" {
		url = "https://api.segmentapis.com"
	}

	if token == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Missing Public API Token",
			"The provider cannot create the Public API client as there is a missing or empty value for the Public API token. "+
				"Set the token value in the configuration or use the PUBLIC_API_TOKEN environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	auth := context.WithValue(context.Background(), api.ContextAccessToken, token)
	configuration := api.NewConfiguration()
	configuration.UserAgent = "Segment (terraform " + p.version + ")"
	configuration.Servers = api.ServerConfigurations{
		{
			URL: url,
		},
	}

	client := api.NewAPIClient(configuration)

	clientInfo := &ClientInfo{
		client:      client,
		authContext: auth,
	}

	resp.DataSourceData = clientInfo
	resp.ResourceData = clientInfo
}

func (p *segmentProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewLabelResource,
		NewDestinationResource,
		NewSourceResource,
		NewWarehouseResource,
		NewSourceWarehouseConnectionResource,
		NewTrackingPlanResource,
		NewUserResource,
		NewUserGroupResource,
		NewFunctionResource,
		NewDestinationFilterResource,
		NewProfilesWarehouseResource,
		NewDestinationSubscriptionResource,
		NewSourceTrackingPlanConnectionResource,
		NewReverseETLModelResource,
	}
}

func (p *segmentProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewWorkspaceDataSource,
		NewSourceDataSource,
		NewSourceMetadataDataSource,
		NewDestinationMetadataDataSource,
		NewWarehouseMetadataDataSource,
		NewDestinationDataSource,
		NewWarehouseDataSource,
		NewTrackingPlanDataSource,
		NewRoleDataSource,
		NewUserDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &segmentProvider{
			version: version,
		}
	}
}
