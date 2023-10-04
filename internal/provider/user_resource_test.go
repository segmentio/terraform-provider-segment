package provider

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccUserResource(t *testing.T) {
	t.Parallel()

	updated := 0
	inviteQueried := 0
	fakeServer := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("content-type", "application/json")
			var payload string

			if req.URL.Path == "/invites" && req.Method == http.MethodPost {
				payload = `
				{
					"data": {
						"emails": [
							"test@twilio.com"
						]
					}
				}
				`
			} else if req.URL.Path == "/invites" && req.Method == http.MethodDelete {
				payload = `
				{
					"data": {
						"status": "SUCCESS"
					}
				}
				`
			} else if req.URL.Path == "/users" && req.Method == http.MethodGet {
				if inviteQueried < 6 {
					payload = `
					{
						"data": {
							"users": [
								{
									"id": "my-other-user-id",
									"name": "My other user",
									"email": "other@segment.com"
								}
							]
						}
					}
					`
					inviteQueried++
				} else {
					payload = `
					{
						"data": {
							"users": [
								{
									"id": "my-other-user-id",
									"name": "My other user",
									"email": "other@segment.com"
								},
								{
									"id": "my-user-id",
									"name": "My user",
									"email": "test@segment.com"
								}
							]
						}
					}
					`
				}
			} else if req.URL.Path == "/users/my-user-id" && req.Method == http.MethodGet {
				if updated == 0 {
					payload = `
					{
						"data": {
							"user": {
								"id": "my-user-id",
								"name": "My user",
								"email": "test@segment.com",
								"permissions": [
									{
										"roleId": "my-role-id",
										"roleName": "My Role Name",
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
					}
				`
				} else {
					payload = `
					{
						"data": {
							"user": {
								"id": "my-user-id",
								"name": "My user",
								"email": "test@segment.com",
								"permissions": [
									{
										"roleId": "my-other-role-id",
										"roleName": "My Other Role Name",
										"resources": [
											{
												"id": "my-workspace-id",
												"type": "WORKSPACE",
												"labels": [
													{
														"key": "my-new-label-key",
														"value": "my-new-label-value"
													}
												]
											}
										]
									}
								]
							}
						}
					}
				`
				}
			} else if req.URL.Path == "/users/my-user-id/permissions" && req.Method == http.MethodPut {
				updated++
				payload = `
				{
					"data": {
						"permissions": [
							{
								"roleId": "my-other-role-id",
								"roleName": "My Role Name",
								"resources": [
									{
										"id": "my-workspace-id",
										"type": "WORKSPACE",
										"labels": [
											{
												"key": "my-new-label-key",
												"value": "my-new-label-value"
											}
										]
									}
								]
							}
						]
					}
				}
				`
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
			// Invite create testing
			{
				Config: providerConfig + `
					resource "segment_user" "test" {
						email = "old_test@segment.com"
						permissions = [
							{
								role_id = "my-role-id"
								resources = [
									{
										id     = "my-workspace-id"
										type   = "WORKSPACE"
										labels = [
											{
												key   = "my-label-key"
												value = "my-label-value"
											}
										]
									}
								]
							}
						]
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("segment_user.test", "id", "old_test@segment.com"),
					resource.TestCheckResourceAttr("segment_user.test", "name", "old_test@segment.com"),
					resource.TestCheckResourceAttr("segment_user.test", "email", "old_test@segment.com"),
					resource.TestCheckResourceAttr("segment_user.test", "is_invite", "true"),
					resource.TestCheckResourceAttr("segment_user.test", "permissions.#", "1"),
					resource.TestCheckResourceAttr("segment_user.test", "permissions.0.role_id", "my-role-id"),
					resource.TestCheckResourceAttr("segment_user.test", "permissions.0.resources.#", "1"),
					resource.TestCheckResourceAttr("segment_user.test", "permissions.0.resources.0.id", "my-workspace-id"),
					resource.TestCheckResourceAttr("segment_user.test", "permissions.0.resources.0.type", "WORKSPACE"),
					resource.TestCheckResourceAttr("segment_user.test", "permissions.0.resources.0.labels.#", "1"),
					resource.TestCheckResourceAttr("segment_user.test", "permissions.0.resources.0.labels.0.key", "my-label-key"),
					resource.TestCheckResourceAttr("segment_user.test", "permissions.0.resources.0.labels.0.value", "my-label-value"),
				),
			},
			// Invite replace testing
			{
				Config: providerConfig + `
					resource "segment_user" "test" {
						email = "test@segment.com"
						permissions = [
							{
								role_id = "my-role-id"
								resources = [
									{
										id     = "my-workspace-id"
										type   = "WORKSPACE"
										labels = [
											{
												key   = "my-label-key"
												value = "my-label-value"
											}
										]
									}
								]
							}
						]
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("segment_user.test", "id", "test@segment.com"),
					resource.TestCheckResourceAttr("segment_user.test", "name", "test@segment.com"),
					resource.TestCheckResourceAttr("segment_user.test", "email", "test@segment.com"),
					resource.TestCheckResourceAttr("segment_user.test", "is_invite", "true"),
					resource.TestCheckResourceAttr("segment_user.test", "permissions.#", "1"),
					resource.TestCheckResourceAttr("segment_user.test", "permissions.0.role_id", "my-role-id"),
					resource.TestCheckResourceAttr("segment_user.test", "permissions.0.resources.#", "1"),
					resource.TestCheckResourceAttr("segment_user.test", "permissions.0.resources.0.id", "my-workspace-id"),
					resource.TestCheckResourceAttr("segment_user.test", "permissions.0.resources.0.type", "WORKSPACE"),
					resource.TestCheckResourceAttr("segment_user.test", "permissions.0.resources.0.labels.#", "1"),
					resource.TestCheckResourceAttr("segment_user.test", "permissions.0.resources.0.labels.0.key", "my-label-key"),
					resource.TestCheckResourceAttr("segment_user.test", "permissions.0.resources.0.labels.0.value", "my-label-value"),
				),
			},
			// User update testing
			{
				Config: providerConfig + `
					resource "segment_user" "test" {
						email = "test@segment.com"
						permissions = [
							{
								role_id = "my-other-role-id"
								resources = [
									{
										id     = "my-workspace-id"
										type   = "WORKSPACE"
										labels = [
											{
												key   = "my-new-label-key"
												value = "my-new-label-value"
											}
										]
									}
								]
							}
						]
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("segment_user.test", "id", "my-user-id"),
					resource.TestCheckResourceAttr("segment_user.test", "name", "My user"),
					resource.TestCheckResourceAttr("segment_user.test", "email", "test@segment.com"),
					resource.TestCheckResourceAttr("segment_user.test", "is_invite", "false"),
					resource.TestCheckResourceAttr("segment_user.test", "permissions.#", "1"),
					resource.TestCheckResourceAttr("segment_user.test", "permissions.0.role_id", "my-other-role-id"),
					resource.TestCheckResourceAttr("segment_user.test", "permissions.0.resources.#", "1"),
					resource.TestCheckResourceAttr("segment_user.test", "permissions.0.resources.0.id", "my-workspace-id"),
					resource.TestCheckResourceAttr("segment_user.test", "permissions.0.resources.0.type", "WORKSPACE"),
					resource.TestCheckResourceAttr("segment_user.test", "permissions.0.resources.0.labels.#", "1"),
					resource.TestCheckResourceAttr("segment_user.test", "permissions.0.resources.0.labels.0.key", "my-new-label-key"),
					resource.TestCheckResourceAttr("segment_user.test", "permissions.0.resources.0.labels.0.value", "my-new-label-value"),
				),
			},
			// Import testing
			{
				ResourceName: "segment_user.test",
				Config: providerConfig + `
					resource "segment_user" "test" {
						email = "test@segment.com"
						permissions = [
							{
								role_id = "my-other-role-id"
								resources = [
									{
										id     = "my-workspace-id"
										type   = "WORKSPACE"
										labels = [
											{
												key   = "my-new-label-key"
												value = "my-new-label-value"
											}
										]
									}
								]
							}
						]
					}
				`,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "my-user-id",
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
