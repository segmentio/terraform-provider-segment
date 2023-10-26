package provider

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTransformationResource(t *testing.T) {
	t.Parallel()

	updated := 0
	fakeServer := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("content-type", "application/json")

			payload := ""
			if req.URL.Path == "/transformations" && req.Method == http.MethodPost {
				payload = `
					{
						"data": {
							"transformation": {
								"id": "my-transformation-id",
								"name": "My transformation name",
								"workspaceId": "my-workspace-id",
								"sourceId": "my-source-id",
								"enabled": true,
								"if": "event = 'Bad Event'",
								"newEventName": "Good Event",
								"propertyRenames": [
									{
										"oldName": "old-name",
										"newName": "new-name"
									}
								],
								"propertyValueTransformations": [
									{
										"propertyPaths": [
											"properties.some-property",
											"context.some-property"
										],
										"propertyValue": "some property value"
									}
								],
								"fqlDefinedProperties": []
							}
						}
					}
				`
			} else if req.URL.Path == "/transformations/my-transformation-id" && req.Method == http.MethodPatch {
				payload = `
					{
						"data": {
							"transformation": {
								"id": "my-transformation-id",
								"name": "My new transformation name",
								"workspaceId": "my-workspace-id",
								"sourceId": "my-other-source-id",
								"enabled": false,
								"if": "event = 'Good Event'",
								"newEventName": "Bad Event",
								"propertyRenames": [],
								"propertyValueTransformations": [],
								"fqlDefinedProperties": [
									{
										"fql": "event = 'Good Event'",
										"propertyName": "some-property"
									}
								]
							}
						}
					}
				`
				updated++
			} else if req.URL.Path == "/transformations/my-transformation-id" && req.Method == http.MethodGet {
				if updated == 0 {
					payload = `
						{
							"data": {
								"transformation": {
									"id": "my-transformation-id",
									"name": "My transformation name",
									"workspaceId": "my-workspace-id",
									"sourceId": "my-source-id",
									"enabled": true,
									"if": "event = 'Bad Event'",
									"newEventName": "Good Event",
									"propertyRenames": [
										{
											"oldName": "old-name",
											"newName": "new-name"
										}
									],
									"propertyValueTransformations": [
										{
											"propertyPaths": [
												"properties.some-property",
												"context.some-property"
											],
											"propertyValue": "some property value"
										}
									],
									"fqlDefinedProperties": []
								}
							}
						}
					`
				} else {
					payload = `
						{
							"data": {
								"transformation": {
									"id": "my-transformation-id",
									"name": "My new transformation name",
									"workspaceId": "my-workspace-id",
									"sourceId": "my-other-source-id",
									"enabled": false,
									"if": "event = 'Good Event'",
									"newEventName": "Bad Event",
									"propertyRenames": [],
									"propertyValueTransformations": [],
									"fqlDefinedProperties": [
										{
											"fql": "event = 'Good Event'",
											"propertyName": "some-property"
										}
									]
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
					resource "segment_transformation" "test" {
						source_id      = "my-source-id"
						name           = "My transformation name"
						enabled        = true
						if             = "event = 'Bad Event'"
						new_event_name = "Good Event"
						property_renames = [
							{
								old_name = "old-name"
								new_name = "new-name"
							}
						]
						property_value_transformations = [
							{
								property_paths = ["properties.some-property", "context.some-property"],
								property_value = "some property value"
							},
						]
						fql_defined_properties = []
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("segment_transformation.test", "id", "my-transformation-id"),
					resource.TestCheckResourceAttr("segment_transformation.test", "source_id", "my-source-id"),
					resource.TestCheckResourceAttr("segment_transformation.test", "name", "My transformation name"),
					resource.TestCheckResourceAttr("segment_transformation.test", "enabled", "true"),
					resource.TestCheckResourceAttr("segment_transformation.test", "if", "event = 'Bad Event'"),
					resource.TestCheckResourceAttr("segment_transformation.test", "new_event_name", "Good Event"),
					resource.TestCheckResourceAttr("segment_transformation.test", "property_renames.#", "1"),
					resource.TestCheckResourceAttr("segment_transformation.test", "property_renames.0.old_name", "old-name"),
					resource.TestCheckResourceAttr("segment_transformation.test", "property_renames.0.new_name", "new-name"),
					resource.TestCheckResourceAttr("segment_transformation.test", "property_value_transformations.#", "1"),
					resource.TestCheckResourceAttr("segment_transformation.test", "property_value_transformations.0.property_paths.#", "2"),
					resource.TestCheckResourceAttr("segment_transformation.test", "property_value_transformations.0.property_paths.1", "properties.some-property"),
					resource.TestCheckResourceAttr("segment_transformation.test", "property_value_transformations.0.property_paths.0", "context.some-property"),
					resource.TestCheckResourceAttr("segment_transformation.test", "property_value_transformations.0.property_value", "some property value"),
					resource.TestCheckResourceAttr("segment_transformation.test", "fql_defined_properties.#", "0"),
				),
			},
			// ImportState testing
			{
				ResourceName: "segment_transformation.test",
				Config: providerConfig + `
					resource "segment_transformation" "test" {
						source_id      = "my-source-id"
						name           = "My transformation name"
						enabled        = true
						if             = "event = 'Bad Event'"
						new_event_name = "Good Event"
						property_renames = [
							{
								old_name = "old-name"
								new_name = "new-name"
							}
						]
						property_value_transformations = [
							{
								property_paths = ["properties.some-property", "context.some-property"],
								property_value = "some property value"
							},
						]
						fql_defined_properties = []
					}
				`,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: providerConfig + `
					resource "segment_transformation" "test" {
						source_id      = "my-other-source-id"
						name           = "My new transformation name"
						enabled        = false
						if             = "event = 'Good Event'"
						new_event_name = "Bad Event"
						property_renames = []
						property_value_transformations = []
						fql_defined_properties = [
							{
								fql           = "event = 'Good Event'"
								property_name = "some-property"
							}
						]
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("segment_transformation.test", "id", "my-transformation-id"),
					resource.TestCheckResourceAttr("segment_transformation.test", "source_id", "my-other-source-id"),
					resource.TestCheckResourceAttr("segment_transformation.test", "name", "My new transformation name"),
					resource.TestCheckResourceAttr("segment_transformation.test", "enabled", "false"),
					resource.TestCheckResourceAttr("segment_transformation.test", "if", "event = 'Good Event'"),
					resource.TestCheckResourceAttr("segment_transformation.test", "new_event_name", "Bad Event"),
					resource.TestCheckResourceAttr("segment_transformation.test", "property_renames.#", "0"),
					resource.TestCheckResourceAttr("segment_transformation.test", "property_value_transformations.#", "0"),
					resource.TestCheckResourceAttr("segment_transformation.test", "fql_defined_properties.#", "1"),
					resource.TestCheckResourceAttr("segment_transformation.test", "fql_defined_properties.0.fql", "event = 'Good Event'"),
					resource.TestCheckResourceAttr("segment_transformation.test", "fql_defined_properties.0.property_name", "some-property"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
