package provider

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTrackingPlanResource(t *testing.T) {
	t.Parallel()

	updatedRules := 0
	updated := 0
	fakeServer := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("content-type", "application/json")
			payload := ""

			// Rules requests
			if req.URL.Path == "/tracking-plans/my-tracking-plan-id/rules" {
				payload = `
				{
					"data": {
						"rules": [
							{
								"key": "Add Rule",
								"type": "TRACK",
								"version": 1,
								"jsonSchema": {
									"properties": {
										"context": {},
										"traits": {},
										"properties": {}
									}
								},
								"createdAt": "2023-09-08T19:02:55.000Z",
								"updatedAt": "2023-09-08T19:02:55.000Z",
								"deprecatedAt": "0001-01-01T00:00:00.000Z"
							}
						],
						"pagination": {
							"current": "MA==",
							"totalEntries": 1
						}
					}
				}
			`

				// After we update the source, return the updated source for subsequent calls (first update is part of the create call)
				if req.Method == http.MethodPost {
					updatedRules++
				}
				if updatedRules > 1 {
					payload = `
					{
						"data": {
							"rules": [
								{
									"type": "IDENTIFY",
									"version": 2,
									"jsonSchema": {
										"properties": {
											"context": {},
											"traits": {},
											"properties": {}
										}
									},
									"createdAt": "2023-09-08T19:02:55.000Z",
									"updatedAt": "2023-09-08T19:02:55.000Z",
									"deprecatedAt": "0001-01-01T00:00:00.000Z"
								}
							],
							"pagination": {
								"current": "MA==",
								"totalEntries": 1
							}
						}
					}
				`
				}
			} else if req.URL.Path == "/tracking-plans/my-tracking-plan-id" || req.URL.Path == "/tracking-plans" { // Tracking Plan requests
				payload = `
				{
					"data": {
						"trackingPlan": {
							"id": "my-tracking-plan-id",
							"name": "My Tracking Plan",
							"resourceSchemaId": "my-resource-schema-id",
							"slug": "my-tracking-plan-slug",
							"description": "My Tracking Plan Description",
							"type": "LIVE",
							"updatedAt": "2021-11-16T00:06:19.000Z",
							"createdAt": "2021-11-16T00:06:19.000Z"
						}
					}
				}
			`

				if req.Method == http.MethodPatch {
					updated++
				}
				if updated > 0 {
					payload = `
					{
						"data": {
							"trackingPlan": {
								"id": "my-tracking-plan-id",
								"name": "My New Tracking Plan",
								"resourceSchemaId": "my-resource-schema-id",
								"slug": "my-tracking-plan-slug",
								"description": "My New Tracking Plan Description",
								"type": "LIVE",
								"updatedAt": "2021-11-16T00:06:19.000Z",
								"createdAt": "2021-11-16T00:06:19.000Z"
							}
						}
					}
				`
				}
			} else {
				payload = `{
				"data": {}
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
					resource "segment_tracking_plan" "test" {
						name = "My Tracking Plan"
						type = "LIVE"
						description = "My Tracking Plan Description"
						rules = [
							{
								key     = "Add Rule"
								type    = "TRACK"
								version = 1
								json_schema = jsonencode({									
									"properties": {
										"context": {},
										"traits": {},
										"properties": {}
									}
								})
							  }
						]
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("segment_tracking_plan.test", "id", "my-tracking-plan-id"),
					resource.TestCheckResourceAttr("segment_tracking_plan.test", "name", "My Tracking Plan"),
					resource.TestCheckResourceAttr("segment_tracking_plan.test", "slug", "my-tracking-plan-slug"),
					resource.TestCheckResourceAttr("segment_tracking_plan.test", "description", "My Tracking Plan Description"),
					resource.TestCheckResourceAttr("segment_tracking_plan.test", "type", "LIVE"),
					resource.TestCheckResourceAttr("segment_tracking_plan.test", "updated_at", "2021-11-16T00:06:19.000Z"),
					resource.TestCheckResourceAttr("segment_tracking_plan.test", "created_at", "2021-11-16T00:06:19.000Z"),
					resource.TestCheckResourceAttr("segment_tracking_plan.test", "rules.#", "1"),
					resource.TestCheckResourceAttr("segment_tracking_plan.test", "rules.0.key", "Add Rule"),
					resource.TestCheckResourceAttr("segment_tracking_plan.test", "rules.0.type", "TRACK"),
					resource.TestCheckResourceAttr("segment_tracking_plan.test", "rules.0.version", "1"),
					resource.TestCheckResourceAttr("segment_tracking_plan.test", "rules.0.json_schema", "{\"properties\":{\"context\":{},\"properties\":{},\"traits\":{}}}"),
				),
			},
			// ImportState testing
			{
				ResourceName: "segment_tracking_plan.test",
				Config: providerConfig + `
					resource "segment_tracking_plan" "test" {
						name = "My Tracking Plan"
						type = "LIVE"
						description = "My Tracking Plan Description"
						rules = [
							{
								key     = "Add Rule"
								type    = "TRACK"
								version = 1
								json_schema = jsonencode({									
									"properties": {
										"context": {},
										"traits": {},
										"properties": {}
									}
								})
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
					resource "segment_tracking_plan" "test" {
						name = "My New Tracking Plan"
						type = "LIVE"
						description = "My New Tracking Plan Description"
						rules = [
							{
								type    = "IDENTIFY"
								version = 2
								json_schema = jsonencode({})
							}
						]
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("segment_tracking_plan.test", "id", "my-tracking-plan-id"),
					resource.TestCheckResourceAttr("segment_tracking_plan.test", "name", "My New Tracking Plan"),
					resource.TestCheckResourceAttr("segment_tracking_plan.test", "slug", "my-tracking-plan-slug"),
					resource.TestCheckResourceAttr("segment_tracking_plan.test", "description", "My New Tracking Plan Description"),
					resource.TestCheckResourceAttr("segment_tracking_plan.test", "type", "LIVE"),
					resource.TestCheckResourceAttr("segment_tracking_plan.test", "updated_at", "2021-11-16T00:06:19.000Z"),
					resource.TestCheckResourceAttr("segment_tracking_plan.test", "created_at", "2021-11-16T00:06:19.000Z"),
					resource.TestCheckResourceAttr("segment_tracking_plan.test", "rules.#", "1"),
					resource.TestCheckNoResourceAttr("segment_tracking_plan.test", "rules.0.key"),
					resource.TestCheckResourceAttr("segment_tracking_plan.test", "rules.0.type", "IDENTIFY"),
					resource.TestCheckResourceAttr("segment_tracking_plan.test", "rules.0.version", "2"),
					resource.TestCheckResourceAttr("segment_tracking_plan.test", "rules.0.json_schema", "{}"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
