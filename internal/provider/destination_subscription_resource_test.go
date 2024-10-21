package provider

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDestinationSubscriptionResource(t *testing.T) {
	t.Parallel()

	updated := 0
	fakeServer := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("content-type", "application/json")

			payload := ""
			if req.URL.Path == "/destinations/my-destination-id/subscriptions" && req.Method == http.MethodPost {
				payload = `
					{
						"data": {
							"destinationSubscription": {
								"id": "my-subscription-id",
								"name": "My subscription name",
								"actionId": "my-action-id",
								"actionSlug": "my-action-slug",
								"destinationId": "my-destination-id",
								"modelId": "",
								"enabled": true,
								"trigger": "type = \"track\"",
								"settings": {}
							}
						}
					}
				`
			} else if req.URL.Path == "/destinations/my-destination-id/subscriptions/my-subscription-id" && req.Method == http.MethodPatch {
				// First update is to set the model id
				if updated < 1 {
					payload = `
					{
						"data": {
							"subscription": {
								"id": "my-subscription-id",
								"name": "My subscription name",
								"actionId": "my-action-id",
								"actionSlug": "my-action-slug",
								"destinationId": "my-destination-id",
								"modelId": "",
								"enabled": true,
								"trigger": "type = \"track\"",
								"settings": {}
							}
						}
					}
				`
				} else {
					payload = `
					{
						"data": {
							"subscription": {
								"id": "my-subscription-id",
								"name": "My new subscription name",
								"actionId": "my-action-id",
								"actionSlug": "my-action-slug",
								"destinationId": "my-destination-id",
								"modelId": "",
								"enabled": false,
								"trigger": "type = \"track\"",
								"settings": {
									"test": "test"
								}
							}
						}
					}
				`
				}

				updated++
			} else if req.URL.Path == "/destinations/my-destination-id/subscriptions/my-subscription-id" && req.Method == http.MethodGet {
				// First update is to set the model id
				if updated <= 1 {
					payload = `
						{
							"data": {
								"subscription": {
									"id": "my-subscription-id",
									"name": "My subscription name",
									"actionId": "my-action-id",
									"actionSlug": "my-action-slug",
									"destinationId": "my-destination-id",
									"modelId": "",
									"enabled": true,
									"trigger": "type = \"track\"",
									"settings": {}
								}
							}
						}
					`
				} else {
					payload = `
						{
							"data": {
								"subscription": {
									"id": "my-subscription-id",
									"name": "My new subscription name",
									"actionId": "my-action-id",
									"actionSlug": "my-action-slug",
									"destinationId": "my-destination-id",
									"modelId": "",
									"enabled": false,
									"trigger": "type = \"track\"",
									"settings": {
										"test": "test"
									}
								}
							}
						}
					`
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
					resource "segment_destination_subscription" "test" {
						destination_id = "my-destination-id"
						name = "My subscription name"
						enabled = true
						action_id = "my-action-id"
						trigger = "type = \"track\""
						settings = jsonencode({})
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("segment_destination_subscription.test", "id", "my-subscription-id"),
					resource.TestCheckResourceAttr("segment_destination_subscription.test", "destination_id", "my-destination-id"),
					resource.TestCheckResourceAttr("segment_destination_subscription.test", "name", "My subscription name"),
					resource.TestCheckResourceAttr("segment_destination_subscription.test", "enabled", "true"),
					resource.TestCheckResourceAttr("segment_destination_subscription.test", "action_id", "my-action-id"),
					resource.TestCheckResourceAttr("segment_destination_subscription.test", "action_slug", "my-action-slug"),
					resource.TestCheckResourceAttr("segment_destination_subscription.test", "trigger", "type = \"track\""),
					resource.TestCheckNoResourceAttr("segment_destination_subscription.test", "model_id"),
					resource.TestCheckResourceAttr("segment_destination_subscription.test", "settings", "{}"),
				),
			},
			// ImportState testing
			{
				ResourceName: "segment_destination_subscription.test",
				Config: providerConfig + `
					resource "segment_destination_subscription" "test" {
						destination_id = "my-destination-id"
						name = "My subscription name"
						enabled = true
						action_id = "my-action-id"
						trigger = "type = \"track\""
						settings = jsonencode({})
					}
				`,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "my-destination-id:my-subscription-id",
			},
			// Update and Read testing
			{
				Config: providerConfig + `
					resource "segment_destination_subscription" "test" {
						destination_id = "my-destination-id"
						name = "My new subscription name"
						enabled = false
						action_id = "my-action-id"
						trigger = "type = \"track\""
						settings = jsonencode({
							"test": "test"
						})
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("segment_destination_subscription.test", "id", "my-subscription-id"),
					resource.TestCheckResourceAttr("segment_destination_subscription.test", "destination_id", "my-destination-id"),
					resource.TestCheckResourceAttr("segment_destination_subscription.test", "name", "My new subscription name"),
					resource.TestCheckResourceAttr("segment_destination_subscription.test", "enabled", "false"),
					resource.TestCheckResourceAttr("segment_destination_subscription.test", "action_id", "my-action-id"),
					resource.TestCheckResourceAttr("segment_destination_subscription.test", "action_slug", "my-action-slug"),
					resource.TestCheckResourceAttr("segment_destination_subscription.test", "trigger", "type = \"track\""),
					resource.TestCheckResourceAttr("segment_destination_subscription.test", "settings", "{\"test\":\"test\"}"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccDestinationSubscriptionResourceWithModel(t *testing.T) {
	t.Parallel()

	updated := 0
	fakeServer := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("content-type", "application/json")

			payload := ""
			if req.URL.Path == "/destinations/my-destination-id/subscriptions" && req.Method == http.MethodPost {
				payload = `
					{
						"data": {
							"destinationSubscription": {
								"id": "my-subscription-id",
								"name": "My subscription name",
								"actionId": "my-action-id",
								"actionSlug": "my-action-slug",
								"destinationId": "my-destination-id",
								"modelId": "my-model-id",
								"reverseETLSchedule": {
									"config": { "interval": "1d" },
									"strategy": "PERIODIC"
								},
								"enabled": true,
								"trigger": "type = \"track\"",
								"settings": {}
							}
						}
					}
				`
			} else if req.URL.Path == "/destinations/my-destination-id/subscriptions/my-subscription-id" && req.Method == http.MethodPatch {
				if updated < 1 {
					payload = `
					{
						"data": {
							"subscription": {
								"id": "my-subscription-id",
								"name": "My subscription name",
								"actionId": "my-action-id",
								"actionSlug": "my-action-slug",
								"destinationId": "my-destination-id",
								"modelId": "my-model-id",
								"reverseETLSchedule": {
									"config": { "interval": "1d" },
									"strategy": "PERIODIC"
								},								"enabled": true,
								"trigger": "type = \"track\"",
								"settings": {}
							}
						}
					}
				`
				} else {
					payload = `
					{
						"data": {
							"subscription": {
								"id": "my-subscription-id",
								"name": "My new subscription name",
								"actionId": "my-action-id",
								"actionSlug": "my-action-slug",
								"destinationId": "my-destination-id",
								"modelId": "my-model-id",
								"reverseETLSchedule": {
									"config": { "interval": "1d" },
									"strategy": "PERIODIC"
								},								"enabled": false,
								"trigger": "type = \"track\"",
								"settings": {
									"test": "test"
								}
							}
						}
					}
				`
				}

				updated++
			} else if req.URL.Path == "/destinations/my-destination-id/subscriptions/my-subscription-id" && req.Method == http.MethodGet {
				if updated <= 1 {
					payload = `
						{
							"data": {
								"subscription": {
									"id": "my-subscription-id",
									"name": "My subscription name",
									"actionId": "my-action-id",
									"actionSlug": "my-action-slug",
									"destinationId": "my-destination-id",
									"modelId": "my-model-id",
									"reverseETLSchedule": {
										"config": { "interval": "1d" },
										"strategy": "PERIODIC"
									},									"enabled": true,
									"trigger": "type = \"track\"",
									"settings": {}
								}
							}
						}
					`
				} else {
					payload = `
						{
							"data": {
								"subscription": {
									"id": "my-subscription-id",
									"name": "My new subscription name",
									"actionId": "my-action-id",
									"actionSlug": "my-action-slug",
									"destinationId": "my-destination-id",
									"modelId": "my-model-id",
									"reverseETLSchedule": {
										"config": { "interval": "1d" },
										"strategy": "PERIODIC"
									},									"enabled": false,
									"trigger": "type = \"track\"",
									"settings": {
										"test": "test"
									}
								}
							}
						}
					`
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
					resource "segment_destination_subscription" "test" {
						destination_id = "my-destination-id"
						name = "My subscription name"
						enabled = true
						action_id = "my-action-id"
						trigger = "type = \"track\""
						settings = jsonencode({})
						model_id = "my-model-id"
						reverse_etl_schedule = {
							config = jsonencode({ interval = "1d" }),
							strategy = "PERIODIC"
						}
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("segment_destination_subscription.test", "id", "my-subscription-id"),
					resource.TestCheckResourceAttr("segment_destination_subscription.test", "destination_id", "my-destination-id"),
					resource.TestCheckResourceAttr("segment_destination_subscription.test", "name", "My subscription name"),
					resource.TestCheckResourceAttr("segment_destination_subscription.test", "enabled", "true"),
					resource.TestCheckResourceAttr("segment_destination_subscription.test", "action_id", "my-action-id"),
					resource.TestCheckResourceAttr("segment_destination_subscription.test", "action_slug", "my-action-slug"),
					resource.TestCheckResourceAttr("segment_destination_subscription.test", "trigger", "type = \"track\""),
					resource.TestCheckResourceAttr("segment_destination_subscription.test", "model_id", "my-model-id"),
					resource.TestCheckResourceAttr("segment_destination_subscription.test", "settings", "{}"),
					resource.TestCheckResourceAttr("segment_destination_subscription.test", "reverse_etl_schedule.strategy", "PERIODIC"),
					resource.TestCheckResourceAttr("segment_destination_subscription.test", "reverse_etl_schedule.config", "{\"interval\":\"1d\"}"),
				),
			},
			// ImportState testing
			{
				ResourceName: "segment_destination_subscription.test",
				Config: providerConfig + `
					resource "segment_destination_subscription" "test" {
						destination_id = "my-destination-id"
						name = "My subscription name"
						enabled = true
						action_id = "my-action-id"
						trigger = "type = \"track\""
						settings = jsonencode({})
						model_id = "my-model-id"
						reverse_etl_schedule = {
							config = jsonencode({ interval = "1d" }),
							strategy = "PERIODIC"
						}
					}
				`,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "my-destination-id:my-subscription-id",
			},
			// Update and Read testing
			{
				Config: providerConfig + `
					resource "segment_destination_subscription" "test" {
						destination_id = "my-destination-id"
						name = "My new subscription name"
						enabled = false
						action_id = "my-action-id"
						trigger = "type = \"track\""
						settings = jsonencode({
							"test": "test"
						})
						model_id = "my-model-id"
						reverse_etl_schedule = {
							config = jsonencode({ interval = "1d" }),
							strategy = "PERIODIC"
						}
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("segment_destination_subscription.test", "id", "my-subscription-id"),
					resource.TestCheckResourceAttr("segment_destination_subscription.test", "destination_id", "my-destination-id"),
					resource.TestCheckResourceAttr("segment_destination_subscription.test", "name", "My new subscription name"),
					resource.TestCheckResourceAttr("segment_destination_subscription.test", "enabled", "false"),
					resource.TestCheckResourceAttr("segment_destination_subscription.test", "action_id", "my-action-id"),
					resource.TestCheckResourceAttr("segment_destination_subscription.test", "action_slug", "my-action-slug"),
					resource.TestCheckResourceAttr("segment_destination_subscription.test", "trigger", "type = \"track\""),
					resource.TestCheckResourceAttr("segment_destination_subscription.test", "settings", "{\"test\":\"test\"}"),
					resource.TestCheckResourceAttr("segment_destination_subscription.test", "model_id", "my-model-id"),
					resource.TestCheckResourceAttr("segment_destination_subscription.test", "reverse_etl_schedule.strategy", "PERIODIC"),
					resource.TestCheckResourceAttr("segment_destination_subscription.test", "reverse_etl_schedule.config", "{\"interval\":\"1d\"}"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
