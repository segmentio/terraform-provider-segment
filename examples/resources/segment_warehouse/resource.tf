# Configures a specific warehouse
resource "segment_warehouse" "example" {
  metadata = {
    id = "abc123"
  }
  enabled = true
  settings = jsonencode({
    token : "zyx321"
  })
  name = "My Terraform Warehouse!"
}
