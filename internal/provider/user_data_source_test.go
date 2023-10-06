package provider

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccUserDataSource(t *testing.T) {
	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		fakeServer := httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("content-type", "application/json")
				_, _ = w.Write([]byte(`
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
				}`))
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
				// Read testing
				{
					Config: providerConfig + `data "segment_user" "test" { id = "my-user-id" }`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.segment_user.test", "id", "my-user-id"),
						resource.TestCheckResourceAttr("data.segment_user.test", "name", "My user"),
						resource.TestCheckResourceAttr("data.segment_user.test", "email", "test@segment.com"),
						resource.TestCheckResourceAttr("data.segment_user.test", "permissions.#", "1"),
						resource.TestCheckResourceAttr("data.segment_user.test", "permissions.0.role_id", "my-role-id"),
						resource.TestCheckResourceAttr("data.segment_user.test", "permissions.0.resources.#", "1"),
						resource.TestCheckResourceAttr("data.segment_user.test", "permissions.0.resources.0.id", "my-workspace-id"),
						resource.TestCheckResourceAttr("data.segment_user.test", "permissions.0.resources.0.type", "WORKSPACE"),
						resource.TestCheckResourceAttr("data.segment_user.test", "permissions.0.resources.0.labels.#", "1"),
						resource.TestCheckResourceAttr("data.segment_user.test", "permissions.0.resources.0.labels.0.key", "my-label-key"),
						resource.TestCheckResourceAttr("data.segment_user.test", "permissions.0.resources.0.labels.0.value", "my-label-value"),
					),
				},
			},
		})
	})

	t.Run("nulls", func(t *testing.T) {
		fakeServer := httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("content-type", "application/json")
				_, _ = w.Write([]byte(`
				{
					"data": {
						"user": {
							"id": "my-user-id",
							"name": "My user",
							"email": "test@segment.com"
						}
					}
				}`))
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
				// Read testing
				{
					Config: providerConfig + `data "segment_user" "test" { id = "my-user-id" }`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.segment_user.test", "id", "my-user-id"),
						resource.TestCheckResourceAttr("data.segment_user.test", "name", "My user"),
						resource.TestCheckResourceAttr("data.segment_user.test", "email", "test@segment.com"),
						resource.TestCheckResourceAttr("data.segment_user.test", "permissions.#", "0"),
					),
				},
			},
		})
	})
}
