# Configures a specific source
resource "segment_source" "my_source" {
  name    = "My destination"
  enabled = true
  metadata = {
    id = "abc123"
  }
}