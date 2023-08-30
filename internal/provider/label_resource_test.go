package provider

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccLabelResource(t *testing.T) {
	t.Parallel()

	postPayload := `
{
  "data": {
    "label": {
      "key": "environment",
      "value": "dev",
      "description": "dev environment"
    }
  }
}`

	getPayload := `
{
  "data": {
    "labels": [
      {
        "key": "environment",
        "value": "dev",
        "description": "dev environment"
      }
    ]
  }
}`

	fakeServer := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("content-type", "application/json")
			switch r.Method {
			case "POST":
				_, _ = w.Write([]byte(postPayload))
			case "GET":
				_, _ = w.Write([]byte(getPayload))
			}
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
resource "segment_label" "test" {
	key = "environment"
	value = "dev"
	description = "dev environment"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("segment_label.test", "id", "environment:dev"),
					resource.TestCheckResourceAttr("segment_label.test", "key", "environment"),
					resource.TestCheckResourceAttr("segment_label.test", "value", "dev"),
					resource.TestCheckResourceAttr("segment_label.test", "description", "dev environment"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "segment_label.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "environment,dev,dev environment",
			},
		},
	})
}
