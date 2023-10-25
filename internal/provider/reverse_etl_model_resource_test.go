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
								"queryIdentifierColumn": "hi",
								"scheduleStrategy": "SPECIFIC_DAYS",
								"scheduleConfig": {"days":[0,1,2,3],"hours":[0,1,3,2],"timezone":"America/Los_Angeles"}
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
								"queryIdentifierColumn": "hello",
								"scheduleStrategy": "SPECIFIC_DAYS",
								"scheduleConfig": {"days":[0,1,2,3,4],"hours":[0,1,5],"timezone":"America/Los_Angeles"}
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
									"queryIdentifierColumn": "hi",
									"scheduleStrategy": "SPECIFIC_DAYS",
									"scheduleConfig": {"days":[0,1,2,3],"hours":[0,1,3,2],"timezone":"America/Los_Angeles"}
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
									"queryIdentifierColumn": "hello",
									"scheduleStrategy": "SPECIFIC_DAYS",
									"scheduleConfig": {"days":[0,1,2,3,4],"hours":[0,1,5],"timezone":"America/Los_Angeles"}
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
						schedule_strategy       = "SPECIFIC_DAYS"
						query                   = "SELECT hi FROM greetings"
						query_identifier_column = "hi"
						schedule_config         = jsonencode({
							"days": [0, 1, 2, 3],
							"hours": [0, 1, 3, 2],
							"timezone": "America/Los_Angeles"
						})
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("segment_reverse_etl_model.test", "id", "my-reverse-etl-model-id"),
					resource.TestCheckResourceAttr("segment_reverse_etl_model.test", "source_id", "my-source-id"),
					resource.TestCheckResourceAttr("segment_reverse_etl_model.test", "name", "My reverse etl model name"),
					resource.TestCheckResourceAttr("segment_reverse_etl_model.test", "enabled", "true"),
					resource.TestCheckResourceAttr("segment_reverse_etl_model.test", "description", "My reverse etl model description"),
					resource.TestCheckResourceAttr("segment_reverse_etl_model.test", "schedule_strategy", "SPECIFIC_DAYS"),
					resource.TestCheckResourceAttr("segment_reverse_etl_model.test", "query", "SELECT hi FROM greetings"),
					resource.TestCheckResourceAttr("segment_reverse_etl_model.test", "query_identifier_column", "hi"),
					resource.TestCheckResourceAttr("segment_reverse_etl_model.test", "schedule_config", "{\"days\":[0,1,2,3],\"hours\":[0,1,3,2],\"timezone\":\"America/Los_Angeles\"}"),
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
						schedule_strategy       = "SPECIFIC_DAYS"
						query                   = "SELECT hi FROM greetings"
						query_identifier_column = "hi"
						schedule_config         = jsonencode({
							"days": [0, 1, 2, 3],
							"hours": [0, 1, 3, 2],
							"timezone": "America/Los_Angeles"
						})
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
						schedule_strategy       = "SPECIFIC_DAYS"
						query                   = "SELECT hello FROM greetings"
						query_identifier_column = "hello"
						schedule_config         = jsonencode({
							"days": [0, 1, 2, 3, 4],
							"hours": [0, 1, 5],
							"timezone": "America/Los_Angeles"
						})
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("segment_reverse_etl_model.test", "id", "my-reverse-etl-model-id"),
					resource.TestCheckResourceAttr("segment_reverse_etl_model.test", "source_id", "my-source-id"),
					resource.TestCheckResourceAttr("segment_reverse_etl_model.test", "name", "My new reverse etl model name"),
					resource.TestCheckResourceAttr("segment_reverse_etl_model.test", "enabled", "false"),
					resource.TestCheckResourceAttr("segment_reverse_etl_model.test", "description", "My new reverse etl model description"),
					resource.TestCheckResourceAttr("segment_reverse_etl_model.test", "schedule_strategy", "SPECIFIC_DAYS"),
					resource.TestCheckResourceAttr("segment_reverse_etl_model.test", "query", "SELECT hello FROM greetings"),
					resource.TestCheckResourceAttr("segment_reverse_etl_model.test", "query_identifier_column", "hello"),
					resource.TestCheckResourceAttr("segment_reverse_etl_model.test", "schedule_config", "{\"days\":[0,1,2,3,4],\"hours\":[0,1,5],\"timezone\":\"America/Los_Angeles\"}"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
