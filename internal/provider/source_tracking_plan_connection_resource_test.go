package provider

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSourceTrackingPlanConnectionResource(t *testing.T) {
	t.Parallel()
	updatedSchemaSettings := 0

	fakeServer := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("content-type", "application/json")

			payload := ""
			if req.URL.Path == "/sources/my-source-id" && req.Method == http.MethodGet {
				payload = `
					{
						"data": {
							"source": {
								"id": "my-source-id",
								"slug": "my-source-slug",
								"name": "My source name",
								"workspaceId": "my-workspace-id",
								"enabled": true,
								"writeKeys": ["my-write-key"],
								"metadata": {
									"id": "my-metadata-id",
									"slug": "my-metadata-slug",
									"name": "My metadata name",
									"categories": ["my-category"],
									"description": "My metadata description",
									"logos": {
										"default": "https://example.segment.com/image.png",
										"alt": "https://example.segment.com/image.png",
										"mark": "https://example.segment.com/image.png"
									},
									"options": [
										{
											"name": "sid",
											"required": true,
											"type": "string",
											"defaultValue": "",
											"description": "Your Segment SID",
											"defaultValue": "default-sid"
										}
									],
									"isCloudEventSource": false
								},
								"settings": {
									"myKey": "myValue"
								},
								"labels": [
									{
										"key": "my-label-key",
										"value": "my-label-value"
									}
								]
							},
							"trackingPlanId": "my-tracking-plan-id"
						}
					}
				`
			} else if req.URL.Path == "/tracking-plans/my-tracking-plan-id/sources" && req.Method == http.MethodPost {
				payload = `{
					"data": {
						"status": "CONNECTED"
					}
				}`
			} else if req.URL.Path == "/tracking-plans/my-tracking-plan-id/sources" && req.Method == http.MethodDelete {
				payload = `{
					"data": {
						"status": "SUCCESS"
					}
				}`
			} else if req.URL.Path == "/sources/my-source-id/settings" && req.Method == http.MethodGet {
				if updatedSchemaSettings <= 1 {
					payload = `
						{
							"data": {
								"sourceId": "my-source-id",
								"settings": {
									"track": {
										"allowUnplannedEvents": true,
										"allowUnplannedEventProperties": true,
										"allowEventOnViolations": true,
										"allowPropertiesOnViolations": true,
										"commonEventOnViolations": "OMIT_PROPERTIES"
									},
									"group": {
										"allowTraitsOnViolations": true,
										"allowUnplannedTraits": true,
										"commonEventOnViolations": "ALLOW"
									},
									"identify": {
										"allowTraitsOnViolations": true,
										"allowUnplannedTraits": true,
										"commonEventOnViolations": "BLOCK"
									},
									"forwardingBlockedEventsTo": "my-other-source-id",
									"forwardingViolationsTo": "my-other-source-id"
								}
							}
						}
					`
				} else {
					payload = `
						{
							"data": {
								"sourceId": "my-source-id",
								"settings": {
									"track": {
										"allowUnplannedEvents": true,
										"allowUnplannedEventProperties": true,
										"allowEventOnViolations": false,
										"allowPropertiesOnViolations": true,
										"commonEventOnViolations": "ALLOW"
									},
									"group": {
										"allowTraitsOnViolations": true,
										"allowUnplannedTraits": true,
										"commonEventOnViolations": "ALLOW"
									},
									"identify": {
										"allowTraitsOnViolations": false,
										"allowUnplannedTraits": false,
										"commonEventOnViolations": "ALLOW"
									},
									"forwardingBlockedEventsTo": "my-other-other-source-id"
								}
							}
						}
					`
				}

			} else if req.URL.Path == "/sources/my-source-id/settings" && req.Method == http.MethodPatch {
				if updatedSchemaSettings == 0 {
					payload = `
						{
							"data": {
								"sourceId": "my-source-id",
								"settings": {
									"track": {
										"allowUnplannedEvents": true,
										"allowUnplannedEventProperties": true,
										"allowEventOnViolations": true,
										"allowPropertiesOnViolations": true,
										"commonEventOnViolations": "OMIT_PROPERTIES"
									},
									"group": {
										"allowTraitsOnViolations": true,
										"allowUnplannedTraits": true,
										"commonEventOnViolations": "ALLOW"
									},
									"identify": {
										"allowTraitsOnViolations": true,
										"allowUnplannedTraits": true,
										"commonEventOnViolations": "BLOCK"
									},
									"forwardingBlockedEventsTo": "my-other-source-id",
									"forwardingViolationsTo": "my-other-source-id"
								}
							}
						}
					`
				} else {
					payload = `
						{
							"data": {
								"sourceId": "my-source-id",
								"settings": {
									"track": {
										"allowUnplannedEvents": true,
										"allowUnplannedEventProperties": true,
										"allowEventOnViolations": false,
										"allowPropertiesOnViolations": true,
										"commonEventOnViolations": "ALLOW"
									},
									"group": {
										"allowTraitsOnViolations": true,
										"allowUnplannedTraits": true,
										"commonEventOnViolations": "ALLOW"
									},
									"identify": {
										"allowTraitsOnViolations": false,
										"allowUnplannedTraits": false,
										"commonEventOnViolations": "ALLOW"
									},
									"forwardingBlockedEventsTo": "my-other-other-source-id"
								}
							}
						}
					`
				}

				updatedSchemaSettings++
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
					resource.TestCheckNoResourceAttr("segment_source_tracking_plan_connection.test", "schema_settings"),
				),
			},
			// Update and Read testing
			{
				Config: providerConfig + `
						resource "segment_source_tracking_plan_connection" "test" {
							source_id = "my-source-id"
							tracking_plan_id = "my-tracking-plan-id"
							schema_settings = {
								forwarding_blocked_events_to = "my-other-source-id"
								forwarding_violations_to = "my-other-source-id"
								track = {
									allow_unplanned_events = true
									allow_event_on_violations = true
									allow_properties_on_violations = true
									common_event_on_violations = "OMIT_PROPERTIES"
								}
								identify = {
									allow_traits_on_violations = true
									common_event_on_violations = "BLOCK"
								}
							}
						}
					`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("segment_source_tracking_plan_connection.test", "source_id", "my-source-id"),
					resource.TestCheckResourceAttr("segment_source_tracking_plan_connection.test", "tracking_plan_id", "my-tracking-plan-id"),
					resource.TestCheckResourceAttr("segment_source_tracking_plan_connection.test", "schema_settings.forwarding_blocked_events_to", "my-other-source-id"),
					resource.TestCheckResourceAttr("segment_source_tracking_plan_connection.test", "schema_settings.forwarding_violations_to", "my-other-source-id"),
					resource.TestCheckResourceAttr("segment_source_tracking_plan_connection.test", "schema_settings.track.allow_unplanned_events", "true"),
					resource.TestCheckResourceAttr("segment_source_tracking_plan_connection.test", "schema_settings.track.allow_event_on_violations", "true"),
					resource.TestCheckResourceAttr("segment_source_tracking_plan_connection.test", "schema_settings.track.allow_properties_on_violations", "true"),
					resource.TestCheckResourceAttr("segment_source_tracking_plan_connection.test", "schema_settings.track.common_event_on_violations", "OMIT_PROPERTIES"),
					resource.TestCheckNoResourceAttr("segment_source_tracking_plan_connection.test", "schema_settings.track.allow_unplanned_event_properties"),
					resource.TestCheckResourceAttr("segment_source_tracking_plan_connection.test", "schema_settings.identify.common_event_on_violations", "BLOCK"),
					resource.TestCheckResourceAttr("segment_source_tracking_plan_connection.test", "schema_settings.identify.allow_traits_on_violations", "true"),
					resource.TestCheckNoResourceAttr("segment_source_tracking_plan_connection.test", "schema_settings.identify.allow_unplanned_traits"),
					resource.TestCheckNoResourceAttr("segment_source_tracking_plan_connection.test", "schema_settings.group"),
				),
			},
			{
				ResourceName: "segment_source_tracking_plan_connection.test",
				Config: providerConfig + `
					resource "segment_source_tracking_plan_connection" "test" {
						source_id = "my-source-id"
						tracking_plan_id = "my-tracking-plan-id"
						schema_settings = {
							forwarding_blocked_events_to = "my-other-source-id"
							forwarding_violations_to = "my-other-source-id"
							track = {
								allow_unplanned_events = true
								allow_event_on_violations = true
								allow_properties_on_violations = true
								common_event_on_violations = "OMIT_PROPERTIES"
							}
							identify = {
								allow_traits_on_violations = true
								common_event_on_violations = "BLOCK"
							}
						}
					}
				`,
				ImportState:   true,
				ImportStateId: "my-source-id:my-tracking-plan-id",
			},
		},
	})
}
