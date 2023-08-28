terraform {
  required_providers {
    segment = {
      source  = "registry.terraform.io/hashicorp/segment"
      version = "0.0.1"
    }
  }
}

provider "segment" {}

data "segment_workspace" "example" {}

output "workspace" {
  value = data.segment_workspace.example
}

data "segment_source" "example" {
  id = "dFwK97Hamo4NspDCEfCr9C"
}

output "source" {
  value = data.segment_source.example
}

resource "segment_destination" "example2" {
  name      = "Dean's Terraform Destination"
  enabled   = true
  source_id = data.segment_source.example.id
  metadata = {
    id = "54521fd725e721e32a72eebf"
  }
  settings = jsonencode({
    "mykey" : "myvalue",
    "mylist" : [1, 2, 3]
  })
}

data "segment_destination" "example" {
  id = segment_destination.example2.id
}

resource "segment_source" "example" {
  metadata = {
    id = "IqDTy1TpoU"
  }
  slug    = "dean-test-terraform"
  enabled = true
  settings = jsonencode({
    "mykey" : "myvalue",
    "mylist" : [1, 2, 3]
  })
}

resource "segment_warehouse" "example" {
  settings = jsonencode({
    username : "segment",
    password : "0WDc1ky5YEpbNsr8N8DM",
    database : "qa",
    hostname : "warehouses-qa-redshift.cxz7mja8ukuc.us-west-2.redshift.amazonaws.com",
    port : "5439"
  })
  name = "My Terraform Warehouse!"
  metadata = {
    id = "aea3c55dsz"
  }
  enabled = true
}

resource "segment_source_warehouse_connection" "example" {
  source_id    = data.segment_source.example.id
  warehouse_id = segment_warehouse.example.id
}
