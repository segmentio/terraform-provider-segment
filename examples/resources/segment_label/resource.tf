terraform {
  required_providers {
    segment = {
      source  = "registry.terraform.io/hashicorp/segment"
      version = "0.0.1"
    }
  }
}

resource "segment_label" "test" {
  key         = "environment"
  value       = "dev"
  description = "dev environment"
}

resource "segment_label" "import_label" {
  key         = "key"
  value       = "value"
  description = "description"
}