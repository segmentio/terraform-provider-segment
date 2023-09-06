# Configures a specific tracking plan
resource "segment_tracking_plan" "my_tracking_plan" {
  name = "my-tracking-plan"
  type = "LIVE"
}