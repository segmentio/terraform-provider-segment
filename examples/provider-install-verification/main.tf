terraform {
  required_providers {
    publicapi = {
      source = "registry.terraform.io/hashicorp/public-api"
      version = "0.0.1"
    }
  }
}

provider "publicapi" {
  token = "sgp_5onYxVMgX7tlhtE67V4dFifGYXIsMeppwVeW4xjo8b4qYd3AcegYYF42pNMLT5Xc"
}

data "publicapi_workspace" "example" {}

output "workspace" {
  value = data.publicapi_workspace.example
}
