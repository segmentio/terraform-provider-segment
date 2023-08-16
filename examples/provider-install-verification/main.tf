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
