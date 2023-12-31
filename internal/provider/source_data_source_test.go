package provider

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSourceDataSource(t *testing.T) {
	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		t.Parallel()
		fakeServer := httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("content-type", "application/json")
				_, _ = w.Write([]byte(`
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
										"defaultValue": "default-sid",
										"description": "Your Segment SID"
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
					Config: providerConfig + `data "segment_source" "test" { id = "my-source-id" }`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.segment_source.test", "id", "my-source-id"),
						resource.TestCheckResourceAttr("data.segment_source.test", "name", "My source name"),
						resource.TestCheckResourceAttr("data.segment_source.test", "slug", "my-source-slug"),
						resource.TestCheckResourceAttr("data.segment_source.test", "workspace_id", "my-workspace-id"),
						resource.TestCheckResourceAttr("data.segment_source.test", "enabled", "true"),
						resource.TestCheckResourceAttr("data.segment_source.test", "write_keys.#", "1"),
						resource.TestCheckResourceAttr("data.segment_source.test", "write_keys.0", "my-write-key"),
						resource.TestCheckResourceAttr("data.segment_source.test", "metadata.id", "my-metadata-id"),
						resource.TestCheckResourceAttr("data.segment_source.test", "metadata.slug", "my-metadata-slug"),
						resource.TestCheckResourceAttr("data.segment_source.test", "metadata.name", "My metadata name"),
						resource.TestCheckResourceAttr("data.segment_source.test", "metadata.categories.#", "1"),
						resource.TestCheckResourceAttr("data.segment_source.test", "metadata.categories.0", "my-category"),
						resource.TestCheckResourceAttr("data.segment_source.test", "metadata.description", "My metadata description"),
						resource.TestCheckResourceAttr("data.segment_source.test", "metadata.logos.default", "https://example.segment.com/image.png"),
						resource.TestCheckResourceAttr("data.segment_source.test", "metadata.logos.alt", "https://example.segment.com/image.png"),
						resource.TestCheckResourceAttr("data.segment_source.test", "metadata.logos.mark", "https://example.segment.com/image.png"),
						resource.TestCheckResourceAttr("data.segment_source.test", "metadata.options.#", "1"),
						resource.TestCheckResourceAttr("data.segment_source.test", "metadata.options.0.name", "sid"),
						resource.TestCheckResourceAttr("data.segment_source.test", "metadata.options.0.required", "true"),
						resource.TestCheckResourceAttr("data.segment_source.test", "metadata.options.0.type", "string"),
						resource.TestCheckResourceAttr("data.segment_source.test", "metadata.options.0.description", "Your Segment SID"),
						resource.TestCheckResourceAttr("data.segment_source.test", "metadata.options.0.default_value", "\"default-sid\""),
						resource.TestCheckResourceAttr("data.segment_source.test", "metadata.is_cloud_event_source", "false"),
						resource.TestCheckResourceAttr("data.segment_source.test", "settings", "{\"myKey\":\"myValue\"}"),
						resource.TestCheckResourceAttr("data.segment_source.test", "labels.#", "1"),
						resource.TestCheckResourceAttr("data.segment_source.test", "labels.0.key", "my-label-key"),
						resource.TestCheckResourceAttr("data.segment_source.test", "labels.0.value", "my-label-value"),
					),
				},
			},
		})
	})

	t.Run("happy path with tracking plan", func(t *testing.T) {
		t.Parallel()
		fakeServer := httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				w.Header().Set("content-type", "application/json")

				payload := ""

				if req.URL.Path == "/sources/my-source-id" {
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
												"defaultValue": "default-sid",
												"description": "Your Segment SID"
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
								},
								"trackingPlanId": "my-tracking-plan-id"
							}
						}
					`
				} else if req.URL.Path == "/sources/my-source-id/settings" && req.Method == http.MethodGet {
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
				// Read testing
				{
					Config: providerConfig + `data "segment_source" "test" { id = "my-source-id" }`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.segment_source.test", "id", "my-source-id"),
						resource.TestCheckResourceAttr("data.segment_source.test", "name", "My source name"),
						resource.TestCheckResourceAttr("data.segment_source.test", "slug", "my-source-slug"),
						resource.TestCheckResourceAttr("data.segment_source.test", "workspace_id", "my-workspace-id"),
						resource.TestCheckResourceAttr("data.segment_source.test", "enabled", "true"),
						resource.TestCheckResourceAttr("data.segment_source.test", "write_keys.#", "1"),
						resource.TestCheckResourceAttr("data.segment_source.test", "write_keys.0", "my-write-key"),
						resource.TestCheckResourceAttr("data.segment_source.test", "metadata.id", "my-metadata-id"),
						resource.TestCheckResourceAttr("data.segment_source.test", "metadata.slug", "my-metadata-slug"),
						resource.TestCheckResourceAttr("data.segment_source.test", "metadata.name", "My metadata name"),
						resource.TestCheckResourceAttr("data.segment_source.test", "metadata.categories.#", "1"),
						resource.TestCheckResourceAttr("data.segment_source.test", "metadata.categories.0", "my-category"),
						resource.TestCheckResourceAttr("data.segment_source.test", "metadata.description", "My metadata description"),
						resource.TestCheckResourceAttr("data.segment_source.test", "metadata.logos.default", "https://example.segment.com/image.png"),
						resource.TestCheckResourceAttr("data.segment_source.test", "metadata.logos.alt", "https://example.segment.com/image.png"),
						resource.TestCheckResourceAttr("data.segment_source.test", "metadata.logos.mark", "https://example.segment.com/image.png"),
						resource.TestCheckResourceAttr("data.segment_source.test", "metadata.options.#", "1"),
						resource.TestCheckResourceAttr("data.segment_source.test", "metadata.options.0.name", "sid"),
						resource.TestCheckResourceAttr("data.segment_source.test", "metadata.options.0.required", "true"),
						resource.TestCheckResourceAttr("data.segment_source.test", "metadata.options.0.type", "string"),
						resource.TestCheckResourceAttr("data.segment_source.test", "metadata.options.0.description", "Your Segment SID"),
						resource.TestCheckResourceAttr("data.segment_source.test", "metadata.options.0.default_value", "\"default-sid\""),
						resource.TestCheckResourceAttr("data.segment_source.test", "metadata.is_cloud_event_source", "false"),
						resource.TestCheckResourceAttr("data.segment_source.test", "settings", "{\"myKey\":\"myValue\"}"),
						resource.TestCheckResourceAttr("data.segment_source.test", "labels.#", "1"),
						resource.TestCheckResourceAttr("data.segment_source.test", "labels.0.key", "my-label-key"),
						resource.TestCheckResourceAttr("data.segment_source.test", "labels.0.value", "my-label-value"),
						resource.TestCheckResourceAttr("data.segment_source.test", "schema_settings.forwarding_blocked_events_to", "my-other-source-id"),
						resource.TestCheckResourceAttr("data.segment_source.test", "schema_settings.forwarding_violations_to", "my-other-source-id"),
						resource.TestCheckResourceAttr("data.segment_source.test", "schema_settings.track.allow_unplanned_events", "true"),
						resource.TestCheckResourceAttr("data.segment_source.test", "schema_settings.track.allow_event_on_violations", "true"),
						resource.TestCheckResourceAttr("data.segment_source.test", "schema_settings.track.allow_properties_on_violations", "true"),
						resource.TestCheckResourceAttr("data.segment_source.test", "schema_settings.track.common_event_on_violations", "OMIT_PROPERTIES"),
						resource.TestCheckResourceAttr("data.segment_source.test", "schema_settings.track.allow_unplanned_event_properties", "true"),
						resource.TestCheckResourceAttr("data.segment_source.test", "schema_settings.identify.common_event_on_violations", "BLOCK"),
						resource.TestCheckResourceAttr("data.segment_source.test", "schema_settings.identify.allow_traits_on_violations", "true"),
						resource.TestCheckResourceAttr("data.segment_source.test", "schema_settings.identify.allow_unplanned_traits", "true"),
						resource.TestCheckResourceAttr("data.segment_source.test", "schema_settings.group.allow_traits_on_violations", "true"),
						resource.TestCheckResourceAttr("data.segment_source.test", "schema_settings.group.allow_unplanned_traits", "true"),
						resource.TestCheckResourceAttr("data.segment_source.test", "schema_settings.group.common_event_on_violations", "ALLOW"),
					),
				},
			},
		})
	})

	t.Run("nulls", func(t *testing.T) {
		t.Parallel()
		fakeServer := httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("content-type", "application/json")
				_, _ = w.Write([]byte(`
				{
					"data": {
						"source": {
							"id": "my-source-id",
							"slug": "my-source-slug",
							"workspaceId": "my-workspace-id",
							"enabled": true,
							"writeKeys": ["my-write-key"],
							"metadata": {
								"id": "my-metadata-id",
								"slug": "my-metadata-slug",
								"name": "My metadata name",
								"description": "My metadata description",
								"logos": {
									"default": "https://example.segment.com/image.png"
								},
								"options": [
									{
										"name": "sid",
										"required": true,
										"type": "string"
									}
								],
								"isCloudEventSource": false
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
					Config: providerConfig + `data "segment_source" "test" { id = "my-source-id" }`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.segment_source.test", "id", "my-source-id"),
						resource.TestCheckNoResourceAttr("data.segment_source.test", "name"),
						resource.TestCheckResourceAttr("data.segment_source.test", "slug", "my-source-slug"),
						resource.TestCheckResourceAttr("data.segment_source.test", "workspace_id", "my-workspace-id"),
						resource.TestCheckResourceAttr("data.segment_source.test", "enabled", "true"),
						resource.TestCheckResourceAttr("data.segment_source.test", "write_keys.#", "1"),
						resource.TestCheckResourceAttr("data.segment_source.test", "write_keys.0", "my-write-key"),
						resource.TestCheckResourceAttr("data.segment_source.test", "metadata.id", "my-metadata-id"),
						resource.TestCheckResourceAttr("data.segment_source.test", "metadata.slug", "my-metadata-slug"),
						resource.TestCheckResourceAttr("data.segment_source.test", "metadata.name", "My metadata name"),
						resource.TestCheckNoResourceAttr("data.segment_source.test", "metadata.categories"),
						resource.TestCheckResourceAttr("data.segment_source.test", "metadata.description", "My metadata description"),
						resource.TestCheckResourceAttr("data.segment_source.test", "metadata.logos.default", "https://example.segment.com/image.png"),
						resource.TestCheckNoResourceAttr("data.segment_source.test", "metadata.logos.alt"),
						resource.TestCheckNoResourceAttr("data.segment_source.test", "metadata.logos.mark"),
						resource.TestCheckResourceAttr("data.segment_source.test", "metadata.options.#", "1"),
						resource.TestCheckResourceAttr("data.segment_source.test", "metadata.options.0.name", "sid"),
						resource.TestCheckResourceAttr("data.segment_source.test", "metadata.options.0.required", "true"),
						resource.TestCheckResourceAttr("data.segment_source.test", "metadata.options.0.type", "string"),
						resource.TestCheckNoResourceAttr("data.segment_source.test", "metadata.options.0.default_value"),
						resource.TestCheckNoResourceAttr("data.segment_source.test", "metadata.options.0.description"),
						resource.TestCheckResourceAttr("data.segment_source.test", "metadata.is_cloud_event_source", "false"),
						resource.TestCheckNoResourceAttr("data.segment_source.test", "settings"),
						resource.TestCheckResourceAttr("data.segment_source.test", "labels.#", "1"),
						resource.TestCheckResourceAttr("data.segment_source.test", "labels.0.key", "my-label-key"),
						resource.TestCheckResourceAttr("data.segment_source.test", "labels.0.value", "my-label-value"),
					),
				},
			},
		})
	})
}
