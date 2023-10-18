# Configures a source warehouse connection
resource "segment_source_warehouse_connection" "example" {
  source_id    = "abc123"
  warehouse_id = "xyz321"
}
