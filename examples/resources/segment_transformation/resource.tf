# Configures a specific transformation
resource "segment_transformation" "example" {
  source_id      = segment_source.example.id
  name           = "My transformation name"
  enabled        = true
  if             = "event = 'Bad Event'"
  new_event_name = "Good Event"
  property_renames = [
    {
      old_name = "old-name"
      new_name = "new-name"
    }
  ]
  property_value_transformations = [
    {
      property_paths = ["properties.some-property", "context.some-property"],
      property_value = "some property value"
    },
  ]
  fql_defined_properties = []
}
