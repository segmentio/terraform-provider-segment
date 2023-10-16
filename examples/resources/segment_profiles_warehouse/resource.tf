# Configures a specific profiles sync warehouse
resource "segment_profiles_warehouse" "example" {
  space_id    = "cba321"
  metadata_id = "abc123"
  enabled     = true
  settings = jsonencode({
    token : "zyx321"
  })
  name = "My Terraform Warehouse!"
}