# Configures a specific tracking plan
resource "segment_tracking_plan" "my_tracking_plan" {
  name        = "my-tracking-plan"
  type        = "LIVE"
  description = "My Tracking Plan Description"
  rules = [
    {
      key     = "Add Rule"
      type    = "TRACK"
      version = 1
      json_schema = jsonencode({
        "properties" : {
          "context" : {},
          "traits" : {},
          "properties" : {}
        }
      })
    }
  ]
}
