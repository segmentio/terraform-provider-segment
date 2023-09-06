package provider

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTrackingPlanResource(t *testing.T) {
	t.Parallel()

	updated := 0
	fakeServer := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("content-type", "application/json")

			payload := `
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

			// After we update the source, return the updated source for subsequent calls (first update is part of the create call)
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
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
