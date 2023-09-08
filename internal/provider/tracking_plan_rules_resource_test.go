package provider

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTrackingPlanRulesResource(t *testing.T) {
	t.Parallel()

	updated := 0
	fakeServer := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("content-type", "application/json")

			payload := `
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
				updated++
			}
			if updated > 1 {
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
					resource "segment_tracking_plan_rules" "test" {
						tracking_plan_id = "my-tracking-plan-id"
						rules = [
							{
								key     = "Add Rule"
								type    = "TRACK"
								version = 1
								json_schema = jsonencode({})
							  }
						]
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("segment_tracking_plan_rules.test", "tracking_plan_id", "my-tracking-plan-id"),
					resource.TestCheckResourceAttr("segment_tracking_plan_rules.test", "rules.#", "1"),
					resource.TestCheckResourceAttr("segment_tracking_plan_rules.test", "rules.0.key", "Add Rule"),
					resource.TestCheckResourceAttr("segment_tracking_plan_rules.test", "rules.0.type", "TRACK"),
					resource.TestCheckResourceAttr("segment_tracking_plan_rules.test", "rules.0.version", "1"),
					resource.TestCheckResourceAttr("segment_tracking_plan_rules.test", "rules.0.json_schema", "{}"),
				),
			},
			// Update and Read testing
			{
				Config: providerConfig + `
					resource "segment_tracking_plan_rules" "test" {
						tracking_plan_id = "my-tracking-plan-id"
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
					resource.TestCheckResourceAttr("segment_tracking_plan_rules.test", "tracking_plan_id", "my-tracking-plan-id"),
					resource.TestCheckResourceAttr("segment_tracking_plan_rules.test", "rules.#", "1"),
					resource.TestCheckNoResourceAttr("segment_tracking_plan_rules.test", "rules.0.key"),
					resource.TestCheckResourceAttr("segment_tracking_plan_rules.test", "rules.0.type", "IDENTIFY"),
					resource.TestCheckResourceAttr("segment_tracking_plan_rules.test", "rules.0.version", "2"),
					resource.TestCheckResourceAttr("segment_tracking_plan_rules.test", "rules.0.json_schema", "{}"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
