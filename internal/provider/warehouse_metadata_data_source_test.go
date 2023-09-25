package provider

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccWarehouseMetadataDataSource(t *testing.T) {
	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		t.Parallel()
		fakeServer := httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("content-type", "application/json")
				_, _ = w.Write([]byte(`
					{
						"data": {
							"warehouseMetadata": {
								"id": "my-warehouse-metadata-id",
								"slug": "my-warehouse-metadata-slug",
								"name": "The name of the warehouse metadata",
								"description": "The description of a warehouse metadata",
								"logos": {
									"default": "the default value of a logo",
									"mark": "the mark value of a logo",
									"alt": "the alt value of a logo"
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
					Config: providerConfig + `data "segment_warehouse_metadata" "test" { id = "my-warehouse-metadata-id" }`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.segment_warehouse_metadata.test", "id", "my-warehouse-metadata-id"),
						resource.TestCheckResourceAttr("data.segment_warehouse_metadata.test", "name", "The name of the warehouse metadata"),
						resource.TestCheckResourceAttr("data.segment_warehouse_metadata.test", "slug", "my-warehouse-metadata-slug"),
						resource.TestCheckResourceAttr("data.segment_warehouse_metadata.test", "description", "The description of a warehouse metadata"),
						resource.TestCheckResourceAttr("data.segment_warehouse_metadata.test", "logos.default", "the default value of a logo"),
						resource.TestCheckResourceAttr("data.segment_warehouse_metadata.test", "logos.mark", "the mark value of a logo"),
						resource.TestCheckResourceAttr("data.segment_warehouse_metadata.test", "logos.alt", "the alt value of a logo"),
						resource.TestCheckResourceAttr("data.segment_warehouse_metadata.test", "options.#", "1"),
						resource.TestCheckResourceAttr("data.segment_warehouse_metadata.test", "options.0.name", "the option name"),
						resource.TestCheckResourceAttr("data.segment_warehouse_metadata.test", "options.0.type", "the option type"),
						resource.TestCheckResourceAttr("data.segment_warehouse_metadata.test", "options.0.required", "true"),
						resource.TestCheckResourceAttr("data.segment_warehouse_metadata.test", "options.0.description", "the option description"),
						resource.TestCheckResourceAttr("data.segment_warehouse_metadata.test", "options.0.label", "the option label"),
						resource.TestCheckResourceAttr("data.segment_warehouse_metadata.test", "options.0.default_value", "\"default\""),
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
							"warehouseMetadata": {
								"id": "my-warehouse-metadata-id",
								"slug": "my-warehouse-metadata-slug",
								"name": "The name of the warehouse metadata",
								"description": "The description of a warehouse metadata",
								"logos": {
									"default": "the default value of a logo"
								},
								"options": [
									{
										"name": "the option name",
										"required": true,
										"type": "the option type"
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
					Config: providerConfig + `data "segment_warehouse_metadata" "test" { id = "my-warehouse-metadata-id" }`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.segment_warehouse_metadata.test", "id", "my-warehouse-metadata-id"),
						resource.TestCheckResourceAttr("data.segment_warehouse_metadata.test", "name", "The name of the warehouse metadata"),
						resource.TestCheckResourceAttr("data.segment_warehouse_metadata.test", "slug", "my-warehouse-metadata-slug"),
						resource.TestCheckResourceAttr("data.segment_warehouse_metadata.test", "description", "The description of a warehouse metadata"),
						resource.TestCheckResourceAttr("data.segment_warehouse_metadata.test", "logos.default", "the default value of a logo"),
						resource.TestCheckNoResourceAttr("data.segment_warehouse_metadata.test", "logos.alt"),
						resource.TestCheckNoResourceAttr("data.segment_warehouse_metadata.test", "logos.mark"),
						resource.TestCheckResourceAttr("data.segment_warehouse_metadata.test", "options.#", "1"),
						resource.TestCheckResourceAttr("data.segment_warehouse_metadata.test", "options.0.name", "the option name"),
						resource.TestCheckResourceAttr("data.segment_warehouse_metadata.test", "options.0.type", "the option type"),
						resource.TestCheckResourceAttr("data.segment_warehouse_metadata.test", "options.0.required", "true"),
						resource.TestCheckNoResourceAttr("data.segment_warehouse_metadata.test", "options.0.description"),
						resource.TestCheckNoResourceAttr("data.segment_warehouse_metadata.test", "options.0.label"),
						resource.TestCheckNoResourceAttr("data.segment_warehouse_metadata.test", "options.0.default_value"),
					),
				},
			},
		})
	})
}
