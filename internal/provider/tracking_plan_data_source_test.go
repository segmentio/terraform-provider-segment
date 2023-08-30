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
		fakeServer := httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("content-type", "application/json")
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
							"trackingPlan": {
								"id": "my-tracking-plan-id",
								"resourceSchemaId": "my-resource-schema-id",
								"type": "LIVE"
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
					Config: providerConfig + `data "segment_tracking_plan" "test" { id = "my-tracking-plan-id" }`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.segment_tracking_plan.test", "id", "my-tracking-plan-id"),
						resource.TestCheckNoResourceAttr("data.segment_tracking_plan.test", "name"),
						resource.TestCheckNoResourceAttr("data.segment_tracking_plan.test", "slug"),
						resource.TestCheckNoResourceAttr("data.segment_tracking_plan.test", "description"),
						resource.TestCheckResourceAttr("data.segment_tracking_plan.test", "type", "LIVE"),
						resource.TestCheckNoResourceAttr("data.segment_tracking_plan.test", "updated_at"),
						resource.TestCheckNoResourceAttr("data.segment_tracking_plan.test", "created_at"),
					),
				},
			},
		})
	})
}
