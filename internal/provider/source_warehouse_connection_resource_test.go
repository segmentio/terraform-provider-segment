package provider

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSourceWarehouseConnectionResource(t *testing.T) {
	t.Parallel()
	created := 0

	fakeServer := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("content-type", "application/json")

			payload := `{
				"data": {
					"warehouses": [
						{
							"id": "my-warehouse-id",
							"workspaceId": "my-workspace-id",
							"enabled": true,
							"metadata": {
								"id": "my-metadata-id",
								"slug": "redshift",
								"name": "Redshift",
								"description": "Powered by Amazon Web Services",
								"logos": {
									"default": "https://d3hotuclm6if1r.cloudfront.net/logos/redshift-default.svg",
									"mark": "",
									"alt": ""
								},
								"options": []
							},
							"settings": {}
						}
					],
					"pagination": {
						"current": "MA==",
						"totalEntries": 1
					}
				}
			}`

			if req.Method == http.MethodGet && created >= 2 {
				payload = `{
					"data": {
						"warehouses": [
							{
								"id": "my-new-warehouse-id",
								"workspaceId": "my-workspace-id",
								"enabled": true,
								"metadata": {
									"id": "my-metadata-id",
									"slug": "redshift",
									"name": "Redshift",
									"description": "Powered by Amazon Web Services",
									"logos": {
										"default": "https://d3hotuclm6if1r.cloudfront.net/logos/redshift-default.svg",
										"mark": "",
										"alt": ""
									},
									"options": []
								},
								"settings": {}
							}
						],
						"pagination": {
							"current": "MA==",
							"totalEntries": 1
						}
					}
				}`
			}

			if req.Method == http.MethodPost {
				created++
				payload = `{
					"data": {
						"status": "CONNECTED"
					}
				}`
			}

			if req.Method == http.MethodDelete {
				payload = `{
					"data": {
						"status": "SUCCESS"
					}
				}`
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
						resource "segment_source_warehouse_connection" "test" {
							source_id = "my-source-id"
							warehouse_id = "my-warehouse-id"
						}
					`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("segment_source_warehouse_connection.test", "source_id", "my-source-id"),
					resource.TestCheckResourceAttr("segment_source_warehouse_connection.test", "warehouse_id", "my-warehouse-id"),
				),
			},
			// Update and Read testing
			{
				Config: providerConfig + `
						resource "segment_source_warehouse_connection" "test" {
							source_id = "my-new-source-id"
							warehouse_id = "my-new-warehouse-id"
						}
					`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("segment_source_warehouse_connection.test", "source_id", "my-new-source-id"),
					resource.TestCheckResourceAttr("segment_source_warehouse_connection.test", "warehouse_id", "my-new-warehouse-id"),
				),
			},
		},
	})
}
