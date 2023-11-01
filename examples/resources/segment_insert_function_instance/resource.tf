# Configures a specific warehouse
resource "segment_insert_function_instance" "example" {
  integration_id = segment_source.my_source.id
  function_id    = segment_function.my_function.id
  name           = "My insert function instance"
  enabled        = true
  settings = jsonencode({
    "apiKey" : "abc123"
  })
}
