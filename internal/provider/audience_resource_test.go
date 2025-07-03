package provider

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	tfstate "github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccAudienceResource_basic(t *testing.T) {
	spaceID := os.Getenv("SEGMENT_TEST_SPACE_ID")
	if spaceID == "" {
		t.Skip("SEGMENT_TEST_SPACE_ID must be set for acceptance tests.")
	}
	resource.Test(t, resource.TestCase{
		// ProviderFactories: testAccProviderFactories, // Uncomment if you use provider factories
		Steps: []resource.TestStep{
			{
				Config: testAccAudienceResourceConfig(spaceID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("segment_audience.test", "id"),
					resource.TestCheckResourceAttr("segment_audience.test", "name", "Test Audience"),
				),
			},
		},
	})
}

func TestAccAudienceResource_update(t *testing.T) {
	spaceID := os.Getenv("SEGMENT_TEST_SPACE_ID")
	if spaceID == "" {
		t.Skip("SEGMENT_TEST_SPACE_ID must be set for acceptance tests.")
	}
	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: testAccAudienceResourceConfig(spaceID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("segment_audience.test", "id"),
					resource.TestCheckResourceAttr("segment_audience.test", "name", "Test Audience"),
				),
			},
			{
				Config: testAccAudienceResourceConfigUpdate(spaceID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("segment_audience.test", "name", "Updated Audience"),
					resource.TestCheckResourceAttr("segment_audience.test", "description", "Updated by Terraform acceptance test"),
				),
			},
		},
	})
}

func TestAccAudienceResource_import(t *testing.T) {
	spaceID := os.Getenv("SEGMENT_TEST_SPACE_ID")
	if spaceID == "" {
		t.Skip("SEGMENT_TEST_SPACE_ID must be set for acceptance tests.")
	}
	resourceName := "segment_audience.test"
	resource.Test(t, resource.TestCase{
		// ProviderFactories: testAccProviderFactories, // Uncomment if you use provider factories
		Steps: []resource.TestStep{
			{
				Config: testAccAudienceResourceConfig(spaceID),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *tfstate.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("not found: %s", resourceName)
					}
					spaceID := rs.Primary.Attributes["space_id"]
					id := rs.Primary.ID
					return fmt.Sprintf("%s:%s", spaceID, id), nil
				},
			},
		},
	})
}

func TestAccAudienceResource_delete(t *testing.T) {
	spaceID := os.Getenv("SEGMENT_TEST_SPACE_ID")
	if spaceID == "" {
		t.Skip("SEGMENT_TEST_SPACE_ID must be set for acceptance tests.")
	}
	resourceName := "segment_audience.test"
	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: testAccAudienceResourceConfig(spaceID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				ResourceName: resourceName,
				Config:       "",
				Check: func(s *tfstate.State) error {
					if _, ok := s.RootModule().Resources[resourceName]; ok {
						return fmt.Errorf("Resource %s still exists after delete", resourceName)
					}
					return nil
				},
				Destroy: true,
			},
		},
	})
}

func TestAccAudienceResource_error_invalidDefinition(t *testing.T) {
	spaceID := os.Getenv("SEGMENT_TEST_SPACE_ID")
	if spaceID == "" {
		t.Skip("SEGMENT_TEST_SPACE_ID must be set for acceptance tests.")
	}
	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config:      testAccAudienceResourceConfigInvalidDefinition(spaceID),
				ExpectError: regexp.MustCompile(`(?i)definition`),
			},
		},
	})
}

func TestAccAudienceResource_complexAttributes(t *testing.T) {
	spaceID := os.Getenv("SEGMENT_TEST_SPACE_ID")
	if spaceID == "" {
		t.Skip("SEGMENT_TEST_SPACE_ID must be set for acceptance tests.")
	}
	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: testAccAudienceResourceConfigComplex(spaceID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("segment_audience.test", "definition.query", "event('Shoes Bought').count() >= 1 && event('Shirt Bought').count() >= 2"),
					resource.TestCheckResourceAttr("segment_audience.test", "options.includeHistoricalData", "true"),
				),
			},
		},
	})
}

func testAccAudienceResourceConfig(spaceID string) string {
	return `
resource "segment_audience" "test" {
  space_id    = "` + spaceID + `"
  name        = "Test Audience"
  description = "Created by Terraform acceptance test"
  definition  = {
    query = "event('Shoes Bought').count() >= 1"
  }
}
`
}

func testAccAudienceResourceConfigUpdate(spaceID string) string {
	return `
resource "segment_audience" "test" {
  space_id    = "` + spaceID + `"
  name        = "Updated Audience"
  description = "Updated by Terraform acceptance test"
  definition  = {
    query = "event('Shoes Bought').count() >= 2"
  }
}
`
}

func testAccAudienceResourceConfigInvalidDefinition(spaceID string) string {
	return `
resource "segment_audience" "test" {
  space_id    = "` + spaceID + `"
  name        = "Invalid Audience"
  description = "Should fail"
  definition  = {
    invalid_field = "not a valid query"
  }
}
`
}

func testAccAudienceResourceConfigComplex(spaceID string) string {
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
`
}
