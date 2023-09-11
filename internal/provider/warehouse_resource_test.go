package provider

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccWarehouseResource(t *testing.T) {
	t.Parallel()

	updated := 0
	fakeServer := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("content-type", "application/json")

			payload := `
				{
					"data": {
						"warehouse": {
							"id": "warehouse-id",
							"workspaceId": "workspace-id",
							"enabled": true,
							"settings": {
								"myKey": "myValue",
								"name": "My warehouse name"
							},
							"metadata": {
								"id": "my-metadata-id",
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
										"label": "the option label"
									}
								]
							}
						}
					}
				}
			`

			// After we update the source, return the updated source for subsequent calls (first update is part of the create call)
			if req.Method == http.MethodPatch {
				updated++
			}
			if updated > 0 {
				payload = `
					{
						"data": {
							"warehouse": {
								"id": "warehouse-id",
								"workspaceId": "workspace-id",
								"enabled": false,
								"settings": {
									"myKey": "myNewValue",
									"name": "My new warehouse name"
								},
								"metadata": {
									"id": "my-metadata-id",
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
											"label": "the option label"
										}
									]
								}
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
					resource "segment_warehouse" "test" {
						metadata = {
							id = "my-metadata-id"
						}
						enabled = true
						settings = jsonencode({
							"myKey": "myValue"
						})
						name = "My warehouse name"
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("segment_warehouse.test", "id", "warehouse-id"),
					resource.TestCheckResourceAttr("segment_warehouse.test", "name", "My warehouse name"),
					resource.TestCheckResourceAttr("segment_warehouse.test", "workspace_id", "workspace-id"),
					resource.TestCheckResourceAttr("segment_warehouse.test", "enabled", "true"),
					resource.TestCheckResourceAttr("segment_warehouse.test", "metadata.id", "my-metadata-id"),
					resource.TestCheckResourceAttr("segment_warehouse.test", "metadata.name", "The name of the warehouse metadata"),
					resource.TestCheckResourceAttr("segment_warehouse.test", "metadata.slug", "my-warehouse-metadata-slug"),
					resource.TestCheckResourceAttr("segment_warehouse.test", "metadata.description", "The description of a warehouse metadata"),
					resource.TestCheckResourceAttr("segment_warehouse.test", "metadata.logos.default", "the default value of a logo"),
					resource.TestCheckResourceAttr("segment_warehouse.test", "metadata.logos.mark", "the mark value of a logo"),
					resource.TestCheckResourceAttr("segment_warehouse.test", "metadata.logos.alt", "the alt value of a logo"),
					resource.TestCheckResourceAttr("segment_warehouse.test", "metadata.options.#", "1"),
					resource.TestCheckResourceAttr("segment_warehouse.test", "metadata.options.0.name", "the option name"),
					resource.TestCheckResourceAttr("segment_warehouse.test", "metadata.options.0.type", "the option type"),
					resource.TestCheckResourceAttr("segment_warehouse.test", "metadata.options.0.required", "true"),
					resource.TestCheckResourceAttr("segment_warehouse.test", "metadata.options.0.description", "the option description"),
					resource.TestCheckResourceAttr("segment_warehouse.test", "metadata.options.0.label", "the option label"),
					resource.TestCheckResourceAttr("segment_warehouse.test", "settings", "{\"myKey\":\"myValue\"}"),
				),
			},
			// ImportState testing
			{
				ResourceName: "segment_warehouse.test",
				Config: providerConfig + `
					resource "segment_warehouse" "test" {
						metadata = {
							id = "my-metadata-id"
						}
						enabled = true
						settings = jsonencode({
							"myKey": "myValue"
						})
						name = "My warehouse name"
					}
				`,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: providerConfig + `
					resource "segment_warehouse" "test" {
						metadata = {
							id = "my-metadata-id"
						}
						enabled = false
						settings = jsonencode({
							"myKey": "myNewValue"
						})
						name = "My new warehouse name"
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("segment_warehouse.test", "id", "warehouse-id"),
					resource.TestCheckResourceAttr("segment_warehouse.test", "name", "My new warehouse name"),
					resource.TestCheckResourceAttr("segment_warehouse.test", "workspace_id", "workspace-id"),
					resource.TestCheckResourceAttr("segment_warehouse.test", "enabled", "false"),
					resource.TestCheckResourceAttr("segment_warehouse.test", "metadata.id", "my-metadata-id"),
					resource.TestCheckResourceAttr("segment_warehouse.test", "metadata.name", "The name of the warehouse metadata"),
					resource.TestCheckResourceAttr("segment_warehouse.test", "metadata.slug", "my-warehouse-metadata-slug"),
					resource.TestCheckResourceAttr("segment_warehouse.test", "metadata.description", "The description of a warehouse metadata"),
					resource.TestCheckResourceAttr("segment_warehouse.test", "metadata.logos.default", "the default value of a logo"),
					resource.TestCheckResourceAttr("segment_warehouse.test", "metadata.logos.mark", "the mark value of a logo"),
					resource.TestCheckResourceAttr("segment_warehouse.test", "metadata.logos.alt", "the alt value of a logo"),
					resource.TestCheckResourceAttr("segment_warehouse.test", "metadata.options.#", "1"),
					resource.TestCheckResourceAttr("segment_warehouse.test", "metadata.options.0.name", "the option name"),
					resource.TestCheckResourceAttr("segment_warehouse.test", "metadata.options.0.type", "the option type"),
					resource.TestCheckResourceAttr("segment_warehouse.test", "metadata.options.0.required", "true"),
					resource.TestCheckResourceAttr("segment_warehouse.test", "metadata.options.0.description", "the option description"),
					resource.TestCheckResourceAttr("segment_warehouse.test", "metadata.options.0.label", "the option label"),
					resource.TestCheckResourceAttr("segment_warehouse.test", "settings", "{\"myKey\":\"myNewValue\"}"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
