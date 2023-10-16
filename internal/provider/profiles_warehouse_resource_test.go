package provider

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProfilesWarehouseResource(t *testing.T) {
	t.Parallel()

	updated := 0
	fakeServer := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("content-type", "application/json")

			payload := ""
			if req.URL.Path == "/spaces/my-space-id/profiles-warehouses" && req.Method == http.MethodPost {
				payload = `
					{
						"data": {
							"profilesWarehouse": {
								"id": "my-warehouse-id",
								"spaceId": "my-space-id",
								"workspaceId": "my-workspace-id",
								"enabled": true,
								"metadata": {
									"id": "my-metadata-id",
									"slug": "snowflake",
									"name": "Snowflake",
									"description": "Data warehouse built for the cloud",
									"logos": {
										"default": "https://cdn.filepicker.io/api/file/JrQWOYvMRRCVvSHp4HL0",
										"mark": "https://cdn.filepicker.io/api/file/OBhrGoCRKaSyvAhDX3fw",
										"alt": ""
									},
									"options": []
								},
								"settings": {
									"name": "My warehouse name",
									"token": "my-token"
								},
								"schemaName": "my-schema-name"
							}
						}
					}
				`
			} else if req.URL.Path == "/spaces/my-space-id/profiles-warehouses/my-warehouse-id" && req.Method == http.MethodPatch {
				payload = `
					{
						"data": {
							"profilesWarehouse": {
								"id": "my-warehouse-id",
								"spaceId": "my-space-id",
								"workspaceId": "my-workspace-id",
								"enabled": false,
								"metadata": {
									"id": "my-metadata-id",
									"slug": "snowflake",
									"name": "Snowflake",
									"description": "Data warehouse built for the cloud",
									"logos": {
										"default": "https://cdn.filepicker.io/api/file/JrQWOYvMRRCVvSHp4HL0",
										"mark": "https://cdn.filepicker.io/api/file/OBhrGoCRKaSyvAhDX3fw",
										"alt": ""
									},
									"options": []
								},
								"settings": {
									"name": "My new warehouse name",
									"token": "my-other-token"
								},
								"schemaName": "my-new-schema-name"
							}
						}
					}
				`
				updated++
			} else if req.URL.Path == "/spaces/my-space-id/profiles-warehouses" && req.Method == http.MethodGet {
				if updated == 0 {
					payload = `
						{
							"data": {
								"profilesWarehouses": [
									{
										"id": "my-warehouse-id",
										"spaceId": "my-space-id",
										"workspaceId": "my-workspace-id",
										"enabled": true,
										"metadata": {
											"id": "my-metadata-id",
											"slug": "snowflake",
											"name": "Snowflake",
											"description": "Data warehouse built for the cloud",
											"logos": {
												"default": "https://cdn.filepicker.io/api/file/JrQWOYvMRRCVvSHp4HL0",
												"mark": "https://cdn.filepicker.io/api/file/OBhrGoCRKaSyvAhDX3fw",
												"alt": ""
											},
											"options": []
										},
										"settings": {
											"name": "My warehouse name",
											"token": "my-token"
										},
										"schemaName": "my-schema-name"
									}
								]
							}
						}
					`
				} else {
					payload = `
						{
							"data": {
								"profilesWarehouses": [
									{
										"id": "my-warehouse-id",
										"spaceId": "my-space-id",
										"workspaceId": "my-workspace-id",
										"enabled": false,
										"metadata": {
											"id": "my-metadata-id",
											"slug": "snowflake",
											"name": "Snowflake",
											"description": "Data warehouse built for the cloud",
											"logos": {
												"default": "https://cdn.filepicker.io/api/file/JrQWOYvMRRCVvSHp4HL0",
												"mark": "https://cdn.filepicker.io/api/file/OBhrGoCRKaSyvAhDX3fw",
												"alt": ""
											},
											"options": []
										},
										"settings": {
											"name": "My new warehouse name",
											"token": "my-other-token"
										},
										"schemaName": "my-new-schema-name"
									}
								]
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
					resource "segment_profiles_warehouse" "test" {
						space_id = "my-space-id"
						metadata_id = "my-metadata-id"
						name = "My warehouse name"
						enabled = true
						settings = jsonencode({
						  "token": "my-token"
						})
						schema_name = "my-schema-name"
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("segment_profiles_warehouse.test", "id", "my-warehouse-id"),
					resource.TestCheckResourceAttr("segment_profiles_warehouse.test", "space_id", "my-space-id"),
					resource.TestCheckResourceAttr("segment_profiles_warehouse.test", "name", "My warehouse name"),
					resource.TestCheckResourceAttr("segment_profiles_warehouse.test", "enabled", "true"),
					resource.TestCheckResourceAttr("segment_profiles_warehouse.test", "metadata_id", "my-metadata-id"),
					resource.TestCheckResourceAttr("segment_profiles_warehouse.test", "settings", "{\"token\":\"my-token\"}"),
					resource.TestCheckResourceAttr("segment_profiles_warehouse.test", "schema_name", "my-schema-name"),
				),
			},
			// ImportState testing
			{
				ResourceName: "segment_profiles_warehouse.test",
				Config: providerConfig + `
					resource "segment_profiles_warehouse" "test" {
						space_id = "my-space-id"
						metadata_id = "my-metadata-id"
						name = "My warehouse name"
						enabled = true
						settings = jsonencode({
							"token": "my-token"
						})
						schema_name = "my-schema-name"
					}
				`,
				ImportState:   true,
				ImportStateId: "my-space-id:my-warehouse-id",
			},
			// Update and Read testing
			{
				Config: providerConfig + `
					resource "segment_profiles_warehouse" "test" {
						space_id = "my-space-id"
						metadata_id = "my-metadata-id"
						name = "My new warehouse name"
						enabled = false
						settings = jsonencode({
							"token": "my-other-token"
						})
						schema_name = "my-new-schema-name"
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("segment_profiles_warehouse.test", "id", "my-warehouse-id"),
					resource.TestCheckResourceAttr("segment_profiles_warehouse.test", "space_id", "my-space-id"),
					resource.TestCheckResourceAttr("segment_profiles_warehouse.test", "name", "My new warehouse name"),
					resource.TestCheckResourceAttr("segment_profiles_warehouse.test", "enabled", "false"),
					resource.TestCheckResourceAttr("segment_profiles_warehouse.test", "metadata_id", "my-metadata-id"),
					resource.TestCheckResourceAttr("segment_profiles_warehouse.test", "settings", "{\"token\":\"my-other-token\"}"),
					resource.TestCheckResourceAttr("segment_profiles_warehouse.test", "schema_name", "my-new-schema-name")),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
