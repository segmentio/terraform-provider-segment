package provider

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSourceMetadataDataSource(t *testing.T) {
	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		fakeServer := httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("content-type", "application/json")
				_, _ = w.Write([]byte(`
					{
						"data": {
							"sourceMetadata": {
								"id": "my-source-metadata-id",
								"slug": "my-source-metadata-slug",
								"name": "The name of the source metadata",
								"categories": [
									"Payments"
								],
								"description": "A description of a source metadata.",
								"logos": {
									"default": "default logo",
									"alt": "alt logo",
									"mark": "mark logo"
								},
								"options": [
									{
										"name": "the option name",
										"required": true,
										"type": "the option type",
										"description": "the option description",
										"label": "the option label",
										"defaultValue": "default"
									}
								],
								"isCloudEventSource": false
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
					Config: providerConfig + `data "segment_source_metadata" "test" {}`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.segment_source_metadata.test", "id", "my-source-metadata-id"),
						resource.TestCheckResourceAttr("data.segment_source_metadata.test", "name", "The name of the source metadata"),
						resource.TestCheckResourceAttr("data.segment_source_metadata.test", "slug", "my-source-metadata-slug"),
						resource.TestCheckResourceAttr("data.segment_source_metadata.test", "description", "A description of a source metadata."),
						resource.TestCheckResourceAttr("data.segment_source_metadata.test", "is_cloud_event_source", "false"),
						resource.TestCheckResourceAttr("data.segment_source_metadata.test", "logos.default", "default logo"),
						resource.TestCheckResourceAttr("data.segment_source_metadata.test", "logos.mark", "mark logo"),
						resource.TestCheckResourceAttr("data.segment_source_metadata.test", "logos.alt", "alt logo"),
						resource.TestCheckResourceAttr("data.segment_source_metadata.test", "categories.#", "1"),
						resource.TestCheckResourceAttr("data.segment_source_metadata.test", "categories.0", "Payments"),
						resource.TestCheckResourceAttr("data.segment_source_metadata.test", "options.#", "1"),
						resource.TestCheckResourceAttr("data.segment_source_metadata.test", "options.0.name", "the option name"),
						resource.TestCheckResourceAttr("data.segment_source_metadata.test", "options.0.type", "the option type"),
						resource.TestCheckResourceAttr("data.segment_source_metadata.test", "options.0.required", "true"),
						resource.TestCheckResourceAttr("data.segment_source_metadata.test", "options.0.description", "the option description"),
						resource.TestCheckResourceAttr("data.segment_source_metadata.test", "options.0.label", "the option label"),
						resource.TestCheckResourceAttr("data.segment_source_metadata.test", "options.0.default_value", "\"default\""),
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
							"sourceMetadata": {
								"id": "my-source-metadata-id",
								"slug": "my-source-metadata-slug",
								"name": "The name of the source metadata",
								"categories": [
									"Payments"
								],
								"description": "A description of a source metadata.",
								"logos": {
									"default": "default logo"
								},
								"options": [
									{
										"name": "the option name",
										"required": true,
										"type": "the option type"
									}
								],
								"isCloudEventSource": false
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
					Config: providerConfig + `data "segment_source_metadata" "test" {}`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.segment_source_metadata.test", "id", "my-source-metadata-id"),
						resource.TestCheckResourceAttr("data.segment_source_metadata.test", "name", "The name of the source metadata"),
						resource.TestCheckResourceAttr("data.segment_source_metadata.test", "slug", "my-source-metadata-slug"),
						resource.TestCheckResourceAttr("data.segment_source_metadata.test", "description", "A description of a source metadata."),
						resource.TestCheckResourceAttr("data.segment_source_metadata.test", "is_cloud_event_source", "false"),
						resource.TestCheckResourceAttr("data.segment_source_metadata.test", "logos.default", "default logo"),
						resource.TestCheckNoResourceAttr("data.segment_source_metadata.test", "logos.alt"),
						resource.TestCheckNoResourceAttr("data.segment_source_metadata.test", "logos.mark"),
						resource.TestCheckResourceAttr("data.segment_source_metadata.test", "categories.#", "1"),
						resource.TestCheckResourceAttr("data.segment_source_metadata.test", "categories.0", "Payments"),
						resource.TestCheckResourceAttr("data.segment_source_metadata.test", "options.#", "1"),
						resource.TestCheckResourceAttr("data.segment_source_metadata.test", "options.0.name", "the option name"),
						resource.TestCheckResourceAttr("data.segment_source_metadata.test", "options.0.type", "the option type"),
						resource.TestCheckResourceAttr("data.segment_source_metadata.test", "options.0.required", "true"),
						resource.TestCheckNoResourceAttr("data.segment_source_metadata.test", "options.0.description"),
						resource.TestCheckNoResourceAttr("data.segment_source_metadata.test", "options.0.label"),
						resource.TestCheckNoResourceAttr("data.segment_source_metadata.test", "options.0.default_value"),
					),
				},
			},
		})
	})
}
