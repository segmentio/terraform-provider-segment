provider "segment" {
    url = "https://api.segmentapis.build"
}

resource "segment_source" "javascript" {
  metadata = {
    id = "IqDTy1TpoU"
  }
  enabled = true
  slug    = "terraform-javascript-source"
  settings = jsonencode({
    "cool_setting" : "wow"
  })
}

resource "segment_destination" "help_scout" {
  metadata = {
    id = "54521fd725e721e32a72eebf"
  }
  enabled   = true
  source_id = segment_source.javascript.id
  settings = jsonencode({
    "really_cool_setting" : "awesome"
  })
  name = "terraform-help-scout-destination"
}

resource "segment_warehouse" "redshift" {
  metadata = {
    id = "aea3c55dsz"
  }
  settings = jsonencode({ "database" : "qa", "hostname" : "warehouses-qa-redshift.cxz7mja8ukuc.us-west-2.redshift.amazonaws.com", "password" : "0WDc1ky5YEpbNsr8N8DM", "port" : "5439", "username" : "segment", "dean" : "dean" })
  name     = "terraform-redshift-warehouse"
  enabled  = true
}

resource "segment_source_warehouse_connection" "example" {
  source_id    = segment_source.javascript.id
  warehouse_id = segment_warehouse.redshift.id
}

output "source" {
    value = segment_source.javascript
}
