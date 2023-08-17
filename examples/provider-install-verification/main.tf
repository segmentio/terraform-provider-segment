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

resource "segment_destination" "example" {
  name      = "Dean's Terraform Destination"
  enabled   = true
  source_id = data.segment_source.example.id
  metadata = {
    id = "54521fd725e721e32a72eebf"
  }
}
