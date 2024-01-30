package provider

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccInsertFunctionInstanceResource(t *testing.T) {
	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		updated := 0
		fakeServer := httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				w.Header().Set("content-type", "application/json")

				payload := ""
				if req.URL.Path == "/insert-function-instances" && req.Method == http.MethodPost {
					payload = `
						{
							"data": {
								"insertFunctionInstance": {
									"id": "my-instance-id",
									"name": "My instance name",
									"integrationId": "my-integration-id",
									"classId": "my-function-id",
									"enabled": true,
									"createdAt": "2023-11-01T18:38:01.349Z",
									"updatedAt": "2023-11-01T18:41:29.318Z",
									"settings": {
										"apiKey": "abc123"
									},
									"encryptedSettings": {}
								}
							}
						}
					`
				} else if req.URL.Path == "/insert-function-instances/my-instance-id" && req.Method == http.MethodPatch {
					payload = `
						{
							"data": {
								"insertFunctionInstance": {
									"id": "my-instance-id",
									"name": "My new instance name",
									"integrationId": "my-integration-id",
									"classId": "my-function-id",
									"enabled": false,
									"createdAt": "2023-11-01T18:38:01.349Z",
									"updatedAt": "2023-11-01T18:41:29.318Z",
									"settings": {
										"apiKey": "cba321"
									},
									"encryptedSettings": {}
								}
							}
						}
					`
					updated++
				} else if req.URL.Path == "/insert-function-instances/my-instance-id" && req.Method == http.MethodGet {
					if updated == 0 {
						payload = `
							{
								"data": {
									"insertFunctionInstance": {
										"id": "my-instance-id",
										"name": "My instance name",
										"integrationId": "my-integration-id",
										"classId": "my-function-id",
										"enabled": true,
										"createdAt": "2023-11-01T18:38:01.349Z",
										"updatedAt": "2023-11-01T18:41:29.318Z",
										"settings": {
											"apiKey": "abc123"
										},
										"encryptedSettings": {}
									}
								}
							}
						`
					} else {
						payload = `
							{
								"data": {
									"insertFunctionInstance": {
										"id": "my-instance-id",
										"name": "My new instance name",
										"integrationId": "my-integration-id",
										"classId": "my-function-id",
										"enabled": false,
										"createdAt": "2023-11-01T18:38:01.349Z",
										"updatedAt": "2023-11-01T18:41:29.318Z",
										"settings": {
											"apiKey": "cba321"
										},
										"encryptedSettings": {}
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
						resource "segment_insert_function_instance" "test" {
							integration_id = "my-integration-id"
							function_id = "my-function-id"
							name = "My instance name"
							enabled = true
							settings = jsonencode({"apiKey": "abc123"})
						}
					`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("segment_insert_function_instance.test", "id", "my-instance-id"),
						resource.TestCheckResourceAttr("segment_insert_function_instance.test", "integration_id", "my-integration-id"),
						resource.TestCheckResourceAttr("segment_insert_function_instance.test", "name", "My instance name"),
						resource.TestCheckResourceAttr("segment_insert_function_instance.test", "enabled", "true"),
						resource.TestCheckResourceAttr("segment_insert_function_instance.test", "function_id", "my-function-id"),
						resource.TestCheckResourceAttr("segment_insert_function_instance.test", "settings", "{\"apiKey\":\"abc123\"}"),
					),
				},
				// ImportState testing
				{
					ResourceName: "segment_insert_function_instance.test",
					Config: providerConfig + `
						resource "segment_insert_function_instance" "test" {
							integration_id = "my-integration-id"
							function_id = "my-function-id"
							name = "My instance name"
							enabled = true
							settings = jsonencode({"apiKey": "abc123"})
						}
					`,
					ImportState:       true,
					ImportStateVerify: true,
				},
				// Update and Read testing
				{
					Config: providerConfig + `
						resource "segment_insert_function_instance" "test" {
							integration_id = "my-integration-id"
							function_id = "my-function-id"
							name = "My new instance name"
							enabled = false
							settings = jsonencode({"apiKey": "cba321"})
						}
					`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("segment_insert_function_instance.test", "id", "my-instance-id"),
						resource.TestCheckResourceAttr("segment_insert_function_instance.test", "integration_id", "my-integration-id"),
						resource.TestCheckResourceAttr("segment_insert_function_instance.test", "name", "My new instance name"),
						resource.TestCheckResourceAttr("segment_insert_function_instance.test", "enabled", "false"),
						resource.TestCheckResourceAttr("segment_insert_function_instance.test", "function_id", "my-function-id"),
						resource.TestCheckResourceAttr("segment_insert_function_instance.test", "settings", "{\"apiKey\":\"cba321\"}"),
					),
				},
				// Delete testing automatically occurs in TestCase
			},
		})
	})

	t.Run("with prefix", func(t *testing.T) {
		updated := 0
		fakeServer := httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				w.Header().Set("content-type", "application/json")

				payload := ""
				if req.URL.Path == "/insert-function-instances" && req.Method == http.MethodPost {
					payload = `
					{
						"data": {
							"insertFunctionInstance": {
								"id": "my-instance-id",
								"name": "My instance name",
								"integrationId": "my-integration-id",
								"classId": "my-function-id",
								"enabled": true,
								"createdAt": "2023-11-01T18:38:01.349Z",
								"updatedAt": "2023-11-01T18:41:29.318Z",
								"settings": {
									"apiKey": "abc123"
								},
								"encryptedSettings": {}
							}
						}
					}
				`
				} else if req.URL.Path == "/insert-function-instances/my-instance-id" && req.Method == http.MethodPatch {
					payload = `
					{
						"data": {
							"insertFunctionInstance": {
								"id": "my-instance-id",
								"name": "My new instance name",
								"integrationId": "my-integration-id",
								"classId": "my-function-id",
								"enabled": false,
								"createdAt": "2023-11-01T18:38:01.349Z",
								"updatedAt": "2023-11-01T18:41:29.318Z",
								"settings": {
									"apiKey": "cba321"
								},
								"encryptedSettings": {}
							}
						}
					}
				`
					updated++
				} else if req.URL.Path == "/insert-function-instances/my-instance-id" && req.Method == http.MethodGet {
					if updated == 0 {
						payload = `
						{
							"data": {
								"insertFunctionInstance": {
									"id": "my-instance-id",
									"name": "My instance name",
									"integrationId": "my-integration-id",
									"classId": "my-function-id",
									"enabled": true,
									"createdAt": "2023-11-01T18:38:01.349Z",
									"updatedAt": "2023-11-01T18:41:29.318Z",
									"settings": {
										"apiKey": "abc123"
									},
									"encryptedSettings": {}
								}
							}
						}
					`
					} else {
						payload = `
						{
							"data": {
								"insertFunctionInstance": {
									"id": "my-instance-id",
									"name": "My new instance name",
									"integrationId": "my-integration-id",
									"classId": "my-function-id",
									"enabled": false,
									"createdAt": "2023-11-01T18:38:01.349Z",
									"updatedAt": "2023-11-01T18:41:29.318Z",
									"settings": {
										"apiKey": "cba321"
									},
									"encryptedSettings": {}
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
				// Create and Read testing for 'ifnd_' prefixed fn
				{
					Config: providerConfig + `
					resource "segment_insert_function_instance" "test_ifnd" {
						integration_id = "my-integration-id"
						function_id = "ifnd_my-function-id"
						name = "My instance name"
						enabled = true
						settings = jsonencode({"apiKey": "abc123"})
					}
				`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("segment_insert_function_instance.test_ifnd", "id", "my-instance-id"),
						resource.TestCheckResourceAttr("segment_insert_function_instance.test_ifnd", "integration_id", "my-integration-id"),
						resource.TestCheckResourceAttr("segment_insert_function_instance.test_ifnd", "name", "My instance name"),
						resource.TestCheckResourceAttr("segment_insert_function_instance.test_ifnd", "enabled", "true"),
						resource.TestCheckResourceAttr("segment_insert_function_instance.test_ifnd", "function_id", "ifnd_my-function-id"),
						resource.TestCheckResourceAttr("segment_insert_function_instance.test_ifnd", "settings", "{\"apiKey\":\"abc123\"}"),
					),
				},
			},
		})
	})

}
