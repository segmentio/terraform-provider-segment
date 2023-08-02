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

data "segment_destinationCatalog" "example" {}

output "destination_catalog" {
  value = data.segment_destinationCatalog.example
}

data "segment_sourceCatalog" "example" {}

output "source_catalog" {
  value = data.segment_sourceCatalog.example
}

data "segment_warehouseCatalog" "example" {}

output "warehouse_catalog" {
  value = data.segment_warehouseCatalog.example
}