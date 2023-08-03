package provider

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccWorkspaceDataSource(t *testing.T) {
	fakeServer := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("content-type", "application/json")
			_, _ = w.Write([]byte(`
				{
					"data": {
						"workspace": {
							"id": "my-workspace-id",
							"name": "My workspace name",
							"slug": "my-workspace-slug"
						}
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
				Config: providerConfig + `data "segment_workspace" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.segment_workspace.test", "id", "my-workspace-id"),
					resource.TestCheckResourceAttr("data.segment_workspace.test", "name", "My workspace name"),
					resource.TestCheckResourceAttr("data.segment_workspace.test", "slug", "my-workspace-slug"),
				),
			},
		},
	})
}
