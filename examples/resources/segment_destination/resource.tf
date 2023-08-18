# Configures a specific destination
resource "segment_destination" "my_destination" {
  name      = "My Destination"
  enabled   = true
  source_id = "s123"
  metadata = {
    id = "dm123"
  }
}