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

data "segment_warehouseMetadata" "example" {}

output "warehouse_metadata" {
  value = data.segment_warehouseMetadata.example
}

data "segment_sourceMetadata" "example" {}

output "source_metadata" {
  value = data.segment_sourceMetadata.example
}

data "segment_destinationMetadata" "example" {}

output "destination_metadata" {
  value = data.segment_destinationMetadata.example
}