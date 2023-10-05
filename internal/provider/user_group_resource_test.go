package provider

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccUserGroupResource(t *testing.T) {
	t.Parallel()

	updated := 0
	updatedPermissions := 0
	updatedMembers := 0
	fakeServer := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("content-type", "application/json")
			var payload string

			// Rules requests
			if req.URL.Path == "/groups" && req.Method == http.MethodPost {
				payload = `
				{
					"data": {
						"userGroup": {
							"id": "my-group-id",
							"name": "my group name",
							"memberCount": 0
						}
					}
				}`
			} else if req.URL.Path == "/groups/my-group-id/permissions" && req.Method == http.MethodPut {
				updatedPermissions++
				payload = `
				{
					"data": {
						"permissions": [
							{
								"roleId": "my-other-role-id",
								"roleName": "My Other Role",
								"resources": [
									{
										"id": "my-workspace-id",
										"type": "WORKSPACE",
										"labels": [
											{
												"key": "my-label-key",
												"value": "my-label-value"
											}
										]
									}
								]
							}
						]
					}
				}`
			} else if req.URL.Path == "/group/my-group-id/users" && req.Method == http.MethodPut {
				updatedMembers++
				payload = `
				{
					"data": {
						"userGroup": {
							"id": "my-group-id",
							"name": "my group name",
							"memberCount": 1
						}
					}
				}`
			} else if req.URL.Path == "/groups/my-group-id/users" && req.Method == http.MethodGet {
				payload = `
				{
					"data": {
						"users": [
							{
								"id": "my-user-id",
								"email": "test@segment.com",
								"name": "test"
							}
						],
						"pagination": {
							"current": "MA==",
							"totalEntries": 1
						}
					}
				}`
			} else if req.URL.Path == "/groups/my-group-id/invites" && req.Method == http.MethodGet {
				if updatedMembers <= 1 {
					payload = `
					{
						"data": {
							"emails": [],
							"pagination": {
								"current": "MA==",
								"totalEntries": 0
							}
						}
					}`
				} else {
					payload = `
					{
						"data": {
							"emails": ["test-invite@segment.com"],
							"pagination": {
								"current": "MA==",
								"totalEntries": 0
							}
						}
					}`
				}

			} else if req.URL.Path == "/groups/my-group-id" && req.Method == http.MethodGet {
				if updated == 0 && updatedPermissions <= 1 {
					payload = `
					{
						"data": {
							"userGroup": {
								"id": "my-group-id",
								"name": "my group name",
								"memberCount": 0,
								"permissions": [
									{
										"roleId": "my-role-id",
										"roleName": "my-role-name",
										"resources": [
											{
												"id": "my-workspace-id",
												"type": "WORKSPACE",
												"labels": []
											}
										]
									}
								]
							}
						}
					}`
				} else {
					payload = `
					{
						"data": {
							"userGroup": {
								"id": "my-group-id",
								"name": "my new group name",
								"memberCount": 0,
								"permissions": [
									{
										"roleId": "my-other-role-id",
										"roleName": "My Other Role",
										"resources": [
											{
												"id": "my-workspace-id",
												"type": "WORKSPACE",
												"labels": [
													{
														"key": "my-label-key",
														"value": "my-label-value"
													}
												]
											}
										]
									}
								]
							}
						}
					}`
				}

			} else if req.URL.Path == "/groups/my-group-id" && req.Method == http.MethodDelete {
				payload = `
				{
					"data": {
						"status": "SUCCESS"
					}
				}`
			} else if req.URL.Path == "/groups/my-group-id" && req.Method == http.MethodPatch {
				updated++
				payload = `
				{
					"data": {
						"userGroup": {
							"id": "my-group-id",
							"name": "my new group name",
							"memberCount": 0
						}
					}
				}`
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
					resource "segment_user_group" "test" {
						name = "my group name"
						permissions = [
						  {
							role_id = "my-role-id"
							resources = [
							  {
								id     = "my-workspace-id"
								type   = "WORKSPACE"
								labels = []
							  }
							]
						  }
						]
						members = ["test@segment.com"]
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("segment_user_group.test", "id", "my-group-id"),
					resource.TestCheckResourceAttr("segment_user_group.test", "name", "my group name"),
					resource.TestCheckResourceAttr("segment_user_group.test", "permissions.#", "1"),
					resource.TestCheckResourceAttr("segment_user_group.test", "permissions.0.role_id", "my-role-id"),
					resource.TestCheckResourceAttr("segment_user_group.test", "permissions.0.resources.#", "1"),
					resource.TestCheckResourceAttr("segment_user_group.test", "permissions.0.resources.0.id", "my-workspace-id"),
					resource.TestCheckResourceAttr("segment_user_group.test", "permissions.0.resources.0.type", "WORKSPACE"),
					resource.TestCheckResourceAttr("segment_user_group.test", "permissions.0.resources.0.labels.#", "0"),
				),
			},
			// ImportState testing
			{
				ResourceName: "segment_user_group.test",
				Config: providerConfig + `
					resource "segment_user_group" "test" {
						name = "my group name"
						permissions = [
						  {
							role_id = "my-role-id"
							resources = [
							  {
								id     = "my-workspace-id"
								type   = "WORKSPACE"
								labels = []
							  }
							]
						  }
						]
						members = ["test@segment.com"]
					}
				`,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: providerConfig + `
					resource "segment_user_group" "test" {
						name = "my new group name"
						permissions = [
						  {
							role_id = "my-other-role-id"
							resources = [
							  {
								id     = "my-workspace-id"
								type   = "WORKSPACE"
								labels = [
									{
										key = "my-label-key"
										value = "my-label-value"
									}
								]
							  }
							]
						  }
						]
						members = ["test@segment.com", "test-invite@segment.com"]
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("segment_user_group.test", "id", "my-group-id"),
					resource.TestCheckResourceAttr("segment_user_group.test", "name", "my new group name"),
					resource.TestCheckResourceAttr("segment_user_group.test", "permissions.#", "1"),
					resource.TestCheckResourceAttr("segment_user_group.test", "permissions.0.role_id", "my-other-role-id"),
					resource.TestCheckResourceAttr("segment_user_group.test", "permissions.0.resources.#", "1"),
					resource.TestCheckResourceAttr("segment_user_group.test", "permissions.0.resources.0.id", "my-workspace-id"),
					resource.TestCheckResourceAttr("segment_user_group.test", "permissions.0.resources.0.type", "WORKSPACE"),
					resource.TestCheckResourceAttr("segment_user_group.test", "permissions.0.resources.0.labels.#", "1"),
					resource.TestCheckResourceAttr("segment_user_group.test", "permissions.0.resources.0.labels.0.key", "my-label-key"),
					resource.TestCheckResourceAttr("segment_user_group.test", "permissions.0.resources.0.labels.0.value", "my-label-value"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
