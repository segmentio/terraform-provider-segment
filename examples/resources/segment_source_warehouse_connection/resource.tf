# Configures a source warehouse connection
resource "segment_warehouse" "example" {
  source_id    = "abc123"
  warehouse_id = "xyz321"
}
