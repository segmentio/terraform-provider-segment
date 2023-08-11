terraform {
  required_providers {
    segment = {
      source  = "registry.terraform.io/hashicorp/segment"
      version = "0.0.1"
    }
  }
}

resource "segment_label" "test" {
  label = {
    key         = "environment"
    value       = "dev"
    description = "dev environment"
  }
}

output "dev_label" {
  value = segment_label.test
}