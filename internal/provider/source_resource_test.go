package provider

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSourceResource(t *testing.T) {
	t.Parallel()

	updated := 0
	updatedSchemaSettings := 0
	fakeServer := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("content-type", "application/json")

			payload := ""
			if req.URL.Path == "/sources" && req.Method == http.MethodPost {
				payload = `
				{
					"data": {
						"source": {
							"id": "my-source-id",
							"slug": "my-source-slug",
							"name": "My source name",
							"workspaceId": "my-workspace-id",
							"enabled": true,
							"writeKeys": ["my-write-key"],
							"metadata": {
								"id": "my-metadata-id",
								"slug": "my-metadata-slug",
								"name": "My metadata name",
								"categories": ["my-category"],
								"description": "My metadata description",
								"logos": {
									"default": "https://example.segment.com/image.png",
									"alt": "https://example.segment.com/image.png",
									"mark": "https://example.segment.com/image.png"
								},
								"options": [
									{
										"name": "sid",
										"required": true,
										"type": "string",
										"defaultValue": "",
										"description": "Your Segment SID",
										"defaultValue": "default-sid"
									}
								],
								"isCloudEventSource": false
							},
							"settings": {
								"myKey": "myValue"
							},
							"labels": [
								{
									"key": "my-label-key",
									"value": "my-label-value"
								}
							]
						}
					}
				}
			`
			} else if req.URL.Path == "/sources/my-source-id" && req.Method == http.MethodGet {
				if updated <= 1 {
					payload = `
						{
							"data": {
								"source": {
									"id": "my-source-id",
									"slug": "my-source-slug",
									"name": "My source name",
									"workspaceId": "my-workspace-id",
									"enabled": true,
									"writeKeys": ["my-write-key"],
									"metadata": {
										"id": "my-metadata-id",
										"slug": "my-metadata-slug",
										"name": "My metadata name",
										"categories": ["my-category"],
										"description": "My metadata description",
										"logos": {
											"default": "https://example.segment.com/image.png",
											"alt": "https://example.segment.com/image.png",
											"mark": "https://example.segment.com/image.png"
										},
										"options": [
											{
												"name": "sid",
												"required": true,
												"type": "string",
												"defaultValue": "",
												"description": "Your Segment SID",
												"defaultValue": "default-sid"
											}
										],
										"isCloudEventSource": false
									},
									"settings": {
										"myKey": "myValue"
									},
									"labels": [
										{
											"key": "my-label-key",
											"value": "my-label-value"
										}
									]
								}
							}
						}
					`
				} else {
					payload = `
					{
						"data": {
							"source": {
								"id": "my-source-id",
								"slug": "my-new-source-slug",
								"name": "My new source name",
								"workspaceId": "my-workspace-id",
								"enabled": false,
								"writeKeys": ["my-write-key"],
								"metadata": {
									"id": "my-metadata-id",
									"slug": "my-metadata-slug",
									"name": "My metadata name",
									"categories": ["my-category"],
									"description": "My metadata description",
									"logos": {
										"default": "https://example.segment.com/image.png",
										"alt": "https://example.segment.com/image.png",
										"mark": "https://example.segment.com/image.png"
									},
									"options": [
										{
											"name": "sid",
											"required": true,
											"type": "string",
											"defaultValue": "",
											"description": "Your Segment SID",
											"defaultValue": "default-sid"
										}
									],
									"isCloudEventSource": false
								},
								"settings": {
									"myKey": "myOtherValue"
								},
								"labels": [
									{
										"key": "my-label-key",
										"value": "my-label-value"
									}
								]
							}
						}
					}
				`
				}
			} else if req.URL.Path == "/sources/my-source-id" && req.Method == http.MethodPatch {
				updated++
				if updated <= 1 {
					payload = `
						{
							"data": {
								"source": {
									"id": "my-source-id",
									"slug": "my-source-slug",
									"name": "My source name",
									"workspaceId": "my-workspace-id",
									"enabled": true,
									"writeKeys": ["my-write-key"],
									"metadata": {
										"id": "my-metadata-id",
										"slug": "my-metadata-slug",
										"name": "My metadata name",
										"categories": ["my-category"],
										"description": "My metadata description",
										"logos": {
											"default": "https://example.segment.com/image.png",
											"alt": "https://example.segment.com/image.png",
											"mark": "https://example.segment.com/image.png"
										},
										"options": [
											{
												"name": "sid",
												"required": true,
												"type": "string",
												"defaultValue": "",
												"description": "Your Segment SID",
												"defaultValue": "default-sid"
											}
										],
										"isCloudEventSource": false
									},
									"settings": {
										"myKey": "myValue"
									},
									"labels": [
										{
											"key": "my-label-key",
											"value": "my-label-value"
										}
									]
								}
							}
						}
					`
				} else {
					payload = `
					{
						"data": {
							"source": {
								"id": "my-source-id",
								"slug": "my-new-source-slug",
								"name": "My new source name",
								"workspaceId": "my-workspace-id",
								"enabled": false,
								"writeKeys": ["my-write-key"],
								"metadata": {
									"id": "my-metadata-id",
									"slug": "my-metadata-slug",
									"name": "My metadata name",
									"categories": ["my-category"],
									"description": "My metadata description",
									"logos": {
										"default": "https://example.segment.com/image.png",
										"alt": "https://example.segment.com/image.png",
										"mark": "https://example.segment.com/image.png"
									},
									"options": [
										{
											"name": "sid",
											"required": true,
											"type": "string",
											"defaultValue": "",
											"description": "Your Segment SID",
											"defaultValue": "default-sid"
										}
									],
									"isCloudEventSource": false
								},
								"settings": {
									"myKey": "myOtherValue"
								},
								"labels": [
									{
										"key": "my-label-key",
										"value": "my-label-value"
									}
								]
							}
						}
					}
				`
				}
			} else if req.URL.Path == "/sources/my-source-id/settings" && req.Method == http.MethodGet {
				if updatedSchemaSettings <= 1 {
					payload = `
						{
							"data": {
								"sourceId": "my-source-id",
								"settings": {
									"track": {
										"allowUnplannedEvents": true,
										"allowUnplannedEventProperties": true,
										"allowEventOnViolations": true,
										"allowPropertiesOnViolations": true,
										"commonEventOnViolations": "OMIT_PROPERTIES"
									},
									"group": {
										"allowTraitsOnViolations": true,
										"allowUnplannedTraits": true,
										"commonEventOnViolations": "ALLOW"
									},
									"identify": {
										"allowTraitsOnViolations": true,
										"allowUnplannedTraits": true,
										"commonEventOnViolations": "BLOCK"
									},
									"forwardingBlockedEventsTo": "my-other-source-id",
									"forwardingViolationsTo": "my-other-source-id"
								}
							}
						}
					`
				} else {
					payload = `
						{
							"data": {
								"sourceId": "my-source-id",
								"settings": {
									"track": {
										"allowUnplannedEvents": true,
										"allowUnplannedEventProperties": true,
										"allowEventOnViolations": false,
										"allowPropertiesOnViolations": true,
										"commonEventOnViolations": "ALLOW"
									},
									"group": {
										"allowTraitsOnViolations": true,
										"allowUnplannedTraits": true,
										"commonEventOnViolations": "ALLOW"
									},
									"identify": {
										"allowTraitsOnViolations": false,
										"allowUnplannedTraits": false,
										"commonEventOnViolations": "ALLOW"
									},
									"forwardingBlockedEventsTo": "my-other-other-source-id"
								}
							}
						}
					`
				}

			} else if req.URL.Path == "/sources/my-source-id/settings" && req.Method == http.MethodPatch {
				if updatedSchemaSettings == 0 {
					payload = `
						{
							"data": {
								"sourceId": "my-source-id",
								"settings": {
									"track": {
										"allowUnplannedEvents": true,
										"allowUnplannedEventProperties": true,
										"allowEventOnViolations": true,
										"allowPropertiesOnViolations": true,
										"commonEventOnViolations": "OMIT_PROPERTIES"
									},
									"group": {
										"allowTraitsOnViolations": true,
										"allowUnplannedTraits": true,
										"commonEventOnViolations": "ALLOW"
									},
									"identify": {
										"allowTraitsOnViolations": true,
										"allowUnplannedTraits": true,
										"commonEventOnViolations": "BLOCK"
									},
									"forwardingBlockedEventsTo": "my-other-source-id",
									"forwardingViolationsTo": "my-other-source-id"
								}
							}
						}
					`
				} else {
					payload = `
						{
							"data": {
								"sourceId": "my-source-id",
								"settings": {
									"track": {
										"allowUnplannedEvents": true,
										"allowUnplannedEventProperties": true,
										"allowEventOnViolations": false,
										"allowPropertiesOnViolations": true,
										"commonEventOnViolations": "ALLOW"
									},
									"group": {
										"allowTraitsOnViolations": true,
										"allowUnplannedTraits": true,
										"commonEventOnViolations": "ALLOW"
									},
									"identify": {
										"allowTraitsOnViolations": false,
										"allowUnplannedTraits": false,
										"commonEventOnViolations": "ALLOW"
									},
									"forwardingBlockedEventsTo": "my-other-other-source-id"
								}
							}
						}
					`
				}

				updatedSchemaSettings++
			}

			_, _ = w.Write([]byte(payload))
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
					resource "segment_source" "test" {
						slug = "my-source-slug"
						name = "My source name"
						metadata = {
							id = "my-metadata-id"
						}
						enabled = true
						settings = jsonencode({
							"myKey": "myValue"
						})
						labels = [
							{
								key = "my-label-key"
								value = "my-label-value"
							}
						]
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("segment_source.test", "id", "my-source-id"),
					resource.TestCheckResourceAttr("segment_source.test", "name", "My source name"),
					resource.TestCheckResourceAttr("segment_source.test", "slug", "my-source-slug"),
					resource.TestCheckResourceAttr("segment_source.test", "workspace_id", "my-workspace-id"),
					resource.TestCheckResourceAttr("segment_source.test", "enabled", "true"),
					resource.TestCheckResourceAttr("segment_source.test", "write_keys.#", "1"),
					resource.TestCheckResourceAttr("segment_source.test", "write_keys.0", "my-write-key"),
					resource.TestCheckResourceAttr("segment_source.test", "metadata.id", "my-metadata-id"),
					resource.TestCheckResourceAttr("segment_source.test", "metadata.slug", "my-metadata-slug"),
					resource.TestCheckResourceAttr("segment_source.test", "metadata.name", "My metadata name"),
					resource.TestCheckResourceAttr("segment_source.test", "metadata.categories.#", "1"),
					resource.TestCheckResourceAttr("segment_source.test", "metadata.categories.0", "my-category"),
					resource.TestCheckResourceAttr("segment_source.test", "metadata.description", "My metadata description"),
					resource.TestCheckResourceAttr("segment_source.test", "metadata.logos.default", "https://example.segment.com/image.png"),
					resource.TestCheckResourceAttr("segment_source.test", "metadata.logos.alt", "https://example.segment.com/image.png"),
					resource.TestCheckResourceAttr("segment_source.test", "metadata.logos.mark", "https://example.segment.com/image.png"),
					resource.TestCheckResourceAttr("segment_source.test", "metadata.options.#", "1"),
					resource.TestCheckResourceAttr("segment_source.test", "metadata.options.0.name", "sid"),
					resource.TestCheckResourceAttr("segment_source.test", "metadata.options.0.required", "true"),
					resource.TestCheckResourceAttr("segment_source.test", "metadata.options.0.type", "string"),
					resource.TestCheckResourceAttr("segment_source.test", "metadata.options.0.description", "Your Segment SID"),
					resource.TestCheckResourceAttr("segment_source.test", "metadata.options.0.default_value", "\"default-sid\""),
					resource.TestCheckResourceAttr("segment_source.test", "metadata.is_cloud_event_source", "false"),
					resource.TestCheckResourceAttr("segment_source.test", "settings", "{\"myKey\":\"myValue\"}"),
					resource.TestCheckResourceAttr("segment_source.test", "labels.#", "1"),
					resource.TestCheckResourceAttr("segment_source.test", "labels.0.key", "my-label-key"),
					resource.TestCheckResourceAttr("segment_source.test", "labels.0.value", "my-label-value"),
					resource.TestCheckNoResourceAttr("segment_source.test", "schema_settings"),
				),
			},
			// ImportState testing
			{
				ResourceName: "segment_source.test",
				Config: providerConfig + `
					resource "segment_source" "test" {
						slug = "my-source-slug"
						name = "My source name"
						metadata = {
							id = "my-metadata-id"
						}
						enabled = true
						settings = jsonencode({
							"myKey": "myValue"
						})
						labels = [
							{
								key = "my-label-key"
								value = "my-label-value"
							}
						]
					}
				`,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: providerConfig + `
					resource "segment_source" "test" {
						slug = "my-new-source-slug"
						name = "My new source name"
						metadata = {
							id = "my-metadata-id"
						}
						enabled = false
						settings = jsonencode({
							"myKey": "myOtherValue"
						})
						labels = [
							{
								key = "my-label-key"
								value = "my-label-value"
							}
						]
						schema_settings = {
							forwarding_blocked_events_to = "my-other-source-id"
							forwarding_violations_to = "my-other-source-id"
							track = {
								allow_unplanned_events = true
								allow_event_on_violations = true
								allow_properties_on_violations = true
								common_event_on_violations = "OMIT_PROPERTIES"
							}
							identify = {
								allow_traits_on_violations = true
								common_event_on_violations = "BLOCK"
							}
						}
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("segment_source.test", "id", "my-source-id"),
					resource.TestCheckResourceAttr("segment_source.test", "name", "My new source name"),
					resource.TestCheckResourceAttr("segment_source.test", "slug", "my-new-source-slug"),
					resource.TestCheckResourceAttr("segment_source.test", "workspace_id", "my-workspace-id"),
					resource.TestCheckResourceAttr("segment_source.test", "enabled", "false"),
					resource.TestCheckResourceAttr("segment_source.test", "write_keys.#", "1"),
					resource.TestCheckResourceAttr("segment_source.test", "write_keys.0", "my-write-key"),
					resource.TestCheckResourceAttr("segment_source.test", "metadata.id", "my-metadata-id"),
					resource.TestCheckResourceAttr("segment_source.test", "metadata.slug", "my-metadata-slug"),
					resource.TestCheckResourceAttr("segment_source.test", "metadata.name", "My metadata name"),
					resource.TestCheckResourceAttr("segment_source.test", "metadata.categories.#", "1"),
					resource.TestCheckResourceAttr("segment_source.test", "metadata.categories.0", "my-category"),
					resource.TestCheckResourceAttr("segment_source.test", "metadata.description", "My metadata description"),
					resource.TestCheckResourceAttr("segment_source.test", "metadata.logos.default", "https://example.segment.com/image.png"),
					resource.TestCheckResourceAttr("segment_source.test", "metadata.logos.alt", "https://example.segment.com/image.png"),
					resource.TestCheckResourceAttr("segment_source.test", "metadata.logos.mark", "https://example.segment.com/image.png"),
					resource.TestCheckResourceAttr("segment_source.test", "metadata.options.#", "1"),
					resource.TestCheckResourceAttr("segment_source.test", "metadata.options.0.name", "sid"),
					resource.TestCheckResourceAttr("segment_source.test", "metadata.options.0.required", "true"),
					resource.TestCheckResourceAttr("segment_source.test", "metadata.options.0.type", "string"),
					resource.TestCheckResourceAttr("segment_source.test", "metadata.options.0.description", "Your Segment SID"),
					resource.TestCheckResourceAttr("segment_source.test", "metadata.options.0.default_value", "\"default-sid\""),
					resource.TestCheckResourceAttr("segment_source.test", "metadata.is_cloud_event_source", "false"),
					resource.TestCheckResourceAttr("segment_source.test", "settings", "{\"myKey\":\"myOtherValue\"}"),
					resource.TestCheckResourceAttr("segment_source.test", "labels.#", "1"),
					resource.TestCheckResourceAttr("segment_source.test", "labels.0.key", "my-label-key"),
					resource.TestCheckResourceAttr("segment_source.test", "labels.0.value", "my-label-value"),
					resource.TestCheckResourceAttr("segment_source.test", "schema_settings.forwarding_blocked_events_to", "my-other-source-id"),
					resource.TestCheckResourceAttr("segment_source.test", "schema_settings.forwarding_violations_to", "my-other-source-id"),
					resource.TestCheckResourceAttr("segment_source.test", "schema_settings.track.allow_unplanned_events", "true"),
					resource.TestCheckResourceAttr("segment_source.test", "schema_settings.track.allow_event_on_violations", "true"),
					resource.TestCheckResourceAttr("segment_source.test", "schema_settings.track.allow_properties_on_violations", "true"),
					resource.TestCheckResourceAttr("segment_source.test", "schema_settings.track.common_event_on_violations", "OMIT_PROPERTIES"),
					resource.TestCheckNoResourceAttr("segment_source.test", "schema_settings.track.allow_unplanned_event_properties"),
					resource.TestCheckResourceAttr("segment_source.test", "schema_settings.identify.common_event_on_violations", "BLOCK"),
					resource.TestCheckResourceAttr("segment_source.test", "schema_settings.identify.allow_traits_on_violations", "true"),
					resource.TestCheckNoResourceAttr("segment_source.test", "schema_settings.identify.allow_unplanned_traits"),
					resource.TestCheckNoResourceAttr("segment_source.test", "schema_settings.group"),
				),
			},
			{
				Config: providerConfig + `
					resource "segment_source" "test" {
						slug = "my-new-source-slug"
						name = "My new source name"
						metadata = {
							id = "my-metadata-id"
						}
						enabled = false
						settings = jsonencode({
							"myKey": "myOtherValue"
						})
						labels = [
							{
								key = "my-label-key"
								value = "my-label-value"
							}
						]
						schema_settings = {
							forwarding_blocked_events_to = "my-other-other-source-id"
							track = {
								allow_unplanned_events = true
								allow_unplanned_event_properties = true
								allow_event_on_violations = false
								allow_properties_on_violations = true
								common_event_on_violations = "ALLOW"
							}
							identify = {
								allow_traits_on_violations = false
								allow_unplanned_traits = false
								common_event_on_violations = "ALLOW"
							}
							group = {
								allow_traits_on_violations = true
								allow_unplanned_traits = true
								common_event_on_violations = "ALLOW"
							}
						}
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("segment_source.test", "id", "my-source-id"),
					resource.TestCheckResourceAttr("segment_source.test", "name", "My new source name"),
					resource.TestCheckResourceAttr("segment_source.test", "slug", "my-new-source-slug"),
					resource.TestCheckResourceAttr("segment_source.test", "workspace_id", "my-workspace-id"),
					resource.TestCheckResourceAttr("segment_source.test", "enabled", "false"),
					resource.TestCheckResourceAttr("segment_source.test", "write_keys.#", "1"),
					resource.TestCheckResourceAttr("segment_source.test", "write_keys.0", "my-write-key"),
					resource.TestCheckResourceAttr("segment_source.test", "metadata.id", "my-metadata-id"),
					resource.TestCheckResourceAttr("segment_source.test", "metadata.slug", "my-metadata-slug"),
					resource.TestCheckResourceAttr("segment_source.test", "metadata.name", "My metadata name"),
					resource.TestCheckResourceAttr("segment_source.test", "metadata.categories.#", "1"),
					resource.TestCheckResourceAttr("segment_source.test", "metadata.categories.0", "my-category"),
					resource.TestCheckResourceAttr("segment_source.test", "metadata.description", "My metadata description"),
					resource.TestCheckResourceAttr("segment_source.test", "metadata.logos.default", "https://example.segment.com/image.png"),
					resource.TestCheckResourceAttr("segment_source.test", "metadata.logos.alt", "https://example.segment.com/image.png"),
					resource.TestCheckResourceAttr("segment_source.test", "metadata.logos.mark", "https://example.segment.com/image.png"),
					resource.TestCheckResourceAttr("segment_source.test", "metadata.options.#", "1"),
					resource.TestCheckResourceAttr("segment_source.test", "metadata.options.0.name", "sid"),
					resource.TestCheckResourceAttr("segment_source.test", "metadata.options.0.required", "true"),
					resource.TestCheckResourceAttr("segment_source.test", "metadata.options.0.type", "string"),
					resource.TestCheckResourceAttr("segment_source.test", "metadata.options.0.description", "Your Segment SID"),
					resource.TestCheckResourceAttr("segment_source.test", "metadata.options.0.default_value", "\"default-sid\""),
					resource.TestCheckResourceAttr("segment_source.test", "metadata.is_cloud_event_source", "false"),
					resource.TestCheckResourceAttr("segment_source.test", "settings", "{\"myKey\":\"myOtherValue\"}"),
					resource.TestCheckResourceAttr("segment_source.test", "labels.#", "1"),
					resource.TestCheckResourceAttr("segment_source.test", "labels.0.key", "my-label-key"),
					resource.TestCheckResourceAttr("segment_source.test", "labels.0.value", "my-label-value"),
					resource.TestCheckResourceAttr("segment_source.test", "schema_settings.forwarding_blocked_events_to", "my-other-other-source-id"),
					resource.TestCheckNoResourceAttr("segment_source.test", "schema_settings.forwarding_violations_to"),
					resource.TestCheckResourceAttr("segment_source.test", "schema_settings.track.allow_unplanned_events", "true"),
					resource.TestCheckResourceAttr("segment_source.test", "schema_settings.track.allow_event_on_violations", "false"),
					resource.TestCheckResourceAttr("segment_source.test", "schema_settings.track.allow_properties_on_violations", "true"),
					resource.TestCheckResourceAttr("segment_source.test", "schema_settings.track.common_event_on_violations", "ALLOW"),
					resource.TestCheckResourceAttr("segment_source.test", "schema_settings.track.allow_unplanned_event_properties", "true"),
					resource.TestCheckResourceAttr("segment_source.test", "schema_settings.identify.common_event_on_violations", "ALLOW"),
					resource.TestCheckResourceAttr("segment_source.test", "schema_settings.identify.allow_traits_on_violations", "false"),
					resource.TestCheckResourceAttr("segment_source.test", "schema_settings.identify.allow_unplanned_traits", "false"),
					resource.TestCheckResourceAttr("segment_source.test", "schema_settings.group.allow_traits_on_violations", "true"),
					resource.TestCheckResourceAttr("segment_source.test", "schema_settings.group.allow_unplanned_traits", "true"),
					resource.TestCheckResourceAttr("segment_source.test", "schema_settings.group.common_event_on_violations", "ALLOW"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
