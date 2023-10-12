package provider

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDestinationFilterResource(t *testing.T) {
	t.Parallel()

	updated := 0
	fakeServer := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("content-type", "application/json")

			payload := ""
			if req.URL.Path == "/destination/my-destination-id/filters" && req.Method == http.MethodPost {
				payload = `
				{
					"data": {
						"filter": {
							"id": "my-filter-id",
							"sourceId": "my-source-id",
							"destinationId": "my-destination-id",
							"if": "type = \"identify\"",
							"actions": [
								{
									"type": "SAMPLE",
									"percent": 0.2
								},
								{
									"type": "DROP_PROPERTIES",
									"fields": {
										"properties": [
											"a"
										]
									}
								}
							],
							"title": "my filter",
							"description": "my filter description",
							"enabled": true,
							"createdAt": "2023-10-10T16:34:53.000Z",
							"updatedAt": "2023-10-10T18:57:31.000Z"
						}
					}
				}`
			}

			if req.URL.Path == "/destination/my-destination-id/filters/my-filter-id" && req.Method == http.MethodPatch {
				payload = `
				{
					"data": {
						"filter": {
							"id": "my-filter-id",
							"sourceId": "my-source-id",
							"destinationId": "my-destination-id",
							"if": "type = \"track\"",
							"actions": [
								{
									"type": "SAMPLE",
									"percent": 0.3
								}
							],
							"title": "my new filter",
							"description": "my new filter description",
							"enabled": false,
							"createdAt": "2023-10-10T16:34:53.000Z",
							"updatedAt": "2023-10-10T18:57:31.000Z"
						}
					}
				}`
				updated++
			}

			if req.URL.Path == "/destination/my-destination-id/filters/my-filter-id" && req.Method == http.MethodGet {
				if updated < 1 {
					payload = `
					{
						"data": {
							"filter": {
								"id": "my-filter-id",
								"sourceId": "my-source-id",
								"destinationId": "my-destination-id",
								"if": "type = \"identify\"",
								"actions": [
									{
										"type": "SAMPLE",
										"percent": 0.2
									},
									{
										"type": "DROP_PROPERTIES",
										"fields": {
											"properties": [
												"a"
											]
										}
									}
								],
								"title": "my filter",
								"description": "my filter description",
								"enabled": true,
								"createdAt": "2023-10-10T16:34:53.000Z",
								"updatedAt": "2023-10-10T18:57:31.000Z"
							}
						}
					}`
				} else {
					payload = `
					{
						"data": {
							"filter": {
								"id": "my-filter-id",
								"sourceId": "my-source-id",
								"destinationId": "my-destination-id",
								"if": "type = \"track\"",
								"actions": [
									{
										"type": "SAMPLE",
										"percent": 0.3
									}
								],
								"title": "my new filter",
								"description": "my new filter description",
								"enabled": false,
								"createdAt": "2023-10-10T16:34:53.000Z",
								"updatedAt": "2023-10-10T18:57:31.000Z"
							}
						}
					}`
				}

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
					resource "segment_destination_filter" "test" {
						if             = "type = \"identify\""
						destination_id = "my-destination-id"
						source_id      = "my-source-id"
						title          = "my filter"
						enabled        = true
						description    = "my filter description"
						actions = [
						{
							type    = "SAMPLE"
							percent = 0.2
						},
						{
							type = "DROP_PROPERTIES"
							fields = jsonencode({
							"properties": ["a"]
							})
						}
						]
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("segment_destination_filter.test", "id", "my-filter-id"),
					resource.TestCheckResourceAttr("segment_destination_filter.test", "destination_id", "my-destination-id"),
					resource.TestCheckResourceAttr("segment_destination_filter.test", "source_id", "my-source-id"),
					resource.TestCheckResourceAttr("segment_destination_filter.test", "title", "my filter"),
					resource.TestCheckResourceAttr("segment_destination_filter.test", "enabled", "true"),
					resource.TestCheckResourceAttr("segment_destination_filter.test", "description", "my filter description"),
					resource.TestCheckResourceAttr("segment_destination_filter.test", "actions.#", "2"),
					resource.TestCheckResourceAttr("segment_destination_filter.test", "actions.1.type", "SAMPLE"),
					resource.TestCheckResourceAttr("segment_destination_filter.test", "actions.1.percent", "0.2"),
					resource.TestCheckResourceAttr("segment_destination_filter.test", "actions.0.type", "DROP_PROPERTIES"),
					resource.TestCheckResourceAttr("segment_destination_filter.test", "actions.0.fields", "{\"properties\":[\"a\"]}"),
				),
			},
			// ImportState testing
			{
				ResourceName: "segment_destination_filter.test",
				Config: providerConfig + `
					resource "segment_destination_filter" "test" {
						if             = "type = \"identify\""
						destination_id = "my-destination-id"
						source_id      = "my-source-id"
						title          = "my filter"
						enabled        = true
						description    = "my filter description"
						actions = [
						{
							type    = "SAMPLE"
							percent = 0.2
						},
						{
							type = "DROP_PROPERTIES"
							fields = jsonencode({
							"properties": ["a"]
							})
						}
						]
					}
				`,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "my-destination-id:my-filter-id",
			},
			// Update and Read testing
			{
				Config: providerConfig + `
					resource "segment_destination_filter" "test" {
						if             = "type = \"track\""
						destination_id = "my-destination-id"
						source_id      = "my-source-id"
						title          = "my new filter"
						enabled        = false
						description    = "my new filter description"
						actions = [
							{
								type    = "SAMPLE"
								percent = 0.3
							},
						]
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("segment_destination_filter.test", "id", "my-filter-id"),
					resource.TestCheckResourceAttr("segment_destination_filter.test", "destination_id", "my-destination-id"),
					resource.TestCheckResourceAttr("segment_destination_filter.test", "source_id", "my-source-id"),
					resource.TestCheckResourceAttr("segment_destination_filter.test", "title", "my new filter"),
					resource.TestCheckResourceAttr("segment_destination_filter.test", "enabled", "false"),
					resource.TestCheckResourceAttr("segment_destination_filter.test", "description", "my new filter description"),
					resource.TestCheckResourceAttr("segment_destination_filter.test", "actions.#", "1"),
					resource.TestCheckResourceAttr("segment_destination_filter.test", "actions.0.type", "SAMPLE"),
					resource.TestCheckResourceAttr("segment_destination_filter.test", "actions.0.percent", "0.3"),
				),
			},
		},
	})
}
