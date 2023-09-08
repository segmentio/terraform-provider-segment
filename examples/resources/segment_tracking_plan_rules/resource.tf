# Configures a specific tracking plan
resource "segment_tracking_plan_rules" "my_tracking_plan" {
  tracking_plan_id = "abc123"
  rules = [
    {
      key     = "Add to Cart"
      type    = "TRACK"
      version = 1
      json_schema = jsonencode({
        "$schema" : "http://json-schema.org/draft-07/schema#",
        "type" : "object",
        "labels" : {},
        "description" : "addToCart",
        "properties" : {
          "context" : {},
          "traits" : {},
          "properties" : {
            "type" : "object",
            "properties" : {
              "user_id" : {
                "description" : "",
                "type" : "string"
              },
              "item_id" : {
                "description" : "",
                "type" : "string"
              }
            },
            "required" : [
              "user_id",
              "item_id"
            ]
          }
        }
      })
    }
  ]
}
