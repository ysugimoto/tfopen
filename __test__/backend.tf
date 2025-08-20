terraform {
  required_version = ">= 1.11.0"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 6.9.0"
    }
  }
  backend "remote" {
    hostname     = "app.terraform.io"
    organization = "dummy_org"

    workspaces {
      name = "dummy_workspace"
    }
  }
}
