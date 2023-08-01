terraform {
  required_providers {
    publicapi = {
      source  = "registry.terraform.io/hashicorp/public-api"
      version = "0.0.1"
    }
  }
}

provider "publicapi" {}

data "publicapi_workspace" "example" {}

output "workspace" {
  value = data.publicapi_workspace.example
}
