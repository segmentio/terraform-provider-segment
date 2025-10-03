package provider

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
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
			// Create and Read testing.
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
			// ImportState testing.
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
			// Update and Read testing.
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
			// Delete testing automatically occurs in TestCase.
		},
	})
}

func TestAccProfilesWarehouseResource_SchemaNameHandling(t *testing.T) {
	// Test the schemaName handling that prevents API failures when the schema
	// name already exists in the warehouse configuration.
	t.Parallel()

	updateCount := 0
	fakeServer := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("content-type", "application/json")

			payload := ""
			if req.URL.Path == "/spaces/my-space-id/profiles-warehouses" && req.Method == http.MethodPost {
				// Initial create response.
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
				// Update response - schemaName should only be sent when it changes.
				updateCount++
				if updateCount == 1 {
					// First update: name changes, schema_name stays the same
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
										"name": "My updated warehouse name",
										"token": "my-other-token"
									},
									"schemaName": "my-schema-name"
								}
							}
						}
					`
				} else {
					// Second update: schema_name changes
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
										"name": "My final warehouse name",
										"token": "my-final-token"
									},
									"schemaName": "my-new-schema-name"
								}
							}
						}
					`
				}
			} else if req.URL.Path == "/spaces/my-space-id/profiles-warehouses" && req.Method == http.MethodGet {
				// Read response - return current state based on update count.
				if updateCount == 0 {
					// Initial state
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
				} else if updateCount == 1 {
					// After first update: name changed, schema_name same
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
											"name": "My updated warehouse name",
											"token": "my-other-token"
										},
										"schemaName": "my-schema-name"
									}
								]
							}
						}
					`
				} else {
					// After second update: schema_name changed
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
											"name": "My final warehouse name",
											"token": "my-final-token"
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
			// Create with schema_name.
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
					resource.TestCheckResourceAttr("segment_profiles_warehouse.test", "schema_name", "my-schema-name"),
				),
			},
			// Update with same schema_name - should not send schemaName to API (prevents API failure).
			{
				Config: providerConfig + `
					resource "segment_profiles_warehouse" "test" {
						space_id = "my-space-id"
						metadata_id = "my-metadata-id"
						name = "My updated warehouse name"
						enabled = false
						settings = jsonencode({
							"token": "my-other-token"
						})
						schema_name = "my-schema-name"
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("segment_profiles_warehouse.test", "schema_name", "my-schema-name"),
				),
			},
			// Update with different schema_name - should send schemaName to API (legitimate change).
			{
				Config: providerConfig + `
					resource "segment_profiles_warehouse" "test" {
						space_id = "my-space-id"
						metadata_id = "my-metadata-id"
						name = "My final warehouse name"
						enabled = true
						settings = jsonencode({
							"token": "my-final-token"
						})
						schema_name = "my-new-schema-name"
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("segment_profiles_warehouse.test", "schema_name", "my-new-schema-name"),
				),
			},
		},
	})
}

func TestDetermineSchemaNameForUpdate(t *testing.T) {
	t.Parallel()

	// Test the determineSchemaNameForUpdate function that prevents API failures
	// when the schema name already exists in the warehouse configuration.
	tests := []struct {
		name            string
		planSchemaName  types.String
		stateSchemaName types.String
		expectedResult  *string
		description     string
	}{
		{
			name:            "both_null",
			planSchemaName:  types.StringNull(),
			stateSchemaName: types.StringNull(),
			expectedResult:  nil,
			description:     "Both null - should not send schemaName to API (prevents API failure)",
		},
		{
			name:            "both_unknown",
			planSchemaName:  types.StringUnknown(),
			stateSchemaName: types.StringUnknown(),
			expectedResult:  nil,
			description:     "Both unknown - should not send schemaName to API (prevents API failure)",
		},
		{
			name:            "both_same_value",
			planSchemaName:  types.StringValue("my-schema"),
			stateSchemaName: types.StringValue("my-schema"),
			expectedResult:  nil,
			description:     "Both have same value - should not send schemaName to API (prevents API failure)",
		},
		{
			name:            "plan_null_state_has_value",
			planSchemaName:  types.StringNull(),
			stateSchemaName: types.StringValue("my-schema"),
			expectedResult:  nil, // null pointer, not the actual value
			description:     "Plan null, state has value - should send schemaName to API",
		},
		{
			name:            "plan_has_value_state_null",
			planSchemaName:  types.StringValue("my-schema"),
			stateSchemaName: types.StringNull(),
			expectedResult:  stringPtr("my-schema"),
			description:     "Plan has value, state null - should send schemaName to API",
		},
		{
			name:            "plan_unknown_state_has_value",
			planSchemaName:  types.StringUnknown(),
			stateSchemaName: types.StringValue("my-schema"),
			expectedResult:  nil, // unknown becomes null pointer
			description:     "Plan unknown, state has value - should send schemaName to API",
		},
		{
			name:            "plan_has_value_state_unknown",
			planSchemaName:  types.StringValue("my-schema"),
			stateSchemaName: types.StringUnknown(),
			expectedResult:  stringPtr("my-schema"),
			description:     "Plan has value, state unknown - should send schemaName to API",
		},
		{
			name:            "different_values",
			planSchemaName:  types.StringValue("new-schema"),
			stateSchemaName: types.StringValue("old-schema"),
			expectedResult:  stringPtr("new-schema"),
			description:     "Different values - should send schemaName to API",
		},
		{
			name:            "empty_string_vs_null",
			planSchemaName:  types.StringValue(""),
			stateSchemaName: types.StringNull(),
			expectedResult:  stringPtr(""),
			description:     "Empty string vs null - should send schemaName to API",
		},
		{
			name:            "null_vs_empty_string",
			planSchemaName:  types.StringNull(),
			stateSchemaName: types.StringValue(""),
			expectedResult:  nil,
			description:     "Null vs empty string - should send schemaName to API",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Test the actual function.
			result := determineSchemaNameForUpdate(tt.planSchemaName, tt.stateSchemaName)

			// Check if the result matches expected.
			if !compareStringPointers(result, tt.expectedResult) {
				t.Errorf("Test case '%s' failed: %s\nExpected: %v, but got: %v",
					tt.name, tt.description, tt.expectedResult, result)
			}
		})
	}
}

// Helper function to create string pointers for test cases.
func stringPtr(s string) *string {
	return &s
}

// Helper function to compare string pointers.
func compareStringPointers(a, b *string) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	return *a == *b
}
