package provider

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccWarehouseDataSource(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		fakeServer := httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("content-type", "application/json")
				_, _ = w.Write([]byte(`
				{
  "data": {
    "warehouse": {
      "id": "warehouse-id",
      "workspaceId": "workspace-id",
      "enabled": false,
      "metadata": {
        "id": "my-warehouse-metadata-id",
        "slug": "my-warehouse-metadata-slug",
        "name": "The name of the warehouse metadata",
        "description": "The description of a warehouse metadata",
        "logos": {
          "default": "the default value of a logo",
          "mark": "the mark value of a logo",
          "alt": "the alt value of a logo"
        },
        "options": [
          {
            "name": "the option name",
            "required": true,
            "type": "the option type",
            "description": "the option description",
            "label": "the option label"
          }
        ]
      }
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
					Config: providerConfig + `data "segment_warehouse" "test" {}`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.segment_warehouse.test", "id", "warehouse-id"),
						resource.TestCheckResourceAttr("data.segment_warehouse.test", "workspace_id", "workspace-id"),
						resource.TestCheckResourceAttr("data.segment_warehouse.test", "enabled", "false"),
						resource.TestCheckResourceAttr("data.segment_warehouse.test", "metadata.id", "my-warehouse-metadata-id"),
						resource.TestCheckResourceAttr("data.segment_warehouse.test", "metadata.name", "The name of the warehouse metadata"),
						resource.TestCheckResourceAttr("data.segment_warehouse.test", "metadata.slug", "my-warehouse-metadata-slug"),
						resource.TestCheckResourceAttr("data.segment_warehouse.test", "metadata.description", "The description of a warehouse metadata"),
						resource.TestCheckResourceAttr("data.segment_warehouse.test", "metadata.logos.default", "the default value of a logo"),
						resource.TestCheckResourceAttr("data.segment_warehouse.test", "metadata.logos.mark", "the mark value of a logo"),
						resource.TestCheckResourceAttr("data.segment_warehouse.test", "metadata.logos.alt", "the alt value of a logo"),
						resource.TestCheckResourceAttr("data.segment_warehouse.test", "metadata.options.#", "1"),
						resource.TestCheckResourceAttr("data.segment_warehouse.test", "metadata.options.0.name", "the option name"),
						resource.TestCheckResourceAttr("data.segment_warehouse.test", "metadata.options.0.type", "the option type"),
						resource.TestCheckResourceAttr("data.segment_warehouse.test", "metadata.options.0.required", "true"),
						resource.TestCheckResourceAttr("data.segment_warehouse.test", "metadata.options.0.description", "the option description"),
						resource.TestCheckResourceAttr("data.segment_warehouse.test", "metadata.options.0.label", "the option label"),
					),
				},
			},
		})
	})
}
