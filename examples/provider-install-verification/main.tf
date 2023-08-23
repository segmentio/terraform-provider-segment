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

resource "segment_source" "example2" {
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


output "destination" {
  value = data.segment_destination.example
}