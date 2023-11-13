# Configures a specific source
resource "segment_source" "my_source" {
  slug    = "my_source_slug"
  name    = "My Source"
  enabled = true
  metadata = {
    id = "abc123"
  }
  settings = jsonencode({
    "token" : "xyz321",
  })
  labels = [
    {
      key   = "env"
      value = "dev"
    },
  ]
}

resource "segment_source" "my_source_with_schema_settings" {
  slug    = "my_source_slug"
  name    = "My Source"
  enabled = true
  metadata = {
    id = "abc123"
  }
  settings = jsonencode({
    "token" : "xyz321",
  })
  labels = [
    {
      key   = "env"
      value = "dev"
    },
  ]
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
