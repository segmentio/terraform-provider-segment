# Configures a specific Reverse ETL model
resource "segment_reverse_etl_model" "example" {
  source_id               = segment_source.javascript.id
  name                    = "Example Reverse ETL model"
  enabled                 = true
  description             = "Example Reverse ETL model"
  query                   = "SELECT good_stuff FROM stuff"
  query_identifier_column = "good_stuff"
}
