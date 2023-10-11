# Configures a specific function
resource "segment_function" "example" {
  code          = "// Learn more about source functions API at https://segment.com/docs/connections/sources/source-functions"
  display_name  = "My source function"
  logo_url      = "https://placekitten.com/200/138"
  resource_type = "SOURCE"
  description   = "This is a source function."
  settings = [
    {
      name        = "apiKey"
      label       = "api key"
      type        = "STRING"
      description = "api key"
      required    = false
      sensitive   = false
    },
  ]
}
