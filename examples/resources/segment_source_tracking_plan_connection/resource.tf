# Configures a source tracking plan connection
resource "segment_source_tracking_plan_connection" "example" {
  source_id        = "abc123"
  tracking_plan_id = "xyz321"
  schema_settings = {
    forwarding_blocked_events_to = segment_source.my_source.id
    track = {
      allow_unplanned_events           = true
      allow_unplanned_event_properties = true
      allow_event_on_violations        = false
      allow_properties_on_violations   = true
      common_event_on_violations       = "ALLOW"
    }
    identify = {
      allow_traits_on_violations = false
      allow_unplanned_traits     = false
      common_event_on_violations = "ALLOW"
    }
    group = {
      allow_traits_on_violations = true
      allow_unplanned_traits     = true
      common_event_on_violations = "ALLOW"
    }
  }
}
