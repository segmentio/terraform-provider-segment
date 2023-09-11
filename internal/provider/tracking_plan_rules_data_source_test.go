package provider

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTrackingPlanRulesDataSource(t *testing.T) {
	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		fakeServer := httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("content-type", "application/json")
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
					Config: providerConfig + `data "segment_tracking_plan_rules" "test" { tracking_plan_id = "my-tracking-plan-id" }`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.segment_tracking_plan_rules.test", "tracking_plan_id", "my-tracking-plan-id"),
						resource.TestCheckResourceAttr("data.segment_tracking_plan_rules.test", "rules.#", "1"),
						resource.TestCheckResourceAttr("data.segment_tracking_plan_rules.test", "rules.0.key", "Add Rule"),
						resource.TestCheckResourceAttr("data.segment_tracking_plan_rules.test", "rules.0.type", "TRACK"),
						resource.TestCheckResourceAttr("data.segment_tracking_plan_rules.test", "rules.0.version", "1"),
						resource.TestCheckResourceAttr("data.segment_tracking_plan_rules.test", "rules.0.created_at", "2023-09-08T19:02:55.000Z"),
						resource.TestCheckResourceAttr("data.segment_tracking_plan_rules.test", "rules.0.updated_at", "2023-09-08T19:02:55.000Z"),
						resource.TestCheckResourceAttr("data.segment_tracking_plan_rules.test", "rules.0.deprecated_at", "0001-01-01T00:00:00.000Z"),
						resource.TestCheckResourceAttr("data.segment_tracking_plan_rules.test", "rules.0.json_schema", "{\"properties\":{\"context\":{},\"properties\":{},\"traits\":{}}}"),
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
					Config: providerConfig + `data "segment_tracking_plan_rules" "test" { tracking_plan_id = "my-tracking-plan-id" }`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.segment_tracking_plan_rules.test", "tracking_plan_id", "my-tracking-plan-id"),
						resource.TestCheckResourceAttr("data.segment_tracking_plan_rules.test", "rules.#", "1"),
						resource.TestCheckNoResourceAttr("data.segment_tracking_plan_rules.test", "rules.0.key"),
						resource.TestCheckResourceAttr("data.segment_tracking_plan_rules.test", "rules.0.type", "IDENTIFY"),
						resource.TestCheckResourceAttr("data.segment_tracking_plan_rules.test", "rules.0.version", "1"),
						resource.TestCheckResourceAttr("data.segment_tracking_plan_rules.test", "rules.0.json_schema", "{\"properties\":{\"context\":{},\"properties\":{},\"traits\":{}}}"),
						resource.TestCheckNoResourceAttr("data.segment_tracking_plan_rules.test", "rules.0.created_at"),
						resource.TestCheckNoResourceAttr("data.segment_tracking_plan_rules.test", "rules.0.updated_at"),
						resource.TestCheckNoResourceAttr("data.segment_tracking_plan_rules.test", "rules.0.deprecated_at"),
					),
				},
			},
		})
	})
}
