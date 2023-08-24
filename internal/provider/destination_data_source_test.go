package provider

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDestinationDataSource(t *testing.T) {
	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		fakeServer := httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("content-type", "application/json")
				_, _ = w.Write([]byte(`
					{
						"data": {
							"destination": {
								"id": "destination-id",
								"enabled": true,
								"name": "my destination name",
								"sourceId": "source-id",
								"settings": {
									"myKey": "myValue"
								},
								"metadata": {
									"id": "destination-metadata-id",
									"name": "Destination Metadata",
									"description": "Description.",
									"slug": "destination-metadata",
									"logos": {
										"default": "default",
										"mark": "mark",
										"alt": "alt"
									},
									"options": [
										{
											"name": "apiKey",
											"type": "string",
											"defaultValue": "default",
											"description": "description",
											"required": true,
											"label": "API Key",
											"defaultValue": "default"
										}
									],
									"status": "PUBLIC",
									"categories": [
										"Analytics"
									],
									"website": "https://test.com",
									"components": [
										{
											"code": "https://github.com/master/integrations/integration-name",
											"owner": "OWNER",
											"type": "BROWSER"
										}
									],
									"previousNames": [
										"destination-metadata"
									],
									"supportedMethods": {
										"track": true,
										"pageview": true,
										"identify": true,
										"group": true,
										"alias": false
									},
									"supportedPlatforms": {
										"browser": true,
										"mobile": true,
										"server": true
									},
									"supportedFeatures": {
										"cloudModeInstances": "0",
										"deviceModeInstances": "0",
										"replay": false,
										"browserUnbundling": true,
										"browserUnbundlingPublic": true
									},
									"actions": [
										{
											"id": "the-id",
											"slug": "action-slug",
											"name": "action-name",
											"description": "action-description",
											"platform": "action-platform",
											"hidden": false,
											"defaultTrigger": "trigger",
											"fields": [
												{
													"id": "field-id",
													"sort_order": "1234",
													"fieldKey": "field-key",
													"label": "field-label",
													"type": "field-type",
													"description": "field-description",
													"placeholder": "field-placeholder",
													"required": false,
													"multiple": false,
													"dynamic": false,
													"allowNull": false,
													"defaultValue": "default",
													"choices": [
														"choice1"
													]
												}
											]
										}
									],
									"presets": [
										{
											"actionId": "id",
											"name": "name",
											"trigger": "trigger"
										}
									],
									"contacts": [
										{
											"name": "Contact Name",
											"email": "contact@contact.com",
											"role": "Product Manager",
											"isPrimary": true
										}
									],
									"partnerOwned": false,
									"supportedRegions": [
										"eu-west-1",
										"us-west-2"
									],
									"regionEndpoints": [
										"US",
										"EU"
									]
								}
							}
						}
					}
				`))
			}),
		)
		defer fakeServer.Close()

		providerConfig := `
	provider "segment" {
		url   = "` + fakeServer.URL + `"
		token = "abc123"
	}
	`

		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// Read testing
				{
					Config: providerConfig + `data "segment_destination" "test" { id = "destination-id" }`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.segment_destination.test", "id", "destination-id"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "enabled", "true"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "name", "my destination name"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "source_id", "source-id"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.id", "destination-metadata-id"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.name", "Destination Metadata"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.slug", "destination-metadata"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.description", "Description."),
						resource.TestCheckResourceAttr("data.segment_destination.test", "settings", "{\"myKey\":\"myValue\"}"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.logos.default", "default"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.logos.mark", "mark"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.logos.alt", "alt"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.options.#", "1"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.options.0.name", "apiKey"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.options.0.type", "string"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.options.0.required", "true"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.options.0.description", "description"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.options.0.label", "API Key"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.options.0.default_value", "\"default\""),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.status", "PUBLIC"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.categories.#", "1"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.categories.0", "Analytics"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.website", "https://test.com"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.components.#", "1"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.components.0.code", "https://github.com/master/integrations/integration-name"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.components.0.owner", "OWNER"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.components.0.type", "BROWSER"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.previous_names.#", "1"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.previous_names.0", "destination-metadata"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.supported_methods.track", "true"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.supported_methods.pageview", "true"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.supported_methods.identify", "true"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.supported_methods.group", "true"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.supported_methods.alias", "false"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.supported_platforms.browser", "true"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.supported_platforms.mobile", "true"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.supported_platforms.server", "true"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.supported_features.cloud_mode_instances", "0"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.supported_features.device_mode_instances", "0"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.supported_features.replay", "false"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.supported_features.browser_unbundling", "true"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.supported_features.browser_unbundling_public", "true"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.actions.#", "1"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.actions.0.id", "the-id"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.actions.0.slug", "action-slug"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.actions.0.name", "action-name"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.actions.0.description", "action-description"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.actions.0.platform", "action-platform"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.actions.0.hidden", "false"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.actions.0.default_trigger", "trigger"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.actions.0.fields.#", "1"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.actions.0.fields.0.id", "field-id"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.actions.0.fields.0.sort_order", "0"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.actions.0.fields.0.field_key", "field-key"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.actions.0.fields.0.label", "field-label"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.actions.0.fields.0.type", "field-type"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.actions.0.fields.0.description", "field-description"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.actions.0.fields.0.placeholder", "field-placeholder"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.actions.0.fields.0.required", "false"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.actions.0.fields.0.multiple", "false"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.actions.0.fields.0.dynamic", "false"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.actions.0.fields.0.allow_null", "false"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.actions.0.fields.0.default_value", "\"default\""),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.actions.0.fields.0.choices", "[\"choice1\"]"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.presets.#", "1"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.presets.0.action_id", "id"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.presets.0.name", "name"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.presets.0.trigger", "trigger"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.contacts.#", "1"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.contacts.0.name", "Contact Name"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.contacts.0.email", "contact@contact.com"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.contacts.0.role", "Product Manager"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.contacts.0.is_primary", "true"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.partner_owned", "false"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.supported_regions.#", "2"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.supported_regions.0", "eu-west-1"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.supported_regions.1", "us-west-2"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.region_endpoints.#", "2"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.region_endpoints.0", "US"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.region_endpoints.1", "EU"),
					),
				},
			},
		})
	})

	t.Run("nulls", func(t *testing.T) {
		fakeServer := httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("content-type", "application/json")
				_, _ = w.Write([]byte(`
				{
  "data": {
    "destination": {
      "id": "destination-id",
      "enabled": true,
      "name": "my destination name",
      "sourceId": "source-id",
      "metadata": {
        "id": "destination-metadata-id",
        "name": "Destination Metadata",
        "description": "Description.",
        "slug": "destination-metadata",
        "logos": {
          "default": "default"
        },
        "options": [
          {
            "name": "apiKey",
            "type": "string",
            "required": true
          }
        ],
        "status": "PUBLIC",
        "categories": [
          "Analytics"
        ],
        "website": "https://test.com",
        "components": [
          {
            "code": "https://github.com/master/integrations/integration-name",
            "type": "BROWSER"
          }
        ],
        "previousNames": [
          "destination-metadata"
        ],
        "supportedMethods": {},
        "supportedPlatforms": {},
        "supportedFeatures": {},
        "actions": [
          {
            "id": "the-id",
            "slug": "action-slug",
            "name": "action-name",
            "description": "action-description",
            "platform": "action-platform",
            "hidden": false,
            "defaultTrigger": "trigger",
            "fields": [
              {
                "id": "field-id",
                "sort_order": "1234",
                "fieldKey": "field-key",
                "label": "field-label",
                "type": "field-type",
                "description": "field-description",
                "required": false,
                "multiple": false,
                "dynamic": false,
                "allowNull": false
              }
            ]
          }
        ],
        "presets": [
          {
            "actionId": "id",
            "name": "name",
            "trigger": "trigger"
          }
        ]
      }
    }
  }
}
			`))
			}),
		)
		defer fakeServer.Close()

		providerConfig := `
	provider "segment" {
		url   = "` + fakeServer.URL + `"
		token = "abc123"
	}
	`

		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// Read testing
				{
					Config: providerConfig + `
						data "segment_destination" "test" {
							id = "destination-id"
						}
					`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.segment_destination.test", "id", "destination-id"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "enabled", "true"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "name", "my destination name"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "source_id", "source-id"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.id", "destination-metadata-id"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.name", "Destination Metadata"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.slug", "destination-metadata"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.description", "Description."),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.logos.default", "default"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.options.#", "1"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.options.0.name", "apiKey"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.options.0.type", "string"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.options.0.required", "true"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.status", "PUBLIC"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.categories.#", "1"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.categories.0", "Analytics"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.website", "https://test.com"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.components.#", "1"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.components.0.code", "https://github.com/master/integrations/integration-name"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.components.0.type", "BROWSER"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.previous_names.#", "1"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.previous_names.0", "destination-metadata"),
						resource.TestCheckNoResourceAttr("data.segment_destination.test", "metadata.supported_methods.track"),
						resource.TestCheckNoResourceAttr("data.segment_destination.test", "metadata.supported_methods.pageview"),
						resource.TestCheckNoResourceAttr("data.segment_destination.test", "metadata.supported_methods.identify"),
						resource.TestCheckNoResourceAttr("data.segment_destination.test", "metadata.supported_methods.group"),
						resource.TestCheckNoResourceAttr("data.segment_destination.test", "metadata.supported_methods.alias"),
						resource.TestCheckNoResourceAttr("data.segment_destination.test", "metadata.supported_platforms.browser"),
						resource.TestCheckNoResourceAttr("data.segment_destination.test", "metadata.supported_platforms.mobile"),
						resource.TestCheckNoResourceAttr("data.segment_destination.test", "metadata.supported_platforms.server"),
						resource.TestCheckNoResourceAttr("data.segment_destination.test", "metadata.supported_features.cloud_mode_instances"),
						resource.TestCheckNoResourceAttr("data.segment_destination.test", "metadata.supported_features.device_mode_instances"),
						resource.TestCheckNoResourceAttr("data.segment_destination.test", "metadata.supported_features.replay"),
						resource.TestCheckNoResourceAttr("data.segment_destination.test", "metadata.supported_features.browser_unbundling"),
						resource.TestCheckNoResourceAttr("data.segment_destination.test", "metadata.supported_features.browser_unbundling_public"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.actions.#", "1"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.actions.0.id", "the-id"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.actions.0.slug", "action-slug"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.actions.0.name", "action-name"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.actions.0.description", "action-description"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.actions.0.platform", "action-platform"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.actions.0.hidden", "false"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.actions.0.default_trigger", "trigger"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.actions.0.fields.#", "1"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.actions.0.fields.0.id", "field-id"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.actions.0.fields.0.sort_order", "0"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.actions.0.fields.0.field_key", "field-key"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.actions.0.fields.0.label", "field-label"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.actions.0.fields.0.type", "field-type"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.actions.0.fields.0.description", "field-description"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.actions.0.fields.0.required", "false"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.actions.0.fields.0.multiple", "false"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.actions.0.fields.0.dynamic", "false"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.actions.0.fields.0.allow_null", "false"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.presets.#", "1"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.presets.0.action_id", "id"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.presets.0.name", "name"),
						resource.TestCheckResourceAttr("data.segment_destination.test", "metadata.presets.0.trigger", "trigger"),
					),
				},
			},
		})
	})
}
