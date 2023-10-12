# Configures a specific destination filter
resource "segment_destination_filter" "sample" {
  if             = "type = \"identify\""
  destination_id = "abc123"
  source_id      = "xyz321"
  title          = "Identify event sampling filter"
  enabled        = true
  description    = "Samples identify events at 5%"
  actions = [
    {
      type    = "SAMPLE"
      percent = 0.05
    },
  ]
}
