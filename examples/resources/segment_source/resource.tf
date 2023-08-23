# Configures a specific source
resource "segment_source" "my_source" {
  slug    = "my_source_slug"
  name    = "My Source"
  enabled = true
  metadata = {
    id = "abc123"
  }
  settings = jsonencode({
    "token": "xyz321",
  })
}
