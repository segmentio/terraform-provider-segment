package provider

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRoleDataSource(t *testing.T) {
	t.Parallel()
	regex, _ := regexp.Compile("Role not found")

	t.Run("happy path", func(t *testing.T) {
		t.Parallel()
		fakeServer := httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("content-type", "application/json")
				_, _ = w.Write([]byte(`
				{
					"data": {
						"roles": [
							{
								"id": "my-role-id-1",
								"name": "Source Read-only",
								"description": "Read-only access to assigned Source(s), Source settings, enabled Destinations, Schema, live data in the Debugger, and connected Tracking Plans."
							},
							{
								"id": "my-role-id-2",
								"name": "Workspace Owner",
								"description": "Owners have full read and edit access to everything in the workspace, including Sources, Destinations, add-on products, and settings. Owners have full edit access to all Team Permissions."
							},
						]
					}
				}
			`))
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
					Config: providerConfig + `data "segment_role" "test" { id = "my-role-id-1" }`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.segment_role.test", "id", "my-role-id-1"),
						resource.TestCheckResourceAttr("data.segment_role.test", "name", "Source Read-only"),
						resource.TestCheckResourceAttr("data.segment_role.test", "slug", "Read-only access to assigned Source(s), Source settings, enabled Destinations, Schema, live data in the Debugger, and connected Tracking Plans."),
					),
				},
			},
		})
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()
		fakeServer := httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("content-type", "application/json")
				_, _ = w.Write([]byte(`
				{
					"data": {
						"roles": [
							{
								"id": "my-role-id-1",
								"name": "Source Read-only",
								"description": "Read-only access to assigned Source(s), Source settings, enabled Destinations, Schema, live data in the Debugger, and connected Tracking Plans."
							},
							{
								"id": "my-role-id-2",
								"name": "Workspace Owner",
								"description": "Owners have full read and edit access to everything in the workspace, including Sources, Destinations, add-on products, and settings. Owners have full edit access to all Team Permissions."
							},
						]
					}
				}
			`))
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
					Config:      providerConfig + `data "segment_role" "test" { id = "my-role-id-3" }`,
					ExpectError: regex,
				},
			},
		})
	})
}
