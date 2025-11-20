package provider

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAudienceDataSource_basic(t *testing.T) {
	spaceID := os.Getenv("SEGMENT_TEST_SPACE_ID")
	if spaceID == "" {
		t.Skip("SEGMENT_TEST_SPACE_ID must be set for acceptance tests.")
	}
	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: testAccAudienceDataSourceConfig(spaceID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.segment_audience.test", "id"),
					resource.TestCheckResourceAttr("data.segment_audience.test", "name", "Test Audience"),
				),
			},
		},
	})
}

func TestAccAudienceDataSource_notFound(t *testing.T) {
	spaceID := os.Getenv("SEGMENT_TEST_SPACE_ID")
	if spaceID == "" {
		t.Skip("SEGMENT_TEST_SPACE_ID must be set for acceptance tests.")
	}
	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config:      testAccAudienceDataSourceConfigNotFound(spaceID),
				ExpectError: regexp.MustCompile(`(?i)Unable to read Audience`),
			},
		},
	})
}

func TestAccAudienceDataSource_complexAttributes(t *testing.T) {
	spaceID := os.Getenv("SEGMENT_TEST_SPACE_ID")
	if spaceID == "" {
		t.Skip("SEGMENT_TEST_SPACE_ID must be set for acceptance tests.")
	}
	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: testAccAudienceDataSourceConfigComplex(spaceID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.segment_audience.test", "definition.query", "event('Shoes Bought').count() >= 1 && event('Shirt Bought').count() >= 2"),
					resource.TestCheckResourceAttr("data.segment_audience.test", "options.includeHistoricalData", "true"),
				),
			},
		},
	})
}

func testAccAudienceDataSourceConfig(spaceID string) string {
	return `
resource "segment_audience" "test" {
  space_id    = "` + spaceID + `"
  name        = "Test Audience"
  description = "Created by Terraform acceptance test"
  definition  = {
    query = "event('Shoes Bought').count() >= 1"
  }
}

data "segment_audience" "test" {
  space_id = segment_audience.test.space_id
  id       = segment_audience.test.id
}
`
}

func testAccAudienceDataSourceConfigNotFound(spaceID string) string {
	return `
data "segment_audience" "test" {
  space_id = "` + spaceID + `"
  id       = "audience_does_not_exist"
}
`
}

func testAccAudienceDataSourceConfigComplex(spaceID string) string {
	return `
resource "segment_audience" "test" {
  space_id    = "` + spaceID + `"
  name        = "Complex Audience"
  description = "Testing complex attributes"
  definition  = {
    query = "event('Shoes Bought').count() >= 1 && event('Shirt Bought').count() >= 2"
  }
  options = {
    includeHistoricalData = true
  }
}

data "segment_audience" "test" {
  space_id = segment_audience.test.space_id
  id       = segment_audience.test.id
}
`
}
