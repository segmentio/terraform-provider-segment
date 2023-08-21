package provider

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccLabelsDataSource(t *testing.T) {
	t.Parallel()

	fakeServer := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("content-type", "application/json")
			_, _ = w.Write([]byte(`
				{
  "data": {
    "labels": [
      {
        "key": "environment",
        "value": "dev",
        "description": "dev environment"
      },
      {
        "key": "environment",
        "value": "prod",
        "description": "prod environment"
      }
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
				Config: providerConfig + `data "segment_labels" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.segment_labels.test", "id", "placeholder"),
					resource.TestCheckResourceAttr("data.segment_labels.test", "labels.#", "2"),
					resource.TestCheckResourceAttr("data.segment_labels.test", "labels.0.key", "environment"),
					resource.TestCheckResourceAttr("data.segment_labels.test", "labels.0.value", "dev"),
					resource.TestCheckResourceAttr("data.segment_labels.test", "labels.0.description", "dev environment"),
					resource.TestCheckResourceAttr("data.segment_labels.test", "labels.1.key", "environment"),
					resource.TestCheckResourceAttr("data.segment_labels.test", "labels.1.value", "prod"),
					resource.TestCheckResourceAttr("data.segment_labels.test", "labels.1.description", "prod environment"),
				),
			},
		},
	})
}
