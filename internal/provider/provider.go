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
	api "github.com/segmentio/public-api-sdk-go/api"
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
	Url   types.String `tfsdk:"url"`
	Token types.String `tfsdk:"token"`
}

func (p *segmentProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "segment"
	resp.Version = p.version
}

func (p *segmentProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The Segment provider.",
		Attributes: map[string]schema.Attribute{
			"url": schema.StringAttribute{
				Optional:    true,
				Description: "The Public API url. If not set, the PUBLIC_API_URL environment variable will be used, or a default of 'api.segmentapis.com'.",
			},
			"token": schema.StringAttribute{
				Required:    true,
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

	if config.Url.IsUnknown() {
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

	if !config.Url.IsNull() {
		url = config.Url.ValueString()
	}

	if !config.Token.IsNull() {
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
	configuration.UserAgent = "Segment (terraform v" + p.version + ")"
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

func (p *segmentProvider) Resources(ctx context.Context) []func() resource.Resource {
	return nil
}

func (p *segmentProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewWorkspaceDataSource,
		NewSourceDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &segmentProvider{
			version: version,
		}
	}
}
