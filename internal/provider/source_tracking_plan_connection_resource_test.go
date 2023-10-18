package provider

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSourceTrackingPlanConnectionResource(t *testing.T) {
	t.Parallel()
	created := 0

	fakeServer := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("content-type", "application/json")

			payload := `{
				"data": {
					"sources": [
						{
							"id": "my-source-id",
							"slug": "my-source-slug",
							"name": "My source",
							"workspaceId": "my-workspace-id",
							"enabled": false,
							"writeKeys": [
								"my-write-key"
							],
							"metadata": {
								"id": "my-metadata-id",
								"slug": "javascript",
								"name": "Javascript",
								"categories": [
									"Website"
								],
								"description": "This is our most flexible and powerful tracking system, using analytics.js.  Track and analyze information about your visitors and customers, and every action that they take, in any of our 140 integrations, business intelligence tools, or directly with SQL tools.",
								"logos": {
									"default": "https://cdn.filepicker.io/api/file/aRgo4XJQZausZxD4gZQq",
									"alt": "https://cdn.filepicker.io/api/file/aRgo4XJQZausZxD4gZQq",
									"mark": "https://cdn.filepicker.io/api/file/kBpmEoSSaakidAvoFmzd"
								},
								"options": [],
								"isCloudEventSource": false
							},
							"settings": {
								"beep": "beep"
							},
							"labels": [
								{
									"key": "product",
									"value": "b"
								}
							]
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
						"sources": [
							{
								"id": "my-other-source-id",
								"slug": "my-other-source-slug",
								"name": "My source",
								"workspaceId": "my-workspace-id",
								"enabled": false,
								"writeKeys": [
									"my-write-key"
								],
								"metadata": {
									"id": "my-metadata-id",
									"slug": "javascript",
									"name": "Javascript",
									"categories": [
										"Website"
									],
									"description": "This is our most flexible and powerful tracking system, using analytics.js.  Track and analyze information about your visitors and customers, and every action that they take, in any of our 140 integrations, business intelligence tools, or directly with SQL tools.",
									"logos": {
										"default": "https://cdn.filepicker.io/api/file/aRgo4XJQZausZxD4gZQq",
										"alt": "https://cdn.filepicker.io/api/file/aRgo4XJQZausZxD4gZQq",
										"mark": "https://cdn.filepicker.io/api/file/kBpmEoSSaakidAvoFmzd"
									},
									"options": [],
									"isCloudEventSource": false
								},
								"settings": {
									"beep": "beep"
								},
								"labels": [
									{
										"key": "product",
										"value": "b"
									}
								]
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
						resource "segment_source_tracking_plan_connection" "test" {
							source_id = "my-source-id"
							tracking_plan_id = "my-tracking-plan-id"
						}
					`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("segment_source_tracking_plan_connection.test", "source_id", "my-source-id"),
					resource.TestCheckResourceAttr("segment_source_tracking_plan_connection.test", "tracking_plan_id", "my-tracking-plan-id"),
				),
			},
			// Update and Read testing
			{
				Config: providerConfig + `
						resource "segment_source_tracking_plan_connection" "test" {
							source_id = "my-other-source-id"
							tracking_plan_id = "my-other-tracking-plan-id"
						}
					`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("segment_source_tracking_plan_connection.test", "source_id", "my-other-source-id"),
					resource.TestCheckResourceAttr("segment_source_tracking_plan_connection.test", "tracking_plan_id", "my-other-tracking-plan-id"),
				),
			},
		},
	})
}
