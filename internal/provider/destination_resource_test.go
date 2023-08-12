package provider

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDestinationResource(t *testing.T) {
	generalPayload := `{
  "data": {
    "destination": {
      "id": "destination-id",
      "enabled": true,
      "sourceId": "s-123",
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
            "label": "API Key"
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
}`

	updatePayload := `{
  "data": {
    "destination": {
      "id": "destination-id",
      "name": "destination-name",
      "enabled": false,
      "sourceId": "s-123",
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
            "label": "API Key"
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
}`

	t.Run("happy path", func(t *testing.T) {
		fakeServer := httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("content-type", "application/json")
				switch r.Method {
				case "POST":
					_, _ = w.Write([]byte(generalPayload))
				case "GET":
					_, _ = w.Write([]byte(generalPayload))
				case "PATCH":
					_, _ = w.Write([]byte(updatePayload))
				}
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
				// Create and Read testing
				{
					Config: providerConfig + `
resource "segment_destination" "test" {
	enabled = true
	source_id = "s-123"
	metadata_id = "destination-metadata-id"
}
`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("segment_destination.test", "id", "destination-id"),
						resource.TestCheckResourceAttr("segment_destination.test", "enabled", "true"),
						resource.TestCheckResourceAttr("segment_destination.test", "source_id", "s-123"),
						resource.TestCheckResourceAttr("segment_destination.test", "metadata_id", "destination-metadata-id"),
						//resource.TestCheckResourceAttr("segment_destination.test", "metadata.name", "Destination Metadata"),
						//resource.TestCheckResourceAttr("segment_destination.test", "metadata.slug", "destination-metadata"),
						//resource.TestCheckResourceAttr("segment_destination.test", "metadata.description", "Description."),
						//resource.TestCheckResourceAttr("segment_destination.test", "metadata.logos.default", "default"),
						//resource.TestCheckResourceAttr("segment_destination.test", "metadata.logos.mark", "mark"),
						//resource.TestCheckResourceAttr("segment_destination.test", "metadata.logos.alt", "alt"),
						//resource.TestCheckResourceAttr("segment_destination.test", "metadata.options.#", "1"),
						//resource.TestCheckResourceAttr("segment_destination.test", "metadata.options.0.name", "apiKey"),
						//resource.TestCheckResourceAttr("segment_destination.test", "metadata.options.0.type", "string"),
						//resource.TestCheckResourceAttr("segment_destination.test", "metadata.options.0.required", "true"),
						//resource.TestCheckResourceAttr("segment_destination.test", "metadata.options.0.description", "description"),
						//resource.TestCheckResourceAttr("segment_destination.test", "metadata.options.0.label", "API Key"),
						//resource.TestCheckResourceAttr("segment_destination.test", "metadata.status", "PUBLIC"),
						//resource.TestCheckResourceAttr("segment_destination.test", "metadata.categories.#", "1"),
						//resource.TestCheckResourceAttr("segment_destination.test", "metadata.categories.0", "Analytics"),
						//resource.TestCheckResourceAttr("segment_destination.test", "metadata.website", "https://test.com"),
						//resource.TestCheckResourceAttr("segment_destination.test", "metadata.components.#", "1"),
						//resource.TestCheckResourceAttr("segment_destination.test", "metadata.components.0.code", "https://github.com/master/integrations/integration-name"),
						//resource.TestCheckResourceAttr("segment_destination.test", "metadata.components.0.owner", "OWNER"),
						//resource.TestCheckResourceAttr("segment_destination.test", "metadata.components.0.type", "BROWSER"),
						//resource.TestCheckResourceAttr("segment_destination.test", "metadata.previous_names.#", "1"),
						//resource.TestCheckResourceAttr("segment_destination.test", "metadata.previous_names.0", "destination-metadata"),
						//resource.TestCheckResourceAttr("segment_destination.test", "metadata.supported_methods.track", "true"),
						//resource.TestCheckResourceAttr("segment_destination.test", "metadata.supported_methods.pageview", "true"),
						//resource.TestCheckResourceAttr("segment_destination.test", "metadata.supported_methods.identify", "true"),
						//resource.TestCheckResourceAttr("segment_destination.test", "metadata.supported_methods.group", "true"),
						//resource.TestCheckResourceAttr("segment_destination.test", "metadata.supported_methods.alias", "false"),
						//resource.TestCheckResourceAttr("segment_destination.test", "metadata.supported_platforms.browser", "true"),
						//resource.TestCheckResourceAttr("segment_destination.test", "metadata.supported_platforms.mobile", "true"),
						//resource.TestCheckResourceAttr("segment_destination.test", "metadata.supported_platforms.server", "true"),
						//resource.TestCheckResourceAttr("segment_destination.test", "metadata.supported_features.cloud_mode_instances", "0"),
						//resource.TestCheckResourceAttr("segment_destination.test", "metadata.supported_features.device_mode_instances", "0"),
						//resource.TestCheckResourceAttr("segment_destination.test", "metadata.supported_features.replay", "false"),
						//resource.TestCheckResourceAttr("segment_destination.test", "metadata.supported_features.browser_unbundling", "true"),
						//resource.TestCheckResourceAttr("segment_destination.test", "metadata.supported_features.browser_unbundling_public", "true"),
						//resource.TestCheckResourceAttr("segment_destination.test", "metadata.actions.#", "1"),
						//resource.TestCheckResourceAttr("segment_destination.test", "metadata.actions.0.id", "the-id"),
						//resource.TestCheckResourceAttr("segment_destination.test", "metadata.actions.0.slug", "action-slug"),
						//resource.TestCheckResourceAttr("segment_destination.test", "metadata.actions.0.name", "action-name"),
						//resource.TestCheckResourceAttr("segment_destination.test", "metadata.actions.0.description", "action-description"),
						//resource.TestCheckResourceAttr("segment_destination.test", "metadata.actions.0.platform", "action-platform"),
						//resource.TestCheckResourceAttr("segment_destination.test", "metadata.actions.0.hidden", "false"),
						//resource.TestCheckResourceAttr("segment_destination.test", "metadata.actions.0.default_trigger", "trigger"),
						//resource.TestCheckResourceAttr("segment_destination.test", "metadata.actions.0.fields.#", "1"),
						//resource.TestCheckResourceAttr("segment_destination.test", "metadata.actions.0.fields.0.id", "field-id"),
						//resource.TestCheckResourceAttr("segment_destination.test", "metadata.actions.0.fields.0.sort_order", "0"),
						//resource.TestCheckResourceAttr("segment_destination.test", "metadata.actions.0.fields.0.field_key", "field-key"),
						//resource.TestCheckResourceAttr("segment_destination.test", "metadata.actions.0.fields.0.label", "field-label"),
						//resource.TestCheckResourceAttr("segment_destination.test", "metadata.actions.0.fields.0.type", "field-type"),
						//resource.TestCheckResourceAttr("segment_destination.test", "metadata.actions.0.fields.0.description", "field-description"),
						//resource.TestCheckResourceAttr("segment_destination.test", "metadata.actions.0.fields.0.placeholder", "field-placeholder"),
						//resource.TestCheckResourceAttr("segment_destination.test", "metadata.actions.0.fields.0.required", "false"),
						//resource.TestCheckResourceAttr("segment_destination.test", "metadata.actions.0.fields.0.multiple", "false"),
						//resource.TestCheckResourceAttr("segment_destination.test", "metadata.actions.0.fields.0.dynamic", "false"),
						//resource.TestCheckResourceAttr("segment_destination.test", "metadata.actions.0.fields.0.allow_null", "false"),
						//resource.TestCheckResourceAttr("segment_destination.test", "metadata.presets.#", "1"),
						//resource.TestCheckResourceAttr("segment_destination.test", "metadata.presets.0.action_id", "id"),
						//resource.TestCheckResourceAttr("segment_destination.test", "metadata.presets.0.name", "name"),
						//resource.TestCheckResourceAttr("segment_destination.test", "metadata.presets.0.trigger", "trigger"),
						//resource.TestCheckResourceAttr("segment_destination.test", "metadata.contacts.#", "1"),
						//resource.TestCheckResourceAttr("segment_destination.test", "metadata.contacts.0.name", "Contact Name"),
						//resource.TestCheckResourceAttr("segment_destination.test", "metadata.contacts.0.email", "contact@contact.com"),
						//resource.TestCheckResourceAttr("segment_destination.test", "metadata.contacts.0.role", "Product Manager"),
						//resource.TestCheckResourceAttr("segment_destination.test", "metadata.contacts.0.is_primary", "true"),
						//resource.TestCheckResourceAttr("segment_destination.test", "metadata.partner_owned", "false"),
						//resource.TestCheckResourceAttr("segment_destination.test", "metadata.supported_regions.#", "2"),
						//resource.TestCheckResourceAttr("segment_destination.test", "metadata.supported_regions.0", "eu-west-1"),
						//resource.TestCheckResourceAttr("segment_destination.test", "metadata.supported_regions.1", "us-west-2"),
						//resource.TestCheckResourceAttr("segment_destination.test", "metadata.region_endpoints.#", "2"),
						//resource.TestCheckResourceAttr("segment_destination.test", "metadata.region_endpoints.0", "US"),
						//resource.TestCheckResourceAttr("segment_destination.test", "metadata.region_endpoints.1", "EU"),
					),
				},

				// Update and Read testing
				{
					Config: providerConfig + `
resource "segment_destination" "test" {
	enabled = false
	source_id = "s-123"
	metadata_id = "destination-metadata-id"
	name = "destination-name"
}
`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("segment_destination.test", "id", "destination-id"),
						resource.TestCheckResourceAttr("segment_destination.test", "enabled", "false"),
						resource.TestCheckResourceAttr("segment_destination.test", "name", "destination-name"),
						resource.TestCheckResourceAttr("segment_destination.test", "source_id", "s-123"),
						resource.TestCheckResourceAttr("segment_destination.test", "metadata_id", "destination-metadata-id"),
					),
				},
			},
		})
	})
}
