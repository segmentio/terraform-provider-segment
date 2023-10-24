# Configures a specific warehouse
resource "segment_source_tracking_plan_connection" "example" {
  source_id               = segment_source.javascript.id
  name                    = "Example Reverse ETL model"
  enabled                 = true
  description             = "Example Reverse ETL model"
  schedule_strategy       = "SPECIFIC_DAYS"
  query                   = "SELECT good_stuff FROM stuff"
  query_identifier_column = "good_stuff"
  schedule_config = jsonencode({
    "days" : [0, 1, 2, 3],
    "hours" : [0, 1, 3],
    "timezone" : "America/Los_Angeles"
  })
}
