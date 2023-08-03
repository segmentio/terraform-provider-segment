package provider

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSourceDataSource(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
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
										"defaultValue": "",
										"description": "Your Segment SID"
									}
								],
								"isCloudEventSource": false
							},
							"settings": {},
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
						resource.TestCheckResourceAttr("data.segment_source.test", "metadata.is_cloud_event_source", "false"),
						resource.TestCheckResourceAttr("data.segment_source.test", "settings.#", "0"),
						resource.TestCheckResourceAttr("data.segment_source.test", "labels.#", "1"),
						resource.TestCheckResourceAttr("data.segment_source.test", "labels.0.key", "my-label-key"),
						resource.TestCheckResourceAttr("data.segment_source.test", "labels.0.value", "my-label-value"),
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
