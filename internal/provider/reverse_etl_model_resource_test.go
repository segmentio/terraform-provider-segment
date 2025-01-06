package provider

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccReverseETLModelResource(t *testing.T) {
	t.Parallel()

	updated := 0
	fakeServer := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("content-type", "application/json")

			payload := ""
			if req.URL.Path == "/reverse-etl-models" && req.Method == http.MethodPost {
				payload = `
					{
						"data": {
							"reverseEtlModel": {
								"id": "my-reverse-etl-model-id",
								"sourceId": "my-source-id",
								"name": "My reverse etl model name",
								"description": "My reverse etl model description",
								"enabled": true,
								"query": "SELECT hi FROM greetings",
								"queryIdentifierColumn": "hi"
							}
						}
					}
				`
			} else if req.URL.Path == "/reverse-etl-models/my-reverse-etl-model-id" && req.Method == http.MethodPatch {
				payload = `
					{
						"data": {
							"reverseEtlModel": {
								"id": "my-reverse-etl-model-id",
								"sourceId": "my-source-id",
								"name": "My new reverse etl model name",
								"description": "My new reverse etl model description",
								"enabled": false,
								"query": "SELECT hello FROM greetings",
								"queryIdentifierColumn": "hello"
							}
						}
					}
				`
				updated++
			} else if req.URL.Path == "/reverse-etl-models/my-reverse-etl-model-id" && req.Method == http.MethodGet {
				if updated == 0 {
					payload = `
						{
							"data": {
								"reverseEtlModel": {
									"id": "my-reverse-etl-model-id",
									"sourceId": "my-source-id",
									"name": "My reverse etl model name",
									"description": "My reverse etl model description",
									"enabled": true,
									"query": "SELECT hi FROM greetings",
									"queryIdentifierColumn": "hi"
								}
							}
						}
					`
				} else {
					payload = `
						{
							"data": {
								"reverseEtlModel": {
									"id": "my-reverse-etl-model-id",
									"sourceId": "my-source-id",
									"name": "My new reverse etl model name",
									"description": "My new reverse etl model description",
									"enabled": false,
									"query": "SELECT hello FROM greetings",
									"queryIdentifierColumn": "hello"
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
					resource "segment_reverse_etl_model" "test" {
						source_id               = "my-source-id"
						name                    = "My reverse etl model name"
						enabled                 = true
						description             = "My reverse etl model description"
						query                   = "SELECT hi FROM greetings"
						query_identifier_column = "hi"
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("segment_reverse_etl_model.test", "id", "my-reverse-etl-model-id"),
					resource.TestCheckResourceAttr("segment_reverse_etl_model.test", "source_id", "my-source-id"),
					resource.TestCheckResourceAttr("segment_reverse_etl_model.test", "name", "My reverse etl model name"),
					resource.TestCheckResourceAttr("segment_reverse_etl_model.test", "enabled", "true"),
					resource.TestCheckResourceAttr("segment_reverse_etl_model.test", "description", "My reverse etl model description"),
					resource.TestCheckResourceAttr("segment_reverse_etl_model.test", "query", "SELECT hi FROM greetings"),
					resource.TestCheckResourceAttr("segment_reverse_etl_model.test", "query_identifier_column", "hi"),
				),
			},
			// ImportState testing
			{
				ResourceName: "segment_reverse_etl_model.test",
				Config: providerConfig + `
					resource "segment_reverse_etl_model" "test" {
						source_id               = "my-source-id"
						name                    = "My reverse etl model name"
						enabled                 = true
						description             = "My reverse etl model description"
						query                   = "SELECT hi FROM greetings"
						query_identifier_column = "hi"
					}
				`,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: providerConfig + `
					resource "segment_reverse_etl_model" "test" {
						source_id               = "my-source-id"
						name                    = "My new reverse etl model name"
						enabled                 = false
						description             = "My new reverse etl model description"
						query                   = "SELECT hello FROM greetings"
						query_identifier_column = "hello"
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("segment_reverse_etl_model.test", "id", "my-reverse-etl-model-id"),
					resource.TestCheckResourceAttr("segment_reverse_etl_model.test", "source_id", "my-source-id"),
					resource.TestCheckResourceAttr("segment_reverse_etl_model.test", "name", "My new reverse etl model name"),
					resource.TestCheckResourceAttr("segment_reverse_etl_model.test", "enabled", "false"),
					resource.TestCheckResourceAttr("segment_reverse_etl_model.test", "description", "My new reverse etl model description"),
					resource.TestCheckResourceAttr("segment_reverse_etl_model.test", "query", "SELECT hello FROM greetings"),
					resource.TestCheckResourceAttr("segment_reverse_etl_model.test", "query_identifier_column", "hello"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
