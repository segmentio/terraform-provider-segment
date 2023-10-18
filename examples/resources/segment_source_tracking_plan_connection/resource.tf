# Configures a source tracking plan connection
resource "segment_source_tracking_plan_connection" "example" {
  source_id        = "abc123"
  tracking_plan_id = "xyz321"
}
