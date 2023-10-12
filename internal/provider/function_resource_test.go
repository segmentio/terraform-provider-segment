package provider

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFunctionResource(t *testing.T) {
	t.Parallel()

	updated := 0
	fakeServer := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("content-type", "application/json")
			var payload string

			// Rules requests
			if req.URL.Path == "/functions" && req.Method == http.MethodPost {
				payload = `
				{
					"data": {
						"function": {
							"id": "my-function-id",
							"workspaceId": "my-workspace-id",
							"displayName": "My test function",
							"description": "My function description",
							"logoUrl": "https://segment.com/cool-logo.png",
							"code": "// My test code!",
							"createdAt": "2023-10-11T18:52:07.087Z",
							"createdBy": "my-user-id",
							"previewWebhookUrl": "",
							"settings": [
								{
									"name": "mySettingName",
									"label": "My setting label",
									"description": "My setting description",
									"type": "STRING",
									"required": false,
									"sensitive": false
								}
							],
							"buildpack": "boreal",
							"catalogId": "my-catalog-id",
							"batchMaxCount": 0,
							"resourceType": "SOURCE"
						}
					}
				}`
			} else if req.URL.Path == "/functions/my-function-id" && req.Method == http.MethodPatch {
				updated++
				payload = `
				{
					"data": {
						"function": {
							"id": "my-function-id",
							"workspaceId": "my-workspace-id",
							"displayName": "My new test function",
							"description": "My new function description",
							"logoUrl": "https://segment.com/cool-other-logo.png",
							"code": "// My new test code!",
							"createdAt": "2023-10-11T18:52:07.087Z",
							"createdBy": "my-user-id",
							"previewWebhookUrl": "",
							"settings": [
								{
									"name": "myNewSettingName",
									"label": "My new setting label",
									"description": "My new setting description",
									"type": "STRING",
									"required": true,
									"sensitive": true
								}
							],
							"buildpack": "boreal",
							"catalogId": "my-catalog-id",
							"batchMaxCount": 0,
							"resourceType": "SOURCE"
						}
					}
				}`
			} else if req.URL.Path == "/functions/my-function-id" && req.Method == http.MethodGet {
				if updated == 0 {
					payload = `
					{
						"data": {
							"function": {
								"id": "my-function-id",
								"workspaceId": "my-workspace-id",
								"displayName": "My test function",
								"description": "My function description",
								"logoUrl": "https://segment.com/cool-logo.png",
								"code": "// My test code!",
								"createdAt": "2023-10-11T18:52:07.087Z",
								"createdBy": "my-user-id",
								"previewWebhookUrl": "",
								"settings": [
									{
										"name": "mySettingName",
										"label": "My setting label",
										"description": "My setting description",
										"type": "STRING",
										"required": false,
										"sensitive": false
									}
								],
								"buildpack": "boreal",
								"catalogId": "my-catalog-id",
								"batchMaxCount": 0,
								"resourceType": "SOURCE"
							}
						}
					}`
				} else {
					payload = `
					{
						"data": {
							"function": {
								"id": "my-function-id",
								"workspaceId": "my-workspace-id",
								"displayName": "My new test function",
								"description": "My new function description",
								"logoUrl": "https://segment.com/cool-other-logo.png",
								"code": "// My new test code!",
								"createdAt": "2023-10-11T18:52:07.087Z",
								"createdBy": "my-user-id",
								"previewWebhookUrl": "",
								"settings": [
									{
										"name": "myNewSettingName",
										"label": "My new setting label",
										"description": "My new setting description",
										"type": "STRING",
										"required": true,
										"sensitive": true
									}
								],
								"buildpack": "boreal",
								"catalogId": "my-catalog-id",
								"batchMaxCount": 0,
								"resourceType": "SOURCE"
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
					resource "segment_function" "test" {
						code = "// My test code!"
						display_name = "My test function"
						logo_url = "https://segment.com/cool-logo.png"
						resource_type = "SOURCE"
						description = "My function description"
						settings = [
						{
							name = "mySettingName"
							label = "My setting label"
							type = "STRING"
							description = "My setting description"
							required = false
							sensitive = false
						},
						]
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("segment_function.test", "id", "my-function-id"),
					resource.TestCheckResourceAttr("segment_function.test", "code", "// My test code!"),
					resource.TestCheckResourceAttr("segment_function.test", "display_name", "My test function"),
					resource.TestCheckResourceAttr("segment_function.test", "logo_url", "https://segment.com/cool-logo.png"),
					resource.TestCheckResourceAttr("segment_function.test", "resource_type", "SOURCE"),
					resource.TestCheckResourceAttr("segment_function.test", "description", "My function description"),
					resource.TestCheckResourceAttr("segment_function.test", "preview_webhook_url", ""),
					resource.TestCheckResourceAttr("segment_function.test", "catalog_id", "my-catalog-id"),
					resource.TestCheckResourceAttr("segment_function.test", "settings.#", "1"),
					resource.TestCheckResourceAttr("segment_function.test", "settings.0.name", "mySettingName"),
					resource.TestCheckResourceAttr("segment_function.test", "settings.0.label", "My setting label"),
					resource.TestCheckResourceAttr("segment_function.test", "settings.0.description", "My setting description"),
					resource.TestCheckResourceAttr("segment_function.test", "settings.0.type", "STRING"),
					resource.TestCheckResourceAttr("segment_function.test", "settings.0.required", "false"),
					resource.TestCheckResourceAttr("segment_function.test", "settings.0.sensitive", "false"),
				),
			},
			// ImportState testing
			{
				ResourceName: "segment_function.test",
				Config: providerConfig + `
					resource "segment_function" "test" {
						code = "// My test code!"
						display_name = "My test function"
						logo_url = "https://segment.com/cool-logo.png"
						resource_type = "SOURCE"
						description = "My function description"
						settings = [
							{
								name = "mySettingName"
								label = "My setting label"
								type = "STRING"
								description = "My setting description"
								required = false
								sensitive = false
							},
						]
					}
				`,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: providerConfig + `
					resource "segment_function" "test" {
						code = "// My new test code!"
						display_name = "My new test function"
						logo_url = "https://segment.com/cool-other-logo.png"
						resource_type = "SOURCE"
						description = "My new function description"
						settings = [
							{
								name = "myNewSettingName"
								label = "My new setting label"
								type = "STRING"
								description = "My new setting description"
								required = true
								sensitive = true
							},
						]
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("segment_function.test", "id", "my-function-id"),
					resource.TestCheckResourceAttr("segment_function.test", "code", "// My new test code!"),
					resource.TestCheckResourceAttr("segment_function.test", "display_name", "My new test function"),
					resource.TestCheckResourceAttr("segment_function.test", "logo_url", "https://segment.com/cool-other-logo.png"),
					resource.TestCheckResourceAttr("segment_function.test", "resource_type", "SOURCE"),
					resource.TestCheckResourceAttr("segment_function.test", "description", "My new function description"),
					resource.TestCheckResourceAttr("segment_function.test", "preview_webhook_url", ""),
					resource.TestCheckResourceAttr("segment_function.test", "catalog_id", "my-catalog-id"),
					resource.TestCheckResourceAttr("segment_function.test", "settings.#", "1"),
					resource.TestCheckResourceAttr("segment_function.test", "settings.0.name", "myNewSettingName"),
					resource.TestCheckResourceAttr("segment_function.test", "settings.0.label", "My new setting label"),
					resource.TestCheckResourceAttr("segment_function.test", "settings.0.description", "My new setting description"),
					resource.TestCheckResourceAttr("segment_function.test", "settings.0.type", "STRING"),
					resource.TestCheckResourceAttr("segment_function.test", "settings.0.required", "true"),
					resource.TestCheckResourceAttr("segment_function.test", "settings.0.sensitive", "true"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
