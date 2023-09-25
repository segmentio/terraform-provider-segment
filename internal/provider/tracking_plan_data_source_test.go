package provider

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTrackingPlanDataSource(t *testing.T) {
	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		t.Parallel()
		fakeServer := httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				w.Header().Set("content-type", "application/json")

				if req.URL.Path == "/tracking-plans/my-tracking-plan-id/rules" {
					_, _ = w.Write([]byte(`
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
				`))
				} else if req.URL.Path == "/tracking-plans/my-tracking-plan-id" || req.URL.Path == "/tracking-plans" { // Tracking Plan requests
					_, _ = w.Write([]byte(`
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
				`))
				}
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
					Config: providerConfig + `data "segment_tracking_plan" "test" { id = "my-tracking-plan-id" }`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.segment_tracking_plan.test", "id", "my-tracking-plan-id"),
						resource.TestCheckResourceAttr("data.segment_tracking_plan.test", "name", "My Tracking Plan"),
						resource.TestCheckResourceAttr("data.segment_tracking_plan.test", "slug", "my-tracking-plan-slug"),
						resource.TestCheckResourceAttr("data.segment_tracking_plan.test", "description", "My Tracking Plan Description"),
						resource.TestCheckResourceAttr("data.segment_tracking_plan.test", "type", "LIVE"),
						resource.TestCheckResourceAttr("data.segment_tracking_plan.test", "updated_at", "2021-11-16T00:06:19.000Z"),
						resource.TestCheckResourceAttr("data.segment_tracking_plan.test", "created_at", "2021-11-16T00:06:19.000Z"),
						resource.TestCheckResourceAttr("data.segment_tracking_plan.test", "rules.#", "1"),
						resource.TestCheckResourceAttr("data.segment_tracking_plan.test", "rules.0.key", "Add Rule"),
						resource.TestCheckResourceAttr("data.segment_tracking_plan.test", "rules.0.type", "TRACK"),
						resource.TestCheckResourceAttr("data.segment_tracking_plan.test", "rules.0.version", "1"),
						resource.TestCheckResourceAttr("data.segment_tracking_plan.test", "rules.0.created_at", "2023-09-08T19:02:55.000Z"),
						resource.TestCheckResourceAttr("data.segment_tracking_plan.test", "rules.0.updated_at", "2023-09-08T19:02:55.000Z"),
						resource.TestCheckResourceAttr("data.segment_tracking_plan.test", "rules.0.deprecated_at", "0001-01-01T00:00:00.000Z"),
						resource.TestCheckResourceAttr("data.segment_tracking_plan.test", "rules.0.json_schema", "{\"properties\":{\"context\":{},\"properties\":{},\"traits\":{}}}"),
					),
				},
			},
		})
	})

	t.Run("nulls", func(t *testing.T) {
		t.Parallel()
		fakeServer := httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				w.Header().Set("content-type", "application/json")
				if req.URL.Path == "/tracking-plans/my-tracking-plan-id/rules" {
					_, _ = w.Write([]byte(`
						{
							"data": {
								"rules": [
									{
										"type": "IDENTIFY",
										"version": 1,
										"jsonSchema": {
											"properties": {
												"context": {},
												"traits": {},
												"properties": {}
											}
										}
									}
								],
								"pagination": {
									"current": "MA==",
									"totalEntries": 1
								}
							}
						}
					`))
				} else if req.URL.Path == "/tracking-plans/my-tracking-plan-id" || req.URL.Path == "/tracking-plans" { // Tracking Plan requests
					_, _ = w.Write([]byte(`
						{
							"data": {
								"trackingPlan": {
									"id": "my-tracking-plan-id",
									"resourceSchemaId": "my-resource-schema-id",
									"type": "LIVE"
								}
							}
						}
					`))
				}
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
					Config: providerConfig + `data "segment_tracking_plan" "test" { id = "my-tracking-plan-id" }`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.segment_tracking_plan.test", "id", "my-tracking-plan-id"),
						resource.TestCheckNoResourceAttr("data.segment_tracking_plan.test", "name"),
						resource.TestCheckNoResourceAttr("data.segment_tracking_plan.test", "slug"),
						resource.TestCheckNoResourceAttr("data.segment_tracking_plan.test", "description"),
						resource.TestCheckResourceAttr("data.segment_tracking_plan.test", "type", "LIVE"),
						resource.TestCheckNoResourceAttr("data.segment_tracking_plan.test", "updated_at"),
						resource.TestCheckNoResourceAttr("data.segment_tracking_plan.test", "created_at"),
						resource.TestCheckResourceAttr("data.segment_tracking_plan.test", "rules.#", "1"),
						resource.TestCheckNoResourceAttr("data.segment_tracking_plan.test", "rules.0.key"),
						resource.TestCheckResourceAttr("data.segment_tracking_plan.test", "rules.0.type", "IDENTIFY"),
						resource.TestCheckResourceAttr("data.segment_tracking_plan.test", "rules.0.version", "1"),
						resource.TestCheckResourceAttr("data.segment_tracking_plan.test", "rules.0.json_schema", "{\"properties\":{\"context\":{},\"properties\":{},\"traits\":{}}}"),
						resource.TestCheckNoResourceAttr("data.segment_tracking_plan.test", "rules.0.created_at"),
						resource.TestCheckNoResourceAttr("data.segment_tracking_plan.test", "rules.0.updated_at"),
						resource.TestCheckNoResourceAttr("data.segment_tracking_plan.test", "rules.0.deprecated_at"),
					),
				},
			},
		})
	})
}
