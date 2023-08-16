package provider

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSourceResource(t *testing.T) {
	updated := 0
	fakeServer := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("content-type", "application/json")

			payload := `
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
			`

			fmt.Println(req.Method)

			// After we update the source, return the updated source for subsequent calls (first update is part of the create call)
			if req.Method == http.MethodPatch {
				updated++
			}
			if updated > 1 {
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
					resource.TestCheckResourceAttr("segment_source.test", "metadata.is_cloud_event_source", "false"),
					resource.TestCheckResourceAttr("segment_source.test", "settings.#", "0"),
					resource.TestCheckResourceAttr("segment_source.test", "labels.#", "1"),
					resource.TestCheckResourceAttr("segment_source.test", "labels.0.key", "my-label-key"),
					resource.TestCheckResourceAttr("segment_source.test", "labels.0.value", "my-label-value"),
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
					resource.TestCheckResourceAttr("segment_source.test", "metadata.is_cloud_event_source", "false"),
					resource.TestCheckResourceAttr("segment_source.test", "settings.#", "0"),
					resource.TestCheckResourceAttr("segment_source.test", "labels.#", "1"),
					resource.TestCheckResourceAttr("segment_source.test", "labels.0.key", "my-label-key"),
					resource.TestCheckResourceAttr("segment_source.test", "labels.0.value", "my-label-value"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
