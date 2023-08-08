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
  id = "0HK6vEPONy"
}

output "source" {
  value = data.segment_source.example
}

data "segment_warehouse_metadata" "example" {}

output "warehouse_metadata" {
  value = data.segment_warehouse_metadata.example
}

data "segment_source_metadata" "example" {}

output "source_metadata" {
  value = data.segment_source_metadata.example
}

data "segment_destination_metadata" "example" {}

output "destination_metadata" {
  value = data.segment_destination_metadata.example
}

data "segment_warehouse" "example" {}

output "warehouse" {
  value = data.segment_warehouse.example
}